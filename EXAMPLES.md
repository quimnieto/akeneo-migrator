# Usage Examples

## Reference Entity Synchronization

### Sync Complete Reference Entity

```bash
# Sync entity definition + attributes + all records
./bin/akeneo-migrator sync brands
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

### Debug Mode

```bash
./bin/akeneo-migrator sync brands --debug
```

Shows detailed information about each record being synced.

## Product Synchronization

### Sync Simple Product Hierarchy (2 levels)

```bash
# Common → Child Products
./bin/akeneo-migrator sync-product SIMPLE-COMMON-001
```

**Structure:**
```
SIMPLE-COMMON-001 (Common Product)
├── SKU-001-RED
├── SKU-001-BLUE
└── SKU-001-GREEN
```

**Output:**
```
🚀 Starting synchronization for product: SIMPLE-COMMON-001
📥 Fetching product hierarchy for 'SIMPLE-COMMON-001' from source...
   📦 Syncing common: SIMPLE-COMMON-001
   👶 Found 3 child products
   ✅ Synced product: SKU-001-RED
   ✅ Synced product: SKU-001-BLUE
   ✅ Synced product: SKU-001-GREEN

📋 Synchronization Summary:
   📦 Models synced: 0
   📦 Products synced: 4
   📊 Total synced: 4

✅ Hierarchy 'SIMPLE-COMMON-001' synchronized successfully!
```

### Sync Configurable Product Hierarchy (3 levels)

```bash
# Common → Models → Variant Products
./bin/akeneo-migrator sync-product CONFIG-COMMON-001
```

**Structure:**
```
CONFIG-COMMON-001 (Common Model)
├── MODEL-001-S (Size S)
│   ├── SKU-001-S-RED
│   └── SKU-001-S-BLUE
├── MODEL-001-M (Size M)
│   ├── SKU-001-M-RED
│   └── SKU-001-M-BLUE
└── MODEL-001-L (Size L)
    ├── SKU-001-L-RED
    └── SKU-001-L-BLUE
```

**Output:**
```
🚀 Starting synchronization for product: CONFIG-COMMON-001
📥 Fetching product hierarchy for 'CONFIG-COMMON-001' from source...
   📦 Syncing common: CONFIG-COMMON-001
   📋 Found 3 child models
   ✅ Synced model: MODEL-001-S
   ✅ Synced model: MODEL-001-M
   ✅ Synced model: MODEL-001-L
   🔸 Found 2 variants for model MODEL-001-S
   ✅ Synced variant: SKU-001-S-RED
   ✅ Synced variant: SKU-001-S-BLUE
   🔸 Found 2 variants for model MODEL-001-M
   ✅ Synced variant: SKU-001-M-RED
   ✅ Synced variant: SKU-001-M-BLUE
   🔸 Found 2 variants for model MODEL-001-L
   ✅ Synced variant: SKU-001-L-RED
   ✅ Synced variant: SKU-001-L-BLUE

📋 Synchronization Summary:
   📦 Models synced: 3
   📦 Products synced: 6
   📊 Total synced: 9

✅ Hierarchy 'CONFIG-COMMON-001' synchronized successfully!
```

### Sync Single Product (No Hierarchy)

```bash
# Sync only one product, ignore children
./bin/akeneo-migrator sync-product SKU-12345 --single
```

**Output:**
```
🚀 Starting synchronization for product: SKU-12345
📥 Fetching product 'SKU-12345' from source...

📋 Synchronization Summary:
   📦 Models synced: 0
   📦 Products synced: 1
   📊 Total synced: 1

✅ Hierarchy 'SKU-12345' synchronized successfully!
```

## Common Workflows

### Initial Migration

```bash
# 1. Sync Reference Entities first (structure)
./bin/akeneo-migrator sync brands
./bin/akeneo-migrator sync colors
./bin/akeneo-migrator sync sizes

