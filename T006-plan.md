# T006: Update configuration documentation for functional options

## Objective
Update all documentation to reflect the new functional options pattern for configuration and remove or update references to the old `ProcessConfig()` method.

## Approach
1. Identify all documentation files that mention configuration or the `ProcessConfig()` method
2. Update these files to emphasize the functional options pattern while noting the deprecated method is maintained for backward compatibility
3. Update examples in the documentation to use the new pattern
4. Ensure documentation accurately reflects how to configure using functional options

## Implementation Plan
1. Review the current documentation:
   - lib/doc.go
   - lib/README.md
   - README.md
   - Any other documentation files

2. For each documentation file:
   - Update or rewrite sections that describe configuration to emphasize functional options
   - Update code examples to use the new pattern 
   - Note that `ProcessConfig()` is deprecated but maintained for backward compatibility

3. Run tests to ensure documentation examples compile correctly

4. Verify all documentation is consistent and accurate
