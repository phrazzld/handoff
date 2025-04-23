# T013: Optimize file counting by pre-calculating candidates

## Objective
Optimize the file discovery and processing in `processPaths` to avoid scanning directories multiple times. The current implementation counts total files first and then processes them separately, which may lead to redundant directory scans.

## Analysis
The current implementation in `processPaths` does two operations that involve scanning directories:
1. First, it counts files for the stats.FilesTotal value
2. Then it processes each file with processPathWithProcessor

This approach is inefficient because it may scan the same directories twice. Additionally, when we do the first scan to count total files, we're not storing the discovered files, so we have to rediscover them during processing.

## Approach
1. Modify `processPaths` to perform file discovery once upfront for all provided paths
2. Store the discovered files in a list
3. Calculate `stats.FilesTotal` based on this pre-discovered list
4. Refactor the processing logic to iterate over the pre-discovered files instead of rediscovering them

## Implementation Plan
1. Refactor `processPaths` to:
   - Discover all files upfront and store them in a slice
   - Calculate `stats.FilesTotal` from the length of this slice
   - Process each file from the pre-discovered list

2. Ensure proper filtering is still applied during both file discovery and processing

3. Update the progress tracking logic to use the pre-discovered file count

## Testing
1. Run unit tests to ensure the functionality still works as expected
2. Verify that the correct file count is still displayed in the output stats
3. Run the full test suite to ensure no regressions

## Risks
- If file discovery logic is complex or different between counting and processing, we need to ensure the refactoring maintains the same behavior
- We need to handle potential edge cases like when a file is discovered but removed before processing
