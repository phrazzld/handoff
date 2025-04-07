# Handoff Project Guidelines

## Build Commands
- Create module: `go mod init github.com/phrazzld/handoff`
- Build: `go build`
- Run: `./handoff [options] [path1] [path2] ...`
- Run single file: `./handoff path/to/file.go`
- Test all: `go test ./...`
- Test specific: `go test -run TestName`
- Verbose tests: `go test -v ./...`
- Format code: `gofmt -w .`
- Lint: `golangci-lint run` (if installed)

## Code Style Guidelines
- **Imports**: Standard library imports first, then third-party packages, separated by blank lines
- **Formatting**: Follow standard Go formatting with `gofmt`
- **Documentation**: All exported functions and types must have comments
- **Error Handling**: Error values should be checked explicitly, use named returns for clarity
- **Naming**: Use CamelCase for exported names, camelCase for unexported
- **Types**: Use concrete types over interfaces when functions only need specific behavior
- **File Organization**: Group related functionality, keep files under 500 lines
- **Error Messages**: Lowercase, no trailing punctuation, descriptive but concise

## Commits
- Use conventional commit messages (`feat:`, `fix:`, `docs:`, `chore:`) to clearly communicate intent
- Make atomic, semantically meaningful commits that encapsulate exactly one logical change

## Testing Principles
- Prioritize high test coverage with unit, integration, and end-to-end tests
- Write deterministic, repeatable, and efficient tests
- Consider Test-Driven Development (TDD) where feasible

## Documentation Practices
- Document the **why** behind design decisions
- Keep documentation close to the codebase in markdown
- Update documentation as part of feature completion
- Balance high-level architectural overviews with practical guides

## Architecture & Design
- Embrace modularity and loose coupling for maintainability
- Separate infrastructure and business logic clearly
- Design for resilience with proper error handling
- Default to explicit error handling with meaningful messages

## Security Practices
- Assume all inputs could be hostile; build secure defaults
- Regularly check dependencies for vulnerabilities
- Adhere to the principle of least privilege
- Use encryption in transit and at rest

## Performance Considerations
- Establish and maintain performance benchmarks
- Monitor critical performance metrics
- Design for horizontal scalability
- Include load and stress testing in development workflow

## Continuous Improvement
- Regularly review code for improvement opportunities
- Actively manage and reduce technical debt
- Foster an environment where team members can freely suggest changes