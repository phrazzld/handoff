# Add file output handling logic in main function

## Goal
Add conditional logic in main() to write content to a file when the -output flag is provided, using the handoff.WriteToFile function.

## Implementation Approach
I will modify the existing logic in the main function to:

1. Replace the current if-else block for dry-run/clipboard with a more comprehensive if-else-if-else block that handles:
   - Dry run mode (highest precedence)
   - File output (medium precedence)
   - Clipboard output (lowest precedence/default)

2. Use the existing handoff.WriteToFile function from the library to write the content to the file when the -output flag is provided.

3. Add appropriate success messages when writing to the file is successful.

This implementation will focus on the core file writing functionality while maintaining the existing clipboard and dry-run behavior. Error handling will be implemented in a later task.

## Key Reasoning
I've chosen this approach because:

1. **Modular structure**: Using the existing handoff.WriteToFile function keeps the code modular and follows the principle of separation of concerns.

2. **Maintainability**: The if-else-if-else structure clearly shows the precedence between different output methods, making the code easier to understand and maintain.

3. **Extension**: This structure will be easy to extend if additional output methods are added in the future.

4. **Minimal changes**: By focusing solely on the file output handling without changing other functionality, we reduce the risk of introducing bugs.

5. **Consistency**: This approach aligns with the structure outlined in the original PLAN.md, particularly Task 3 "Implement File Writing Logic" and Task 4 "Adjust Existing Logic".

6. **Reuse**: We leverage the existing WriteToFile function in the library rather than reimplementing file writing functionality in the main package.