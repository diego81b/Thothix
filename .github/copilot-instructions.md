# Copilot Instructions for Thothix Project

This file provides coding guidelines and best practices for GitHub Copilot when working on the Thothix project.

## 🏗️ Project Architecture

### Technology Stack

- **Backend**: Go 1.21+ with Gin framework
- **Database**: PostgreSQL 17
- **Authentication**: Clerk
- **Secrets Management**: HashiCorp Vault
- **Containerization**: Docker with multi-stage builds
- **Documentation**: Markdown with links to source files

### Project Structure

```
thothix/
├── backend/                 # Go API source code
├── vault/                   # Vault configuration and scripts
├── scripts/                 # Development and deployment scripts
├── docs/                    # Additional documentation
├── Dockerfile.backend       # Backend container (multi-stage)
├── Dockerfile.postgres      # PostgreSQL container (multi-stage)
├── Dockerfile.vault        # Vault container (multi-stage)
├── docker-compose.yml       # Development environment
├── docker-compose.prod.yml  # Production overrides
└── .env.example            # Single source for all environment variables
```

### Design Patterns & Architecture

- Follow the SOLID principles.
- Prefer composition over inheritance.
- Use Dependency Injection for decoupling and testability.
- Apply appropriate design patterns such as Repository, Strategy, Factory, and Adapter where applicable.
- Follow Clean Architecture when structuring large applications.

## 🔧 Coding Standards

### Go Backend Guidelines

#### Code Structure

- Organize code in packages: `handlers`, `services`, `models`, `middleware`
- Keep handlers thin, business logic in services
- Keep methods small and focused; one method should do one thing only.
- Avoid magic numbers and strings—use constants or enums.
- Prefer modern Go features like type aliases, generics, and error wrapping.
- Use `context.Context` for request-scoped values and cancellation.
- Use `sync.WaitGroup` for concurrent operations where appropriate.
- Use `defer` for cleanup tasks (e.g., closing database connections).

#### Naming Conventions

```go
// Structs: PascalCase
type UserService struct {}

// Interfaces: PascalCase with -er suffix when possible
type UserRepository interface {}

// Functions/Methods: PascalCase for exported, camelCase for private
func (s *UserService) GetUser() {}
func (s *UserService) validateInput() {}

// Variables: camelCase
var userCount int
var maxRetries = 3

// Constants: PascalCase or SCREAMING_SNAKE_CASE for exported
const MaxUsers = 100
const DEFAULT_TIMEOUT = 30
```

#### Error Handling

```go
// Always handle errors explicitly
result, err := service.GetUser(id)
if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
}

// Use custom error types for business logic
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}
```

#### API Responses

```go
// Use consistent response structures
type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
    Message string      `json:"message,omitempty"`
}

// HTTP status codes
// 200: Successful GET, PUT, PATCH
// 201: Successful POST
// 204: Successful DELETE
// 400: Bad request/validation errors
// 401: Unauthorized
// 403: Forbidden
// 404: Not found
// 500: Internal server error
```

### Docker Guidelines

#### Multi-Stage Builds

- Always use multi-stage builds with `dev` and `prod` targets
- Development stage: Include debugging tools, hot reload
- Production stage: Minimal, secure, optimized

#### Naming Conventions

```yaml
# Development
container_name: thothix-service-dev
image: thothix/service:version-dev

# Production
container_name: thothix-service-prod
image: thothix/service:version-prod
```

#### Health Checks

```dockerfile
# Always include health checks
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1
```

## 📝 Documentation Standards

### Principle: Single Source of Truth

- **Never hardcode configurations** in documentation
- **Always reference source files** instead of copying content
- Use relative links to actual files in the repository
- use of English
- link official documentation for external tools (e.g., Docker, Go, Vault)
- Use emojis for clarity and visual appeal
- Use consistent formatting and structure
- Use clear, concise language
- Use headings and bullet points for readability
- avoid duplication of information

#### Good Examples

