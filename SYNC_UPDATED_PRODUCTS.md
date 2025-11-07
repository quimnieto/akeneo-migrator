# Sync Updated Products Feature

## Overview

The sync updated products feature allows you to synchronize all products and product models that have been updated since a specific date, including their complete hierarchies.

## How It Works

### Optimized Hierarchy Sync

The system uses an intelligent approach to avoid duplicate syncs:

1. **Fetch all updated products/models** since the specified date
2. **Filter to get only root/common items** (those without parent)
3. **Sync complete hierarchies** for each root item
4. **Track synced hierarchies** to avoid duplicates

**Why this is efficient:**
- If a variant product is updated, we detect its root and sync the entire hierarchy once
- If multiple products in the same hierarchy are updated, we only sync the hierarchy once
- Reduces API calls and processing time significantly

**Example:**
```
Updated items: VARIANT-001, VARIANT-002, MODEL-001 (all in same hierarchy)
Root detected: COMMON-001
Action: Sync COMMON-001 hierarchy once (includes all variants and models)
```

## Usage

### Basic Command

```bash
./akeneo-migrator sync-updated-products 2024-01-01T00:00:00
```

### With Debug Mode

```bash
./akeneo-migrator sync-updated-products 2024-01-15T10:30:00 --debug
```

## Date Format

**IMPORTANT: All dates are interpreted and processed in UTC timezone.**

The date can be in any of these formats:

**ISO 8601 with timezone (recommended):**
```
YYYY-MM-DDTHH:MM:SS+00:00
```

**ISO 8601 without timezone (assumes UTC):**
```
YYYY-MM-DDTHH:MM:SS
```

**Akeneo format (assumes UTC):**
```
YYYY-MM-DD HH:MM:SS
```

Examples:
- `2024-01-01T00:00:00+00:00` - Start of January 1st, 2024 UTC (explicit timezone)
- `2024-01-01T00:00:00` - Start of January 1st, 2024 UTC (assumed)
- `2024-01-01 00:00:00` - Start of January 1st, 2024 UTC (assumed)
- `2024-01-15T10:30:00` - January 15th, 2024 at 10:30 AM UTC
- `2024-12-31T23:59:59` - End of December 31st, 2024 UTC

The tool automatically:
1. Parses the input date in any of the above formats
2. Converts it to UTC if a timezone is specified
3. Formats it as `YYYY-MM-DD HH:MM:SS` for Akeneo's API (in UTC)
4. Filters only root products/models (without parent) to avoid duplicates

## Examples

### Sync Last 24 Hours

```bash
# Get yesterday's date
YESTERDAY=$(date -u -d '1 day ago' '+%Y-%m-%dT%H:%M:%S')
./akeneo-migrator sync-updated-products $YESTERDAY
```

### Sync Last Week

```bash
# Get date from 7 days ago
LAST_WEEK=$(date -u -d '7 days ago' '+%Y-%m-%dT%H:%M:%S')
./akeneo-migrator sync-updated-products $LAST_WEEK
```

### Sync Last Month

```bash
# Get date from 30 days ago
LAST_MONTH=$(date -u -d '30 days ago' '+%Y-%m-%dT%H:%M:%S')
./akeneo-migrator sync-updated-products $LAST_MONTH
```

### Sync Specific Date Range

```bash
# Sync everything updated since January 1st, 2024
./akeneo-migrator sync-updated-products 2024-01-01T00:00:00
```

## Output Example

