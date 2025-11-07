package syncing

import "akeneo-migrator/kit/bus"

const (
	SyncProductCommandType          bus.Type = "product.sync"
	SyncProductHierarchyCommandType bus.Type = "product.sync_hierarchy"
)

// SyncProductCommand represents a command to sync a single product
type SyncProductCommand struct {
	Identifier string
	Debug      bool
}

// Type returns the command type
func (c SyncProductCommand) Type() bus.Type {
	return SyncProductCommandType
}

// SyncProductHierarchyCommand represents a command to sync a product hierarchy
type SyncProductHierarchyCommand struct {
	Identifier string
	Debug      bool
}

// Type returns the command type
func (c SyncProductHierarchyCommand) Type() bus.Type {
	return SyncProductHierarchyCommandType
}
