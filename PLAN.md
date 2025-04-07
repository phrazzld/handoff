```markdown
# PLAN.md: Add File Output Option

## 1. Overview

This plan outlines the steps required to add functionality to the `handoff` CLI tool, allowing users to write the aggregated and formatted content directly to a specified file (e.g., `HANDOFF.md`) as an alternative to the current default behavior of copying to the system clipboard. This involves adding a new command-line flag and modifying the main execution flow to handle file writing.

## 2. Task Breakdown

| Task                                                    | Description                                                                                                                               | Effort | Files Affected                 |
| :------------------------------------------------------ | :---------------------------------------------------------------------------------------------------------------------------------------- | :----- | :----------------------------- |
| **1. Add `-output` Flag**                               | Define a new string flag `-output` in `main.go` to accept the desired output filename.                                                    | S      | `main.go`                      |
| **2. Update Flag Parsing**                              | Modify `parseConfig` in `main.go` to handle the new `-output` flag and return its value.                                                  | S      | `main.go`                      |
| **3. Implement File Writing Logic**                     | In `main()`, check if the `-output` flag was provided. If so, call `handoff.WriteToFile` with the processed content and the specified path. | M      | `main.go`                      |
| **4. Adjust Existing Logic (Clipboard/Dry-Run)**        | Ensure that if `-output` is specified, the content is *not* copied to the clipboard. Decide precedence if `-dry-run` is also used.         | S      | `main.go`                      |
| **5. Error Handling**                                   | Add appropriate error handling for file writing operations (e.g., path invalid, permissions denied) and report errors via the logger.     | S      | `main.go`                      |
| **6. Update Documentation**                             | Update `README.md` to document the new `-output` flag, its usage, and examples.                                                           | S      | `README.md`                    |
| **7. Add/Update Tests**                                 | Add tests to verify the new functionality, including file creation, content correctness, and interaction with other flags (`-dry-run`). | M      | `handoff_test.go` (potentially) |
| **8. Final Logging Adjustment**                         | Ensure statistics are still logged to stderr appropriately when using the `-output` flag.                                                 | XS     | `main.go`                      |

**Effort Estimation:** S = Small (<= 2 hours), M = Medium (2-4 hours), L = Large (4+ hours)

## 3. Implementation Details

### 3.1. Add `-output` Flag (`main.go`)

Modify the `parseConfig` function to include the new flag:

```go
// parseConfig defines and parses command-line flags...
func parseConfig() (*handoff.Config, string, bool) { // Return output file path too
	config := handoff.NewConfig()

	// Define flags
	var outputFile string // New variable for output file path
	var dryRun bool
	flag.BoolVar(&config.Verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&dryRun, "dry-run", false, "Preview output instead of copying or writing to file")
	flag.StringVar(&config.Include, "include", "", "Comma-separated list of file extensions to include")
	flag.StringVar(&config.Exclude, "exclude", "", "Comma-separated list of file extensions to exclude")
	flag.StringVar(&config.ExcludeNamesStr, "exclude-names", "", "Comma-separated list of file names to exclude")
	flag.StringVar(&config.Format, "format", "<{path}>\n`` `\n{content}\n`` `\n</{path}>\n\n", "Custom format for output")
	flag.StringVar(&outputFile, "output", "", "Write output to the specified file instead of clipboard (e.g., HANDOFF.md)") // New flag

	// Parse command-line flags
	flag.Parse()

	// Process config
	config.ProcessConfig()

	return config, outputFile, dryRun // Return outputFile
}
```

### 3.2. Implement File Writing Logic (`main.go`)

Modify the `main` function to handle the output destination based on the flags:

```go
func main() {
	// Parse command-line flags and get configuration
	config, outputFile, dryRun := parseConfig() // Get outputFile from parseConfig
	logger := handoff.NewLogger(config.Verbose)

	// Check for paths
	if flag.NArg() < 1 {
		logger.Error("usage: %s [options] path1 [path2 ...]", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Process paths and get content
	formattedContent, err := handoff.ProcessProject(flag.Args(), config)
	if err != nil {
		logger.Error("Failed to process project: %v", err)
		os.Exit(1)
	}

	// Determine output action: Dry Run > Output File > Clipboard
	if dryRun {
		fmt.Println("### DRY RUN: Content that would be generated ###")
		fmt.Println(formattedContent)
		logger.Info("Dry run complete. No file written or clipboard modified.")
	} else if outputFile != "" {
		// Write to file using the library function
		// Consider resolving the path relative to the current working directory
		absPath, err := filepath.Abs(outputFile)
		if err != nil {
			logger.Error("Failed to determine absolute path for output file %s: %v", outputFile, err)
			os.Exit(1)
		}
		logger.Verbose("Writing output to file: %s", absPath)
		if err := handoff.WriteToFile(formattedContent, absPath); err != nil {
			logger.Error("Failed to write to file %s: %v", absPath, err)
			os.Exit(1)
		}
		logger.Info("Output successfully written to %s", absPath)
	} else {
		// Copy to clipboard (existing behavior)
		if err := copyToClipboard(formattedContent); err != nil {
			logger.Error("Failed to copy to clipboard: %v", err)
			// Optionally, print to stdout as a fallback? Or just exit.
			// fmt.Println("Failed to copy to clipboard. Content:\n", formattedContent)
			os.Exit(1)
		}
		logger.Info("Content successfully copied to clipboard.")
	}

	// Calculate and log statistics (This part remains largely the same)
	charCount, lineCount, tokenCount := handoff.CalculateStatistics(formattedContent)
	processedFiles := strings.Count(formattedContent, "</") // Simple count based on default format end tag

	logger.Info("Handoff complete:")
	logger.Info("- Files: %d", processedFiles)
	logger.Info("- Lines: %d", lineCount)
	logger.Info("- Characters: %d", charCount)
	logger.Info("- Estimated tokens: %d", tokenCount)
}
```

### 3.3. Library Function (`lib/handoff.go`)

The required `WriteToFile` function already exists in `lib/handoff.go`. No changes are needed there.

```go
// WriteToFile writes the content to a file
func WriteToFile(content, filePath string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}
```

## 4. Potential Challenges & Considerations

*   **Flag Precedence:** The implementation should clearly define the order of operations if multiple output-related flags are used. The proposed order is: `-dry-run` (prints to stdout, does nothing else) > `-output <file>` (writes to file) > default (copies to clipboard).
*   **File Path Resolution:** The `-output` path should likely be interpreted relative to the current working directory. Using `filepath.Abs` is recommended before writing.
*   **File Overwriting:** The current `handoff.WriteToFile` uses `os.WriteFile`, which will truncate and overwrite existing files. This is standard behavior but should be noted.
*   **Error Handling:** File system errors (permissions, invalid path, disk full) during the `WriteToFile` call need to be caught and reported clearly to the user.
*   **Large Files:** For very large projects, the entire content is held in memory (`formattedContent`). While this is the current architecture, adding file output doesn't change this potential limitation.
*   **Statistics Logging:** Ensure the final statistics summary is always printed to stderr, regardless of whether the output went to clipboard, a file, or stdout (dry-run). The current logging seems appropriate.

## 5. Testing Strategy

*   **Unit Tests:**
    *   Test `parseConfig` to ensure the `-output` flag is correctly parsed and returned.
*   **Integration/CLI Tests:**
    *   Run `handoff -output test_output.md .` and verify `test_output.md` is created with the expected content.
    *   Run `handoff -output existing_file.md .` and verify the file is overwritten.
    *   Run `handoff -output /nonexistent_dir/test.md .` and verify a file writing error is reported.
    *   Run `handoff -output test_output.md -dry-run .` and verify the output goes to stdout and *no* file is created/modified.
    *   Run `handoff -dry-run .` and verify output goes to stdout.
    *   Run `handoff .` (no flags) and verify content is copied to clipboard (requires manual check or platform-specific test setup).
    *   Test with various combinations of include/exclude flags alongside `-output`.
*   **Manual Testing:**
    *   Perform the CLI tests above on different operating systems (Linux, macOS, Windows) if possible.
    *   Verify file permissions and content integrity after writing.

## 6. Open Questions

*   None currently identified. The plan seems straightforward, leveraging existing library functionality. The main work is integrating the flag and conditional logic into `main.go`.
```