# Todo

## Core Logic Consolidation (lib)
- [x] **T001 · Refactor · P1: review logic in main/files.go and main/output.go**
    - **Context:** PLAN.md § 2.A.1
    - **Action:**
        1. Identify all logic currently residing in `main/files.go` and `main/output.go`.
        2. Determine which parts are core processing logic vs. CLI-specific concerns.
    - **Done‑when:**
        1. A clear map of functions/logic and their destination (lib or main) exists.
    - **Depends‑on:** none
- [x] **T002 · Refactor · P1: move core file/output logic from main to lib/handoff.go**
    - **Context:** PLAN.md § 2.A.1
    - **Action:**
        1. Relocate identified core logic (file collection, filtering, formatting) from `main/files.go` and `main/output.go` into `lib/handoff.go` or new files within `lib`.
        2. Adapt function signatures and dependencies as needed for the library context.
    - **Done‑when:**
        1. Core file processing and formatting logic resides entirely within the `lib` package.
        2. Code compiles successfully.
    - **Depends‑on:** [T001]
- [x] **T003 · Refactor · P1: refine exported vs unexported identifiers in lib package**
    - **Context:** PLAN.md § 2.A.2, § 5
    - **Action:**
        1. Review all functions, types, constants, and variables in the `lib` package.
        2. Ensure only necessary identifiers intended for public use are exported.
        3. Confirm API aligns with the example structure (e.g., `Config`, `ProcessProject`, `WriteToFile`, `CalculateStatistics`).
    - **Done‑when:**
        1. The public API surface of the `lib` package is intentionally defined and minimized.
        2. Unnecessary identifiers are unexported.
    - **Depends‑on:** [T002]
- [x] **T004 · Chore · P2: add/update godoc comments for exported lib identifiers**
    - **Context:** PLAN.md § 2.A.2, § 2.D.1
    - **Action:**
        1. Write clear, concise GoDoc comments for all exported functions, types, constants, and variables in the `lib` package.
        2. Explain purpose, parameters, return values, and any nuances.
    - **Done‑when:**
        1. All exported identifiers in `lib` have comprehensive GoDoc comments.
        2. `go doc ./lib/...` shows complete documentation.
    - **Depends‑on:** [T003]
- [x] **T005 · Chore · P2: add package-level documentation for lib package**
    - **Context:** PLAN.md § 2.A.2
    - **Action:**
        1. Create or update `lib/doc.go` with package-level documentation explaining the library's purpose and usage overview.
    - **Done‑when:**
        1. `go doc ./lib` displays clear package documentation.
    - **Depends‑on:** [T003]

## CLI Refactoring (main)
- [x] **T006 · Refactor · P1: modify main.parseConfig to return *lib.Config**
    - **Context:** PLAN.md § 2.B.1
    - **Action:**
        1. Update the `parseConfig` function (or equivalent) in `main.go` to populate and return a `*lib.Config` struct.
        2. Ensure flag parsing correctly maps to `lib.Config` fields.
    - **Done‑when:**
        1. `parseConfig` returns a `*lib.Config` instance based on CLI flags.
    - **Depends‑on:** [T003]
- [x] **T007 · Refactor · P1: update main.go to use lib.ProcessProject for core logic**
    - **Context:** PLAN.md § 2.B.1
    - **Action:**
        1. Import the `lib` package in `main.go`.
        2. Replace direct file processing/formatting calls with a call to `lib.ProcessProject`, passing the `lib.Config` and paths.
    - **Done‑when:**
        1. `main.go` uses `lib.ProcessProject` to get the final formatted content string.
        2. CLI execution flow utilizes the library for core work.
    - **Depends‑on:** [T003, T006]
- [x] **T008 · Refactor · P1: update main.go to use lib.WriteToFile for file output**
    - **Context:** PLAN.md § 2.B.1, § 5
    - **Action:**
        1. If file output logic exists in `main.go`, replace it with a call to `lib.WriteToFile` (assuming `WriteToFile` is part of the lib API as per §5).
        2. Pass the content from `lib.ProcessProject` and the output file path.
    - **Done‑when:**
        1. CLI output file writing uses the `lib.WriteToFile` function.
    - **Depends‑on:** [T003, T007]
