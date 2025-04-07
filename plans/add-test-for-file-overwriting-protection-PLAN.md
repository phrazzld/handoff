# Add test for file overwriting protection

## Goal
Create a test to verify that existing files are NOT overwritten when `-output` points to an existing file and `-force` is not specified. Also verify that files ARE overwritten when `-force` is specified.

## Implementation Approach
I'll create a new integration test function named `TestFileOverwriteProtection` that will:

1. Set up a temporary test directory and file path
2. Create an initial output file with known content at the test path
3. Run the handoff command with `-output` flag pointing to the existing file (without `-force`)
4. Verify that:
   - The command returns an error about the existing file
   - The original file content remains unchanged
5. Run the handoff command again with both `-output` and `-force` flags
6. Verify that:
   - The file is successfully overwritten
   - The content is updated to the new formatted output

The test will reuse much of the infrastructure from the existing `TestFileCreation` test but will focus specifically on the overwrite protection behavior. It will test both the negative case (file exists, no force flag) and the positive case (file exists, force flag provided).

## Reasoning
Testing the file overwrite protection is essential as it's a critical safety feature that prevents accidental data loss. Users need to trust that their existing files won't be silently overwritten unless they explicitly use the `-force` flag.

This approach is efficient because it builds on the existing test infrastructure while focusing specifically on the overwrite protection behavior. By testing both the rejection (without `-force`) and the override (with `-force`), we ensure that both paths of the protection mechanism work correctly.

I considered creating separate tests for each case (with and without `-force`), but a combined approach is more efficient as it avoids duplicating the setup code while still clearly testing both behaviors.