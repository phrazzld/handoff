# Todo - Remediation Plan Implementation

## Repository Hygiene
- [x] **T001 · Bugfix · P0: unstage committed binary executable**
    - **Context:** A compiled binary file is committed to the repository
    - **Action:**
        1. Execute `git rm --cached tools/coverage-check/coverage-check`
    - **Done-when:**
        1. `git status` shows the binary as deleted (unstaged)
        2. The binary file still exists locally but is not tracked by Git
    - **Depends-on:** none

- [ ] **T002 · Chore · P0: ignore binary executable file pattern**
    - **Context:** Prevent binary files from being committed in the future
    - **Action:**
        1. Add the pattern `tools/coverage-check/coverage-check*` to the root `.gitignore` file
    - **Done-when:**
        1. The pattern exists in `.gitignore`
        2. `git status --ignored` shows the local binary matching the ignore rule
    - **Depends-on:** none

- [ ] **T003 · Chore · P0: build coverage-check binary in CI**
    - **Context:** Ensure CI workflow can still use the coverage-check tool
    - **Action:**
        1. Identify the CI workflow step(s) that use `tools/coverage-check/coverage-check`
        2. Add a `go build` command before these steps to compile the tool into `tools/coverage-check/`
    - **Done-when:**
        1. CI pipeline successfully builds the `coverage-check` binary
        2. CI pipeline successfully executes the coverage check using the built binary
    - **Depends-on:** [T001, T002]

- [ ] **T004 · Chore · P3: evaluate removing binary from git history**
    - **Context:** Potential repository size optimization
    - **Action:**
        1. Assess repository size and impact of the committed binary on history
        2. Determine if using `git-filter-repo` (or similar) is necessary and safe
    - **Done-when:**
        1. A decision (yes/no) is made and documented on whether to rewrite history
        2. (If yes) A separate plan/ticket is created for the history rewrite
    - **Depends-on:** [T001, T002, T003]

## Configuration Refactor
- [ ] **T005 · Refactor · P1: implement functional options for configuration**
    - **Context:** Current configuration requires manual call to `ProcessConfig()` after setting options
    - **Action:**
        1. Define `type Option func(*Config)`, modify `NewConfig` to accept `...Option`
        2. Create option functions (e.g., `WithInclude`, `WithExclude`) that handle string parsing internally
        3. Make original string config fields unexported and remove the `ProcessConfig()` method
    - **Done-when:**
        1. `Config` struct uses unexported fields for settings configured by options
        2. `ProcessConfig()` method is removed
        3. All internal usages are updated; `go test ./...` passes
    - **Depends-on:** none

- [ ] **T006 · Chore · P2: update configuration documentation for functional options**
    - **Context:** Current documentation refers to old configuration approach
    - **Action:**
        1. Update all documentation to reflect the new functional options pattern
        2. Remove any references to the old `ProcessConfig()` method
    - **Done-when:**
        1. Documentation accurately describes how to configure using functional options
        2. Examples in documentation use the new pattern
    - **Depends-on:** [T005]

## API Layering & Consistency
- [ ] **T007 · Refactor · P2: consolidate API entry point to ProcessProject**
    - **Context:** The library exports multiple processing functions with overlapping responsibilities
    - **Action:**
        1. Identify helper functions that should be internal implementation details
        2. Unexport these helper functions
        3. Update any internal calls to use the now-unexported functions
    - **Done-when:**
        1. Helper functions are unexported
        2. `ProcessProject` remains the primary exported function
        3. `go test ./...` passes
    - **Depends-on:** none

- [ ] **T008 · Chore · P2: update documentation for consolidated API**
    - **Context:** Documentation needs to reflect API changes from T007
    - **Action:**
        1. Update documentation to emphasize `ProcessProject` as the main entry point
        2. Remove documentation for functions that were unexported in T007
    - **Done-when:**
        1. Public API documentation focuses on `ProcessProject`
        2. Documentation for unexported helpers is removed or marked internal
    - **Depends-on:** [T007]

- [ ] **T009 · Refactor · P2: unexport ProcessFile and related low-level functions**
    - **Context:** `ProcessFile` is exported while related functions aren't, creating inconsistency
    - **Action:**
        1. Unexport the `ProcessFile` function
        2. Identify and unexport any other similar low-level processing functions
    - **Done-when:**
        1. `ProcessFile` is unexported
        2. Other related low-level functions are unexported
        3. `go test ./...` passes
    - **Depends-on:** [T007]

- [ ] **T010 · Chore · P2: update documentation for consistent API surface**
    - **Context:** Documentation needs to reflect API changes from T009
    - **Action:**
        1. Remove documentation for `ProcessFile` and any other functions unexported in T009
        2. Ensure documentation clearly distinguishes between the main API and utility functions
    - **Done-when:**
        1. Documentation accurately reflects the minimal, consistent exported API
    - **Depends-on:** [T009]

## File Writing
- [ ] **T011 · Feature · P1: add overwrite control to WriteToFile**
    - **Context:** WriteToFile currently always overwrites existing files
    - **Action:**
        1. Change signature to `WriteToFile(content, path string, overwrite bool) error`
        2. Implement logic to check for file existence and return error if exists and overwrite is false
        3. Update tests and documentation for the new parameter
    - **Done-when:**
        1. Function correctly prevents or allows overwriting based on the flag
        2. Tests verify both behaviors
        3. Documentation is updated
    - **Depends-on:** none

- [ ] **T012 · Refactor · P1: update CLI to use WriteToFile overwrite control**
    - **Context:** CLI needs to pass appropriate overwrite flag to WriteToFile
    - **Action:**
        1. Locate CLI code calling `WriteToFile`
        2. Pass the value of the `-force` flag as the `overwrite` parameter
    - **Done-when:**
        1. CLI correctly respects `-force` flag when writing output files
        2. Manual testing confirms expected behavior
    - **Depends-on:** [T011]

