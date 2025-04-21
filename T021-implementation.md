# T021 Implementation: Add Simple Library Usage Example

## Implementation Summary

A simple, standalone example `examples/simple_usage.go` has been created to demonstrate the basic usage of the handoff library. This example serves as a minimal but complete reference for using the library's core functionality.

## Implementation Details

### 1. Created New Example File

Created a new file at `examples/simple_usage.go` that demonstrates:
- Creating and configuring a Config object
- Processing files with the library
- Calculating statistics on the output
- Writing content to a file or displaying on console

### 2. Implemented Command-line Interface

Added command-line flags to demonstrate common configuration options:
- `--dir`: Directory or file to process (defaults to current directory)
- `--output`: File to write the output to (if not specified, prints to console)
- `--include`: File extensions to include
- `--exclude`: File extensions to exclude (with sensible defaults)
- `--verbose`: Enable verbose output

### 3. Added Comprehensive Documentation

The example includes:
- Detailed header comments explaining the purpose
- Multiple usage examples showing different flag combinations
- Inline comments explaining each step of the process
- Comments explaining why certain steps (like `ProcessConfig()`) are necessary

### 4. Featured Core Library Functionality

Demonstrated all key library functions:
- `lib.NewConfig()`: Creating a configuration
- `config.ProcessConfig()`: Processing string-based settings
- `lib.ProcessProject()`: Processing files to get formatted content
- `lib.CalculateStatistics()`: Getting content statistics
- `lib.WriteToFile()`: Writing output to a file

### 5. Implemented Error Handling

The example shows proper error handling for:
- Processing errors
- File writing errors

## Verification

The example was verified by:
1. Building the file with `go build ./examples/simple_usage.go`
2. Running it with various options to confirm it works correctly:
   - Processing a specific directory
   - Writing output to a file
   - Displaying statistics
3. Verifying the code is clean, readable, and follows best practices

## Success Criteria Met

1. ✅ Created a simple but complete example
2. ✅ Demonstrated all core steps: configuration, processing, output, statistics
3. ✅ Included clear documentation and comments
4. ✅ Code compiles and runs successfully
5. ✅ Followed clean code and best practices