# T011 Implementation Plan: Add Overwrite Control to WriteToFile

## Context

Currently, the `WriteToFile` function in the `handoff/lib` package always overwrites existing files without any control to prevent accidental overwrites. This can lead to unintended data loss if users aren't careful. Adding an `overwrite` parameter will give users explicit control over this behavior.

## Analysis

The current signature is: `WriteToFile(content, filePath string) error`

This will be changed to: `WriteToFile(content, filePath string, overwrite bool) error`

The implementation will:
1. Check if the file exists
2. If it exists and `overwrite` is false, return an error
3. Otherwise, proceed with writing the file

## Implementation Steps

1. Locate the `WriteToFile` function in the handoff package.
2. Update the function signature to add the `overwrite` bool parameter.
3. Add file existence check logic at the beginning of the function.
4. Return a descriptive error if the file exists and overwrite is false.
5. Update existing tests to work with the new parameter.
6. Add new tests to verify both overwrite true/false behaviors.
7. Update documentation to reflect the new parameter.

## Testing Plan

1. Test with `overwrite=true` on an existing file (should succeed)
2. Test with `overwrite=false` on an existing file (should fail with appropriate error)
3. Test with `overwrite=true` on a new file (should succeed)
4. Test with `overwrite=false` on a new file (should succeed - no existing file to overwrite)
5. Ensure all existing tests pass with the updated signature

## Notes

- Need to use a clear, descriptive error message for the overwrite protection case
- Consider defining a custom error type for this specific error for better error handling
- Ensure backwards compatibility is considered for any existing calls to this function