- [x] **T009 · Refactor · P1: verify clipboard handling remains functional in main.go**
    - **Context:** PLAN.md § 2.B.1
    - **Action:**
        1. Ensure the logic for copying output to the clipboard remains within `main.go`.
        2. Verify it correctly uses the content string returned by `lib.ProcessProject`.
    - **Done‑when:**
        1. Clipboard functionality works as expected using the library output.
    - **Depends‑on:** [T007]
- [x] **T010 · Refactor · P2: remove redundant main/files.go and main/output.go**
    - **Context:** PLAN.md § 2.A.1
    - **Action:**
        1. Confirm all necessary logic from `files.go` and `output.go` has been moved to `lib` or is handled by the refactored `main.go`.
        2. Delete the `files.go` and `output.go` files from the main package.
    - **Done‑when:**
        1. `files.go` and `output.go` are removed from the project root/main package.
        2. The project compiles and runs without them.
    - **Depends‑on:** [T002, T007, T008, T009, T012]

## Testing (lib & main)
- [x] **T011 · Test · P1: create lib/handoff_test.go**
    - **Context:** PLAN.md § 2.C.1
    - **Action:**
        1. Create the test file `lib/handoff_test.go`.
    - **Done‑when:**
        1. `lib/handoff_test.go` exists.
    - **Depends‑on:** [T002]
- [x] **T012 · Test · P1: migrate relevant tests from root handoff_test.go to lib/handoff_test.go**
    - **Context:** PLAN.md § 2.C.1
    - **Action:**
        1. Identify tests in the root `handoff_test.go` that primarily test logic now residing in `lib`.
        2. Move these tests to `lib/handoff_test.go`, adapting them to test the library's public API directly.
    - **Done‑when:**
        1. Library-specific tests are moved to `lib/handoff_test.go`.
        2. Moved tests pass against the library code.
    - **Depends‑on:** [T011]
- [x] **T013 · Test · P1: add unit tests for new/existing lib helper functions**
    - **Context:** PLAN.md § 2.C.1, § 3
    - **Action:**
        1. Write unit tests for helper functions within `lib` (e.g., `NewConfig`, `Config.ProcessConfig`, `CalculateStatistics`, any internal helpers).
        2. Focus on testing individual functions in isolation.
    - **Done‑when:**
        1. Core library helper functions have adequate unit test coverage.
        2. Tests pass.
    - **Depends‑on:** [T011]
- [ ] **T014 · Test · P1: add integration tests for lib.ProcessProject**
    - **Context:** PLAN.md § 2.C.1, § 3
    - **Action:**
        1. Write integration tests for `lib.ProcessProject` in `lib/handoff_test.go`.
        2. Test various `lib.Config` combinations (include/exclude patterns, formats).
        3. Use sample file structures for testing input.
    - **Done‑when:**
        1. `lib.ProcessProject` behavior is verified with different configurations.
        2. Tests pass.
    - **Depends‑on:** [T011]
- [ ] **T015 · Test · P1: refactor root handoff_test.go for cli-specific tests**
    - **Context:** PLAN.md § 2.C.2, § 3
    - **Action:**
        1. Update the root `handoff_test.go` to focus solely on testing CLI behavior.
        2. Remove tests that were migrated to `lib/handoff_test.go`.
        3. Ensure tests cover flag parsing, argument handling, and exit codes.
    - **Done‑when:**
        1. Root `handoff_test.go` contains only CLI-level tests.
        2. Existing relevant CLI tests pass after refactoring.
    - **Depends‑on:** [T007, T008, T012]
- [ ] **T016 · Test · P1: add cli integration tests verifying library usage**
    - **Context:** PLAN.md § 2.C.2, § 3
    - **Action:**
        1. Add tests in root `handoff_test.go` that execute the compiled CLI binary.
        2. Verify end-to-end behavior, ensuring flags correctly influence the output produced via the library calls (checking final output, clipboard, file writing).
    - **Done‑when:**
        1. CLI integration tests confirm correct interaction between `main` and `lib`.
        2. All CLI functionality (flags, output, filtering) is verified via tests.
        3. Tests pass.
    - **Depends‑on:** [T015]
