# T021 Implementation Plan: Add Simple Library Usage Example

## Overview

This task involves creating a simple, standalone example that demonstrates the basic usage of the handoff library. The example should be minimal but complete, showing the essential steps for using the library's core functionality.

## Implementation Steps

1. **Create New Example File**
   - Create a new file at `examples/simple_usage.go`
   - Structure it as a complete, runnable Go program

2. **Implement Basic Library Usage**
   - Import the `lib` package
   - Demonstrate creating and configuring a Config object
   - Show how to call ProcessProject
   - Handle errors appropriately
   - Process the results (output to console/file)

3. **Add Command-line Options**
   - Include basic command-line flags for:
     - Input directory/files
     - Output file
     - Basic filtering options
   - Parse and use these flags in the example

4. **Add Comprehensive Documentation**
   - Include clear comments explaining each step
   - Document command-line options
   - Add usage examples in comments

5. **Verify the Example**
   - Ensure it builds successfully
   - Test run the example with different inputs
   - Verify output is as expected

## Technical Implementation Details

The example will follow this general structure:

```go
// Simple example demonstrating basic usage of the handoff library
package main

import (
    "flag"
    "fmt"
    "os"
    
    "github.com/phrazzld/handoff/lib"
)

func main() {
    // Parse command-line flags
    inputDir := flag.String("dir", ".", "Directory to process")
    outputFile := flag.String("output", "", "Output file (if empty, prints to stdout)")
    includeExts := flag.String("include", "", "File extensions to include (e.g., '.go,.txt')")
    excludeExts := flag.String("exclude", "", "File extensions to exclude (e.g., '.exe,.bin')")
    flag.Parse()
    
    // Create and configure handoff
    config := lib.NewConfig()
    config.Include = *includeExts
    config.Exclude = *excludeExts
    config.ProcessConfig()
    
    // Process files
    content, err := lib.ProcessProject([]string{*inputDir}, config)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
    
    // Output results
    if *outputFile != "" {
        if err := lib.WriteToFile(content, *outputFile); err != nil {
            fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Content written to %s\n", *outputFile)
    } else {
        fmt.Println(content)
    }
    
    // Display statistics
    chars, lines, tokens := lib.CalculateStatistics(content)
    fmt.Printf("Stats: %d chars, %d lines, ~%d tokens\n", chars, lines, tokens)
}
```

## Verification

To verify the example works correctly:
1. Build the example: `go build ./examples/simple_usage.go`
2. Run with various options:
   - Basic: `./simple_usage --dir ./some_directory`
   - With output file: `./simple_usage --dir ./some_directory --output output.md`
   - With filtering: `./simple_usage --dir ./some_directory --include .go,.md`

## Success Criteria

1. The example is simple but complete
2. It demonstrates all core steps: configuration, processing, output, statistics
3. It includes clear documentation and comments
4. It compiles and runs successfully
5. The code is clean and follows best practices