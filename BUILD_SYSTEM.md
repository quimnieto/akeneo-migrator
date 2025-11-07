# Build System Documentation

## Overview

This project uses Make for build automation and GitHub Actions for CI/CD.

## Makefile

### Available Targets

| Target | Description |
|--------|-------------|
| `help` | Display help message with all available targets |
| `build` | Build the application binary to `./bin/akeneo-migrator` |
| `test` | Run all tests |
| `test-coverage` | Run tests with coverage report (generates `coverage.html`) |
| `clean` | Remove build artifacts and coverage files |
| `run` | Build and run the application |
| `install` | Install and tidy dependencies |
| `deps` | Download dependencies |
| `lint` | Run golangci-lint (requires installation) |
| `fmt` | Format code with `go fmt` |
| `vet` | Run `go vet` |
| `check` | Run all checks (fmt, vet, test) |

### Usage Examples

```bash
# Build the project
make build

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format and check code
make check

# Clean build artifacts
make clean
```

## GitHub Actions CI/CD

### Workflow: `.github/workflows/ci.yml`

The CI pipeline runs on:
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches

### Jobs

#### 1. Test Job
- **Matrix**: Tests on Go 1.21, 1.22, and 1.23
- **Steps**:
  - Checkout code
  - Set up Go
  - Cache Go modules
  - Download dependencies
  - Verify dependencies
  - Run tests with race detector and coverage
  - Upload coverage to Codecov

#### 2. Lint Job
- **Go Version**: 1.23
- **Steps**:
  - Checkout code
  - Set up Go
  - Run golangci-lint with 5-minute timeout

#### 3. Build Job
- **Depends on**: test, lint
- **Steps**:
  - Checkout code
  - Set up Go
  - Build binary with optimizations (`-ldflags="-s -w"`)
  - Upload binary as artifact (7-day retention)

#### 4. Security Job
- **Steps**:
  - Checkout code
  - Run Gosec security scanner
  - Upload SARIF results to GitHub Security

### Caching

The pipeline caches:
- Go build cache (`~/.cache/go-build`)
- Go modules (`~/go/pkg/mod`)

Cache key: `${{ runner.os }}-go-${{ matrix.go-version }}-${{ hashFiles('**/go.sum') }}`

## Linting Configuration

### `.golangci.yml`

Enabled linters:
- `errcheck` - Check for unchecked errors
- `gosimple` - Simplify code
- `govet` - Vet examines Go source code
- `ineffassign` - Detect ineffectual assignments
- `staticcheck` - Static analysis
- `unused` - Check for unused code
- `gofmt` - Check code formatting
- `goimports` - Check import formatting
- `misspell` - Check for misspelled words
- `unconvert` - Remove unnecessary type conversions
- `unparam` - Check for unused function parameters
- `gosec` - Security checks
- `gocritic` - Comprehensive Go linter

### Settings

- Timeout: 5 minutes
- Tests: Included in linting
- Locale: US English
- Shadow checking: Enabled

## Build Artifacts

### Generated Files

| File | Description | Gitignored |
|------|-------------|------------|
| `bin/akeneo-migrator` | Compiled binary | ✅ Yes |
| `coverage.out` | Coverage data | ✅ Yes |
| `coverage.html` | HTML coverage report | ✅ Yes |
| `*.test` | Test binaries | ✅ Yes |

### Binary Optimization

The build uses these flags:
- `-s`: Omit symbol table
- `-w`: Omit DWARF symbol table

This reduces binary size significantly.

## Local Development

### Prerequisites

```bash
# Install Go (1.21+)
# macOS
brew install go

# Install golangci-lint (optional)
brew install golangci-lint

# Install Make (usually pre-installed on Unix systems)
```

### First Time Setup

```bash
# Clone repository
git clone https://github.com/YOUR_USERNAME/akeneo-migrator.git
cd akeneo-migrator

# Install dependencies
make install

# Run tests
make test

# Build
make build
```

### Development Cycle

```bash
# 1. Make changes to code

# 2. Format code
make fmt

# 3. Run tests
make test

# 4. Run linter (if installed)
make lint

# 5. Build
make build

# 6. Test manually
./bin/akeneo-migrator sync test --debug
```

## Continuous Integration

### Pull Request Checks

When you create a PR, the following checks run:
1. ✅ Tests on 3 Go versions
2. ✅ Linting
3. ✅ Build verification
4. ✅ Security scan

All checks must pass before merging.

### Branch Protection

Recommended branch protection rules for `main`:
- Require pull request reviews
- Require status checks to pass
- Require branches to be up to date
- Include administrators

## Troubleshooting

### Make Command Not Found

```bash
# macOS
brew install make

# Ubuntu/Debian
sudo apt-get install make

# Or use Go commands directly
go build -o bin/akeneo-migrator ./cmd/app
```

### Linter Not Found

```bash
# Install golangci-lint
make lint  # Will show installation instructions

# Or install manually
brew install golangci-lint  # macOS
```

### Tests Failing in CI

1. Check if tests pass locally: `make test`
2. Check Go version compatibility
3. Review CI logs in GitHub Actions
4. Ensure all dependencies are in `go.mod`

### Build Failing

```bash
# Clean everything
make clean

# Reinstall dependencies
make install

# Try building again
make build
```

## Performance

### Build Times

Typical build times on GitHub Actions:
- Test job: ~2-3 minutes per Go version
- Lint job: ~1-2 minutes
- Build job: ~1 minute
- Security job: ~1-2 minutes

Total pipeline time: ~5-8 minutes

### Optimization Tips

1. **Use caching**: Already configured in CI
2. **Parallel jobs**: Test matrix runs in parallel
3. **Minimal dependencies**: Keep `go.mod` clean
4. **Fast tests**: Write efficient unit tests

## Future Enhancements

Potential improvements:
- [ ] Add release workflow with GoReleaser
- [ ] Add Docker build
- [ ] Add integration tests
- [ ] Add benchmark tests
- [ ] Add code quality badges
- [ ] Add automatic changelog generation
