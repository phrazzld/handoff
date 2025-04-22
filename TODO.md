# Todo

## Library API Improvements

- [x] **T001 · Bugfix · P0: correct `ProcessorFunc` signature example in README**
    - **Context:** cr-01 Fix Incorrect `ProcessorFunc` Signature in README
    - **Action:**
        1. Edit `lib/README.md` Advanced Usage section.
        2. Update example signature from `func(content string, filePath string)` to `func(filePath string, content []byte)`.
        3. Update example code to match the new signature using `string(content)` if needed.
    - **Done-when:**
        1. Example signature matches actual type definition in `lib/handoff.go`.
        2. Example code is valid Go that would compile if used.
    - **Depends-on:** none

- [x] **T002 · Refactor · P1: define Stats struct in library**
    - **Context:** cr-07 Refactor `ProcessProject` to Return Stats, Not Log - Step 1
    - **Action:**
        1. Add `type Stats struct { FilesProcessed, FilesTotal, Lines, Chars, Tokens int }` to `lib/handoff.go`.
        2. Add appropriate documentation comments for the type and fields.
    - **Done-when:**
        1. Stats struct is defined in `lib/handoff.go` with proper documentation.
    - **Depends-on:** none

- [x] **T003 · Refactor · P1: update `ProcessPaths` to return Stats**
    - **Context:** cr-07 Refactor `ProcessProject` to Return Stats, Not Log - Step 2
    - **Action:**
        1. Modify `ProcessPaths` function signature to return `(string, Stats, error)`.
        2. Update internal logic to populate Stats struct with file counts and content stats.
        3. Remove any internal statistics logging code.
    - **Done-when:**
        1. `ProcessPaths` successfully returns a populated Stats struct.
        2. No direct logging of statistics occurs in this function.
    - **Depends-on:** T002

- [x] **T004 · Refactor · P1: update `ProcessProject` to return Stats**
    - **Context:** cr-07 Refactor `ProcessProject` to Return Stats, Not Log - Step 3/4
    - **Action:**
        1. Modify `ProcessProject` function signature to return `(string, Stats, error)`.
        2. Update it to receive and return `Stats` from `ProcessPaths`.
        3. Remove all statistics logging code from `ProcessProject`.
    - **Done-when:**
        1. `ProcessProject` returns content, Stats struct, and error.
        2. No direct logging of statistics occurs in this function.
    - **Depends-on:** T003

- [x] **T005 · Feature · P2: implement directory creation in `WriteToFile`**
    - **Context:** cr-04 Align `WriteToFile` Docs/Implementation (Dir Creation)
    - **Action:**
        1. Edit `WriteToFile` function in `lib/handoff.go`.
        2. Add directory creation logic using `filepath.Dir` and `os.MkdirAll`.
        3. Ensure proper error handling for directory creation failures.
    - **Done-when:**
        1. `WriteToFile` successfully creates parent directories before writing files.
        2. Tests pass, including added tests for this feature.
    - **Depends-on:** none

- [x] **T006 · Test · P2: add test for directory creation in `WriteToFile`**
    - **Context:** cr-04 Align `WriteToFile` Docs/Implementation (Dir Creation) - Step 6
    - **Action:**
        1. Add test case to `lib/handoff_test.go` for `WriteToFile` creating directories.
        2. Create test with non-existent nested path and verify file is written properly.
        3. Clean up test directories after test.
    - **Done-when:**
        1. Test passes, showing `WriteToFile` creates required directories.
    - **Depends-on:** T005

- [x] **T007 · Docs · P2: remove examples using unexported functions from README**
    - **Context:** cr-06 Remove/Update Doc Examples Referencing Unexported Functions
    - **Action:**
        1. Edit `lib/README.md`.
        2. Remove sections showing examples for `GetFilesFromDir`, `ShouldProcess`, `ProcessFile`, and `IsBinaryFile`.
        3. Ensure any remaining examples only use exported functions.
    - **Done-when:**
        1. `lib/README.md` contains only examples of exported functions.
        2. Advanced Usage section is revised to show only appropriate public API usage.
    - **Depends-on:** T001

## CLI Improvements

- [x] **T008 · Refactor · P1: update main to use returned Stats struct**
    - **Context:** cr-07/08 Refactor `ProcessProject` and Fix Fragile File Count Calc - Steps 5/2-4
    - **Action:**
        1. Update `main` function in `main.go` to capture the Stats struct returned by `handoff.ProcessProject`.
        2. Pass the Stats struct to a refactored `logStatisticsUsingLib` function.
    - **Done-when:**
        1. `main.go` captures and uses the Stats struct returned by the library.
        2. Tests pass.
    - **Depends-on:** T004