```markdown
📋 **Configuration**: See [.env.example](./.env.example) for all options
🐳 **Docker Setup**: Check [docker-compose.yml](./docker-compose.yml)
🔧 **Build Process**: View [Dockerfile.backend](./Dockerfile.backend)
```

#### Bad Examples (Don't Do This)

````markdown
# ❌ Don't hardcode configurations

```yaml
version: '3.8'
services:
  app:
    build: .
    # ... 50 lines of config
```
````

### Documentation Structure

- Keep main README.md focused and concise
- Use dedicated files for complex topics (VAULT_INTEGRATION.md, etc.)
- Include table of contents for long documents
- Use consistent emoji and formatting

#### File Naming

- Use SCREAMING_SNAKE_CASE for important docs: `README.md`, `VAULT_INTEGRATION.md`
- Use kebab-case for specific guides: `docker-modernization.md`
- Always use `.md` extension

## 🔐 Security Guidelines

- Never trust client-side validation—always validate on the server.
- Avoid exposing sensitive logic or secrets
- Use HTTPS for all API endpoints
- Implement rate limiting and IP whitelisting where applicable
- Use secure headers (CORS, CSP, etc.)
- Regularly review and update dependencies for security patches
- Use environment variables for sensitive configurations
- Use Vault for production secrets management
- Use Clerk for authentication and authorization
- Regularly audit code for security vulnerabilities
- least privilege principle: only give permissions that are necessary

### Environment Variables

- Never commit real secrets to repository
- Use `.env.example` as template with fake/example values
- Use descriptive comments for each variable
- Group related variables with headers

```bash
# === DATABASE CONFIGURATION ===
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=thothix_dev
POSTGRES_USER=thothix_user
POSTGRES_PASSWORD=change_me_in_production

# === VAULT CONFIGURATION ===
USE_VAULT=true
VAULT_ADDR=http://vault:8200
VAULT_ROOT_TOKEN=thothix-secure-root-token-2025-v1
```

### Vault Integration

- Use Vault for all production secrets
- Keep development secrets simple but not hardcoded
- Document vault paths and policies clearly
- Use appropriate token policies (least privilege)

## 🛠️ Development Workflow

### Script Development Guidelines

- **Single Version Policy**: Create only ONE version of each script per functionality
- **Windows Priority**: Prefer `.bat` files for Windows environments (primary development OS)
- **Cross-Platform**: Only create multiple versions if truly necessary for cross-platform support
- **Naming Convention**: Use descriptive names (e.g., `version-bump.bat`, `dev.bat`, `setup-hooks.ps1`)
- **Documentation**: Each script must be documented in `scripts/README.md`
- **Functionality**: Each script should have a single, clear purpose
- **Error Handling**: Include proper error handling and user feedback

### Git Practices

- Use conventional commits: `feat:`, `fix:`, `docs:`, `refactor:`
- Keep commits focused and atomic
- Update documentation in the same commit as code changes

### Testing

- Write unit tests for all business logic
- Use table-driven tests in Go
- Mock external dependencies (database, vault, etc.)
- Aim for >80% code coverage

### Local Development

```bash
# Standard workflow
cp .env.example .env
# Edit .env with your local values
docker compose up -d --build
```

## 🚀 Deployment Guidelines

### Environment Separation

- **Development**: Local docker-compose with dev mode services
- **Staging**: Production-like setup with test data
- **Production**: Secure, optimized, monitored

### Configuration Management

- Development: Use `.env` file (git-ignored)
- Staging: Use `.env.staging` (git-ignored)
- Production: Use `.env.prod` (git-ignored) or external secret management

## 📊 Monitoring and Logging

### Logging Standards

```go
// Use structured logging
log.Info("User created successfully",
    "user_id", user.ID,
    "email", user.Email,
    "timestamp", time.Now())

// Log levels
// DEBUG: Detailed information for debugging
// INFO: General information about application flow
// WARN: Warning conditions that don't stop the application
// ERROR: Error conditions that might stop specific operations
// FATAL: Critical errors that stop the application
```

### Health Checks

