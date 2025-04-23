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
	// Create configuration using functional options
	config := lib.NewConfig(
		lib.WithVerbose(true),
		lib.WithExclude(".exe,.bin,.jpg,.png"),
		lib.WithExcludeNames("node_modules,package-lock.json"),
	)

	// Process project files
	content, stats, err := lib.ProcessProject([]string{"./my-project"}, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing project: %v\n", err)
		os.Exit(1)
	}

	// Use the content - write to file
	if err := lib.WriteToFile(content, "output.md", true); err != nil {
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

## Core API

The library provides a focused API with ProcessProject as the main entry point:

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
  - When using functional options pattern (recommended), no additional configuration processing is needed
  - For backward compatibility, ProcessProject will call ProcessConfig() if needed
  - The recommended approach is to use functional options for a cleaner, more maintainable codebase

### WriteToFile

```go
func WriteToFile(content, filePath string, overwrite bool) error
```

Utility to write content to a file with overwrite control.

- **Parameters:**
  - `content string`: The content to write
  - `filePath string`: The path where the file should be written
  - `overwrite bool`: Whether to overwrite the file if it already exists
- **Returns:**
  - `error`: Any error encountered while writing, including `ErrFileExists` if the file exists and `overwrite` is false
- **Notes:**
  - Creates parent directories if they don't exist
  - Controls overwriting behavior with the `overwrite` parameter
  - Returns `ErrFileExists` when trying to write to an existing file with `overwrite=false`

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
  - **Token count is a simple approximation** based on whitespace boundaries:
    - Counts transitions between whitespace and non-whitespace characters
    - Treats any continuous sequence of non-whitespace characters as one token
    - Is significantly less sophisticated than actual LLM tokenizers
  - **Limitations compared to real LLM tokenizers:**
    - Real tokenizers use subword tokenization algorithms with trained vocabularies
    - They have special handling for punctuation, common words, and different languages
    - The approximation may undercount tokens (e.g., punctuation often gets separate tokens)
    - The approximation may overcount tokens (e.g., common words often get single tokens)
  - **Usage guidance:**
    - For LLM usage planning, consider adding a 30-50% safety margin to these estimates
    - When precise token counts matter, use the tokenizer specific to your LLM provider

## Configuration

There are two ways to configure the library: the recommended functional options pattern and the traditional approach (maintained for backward compatibility).

### Recommended: Functional Options Pattern

The functional options pattern is the recommended way to configure the library. It's cleaner, less error-prone, and handles all processing internally:

```go
// Create configuration with functional options
config := lib.NewConfig(
    lib.WithVerbose(true),                        // Enable verbose logging
    lib.WithInclude(".go,.ts,.js"),               // Only include these extensions
    lib.WithExclude(".exe,.dll,.jpg,.png,.gif"),  // Exclude these extensions
    lib.WithExcludeNames("node_modules,dist"),    // Exclude these file/directory names
    lib.WithFormat("File: {path}\n```\n{content}\n```\n\n"), // Custom format
)

// Use the configuration directly - no additional processing needed
content, stats, err := lib.ProcessProject(paths, config)
```

This approach:
- Eliminates the need to call `ProcessConfig()`
- Processes string inputs automatically
- Provides better type safety and discoverability
- Works directly with ProcessProject with no intermediate steps

### Legacy: Traditional Configuration

For backward compatibility, you can still use the direct field approach, but it requires an additional step:

```go
// Create a new config with default values
config := lib.NewConfig()

// Configure options
config.Verbose = true
config.Include = ".go,.ts,.js"
config.Exclude = ".exe,.dll,.jpg,.png,.gif"
config.ExcludeNamesStr = "node_modules,dist,build"
config.Format = "File: {path}\n```\n{content}\n```\n\n"

// REQUIRED for the traditional approach: Process string-based configs into slices
config.ProcessConfig()
```

**Note**: When using the traditional approach, you MUST call `ProcessConfig()` after setting string-based options. This is not required when using the functional options pattern.

### Key Configuration Options

- **Format**: Template for formatting each file's output
  - Functional option: `WithFormat("template string")`
  - Uses `{path}` and `{content}` placeholders
  - Default: `<{path}>\n```\n{content}\n```\n</{path}>\n\n`

- **Include**: File extensions to include
  - Functional option: `WithInclude(".go,.txt")` 
  - If specified, only files with these extensions will be processed
  - Can be provided with or without dots: `.go,.txt` or `go,txt`
  - If not specified, all non-excluded files are processed

- **Exclude**: File extensions to exclude
  - Functional option: `WithExclude(".bin,.exe")`
  - Files with these extensions will be skipped
  - Default excludes common binary files and images

- **ExcludeNames**: File names to exclude
  - Functional option: `WithExcludeNames("package-lock.json,node_modules")`
  - Files with these exact names will be skipped
  - Useful for excluding specific files or directories

- **Verbose**: Enable detailed logging
  - Functional option: `WithVerbose(true)`
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
	content1, stats1, err := lib.ProcessProject([]string{"./src"}, defaultConfig)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Default config: processed %d files\n", stats1.FilesProcessed)
	
	// Example 2: Include only specific file types
	codeConfig := lib.NewConfig(
		lib.WithInclude(".go,.ts,.js"),  // Only process code files
	)
	
	content2, stats2, err := lib.ProcessProject([]string{"./src"}, codeConfig)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Code files only: processed %d files\n", stats2.FilesProcessed)
	
	// Example 3: Custom format
	markdownConfig := lib.NewConfig(
		lib.WithFormat("## {path}\n\n```go\n{content}\n```\n\n"),
	)
	
	content3, stats3, err := lib.ProcessProject([]string{"./main.go"}, markdownConfig)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Custom format: processed %d lines\n", stats3.Lines)
	
	// Write results to separate files
	lib.WriteToFile(content1, "all_files.md", true)
	lib.WriteToFile(content2, "code_files.md", true) 
	lib.WriteToFile(content3, "markdown_format.md", true)
}
```

### Handling Multiple Paths and Error Checking

```go
func processMultiplePaths() {
	// Create configuration with functional options (if needed)
	config := lib.NewConfig(
		lib.WithVerbose(true),  // Enable verbose logging
	)
	
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


## Advanced Features

The library includes additional types and utilities to support more complex use cases:

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
