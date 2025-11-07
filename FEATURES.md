# Features

## Reference Entity Synchronization

The `sync` command performs a complete synchronization of a Reference Entity from source to destination.

### What Gets Synchronized

#### 1. Reference Entity Definition
- Entity code
- Labels (all locales)
- Image attribute (if configured)
- Attributes configuration:
  - Attribute codes
  - Attribute types
  - Labels
  - Value per locale/channel settings
  - Validation rules
  - Options (for select/multi-select attributes)

#### 2. Reference Entity Records
- All record codes
- All attribute values
- Labels
- Images
- Localized values
- Channel-specific values

### Behavior

#### Creating New Reference Entity
If the Reference Entity doesn't exist in the destination:
- The entity definition is created first
- Then all records are created

#### Updating Existing Reference Entity
If the Reference Entity already exists in the destination:
- The entity definition is updated (attributes, labels, etc.)
- Records are created or updated (PATCH operation)
- Existing records not in source remain unchanged

### Example

```bash
# Synchronize the 'brands' Reference Entity
./akeneo-migrator sync brands
```

**Output:**
```
🚀 Starting synchronization for entity: brands
📋 Synchronizing Reference Entity 'brands'...
   1️⃣  Syncing entity definition...
   2️⃣  Syncing attributes...
   3️⃣  Syncing records...
📊 Found 150 records to synchronize
✅ Successfully synchronized records: 150

📋 Synchronization summary:
   ✅ Successfully synchronized records: 150
   ❌ Records with errors: 0
   📊 Total processed: 150

🎉 Synchronization completed successfully!
```

### Use Cases

1. **Initial Migration**: Copy a complete Reference Entity from one Akeneo to another
2. **Continuous Sync**: Keep Reference Entities in sync between environments (dev → staging → prod)
3. **Backup/Restore**: Create backups by syncing to a backup instance
4. **Multi-tenant**: Sync Reference Entities across multiple Akeneo instances

### Limitations

- Only synchronizes one Reference Entity at a time
- Does not delete records in destination that don't exist in source
- Requires the same attribute types in both instances
- Image files are referenced by URL (not copied)

### Future Enhancements

Planned features:
- [ ] Sync multiple Reference Entities at once
- [ ] Delete records in destination not present in source (--delete flag)
- [ ] Dry-run mode to preview changes
- [ ] Incremental sync (only changed records)
- [ ] Sync Reference Entity attributes separately
- [ ] Image file copying
