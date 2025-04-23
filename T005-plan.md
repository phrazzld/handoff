# T005 Implementation Plan: Implement Functional Options for Configuration

## Context
Currently, the handoff library requires users to manually call `ProcessConfig()` after setting configuration options. This is error-prone and easy to forget. We should refactor to use the functional options pattern, which is more idiomatic in Go and provides a better developer experience.

## Analysis
- The current configuration system uses string fields that need manual conversion to internal representation.
- Users must remember to call `ProcessConfig()` after setting options.
- The functional options pattern will simplify configuration and eliminate the need for the ProcessConfig() call.

## Changes Required

### 1. Define Option function type
- Create a type `Option func(*Config)` 
- Modify `NewConfig` to accept variadic Options

### 2. Create option functions
- Implement functions like `WithInclude`, `WithExclude`, etc.
- Each function should handle parsing of string input internally

### 3. Update Config struct
- Make string-based fields unexported
- Keep processed slice fields as they are (potentially rename for clarity)
- Remove ProcessConfig() method

### 4. Update all internal usages
- Modify all places that use the configuration
- Update tests to use the new pattern

## Implementation Steps

1. Modify the Config struct to make string fields unexported
2. Define the Option type and modify NewConfig
3. Implement all option functions
4. Update internal usages of Config
5. Update tests
6. Run tests to ensure everything works correctly

## Testing Plan
- Ensure all existing tests pass with the new implementation
- Add tests for each option function
- Verify that ProcessConfig() is no longer needed or called internally
