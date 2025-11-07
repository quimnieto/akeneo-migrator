package syncing

import (
	"context"

	"akeneo-migrator/kit/bus"
)

// CommandHandler handles product sync commands
type CommandHandler struct {
	service *Service
}

// NewCommandHandler creates a new command handler
func NewCommandHandler(service *Service) *CommandHandler {
	return &CommandHandler{
		service: service,
	}
}

// Handle executes the sync command
func (h *CommandHandler) Handle(ctx context.Context, msg bus.Message) (bus.Response, error) {
	switch cmd := msg.(type) {
	case SyncProductCommand:
		result, err := h.service.Sync(ctx, cmd.Identifier)
		if err != nil {
			return bus.Response{Error: err}, err
		}
		return bus.Response{Data: result}, nil

	case SyncProductHierarchyCommand:
		result, err := h.service.SyncHierarchy(ctx, cmd.Identifier)
		if err != nil {
			return bus.Response{Error: err}, err
		}
		return bus.Response{Data: result}, nil

	default:
		return bus.Response{}, nil
	}
}
