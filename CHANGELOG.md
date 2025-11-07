# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Changed
- **Refactored product module structure**
  - Split into `syncing` and `syncing_since` submodules
  - Removed single product sync (always syncs complete hierarchies)
  - Simplified API: `Sync()` always syncs complete hierarchies
  - Removed `--single` flag from CLI
  - Better naming: `syncing` for hierarchy sync, `syncing_since` for date-based sync
  - Cleaner separation of concerns
  - Each module has single responsibility

### Added
- **Command Bus Architecture**
  - In-memory command bus implementation
  - Command/Handler pattern for all sync operations
  - Middleware support (Logging)
  - Centralized command dispatching
  - Documentation in COMMAND_BUS.md
  
- **Sync Updated Products feature**
  - New `sync-updated-products` command to sync products by update date
  - Optimized root-only sync (only syncs common/root items)
  - Automatic hierarchy deduplication
  - Support for incremental synchronization
  - Date-based filtering with ISO 8601 format
  - Automatic date format conversion (ISO 8601 → Akeneo format)
  - Handles both products and product models
  - Documentation in SYNC_UPDATED_PRODUCTS.md
  
- **Category synchronization feature**
  - New `sync-category` command to sync individual categories
  - Domain interfaces for category operations
  - Akeneo client methods for category API
  - Complete test coverage for category syncing
  - Documentation in CATEGORY_SYNC_FEATURE.md
  
- **Attribute synchronization feature**
  - New `sync-attribute` command to sync individual attributes
  - Domain interfaces for attribute operations
  - Akeneo client methods for attribute API
  - Complete test coverage for attribute syncing
  - Documentation in ATTRIBUTE_SYNC_FEATURE.md

### Changed
- **Refactored Bootstrap**: Now uses Command Bus instead of direct service calls
- **Simplified Application struct**: Only contains Config and CommandBus
- **Reorganized config module**: Moved from `internal/config` to `internal/platform/config` (infrastructure layer)
- Updated golangci-lint configuration to v2 format
- Improved error handling with explicit blank identifier usage
- Fixed shadow variable errors in syncing services

### Documentation
- Added ARCHITECTURE.md with complete architecture documentation
- Updated README.md with improved layer descriptions
- Enhanced COMMAND_BUS.md with usage examples

### Fixed
- CI pipeline compatibility with Go 1.25.0 (local) and Go 1.23 (CI)
- golangci-lint configuration for latest version (v1.64.8)
- Error checking for all defer Close() operations
- Type assertion error checking

## [0.2.0] - 2024-01-XX

### Added
- Product hierarchy synchronization
  - Support for simple products (2-level hierarchy)
  - Support for configurable products (3-level hierarchy)
  - `--single` flag to sync individual products
  - Recursive hierarchy traversal
  - Progress reporting for large hierarchies

### Changed
- Split Repository interface into SourceRepository and DestRepository
- Improved error messages with detailed validation errors
- Enhanced field cleaning for API compatibility

## [0.1.0] - 2024-01-XX

### Added
- Initial release
- Reference Entity synchronization
  - Entity definition sync
  - Attribute sync with field normalization
  - Record sync with metadata cleaning
- Hexagonal architecture implementation
- OAuth2 authentication with token refresh
- Pagination handling for large datasets
- Comprehensive error handling
- CI/CD pipeline with GitHub Actions
- Build system with Makefile
- golangci-lint integration
- Test coverage reporting

### Features
- `sync` command for Reference Entities
- Debug mode for troubleshooting
- JSON configuration support
- Environment variable configuration
