# Implement file existence check

## Goal
Add code to check if the output file already exists before writing to it, and refuse to overwrite unless the -force flag is specified.

## Implementation Approach
I will implement a new function called `checkFileExists` that will:

1. Check if a file exists at the specified path
2. Return a boolean indicating whether the file exists
3. Provide clear error handling for access issues

Then, I will modify the main function to:

1. Call this function when a file output path is specified
2. If the file exists and the -force flag is not set, display a warning and exit
3. Allow the program to continue if the file doesn't exist or if the -force flag is set

This approach separates the file existence check from the main program flow, making the code more maintainable and easier to test.

## Key Reasoning
I've chosen this approach because:

1. **Separation of concerns**: Creating a separate function for the file existence check makes the code more modular and easier to maintain. It also allows for better testing of this specific functionality.

2. **User protection**: By default, users should be protected from accidentally overwriting existing files. This is a common safety mechanism in command-line tools.

3. **Explicit opt-in**: Requiring the -force flag to overwrite files means users must explicitly opt into potentially destructive operations.

4. **Clear error messaging**: By checking for file existence early in the process, we can provide clear error messages before attempting any write operations.

5. **Performance considerations**: The file existence check is lightweight and is only performed when necessary (when -output flag is used).

6. **Error handling**: We need to distinguish between a file actually existing versus permission errors that might prevent us from determining if a file exists.

7. **Compatibility with existing code**: This implementation works seamlessly with the already implemented path resolution function.