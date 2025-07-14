# GitHub Actions - CI/CD Documentation

## Prerequisites

### Local Development Environment
- **Node.js**: v22.13.0+ (matches CI environment)
- **Go**: v1.24.3+ (matches CI environment)
- **npm**: Latest version compatible with Node 22

### Environment Alignment
The CI/CD workflows are configured to match your local environment:
- All workflows use Node.js v22
- All workflows use Go v1.24.3
- Package.json requires Node >=22.0.0

## Available Workflows

### 1. Basic Checks (`basic-checks.yml`)
**Trigger**: Push/PR to `main` or `dev`
**Purpose**: Fast feedback for basic validation
- Go tests with PostgreSQL
- Go linting with golangci-lint
- npm scripts validation

### 2. CI/CD Pipeline (`ci.yml`)
**Trigger**: Push/PR to `main` or `dev`
**Purpose**: Complete testing and building
- Backend tests + build
- Docker image building
- Security scanning
- Docker Compose integration test

### 3. Auto Version Bump (`auto-version-bump.yml`)
**Trigger**: PR merged to `main`
**Purpose**: Automated versioning
- Bumps package.json version
- Creates PR with version update
- Extracts feature authors from git log

### 4. PR Merge Checks (`pr-merge-checks.yml`)
**Trigger**: PR merged to `main` or `dev`
**Purpose**: Post-merge validation
- Runs full pre-commit checks
- Ensures code quality after merge

### 5. Sync Main to Dev (`sync-main-to-dev.yml`)
**Trigger**: Push to `main`
**Purpose**: Keep dev branch updated
- Creates PR from main to dev
- Automated branch synchronization

### 6. Debug Workflow (`debug.yml`)
**Trigger**: Manual dispatch
**Purpose**: Troubleshooting CI issues
- Three debug levels: basic, detailed, full
- Environment inspection
- Dependency validation

## Troubleshooting

### Local Testing
Use existing npm scripts to validate your environment:
```bash
# Test core functionality
npm run pre-commit

# Individual checks
npm run format
npm run lint
npm run test

# Full validation
npm run pre-commit:full
```

### CI Debugging
1. Go to Actions tab in GitHub
2. Run "Debug Workflow" manually
3. Select appropriate debug level
4. Review output for issues

### Common Issues

#### Node.js Version Mismatch
- **Problem**: CI fails with Node.js compatibility errors
- **Solution**: Update local Node.js to v22.13.0+

#### Go Version Mismatch
- **Problem**: Go build/test failures in CI
- **Solution**: Update local Go to v1.24.3+

#### npm Script Failures
- **Problem**: npm commands timeout or fail
- **Solution**: Run `npm run pre-commit:full` locally to identify issues

#### Docker Compose Issues
- **Problem**: Health check failures in CI
- **Solution**: Check service logs and ensure all containers start properly

## Best Practices

1. **Keep versions aligned**: Local environment should match CI
2. **Test locally first**: Use `npm run pre-commit` before pushing
3. **Use debug workflow**: For persistent CI issues
4. **Check logs**: Always review failure logs in Actions tab

## Maintenance

When updating versions:
1. Update all workflow files consistently
2. Update package.json engines field
3. Update this documentation
4. Test locally with npm scripts
