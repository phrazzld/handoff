# Add test for error handling on invalid paths

## Goal
Create a test to verify proper error handling when `-output` points to an invalid or inaccessible path.

## Implementation Approach
I'll create a new integration test function named `TestInvalidPathErrorHandling` that will test error handling for invalid output paths in two scenarios:

1. **Invalid directory path** - Test when `-output` points to a file in a non-existent directory
   - Create a temporary test directory
   - Set up a path that points to a file in a subdirectory that doesn't exist
   - Attempt to use this path with the `-output` flag
   - Verify that an appropriate error is detected and handled correctly

2. **Inaccessible path** - Test when `-output` points to a location that can't be written to due to permissions
   - Create a read-only directory within the temporary test directory
   - Set up a path that points to a file within this read-only directory
   - Attempt to use this path with the `-output` flag
   - Verify that a permission error is correctly detected and handled

For each scenario, I'll need to modify the testing approach slightly compared to previous tests, because the `main()` function exits with `os.Exit(1)` when it encounters these errors. Since we can't capture the exit in a test, I'll:

1. Mock the necessary components to test the error detection logic
2. Focus on testing that the application correctly identifies the error conditions without actually calling `os.Exit()`
3. Check that appropriate error messages would be logged before the exit

## Reasoning
Thoroughly testing error handling is crucial for ensuring the application behaves predictably and provides useful feedback when things go wrong. The two scenarios (non-existent directory and permission issues) represent the most common real-world error cases users might encounter when using the `-output` flag.

This approach lets us test the error handling logic in isolation without having our tests terminate due to `os.Exit()` calls. By focusing on the error detection rather than the actual termination of the program, we can verify that the application correctly identifies and reports the error conditions.

An alternative approach would be to refactor the main code to make it more testable (e.g., by returning errors instead of calling `os.Exit()`), but that would require more invasive changes to the codebase. The chosen approach minimizes changes to the production code while still effectively testing the error handling behavior.