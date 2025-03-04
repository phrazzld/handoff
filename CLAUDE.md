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