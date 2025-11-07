package syncing_since

import "akeneo-migrator/kit/bus"

const SyncProductsSinceCommandType bus.Type = "product.sync_updated"

// SyncProductsSinceCommand represents a command to sync products updated since a date
type SyncProductsSinceCommand struct {
	UpdatedSince string
	Debug        bool
}

// Type returns the command type
func (c SyncProductsSinceCommand) Type() bus.Type {
	return SyncProductsSinceCommandType
}
