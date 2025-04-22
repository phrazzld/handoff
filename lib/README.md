# Handoff Library

[![Go Reference](https://pkg.go.dev/badge/github.com/phrazzld/handoff/lib.svg)](https://pkg.go.dev/github.com/phrazzld/handoff/lib)
[![Test Coverage](https://img.shields.io/badge/coverage-85%2B%25-brightgreen)](.github/workflows/test-coverage.yml)

This package provides programmatic access to Handoff's core functionality for collecting and formatting file contents. It can be used by other Go programs to gather code from a project and use it for various purposes like documentation, AI assistant input, code analysis, and more.

## Table of Contents

- [Installation](#installation)
- [Basic Usage](#usage)
- [Configuration](#configuration)
- [Core Functions](#core-functions)
- [Examples](#examples)
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
	config.ProcessConfig() // This converts the string settings into slices

	// Process project files
	content, err := lib.ProcessProject([]string{"./my-project"}, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing project: %v\n", err)
		os.Exit(1)
	}

	// Use the content - write to file
	if err := lib.WriteToFile(content, "output.md"); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
		os.Exit(1)
	}
	
	// Get statistics about the processed content
	chars, lines, tokens := lib.CalculateStatistics(content)
	fmt.Printf("Generated content with %d characters, %d lines, and approximately %d tokens\n", 
		chars, lines, tokens)
}
```

## Core Functions

The library provides several key functions:

### ProcessProject

```go
func ProcessProject(paths []string, config *Config) (string, error)
```

The main function that processes one or more files or directories and returns their formatted content.

- **Parameters:**
  - `paths []string`: File or directory paths to process
  - `config *Config`: Configuration options for processing
- **Returns:**
  - `string`: Formatted content from all processed files
  - `error`: Any error encountered during processing
- **Error handling:**
  - Returns errors for inaccessible paths or problems reading files
  - Non-critical errors (like skipping a single file) are logged but don't stop processing

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
  - Returns an error if the file already exists (use the CLI with `-force` for overwriting)

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

// Process string-based configs into slices (required before use)
config.ProcessConfig()
```

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

## Examples

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
	defaultConfig.ProcessConfig()
	
	content1, err := lib.ProcessProject([]string{"./src"}, defaultConfig)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	
	// Example 2: Include only specific file types
	codeConfig := lib.NewConfig()
	codeConfig.Include = ".go,.ts,.js"
	codeConfig.ProcessConfig()
	
	content2, err := lib.ProcessProject([]string{"./src"}, codeConfig)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	
	// Example 3: Custom format
	markdownConfig := lib.NewConfig()
	markdownConfig.Format = "## {path}\n\n```go\n{content}\n```\n\n"
	markdownConfig.ProcessConfig()
	
	content3, err := lib.ProcessProject([]string{"./main.go"}, markdownConfig)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	
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
	config.ProcessConfig()
	
	// Process multiple specific files
	paths := []string{
		"./main.go",
		"./lib/handoff.go",
		"./README.md",
	}
	
	content, err := lib.ProcessProject(paths, config)
	if err != nil {
		// Still check the content - ProcessProject may return partial results
		// even if some files failed to process
		if content == "" {
			fmt.Println("No content was processed successfully")
			return
		}
		fmt.Printf("Warning: Some files could not be processed: %v\n", err)
	}
	
	// Get statistics
	chars, lines, tokens := lib.CalculateStatistics(content)
	fmt.Printf("Processed %d lines (%d chars, ~%d tokens)\n", lines, chars, tokens)
}
```

### Gemini Integration

For a complete example of using this library with Google's Gemini API, see `examples/gemini_planner.go` in the main project repository.

## Advanced Usage

For advanced use cases, you can access lower-level functions to implement custom processing logic:

### GetFilesFromDir

```go
func GetFilesFromDir(dir string) ([]string, error)
```

Gets a list of files from a directory, respecting Git ignore rules if applicable.

```go
// Example: Get files and process them individually
files, err := lib.GetFilesFromDir("./src")
if err != nil {
    return err
}

for _, file := range files {
    // Custom processing logic for each file
    fmt.Println("Processing:", file)
}
```

### ShouldProcess

```go
func ShouldProcess(file string, config *Config) bool
```

Determines if a file should be processed based on configuration filters.

```go
// Example: Custom file filtering
files, _ := lib.GetFilesFromDir("./")
config := lib.NewConfig()
config.Include = ".md"
config.ProcessConfig()

for _, file := range files {
    if lib.ShouldProcess(file, config) {
        fmt.Println("Will process:", file)
    } else {
        fmt.Println("Will skip:", file)
    }
}
```

### ProcessFile

```go
func ProcessFile(filePath string, logger *Logger, config *Config, processor ProcessorFunc) string
```

Processes a single file with a custom processor function.

```go
// Example: Custom file processing
logger := lib.NewLogger(true)
config := lib.NewConfig()
config.ProcessConfig()

// Define a custom processor function
processor := func(filePath string, content []byte) string {
    // Custom content transformation
    return fmt.Sprintf("PROCESSED: %s\n%s", filePath, string(content))
}

result := lib.ProcessFile("./main.go", logger, config, processor)
fmt.Println(result)
```

### IsBinaryFile

```go
func IsBinaryFile(content []byte) bool
```

Determines if file content appears to be binary.

```go
// Example: Check if a file is binary
data, _ := os.ReadFile("./some-file")
if lib.IsBinaryFile(data) {
    fmt.Println("This is a binary file")
} else {
    fmt.Println("This is a text file")
}
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