# Product Synchronization Feature

## Overview

The product synchronization feature allows you to sync individual products from a source Akeneo instance to a destination instance.

## Architecture

Following the same hexagonal architecture pattern as Reference Entities:

```
┌─────────────────────────────────────────────────────────────┐
│                    CLI Command                               │
│              (sync-product SKU-12345)                        │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                  Product Syncing Service                     │
│            (internal/product/syncing)                        │
│                                                              │
│  - Fetches product from source                              │
│  - Saves product to destination                             │
└──────────────┬──────────────────────────┬───────────────────┘
               │                          │
               ▼                          ▼
┌──────────────────────────┐  ┌──────────────────────────┐
│   Source Repository      │  │    Dest Repository       │
│  (interface)             │  │   (interface)            │
└──────────┬───────────────┘  └──────────┬───────────────┘
           │                              │
           ▼                              ▼
┌──────────────────────────┐  ┌──────────────────────────┐
│ Akeneo Repository Impl   │  │ Akeneo Repository Impl   │
│ (platform/storage)       │  │ (platform/storage)       │
└──────────┬───────────────┘  └──────────┬───────────────┘
           │                              │
           ▼                              ▼
┌──────────────────────────┐  ┌──────────────────────────┐
│   Akeneo HTTP Client     │  │   Akeneo HTTP Client     │
│   (platform/client)      │  │   (platform/client)      │
└──────────┬───────────────┘  └──────────┬───────────────┘
           │                              │
           ▼                              ▼
     Source Akeneo API            Dest Akeneo API
```

## Files Created

### Domain Layer
- `internal/product/repository.go` - Repository interfaces
- `internal/product/syncing/service.go` - Synchronization service
- `internal/product/syncing/service_test.go` - Unit tests

### Infrastructure Layer
- `internal/platform/client/akeneo/client.go` - Added product methods:
  - `GetProduct(identifier string)` - Fetch product
  - `PatchProduct(identifier, product)` - Create/update product
  - `cleanProduct(product)` - Clean metadata fields
- `internal/platform/storage/akeneo/product_repository.go` - Repository implementations

### Application Layer
- `cmd/app/bootstrap/bootstrap.go` - Added:
  - Product repositories initialization
  - Product syncer service
  - `sync-product` command

## Usage

### Basic Sync

```bash
./akeneo-migrator sync-product SKU-12345
```

### With Debug Mode

```bash
./akeneo-migrator sync-product SKU-12345 --debug
```

## What Gets Synchronized

When syncing a product, the following data is transferred:

- **Identifier** (SKU)
- **Family**
- **Categories**
- **Enabled status**
- **Values** (all attribute values)
  - Text attributes
  - Number attributes
  - Date attributes
  - Boolean attributes
  - Select/Multi-select options
  - Images
  - Files
  - Prices
  - Metrics
- **Associations**
- **Quantified associations**
- **Parent** (for product variants)
- **Groups**

## Excluded Fields

The following metadata fields are automatically excluded:
- `_links` - API navigation links
- `created` - Creation timestamp
- `updated` - Last update timestamp

## Example Output

```bash
$ ./akeneo-migrator sync-product SKU-12345

🚀 Starting synchronization for product: SKU-12345
📥 Fetching product 'SKU-12345' from source...
✅ Product 'SKU-12345' synchronized successfully!
```

## Error Handling

### Product Not Found

```bash
❌ Synchronization error: error fetching product from source: product 'SKU-12345' not found
```

### Validation Error

```bash
❌ Synchronization error: error saving product to destination: validation error in product SKU-12345: ...
```

## Testing

### Run Tests

```bash
# Run product sync tests
go test ./internal/product/syncing/...

# Run with coverage
go test -cover ./internal/product/syncing/...
```

### Test Coverage

The service includes tests for:
- ✅ Successful synchronization
- ✅ Source repository errors
- ✅ Destination repository errors

## Use Cases

1. **Single Product Migration** - Move a specific product between environments
2. **Product Testing** - Test product sync before bulk operations
3. **Product Recovery** - Restore a specific product from backup
4. **Selective Sync** - Sync only specific products instead of all

## Limitations

- Syncs one product at a time
- Does not sync product models (use for simple and variant products)
- Does not sync associated media files (only references)
- Requires product family to exist in destination

## Future Enhancements

Potential improvements:
- [ ] Bulk product sync (multiple SKUs)
- [ ] Product model sync
- [ ] Media file copying
- [ ] Dry-run mode
- [ ] Progress bar for large products
- [ ] Sync product with dependencies (family, attributes, etc.)
- [ ] Incremental sync (only changed values)

## Comparison with Reference Entity Sync

| Feature | Reference Entity | Product |
|---------|-----------------|---------|
| Scope | All records | Single product |
| Structure | Entity + Attributes + Records | Product data only |
| Dependencies | Creates entity if missing | Requires family to exist |
| Use Case | Complete entity migration | Selective product sync |

## Integration with CI/CD

The product sync feature is fully tested in the CI pipeline:
- Unit tests run on every commit
- Code coverage tracked
- Linting enforced

## API Endpoints Used

### Source Akeneo
- `GET /api/rest/v1/products/{identifier}` - Fetch product

### Destination Akeneo
- `PATCH /api/rest/v1/products/{identifier}` - Create/update product
