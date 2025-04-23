# T009: Unexport ProcessFile and related low-level functions

## Objective
Unexport the `ProcessFile` function and identify and unexport any other similar low-level processing functions to maintain API consistency, following the work done in T007.

## Analysis
Based on previous work in T007, several processing functions were unexported:
- `ProcessDirectory` → `processDirectory`
- `ProcessPathWithProcessor` → `processPathWithProcessor`  
- `ProcessPaths` → `processPaths`

However, `ProcessFile` is still exported, creating inconsistency in the API surface.

## Approach
1. Identify all exported low-level functions in `lib/handoff.go` that should be implementation details
2. Unexport these functions by changing their first letter to lowercase
3. Verify no external code relies on these functions
4. Update any internal references to use the unexported names
5. Ensure all tests still pass

## Implementation Plan
1. Unexport `ProcessFile` by renaming it to `processFile`
2. Scan the codebase for other similar low-level functions that should be unexported
3. Check for any internal references to the exported name and update them
4. Run tests to verify everything still works correctly

## Testing
After implementation:
1. Run tests with `go test ./...` to ensure all tests still pass
2. Verify the package still builds correctly
