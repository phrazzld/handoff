# T015: Update README examples for API changes

## Objective
Update the examples in README.md to reflect API changes made in T005 (functional options for configuration), T007 (consolidated API entry point to ProcessProject), and T009 (unexported ProcessFile).

## Analysis
The README contains code examples that need to be verified and updated to:
1. Use the functional options pattern instead of directly setting configuration fields and calling ProcessConfig()
2. Reflect that ProcessProject is the main entry point (no direct references to ProcessFile)
3. Use the correct function signatures and parameters, especially for WriteToFile which now requires an overwrite parameter

## Approach
1. Carefully review all code examples in README.md
2. Update examples to use the functional options pattern
3. Verify all function calls use the correct signatures
4. Remove any references to unexported functions
5. Ensure examples are formatted consistently and clearly demonstrate current best practices

## Implementation Plan
1. Locate all code examples in README.md
2. Check and update each example to:
   - Use WithXxx() functional options instead of direct field assignment and ProcessConfig()
   - Use correct function signatures
   - Remove references to any unexported functions
3. Verify the examples accurately represent the current API
4. Ensure the examples demonstrate the functional options pattern clearly

## Testing
- Examples should compile if copied into a Go file
- Examples should follow idiomatic Go practices
- Examples should demonstrate current best practices for using the library