# 2. Sync product hierarchies
./bin/akeneo-migrator sync-product COMMON-SHOES-001
./bin/akeneo-migrator sync-product COMMON-SHIRTS-001
```

### Selective Sync

```bash
# Sync only specific products
./bin/akeneo-migrator sync-product SKU-SPECIAL-001 --single
./bin/akeneo-migrator sync-product SKU-SPECIAL-002 --single
```

### Testing Before Production

```bash
# Use debug mode to verify data
./bin/akeneo-migrator sync brands --debug
./bin/akeneo-migrator sync-product COMMON-001 --debug
```

## Batch Operations

### Sync Multiple Reference Entities

```bash
#!/bin/bash
# sync-all-entities.sh

entities=("brands" "colors" "sizes" "materials" "features")

for entity in "${entities[@]}"; do
  echo "Syncing $entity..."
  ./bin/akeneo-migrator sync "$entity"
  echo "---"
done
```

### Sync Multiple Products

```bash
#!/bin/bash
# sync-products.sh

products=("COMMON-001" "COMMON-002" "COMMON-003")

for product in "${products[@]}"; do
  echo "Syncing $product..."
  ./bin/akeneo-migrator sync-product "$product"
  echo "---"
done
```

## Error Handling

### Reference Entity Not Found

```bash
$ ./bin/akeneo-migrator sync nonexistent

❌ Synchronization error: error fetching reference entity definition from source: 
reference entity 'nonexistent' not found
```

### Product Not Found

```bash
$ ./bin/akeneo-migrator sync-product INVALID-SKU

❌ Synchronization error: error fetching product from source: 
product 'INVALID-SKU' not found
```

### Validation Errors

```bash
$ ./bin/akeneo-migrator sync-product COMMON-001

⚠️  Error syncing product SKU-001: validation error in product SKU-001: 
Field 'family': Family 'shoes' does not exist
```

## Configuration Examples

### Development Environment

```json
{
  "akeneoSource": {
    "api": {
      "url": "https://dev-akeneo.example.com",
      "credentials": {
        "clientId": "dev_client_id",
        "secret": "dev_secret",
        "username": "dev_user",
        "password": "dev_pass"
      }
    }
  },
  "akeneoDest": {
    "api": {
      "url": "https://staging-akeneo.example.com",
      "credentials": {
        "clientId": "staging_client_id",
        "secret": "staging_secret",
        "username": "staging_user",
        "password": "staging_pass"
      }
    }
  }
}
```

### Production Migration

```json
{
  "akeneoSource": {
    "api": {
      "url": "https://old-akeneo.example.com",
      "credentials": {
        "clientId": "old_client_id",
        "secret": "old_secret",
        "username": "migration_user",
        "password": "migration_pass"
      }
    }
  },
  "akeneoDest": {
    "api": {
      "url": "https://new-akeneo.example.com",
      "credentials": {
        "clientId": "new_client_id",
        "secret": "new_secret",
        "username": "migration_user",
        "password": "migration_pass"
      }
    }
  }
}
```

## Performance Tips

### Large Reference Entities

For Reference Entities with thousands of records:
- Use `--debug` sparingly (generates lots of output)
- Monitor API rate limits
- Consider syncing during off-peak hours

### Large Product Hierarchies

For products with many variants:
- Sync common products during maintenance windows
- Monitor memory usage
- Check destination storage capacity

## Troubleshooting

### Slow Synchronization

```bash
# Check network latency
ping your-akeneo-instance.com

# Check API response time
time curl -X GET "https://your-akeneo.com/api/rest/v1/products/SKU-001" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Memory Issues

```bash
# Monitor memory usage
top -p $(pgrep akeneo-migrator)

# For large hierarchies, sync in smaller batches
```

### Rate Limiting

If you hit API rate limits:
- Add delays between requests
- Reduce batch sizes
- Contact Akeneo support for higher limits

## Advanced Usage

### Combining Commands

```bash
# Sync entity and then products
./bin/akeneo-migrator sync brands && \
./bin/akeneo-migrator sync-product BRAND-COMMON-001
```