## Performance Optimization
- [ ] **T013 · Refactor · P2: optimize file counting by pre-calculating candidates**
    - **Context:** Current implementation may re-scan directories multiple times
    - **Action:**
        1. Modify `ProcessPaths` to perform file discovery once upfront
        2. Calculate `stats.FilesTotal` based on this initial list
        3. Refactor processing logic to iterate over the pre-discovered list
    - **Done-when:**
        1. Directory scanning happens only once per provided path
        2. `stats.FilesTotal` accurately reflects the candidate file count
        3. `go test ./...` passes
    - **Depends-on:** none

## Coverage Tooling
- [ ] **T014 · Refactor · P1: parse coverage profiles using Go's cover package**
    - **Context:** Current coverage tool relies on parsing `go tool cover -func` output
    - **Action:**
        1. Add `golang.org/x/tools/cover` as a dependency
        2. Rewrite parsing logic to use `cover.ParseProfiles` directly
        3. Remove code related to executing and parsing `go tool cover -func`
    - **Done-when:**
        1. Tool uses `golang.org/x/tools/cover` instead of command output parsing
        2. Tool produces the same correct results as before
        3. Tests pass
    - **Depends-on:** none

## Documentation Updates
- [ ] **T015 · Chore · P2: update README examples for API changes**
    - **Context:** Examples in README need updating for API changes
    - **Action:**
        1. Update all Go code examples in README.md
        2. Ensure examples use correct signatures and functional options
        3. Remove references to deleted examples or outdated concepts
    - **Done-when:**
        1. Examples compile and accurately reflect the current API
        2. Examples demonstrate the functional options pattern
    - **Depends-on:** [T005, T007, T009]

- [ ] **T016 · Chore · P2: update lib/doc.go examples for API changes**
    - **Context:** Examples in package documentation need updating
    - **Action:**
        1. Update all Go code examples in lib/doc.go
        2. Ensure examples use correct signatures and functional options
        3. Remove references to deleted examples or outdated concepts
    - **Done-when:**
        1. Examples in GoDoc compile and accurately reflect current API
        2. Examples demonstrate the functional options pattern
    - **Depends-on:** [T005, T007, T009]

## Testing Improvements
- [ ] **T017 · Test · P2: refactor git availability tests using dependency injection**
    - **Context:** Tests checking for git availability modify global state
    - **Action:**
        1. Identify tests interacting with git commands
        2. Refactor to allow injecting a mock or interface for git interactions
    - **Done-when:**
        1. Tests no longer rely on actual git executable or global state
        2. Tests use mocks or fakes for git interactions
        3. Tests pass
    - **Depends-on:** none

- [ ] **T018 · Test · P2: replace string replacement with filepath.ToSlash for path normalization**
    - **Context:** Tests use `strings.ReplaceAll` for path normalization
    - **Action:**
        1. Find test code using string replacement for path separators
        2. Replace with `filepath.ToSlash` for canonical path normalization
    - **Done-when:**
        1. `strings.ReplaceAll` is no longer used for path separator normalization
        2. `filepath.ToSlash` is used where appropriate
        3. Tests pass
    - **Depends-on:** none

- [ ] **T019 · Test · P2: improve error type checking in clipboard tests**
    - **Context:** Clipboard tests rely on error message strings
    - **Action:**
        1. Review tests involving clipboard operations that check for errors
        2. Use `errors.Is` or `errors.As` instead of string matching
    - **Done-when:**
        1. Clipboard tests use robust error type checking
        2. Tests pass
    - **Depends-on:** none

- [ ] **T020 · Test · P2: improve robustness of error message checking**
    - **Context:** Tests that assert specific error message strings are brittle
    - **Action:**
        1. Review tests that assert exact error message strings
        2. Check for error types or key substrings rather than exact matches
    - **Done-when:**
        1. Error message assertions are less sensitive to minor wording changes
        2. Tests pass
    - **Depends-on:** none

## Miscellaneous Cleanup
- [ ] **T021 · Chore · P3: align Go versions across modules**
    - **Context:** Different modules specify different Go versions
    - **Action:**
        1. Check `go.mod` files in the repository
        2. Ensure the specified Go version is consistent across all modules
    - **Done-when:**
        1. All `go.mod` files declare the same Go version
    - **Depends-on:** none

- [ ] **T022 · Chore · P3: verify and clean up .gitignore entries**
    - **Context:** .gitignore may contain outdated entries
    - **Action:**
        1. Review entries in `.gitignore`
        2. Remove any entries that are no longer relevant
        3. Ensure necessary entries are present
    - **Done-when:**
        1. `.gitignore` file contains only relevant patterns
    - **Depends-on:** [T002]

- [ ] **T023 · Chore · P3: update CI workflow to use standard Go version format**
    - **Context:** CI configuration may use non-standard Go version format
    - **Action:**
        1. Review CI configuration files
        2. Update to use standard `x.y.z` format for Go versions
    - **Done-when:**
        1. CI configuration uses the `x.y.z` Go version format
        2. CI pipeline runs successfully with the updated format
    - **Depends-on:** none

- [ ] **T024 · Chore · P3: update BACKLOG.md regarding verbose flag status**
    - **Context:** BACKLOG.md has contradictory information about verbose flag
    - **Action:**
        1. Edit `BACKLOG.md`
        2. Clarify the current status of the verbose flag
    - **Done-when:**
        1. `BACKLOG.md` accurately reflects the status of the verbose flag
    - **Depends-on:** none
