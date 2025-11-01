# Outbox Pattern - PostgreSQL Implementation

PostgreSQL implementation of the Outbox Pattern.

## Installation

```bash
go get github.com/arash/outbox-sql
```

## Quick Start

```go
import (
    "github.com/arash/outbox-abstraction/abstraction"
    "github.com/arash/outbox-sql/sql_channel"
)

// Create channel
channel := sqlchannel.NewsqlEventChannel(...)

// Register with manager
manager := abstraction.NewOutboxEventManager()
manager.RegisterEventChannel(eventType, channel)
```

## Testing

```bash
# Unit tests
make test-unit

# Integration tests with Docker
make docker-integration
```

## Documentation

- [Quick Start](QUICKSTART.md)
- [Development Guide](DEVELOPMENT.md)
- [Testing Guide](TESTING.md)
- [Core Abstraction](https://github.com/ali63yavari/outbox-abstraction)