### Logging to File

```bash
# Save output to log file
./bin/akeneo-migrator sync brands 2>&1 | tee sync-brands.log

# Save only errors
./bin/akeneo-migrator sync brands 2> errors.log
```

### Scheduled Sync

```bash
# Add to crontab for daily sync at 2 AM
0 2 * * * cd /path/to/akeneo-migrator && ./bin/akeneo-migrator sync brands >> /var/log/akeneo-sync.log 2>&1
```

## Attribute Synchronization

### Sync Single Attribute

```bash
# Sync SKU attribute
./bin/akeneo-migrator sync-attribute sku
```

**Output:**
```
🚀 Starting synchronization for attribute: sku
✅ Attribute 'sku' synchronized successfully!
```

### Sync with Debug Mode

```bash
./bin/akeneo-migrator sync-attribute description --debug
```

**Output:**
```
🚀 Starting synchronization for attribute: description
🔍 Debug mode enabled
✅ Attribute 'description' synchronized successfully!
```

### Common Attributes to Sync

```bash
# Core attributes
./bin/akeneo-migrator sync-attribute sku
./bin/akeneo-migrator sync-attribute name
./bin/akeneo-migrator sync-attribute description

# Custom attributes
./bin/akeneo-migrator sync-attribute my_custom_field
./bin/akeneo-migrator sync-attribute brand_reference
```

### Batch Attribute Sync

```bash
#!/bin/bash
# sync-attributes.sh

attributes=("sku" "name" "description" "price" "weight")

for attr in "${attributes[@]}"; do
  echo "Syncing attribute: $attr"
  ./bin/akeneo-migrator sync-attribute "$attr"
  echo "---"
done
```

### Attribute Not Found Error

```bash
$ ./bin/akeneo-migrator sync-attribute nonexistent

❌ Synchronization error: error fetching attribute from source: 
attribute 'nonexistent' not found
```

### Complete Migration Workflow

```bash
# 1. Sync attributes first (structure)
./bin/akeneo-migrator sync-attribute sku
./bin/akeneo-migrator sync-attribute name
./bin/akeneo-migrator sync-attribute description

# 2. Sync Reference Entities (data)
./bin/akeneo-migrator sync brands
./bin/akeneo-migrator sync colors

# 3. Sync products (final data)
./bin/akeneo-migrator sync-product COMMON-001
```

## Category Synchronization

### Sync Single Category

```bash
# Sync master category
./bin/akeneo-migrator sync-category master
```

**Output:**
```
🚀 Starting synchronization for category: master
✅ Category 'master' synchronized successfully!
```

### Sync with Debug Mode

```bash
./bin/akeneo-migrator sync-category clothing --debug
```

**Output:**
```
🚀 Starting synchronization for category: clothing
🔍 Debug mode enabled
✅ Category 'clothing' synchronized successfully!
```

### Common Categories to Sync

```bash
# Root categories
./bin/akeneo-migrator sync-category master
./bin/akeneo-migrator sync-category products

# Sub-categories
./bin/akeneo-migrator sync-category clothing
./bin/akeneo-migrator sync-category electronics
./bin/akeneo-migrator sync-category home
```

### Batch Category Sync

```bash
#!/bin/bash
# sync-categories.sh

categories=("master" "clothing" "electronics" "home" "sports")

for cat in "${categories[@]}"; do
  echo "Syncing category: $cat"
  ./bin/akeneo-migrator sync-category "$cat"
  echo "---"
done
```

### Category Not Found Error

```bash
$ ./bin/akeneo-migrator sync-category nonexistent

❌ Synchronization error: error fetching category from source: 
category 'nonexistent' not found
```

### Complete Migration Workflow with Categories

