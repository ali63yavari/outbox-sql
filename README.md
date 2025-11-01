# Outbox Pattern with Unit of Work

A Go implementation of the Outbox Pattern with Unit of Work for reliable event publishing in distributed systems.

## 📚 Documentation

- **[QUICKSTART.md](QUICKSTART.md)** - Get started in 5 minutes
- **[ARCHITECTURE.md](ARCHITECTURE.md)** - Architecture and design principles ⭐
- **[EXTENSION_GUIDE.md](EXTENSION_GUIDE.md)** - Creating custom implementations
- **[DEVELOPMENT.md](DEVELOPMENT.md)** - Complete developer guide
- **[TESTING.md](TESTING.md)** - Testing strategies and CI/CD
- **[PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md)** - Project structure
- **[SUMMARY.md](SUMMARY.md)** - Implementation summary

## 🚀 Quick Start

```bash
# Install dependencies
make install

# Run unit tests (no database required)
make test-unit

# Run integration tests with Docker PostgreSQL
make docker-integration

# Show all available commands
make help
```

## Architecture

The project consists of two main modules:

### 1. Abstraction Module (`outbox_uow_abstraction`)

Defines the core interfaces and types:
- `OutboxEvent`: Event structure with metadata
- `OutboxEventType`: Interface for event type definitions
- `OutboxEventManager`: Manages multiple event channels
- `OutboxEventChannel`: Interface for channel implementations

### 2. PostgreSQL Implementation (`outbox_uow_pgsql`)

PostgreSQL-backed implementation with:
- Automatic table migration
- Batch processing with configurable intervals
- Retry mechanism with exponential backoff
- Row-level locking with `SKIP LOCKED`
- Per-event transaction handling
- Support for multiple event types on a single table

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd outbox_uow

# Install dependencies
cd outbox_uow_abstraction && go mod tidy
cd ../outbox_uow_pgsql && go mod tidy
```

## Usage

### Basic Example

```go
package main

import (
    "context"
    "time"

    "github.com/arash/outbox_abstraction/abstraction"
    "outbox_uow_pgsql/pgsql_channel"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// Define your event type
type UserCreatedEvent struct{}

func (e UserCreatedEvent) GetName() string { return "UserCreated" }
func (e UserCreatedEvent) GetID() int      { return 1 }

func main() {
    // Setup database
    dsn := "host=localhost user=postgres password=postgres dbname=mydb sslmode=disable"
    db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

    // Create context
    ctx := context.Background()

    // Define event handler
    handler := func(ctx context.Context, event *abstraction.OutboxEvent) error {
        // Publish to message broker, call external API, etc.
        println("Processing event:", event.EventID)
        return nil
    }

    // Configure channel
    maxRetries := 3
    batchSize := 10
    interval := 60 * time.Second
    eventType := UserCreatedEvent{}

    // Create channel
    channel := pgsqlchannel.NewPgSqlEventChannel(
        ctx,
        db,
        eventType,
        &maxRetries,
        &batchSize,
        &interval,
        handler,
    )

    // Create and register event
    event := abstraction.CreateNewEvent(
        eventType,
        "User",
        "user-123",
        map[string]interface{}{
            "email": "user@example.com",
            "name":  "John Doe",
        },
    )

    channel.RegisterEvent(event)
}
```

### Event Manager with Multiple Channels

```go
// Create event manager
manager := abstraction.NewOutboxEventManager()

// Register multiple channels for different event types
manager.RegisterEventChannel(UserCreatedEvent{}, userCreatedChannel)
manager.RegisterEventChannel(UserUpdatedEvent{}, userUpdatedChannel)
manager.RegisterEventChannel(OrderPlacedEvent{}, orderPlacedChannel)

// Get channel and register event
channel, _ := manager.GetChannel(UserCreatedEvent{})
channel.RegisterEvent(event)
```

## Running Tests

### Quick Start with Make

```bash
# Install dependencies
make install

# Run unit tests only (fastest, no database required)
make test-unit

# Run integration tests with Docker (recommended)
make docker-integration

# Run all tests
make docker-test

# Show all available commands
make help
```

### Unit Tests

```bash
# Using Make (recommended)
make test-unit

# Or directly
cd outbox_uow_abstraction && go test -v ./abstraction
cd outbox_uow_pgsql && go test -v ./pgsql_channel
```

### Integration Tests

Integration tests require PostgreSQL. The easiest way is using Docker Compose:

```bash
# Option 1: Using Make (recommended)
make docker-integration

# Option 2: Manual Docker Compose
docker-compose up -d postgres  # Start PostgreSQL
make test-integration          # Run integration tests
docker-compose down            # Stop PostgreSQL

# Option 3: Fully isolated Docker environment
make docker-test-full
```

For more testing options and troubleshooting, see [TESTING.md](TESTING.md) and [DEVELOPMENT.md](DEVELOPMENT.md).

## Configuration

### Channel Parameters

- **maxRetries**: Maximum number of retry attempts (default: 3)
- **batchSize**: Number of events to process per batch (default: 10)
- **interval**: Polling interval for new events (default: 60 seconds)

### Event Status

- `pending`: Event is waiting to be processed
- `closed`: Event processed successfully
- `failed`: Event failed after max retries

## Features

- ✅ Automatic EventID generation (UUID)
- ✅ Retry mechanism with exponential backoff
- ✅ Per-event transaction handling
- ✅ Multiple event types on single table
- ✅ Row-level locking to prevent conflicts
- ✅ Configurable batch processing
- ✅ Graceful shutdown with context cancellation
- ✅ Comprehensive test coverage

## Database Schema

The PostgreSQL implementation creates the following table:

```sql
CREATE TABLE event_models (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id VARCHAR NOT NULL,
    event_type VARCHAR,
    event_status VARCHAR,
    aggregate_type VARCHAR,
    aggregate_id VARCHAR,
    payload JSONB,
    attempt_count INTEGER,
    last_error TEXT,
    next_attempt_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_event_type ON event_models(event_type);
CREATE INDEX idx_event_status ON event_models(event_status);
CREATE INDEX idx_aggregate_type ON event_models(aggregate_type);
CREATE INDEX idx_aggregate_id ON event_models(aggregate_id);
```

## License

MIT