- [x] **T009 · Refactor · P1: refactor CLI statistics logging**
    - **Context:** cr-08/09 Update logStatisticsUsingLib to use Stats
    - **Action:**
        1. Modify `logStatisticsUsingLib` function to accept Stats struct.
        2. Remove `strings.Count` calculation for processed files.
        3. Use Stats fields directly for all statistics reporting.
    - **Done-when:**
        1. `logStatisticsUsingLib` no longer calls `handoff.CalculateStatistics`.
        2. Function uses only the Stats struct for reporting.
    - **Depends-on:** T008

- [x] **T010 · Test · P2: remove redundant CLI tests**
    - **Context:** cr-03 Remove Redundant CLI Tests from Unit Tests
    - **Action:**
        1. Review and identify redundant tests in `handoff_test.go`.
        2. Verify coverage for these functions exists in library or integration tests.
        3. Remove the redundant test functions.
    - **Done-when:**
        1. Redundant tests are removed from `handoff_test.go`.
        2. All test coverage is maintained through other tests.
        3. `go test ./...` passes with no errors.
    - **Depends-on:** none

- [x] **T011 · Test · P2: remove duplicate test helper functions**
    - **Context:** cr-05 Remove Duplicate Test Helper Functions
    - **Action:**
        1. Identify helper functions duplicated in `cli_integration_test.go` and main code.
        2. Refactor tests to rely on CLI binary execution instead of reimplementing logic.
        3. Remove the duplicated helper functions.
    - **Done-when:**
        1. Duplicated helpers are removed.
        2. Tests still pass and functionality is maintained.
    - **Depends-on:** none

- [ ] **T012 · Chore · P3: remove unused function `processPathUsingLib`**
    - **Context:** "The `processPathUsingLib` helper function is unused after the refactoring." (cr code review)
    - **Action:**
        1. Remove the unused `processPathUsingLib` function from `main.go`.
    - **Done-when:**
        1. Function is removed from codebase.
        2. Code still compiles and tests pass.
    - **Depends-on:** none

## Examples Update

- [ ] **T013 · Refactor · P2: update `gemini_planner.go` example to use Stats**
    - **Context:** cr-07 Update examples that call `ProcessProject` - Step 6
    - **Action:**
        1. Update `examples/gemini_planner.go` to use the new `ProcessProject` return values.
        2. Replace any calls to `CalculateStatistics` with values from the returned Stats.
    - **Done-when:**
        1. Example compiles with updated library API.
        2. Example uses Stats struct for statistics instead of recalculating.
    - **Depends-on:** T004

- [ ] **T014 · Refactor · P2: update `simple_usage.go` example to use Stats**
    - **Context:** cr-07 Update examples that call `ProcessProject` - Step 6
    - **Action:**
        1. Update `examples/simple_usage.go` to use the new `ProcessProject` return values.
        2. Replace any calls to `CalculateStatistics` with values from the returned Stats.
    - **Done-when:**
        1. Example compiles with updated library API.
        2. Example uses Stats struct for statistics instead of recalculating.
    - **Depends-on:** T004

## CI Improvements

- [ ] **T015 · Feature · P2: create Go coverage checker tool**
    - **Context:** cr-02 Replace Brittle CI Shell Script - Steps 1-3
    - **Action:**
        1. Create directory `tools/coverage-check`.
        2. Create a simple Go program to parse coverage output and check against threshold.
        3. Add unit tests for the coverage checker.
    - **Done-when:**
        1. Coverage checker tool works correctly with various inputs.
        2. Tests for the tool pass.
    - **Depends-on:** none

- [ ] **T016 · CI · P2: update CI workflow to use coverage checker**
    - **Context:** cr-02 Replace Brittle CI Shell Script - Steps 4-5
    - **Action:**
        1. Update `.github/workflows/test-coverage.yml` to build and use coverage checker.
        2. Replace shell script block with Go tool invocation.
    - **Done-when:**
        1. CI workflow uses the new Go tool for coverage checking.
        2. CI passes/fails correctly based on coverage thresholds.
    - **Depends-on:** T015

- [ ] **T017 · CI · P3: update Go version in CI workflow**
    - **Context:** cr-10 Update Go Version in CI Workflow
    - **Action:**
        1. Edit `.github/workflows/test-coverage.yml` to set go-version to `1.22.3`.
    - **Done-when:**
        1. CI workflow uses Go 1.22.3.
    - **Depends-on:** none

## Documentation

- [ ] **T018 · Docs · P3: standardize Go version across documentation**
    - **Context:** cr-11 Standardize Go Version Info Across Project
    - **Action:**
        1. Verify `go.mod` specifies Go 1.22.3.
        2. Update README.md prerequisites section to specify Go 1.22.3.
        3. Check for any other version references that may need updating.
    - **Done-when:**
        1. All project documentation consistently references Go 1.22.3.
    - **Depends-on:** T017

### Clarifications & Assumptions
- [ ] **Issue:** Determine specific test files to remove from handoff_test.go
    - **Context:** cr-03 lists `TestFileCreation`, `TestFileOverwriteProtection`, `TestInvalidPathErrorHandling` as redundant, but we should confirm with full source review
    - **Blocking?:** no