```bash
# 1. Sync categories first (structure)
./bin/akeneo-migrator sync-category master
./bin/akeneo-migrator sync-category clothing
./bin/akeneo-migrator sync-category electronics

# 2. Sync attributes (metadata)
./bin/akeneo-migrator sync-attribute sku
./bin/akeneo-migrator sync-attribute name

# 3. Sync Reference Entities (data)
./bin/akeneo-migrator sync brands
./bin/akeneo-migrator sync colors

# 4. Sync products (final data)
./bin/akeneo-migrator sync-product COMMON-001
```

## Sync Updated Products

### Sync Products Updated Since Specific Date

```bash
# Sync everything updated since January 1st, 2024
./bin/akeneo-migrator sync-updated-products 2024-01-01T00:00:00
```

**Output:**
```
🚀 Starting synchronization of products updated since: 2024-01-01T00:00:00
📅 Fetching products updated since: 2024-01-01T00:00:00
   📦 Found 5 updated models
   📦 Found 12 updated products
   🔄 Syncing model hierarchy: MODEL-001
   🔄 Syncing parent hierarchy for product: SKU-123 (parent: MODEL-001)

📋 Synchronization Summary:
   📅 Updated since: 2024-01-01T00:00:00
   📦 Models synced: 8
   📦 Products synced: 24
   📊 Total synced: 32

✅ Synchronization completed successfully!
```

### Sync Last 24 Hours

```bash
# Get yesterday's date and sync
YESTERDAY=$(date -u -d '1 day ago' '+%Y-%m-%dT%H:%M:%S')
./bin/akeneo-migrator sync-updated-products $YESTERDAY
```

### Sync Last Week

```bash
# Get date from 7 days ago
LAST_WEEK=$(date -u -d '7 days ago' '+%Y-%m-%dT%H:%M:%S')
./bin/akeneo-migrator sync-updated-products $LAST_WEEK
```

### Incremental Sync with Timestamp Tracking

```bash
#!/bin/bash
# incremental-sync.sh

# Read last sync timestamp
if [ -f .last_sync ]; then
  LAST_SYNC=$(cat .last_sync)
else
  # Default to 1 day ago if no previous sync
  LAST_SYNC=$(date -u -d '1 day ago' '+%Y-%m-%dT%H:%M:%S')
fi

echo "Syncing products updated since: $LAST_SYNC"
./bin/akeneo-migrator sync-updated-products $LAST_SYNC

# Save current timestamp for next sync
date -u '+%Y-%m-%dT%H:%M:%S' > .last_sync
echo "Next sync will start from: $(cat .last_sync)"
```

### Scheduled Hourly Sync

```bash
#!/bin/bash
# sync-hourly.sh

ONE_HOUR_AGO=$(date -u -d '1 hour ago' '+%Y-%m-%dT%H:%M:%S')
./bin/akeneo-migrator sync-updated-products $ONE_HOUR_AGO
```

Add to crontab:
```bash
# Run every hour
0 * * * * cd /path/to/akeneo-migrator && ./sync-hourly.sh >> /var/log/akeneo-sync.log 2>&1
```

### Real-time Sync (Every 5 Minutes)

```bash
#!/bin/bash
# sync-realtime.sh

FIVE_MINUTES_AGO=$(date -u -d '5 minutes ago' '+%Y-%m-%dT%H:%M:%S')
./bin/akeneo-migrator sync-updated-products $FIVE_MINUTES_AGO
```

Add to crontab:
```bash
# Run every 5 minutes
*/5 * * * * cd /path/to/akeneo-migrator && ./sync-realtime.sh >> /var/log/akeneo-sync.log 2>&1
```

### Complete Sync Strategy

```bash
#!/bin/bash
# complete-sync-strategy.sh

echo "=== Initial Full Sync (Run Once) ==="
./bin/akeneo-migrator sync brands
./bin/akeneo-migrator sync-attribute sku
./bin/akeneo-migrator sync-category master

echo ""
echo "=== Incremental Sync (Run Regularly) ==="
ONE_HOUR_AGO=$(date -u -d '1 hour ago' '+%Y-%m-%dT%H:%M:%S')
./bin/akeneo-migrator sync-updated-products $ONE_HOUR_AGO
```
