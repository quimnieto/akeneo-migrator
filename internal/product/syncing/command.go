package syncing

import "akeneo-migrator/kit/bus"

const SyncProductCommandType bus.Type = "product.sync"

// SyncProductCommand represents a command to sync a product hierarchy
type SyncProductCommand struct {
	Identifier string
	Debug      bool
}

// Type returns the command type
func (c SyncProductCommand) Type() bus.Type {
	return SyncProductCommandType
}
