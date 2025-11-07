package syncing

import "akeneo-migrator/kit/bus"

const SyncCategoryCommandType bus.Type = "category.sync"

// SyncCategoryCommand represents a command to sync a category
type SyncCategoryCommand struct {
	Code  string
	Debug bool
}

// Type returns the command type
func (c SyncCategoryCommand) Type() bus.Type {
	return SyncCategoryCommandType
}