- Implement `/health` endpoint for basic service status
- Implement `/readiness` endpoint for dependency checks
- Include version information in health responses

## 🔄 Maintenance Guidelines

### Dependency Management

- Keep dependencies up to date
- Use specific versions in production
- Document any version constraints or compatibility issues

### Documentation Maintenance

- Review documentation with every major feature
- Ensure links still work (use relative paths)
- Update examples when APIs change
- Remove outdated information promptly

## 🎯 Best Practices Summary

1. **Code Quality**: Clean, well-tested, documented
2. **Security**: Vault integration, no hardcoded secrets
3. **Documentation**: Reference source files, single source of truth
4. **Docker**: Multi-stage builds, consistent naming
5. **Development**: Local-first, containerized, reproducible
6. **Deployment**: Environment-specific configurations
7. **Maintenance**: Keep everything up to date and working

---

```markdown
## 📝 Git Commit Best Practices

### Conventional Commits

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification for all commit messages:
```

<type>[optional scope]: <description>

[optional body]

[optional footer(s)]

````

#### Commit Types

- **feat**: A new feature for the user
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **build**: Changes that affect the build system or external dependencies
- **ci**: Changes to CI configuration files and scripts
- **chore**: Other changes that don't modify src or test files
- **revert**: Reverts a previous commit

#### Examples

```bash
# Good commit messages
feat: add user authentication with Clerk integration
fix: resolve database connection timeout in production
docs: update README with new environment variables
refactor: extract user validation logic into separate service
perf: optimize database queries for user retrieval
test: add unit tests for user service validation
build: update Go dependencies to latest versions
ci: add Docker build step to GitHub Actions
chore: update .gitignore for IDE files

# Bad commit messages (avoid these)
fix: stuff
update: changes
wip
asdf
````

### Commit Message Guidelines

#### Subject Line (First Line)

- **Length**: Keep under 50 characters
- **Style**: Use imperative mood ("add" not "added" or "adds")
- **Capitalization**: Start with lowercase after the type
- **Punctuation**: No period at the end
- **Language**: Use English only

#### Body (Optional)

- **Purpose**: Explain what and why, not how
- **Length**: Wrap lines at 72 characters
- **Details**: Include motivation for the change and contrast with previous behavior

#### Footer (Optional)

- **Breaking Changes**: Start with `BREAKING CHANGE:`
- **Issues**: Reference issues with `Closes #123` or `Fixes #456`
- **Co-authors**: Use `Co-authored-by: name <email>`

### Advanced Commit Examples

#### Simple Feature Addition

```bash
feat: implement user role management endpoints

Add CRUD operations for user roles with proper validation
and authorization middleware integration.
```

#### Bug Fix with Context

```bash
fix: resolve Vault connection timeout in development

The Vault health check was failing on slower machines due to
insufficient timeout values. Increased timeouts and improved
error handling for better development experience.

Fixes #123
```

#### Breaking Change

```bash
feat!: migrate authentication from JWT to Clerk

BREAKING CHANGE: JWT authentication is no longer supported.
All API clients must now use Clerk tokens for authentication.

Migration guide available in AUTHENTICATION_MIGRATION.md
```

#### Documentation Update

```bash
docs: simplify configuration references

Replace hardcoded configuration examples with links to source
files to maintain single source of truth principle.

- Remove duplicate .env examples from README
- Add references to .env.example
- Update Docker documentation links
```

### Commit Frequency and Size

#### When to Commit

- **Logical Units**: Each commit should represent one logical change
- **Compilable**: Code should compile and pass basic tests
- **Atomic**: Changes should be self-contained and reversible
- **Frequent**: Commit often to avoid losing work

#### Commit Size Guidelines

- **Small Changes**: Prefer small, focused commits
- **Single Responsibility**: One commit = one purpose
- **Related Changes**: Group related changes together
- **Separate Concerns**: Split unrelated changes into separate commits

#### Examples of Good Commit Grouping

```bash
# Good: Separate commits for different concerns
feat: add user validation middleware
test: add tests for user validation
docs: update API documentation for user endpoints

# Bad: Mixed concerns in one commit
feat: add user validation, fix typo, update README
```

