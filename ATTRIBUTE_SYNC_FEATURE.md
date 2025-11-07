# Attribute Synchronization Feature

## Overview

The attribute synchronization feature allows you to sync individual attributes from a source Akeneo PIM instance to a destination instance.

## Architecture

The implementation follows the hexagonal architecture pattern:

```
internal/attribute/
├── repository.go              # Domain interfaces
└── syncing/
    └── service.go            # Business logic

internal/platform/storage/akeneo/
└── attribute_repository.go   # Akeneo implementation

internal/platform/client/akeneo/
└── client.go                 # HTTP client methods
```

### Components

1. **Domain Layer** (`internal/attribute/`)
   - `Attribute`: Type alias for `map[string]interface{}`
   - `SourceRepository`: Read-only interface for source operations
   - `DestRepository`: Write-only interface for destination operations

2. **Application Layer** (`internal/attribute/syncing/`)
   - `Service`: Orchestrates the sync process
   - `SyncResult`: Contains sync operation results

3. **Infrastructure Layer** (`internal/platform/`)
   - `SourceAttributeRepository`: Implements read operations
   - `DestAttributeRepository`: Implements write operations
   - `Client.GetAttribute()`: Fetches attribute from API
   - `Client.PatchAttribute()`: Creates/updates attribute via API

## Usage

### Command Line

Sync a single attribute:

```bash
akeneo-migrator sync-attribute sku
```

With debug mode:

```bash
akeneo-migrator sync-attribute description --debug
```

### Examples

```bash
# Sync SKU attribute
akeneo-migrator sync-attribute sku

# Sync description attribute with debug output
akeneo-migrator sync-attribute description --debug

# Sync custom attribute
akeneo-migrator sync-attribute my_custom_attribute
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

- **GET** `/api/rest/v1/attributes/{code}` - Fetch attribute
- **PATCH** `/api/rest/v1/attributes/{code}` - Create/update attribute

## Error Handling

The service handles:
- Attribute not found (404)
- Validation errors (422)
- Network errors
- Authentication errors

## Field Cleaning

The client automatically removes metadata fields before sending:
- `_links`

## Future Enhancements

Possible improvements:
- Bulk attribute synchronization
- Attribute group synchronization
- Attribute option synchronization
- Dry-run mode
- Conflict resolution strategies
