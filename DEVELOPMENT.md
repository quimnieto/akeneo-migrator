# Development Guide

## Prerequisites

- Go 1.21 or higher
- Make (optional, but recommended)
- golangci-lint (optional, for linting)

## Quick Start

### Using Make

```bash
# Show all available commands
make help

# Install dependencies
make install

# Run tests
make test

# Build the application
make build

# Run the application
make run
```

### Without Make

```bash
# Install dependencies
go mod download
go mod tidy

# Run tests
go test ./...

# Build the application
go build -o bin/akeneo-migrator ./cmd/app

# Run the application
./bin/akeneo-migrator
```

## Development Workflow

### 1. Install Dependencies

```bash
make install
```

### 2. Run Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### 3. Format and Lint

```bash
# Format code
make fmt

# Run linter (requires golangci-lint)
make lint

# Run go vet
make vet

# Run all checks
make check
```

### 4. Build

```bash
# Build binary
make build

# Binary will be in ./bin/akeneo-migrator
```

### 5. Clean

```bash
# Remove build artifacts
make clean
```

## Testing

### Unit Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/reference_entity/syncing/...
```

### Coverage Report

```bash
# Generate HTML coverage report
make test-coverage

# Open coverage.html in your browser
```

## Linting

### Install golangci-lint

```bash
# macOS
brew install golangci-lint

# Linux
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Windows
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Run Linter

```bash
make lint
```

## CI/CD

The project uses GitHub Actions for continuous integration. The pipeline runs on:
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches

### Pipeline Steps

1. **Test** - Runs tests on Go 1.21, 1.22, and 1.23
2. **Lint** - Runs golangci-lint
3. **Build** - Builds the binary
4. **Security** - Runs security scan with gosec

### Local CI Simulation

```bash
# Run the same checks as CI
make check
make build
```

## Project Structure

```
akeneo-migrator/
├── cmd/app/              # Application entry point
├── internal/             # Private application code
│   ├── config/          # Configuration management
│   ├── reference_entity/ # Domain logic
│   └── platform/        # Infrastructure
├── kit/                 # Shared utilities
├── configs/             # Configuration files
├── .github/workflows/   # GitHub Actions
├── Makefile            # Build automation
└── go.mod              # Go modules
```

## Debugging

### Enable Debug Mode

```bash
./bin/akeneo-migrator sync brands --debug
```

### View Logs

The application outputs logs to stdout. You can redirect to a file:

```bash
./bin/akeneo-migrator sync brands 2>&1 | tee sync.log
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Format code (`make fmt`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## Troubleshooting

### Tests Failing

```bash
# Clean and rebuild
make clean
make install
make test
```

### Build Errors

```bash
# Verify Go version
go version

# Clean Go cache
go clean -cache -modcache -i -r

# Reinstall dependencies
make install
```

### Linter Issues

```bash
# Update golangci-lint
brew upgrade golangci-lint  # macOS
# or
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```
