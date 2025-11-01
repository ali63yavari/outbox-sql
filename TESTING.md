# Testing Guide

This document describes the testing strategy and how to run tests for the Outbox Pattern implementation.

## Test Coverage

### Abstraction Module Tests (`outbox_uow_abstraction/abstraction`)

#### `event_test.go`
- ✅ Event creation with UUID generation
- ✅ Unique EventID generation across multiple events
- ✅ Event status string conversion
- ✅ JSON marshaling/unmarshaling of event status
- ✅ Full event JSON serialization

**Tests:**
- `TestCreateNewEvent`
- `TestCreateNewEvent_UniqueEventIDs`
- `TestOutboxEventStatus_String`
- `TestOutboxEventStatus_MarshalJSON`
- `TestOutboxEventStatus_UnmarshalJSON`
- `TestOutboxEvent_JSONMarshaling`

#### `event_manager_test.go`
- ✅ Event manager creation
- ✅ Channel registration
- ✅ Channel overwrite behavior
- ✅ Error handling for missing channels
- ✅ Multiple channel management
- ✅ Channel isolation

**Tests:**
- `TestNewOutboxEventManager`
- `TestOutboxEventManager_RegisterEventChannel`
- `TestOutboxEventManager_RegisterEventChannel_Overwrite`
- `TestOutboxEventManager_GetChannel_NotFound`
- `TestOutboxEventManager_MultipleChannels`
- `TestOutboxEventManager_ChannelIsolation`

### PostgreSQL Module Tests (`outbox_uow_pgsql/pgsql_channel`)

#### `channel_test.go` (Unit Tests)
- ✅ EventModel to/from OutboxEvent conversion
- 📝 Database integration tests (skipped, require PostgreSQL)

**Tests:**
- `TestEventModel_FromOutboxEvent`
- `TestEventModel_ToOutboxEvent`
- `TestRegisterEvent` (skipped)
- `TestRegisterEvent_MultipleEvents` (skipped)
- `TestProcessBatch_SuccessfulEvents` (skipped)
- `TestProcessBatch_DifferentEventTypes` (skipped)
- `TestNewPgSqlEventChannel_DefaultValues` (skipped)
- `TestNewPgSqlEventChannel_CustomValues` (skipped)

#### `integration_test.go` (Integration Tests)

Requires PostgreSQL database. Run with: `go test -tags=integration -v`

**Tests:**
- `TestIntegration_RegisterEvent` - Event registration and persistence
- `TestIntegration_ProcessBatch_SuccessfulEvents` - Successful event processing
- `TestIntegration_ProcessBatch_DifferentEventTypes` - Multi-channel event isolation
- `TestIntegration_FailedEvent_Retry` - Retry mechanism with eventual success
- `TestIntegration_MaxRetries_MarkAsFailed` - Failed events after max retries

## Running Tests

### Quick Test (Unit Tests Only)

```bash
# Test abstraction module
cd outbox_uow_abstraction
go test -v ./abstraction

# Test PostgreSQL module (unit tests)
cd ../outbox_uow_pgsql
go test -v ./pgsql_channel
```

### Integration Tests with PostgreSQL

#### Option 1: Using Docker

```bash
# Start PostgreSQL container
docker run -d \
  --name outbox-test-db \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=outbox_test \
  -p 5432:5432 \
  postgres:15-alpine

# Wait for PostgreSQL to be ready
sleep 3

# Run integration tests
cd outbox_uow_pgsql
go test -tags=integration -v ./pgsql_channel

# Cleanup
docker stop outbox-test-db && docker rm outbox-test-db
```

#### Option 2: Using Existing PostgreSQL

```bash
# Set environment variables
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5432
export TEST_DB_USER=postgres
export TEST_DB_PASSWORD=postgres
export TEST_DB_NAME=outbox_test

# Create test database (if needed)
createdb outbox_test

# Run integration tests
cd outbox_uow_pgsql
go test -tags=integration -v ./pgsql_channel

# Cleanup
dropdb outbox_test
```

### Test with Coverage

```bash
# Abstraction module
cd outbox_uow_abstraction
go test -coverprofile=coverage.out ./abstraction
go tool cover -html=coverage.out

# PostgreSQL module
cd ../outbox_uow_pgsql
go test -coverprofile=coverage.out ./pgsql_channel
go tool cover -html=coverage.out
```

### Run All Tests

```bash
# From project root
./run_tests.sh
```

## Test Environment Variables

For integration tests, the following environment variables can be configured:

| Variable | Default | Description |
|----------|---------|-------------|
| `TEST_DB_HOST` | localhost | PostgreSQL host |
| `TEST_DB_PORT` | 5432 | PostgreSQL port |
| `TEST_DB_USER` | postgres | Database user |
| `TEST_DB_PASSWORD` | postgres | Database password |
| `TEST_DB_NAME` | outbox_test | Database name |

## Continuous Integration

For CI/CD pipelines, use the following workflow:

```yaml
name: Tests

on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      
      - name: Test abstraction
        run: |
          cd outbox_uow_abstraction
          go test -v ./abstraction
      
      - name: Test PostgreSQL (unit)
        run: |
          cd outbox_uow_pgsql
          go test -v ./pgsql_channel

  integration-tests:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15-alpine
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: outbox_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
    
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      
      - name: Integration tests
        env:
          TEST_DB_HOST: localhost
          TEST_DB_PORT: 5432
          TEST_DB_USER: postgres
          TEST_DB_PASSWORD: postgres
          TEST_DB_NAME: outbox_test
        run: |
          cd outbox_uow_pgsql
          go test -tags=integration -v ./pgsql_channel
```

## Test Results

All unit tests pass successfully:

```
✅ Abstraction Module: 12 tests, 0 failures
✅ PostgreSQL Module: 2 tests, 6 skipped (require PostgreSQL)
```

Integration tests (with PostgreSQL):
```
✅ 5 integration tests covering:
   - Event registration
   - Batch processing
   - Multi-channel isolation
   - Retry mechanism
   - Failure handling
```

## Troubleshooting

### SQLite Compatibility Issues

The PostgreSQL implementation uses PostgreSQL-specific features (UUID, JSONB) that are not compatible with SQLite. Unit tests that require database operations are skipped. Use integration tests with PostgreSQL for full testing.

### Connection Errors

If integration tests fail to connect:
1. Verify PostgreSQL is running: `pg_isready`
2. Check connection parameters
3. Ensure test database exists
4. Verify firewall/network settings

### Permission Errors

If you get permission errors:
```bash
# Grant necessary permissions
psql -c "GRANT ALL PRIVILEGES ON DATABASE outbox_test TO postgres;"
```

## Adding New Tests

When adding new features:

1. **Unit tests**: Add to `*_test.go` files without database dependencies
2. **Integration tests**: Add to `integration_test.go` with `// +build integration` tag
3. **Update this document** with new test descriptions

