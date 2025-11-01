# Quick Start Guide

Get up and running with the Outbox Pattern implementation in 5 minutes.

## Prerequisites

- Go 1.23+
- Docker & Docker Compose
- Make

## Installation

```bash
# Clone the repository
cd /path/to/outbox_uow

# Install dependencies
make install
```

## Running Tests

### 1. Unit Tests (No Database Required)

```bash
make test-unit
```

Expected output:

```
✓ 12 tests in abstraction module
✓ 2 tests in PostgreSQL module
```

### 2. Integration Tests (With Docker)

```bash
make docker-integration
```

This will:

1. ✓ Start PostgreSQL in Docker
2. ✓ Wait for database to be ready
3. ✓ Run 5 integration tests
4. ✓ Keep PostgreSQL running

**Stop PostgreSQL:**

```bash
make docker-down
```

### 3. All Tests at Once

```bash
make docker-test
```

## Common Commands

| What You Want         | Command                   |
| --------------------- | ------------------------- |
| See all commands      | `make help`               |
| Install dependencies  | `make install`            |
| Run unit tests        | `make test-unit`          |
| Run integration tests | `make docker-integration` |
| Start PostgreSQL      | `make docker-up`          |
| Stop PostgreSQL       | `make docker-down`        |
| Format code           | `make fmt`                |
| Check code quality    | `make verify`             |
| Clean up              | `make clean`              |

## Basic Usage Example

```go
package main

import (
    "context"
    "time"
    "github.com/ali63yavari/outbox-abstraction/abstraction"
    "outbox_uow_sql/sql_channel"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// Define event type
type UserCreated struct{}
func (e UserCreated) GetName() string { return "UserCreated" }
func (e UserCreated) GetID() int      { return 1 }

func main() {
    // Connect to database
    dsn := "host=localhost user=postgres password=postgres dbname=mydb"
    db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

    // Create event handler
    handler := func(ctx context.Context, event *abstraction.OutboxEvent) error {
        println("Processing:", event.EventID)
        return nil
    }

    // Configure channel
    maxRetries := 3
    batchSize := 10
    interval := 60 * time.Second

    // Create channel
    channel := sqlchannel.NewsqlEventChannel(
        context.Background(),
        db,
        UserCreated{},
        &maxRetries,
        &batchSize,
        &interval,
        handler,
    )

    // Register event
    event := abstraction.CreateNewEvent(
        UserCreated{},
        "User",
        "user-123",
        map[string]interface{}{"email": "user@example.com"},
    )

    channel.RegisterEvent(event)
}
```

## Troubleshooting

### Port 5432 Already in Use

```bash
# Stop local PostgreSQL
sudo systemctl stop postgresql  # Linux
brew services stop postgresql   # macOS

# Or use different port
TEST_DB_PORT=5433 make docker-up
```

### Tests Failing

```bash
# Clean everything and try again
make clean
make docker-clean
make install
make docker-integration
```

### Docker Issues

```bash
# Check Docker is running
docker ps

# Restart Docker services
make docker-down
make docker-up

# View logs
make docker-logs
```

## Next Steps

- Read [README.md](README.md) for detailed usage
- Check [TESTING.md](TESTING.md) for testing guide
- See [DEVELOPMENT.md](DEVELOPMENT.md) for development workflow
- Explore `make help` for all commands

## Getting Help

```bash
# Show configuration
make info

# Show all commands
make help

# Check versions
make version
```

## Quick Verification

Run this to verify everything works:

```bash
make verify
```

This runs:

- Code formatting
- Static analysis
- Build check
- Unit tests

If all pass, you're ready to develop! 🚀
