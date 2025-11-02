# Contributing to DKL Email Service

Thank you for considering contributing to the DKL Email Service! This document provides guidelines and instructions for contributing.

## 📋 Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Code Structure](#code-structure)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Git Commit Messages](#git-commit-messages)

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 14+
- Redis (optional, for rate limiting)
- Docker & Docker Compose (for local development)

### Development Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Jeffreasy/dklemailservice.git
   cd dklemailservice
   ```

2. **Copy environment file:**
   ```bash
   cp .env.example .env
   ```

3. **Configure environment variables:**
   Edit `.env` with your local settings

4. **Start development environment:**
   ```bash
   docker-compose -f docker-compose.dev.yml up -d
   ```

5. **Run the application:**
   ```bash
   go run main.go
   ```

## 📁 Code Structure

```
├── handlers/     # HTTP request handlers
├── services/     # Business logic layer
├── repository/   # Data access layer
├── models/       # Data structures
├── config/       # Configuration
├── database/     # Migrations & scripts
└── tests/        # Test files
```

See [ARCHITECTURE.md](ARCHITECTURE.md) for detailed architecture documentation.

## 📝 Coding Standards

### General Guidelines

1. **Follow Go conventions:**
   - Use `gofmt` for formatting
   - Follow [Effective Go](https://golang.org/doc/effective_go)
   - Use meaningful variable names

2. **Package organization:**
   - One concern per file
   - Keep files under 500 lines
   - Group related functionality

3. **Documentation:**
   - Add comments for exported functions
   - Document complex logic
   - Update relevant docs

### File Naming

- Use snake_case: `user_handler.go`
- Test files: `user_handler_test.go`
- Interfaces: `interfaces.go` or `{package}_interfaces.go`

### Code Examples

#### Handler Structure
```go
func (h *Handler) HandleRequest(c *fiber.Ctx) error {
    // 1. Parse & validate input
    // 2. Call service layer
    // 3. Format & return response
}
```

#### Service Structure
```go
func (s *Service) BusinessLogic(params) (result, error) {
    // 1. Validate business rules
    // 2. Call repository
    // 3. Transform data
    return result, nil
}
```

#### Repository Structure
```go
func (r *Repository) Create(ctx context.Context, entity *Model) error {
    // 1. Build query
    // 2. Execute
    // 3. Handle errors
}
```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests in a specific package
go test ./handlers

# Run specific test
go test -run TestFunctionName ./package
```

### Writing Tests

1. **Location:** Place tests in `tests/` directory
2. **Naming:** `{package}_{functionality}_test.go`
3. **Structure:**
   ```go
   func TestFeatureName(t *testing.T) {
       // Arrange
       // Act
       // Assert
   }
   ```

4. **Coverage:** Aim for >80% coverage on business logic

### Test Helpers

Use test helpers from `tests/test_helper.go`:
- `SetupTestDB()` - Initialize test database
- `SeedTestData()` - Add test data
- `CleanupTestDB()` - Clean up after tests

## 🔄 Pull Request Process

### Before Submitting

1. **Update your branch:**
   ```bash
   git pull origin master
   git rebase master
   ```

2. **Run tests:**
   ```bash
   go test ./...
   ```

3. **Format code:**
   ```bash
   gofmt -w .
   ```

4. **Check for issues:**
   ```bash
   go vet ./...
   ```

### PR Guidelines

1. **Clear description:**
   - What changed
   - Why it changed
   - How to test

2. **Small, focused changes:**
   - One feature/fix per PR
   - Keep PRs under 500 lines

3. **Documentation:**
   - Update relevant docs
   - Add code comments
   - Update CHANGELOG

4. **Tests:**
   - Add tests for new features
   - Verify existing tests pass
   - Update test documentation

### PR Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
How to test these changes

## Checklist
- [ ] Tests pass
- [ ] Code formatted
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
```

## 💾 Git Commit Messages

### Format

```
<emoji> <type>: <subject>

<body>

<footer>
```

### Types

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation
- `style:` Formatting
- `refactor:` Code restructuring
- `test:` Tests
- `chore:` Maintenance

### Emojis

- 🎨 `:art:` - Improve structure
- ✨ `:sparkles:` - New feature
- 🐛 `:bug:` - Fix bug
- 📝 `:memo:` - Documentation
- 🚀 `:rocket:` - Performance
- 🔒 `:lock:` - Security
- ♻️ `:recycle:` - Refactor
- 🧹 `:broom:` - Cleanup

### Examples

```bash
✨ feat: Add user profile endpoint

- Added GET /api/users/:id/profile
- Includes avatar URL from Cloudinary
- Requires authentication

Closes #123
```

```bash
🐛 fix: Correct email template rendering

Fixed issue where variables weren't being replaced
in newsletter templates.

Fixes #456
```

## 🔒 Security

### Reporting Vulnerabilities

**DO NOT** open public issues for security vulnerabilities.

Instead, email: security@dekoninklijkeloop.nl

### Security Guidelines

1. **Never commit secrets:**
   - Use environment variables
   - Keep `.env` in `.gitignore`
   - Use secret management tools

2. **Validate input:**
   - Sanitize user input
   - Validate data types
   - Check boundaries

3. **Authentication:**
   - Use JWT tokens
   - Implement refresh tokens
   - Secure password hashing

## 📚 Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Fiber Framework](https://docs.gofiber.io/)
- [PostgreSQL Docs](https://www.postgresql.org/docs/)
- [Architecture Guide](ARCHITECTURE.md)

## 💬 Questions?

- Check existing [Issues](https://github.com/Jeffreasy/dklemailservice/issues)
- Create a new issue
- Contact: dev@dekoninklijkeloop.nl

## 📄 License

By contributing, you agree that your contributions will be licensed under the project's license.

---

**Thank you for contributing!** 🎉