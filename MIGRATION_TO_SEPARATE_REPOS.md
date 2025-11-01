# Migration Guide: Separate Repositories

Guide for splitting the monorepo into separate repositories.

## Current Structure (Monorepo)

```
outbox_uow/
├── outbox_uow_abstraction/
└── outbox_uow_sql/
```

## Target Structure (Multiple Repos)

```
github.com/arash/outbox-abstraction/    (Repository 1)
github.com/arash/outbox-sql/          (Repository 2)
```

## Step-by-Step Migration

### Phase 1: Prepare Abstraction Repository

#### 1. Create New Repository

```bash
# On GitHub, create: github.com/arash/outbox-abstraction
```

#### 2. Extract Abstraction Code

```bash
# Create new directory
mkdir -p ~/temp/outbox-abstraction
cd ~/temp/outbox-abstraction

# Initialize git
git init

# Copy abstraction code
cp -r /path/to/outbox_uow/outbox_uow_abstraction/* .

# Update go.mod
cat > go.mod << 'EOF'
module github.com/arash/outbox-abstraction

go 1.23

require github.com/google/uuid v1.6.0
EOF

# Copy documentation
cp /path/to/outbox_uow/ARCHITECTURE.md .
cp /path/to/outbox_uow/EXTENSION_GUIDE.md .

# Create README
cat > README.md << 'EOF'
# Outbox Pattern - Core Abstraction

Core interfaces and types for the Outbox Pattern implementation.

## Installation

```bash
go get github.com/arash/outbox-abstraction
```

## Usage

```go
import "github.com/arash/outbox-abstraction/abstraction"

// Create event
event := abstraction.CreateNewEvent(...)

// Create manager
manager := abstraction.NewOutboxEventManager()
```

## Implementations

- [PostgreSQL](https://github.com/arash/outbox-sql) - Official implementation
- [NATS](https://github.com/arash/outbox-nats) - Community implementation
- See [EXTENSION_GUIDE.md](EXTENSION_GUIDE.md) for creating your own

## Documentation

- [Architecture](ARCHITECTURE.md) - Design principles
- [Extension Guide](EXTENSION_GUIDE.md) - Creating implementations
EOF

# Add and commit
git add .
git commit -m "Initial commit: core abstraction"

# Add remote and push
git remote add origin git@github.com:arash/outbox-abstraction.git
git branch -M main
git push -u origin main
```

#### 3. Create Release Tag

```bash
cd ~/temp/outbox-abstraction
git tag v1.0.0
git push origin v1.0.0
```

### Phase 2: Create PostgreSQL Repository

#### 1. Create New Repository

```bash
# On GitHub, create: github.com/arash/outbox-sql
```

#### 2. Extract PostgreSQL Code

```bash
# Create new directory
mkdir -p ~/temp/outbox-sql
cd ~/temp/outbox-sql

# Initialize git
git init

# Copy PostgreSQL code
cp -r /path/to/outbox_uow/outbox_uow_sql/* .

# Update go.mod to reference public abstraction
cat > go.mod << 'EOF'
module github.com/arash/outbox-sql

go 1.23

require (
    github.com/arash/outbox-abstraction v1.0.0
    github.com/google/uuid v1.6.0
    gorm.io/gorm v1.31.0
)

require (
    github.com/jinzhu/inflection v1.0.0 // indirect
    github.com/jinzhu/now v1.1.5 // indirect
    golang.org/x/text v0.20.0 // indirect
)
EOF

# Copy build files
cp /path/to/outbox_uow/docker-compose.yml .
cp /path/to/outbox_uow/Dockerfile.test .
cp /path/to/outbox_uow/Makefile .
cp -r /path/to/outbox_uow/scripts .

# Copy documentation
cp /path/to/outbox_uow/QUICKSTART.md .
cp /path/to/outbox_uow/DEVELOPMENT.md .
cp /path/to/outbox_uow/TESTING.md .

# Create README
cat > README.md << 'EOF'
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
- [Core Abstraction](https://github.com/arash/outbox-abstraction)
EOF

# Download dependencies
go mod tidy

# Add and commit
git add .
git commit -m "Initial commit: PostgreSQL implementation"

# Add remote and push
git remote add origin git@github.com:arash/outbox-sql.git
git branch -M main
git push -u origin main
```

#### 3. Create Release Tag

```bash
cd ~/temp/outbox-sql
git tag v1.0.0
git push origin v1.0.0
```

### Phase 3: Update Import Paths

#### In PostgreSQL Implementation

```bash
cd ~/temp/outbox-sql

# Update import paths in all Go files
find . -name "*.go" -type f -exec sed -i '' \
  's|github.com/arash/outbox_abstraction/abstraction|github.com/arash/outbox-abstraction/abstraction|g' {} \;

# Verify
go mod tidy
go test ./...

# Commit
git add .
git commit -m "Update import paths to use public abstraction repository"
git push
```

### Phase 4: Verify Everything Works

