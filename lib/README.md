# Handoff Library

This package provides programmatic access to Handoff's core functionality for collecting and formatting file contents. It can be used by other Go programs to gather code from a project and use it for various purposes.

## Usage

```go
import "github.com/phrazzld/handoff/lib"

// Create configuration
config := handoff.NewConfig()
config.Verbose = true
config.Exclude = ".exe,.bin,.jpg,.png"
config.ExcludeNamesStr = "node_modules,package-lock.json"
config.ProcessConfig() // This converts the string settings into slices

// Process project files
content, err := handoff.ProcessProject([]string{"./my-project"}, config)
if err != nil {
    panic(err)
}

// Use the content (send to API, write to file, etc.)
fmt.Println(content)
```

## Core Functions

- `ProcessProject(paths []string, config *Config) (string, error)`: Main function to process a list of paths and get formatted content
- `WriteToFile(content, filePath string) error`: Utility to write content to a file
- `CalculateStatistics(content string) (charCount, lineCount, tokenCount int)`: Get content statistics

## Configuration

Use the `Config` struct to configure how files are processed:

```go
type Config struct {
    Verbose        bool
    Include        string
    Exclude        string
    ExcludeNamesStr string
    Format         string
    IncludeExts    []string
    ExcludeExts    []string
    ExcludeNames   []string
}
```

Key settings:
- `Format`: Template for formatting each file with `{path}` and `{content}` placeholders
- `Include`/`IncludeExts`: File extensions to include (e.g., ".go,.txt")
- `Exclude`/`ExcludeExts`: File extensions to exclude (e.g., ".bin,.exe")
- `ExcludeNamesStr`/`ExcludeNames`: File names to exclude (e.g., "package-lock.json")
- `Verbose`: Enable detailed logging

## Example: Gemini Integration

For a complete example of using this library with Gemini API, see `examples/gemini_planner.go` in the main project repository.

## Advanced Usage

For advanced use cases, you can access lower-level functions:

- `GetFilesFromDir(dir string) ([]string, error)`: Get filtered list of files from a directory
- `ShouldProcess(file string, config *Config) bool`: Check if a file should be processed
- `ProcessFile(filePath string, logger *Logger, config *Config, processor ProcessorFunc) string`: Process a single file with custom processor
- `IsBinaryFile(content []byte) bool`: Detect if content is binary

## Development

### Test Coverage

The library package maintains a minimum test coverage threshold of 85%. This is enforced via GitHub Actions for all pull requests and pushes to the main branch. If you contribute to this library, ensure your changes include appropriate test coverage.

To check coverage locally:

```bash
go test -coverprofile=coverage.out ./lib/...
go tool cover -func=coverage.out
```

For a visual report in your browser:

```bash
go tool cover -html=coverage.out
```