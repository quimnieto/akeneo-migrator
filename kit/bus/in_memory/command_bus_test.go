package inmemory

import (
	"context"
	"errors"
	"testing"

	"akeneo-migrator/kit/bus"
)

// Mock command
type mockCommand struct {
	value string
}

func (m mockCommand) Type() bus.Type {
	return "mock.command"
}

// Mock handler
type mockHandler struct {
	handleFunc func(ctx context.Context, msg bus.Message) (bus.Response, error)
}

func (m *mockHandler) Handle(ctx context.Context, msg bus.Message) (bus.Response, error) {
	if m.handleFunc != nil {
		return m.handleFunc(ctx, msg)
	}
	return bus.Response{Data: "handled"}, nil
}

func TestCommandBus_Dispatch_Success(t *testing.T) {
	handler := &mockHandler{
		handleFunc: func(ctx context.Context, msg bus.Message) (bus.Response, error) {
			cmd := msg.(mockCommand)
			return bus.Response{Data: "handled: " + cmd.value}, nil
		},
	}

	commandBus := NewCommandBus()
	commandBus.Register("mock.command", handler)

	response, err := commandBus.Dispatch(context.Background(), mockCommand{value: "test"})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Data != "handled: test" {
		t.Errorf("Expected 'handled: test', got %v", response.Data)
	}
}

func TestCommandBus_Dispatch_HandlerNotFound(t *testing.T) {
	commandBus := NewCommandBus()

	response, err := commandBus.Dispatch(context.Background(), mockCommand{value: "test"})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Data != nil {
		t.Errorf("Expected nil data, got %v", response.Data)
	}
}

func TestCommandBus_Dispatch_HandlerError(t *testing.T) {
	expectedError := errors.New("handler error")
	handler := &mockHandler{
		handleFunc: func(ctx context.Context, msg bus.Message) (bus.Response, error) {
			return bus.Response{}, expectedError
		},
	}

	commandBus := NewCommandBus()
	commandBus.Register("mock.command", handler)

	_, err := commandBus.Dispatch(context.Background(), mockCommand{value: "test"})

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}
}

func TestCommandBus_WithMiddleware(t *testing.T) {
	middlewareCalled := false

	middleware := func(ctx context.Context, msg bus.Message, next NextFunc) (bus.Response, error) {
		middlewareCalled = true
		return next(ctx, msg)
	}

	handler := &mockHandler{}

	commandBus := NewCommandBus(middleware)
	commandBus.Register("mock.command", handler)

	_, err := commandBus.Dispatch(context.Background(), mockCommand{value: "test"})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !middlewareCalled {
		t.Error("Expected middleware to be called")
	}
}

func TestCommandBus_MultipleMiddlewares(t *testing.T) {
	var executionOrder []string

	middleware1 := func(ctx context.Context, msg bus.Message, next NextFunc) (bus.Response, error) {
		executionOrder = append(executionOrder, "middleware1-before")
		response, err := next(ctx, msg)
		executionOrder = append(executionOrder, "middleware1-after")
		return response, err
	}

	middleware2 := func(ctx context.Context, msg bus.Message, next NextFunc) (bus.Response, error) {
		executionOrder = append(executionOrder, "middleware2-before")
		response, err := next(ctx, msg)
		executionOrder = append(executionOrder, "middleware2-after")
		return response, err
	}

	handler := &mockHandler{
		handleFunc: func(ctx context.Context, msg bus.Message) (bus.Response, error) {
			executionOrder = append(executionOrder, "handler")
			return bus.Response{}, nil
		},
	}

	commandBus := NewCommandBus(middleware1, middleware2)
	commandBus.Register("mock.command", handler)

	_, err := commandBus.Dispatch(context.Background(), mockCommand{value: "test"})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := []string{
		"middleware1-before",
		"middleware2-before",
		"handler",
		"middleware2-after",
		"middleware1-after",
	}

	if len(executionOrder) != len(expected) {
		t.Fatalf("Expected %d executions, got %d", len(expected), len(executionOrder))
	}

	for i, exp := range expected {
		if executionOrder[i] != exp {
			t.Errorf("Expected execution[%d] to be %s, got %s", i, exp, executionOrder[i])
		}
	}
}