```bash
# Test abstraction
cd ~/temp/outbox-abstraction
go test ./...

# Test PostgreSQL implementation
cd ~/temp/outbox-sql
go mod tidy
make test-unit
make docker-integration
```

## Usage After Migration

### For Users

```bash
# Install abstraction
go get github.com/arash/outbox-abstraction

# Install PostgreSQL implementation
go get github.com/arash/outbox-sql
```

### In Code

```go
package main

import (
    "github.com/arash/outbox-abstraction/abstraction"
    "github.com/arash/outbox-sql/sql_channel"
)

func main() {
    // Use abstraction
    manager := abstraction.NewOutboxEventManager()
    
    // Use PostgreSQL implementation
    channel := sqlchannel.NewsqlEventChannel(...)
    manager.RegisterEventChannel(eventType, channel)
}
```

## Benefits After Migration

### 1. Independent Versioning

```bash
# Stable abstraction
github.com/arash/outbox-abstraction v1.0.0  (rarely changes)

# Active development on implementation
github.com/arash/outbox-sql v1.0.0
github.com/arash/outbox-sql v1.1.0  (bug fix)
github.com/arash/outbox-sql v1.2.0  (optimization)
github.com/arash/outbox-sql v2.0.0  (major update)
```

### 2. Clear Dependencies

```
Users' Applications
    ↓ depends on
Abstraction (stable interface)
    ↑ implemented by
Implementations (can change)
```

### 3. Community Contributions

```
Official:
├── github.com/arash/outbox-abstraction
└── github.com/arash/outbox-sql

Community can create:
├── github.com/someone/outbox-nats
├── github.com/another/outbox-kafka
└── github.com/team/outbox-redis

All reference: github.com/arash/outbox-abstraction
```

### 4. Easier Maintenance

```
Bug in PostgreSQL?
├─ Fix in github.com/arash/outbox-sql
├─ Release v1.0.1
└─ Abstraction unchanged ✓

New feature in abstraction?
├─ Update github.com/arash/outbox-abstraction
├─ Release v2.0.0
├─ Update implementations to v2.0.0 when ready
└─ Backward compatible if designed well ✓
```

## Repository Settings

### Abstraction Repository

```yaml
Name: outbox-abstraction
Description: Core interfaces for Outbox Pattern implementation
Topics: outbox-pattern, event-sourcing, go, interface, abstraction
License: MIT

Branches:
  main: Protected, require PR reviews
  
Tags:
  v1.0.0: Initial release
```

### PostgreSQL Repository

```yaml
Name: outbox-sql
Description: PostgreSQL implementation of Outbox Pattern
Topics: outbox-pattern, postgresql, event-sourcing, go, gorm
License: MIT

Branches:
  main: Protected, require PR reviews
  
Tags:
  v1.0.0: Initial release
```

## CI/CD Configuration

### Abstraction Repository (.github/workflows/test.yml)

```yaml
name: Test
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      - run: go test -v ./...
      - run: go test -coverprofile=coverage.out ./...
      - run: go tool cover -html=coverage.out -o coverage.html
```

### PostgreSQL Repository (.github/workflows/test.yml)

```yaml
name: Test
on: [push, pull_request]
jobs:
  test:
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
      - run: make install
      - run: make test-unit
      - run: make test-integration
```

## Documentation Updates

### Abstraction README

Focus on:
- Interface definitions
- Core concepts
- How to implement
- List of known implementations

### PostgreSQL README

Focus on:
- Installation
- Configuration
- Usage examples
- Performance tips
- Docker setup

## Release Process

### Abstraction

```bash
# Rarely changes
# When it does, it's a big deal
git tag v2.0.0  # Major version bump
git push origin v2.0.0
```

### PostgreSQL Implementation

```bash
# Regular updates
git tag v1.1.0  # Minor version for features
git tag v1.0.1  # Patch version for fixes
git push origin --tags
```

## Migration Checklist

- [ ] Create abstraction repository
- [ ] Create PostgreSQL repository
- [ ] Update import paths
- [ ] Test abstraction standalone
- [ ] Test PostgreSQL with public abstraction
- [ ] Update documentation
- [ ] Create release tags
- [ ] Set up CI/CD
- [ ] Update READMEs
- [ ] Announce migration to users

## Backward Compatibility

If you need to maintain the old monorepo temporarily:

```bash
# In old monorepo, update go.mod
cd outbox_uow/outbox_uow_sql

# Replace local reference with public
module outbox_uow_sql

require (
    github.com/arash/outbox-abstraction v1.0.0  // Changed!
    // ... other deps
)

# Remove local replace directive
# replace github.com/arash/outbox_abstraction => ../outbox_uow_abstraction
```

## Summary

After migration:

✅ **Abstraction** = Stable, versioned interface  
✅ **Implementations** = Independent, evolve freely  
✅ **Community** = Can contribute implementations  
✅ **Maintenance** = Easier, isolated changes  
✅ **Professional** = Industry-standard approach  

This is the **correct way** to structure Go libraries! 🎉

