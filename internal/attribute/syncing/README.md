# Attribute Synchronization

## Overview

Synchronizes individual attributes from source to destination Akeneo instance.

## Usage

```bash
# Sync a single attribute
./akeneo-migrator sync-attribute sku

# With debug mode
./akeneo-migrator sync-attribute description --debug
```

## What Gets Synchronized

- Attribute code
- Attribute type (text, number, date, etc.)
- Labels (all locales)
- Group
- Sort order
- Required flag
- Unique flag
- Localizable flag
- Scopable flag
- Available locales
- Type-specific options

## Components

- **Service** (`service.go`): Sync orchestration
- **Repository** (`internal/attribute/repository.go`): Data access interface
- **Client** (`internal/platform/client/akeneo/client.go`): API calls

## API Endpoints

### Source
- `GET /api/rest/v1/attributes/{code}`

### Destination
- `PATCH /api/rest/v1/attributes/{code}`

## Limitations

- Syncs one attribute at a time
- Does not sync attribute options (for select/multiselect)
- Requires attribute group to exist in destination
