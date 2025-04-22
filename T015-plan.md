# T015 · Feature · P2: Create Go Coverage Checker Tool

## Task Description
Create a simple Go tool that parses Go's coverage output and checks if it meets a specified threshold.

## Current State
Currently, the CI workflow likely uses a shell script to check coverage thresholds, which can be brittle.

## Implementation Plan

1. Create the `tools/coverage-check` directory 
2. Design a simple Go program that:
   - Accepts coverage output (either from stdin or file)
   - Parses the coverage percentage
   - Compares against a threshold
   - Returns appropriate exit code based on pass/fail
   - Provides clear output for CI logs

3. Add unit tests that verify:
   - Correct parsing of various coverage report formats
   - Proper threshold comparison logic
   - Appropriate exit code handling
   - Edge cases (empty input, malformed input, etc.)

## Implementation Details

### Program Structure
- `main.go`: Entry point 
- `parser.go`: Coverage report parsing logic
- `checker.go`: Threshold checking logic
- `parser_test.go` and `checker_test.go`: Unit tests

### CLI Interface
```
Usage: coverage-check [options]

Options:
  -file string
        Coverage profile file (default reads from stdin)
  -threshold float
        Minimum coverage percentage required (default 85.0)
  -verbose
        Show detailed output
```

### Exit Codes
- 0: Coverage meets or exceeds threshold
- 1: Coverage is below threshold
- 2: Error parsing coverage data or other issue

## Testing Strategy
- Create sample coverage output files with various percentages
- Test parsing logic with different formats of coverage output
- Test threshold comparison with edge cases (exact match, slightly below, etc.)
- Mock stdin for testing reading from standard input