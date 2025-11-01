# Outbox Pattern - PostgreSQL Implementation

[![CI](https://github.com/ali63yavari/outbox-sql/workflows/CI/badge.svg)](https://github.com/ali63yavari/outbox-sql/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ali63yavari/outbox-sql)](https://goreportcard.com/report/github.com/ali63yavari/outbox-sql)
[![Go Reference](https://pkg.go.dev/badge/github.com/ali63yavari/outbox-sql.svg)](https://pkg.go.dev/github.com/ali63yavari/outbox-sql)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

<!-- 
To enable coverage badge, follow instructions in .github/SETUP.md to set up GIST_SECRET and GIST_ID, then uncomment:
[![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/USERNAME/GIST_ID/raw/outbox-coverage.json)](https://github.com/ali63yavari/outbox-sql/actions/workflows/coverage.yml)

Or use Codecov badge after setting up CODECOV_TOKEN:
[![codecov](https://codecov.io/gh/ali63yavari/outbox-sql/branch/main/graph/badge.svg)](https://codecov.io/gh/ali63yavari/outbox-sql)
-->

PostgreSQL implementation of the Outbox Pattern.

## Installation

```bash
go get github.com/ali63yavari/outbox-sql
```

## Quick Start

```go
import (
    "github.com/ali63yavari/outbox-abstraction/abstraction"
    "github.com/ali63yavari/outbox-sql/sql_channel"
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