# Handoff Library

[![Go Reference](https://pkg.go.dev/badge/github.com/phrazzld/handoff/lib.svg)](https://pkg.go.dev/github.com/phrazzld/handoff/lib)
[![Test Coverage](https://img.shields.io/badge/coverage-85%2B%25-brightgreen)](.github/workflows/test-coverage.yml)

This package provides programmatic access to Handoff's core functionality for collecting and formatting file contents. It can be used by other Go programs to gather code from a project and use it for various purposes like documentation, AI assistant input, code analysis, and more.

## Table of Contents

- [Installation](#installation)
- [Basic Usage](#usage)
- [Configuration](#configuration)
- [Core Functions](#core-functions)
- [Advanced Usage](#advanced-usage)
- [Development](#development)

## Installation

```bash
go get github.com/phrazzld/handoff/lib
```

## Usage

Basic usage of the library involves creating a configuration, processing files, and using the resulting content:

```go
package main

import (
	"fmt"
	"github.com/phrazzld/handoff/lib"
	"os"
)

func main() {
	// Create configuration
	config := lib.NewConfig()
	config.Verbose = true
	config.Exclude = ".exe,.bin,.jpg,.png"
	config.ExcludeNamesStr = "node_modules,package-lock.json"
	
	// REQUIRED: Call ProcessConfig() to convert string-based settings into slices
	// Skipping this step will cause your include/exclude filters not to work!
	config.ProcessConfig()

	// Process project files
	content, stats, err := lib.ProcessProject([]string{"./my-project"}, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing project: %v\n", err)
		os.Exit(1)
	}

	// Use the content - write to file
	if err := lib.WriteToFile(content, "output.md"); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
		os.Exit(1)
	}
	
	// Use the statistics returned from ProcessProject
	fmt.Printf("Generated content with %d characters, %d lines, and approximately %d tokens\n", 
		stats.Chars, stats.Lines, stats.Tokens)
	fmt.Printf("Processed %d out of %d total files\n", 
		stats.FilesProcessed, stats.FilesTotal)
}
```

## Core Functions

The library provides several key functions:

### ProcessProject

```go
func ProcessProject(paths []string, config *Config) (string, Stats, error)
```

The main function that processes one or more files or directories and returns their formatted content along with statistics.

- **Parameters:**
  - `paths []string`: File or directory paths to process
  - `config *Config`: Configuration options for processing (can be nil for defaults)
- **Returns:**
  - `string`: Formatted content from all processed files
  - `Stats`: Statistics about processed files and content
  - `error`: Any error encountered during processing
- **Error handling:**
  - Returns errors for inaccessible paths or problems reading files
  - Non-critical errors (like skipping a single file) are logged but don't stop processing
- **Notes:**
  - ProcessProject automatically calls ProcessConfig() on the provided configuration
  - This means you can set string-based filters directly before calling ProcessProject
  - However, for clarity and consistent practice, it's still recommended to call ProcessConfig() yourself after setting configuration options

### WriteToFile

```go
func WriteToFile(content, filePath string) error
```

Utility to write content to a file.

- **Parameters:**
  - `content string`: The content to write
  - `filePath string`: The path where the file should be written
- **Returns:**
  - `error`: Any error encountered while writing
- **Notes:**
  - Creates parent directories if they don't exist
  - Will overwrite existing files without warning

### CalculateStatistics

```go
func CalculateStatistics(content string) (charCount, lineCount, tokenCount int)
```

Analyzes content and returns statistics.

- **Parameters:**
  - `content string`: The text content to analyze
- **Returns:**
  - `charCount int`: Number of characters in the content
  - `lineCount int`: Number of lines in the content
  - `tokenCount int`: Estimated token count (helpful for LLM context limits)
- **Notes:**
  - Token count is an approximation based on common tokenization rules (roughly 4 characters per token)

## Configuration

Use the `Config` struct to configure how files are processed:

```go
type Config struct {
    // String-based config options (for CLI flags)
    Verbose        bool   // Enable verbose logging
    Include        string // Comma-separated extensions to include (e.g., ".go,.txt")
    Exclude        string // Comma-separated extensions to exclude (e.g., ".bin,.exe")
    ExcludeNamesStr string // Comma-separated filenames to exclude (e.g., "package-lock.json")
    Format         string // Template for formatting output

    // Processed slice versions (used internally)
    IncludeExts    []string // Processed Include string as slice
    ExcludeExts    []string // Processed Exclude string as slice
    ExcludeNames   []string // Processed ExcludeNamesStr as slice
}
```

### Creating and Configuring

```go
// Create a new config with default values
config := lib.NewConfig()

// Configure options
config.Verbose = true
config.Include = ".go,.ts,.js"
config.Exclude = ".exe,.dll,.jpg,.png,.gif"
config.ExcludeNamesStr = "node_modules,dist,build"
config.Format = "File: {path}\n```\n{content}\n```\n\n"

// IMPORTANT: Process string-based configs into slices (required before use)
config.ProcessConfig()
```

### IMPORTANT: Processing Configuration

**You MUST call `ProcessConfig()` after setting any string-based configuration options.**

The `ProcessConfig()` method converts the string-based configuration fields (`Include`, `Exclude`, `ExcludeNamesStr`) into their corresponding slice fields (`IncludeExts`, `ExcludeExts`, `ExcludeNames`) that are used during file processing.

If you forget to call `ProcessConfig()`:
- Your string-based include/exclude patterns will NOT be applied
- No files will be filtered by extension or name as intended
- The library will use whatever was previously in the slice fields (usually empty)

### Key Configuration Options

- **Format**: Template for formatting each file's output
  - Uses `{path}` and `{content}` placeholders
  - Default: `<{path}>\n```\n{content}\n```\n</{path}>\n\n`

- **Include/IncludeExts**: File extensions to include
  - If specified, only files with these extensions will be processed
  - Can be provided with or without dots: `.go,.txt` or `go,txt`
  - If not specified, all non-excluded files are processed

- **Exclude/ExcludeExts**: File extensions to exclude
  - Files with these extensions will be skipped
  - Default excludes common binary files and images

- **ExcludeNamesStr/ExcludeNames**: File names to exclude
  - Files with these exact names will be skipped
  - Useful for excluding specific files like `package-lock.json` or directories like `node_modules`

- **Verbose**: Enable detailed logging
  - When true, shows verbose information about file processing

## Additional Examples

### Processing Files with Different Configurations

```go
package main

import (
	"fmt"
	"github.com/phrazzld/handoff/lib"
)

func main() {
	// Example 1: Default configuration
	defaultConfig := lib.NewConfig()
	defaultConfig.ProcessConfig() // Required step to process any string-based settings
	
	content1, stats1, err := lib.ProcessProject([]string{"./src"}, defaultConfig)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Default config: processed %d files\n", stats1.FilesProcessed)
	
	// Example 2: Include only specific file types
	codeConfig := lib.NewConfig()
	codeConfig.Include = ".go,.ts,.js"
	codeConfig.ProcessConfig() // Required to convert string filters to slices
	
	content2, stats2, err := lib.ProcessProject([]string{"./src"}, codeConfig)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Code files only: processed %d files\n", stats2.FilesProcessed)
	
	// Example 3: Custom format
	markdownConfig := lib.NewConfig()
	markdownConfig.Format = "## {path}\n\n```go\n{content}\n```\n\n"
	markdownConfig.ProcessConfig() // Note: ProcessProject calls this internally, but it's good practice to call it explicitly
	
	content3, stats3, err := lib.ProcessProject([]string{"./main.go"}, markdownConfig)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Custom format: processed %d lines\n", stats3.Lines)
	
	// Write results to separate files
	lib.WriteToFile(content1, "all_files.md")
	lib.WriteToFile(content2, "code_files.md") 
	lib.WriteToFile(content3, "markdown_format.md")
}
```

### Handling Multiple Paths and Error Checking

```go
func processMultiplePaths() {
	config := lib.NewConfig()
	// No string-based settings to process in this example, but it's
	// still good practice to call ProcessConfig() before using the config
	config.ProcessConfig()
	
	// Process multiple specific files
	paths := []string{
		"./main.go",
		"./lib/handoff.go",
		"./README.md",
	}
	
	content, stats, err := lib.ProcessProject(paths, config)
	if err != nil {
		// Still check the content - ProcessProject may return partial results
		// even if some files failed to process
		if content == "" {
			fmt.Println("No content was processed successfully")
			return
		}
		fmt.Printf("Warning: Some files could not be processed: %v\n", err)
	}
	
	// Use the statistics returned from ProcessProject
	fmt.Printf("Processed %d/%d files\n", stats.FilesProcessed, stats.FilesTotal)
	fmt.Printf("Content has %d lines (%d chars, ~%d tokens)\n", 
		stats.Lines, stats.Chars, stats.Tokens)
}
```


## Advanced Usage

For advanced use cases, the library provides a few additional exported functions and types that can be useful for custom workflows:

### ProcessorFunc Type

```go
type ProcessorFunc func(filePath string, content []byte) string
```

The `ProcessorFunc` type represents a function that can process a file's content. This is the signature for custom processors that can be used with the library's internal processing mechanisms. While you can't directly use this with the exported API currently, understanding this signature is helpful if you're implementing your own file processing logic.

### Logger

```go
type Logger struct {
    // Unexported fields
}

func NewLogger(verbose bool) *Logger
```

The library provides a simple logger type for controlling output verbosity. While primarily used internally, it's available if you need consistent logging.

```go
// Create a new logger with verbose output enabled
logger := lib.NewLogger(true) 

// Use the logger methods
logger.Info("Processing started")
logger.Verbose("Detailed information that only shows in verbose mode")
logger.Warn("Warning message")
logger.Error("Error message")
```

### Stats Struct

```go
type Stats struct {
    FilesProcessed int
    FilesTotal int
    Lines int
    Chars int
    Tokens int
}
```

The `Stats` struct provides detailed information about processed content. It's returned by `ProcessProject` and contains metrics about the files and content processed.

```go
// Get content and stats from processing
content, stats, err := lib.ProcessProject(paths, config)
if err != nil {
    // Handle error
}

// Use the stats in your application
fmt.Printf("Processed %d/%d files\n", stats.FilesProcessed, stats.FilesTotal)
fmt.Printf("Content has %d lines, %d characters, and approximately %d tokens\n", 
    stats.Lines, stats.Chars, stats.Tokens)
```

### WrapInContext

```go
func WrapInContext(content string) string
```

Utility function to wrap content in context tags, which can be useful for formatting output for LLMs or other processors.

```go
// Wrap some content in context tags
rawContent := "Some content to be wrapped"
wrappedContent := lib.WrapInContext(rawContent)
fmt.Println(wrappedContent)
// Output:
// <context>
// Some content to be wrapped
// </context>
```

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