# Sync Updated Products Feature

## Overview

The sync updated products feature allows you to synchronize all products and product models that have been updated since a specific date, including their complete hierarchies.

## How It Works

### Smart Hierarchy Detection

The system uses an intelligent approach to sync complete hierarchies efficiently:

**1. Fetch ALL Updated Items**
- Gets all products and models updated since the specified date
- Includes root items, child models, and variant products

**2. Navigate to Root**
- For each updated item, navigates up the hierarchy to find the root
- Example: If variant `VAR-001` is updated, finds its parent model `MODEL-001`, then finds the root `COMMON-001`

**3. Sync from Root**
- Syncs the entire hierarchy starting from the root
- Ensures all related items are synchronized together
- Tracks synced hierarchies to avoid duplicates

**4. Memory-Efficient Streaming**
- Processes items in batches of 100
- Never loads all items into memory at once
- Constant memory usage regardless of dataset size

### Example Optimization

**Scenario: Multiple updates in same hierarchy**

```
Updated items detected:
- VARIANT-001 (updated) → parent: MODEL-001
- VARIANT-002 (updated) → parent: MODEL-001  
- MODEL-001 (updated) → parent: COMMON-001
- COMMON-001 (updated) → parent: null

Smart hierarchy detection:
1. Process VARIANT-001 → Navigate up → Find root: COMMON-001 → Sync COMMON-001 hierarchy
2. Process VARIANT-002 → Navigate up → Find root: COMMON-001 → Already synced, skip
3. Process MODEL-001 → Navigate up → Find root: COMMON-001 → Already synced, skip
4. Process COMMON-001 → No parent → Root is COMMON-001 → Already synced, skip

Result:
✅ Sync COMMON-001 hierarchy once → 10 API calls
✅ All 4 updated items are included in the sync
✅ No duplicates, no wasted API calls
Total: 10 API calls (75% reduction!)
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

## Performance Considerations

### Memory-Efficient Streaming

The tool uses a **streaming architecture** to handle large datasets efficiently:

**How it works:**
1. Fetches products/models in batches of 100
2. Processes each batch immediately
3. Releases memory before fetching the next batch
4. Never loads all products into memory at once

**Benefits:**
- ✅ Can handle millions of products without running out of memory
- ✅ Constant memory usage regardless of dataset size
- ✅ Starts syncing immediately (no waiting to fetch all data)
- ✅ Resilient to interruptions (can resume from where it stopped)

### Scalability

**Memory usage:** ~Constant (only holds 100 products/models at a time)  
**Processing speed:** ~100-200 products per minute (depends on hierarchy complexity)  
**Recommended for:** Any dataset size, from hundreds to millions of products

## Examples

### Sync Last 24 Hours

```bash
YESTERDAY=$(date -u -d '1 day ago' '+%Y-%m-%dT%H:%M:%S')
./akeneo-migrator sync-updated-products $YESTERDAY
```

### Sync Last Week

```bash
LAST_WEEK=$(date -u -d '7 days ago' '+%Y-%m-%dT%H:%M:%S')
./akeneo-migrator sync-updated-products $LAST_WEEK
```

### Incremental Sync

```bash
# Store last sync time
echo "2024-01-15T10:00:00" > .last_sync

# Next sync uses that timestamp
LAST_SYNC=$(cat .last_sync)
./akeneo-migrator sync-updated-products $LAST_SYNC

# Update timestamp
date -u '+%Y-%m-%dT%H:%M:%S' > .last_sync
```

## Architecture

### Components

- **Service** (`service.go`): Orchestrates the sync process
- **Repository** (`internal/product/repository.go`): Defines data access interface
- **Client** (`internal/platform/client/akeneo/client.go`): Implements Akeneo API calls
- **Command Handler** (`command_handler.go`): Handles CLI commands

### Flow

```
CLI Command
    ↓
Command Handler
    ↓
Service.Sync()
    ↓
Repository.StreamProductsUpdatedSince() ← Batches of 100
    ↓
For each product:
    - Find root (navigate up hierarchy)
    - Sync entire hierarchy from root
    - Track to avoid duplicates
```

## Best Practices

1. **Start with recent changes** - Test with last hour first
2. **Use debug mode initially** - Understand what's being synced
3. **Monitor first run** - Check timing and errors
4. **Automate gradually** - Daily → Hourly → Real-time
5. **Keep logs** - Use `tee` to save output

## Troubleshooting

### Missing Products

- Check date format is correct
- Verify product was updated after specified date
- Use `--debug` to see API responses

### Performance Issues

- Reduce time window
- Check network latency
- Monitor API rate limits

### Memory Issues

- Should not occur due to streaming architecture
- If it does, reduce batch size in code (default: 100)