```
🚀 Starting synchronization of products updated since: 2024-01-15T00:00:00
📅 Fetching products updated since: 2024-01-15T00:00:00
   📦 Found 5 updated models
   📦 Found 12 updated products
   🎯 Identified 2 common models (root level)
   🎯 Identified 3 common products (root level)
   🔄 Syncing model hierarchy: COMMON-MODEL-001
   🔄 Syncing model hierarchy: COMMON-MODEL-002
   🔄 Syncing product hierarchy: COMMON-PRODUCT-001
   🔄 Syncing product hierarchy: COMMON-PRODUCT-002
   🔄 Syncing product hierarchy: COMMON-PRODUCT-003

📋 Synchronization Summary:
   📅 Updated since: 2024-01-15T00:00:00
   📦 Models synced: 8
   📦 Products synced: 24
   📊 Total synced: 32

✅ Synchronization completed successfully!
```

## Use Cases

### 1. Incremental Sync

Sync only changes since last sync:

```bash
# Store last sync time
echo "2024-01-15T10:00:00" > .last_sync

# Next sync uses that timestamp
LAST_SYNC=$(cat .last_sync)
./akeneo-migrator sync-updated-products $LAST_SYNC

# Update timestamp
date -u '+%Y-%m-%dT%H:%M:%S' > .last_sync
```

### 2. Scheduled Sync

Run hourly via cron:

```bash
# Add to crontab
0 * * * * cd /path/to/akeneo-migrator && ./sync-hourly.sh >> /var/log/akeneo-sync.log 2>&1
```

`sync-hourly.sh`:
```bash
#!/bin/bash
ONE_HOUR_AGO=$(date -u -d '1 hour ago' '+%Y-%m-%dT%H:%M:%S')
./bin/akeneo-migrator sync-updated-products $ONE_HOUR_AGO
```

### 3. Real-time Sync

Run every 5 minutes:

```bash
# Add to crontab
*/5 * * * * cd /path/to/akeneo-migrator && ./sync-recent.sh
```

`sync-recent.sh`:
```bash
#!/bin/bash
FIVE_MINUTES_AGO=$(date -u -d '5 minutes ago' '+%Y-%m-%dT%H:%M:%S')
./bin/akeneo-migrator sync-updated-products $FIVE_MINUTES_AGO
```

## Optimization: Root-Only Sync

### How It Works

The system uses a two-level optimization approach:

**Level 1: API Filtering (New)**
- Filters at the API level using `parent:[{"operator":"EMPTY"}]`
- Only fetches root products/models from Akeneo
- Dramatically reduces API calls and data transfer

**Level 2: Application Logic**
1. **Identifies root items**: Filters products/models without parent
2. **Syncs hierarchies once**: Each root triggers a complete hierarchy sync
3. **Avoids duplicates**: Tracks synced hierarchies to prevent re-syncing

### Example Optimization

**Scenario: Multiple updates in same hierarchy**

```
Updated items detected:
- VARIANT-001 (updated)
- VARIANT-002 (updated)
- MODEL-001 (updated)
- COMMON-001 (updated)

Traditional approach:
❌ Sync VARIANT-001 hierarchy → 10 API calls
❌ Sync VARIANT-002 hierarchy → 10 API calls (duplicate!)
❌ Sync MODEL-001 hierarchy → 10 API calls (duplicate!)
❌ Sync COMMON-001 hierarchy → 10 API calls (duplicate!)
Total: 40 API calls

Optimized approach:
✅ Filter: Only COMMON-001 has no parent
✅ Sync COMMON-001 hierarchy once → 10 API calls
Total: 10 API calls (75% reduction!)
```

### Performance Benefits

- **Fewer API calls**: Only syncs each hierarchy once
- **Faster execution**: No duplicate processing
- **Lower load**: Reduces stress on both source and destination APIs
- **Same result**: Complete hierarchies are still fully synchronized

## How Hierarchies Are Synced

### Scenario 1: Updated Model

```
MODEL-001 (updated) ← Detected
├── VARIANT-001
├── VARIANT-002
└── VARIANT-003

Syncs: MODEL-001 + all 3 variants
```

### Scenario 2: Updated Variant

```
MODEL-001
├── VARIANT-001
├── VARIANT-002 (updated) ← Detected
└── VARIANT-003

Syncs: MODEL-001 + all 3 variants (entire hierarchy)
```

