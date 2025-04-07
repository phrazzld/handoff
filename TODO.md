# TODO

## File Output Option Implementation

- [x] **Update parseConfig function signature and return values**
  - **Action:** Modify the function signature in main.go to return the output file path and force flag values in addition to the config and dryRun values.
  - **Depends On:** None
  - **AC Ref:** Task 2 "Update Flag Parsing"

- [x] **Add -output command line flag**
  - **Action:** Define a new string flag `-output` in the flags section of parseConfig function that will accept a file path as input.
  - **Depends On:** "Update parseConfig function signature and return values"
  - **AC Ref:** Task 1 "Add `-output` Flag"

- [x] **Add -force command line flag**
  - **Action:** Define a new boolean flag `-force` in the flags section of parseConfig function that will allow overwriting existing files.
  - **Depends On:** "Add -output command line flag"
  - **AC Ref:** Updated requirement for file overwrite protection

- [ ] **Implement path resolution and validation**
  - **Action:** Add code to ensure the provided output path is correctly resolved to an absolute path using filepath.Abs when the flag is used.
  - **Depends On:** "Add -output command line flag"
  - **AC Ref:** Task 3 "Implement File Writing Logic", Consideration "File Path Resolution"

- [ ] **Implement file existence check**
  - **Action:** Add code to check if the output file already exists before writing to it, and refuse to overwrite unless -force flag is specified.
  - **Depends On:** "Implement path resolution and validation", "Add -force command line flag"
  - **AC Ref:** Task 3 "Implement File Writing Logic", Updated requirement for file overwrite protection

- [ ] **Add file output handling logic in main function**
  - **Action:** Add conditional logic in main() to write content to a file when -output flag is provided, using handoff.WriteToFile function.
  - **Depends On:** "Implement file existence check"
  - **AC Ref:** Task 3 "Implement File Writing Logic"

- [ ] **Implement output precedence logic**
  - **Action:** Ensure proper precedence between -dry-run, -output, and default clipboard behavior. Order should be: Dry Run > Output File > Clipboard.
  - **Depends On:** "Add file output handling logic in main function"
  - **AC Ref:** Task 4 "Adjust Existing Logic", Consideration "Flag Precedence"

- [ ] **Implement error handling for file operations**
  - **Action:** Add appropriate error handling for file writing operations (path invalid, permissions, disk full) with user-friendly error messages.
  - **Depends On:** "Add file output handling logic in main function"
  - **AC Ref:** Task 5 "Error Handling", Consideration "Error Handling"

- [ ] **Update logging for file output operations**
  - **Action:** Add informative log messages for file operations (verbose log for file path, confirmation when writing succeeds).
  - **Depends On:** "Add file output handling logic in main function"
  - **AC Ref:** Task 8 "Final Logging Adjustment"

- [ ] **Ensure statistics are logged regardless of output mode**
  - **Action:** Verify that the statistics summary is always printed to stderr, regardless of whether output goes to clipboard, file, or stdout.
  - **Depends On:** "Update logging for file output operations"
  - **AC Ref:** Task 8 "Final Logging Adjustment", Consideration "Statistics Logging"

- [ ] **Update README.md with new -output and -force flag documentation**
  - **Action:** Add documentation in README.md for the new -output and -force flags, including descriptions, usage examples, and behavior with other flags. Include clear explanation of file overwrite protection.
  - **Depends On:** "Implement output precedence logic"
  - **AC Ref:** Task 6 "Update Documentation"

- [ ] **Add unit test for parseConfig with -output and -force flags**
  - **Action:** Update or add unit tests to verify the -output and -force flags are correctly parsed and returned by parseConfig.
  - **Depends On:** "Add -output command line flag", "Add -force command line flag"
  - **AC Ref:** Task 7 "Add/Update Tests", Testing Strategy "Unit Tests"

- [ ] **Add integration test for file creation**
  - **Action:** Add test to verify file is created with correct content when -output flag is used.
  - **Depends On:** "Add file output handling logic in main function"
  - **AC Ref:** Task 7 "Add/Update Tests", Testing Strategy "Integration/CLI Tests"

- [ ] **Add test for file overwriting protection**
  - **Action:** Create test to verify existing files are NOT overwritten when -output points to an existing file and -force is not specified. Also verify files ARE overwritten when -force is specified.
  - **Depends On:** "Add integration test for file creation"
  - **AC Ref:** Task 7 "Add/Update Tests", Testing Strategy "Integration/CLI Tests", Updated requirement for file overwrite protection

- [ ] **Add test for error handling on invalid paths**
  - **Action:** Create test to verify proper error handling when -output points to an invalid or inaccessible path.
  - **Depends On:** "Implement error handling for file operations"
  - **AC Ref:** Task 7 "Add/Update Tests", Testing Strategy "Integration/CLI Tests"

- [ ] **Add test for flag interaction between -output, -force, and -dry-run**
  - **Action:** Add test to verify correct precedence when various combinations of -output, -force, and -dry-run flags are used together.
  - **Depends On:** "Implement output precedence logic"
  - **AC Ref:** Task 7 "Add/Update Tests", Testing Strategy "Integration/CLI Tests"

- [ ] **Manually test across different operating systems**
  - **Action:** Manually test the -output functionality on different platforms (if possible) to verify cross-platform compatibility.
  - **Depends On:** All implementation tasks
  - **AC Ref:** Testing Strategy "Manual Testing"

## [!] CLARIFICATIONS NEEDED / ASSUMPTIONS

- [ ] **Assumption: Default flag description is sufficient**
  - **Context:** The plan provides a suggestion for the flag description: "Write output to the specified file instead of clipboard (e.g., HANDOFF.md)". Assuming this is sufficient and no additional clarification is needed.

- [ ] **Clarification: Must implement file overwrite protection**
  - **Context:** The plan notes in section 4: "File Overwriting: The current handoff.WriteToFile uses os.WriteFile, which will truncate and overwrite existing files." However, we need to change this behavior to protect existing files by default. If the output file already exists, warn the user and quit unless a -force flag is used.

- [ ] **Assumption: Statistics log format remains unchanged**
  - **Context:** The plan mentions ensuring statistics are still logged, but doesn't specify any changes to the statistics format. Assuming the current format is sufficient.

- [ ] **Clarification: Need to add -force flag for file overwriting**
  - **Context:** The plan doesn't mention additional flags for controlling file behavior, but we need to add a -force flag to allow overwriting of existing files. Without this flag, the tool should refuse to overwrite existing files.

- [ ] **Assumption: Default file permissions (0644) are appropriate**
  - **Context:** The existing WriteToFile function uses default 0644 permissions. Assuming these permissions are appropriate for all platforms and use cases.