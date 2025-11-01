# GitHub Actions Setup Guide

This document explains how to set up the GitHub Actions workflows for this repository.

## Workflows Overview

This repository includes three GitHub Actions workflows:

### 1. **CI Workflow** (`.github/workflows/ci.yml`)
- **Triggers**: Push to `main`/`develop` branches, Pull Requests
- **Jobs**:
  - **Lint**: Runs `gofmt`, `go vet`, and `golangci-lint`
  - **Test**: Runs unit and integration tests with PostgreSQL
  - **Build**: Verifies the code builds successfully
- **Coverage**: Uploads coverage reports to Codecov and as artifacts

### 2. **Coverage Workflow** (`.github/workflows/coverage.yml`)
- **Triggers**: Push to `main` branch only
- **Jobs**:
  - Generates test coverage reports
  - Creates a dynamic coverage badge
  - Uploads coverage to Codecov

## Required Secrets

To enable all features, you need to set up the following GitHub secrets:

### 1. Codecov Token (Optional but Recommended)

Codecov provides detailed coverage reports and trends.

**Steps:**
1. Go to [codecov.io](https://codecov.io/) and sign in with GitHub
2. Add your repository
3. Copy the upload token
4. In your GitHub repository, go to **Settings** → **Secrets and variables** → **Actions**
5. Click **New repository secret**
6. Name: `CODECOV_TOKEN`
7. Value: Paste the token from Codecov

### 2. Coverage Badge Secrets (Optional)

To enable dynamic coverage badges in your README:

**Steps:**

1. **Create a GitHub Gist:**
   - Go to https://gist.github.com/
   - Create a new **public** gist
   - Filename: `outbox-coverage.json`
   - Content: `{"schemaVersion": 1, "label": "coverage", "message": "0%", "color": "red"}`
   - Click **Create public gist**
   - Copy the Gist ID from the URL (e.g., `https://gist.github.com/username/GIST_ID`)

2. **Create a Personal Access Token:**
   - Go to GitHub **Settings** → **Developer settings** → **Personal access tokens** → **Tokens (classic)**
   - Click **Generate new token** → **Generate new token (classic)**
   - Name: `Gist Token for Coverage Badge`
   - Select scope: **gist** (only)
   - Click **Generate token**
   - Copy the token

3. **Add Secrets to Repository:**
   - Go to your repository **Settings** → **Secrets and variables** → **Actions**
   - Add two secrets:
     - `GIST_SECRET`: The personal access token you just created
     - `GIST_ID`: The Gist ID from step 1

4. **Update README Badge URL:**
   - In `README.md`, replace `GIST_ID` in the coverage badge URL with your actual Gist ID:
   ```markdown
   [![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/ali63yavari/YOUR_GIST_ID/raw/outbox-coverage.json)]
   ```

## Testing the Workflows

### Locally Test Before Push

```bash
# Run linting
make lint

# Run unit tests
make test-unit

# Run integration tests (requires Docker)
make docker-integration

# Run all checks
make verify
```

### After Push

1. Go to your repository on GitHub
2. Click the **Actions** tab
3. You should see your workflows running
4. Click on a workflow run to see detailed logs

## Workflow Customization

### Changing Go Version

Edit the `go-version` in workflow files:

```yaml
- name: Set up Go
  uses: actions/setup-go@v5
  with:
    go-version: '1.23'  # Change this
```

### Changing PostgreSQL Version

Edit the `postgres` service in workflow files:

```yaml
services:
  postgres:
    image: postgres:16-alpine  # Change this
```

### Adding More Branches

Edit the trigger branches in workflow files:

```yaml
on:
  push:
    branches: [ main, develop, staging ]  # Add more branches
  pull_request:
    branches: [ main, develop, staging ]
```

## Troubleshooting

### Coverage Badge Not Updating

1. Check that `GIST_SECRET` and `GIST_ID` are correctly set
2. Verify the Gist is public
3. Check the workflow logs for errors
4. The badge updates only on push to `main` branch

### Tests Failing in CI but Passing Locally

1. Check Go version matches
2. Verify PostgreSQL service is healthy
3. Check environment variables are set correctly
4. Look at the full test output in Actions logs

### Codecov Upload Failing

1. The workflow is set to `fail_ci_if_error: false`, so it won't fail the build
2. Check if `CODECOV_TOKEN` is set correctly
3. Verify repository is added to Codecov

## Badge Status

Once set up, your README will display:

- **CI Badge**: Shows if the latest build passed
- **Coverage Badge**: Shows current test coverage percentage
- **Go Report Card**: Shows code quality score
- **Go Reference**: Links to package documentation
- **License Badge**: Shows the project license

## Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Codecov Documentation](https://docs.codecov.io/)
- [Dynamic Badges Action](https://github.com/schneegans/dynamic-badges-action)