### Scenario 3: Updated Common Product

```
COMMON-001 (updated) ← Detected
├── CHILD-001
├── CHILD-002
└── CHILD-003

Syncs: COMMON-001 + all 3 children
```

### Scenario 4: Complex Hierarchy

```
COMMON-001 (updated) ← Detected
├── MODEL-001
│   ├── VARIANT-001
│   └── VARIANT-002
└── MODEL-002
    ├── VARIANT-003
    └── VARIANT-004

Syncs: COMMON-001 + 2 models + 4 variants (entire tree)
```

## Performance Considerations

### Large Datasets

For many updated products:

```bash
# Use debug mode to monitor progress
./akeneo-migrator sync-updated-products 2024-01-01T00:00:00 --debug
```

### API Rate Limits

The tool respects Akeneo API rate limits through:
- Token refresh mechanism
- Pagination for large result sets
- Sequential processing to avoid overwhelming the API

### Memory Usage

For very large sync operations, consider:
- Syncing smaller time windows
- Running during off-peak hours
- Monitoring system resources

## Error Handling

### Common Errors

**Invalid Date Format:**
```bash
$ ./akeneo-migrator sync-updated-products 2024-01-01

❌ Date must be in ISO 8601 format: YYYY-MM-DDTHH:MM:SS
```

**No Products Found:**
```bash
$ ./akeneo-migrator sync-updated-products 2024-12-31T00:00:00

📅 Fetching products updated since: 2024-12-31T00:00:00
   📦 Found 0 updated models
   📦 Found 0 updated products

📋 Synchronization Summary:
   📊 Total synced: 0

✅ Synchronization completed successfully!
```

**Partial Failures:**
```bash
📋 Synchronization Summary:
   📦 Models synced: 5
   📦 Products synced: 10
   ⚠️  Errors: 2

⚠️  Synchronization completed with errors
```

## Best Practices

### 1. Start with Recent Changes

Test with a recent date first:

```bash
# Test with last hour
ONE_HOUR_AGO=$(date -u -d '1 hour ago' '+%Y-%m-%dT%H:%M:%S')
./akeneo-migrator sync-updated-products $ONE_HOUR_AGO
```

### 2. Use Debug Mode Initially

```bash
./akeneo-migrator sync-updated-products 2024-01-15T00:00:00 --debug
```

### 3. Monitor First Run

Watch the output to understand:
- How many products are affected
- How long it takes
- If there are any errors

### 4. Automate Gradually

Start with manual runs, then:
1. Daily sync
2. Hourly sync
3. Real-time sync (every 5-15 minutes)

### 5. Keep Logs

```bash
./akeneo-migrator sync-updated-products $DATE 2>&1 | tee -a sync.log
```

## Integration with Other Commands

### Complete Sync Strategy

```bash
# 1. Initial full sync (one-time)
./akeneo-migrator sync brands
./akeneo-migrator sync-attribute sku
./akeneo-migrator sync-product COMMON-001

# 2. Ongoing incremental sync (automated)
./akeneo-migrator sync-updated-products $(date -u -d '1 hour ago' '+%Y-%m-%dT%H:%M:%S')
```

## Troubleshooting

### Duplicate Syncs

If the same product is synced multiple times:
- This is expected when multiple products in a hierarchy are updated
- The tool ensures the entire hierarchy is consistent

### Missing Products

If expected products aren't synced:
- Check the date format
- Verify the product was actually updated after that date
- Use `--debug` to see what's being fetched

### Performance Issues

If sync is slow:
- Reduce the time window
- Check network latency
- Monitor API rate limits

## Future Enhancements

Planned improvements:
- Parallel processing for faster syncs
- Deduplication to avoid syncing same hierarchy multiple times
- Progress bar for large operations
- Dry-run mode to preview changes
- Webhook integration for real-time triggers
