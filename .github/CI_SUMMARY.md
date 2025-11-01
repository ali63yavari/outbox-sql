# CI/CD Summary

This document provides a quick overview of the CI/CD setup for this repository.

## Files Created

### GitHub Actions Workflows

1. **`.github/workflows/ci.yml`**
   - Main CI pipeline
   - Runs: Lint, Test (unit + integration), Build
   - Triggers: Push and Pull Requests to `main`/`develop`
   - Uploads coverage to Codecov and GitHub artifacts

2. **`.github/workflows/coverage.yml`**
   - Coverage reporting and badge generation
   - Triggers: Push to `main` branch only
   - Generates dynamic coverage badge via GitHub Gist

3. **`.github/workflows/release.yml`**
   - Automated release creation
   - Triggers: Push of version tags (e.g., `v1.0.0`)
   - Creates GitHub release with changelog

### Configuration Files

4. **`.golangci.yml`**
   - golangci-lint configuration
   - Defines linting rules and settings
   - Used by both local development and CI

### Documentation

5. **`.github/SETUP.md`**
   - Detailed setup instructions
   - Secret configuration guide
   - Troubleshooting tips

6. **`.github/CI_SUMMARY.md`** (this file)
   - Quick reference for CI/CD setup

### Updated Files

7. **`README.md`**
   - Added CI/CD badges:
     - Build status
     - Coverage percentage
     - Go Report Card
     - Go Reference
     - License

8. **`TESTING.md`**
   - Updated with CI/CD information
   - Links to setup guide

## Quick Start

### 1. Push to GitHub

```bash
# Add all files
git add .github/ .golangci.yml README.md TESTING.md

# Commit
git commit -m "Add GitHub Actions CI/CD pipeline"

# Push
git push origin main
```

### 2. Verify Workflows

1. Go to your repository on GitHub
2. Click **Actions** tab
3. You should see workflows running

### 3. Optional: Enable Coverage Badge

Follow the instructions in `.github/SETUP.md` to set up:
- Codecov token
- GitHub Gist for coverage badge

### 4. Create a Release

```bash
# Tag a release
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

This will automatically:
- Run tests
- Create a GitHub release
- Generate changelog

## Features

### ✅ Automated Testing
- Unit tests on every push/PR
- Integration tests with PostgreSQL
- Test results visible in Actions tab

### ✅ Code Quality
- Automated linting with golangci-lint
- Code formatting checks (gofmt)
- Static analysis (go vet)

### ✅ Coverage Reporting
- Coverage reports uploaded to Codecov
- Coverage artifacts downloadable from Actions
- Dynamic coverage badge in README

### ✅ Build Verification
- Ensures code builds successfully
- Verifies no uncommitted changes after build

### ✅ Automated Releases
- Create releases by pushing tags
- Automatic changelog generation
- Release notes with installation instructions

## Badge URLs

Update these in `README.md` after setting up:

### CI Badge
```markdown
[![CI](https://github.com/ali63yavari/outbox-sql/workflows/CI/badge.svg)](https://github.com/ali63yavari/outbox-sql/actions/workflows/ci.yml)
```

### Coverage Badge (requires Gist setup)
```markdown
[![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/ali63yavari/YOUR_GIST_ID/raw/outbox-coverage.json)](https://github.com/ali63yavari/outbox-sql/actions/workflows/coverage.yml)
```

### Codecov Badge (alternative)
```markdown
[![codecov](https://codecov.io/gh/ali63yavari/outbox-sql/branch/main/graph/badge.svg)](https://codecov.io/gh/ali63yavari/outbox-sql)
```

## Secrets Required

| Secret | Required | Purpose |
|--------|----------|---------|
| `CODECOV_TOKEN` | Optional | Upload coverage to Codecov |
| `GIST_SECRET` | Optional | Update coverage badge |
| `GIST_ID` | Optional | Coverage badge location |
| `GITHUB_TOKEN` | Auto | GitHub releases (auto-provided) |

## Workflow Triggers

| Workflow | Trigger | Branches |
|----------|---------|----------|
| CI | Push, Pull Request | `main`, `develop` |
| Coverage | Push | `main` only |
| Release | Tag Push | `v*` tags |

## Test Strategy

The CI pipeline uses test name patterns to separate unit and integration tests:

- **Unit Tests**: `Test[^I]` - All tests NOT starting with "TestI"
- **Integration Tests**: `TestI` - Tests starting with "TestI"

### Example Test Names

```go
// Unit test - runs in CI without database
func TestEventModel_ToOutboxEvent(t *testing.T) { ... }

// Integration test - runs with PostgreSQL service
func TestIntegration_RegisterEvent(t *testing.T) { ... }
```

## Local Testing Before Push

```bash
# Run what CI will run
make lint
make test-unit
make docker-integration  # Requires Docker
make build
```

## Monitoring

### View Test Results
1. Go to **Actions** tab
2. Click on a workflow run
3. Expand job steps to see logs

### Download Coverage Report
1. Go to completed workflow run
2. Scroll to **Artifacts** section
3. Download `coverage-report`

### Check Coverage Trends
- Visit Codecov dashboard (after setup)
- View coverage over time
- See coverage per file/function

## Next Steps

1. ✅ Push workflows to GitHub
2. ⬜ Set up Codecov (optional)
3. ⬜ Set up coverage badge (optional)
4. ⬜ Create first release tag
5. ⬜ Add branch protection rules (recommended)

## Branch Protection (Recommended)

To require CI to pass before merging:

1. Go to **Settings** → **Branches**
2. Add rule for `main` branch
3. Enable:
   - ✅ Require status checks to pass
   - ✅ Require branches to be up to date
   - Select: `Lint`, `Test`, `Build`
4. Save

This prevents merging code that fails tests or linting.

## Support

For issues or questions:
- Check `.github/SETUP.md` for detailed setup
- Review workflow logs in Actions tab
- Consult [GitHub Actions Documentation](https://docs.github.com/en/actions)

