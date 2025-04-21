# T020 Implementation Plan: Update examples/gemini_planner.go to Use Library API

## Overview

This task involves updating the `examples/gemini_planner.go` example to use the new library API instead of any previous methods. This ensures the example properly demonstrates the correct usage of the refactored library.

## Current State Assessment

First, we'll examine the current implementation to understand how it works and what needs to be updated.

## Implementation Steps

1. **Examine Current Implementation**
   - Review the existing `examples/gemini_planner.go` file
   - Identify how it currently processes files and any direct usage of files.go or output.go

2. **Update Import Statements**
   - Replace any imports from the main package with the lib package
   - Ensure all necessary library functions are imported

3. **Update Code to Use Library API**
   - Replace direct file processing calls with lib.ProcessProject
   - Update configuration handling to use lib.Config
   - Ensure any file writing uses lib.WriteToFile
   - Maintain functionality while adapting to the new API

4. **Test the Updated Example**
   - Verify the example builds and runs without errors
   - Test with a small project to ensure it still functions as expected

## Technical Implementation Details

### Import Updates
```go
// BEFORE:
import (
    "github.com/phrazzld/handoff/files"  // old import
    "github.com/phrazzld/handoff/output" // old import
)

// AFTER:
import (
    "github.com/phrazzld/handoff/lib"  // new import
)
```

### API Usage Updates

```go
// BEFORE (example, actual code may differ):
config := &files.Config{
    Verbose: verbose,
    Include: include,
    // ...
}
files.ProcessConfig(config)
content := files.ProcessProject([]string{projectDir}, config)

// AFTER:
config := lib.NewConfig()
config.Verbose = verbose
config.Include = include
// ...
config.ProcessConfig()
content, err := lib.ProcessProject([]string{projectDir}, config)
```

## Verification

- Ensure the example compiles successfully
- Run the example to verify it produces the expected output
- Confirm it demonstrates proper usage of the library package
- Check the generated plan for correctness

## Success Criteria

1. The example compiles and runs without errors
2. All references to the old implementation are replaced with the new lib package
3. The example properly demonstrates the intended usage of the library
4. Functionality remains unchanged