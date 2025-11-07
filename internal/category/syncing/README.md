# Category Synchronization

## Overview

Synchronizes individual categories from source to destination Akeneo instance.

## Usage

```bash
# Sync a single category
./akeneo-migrator sync-category master

# With debug mode
./akeneo-migrator sync-category clothing --debug
```

## What Gets Synchronized

- Category code
- Labels (all locales)
- Parent category
- Sort order

## Components

- **Service** (`service.go`): Sync orchestration
- **Repository** (`internal/category/repository.go`): Data access interface
- **Client** (`internal/platform/client/akeneo/client.go`): API calls

## API Endpoints

### Source
- `GET /api/rest/v1/categories/{code}`

### Destination
- `PATCH /api/rest/v1/categories/{code}`

## Limitations

- Syncs one category at a time
- Does not sync category tree (only single category)
- Parent category must exist in destination
