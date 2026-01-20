package syncing

import "akeneo-migrator/kit/bus"

const SyncFamilyCommandType bus.Type = "family.sync"

// SyncFamilyCommand represents a command to sync a family
type SyncFamilyCommand struct {
	Code  string
	Debug bool
}

// Type returns the command type
func (c SyncFamilyCommand) Type() bus.Type {
	return SyncFamilyCommandType
}