### Branch and Merge Practices

#### Branch Naming

```bash
# Feature branches
feat/user-authentication
feat/vault-integration
feat/docker-modernization

# Bug fix branches
fix/database-timeout
fix/auth-middleware

# Documentation branches
docs/api-documentation
docs/deployment-guide

# Maintenance branches
chore/dependency-updates
refactor/clean-architecture
```

#### Merge Strategies

- **Squash and Merge**: For feature branches with many commits
- **Merge Commit**: For important milestones and releases
- **Rebase and Merge**: For clean, linear history when appropriate

### Pre-Commit Checklist

Before committing, ensure:

- [ ] **Code compiles** without errors
- [ ] **Tests pass** (run `go test ./...`)
- [ ] **Linting passes** (if pre-commit hooks are configured)
- [ ] **Documentation updated** if API or configuration changes
- [ ] **Environment variables** added to `.env.example` if new ones introduced
- [ ] **No secrets** or sensitive data in the commit
- [ ] **Related files updated** (Docker files, configs, etc.)
- [ ] **CHANGELOG.md updated** with commit details (see CHANGELOG Guidelines below)

### CHANGELOG Guidelines

#### Update Process

Always update `CHANGELOG.md` before committing significant changes:

1. **Add entry under "Unreleased"** section
2. **Use same format as commit message** but with more details
3. **Include impact and benefits** for users/developers
4. **Link to related issues/PRs** when applicable

#### CHANGELOG Entry Format

```markdown
## [Unreleased]

### Added

- feat: new feature description with benefits and usage details
- feat: another feature with comprehensive explanation

### Changed

- refactor: refactoring description with impact on existing functionality
- perf: performance improvement details with metrics if available

### Fixed

- fix: bug fix description with root cause and resolution details
- fix: security fix description (without exposing vulnerability details)

### Documentation

- docs: documentation update description with scope and improvements

### Infrastructure

- build: build system changes with impact on development workflow
- ci: CI/CD improvements with deployment and testing benefits
```

#### CHANGELOG Best Practices

- **User-Focused**: Write for end users and developers using the project
- **Detailed Context**: Include more context than commit messages
- **Impact Description**: Explain how changes affect users/developers
- **Migration Notes**: Include breaking changes and migration steps
- **Links**: Reference issues, PRs, and documentation when helpful
- **Consistent Format**: Use format `v1.2.3 - Description (YYYY-MM-DD)` for releases
- **Recent First**: Keep most recent versions at the top

#### Semantic Versioning Guidelines

