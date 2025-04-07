# Add -output command line flag

## Goal
Define a new string flag `-output` in the flags section of parseConfig function that will accept a file path as input.

## Implementation Approach
I will implement this task by:

1. Adding the flag definition in the parseConfig function:
   ```go
   flag.StringVar(&outputFile, "output", "", "Write output to the specified file instead of clipboard (e.g., HANDOFF.md)")
   ```

2. Ensuring the flag follows the existing pattern and style in the codebase for consistency.

3. Verifying that the outputFile variable is properly initialized and returned by the function.

This approach is straightforward and minimal, as the actual handling of the flag's value will be implemented in subsequent tasks. In the current task, we're only concerned with defining the flag itself to be available for command-line parsing.

## Key Reasoning
I've chosen this approach because:

1. It's the most direct way to implement the task as specified in the TODO.md file.

2. It follows the existing pattern of flag definition in the codebase, maintaining consistency with how other flags are defined.

3. It's a minimal change that focuses specifically on adding the flag definition, without introducing any additional complexity or logic at this stage.

4. The description text is clear and follows the format used for other flags, providing users with appropriate information about what the flag does.

5. The empty default value is appropriate, as we want the flag to be explicitly set by the user rather than having a pre-defined default output file.