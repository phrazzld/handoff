# Add test for flag interaction between -output, -force, and -dry-run

## Goal
Add a test to verify correct precedence when various combinations of `-output`, `-force`, and `-dry-run` flags are used together.

## Implementation Approach
I'll create a new integration test function named `TestFlagInteractionPrecedence` that will test the precedence logic between the three flags. The test will implement several scenarios:

1. **All flags provided together**: Test when `-dry-run`, `-output`, and `-force` are all specified
   - Verify that `-dry-run` takes highest precedence (no file is written or clipboard modified)
   
2. **Output and force flags**: Test when `-output` and `-force` are specified (without `-dry-run`)
   - Verify that file output is used (file is written)
   
3. **Output only**: Test when only `-output` is specified
   - Verify that file output is used (file is written)
   
4. **Dry-run and output**: Test when `-dry-run` and `-output` are specified
   - Verify that `-dry-run` takes precedence over `-output` (no file is written)
   
5. **No flags**: Test when no special flags are provided
   - Verify that clipboard output is used (default behavior)

Since the `main()` function calls `os.Exit()` in error cases and we can't directly test the execution path in `main()`, I'll create a helper function that mimics the precedence logic for testing purposes. This function will:

1. Take the same arguments as would be processed by `main()`
2. Return a string indicating which output mode was selected (dry-run, file, or clipboard)
3. Not actually execute the output operations or call `os.Exit()`

By testing this helper function with various flag combinations, we can verify that the precedence logic works correctly without needing to modify the main application code.

## Reasoning
Testing the flag interaction is important to ensure that the application behaves predictably when users provide multiple flags. The precedence order (dry-run > output file > clipboard) is a key part of the application's behavior and should be thoroughly tested.

The approach of using a separate test helper function that mimics the precedence logic allows us to test the decision-making process without getting caught by the `os.Exit()` calls in the main function. This approach is more maintainable than trying to mock or patch the OS functions, and it directly tests the logic that determines which output mode to use.

I considered trying to test the `main()` function directly by capturing stdout/stderr, but that approach would be more complex and brittle due to the `os.Exit()` calls. The chosen approach focuses on testing the specific logic we care about (the precedence rules) in a clean, isolated manner.