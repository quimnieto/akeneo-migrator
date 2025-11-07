# Product Synchronization Feature

## Overview

The product synchronization feature allows you to sync complete product hierarchies from a source Akeneo instance to a destination instance.

## Architecture

```
CLI Command (sync-product COMMON-001)
    ↓
Command Handler
    ↓
Service.Sync()
    ↓
1. Sync common product/model
2. Sync child models (if configurable)
3. Sync variant products
```

## Usage

### Sync a Product Hierarchy

```bash
# Sync a simple product with children
./akeneo-migrator sync-product COMMON-001

# Sync a configurable product (common → models → variants)
./akeneo-migrator sync-product COMMON-CONFIG-001

# With debug mode
./akeneo-migrator sync-product COMMON-001 --debug
```

## How It Works

### For Simple Products

```
COMMON-001 (simple)
├── CHILD-001
├── CHILD-002
└── CHILD-003

Syncs: COMMON-001 + all 3 children
```

### For Configurable Products

```
COMMON-001 (configurable)
├── MODEL-001
│   ├── VARIANT-001
│   └── VARIANT-002
└── MODEL-002
    ├── VARIANT-003
    └── VARIANT-004

Syncs: COMMON-001 + 2 models + 4 variants (entire tree)
```

## What Gets Synchronized

- **Identifier** (SKU)
- **Family**
- **Categories**
- **Enabled status**
- **Values** (all attribute values)
- **Associations**
- **Quantified associations**
- **Parent** (for variants)
- **Groups**

## Excluded Fields

Metadata fields are automatically excluded:
- `_links` - API navigation links
- `created` - Creation timestamp
- `updated` - Last update timestamp

## Components

- **Service** (`service.go`): Orchestrates hierarchy sync
- **Repository** (`internal/product/repository.go`): Data access interface
- **Client** (`internal/platform/client/akeneo/client.go`): Akeneo API calls
- **Command Handler** (`command_handler.go`): CLI command handling

## Testing

```bash
# Run tests
go test ./internal/product/syncing/...

# With coverage
go test -cover ./internal/product/syncing/...
```

## Use Cases

1. **Product Migration** - Move products between environments
2. **Product Recovery** - Restore products from backup
3. **Selective Sync** - Sync specific product hierarchies
4. **Testing** - Test product sync before bulk operations

## Limitations

- Requires product family to exist in destination
- Does not sync media files (only references)
- Does not sync dependencies (families, attributes, etc.)

## API Endpoints Used

### Source Akeneo
- `GET /api/rest/v1/products/{identifier}`
- `GET /api/rest/v1/product-models/{code}`
- `GET /api/rest/v1/products?search={"parent":[{"operator":"=","value":"..."}]}`
- `GET /api/rest/v1/product-models?search={"parent":[{"operator":"=","value":"..."}]}`

### Destination Akeneo
- `PATCH /api/rest/v1/products/{identifier}`
- `PATCH /api/rest/v1/product-models/{code}`