Follow [Semantic Versioning 2.0.0](https://semver.org/) (MAJOR.MINOR.PATCH):

- **MAJOR** (1.0.0 → 2.0.0): Breaking changes, incompatible API changes
- **MINOR** (1.0.0 → 1.1.0): New features, backward compatible functionality
- **PATCH** (1.0.0 → 1.0.1): Bug fixes, backward compatible fixes

#### Version Bump Decision Guide

Based on commit type, determine version increment:

- **feat**: MINOR version (new functionality)
- **fix**: PATCH version (bug fixes)
- **BREAKING CHANGE**: MAJOR version (breaking changes)
- **perf**: PATCH version (performance improvements)
- **refactor**: PATCH version (code refactoring)
- **docs**: No version bump (documentation only)
- **style**: No version bump (formatting, no code changes)
- **test**: No version bump (tests only)
- **build/ci**: PATCH version if affects production builds
- **chore**: No version bump (maintenance tasks)

#### Example CHANGELOG Entry

```markdown
### Added

- feat: comprehensive Copilot instructions with Git best practices
  - Added complete development guidelines in `.github/copilot-instructions.md`
  - Includes Go coding standards with clean architecture patterns
  - Provides Docker multi-stage build conventions and naming standards
  - Documents security guidelines for environment variables and Vault integration
  - Establishes documentation standards following single source of truth principle
  - Defines comprehensive Git commit best practices with conventional commits
  - Includes pre-commit checklist and Git hooks examples
  - Provides branch naming conventions and collaborative workflows
  - **Impact**: GitHub Copilot now has complete context for maintaining code quality,
    security standards, and consistent development practices across the project
```

#### Release Process

When creating a release:

1. **Determine version bump** based on changes in [Unreleased] section
2. **Move entries** from "Unreleased" to new version section with format: `v1.2.3 - Description (YYYY-MM-DD)`
3. **Create Git tag** with semantic versioning
4. **Generate release notes** from CHANGELOG entries

```markdown
## v1.2.0 - Enhanced Development Workflow (2025-06-15)

### Added

- feat: comprehensive Copilot instructions with Git best practices
  [... detailed description ...]

### Changed

- refactor: Docker configuration modernization with multi-stage builds
  [... detailed description ...]
```

- [ ] **CHANGELOG.md updated** with commit details (see CHANGELOG Guidelines below)

### CHANGELOG Guidelines

#### Update Process

Always update `CHANGELOG.md` before committing significant changes:

1. **Add entry under "Unreleased"** section
2. **Use same format as commit message** but with more details
3. **Include impact and benefits** for users/developers
4. **Link to related issues/PRs** when applicable

### Git Hooks Integration

#### Pre-commit Hook Example

```bash
#!/bin/sh
# Run tests before commit
go test ./...
if [ $? -ne 0 ]; then
    echo "Tests failed. Commit aborted."
    exit 1
fi

# Run linting
golangci-lint run
if [ $? -ne 0 ]; then
    echo "Linting failed. Commit aborted."
    exit 1
fi
```

#### Commit Message Hook

```bash
#!/bin/sh
# Validate commit message format
commit_regex='^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\(.+\))?: .{1,50}'

if ! grep -qE "$commit_regex" "$1"; then
    echo "Invalid commit message format. Please use conventional commits."
    echo "Example: feat: add user authentication"
    exit 1
fi
```

### Working with Multiple Contributors

#### Collaborative Commits

```bash
# When pair programming or collaborating
git commit -m "feat: implement user service

Co-authored-by: John Doe <john@example.com>
Co-authored-by: Jane Smith <jane@example.com>"
```

#### Reviewing and Amending

```bash
# Amend last commit (only if not pushed)
git commit --amend -m "feat: implement user service with proper validation"

# Interactive rebase for multiple commits (only if not pushed)
git rebase -i HEAD~3
```

### Troubleshooting Common Issues

#### Fixing Commit Messages

```bash
# Fix last commit message (not yet pushed)
git commit --amend -m "fix: correct typo in user validation"

# Fix older commit messages (not yet pushed)
git rebase -i HEAD~3  # Then edit the commits you want to change
```

#### Handling Large Changes

```bash
# Split large changes into smaller commits
git add -p  # Interactive staging
git commit -m "feat: add user model structure"
git add -p  # Stage more changes
git commit -m "feat: add user validation logic"
```

#### Emergency Fixes

```bash
# Hotfix workflow
git checkout main
git checkout -b hotfix/critical-security-fix
# Make minimal fix
git commit -m "fix: patch security vulnerability in auth middleware

CVE-2024-XXXX: Improper input validation could lead to
unauthorized access. Added proper sanitization.

Closes #urgent-security-issue"
```

### Commit History Maintenance

#### Keeping Clean History

- Use `git rebase` instead of `git merge` for feature branches when possible
- Squash related commits before merging to main
- Use descriptive commit messages for better project archaeology
- Remove or squash "WIP" and "fix typo" commits before merging

#### Release Management

```bash
# Tag releases with semantic versioning
git tag -a v1.2.0 -m "Release version 1.2.0

- Added user authentication with Clerk
- Implemented Vault integration
- Updated Docker configuration
- Various bug fixes and improvements"
```

Remember: **Good commits tell a story** of your project's evolution and make it easier for team members (including future you) to understand the codebase changes.

```

---

_Last Updated: When making changes to this file, update this timestamp and document major changes in commit messages._
```
