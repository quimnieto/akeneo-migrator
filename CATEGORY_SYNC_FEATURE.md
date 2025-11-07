# Category Synchronization Feature

## Overview

The category synchronization feature allows you to sync individual categories from a source Akeneo PIM instance to a destination instance.

## Architecture

The implementation follows the hexagonal architecture pattern with Command Bus:

```
internal/category/
├── repository.go              # Domain interfaces
└── syncing/
    ├── service.go            # Business logic
    ├── command.go            # Command definition
    ├── command_handler.go    # Command handler
    └── service_test.go       # Unit tests

internal/platform/storage/akeneo/
└── category_repository.go    # Akeneo implementation

internal/platform/client/akeneo/
└── client.go                 # HTTP client methods
```

### Components

1. **Domain Layer** (`internal/category/`)
   - `Category`: Type alias for `map[string]interface{}`
   - `SourceRepository`: Read-only interface for source operations
   - `DestRepository`: Write-only interface for destination operations

2. **Application Layer** (`internal/category/syncing/`)
   - `Service`: Orchestrates the sync process
   - `SyncResult`: Contains sync operation results
   - `SyncCategoryCommand`: Command definition
   - `CommandHandler`: Executes the command

3. **Infrastructure Layer** (`internal/platform/`)
   - `SourceCategoryRepository`: Implements read operations
   - `DestCategoryRepository`: Implements write operations
   - `Client.GetCategory()`: Fetches category from API
   - `Client.PatchCategory()`: Creates/updates category via API

## Usage

### Command Line

Sync a single category:

```bash
akeneo-migrator sync-category master
```

With debug mode:

```bash
akeneo-migrator sync-category clothing --debug
```

### Examples

```bash
# Sync master category
akeneo-migrator sync-category master

# Sync clothing category with debug output
akeneo-migrator sync-category clothing --debug

# Sync custom category
akeneo-migrator sync-category my_custom_category
```

## Configuration

Uses the same configuration as other sync commands:

```bash
export SOURCE_HOST="https://source.akeneo.com"
export SOURCE_CLIENT_ID="your_client_id"
export SOURCE_SECRET="your_secret"
export SOURCE_USERNAME="your_username"
export SOURCE_PASSWORD="your_password"

export DEST_HOST="https://destination.akeneo.com"
export DEST_CLIENT_ID="your_client_id"
export DEST_SECRET="your_secret"
export DEST_USERNAME="your_username"
export DEST_PASSWORD="your_password"
```

## API Endpoints Used

- **GET** `/api/rest/v1/categories/{code}` - Fetch category
- **PATCH** `/api/rest/v1/categories/{code}` - Create/update category

## Category Structure

A category in Akeneo typically contains:

```json
{
  "code": "clothing",
  "parent": "master",
  "labels": {
    "en_US": "Clothing",
    "es_ES": "Ropa"
  }
}
```

## Error Handling

The service handles:
- Category not found (404)
- Validation errors (422)
- Network errors
- Authentication errors

## Field Cleaning

The client automatically removes metadata fields before sending:
- `_links`

## Sync Order

When migrating a complete category tree, sync in hierarchical order:

```bash
# 1. Root category first
./akeneo-migrator sync-category master

# 2. Then first-level children
./akeneo-migrator sync-category clothing
./akeneo-migrator sync-category electronics

# 3. Then second-level children
./akeneo-migrator sync-category mens_clothing
./akeneo-migrator sync-category womens_clothing
```

## Batch Synchronization

### Sync Multiple Categories

```bash
#!/bin/bash
# sync-categories.sh

categories=(
  "master"
  "clothing"
  "mens_clothing"
  "womens_clothing"
  "electronics"
  "computers"
  "phones"
)

for category in "${categories[@]}"; do
  echo "Syncing category: $category"
  ./bin/akeneo-migrator sync-category "$category"
  
  if [ $? -eq 0 ]; then
    echo "✅ Success"
  else
    echo "❌ Failed"
  fi
  
  echo "---"
done
```

### Sync Category Tree

```bash
#!/bin/bash
# sync-category-tree.sh

# Function to sync a category and its children
sync_tree() {
  local category=$1
  local indent=$2
  
  echo "${indent}Syncing: $category"
  ./bin/akeneo-migrator sync-category "$category"
  
  # Add child categories here
  # sync_tree "child_category" "  $indent"
}

# Start with root
sync_tree "master" ""
sync_tree "clothing" "  "
sync_tree "mens_clothing" "    "
sync_tree "womens_clothing" "    "
```

## Integration with Other Features

### Complete Migration Workflow

```bash
# 1. Categories (structure)
./bin/akeneo-migrator sync-category master
./bin/akeneo-migrator sync-category clothing

# 2. Attributes (metadata)
./bin/akeneo-migrator sync-attribute sku
./bin/akeneo-migrator sync-attribute name

# 3. Reference Entities (data)
./bin/akeneo-migrator sync brands

# 4. Products (final data)
./bin/akeneo-migrator sync-product COMMON-001
```

## Testing

### Unit Tests

```bash
# Run category sync tests
go test ./internal/category/syncing/...
```

### Manual Testing

```bash
# Test with a simple category
./bin/akeneo-migrator sync-category test_category --debug

# Verify in destination Akeneo
# Check that category exists with correct parent and labels
```

## Troubleshooting

### Parent Category Not Found

If you get an error about parent category not existing:

```bash
# Sync parent first
./bin/akeneo-migrator sync-category master

# Then sync child
./bin/akeneo-migrator sync-category clothing
```

### Circular Dependencies

Akeneo prevents circular category references. Ensure your category tree is acyclic.

### Label Validation

Categories require at least one label. Ensure your source categories have labels defined.

## Future Enhancements

Possible improvements:
- Bulk category synchronization
- Category tree synchronization (recursive)
- Dry-run mode
- Conflict resolution strategies
- Category mapping/transformation
- Automatic parent resolution
