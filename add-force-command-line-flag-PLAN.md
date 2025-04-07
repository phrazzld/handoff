# Add -force command line flag

## Goal
Define a new boolean flag `-force` in the flags section of parseConfig function that will allow overwriting existing files when used with the -output flag.

## Implementation Approach
I will implement this task by:

1. Adding the flag definition to the parseConfig function:
   ```go
   flag.BoolVar(&force, "force", false, "Allow overwriting existing files when using -output flag")
   ```

2. Following the existing code style and patterns for consistency.

3. Ensuring the flag variable is properly initialized and returned from the parseConfig function (which is already set up to return the force flag value).

This is a straightforward implementation that focuses solely on defining the flag. The actual logic for using this flag to control file overwriting behavior will be implemented in later tasks.

## Key Reasoning
I've chosen this approach because:

1. It's the simplest way to implement the task as specified in the TODO.md file.

2. It follows the same pattern as other boolean flags in the codebase, maintaining consistency.

3. The flag description clearly communicates its purpose and its relationship to the -output flag.

4. The default value of 'false' is appropriate for safety since we want to prevent accidental file overwriting by default.

5. It keeps the implementation focused on just adding the flag definition, which is the specific requirement for this task. The actual file overwrite protection logic will be implemented in the "Implement file existence check" task.