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

### PostgreSQL Module Tests (`outbox_uow_sql/sql_channel`)

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
- `TestNewsqlEventChannel_DefaultValues` (skipped)
- `TestNewsqlEventChannel_CustomValues` (skipped)

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
cd ../outbox_uow_sql
go test -v ./sql_channel
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
cd outbox_uow_sql
go test -tags=integration -v ./sql_channel

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
cd outbox_uow_sql
go test -tags=integration -v ./sql_channel

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
cd ../outbox_uow_sql
go test -coverprofile=coverage.out ./sql_channel
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

This repository uses GitHub Actions for automated testing and coverage reporting.

### Workflows

1. **CI Workflow** (`.github/workflows/ci.yml`)
   - Runs on every push and pull request
   - Jobs: Lint → Test → Build
   - Includes unit and integration tests
   - Generates coverage reports
   - Uploads coverage to Codecov

2. **Coverage Workflow** (`.github/workflows/coverage.yml`)
   - Runs on push to `main` branch
   - Generates coverage badge
   - Updates coverage statistics

3. **Release Workflow** (`.github/workflows/release.yml`)
   - Runs when a version tag is pushed
   - Creates GitHub releases automatically
   - Generates changelog

### CI Test Commands

The CI pipeline uses the following test commands:

```bash
# Unit tests
go test ./sqlchannel -run 'Test[^I]' -short -v

# Integration tests (with PostgreSQL service)
go test -v ./sqlchannel -run 'TestI'

# Coverage
go test -coverprofile=coverage.out -covermode=atomic ./sqlchannel
```

### Setting Up CI

See [GitHub Actions Setup Guide](.github/SETUP.md) for detailed instructions on:
- Required secrets (Codecov, coverage badges)
- Customizing workflows
- Troubleshooting CI issues

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

