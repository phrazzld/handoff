# PLAN: Convert Handoff to a Proper Go Package

This plan outlines the steps needed to refactor Handoff into a proper Go package that other programs can easily use, while maintaining the existing CLI functionality.

## 1. Current State Analysis

The codebase already has some modular structure with a `lib` package, but there's redundancy between the main package files (`files.go`, `output.go`) and the library (`lib/handoff.go`). The core functionality needs to be fully encapsulated in the library with the CLI becoming a thin wrapper around it.

## 2. Implementation Steps

### A. Consolidate Core Logic into `lib` Package

1. **Review and move remaining logic from main package to `lib`**:
   - Identify unique logic in `files.go` and `output.go`
   - Move any missing functionality to `lib/handoff.go`
   - Remove redundant files after confirmation (`files.go`, `output.go`)

2. **Define clear public API**:
   - Ensure all exported functions have proper documentation
   - Add package-level documentation
   - Refine exported vs unexported functions

### B. Refactor CLI to Use Library

1. **Update `main.go` to be a thin wrapper**:
   - Import `github.com/phrazzld/handoff/lib`
   - Modify `parseConfig()` to return a `*handoff.Config`
   - Replace direct operations with library calls
   - Keep CLI-specific operations (clipboard, output file handling)

2. **Ensure CLI functionality is preserved**:
   - Maintain all existing flags and behavior
   - Support the same output formats and filtering options

### C. Update and Migrate Tests

1. **Create `lib/handoff_test.go`**:
   - Move relevant tests from root `handoff_test.go`
   - Add new tests for library API
   - Ensure comprehensive coverage of library functions

2. **Update main package tests**:
   - Keep CLI-specific tests in root `handoff_test.go`
   - Test flag parsing and CLI behaviors
   - Verify integration between CLI and library

### D. Documentation and Examples

1. **Update project documentation**:
   - Revise main `README.md` to highlight library and CLI usage
   - Update `lib/README.md` with detailed API usage
   - Add comprehensive GoDoc comments to all exported identifiers

2. **Create example code**:
   - Update `examples/gemini_planner.go` if needed
   - Consider adding additional simple examples

### E. Versioning and Go Modules

1. **Finalize Go module configuration**:
   - Ensure `go.mod` is correct
   - Run `go mod tidy`
   - Plan for semantic versioning

## 3. Testing Strategy

- **Library unit tests**: Test individual functions in isolation
- **Library integration tests**: Test `ProcessProject` with various configurations
- **CLI integration tests**: Test CLI behaviors and flag interactions
- **Documentation tests**: Ensure examples compile and run

## 4. Deliverables

1. A clean library package (`github.com/phrazzld/handoff/lib`) with:
   - Well-documented API
   - Comprehensive test coverage
   - Examples of usage

2. A CLI tool (`github.com/phrazzld/handoff`) that:
   - Uses the library package
   - Maintains existing functionality
   - Has clean separation of concerns

## 5. Example Library API

```go
// Config holds configuration for file collection/formatting
type Config struct {
    Include       string
    Exclude       string
    ExcludeNames  string
    Format        string
    Verbose       bool
    // ...
}

// NewConfig returns a Config with defaults
func NewConfig() *Config

// ProcessConfig processes string-based configs into slices
func (c *Config) ProcessConfig()

// ProcessProject collects and formats files from given paths
func ProcessProject(paths []string, config *Config) (string, error)

// WriteToFile writes content to a file
func WriteToFile(content, filePath string) error

// CalculateStatistics calculates statistics about the content
func CalculateStatistics(content string) (charCount, lineCount, tokenCount int)
```

## 6. Potential Challenges

- Maintaining backward compatibility for existing users
- Ensuring proper error handling in the library vs. CLI
- Consistent testing across different environments (Git dependency)
- Finding the right balance of API surface area (not too broad, not too narrow)