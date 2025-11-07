package bus

import "context"

// Type represents the type of a message
type Type string

// Message represents a command or query
type Message interface {
	Type() Type
}

// Response represents the response from a handler
type Response struct {
	Data  interface{}
	Error error
}

// Handler handles a message
type Handler interface {
	Handle(ctx context.Context, msg Message) (Response, error)
}

// Bus dispatches messages to handlers
type Bus interface {
	Dispatch(ctx context.Context, msg Message) (Response, error)
	Register(msgType Type, handler Handler)
}
