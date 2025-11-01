# Development Guide

This guide covers the development workflow, testing, and Docker integration for the Outbox Pattern implementation.

## Prerequisites

- **Go 1.23+** - [Download](https://golang.org/dl/)
- **Docker** - [Download](https://www.docker.com/get-started)
- **Docker Compose** - Usually included with Docker Desktop
- **Make** - Pre-installed on macOS/Linux, [Windows Guide](https://gnuwin32.sourceforge.net/packages/make.htm)
- **PostgreSQL Client** (optional) - For database interaction

## Quick Start

### 1. Install Dependencies

```bash
make install
```

### 2. Run Unit Tests

```bash
make test-unit
```

### 3. Run Integration Tests

Start PostgreSQL and run integration tests:

```bash
make docker-integration
```

This command will:
1. Start PostgreSQL in Docker
2. Wait for it to be ready
3. Run integration tests
4. Keep PostgreSQL running (use `make docker-down` to stop)

## Makefile Commands

Run `make help` to see all available commands.

### Essential Commands

| Command | Description |
|---------|-------------|
| `make help` | Show all available commands |
| `make install` | Install Go dependencies |
| `make test-unit` | Run unit tests only (no database required) |
| `make test-integration` | Run integration tests (requires PostgreSQL) |
| `make test-all` | Run all tests (unit + integration) |
| `make docker-up` | Start PostgreSQL with Docker Compose |
| `make docker-down` | Stop and remove Docker containers |
| `make docker-integration` | Start PostgreSQL + run integration tests |

### Docker Commands

| Command | Description |
|---------|-------------|
| `make docker-up` | Start PostgreSQL container |
| `make docker-down` | Stop PostgreSQL container |
| `make docker-clean` | Stop containers and remove volumes |
| `make docker-test` | Start PostgreSQL and run all tests |
| `make docker-integration` | Start PostgreSQL and run integration tests |
| `make docker-test-full` | Run tests in isolated Docker container |
| `make docker-logs` | Show container logs |
| `make docker-ps` | Show running containers |

### Database Commands

| Command | Description |
|---------|-------------|
| `make db-shell` | Open PostgreSQL shell |
| `make db-reset` | Reset the test database |

### Code Quality Commands

| Command | Description |
|---------|-------------|
| `make fmt` | Format Go code |
| `make vet` | Run go vet |
| `make lint` | Run golangci-lint |
| `make build` | Build all modules |
| `make coverage` | Generate coverage reports |
| `make verify` | Run fmt, vet, build, and unit tests |

### Workflow Commands

| Command | Description |
|---------|-------------|
| `make quick` | Quick unit test run |
| `make ci` | Full CI pipeline (install, verify, test-all) |
| `make all` | Run everything from scratch |
| `make clean` | Clean build artifacts |
| `make info` | Show configuration |

## Development Workflows

### Standard Development Workflow

```bash
# 1. Install dependencies
make install

# 2. Make your changes
# ... edit code ...

# 3. Format code
make fmt

# 4. Run checks
make vet

# 5. Run unit tests
make test-unit

# 6. Run integration tests (starts Docker automatically)
make docker-integration

# 7. Check coverage
make coverage
```

### Quick Test Loop

```bash
# Quick unit tests (fastest)
make quick

# Or just
make test
```

### Full Test Suite

```bash
# Option 1: Using Make (recommended)
make docker-integration

# Option 2: Manual Docker setup
make docker-up              # Start PostgreSQL
make test-integration       # Run integration tests
make docker-down           # Stop PostgreSQL

# Option 3: Fully isolated Docker environment
make docker-test-full      # Everything runs in Docker
```

### Working with the Database

```bash
# Start PostgreSQL
make docker-up

# Open psql shell
make db-shell

# In the shell, you can:
outbox_test=# \dt                    # List tables
outbox_test=# SELECT * FROM event_models;  # Query events
outbox_test=# \q                     # Exit

# Reset database
make db-reset

# Stop PostgreSQL
make docker-down
```

## Docker Compose

The project includes a `docker-compose.yml` with two services:

### Services

1. **postgres** - PostgreSQL 15 database
   - Port: 5432
   - Database: outbox_test
   - Username: postgres
   - Password: postgres

2. **test-runner** - Go test environment
   - Runs integration tests automatically
   - Connected to postgres service

### Manual Docker Compose Usage

```bash
# Start only PostgreSQL
docker-compose up -d postgres

# Start all services (PostgreSQL + test runner)
docker-compose up

# Run in background
docker-compose up -d

# View logs
docker-compose logs -f postgres

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v

# Rebuild and run tests
docker-compose up --build test-runner
```

## Environment Configuration

### Environment Variables

Create a `.env` file (based on `.env.example`):

```bash
cp .env.example .env
```

Default values:
```env
TEST_DB_HOST=localhost
TEST_DB_PORT=5432
TEST_DB_USER=postgres
TEST_DB_PASSWORD=postgres
TEST_DB_NAME=outbox_test
```

### Override Variables

```bash
# Override for a single command
TEST_DB_PORT=5433 make test-integration

# Export for session
export TEST_DB_HOST=myhost.com
make test-integration
```

## Integration Test Details

Integration tests are tagged with `// +build integration` and require:

1. Running PostgreSQL instance
2. Proper connection configuration
3. Test database permissions

### What Integration Tests Cover

- ✅ Event registration and persistence
- ✅ Batch processing with real database
- ✅ Multi-channel event isolation
- ✅ Retry mechanism with exponential backoff
- ✅ Failed event handling after max retries
- ✅ Transaction handling
- ✅ Database locking (SKIP LOCKED)

### Running Specific Integration Tests

```bash
# Start PostgreSQL
make docker-up

# Run specific test
cd outbox_uow_sql
go test -tags=integration -v ./sql_channel -run TestIntegration_RegisterEvent

# Run with timeout
go test -tags=integration -v -timeout 30s ./sql_channel
```

## Continuous Integration

### GitHub Actions Example

```yaml
name: CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      
      - name: Install dependencies
        run: make install
      
      - name: Run unit tests
        run: make test-unit
      
      - name: Start PostgreSQL
        run: make docker-up
      
      - name: Run integration tests
        run: make test-integration
      
      - name: Stop PostgreSQL
        run: make docker-down
```

### GitLab CI Example

```yaml
stages:
  - test

variables:
  TEST_DB_HOST: postgres
  TEST_DB_PORT: "5432"
  TEST_DB_USER: postgres
  TEST_DB_PASSWORD: postgres
  TEST_DB_NAME: outbox_test

test:
  image: golang:1.23
  stage: test
  services:
    - postgres:15-alpine
  variables:
    POSTGRES_DB: outbox_test
    POSTGRES_USER: postgres
    POSTGRES_PASSWORD: postgres
  script:
    - make install
    - make test-unit
    - make test-integration
```

## Troubleshooting

### PostgreSQL Connection Issues

```bash
# Check if PostgreSQL is running
make docker-ps

# Check PostgreSQL logs
make docker-logs

# Test connection manually
docker-compose exec postgres pg_isready -U postgres

# Restart PostgreSQL
make docker-down
make docker-up
```

### Port Already in Use

If port 5432 is already in use:

```bash
# Option 1: Stop existing PostgreSQL
sudo systemctl stop postgresql  # Linux
brew services stop postgresql   # macOS

# Option 2: Use different port
TEST_DB_PORT=5433 make docker-up
```

### Permission Issues

```bash
# Reset Docker volumes
make docker-clean

# Rebuild containers
docker-compose build --no-cache
```

### Test Failures

```bash
# Clean test cache
make clean

# Reset database
make db-reset

# Run tests with verbose output
cd outbox_uow_sql
go test -tags=integration -v ./sql_channel
```

### Docker Issues

```bash
# Clean everything
make docker-clean
docker system prune -a

# Rebuild from scratch
make docker-test-full
```

## Code Quality

### Pre-commit Checklist

```bash
# Format code
make fmt

# Run static checks
make vet

# Run linter (if available)
make lint

# Run tests
make test-unit

# Check integration
make docker-integration

# Verify everything
make verify
```

### Coverage Reports

```bash
# Generate coverage for all tests
make coverage

# Open coverage report in browser
open outbox_uow_abstraction/coverage.html
open outbox_uow_sql/coverage.html

# Integration test coverage
make coverage-integration
```

## Performance Testing

### Benchmark Tests

```bash
# Run benchmarks
cd outbox_uow_abstraction
go test -bench=. -benchmem ./abstraction

cd ../outbox_uow_sql
go test -tags=integration -bench=. -benchmem ./sql_channel
```

### Load Testing

Start PostgreSQL and create a load test script:

```bash
make docker-up

# Example load test
cd outbox_uow_sql
go test -tags=integration -v ./sql_channel -run TestLoad
```

## Best Practices

1. **Always run unit tests first** - They're faster and catch most issues
2. **Use `make verify`** before committing - Ensures code quality
3. **Run integration tests locally** - Before pushing changes
4. **Clean up Docker resources** - Use `make docker-clean` periodically
5. **Check test coverage** - Aim for >80% coverage
6. **Use meaningful commit messages** - Follow conventional commits
7. **Keep dependencies updated** - Run `make tidy` regularly

## Getting Help

```bash
# Show all Make commands
make help

# Show configuration
make info

# Show version information
make version
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run `make verify`
5. Run `make docker-integration`
6. Submit a pull request

For more details, see [TESTING.md](TESTING.md) and [README.md](README.md).

