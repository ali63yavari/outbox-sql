#!/bin/bash

set -e

echo "========================================="
echo "Running Outbox Pattern Tests"
echo "========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Test sqlchannel module (unit tests)
echo -e "${GREEN}Testing sqlchannel (Unit Tests)...${NC}"
go test ./sqlchannel -run 'Test[^I]' -short
echo ""

# Check if integration tests should run
if [ "$RUN_INTEGRATION" = "true" ]; then
    echo -e "${GREEN}Running Integration Tests...${NC}"
    echo -e "${YELLOW}Note: Requires PostgreSQL to be running${NC}"
    echo ""
    
    # Check if PostgreSQL is available
    if command -v pg_isready &> /dev/null; then
        if pg_isready -h ${TEST_DB_HOST:-localhost} -p ${TEST_DB_PORT:-5432} -U ${TEST_DB_USER:-postgres} &> /dev/null; then
            go test -v ./sqlchannel -run 'TestI'
        else
            echo -e "${YELLOW}PostgreSQL is not ready. Skipping integration tests.${NC}"
            echo "To run integration tests, start PostgreSQL and set RUN_INTEGRATION=true"
        fi
    else
        echo -e "${YELLOW}pg_isready not found. Skipping integration tests.${NC}"
        echo "To run integration tests, install PostgreSQL and set RUN_INTEGRATION=true"
    fi
else
    echo -e "${YELLOW}Integration tests skipped.${NC}"
    echo "To run integration tests: RUN_INTEGRATION=true ./run_tests.sh"
fi

echo ""
echo -e "${GREEN}=========================================${NC}"
echo -e "${GREEN}All Tests Completed Successfully!${NC}"
echo -e "${GREEN}=========================================${NC}"

