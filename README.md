# Akeneo Migrator

CLI tool to synchronize data between Akeneo PIM instances.

## Installation

```bash
go mod tidy
go build -o akeneo-migrator ./cmd/app
```

## Configuration

### Option 1: JSON File (Recommended)

1. Copy the example configuration file:
```bash
cp configs/akeneo-migrator/settings.local.json.example configs/akeneo-migrator/settings.local.json
```

2. Edit `configs/akeneo-migrator/settings.local.json` with your credentials:
```json
{
  "akeneoSource": {
    "api": {
      "url": "https://your-source-akeneo.com",
      "credentials": {
        "clientId": "your_source_client_id",
        "secret": "your_source_secret",
        "username": "your_source_username",
        "password": "your_source_password"
      }
    }
  },
  "akeneoDest": {
    "api": {
      "url": "https://your-dest-akeneo.com",
      "credentials": {
        "clientId": "your_dest_client_id",
        "secret": "your_dest_secret",
        "username": "your_dest_username",
        "password": "your_dest_password"
      }
    }
  }
}
```

### Option 2: Environment Variables

1. Copy the example .env file:
```bash
cp .env.example .env
```

2. Configure the variables in the `.env` file

## Usage

### Synchronize a Reference Entity

```bash
./akeneo-migrator sync brands
```

This will:
1. **Synchronize the Reference Entity definition** (structure, attributes, etc.)
   - If the entity doesn't exist in destination, it will be created
   - If it exists, it will be updated with the source definition
2. **Synchronize all records** from the "brands" Reference Entity from source to destination

### Debug Mode

```bash
./akeneo-migrator sync brands --debug
```

Debug mode shows:
- Record contents before sending
- Detailed error messages
- Validation issues

## Project Structure

```
akeneo-migrator/
├── cmd/
│   └── app/
│       ├── main.go              # Entry point (only calls bootstrap)
│       └── bootstrap/           # Bootstrap with all CLI logic
├── configs/                     # Configuration files
│   └── akeneo-migrator/
│       └── settings.local.json  # Local configuration
├── internal/
│   ├── config/                  # Configuration management
│   ├── reference_entity/        # Reference Entities domain
│   │   ├── repository.go        # Repository interface
│   │   └── syncing/             # Synchronization service
│   └── platform/
│       ├── client/
│       │   └── akeneo/          # Akeneo HTTP client
│       └── storage/
│           └── akeneo/          # Repository implementation
└── kit/                         # Shared utilities
    └── config/
        └── static/
            └── viper/           # Configuration with Viper
```

## Architecture

The project follows the **Hexagonal Architecture** pattern (Ports & Adapters):

- **Domain** (`internal/reference_entity/`): Contains business logic and interfaces (ports)
- **Services** (`internal/reference_entity/syncing/`): Use cases and application logic
- **Infrastructure** (`internal/platform/`): Concrete implementations (adapters)
  - `client/akeneo/`: HTTP client for Akeneo API
  - `storage/akeneo/`: Repository implementation using the client
- **Bootstrap** (`cmd/app/bootstrap/`): Dependency injection and configuration

### Architecture Advantages:

- ✅ **Testable**: You can easily mock repositories
- ✅ **Decoupled**: Business logic doesn't depend on Akeneo directly
- ✅ **Extensible**: Easy to add new implementations (e.g., another PIM)
- ✅ **Maintainable**: Clear separation of responsibilities

## Development

To add new commands:

1. Define the interface in the domain (`internal/[domain]/repository.go`)
2. Create the service with business logic (`internal/[domain]/[action]/service.go`)
3. Implement the repository in `internal/platform/storage/`
4. Register dependencies in bootstrap
5. Create the command in bootstrap

## Testing

```bash
# Run unit tests
go test ./internal/reference_entity/syncing/...

# Run all tests
go test ./...

# Run with coverage
go test -cover ./...
```

## Common Issues

### 422 Unprocessable Entity

This error usually means:
- Invalid field format
- Missing required attributes
- Invalid attribute values

Use `--debug` flag to see detailed error messages.

### Authentication Errors

Verify your credentials in the configuration file:
- Client ID and Secret are correct
- Username and Password are valid
- API URL is accessible

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

See LICENSE file for details.
