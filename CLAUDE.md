# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands
- Build: `go build`
- Install: `go install`
- Run: `./handoff [options] [path1] [path2] ...`
  - Key options: `-verbose`, `-output`, `-include`, `-exclude`, `-ignore-gitignore`
  - Use `--help` for complete flag reference
- Test all: `go test ./...`
- Test specific: `go test -run TestName`
- Test with verbosity: `go test -v ./...`
- Test single file: `go test -run TestFunctionName path/to/file_test.go`
- Format code: `gofmt -w .`
- Lint: `golangci-lint run`

## Code Style Guidelines
- **Imports**: Standard library first, then third-party, separated by blank lines
- **Formatting**: Use `gofmt` for standard Go formatting
- **Documentation**: All exported functions/types must have comments
- **Error Handling**: Check errors explicitly, use named returns for clarity
- **Naming**: CamelCase for exported names, camelCase for unexported
- **Error Messages**: Lowercase, no trailing punctuation, descriptive but concise
- **Functions**: Keep functions focused and under 50 lines when possible
- **Files**: Group related functionality, keep files under 500 lines

## Testing Principles
- Write deterministic, repeatable tests
- Use table-driven tests for testing multiple cases
- Test both success and error paths
- Provide clear error messages in test failures
- Consider Test-Driven Development (TDD) where appropriate

## Git Workflow
- Use conventional commit prefixes: `feat:`, `fix:`, `docs:`, `test:`, `chore:`
- Make atomic, focused commits that implement one logical change
- Update documentation as part of feature implementation
- NEVER include AI attribution in commit messages (e.g., "Generated by Claude"). Commit messages should ONLY contain detailed multiline conventional commit messages about the work done
