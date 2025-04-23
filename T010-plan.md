# T010: Update documentation for consistent API surface

## Objective
Update all documentation to reflect API changes from T009, removing references to ProcessFile (which is already unexported), and ensure documentation clearly distinguishes between the main API and utility functions.

## Analysis
Based on T009, we've confirmed that ProcessFile was already unexported as `processFile`. Our task is to update documentation to ensure it reflects the current API surface, with ProcessProject as the main entry point and all other processing functions properly marked as internal.

## Approach
1. Review all documentation files for references to ProcessFile or other unexported functions
2. Update any mentions of ProcessFile to reflect that it's an internal function
3. Ensure documentation clearly distinguishes between the public API and utility functions
4. Update any remaining code examples to match the current API

## Implementation Plan
1. Check `lib/doc.go` for any references to ProcessFile
2. Check `lib/README.md` for any remaining references to ProcessFile
3. Update documentation to emphasize ProcessProject as the main API entry point
4. Ensure all code examples in documentation use only exported functions
5. Run tests to verify everything still works correctly

## Testing
After implementation:
1. Run tests with `go test ./...` to ensure all tests still pass
2. Verify documentation is consistent and accurately represents the API
