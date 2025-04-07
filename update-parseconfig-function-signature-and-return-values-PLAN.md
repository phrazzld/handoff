# Update parseConfig function signature and return values

## Goal
Modify the parseConfig function signature in main.go to return the output file path and force flag values in addition to the config and dryRun values.

## Implementation Approach
I will update the parseConfig function signature and implementation to:
1. Add a new string variable to store the output file path
2. Add a new boolean variable to store the force flag value
3. Update the function signature to return these additional values
4. Keep the existing functionality intact

This is the simplest, most direct approach to implement this change. It maintains compatibility with the existing codebase while preparing for the subsequent tasks of actually defining and using the new flags.

## Key Reasoning
This implementation approach was chosen because:
1. It's a minimal change that doesn't modify existing behavior
2. It provides a clear foundation for the upcoming tasks that will define and use the new flags
3. It follows the established pattern of how flags are defined and returned in the existing codebase
4. It maintains the separation of concerns - parseConfig handles flag parsing and returns the parsed values, while the actual file output logic will be implemented separately