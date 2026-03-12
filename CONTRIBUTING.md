# Contributing to Airport Data Go

Thank you for your interest in contributing to Airport Data Go! This document explains the branching strategy, CI/CD setup, and how to make a release.

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/<your-username>/airport-data-go.git
   cd airport-data-go
   ```
3. Make sure tests pass:
   ```bash
   go test -v -race ./...
   ```

## Development Workflow

### Branch Structure

| Branch | Purpose |
|---|---|
| `main` | Active development. All new work goes here. |
| `release` | Production releases only. Merging into this branch triggers a release. |

### Making Changes

1. Create a feature branch from `main`:
   ```bash
   git checkout main
   git pull origin main
   git checkout -b feature/your-feature-name
   ```
2. Make your changes
3. Run tests locally:
   ```bash
   go test -v -race ./...
   go vet ./...
   ```
4. Commit your changes:
   ```bash
   git add .
   git commit -m "feat: description of your change"
   ```
5. Push and open a Pull Request against `main`:
   ```bash
   git push origin feature/your-feature-name
   ```

### Commit Message Convention

Use conventional commit prefixes:

- `feat:` — New feature
- `fix:` — Bug fix
- `docs:` — Documentation only
- `test:` — Adding or updating tests
- `ci:` — CI/CD changes
- `refactor:` — Code refactoring (no feature or fix)
- `chore:` — Maintenance tasks

## CI/CD

### CI (Continuous Integration)

Triggered on every push or pull request to `main`.

**What it does:**
- Runs `go mod verify`
- Runs `go vet ./...`
- Runs tests with race detection and coverage (`go test -v -race -coverprofile=coverage.out ./...`)
- Tests on Go 1.26 (minimum version required by `go.mod`)

### Release Workflow

Triggered when code is pushed to the `release` branch.

**What it does:**
1. Runs the full test suite across all Go versions
2. Reads the version from the `VERSION` file
3. Checks if a git tag for that version already exists
4. If the tag is new:
   - Creates an annotated git tag `vX.Y.Z`
   - Indexes the module on the Go module proxy (`proxy.golang.org`)
   - Creates a GitHub Release

## Making a Release

Only maintainers with push access to the `release` branch can make releases.

### Steps

1. **Ensure `main` is stable.** All tests should be passing on CI.

2. **Bump the version** in the `VERSION` file on `main`:
   ```bash
   git checkout main
   echo "1.2.0" > VERSION
   git add VERSION
   git commit -m "chore: bump version to 1.2.0"
   git push origin main
   ```

3. **Merge `main` into `release`:**
   ```bash
   git checkout release
   git merge main
   git push origin release
   ```

4. **Switch back to `main`:**
   ```bash
   git checkout main
   ```

The release workflow will automatically create the tag, GitHub Release, and index the new version on the Go module proxy. After a few minutes, users can install it with:

```bash
go get github.com/aashishvanand/airport-data-go@v1.2.0
```

### Versioning

This project follows [Semantic Versioning](https://semver.org/):

- **MAJOR** (X.0.0) — Breaking API changes
- **MINOR** (0.X.0) — New features, backward compatible
- **PATCH** (0.0.X) — Bug fixes, backward compatible

The current version is tracked in the `VERSION` file at the repository root.

## Project Structure

```
airport-data-go/
├── .github/
│   ├── dependabot.yml          # Automated dependency updates
│   └── workflows/
│       ├── ci.yml              # CI on push/PR to main
│       └── release.yml         # Release on push to release branch
├── data/
│   └── airports.json           # Embedded airport database
├── airport_data.go             # Library source code
├── airport_data_test.go        # Test suite
├── go.mod                      # Go module definition
├── VERSION                     # Current release version
├── LICENSE                     # CC BY 4.0
├── README.md                   # Documentation and API reference
└── CONTRIBUTING.md             # This file
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output and race detection
go test -v -race ./...

# Run with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Run a specific test
go test -run TestGetAirportByIata ./...
```

## Code Guidelines

- All exported functions and types must have godoc comments
- Keep zero external dependencies unless absolutely necessary
- Maintain backward compatibility within a major version
- Aim for test coverage above 85%

## Reporting Issues

If you find a bug or have a feature request, please [open an issue](https://github.com/aashishvanand/airport-data-go/issues).

## License

By contributing to this project, you agree that your contributions will be licensed under the [Creative Commons Attribution 4.0 International (CC BY 4.0)](LICENSE).
