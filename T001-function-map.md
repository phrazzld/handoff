# Function Mapping for Handoff Package Refactoring

## 1. Functions in files.go

| Function | Purpose | Lib Equivalent | Destination | Notes |
|----------|---------|----------------|-------------|-------|
| `gitAvailable` (var) | Checks for git command | `GitAvailable` (in lib) | Remove | Already imported from lib |
| `isGitIgnored` | Checks if file is git-ignored | `IsGitIgnored` | Remove | Duplicate of lib function |
| `getGitFiles` | Gets files from git repo | `GetGitFiles` | Remove | Duplicate of lib function |
| `getFilesWithFilepathWalk` | Gets files by walking dirs | `GetFilesWithFilepathWalk` | Remove | Duplicate of lib function |
| `getFilesFromDir` | Gets files using git or walking | `GetFilesFromDir` | Remove | Duplicate of lib function |
| `isBinaryFile` | Detects binary files | `IsBinaryFile` | Remove | Duplicate of lib function |
| `isWhitespace` | Checks if byte is whitespace | `isWhitespace` | Remove | Duplicate of lib function |
| `minInt` | Returns minimum of two integers | `minInt` | Remove | Duplicate of lib function |
| `shouldProcess` | Filters files by config | `ShouldProcess` | Remove | Duplicate of lib function |
| `processFile` | Processes single file | `ProcessFile` | Remove | Duplicate of lib function |
| `processDirectory` | Processes directory | `ProcessDirectory` | Remove | Duplicate of lib function |
| `processPathWithProcessor` | Processes path with custom processor | `ProcessPathWithProcessor` | Remove | Duplicate of lib function |
| `processPath` | Processes path with default processor | Not directly in lib | Keep in main | CLI-specific function that forwards to lib functions |

## 2. Functions in output.go

| Function | Purpose | Lib Equivalent | Destination | Notes |
|----------|---------|----------------|-------------|-------|
| `estimateTokenCount` | Counts tokens in text | `EstimateTokenCount` | Remove | Duplicate of lib function |
| `wrapInContext` | Wraps content in context tags | `WrapInContext` | Remove | Duplicate of lib function |
| `logStatistics` | Logs content statistics | Partially in `ProcessProject` | Keep in main | This function is CLI-specific with direct logging |

## 3. Summary of Changes

### Functions to Remove (Already in lib)
- `gitAvailable` - Already imported from lib
- All functions in `files.go` except `processPath`
- `estimateTokenCount` and `wrapInContext` from `output.go`

### Functions to Keep in main Package
- `processPath` - CLI-specific function that forwards to lib functions
- `logStatistics` - CLI-specific function that handles log output

### Other Observations
1. Almost all core logic already exists in the `lib` package
2. Most functions in `files.go` and `output.go` are duplicates of lib functions
3. Only CLI-specific formatting and logging needs to remain in main