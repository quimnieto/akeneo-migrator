package syncing

import "akeneo-migrator/kit/bus"

const SyncAttributeCommandType bus.Type = "attribute.sync"

// SyncAttributeCommand represents a command to sync an attribute
type SyncAttributeCommand struct {
	Code  string
	Debug bool
}

// Type returns the command type
func (c SyncAttributeCommand) Type() bus.Type {
	return SyncAttributeCommandType
}
