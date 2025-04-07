# Add unit test for parseConfig with -output and -force flags

## Goal
Update or add unit tests to verify the `-output` and `-force` flags are correctly parsed and returned by parseConfig.

## Implementation Approach
I'll create a new unit test function named `TestParseConfigOutputAndForceFlags` that will:

1. Set up mock command-line arguments that include various combinations of the `-output` and `-force` flags
2. Call the `parseConfig()` function and capture its return values
3. Verify that the returned values correctly reflect the provided flags

This approach will involve temporarily modifying `os.Args` to simulate command-line arguments, then restoring the original values after the test. I'll test several scenarios:

1. Basic usage with only `-output` flag
2. Using both `-output` and `-force` flags
3. Verifying default values when flags aren't provided

## Reasoning
Unit testing the flag parsing is important to ensure the command-line interface works correctly. The approach of temporarily modifying `os.Args` is a common practice for testing flag parsing in Go, as it allows us to simulate different command-line inputs without actually running the program with those arguments.

I considered alternative approaches such as refactoring `parseConfig()` to accept a custom `flag.FlagSet` for testing, but that would require changing the existing function signature and potentially affect other parts of the code. The chosen approach is non-invasive and matches the existing testing patterns in the codebase.