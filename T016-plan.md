# T016 Plan: Add CLI Integration Tests Verifying Library Usage

## Task Analysis
- Need to create integration tests that execute the compiled CLI binary
- Tests should verify that CLI flags correctly influence output via the library calls
- Need to ensure all CLI functionality is properly tested
- Tests should be idempotent and not rely on specific environment configurations

## Implementation Approach

1. **Create Test Helper Functions**
   - Create a function to build the CLI binary for testing
   - Create functions to run the binary with different arguments
   - Create functions to capture and verify outputs

2. **Test Cases to Implement**
   - Test CLI help output (running with --help flag)
   - Test basic file processing (processing a single file)
   - Test directory processing (processing a directory of files)
   - Test include/exclude filtering options
   - Test output to file functionality
   - Test file overwrite protection (with and without -force)
   - Test dry-run functionality
   - Test error conditions (non-existent paths, etc.)

3. **Implementation Strategy**
   - Use `os/exec` package to run the compiled binary
   - Create temporary test directories and files for testing
   - Capture standard output, standard error, and exit codes
   - Verify outputs match expected results
   - Clean up temporary files and binary after tests

## Technical Considerations
- Test execution will require compiling the binary, which may add time to tests
- Need to gracefully handle platform differences (Windows vs Unix)
- Need to avoid actual clipboard operations in tests
- Should mock or skip clipboard operations in CI environments

## Testing Structure
- Create a new test file or extend handoff_test.go with integration-specific tests
- Group tests logically by functionality being tested
- Use subtests for different variations of the same functionality 
- Ensure proper setup and teardown for all test resources