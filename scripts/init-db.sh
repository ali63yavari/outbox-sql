#!/bin/bash
set -e

# Initialize PostgreSQL for outbox tests

echo "Initializing test database..."

# Create the test database if it doesn't exist
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    SELECT 'CREATE DATABASE outbox_test'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'outbox_test')\gexec

    \c outbox_test

    -- Enable UUID extension
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

    -- Grant necessary permissions
    GRANT ALL PRIVILEGES ON DATABASE outbox_test TO $POSTGRES_USER;
EOSQL

echo "Database initialization completed!"

