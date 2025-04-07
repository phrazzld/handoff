# Implement path resolution and validation

## Goal
Add code to ensure the provided output path is correctly resolved to an absolute path using filepath.Abs when the -output flag is used.

## Implementation Approach
I will implement this task by adding a new function `resolveOutputPath` that will:

1. Accept a path string as input
2. Use filepath.Abs to convert the path to an absolute path
3. Return the absolute path and any errors encountered

This function will be called in the main function when the -output flag is used, before any file operations are attempted. The implementation will include:

1. Verifying the output path is not empty
2. Converting a relative path to an absolute path
3. Handling and reporting any errors that occur during path resolution

This approach keeps the path resolution logic separate from other concerns, making the code more maintainable and easier to test.

## Key Reasoning
I've chosen this approach because:

1. It maintains separation of concerns by isolating the path resolution logic in its own function.

2. The filepath.Abs function from the standard library is the idiomatic way to resolve paths in Go, handling edge cases like relative paths, symbolic links, and platform-specific differences.

3. Creating a separate function is more maintainable and makes it easier to add unit tests for this specific functionality in the future.

4. Resolving the path early in the execution flow allows us to fail fast and provide clear error messages to the user before attempting any file operations.

5. This approach aligns with the "Consideration: File Path Resolution" mentioned in the PLAN.md, which recommends using filepath.Abs for path resolution.