package inmemory

import (
	"context"

	"akeneo-migrator/kit/bus"
)

type Middleware func(ctx context.Context, msg bus.Message, next NextFunc) (bus.Response, error)

type NextFunc func(ctx context.Context, msg bus.Message) (bus.Response, error)

type CommandBus struct {
	handlers    map[bus.Type]bus.Handler
	middlewares []Middleware
}

func NewCommandBus(middlewares ...Middleware) *CommandBus {
	return &CommandBus{
		handlers:    make(map[bus.Type]bus.Handler),
		middlewares: middlewares,
	}
}

func (b *CommandBus) Dispatch(ctx context.Context, cmd bus.Message) (bus.Response, error) {
	handler, ok := b.handlers[cmd.Type()]
	if !ok {
		return bus.Response{}, nil
	}

	finalHandler := func(ctx context.Context, m bus.Message) (bus.Response, error) {
		return handler.Handle(ctx, m)
	}

	chain := finalHandler
	for i := len(b.middlewares) - 1; i >= 0; i-- {
		m := b.middlewares[i]
		next := chain
		chain = func(ctx context.Context, m_ bus.Message) (bus.Response, error) {
			return m(ctx, m_, next)
		}
	}

	return chain(ctx, cmd)
}

// Register implements the bus.Bus interface
func (b *CommandBus) Register(cmdType bus.Type, handler bus.Handler) {
	b.handlers[cmdType] = handler
}
