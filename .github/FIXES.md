# CI/CD Fixes Applied

## Issues Fixed

### 1. PostgreSQL UUID Extension Missing

**Problem:**
```
integration_test.go:80: Failed to migrate database: ERROR: function uuid_generate_v4() does not exist (SQLSTATE 42883)
```

The integration tests were failing because the PostgreSQL database didn't have the `uuid-ossp` extension enabled, which is required for the `uuid_generate_v4()` function used in the database schema.

**Solution:**
Added a "Setup test database" step to both CI and coverage workflows that:
1. Creates the `outbox_test` database (if it doesn't exist)
2. Enables the `uuid-ossp` extension

**Code added to workflows:**
```yaml
- name: Setup test database
  env:
    PGPASSWORD: postgres
  run: |
    psql -h localhost -U postgres -d postgres -c "CREATE DATABASE outbox_test;" || true
    psql -h localhost -U postgres -d outbox_test -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"
```

This step runs before the integration tests and ensures the database is properly configured.

### 2. Cache Restoration Errors

**Problem:**
```
Error: /usr/bin/tar: ../../../go/pkg/mod/gorm.io/driver/sqlite@v1.6.0/*.go: Cannot open: File exists
```

The manual cache configuration using `actions/cache@v4` was causing conflicts when restoring cached Go modules, resulting in "File exists" errors during tar extraction.

**Solution:**
Replaced manual cache configuration with built-in caching from `actions/setup-go@v5`:

**Before:**
```yaml
- name: Set up Go
  uses: actions/setup-go@v5
  with:
    go-version: '1.23'

- name: Cache Go modules
  uses: actions/cache@v4
  with:
    path: |
      ~/.cache/go-build
      ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
    restore-keys: |
      ${{ runner.os }}-go-
```

**After:**
```yaml
- name: Set up Go
  uses: actions/setup-go@v5
  with:
    go-version: '1.23'
    cache: true
    cache-dependency-path: go.sum
```

**Benefits:**
- ✅ No more cache conflicts
- ✅ Automatic cache key management
- ✅ Better cache invalidation
- ✅ Simpler configuration
- ✅ Applied to all workflows (ci.yml, coverage.yml, release.yml)

### 3. golangci-lint Configuration Warnings

**Problem:**
```
WARN [config_reader] The configuration option `run.skip-files` is deprecated
WARN [config_reader] The configuration option `run.skip-dirs` is deprecated
WARN [config_reader] The configuration option `output.format` is deprecated
WARN [lintersdb] The linter "varcheck" is deprecated
WARN [lintersdb] The linter "structcheck" is deprecated
WARN [lintersdb] The linter "deadcode" is deprecated
```

**Solution:**
Updated `.golangci.yml`:
- Moved `skip-files` from `run` section to `issues.exclude-files`
- Moved `skip-dirs` from `run` section to `issues.exclude-dirs`
- Changed `output.format` to `output.formats` array
- Removed deprecated linters (`varcheck`, `structcheck`, `deadcode`) from disable list

### 4. Coverage Pipeline Failure

**Problem:**
```
Error: Failed to get gist: 404 Not Found
```

The coverage badge generation step was failing because GitHub Gist secrets (`GIST_SECRET` and `GIST_ID`) were not configured.

**Solution:**
Updated workflows to make optional features fail gracefully:

1. **`.github/workflows/coverage.yml`**:
   - Added `continue-on-error: true` to "Create coverage badge" step
   - Added `continue-on-error: true` to "Upload coverage to Codecov" step
   - Pipeline now succeeds even if these optional services aren't configured

2. **`.github/workflows/ci.yml`**:
   - Added `continue-on-error: true` to "Upload coverage to Codecov" step

3. **`README.md`**:
   - Commented out coverage badge (requires setup)
   - Added instructions for enabling badges when ready
   - Provided alternative Codecov badge option

## What Works Now

### ✅ Without Any Secrets
- CI pipeline runs successfully
- Lint, test, and build jobs complete
- Coverage reports generated and uploaded as artifacts
- Coverage percentage displayed in workflow logs

### ✅ With Optional Secrets (when configured)

#### Codecov Integration (optional)
- Set `CODECOV_TOKEN` secret
- Coverage uploaded to codecov.io
- Detailed coverage analysis and trends

#### Coverage Badge (optional)
- Set `GIST_SECRET` and `GIST_ID` secrets
- Dynamic coverage badge in README
- Auto-updates on push to main

## How to Use

### Minimal Setup (No Secrets Required)
Just push your code! The CI will:
1. Run linting
2. Run all tests
3. Generate coverage reports
4. Upload coverage artifacts (downloadable from Actions tab)

### Full Setup (With Optional Features)

Follow instructions in `.github/SETUP.md` to enable:
1. **Codecov**: Detailed coverage reporting
2. **Coverage Badge**: Dynamic badge in README

## Coverage Report Access

Even without secrets, you can always:
1. Go to **Actions** tab → select workflow run
2. Scroll to **Artifacts** section
3. Download `coverage-report` (contains coverage.out and coverage.html)
4. Open `coverage.html` in your browser for detailed coverage view

Or check coverage percentage in workflow logs:
- Coverage workflow → "Extract coverage percentage" step
- Shows: `Coverage: X.X%`

## Branch Configuration

The workflows now run on:
- `main` branch
- `dev` branch (updated from `develop`)
- All pull requests to these branches

## Status

✅ All issues resolved
✅ CI pipeline working without secrets
✅ Optional features fail gracefully
✅ Coverage reports always generated
✅ golangci-lint warnings eliminated

