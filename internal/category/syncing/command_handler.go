package syncing

import (
	"context"

	"akeneo-migrator/kit/bus"
)

// CommandHandler handles SyncCategoryCommand
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
	cmd, ok := msg.(SyncCategoryCommand)
	if !ok {
		return bus.Response{}, nil
	}

	result, err := h.service.Sync(ctx, cmd.Code)
	if err != nil {
		return bus.Response{Error: err}, err
	}

	return bus.Response{Data: result}, nil
}
