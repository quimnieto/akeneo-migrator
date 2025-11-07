package syncing

import "akeneo-migrator/kit/bus"

const SyncReferenceEntityCommandType bus.Type = "reference_entity.sync"

// SyncReferenceEntityCommand represents a command to sync a reference entity
type SyncReferenceEntityCommand struct {
	EntityName string
	Debug      bool
}

// Type returns the command type
func (c SyncReferenceEntityCommand) Type() bus.Type {
	return SyncReferenceEntityCommandType
}
