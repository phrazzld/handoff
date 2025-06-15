# Pre-Merge TODO: --ignore-gitignore Flag

## Critical Items (Must Complete Before Merge)

- [x] Update flag help text in `main.go` to clarify behavior and show default value
- [x] Change flag description from `"Process files even if they are gitignored"` to `"Process files even if they are gitignored (bypasses .gitignore rules; default: false)"`
- [x] Test updated help text displays correctly with `handoff --help`
- [x] Verify no regression in existing CLI behavior with updated help text

## Code Quality Improvements (Should Complete Before Merge)

- [x] Add inline code comment above `processFile` gitignore check explaining the bypass logic
- [x] Consider adding verbose log message when `--ignore-gitignore` is active and files are being processed despite gitignore status
- [x] Verify all existing tests pass with the enhanced flag description

## Documentation Updates (Nice to Have Before Merge)

- [x] Add usage example to README.md showing `--ignore-gitignore` flag usage
- [x] Document the flag in CLAUDE.md build commands section
- [x] Add brief mention of the flag's purpose in the main project description

## Post-Merge Enhancements (Future Work)

- [ ] Add integration test specifically for `--ignore-gitignore` functionality
- [ ] Consider adding debug output showing count of gitignored files processed when flag is used
- [ ] Evaluate adding more granular gitignore bypass options in future versions