- [ ] **T017 · Test · P2: configure and enforce test coverage check for lib package**
    - **Context:** PLAN.md § 2.C.1
    - **Action:**
        1. Configure CI (e.g., GitHub Actions) to calculate test coverage for the `lib` package.
        2. Set a minimum coverage threshold (e.g., 85%).
        3. Ensure the CI job fails if coverage drops below the threshold.
    - **Done‑when:**
        1. Test coverage for the `lib` package is reported in CI.
        2. CI enforces the minimum coverage threshold.
    - **Depends‑on:** [T012, T013, T014]

## Documentation & Examples
- [ ] **T018 · Chore · P2: update main README.md for library/cli structure**
    - **Context:** PLAN.md § 2.D.1
    - **Action:**
        1. Revise the main `README.md`.
        2. Clearly explain that the project provides both a library (`lib`) and a CLI tool.
        3. Update installation/usage instructions accordingly.
    - **Done‑when:**
        1. Main `README.md` accurately reflects the project's structure and usage for both library and CLI consumers.
    - **Depends‑on:** [T003, T007]
- [ ] **T019 · Chore · P2: create/update lib/README.md with api usage guide**
    - **Context:** PLAN.md § 2.D.1
    - **Action:**
        1. Create or update `lib/README.md`.
        2. Provide a detailed guide on how to import and use the `lib` package's public API.
        3. Include code snippets demonstrating common use cases.
    - **Done‑when:**
        1. `lib/README.md` exists and provides clear instructions for library users.
    - **Depends‑on:** [T004, T005]
- [ ] **T020 · Chore · P2: update examples/gemini_planner.go to use library api**
    - **Context:** PLAN.md § 2.D.2
    - **Action:**
        1. Review the existing `examples/gemini_planner.go`.
        2. Update it to import and use the new `lib` package API instead of any previous methods.
    - **Done‑when:**
        1. `examples/gemini_planner.go` correctly demonstrates usage of the refactored library.
        2. The example compiles and runs.
    - **Depends‑on:** [T003]
- [ ] **T021 · Chore · P2: add simple library usage example**
    - **Context:** PLAN.md § 2.D.2
    - **Action:**
        1. Create a new, minimal example file (e.g., `examples/simple_usage.go`).
        2. Demonstrate the basic steps: create config, call `ProcessProject`, handle output/error.
    - **Done‑when:**
        1. A simple, clear example of library usage exists in the `examples` directory.
        2. The example compiles and runs.
    - **Depends‑on:** [T003]

## Go Modules & Versioning
- [ ] **T022 · Chore · P2: verify go.mod module path**
    - **Context:** PLAN.md § 2.E.1
    - **Action:**
        1. Ensure the `module` directive in `go.mod` is set correctly to `github.com/phrazzld/handoff`.
    - **Done‑when:**
        1. `go.mod` specifies the correct canonical module path.
    - **Depends‑on:** none
- [ ] **T023 · Chore · P2: run go mod tidy**
    - **Context:** PLAN.md § 2.E.1
    - **Action:**
        1. Run `go mod tidy` in the project root.
        2. Commit any changes to `go.mod` and `go.sum`.
    - **Done‑when:**
        1. Go module dependencies are cleaned up and consistent.
    - **Depends‑on:** [T010, T016] (Run after all code/test changes potentially affecting dependencies)
- [ ] **T024 · Chore · P3: plan initial semantic version tag**
    - **Context:** PLAN.md § 2.E.1
    - **Action:**
        1. Decide on the initial semantic version for the library release (e.g., `v0.1.0` or `v1.0.0`).
        2. Document the decision (e.g., in project planning notes or an issue).
    - **Done‑when:**
        1. An initial SemVer tag strategy is decided upon.
    - **Depends‑on:** none

### Clarifications & Assumptions
- [ ] **Issue:** Define strategy for error handling between library and CLI.
    - **Context:** PLAN.md § 6 (Potential Challenges) - Ensuring proper error handling in the library vs. CLI.
    - **Blocking?:** no (Decision needed before finalizing T007/T008 and related tests)
- [ ] **Issue:** Confirm `WriteToFile` function belongs in `lib` package API.
    - **Context:** PLAN.md § 5 (Example Library API) includes `WriteToFile`. T008 assumes this is correct.
    - **Blocking?:** no (Assumption made, but confirmation preferred)