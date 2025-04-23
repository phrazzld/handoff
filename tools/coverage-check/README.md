# Coverage Checker

A simple Go tool to check if test coverage meets a specified threshold.

## Usage

```bash
# Check coverage from a file
coverage-check -file coverage.out -threshold 85.0

# Read coverage profile from stdin
cat coverage.out | coverage-check -threshold 90.0

# Show detailed output
coverage-check -file coverage.out -threshold 85.0 -verbose
```

## Command Line Options

- `-file string`: Coverage profile file (default reads from stdin)
- `-threshold float`: Minimum coverage percentage required (default 85.0)
- `-verbose`: Show detailed output

## Exit Codes

- 0: Coverage meets or exceeds threshold
- 1: Coverage is below threshold
- 2: Error parsing coverage data or other issue

## Building

```bash
go build -o coverage-check
```

## Running Tests

```bash
go test ./...
```
