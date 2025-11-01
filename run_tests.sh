#!/bin/bash

set -e

echo "========================================="
echo "Running Outbox Pattern Tests"
echo "========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test abstraction module
echo "${GREEN}Testing Abstraction Module...${NC}"
cd outbox_uow_abstraction
go test -v ./abstraction
echo ""

# Test PostgreSQL module (unit tests)
echo "${GREEN}Testing PostgreSQL Module (Unit Tests)...${NC}"
cd ../outbox_uow_sql
go test -v ./sql_channel
echo ""

# Check if integration tests should run
if [ "$RUN_INTEGRATION" = "true" ]; then
    echo "${GREEN}Running Integration Tests...${NC}"
    echo "${YELLOW}Note: Requires PostgreSQL to be running${NC}"
    echo ""
    
    # Check if PostgreSQL is available
    if command -v pg_isready &> /dev/null; then
        if pg_isready -h ${TEST_DB_HOST:-localhost} -p ${TEST_DB_PORT:-5432} &> /dev/null; then
            go test -tags=integration -v ./sql_channel
        else
            echo "${YELLOW}PostgreSQL is not ready. Skipping integration tests.${NC}"
            echo "To run integration tests, start PostgreSQL and set RUN_INTEGRATION=true"
        fi
    else
        echo "${YELLOW}pg_isready not found. Skipping integration tests.${NC}"
        echo "To run integration tests, install PostgreSQL and set RUN_INTEGRATION=true"
    fi
else
    echo "${YELLOW}Integration tests skipped.${NC}"
    echo "To run integration tests: RUN_INTEGRATION=true ./run_tests.sh"
fi

echo ""
echo "${GREEN}=========================================${NC}"
echo "${GREEN}All Tests Completed Successfully!${NC}"
echo "${GREEN}=========================================${NC}"

