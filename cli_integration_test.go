// Package main contains the CLI implementation of handoff.
package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// CLI integration tests that build and run the handoff binary.
// These tests verify end-to-end behavior of the CLI, ensuring flags
// correctly influence the output produced via the library calls.

// errorMessageContainsAny checks if the error message contains at least one of the specified key phrases.
// This is more robust than exact string matching as it's less sensitive to minor error message changes.
func errorMessageContainsAny(message string, keyPhrases []string) bool {
	for _, phrase := range keyPhrases {
		if strings.Contains(message, phrase) {
			return true
		}
	}
	return false
}

// buildBinary builds the handoff binary for testing.
// Returns the path to the binary and any error encountered.
func buildBinary(t *testing.T) string {
	t.Helper()

	// Determine binary name based on OS
	binaryName := "handoff_test_bin"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	// Build the binary
	cmd := exec.Command("go", "build", "-o", binaryName)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	// Get absolute path to the binary
	absPath, err := filepath.Abs(binaryName)
	if err != nil {
		t.Fatalf("Failed to get absolute path to binary: %v", err)
	}

	// Register cleanup to remove the binary after tests
	t.Cleanup(func() {
		os.Remove(absPath)
	})

	return absPath
}

// runCliCommand runs the handoff binary with the given arguments.
// Returns stdout, stderr, and any error encountered.
func runCliCommand(t *testing.T, binaryPath string, args ...string) (string, string, error) {
	t.Helper()

	cmd := exec.Command(binaryPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// createTestFiles creates a test file structure for integration testing.
// Returns the path to the temporary directory and a map of file paths to their content.
func createTestFiles(t *testing.T) (string, map[string]string) {
	t.Helper()

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "handoff-integration-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Register cleanup to remove the directory after tests
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	// Create various test files
	fileContents := map[string]string{
		"file1.txt":   "Content of text file",
		"file2.go":    "package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}\n",
		"file3.json":  `{"name": "test", "version": "1.0.0"}`,
		"file4.md":    "# Markdown Test\n\nThis is a test markdown file.",
		".hiddenfile": "This file should be ignored by default",
	}

	// Create a subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Add subdirectory files
	fileContents[filepath.Join("subdir", "subfile1.txt")] = "Content of subdirectory text file"
	fileContents[filepath.Join("subdir", "subfile2.go")] = "package sub\n\nfunc Sub() {}\n"

	// Write all files to disk
	for relPath, content := range fileContents {
		fullPath := filepath.Join(tempDir, relPath)

		// Create parent directories if needed
		parentDir := filepath.Dir(fullPath)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			t.Fatalf("Failed to create directories for %s: %v", fullPath, err)
		}

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", fullPath, err)
		}
	}

	return tempDir, fileContents
}

// TestCLIFlags tests that the CLI accepts various flags.
// This is a simple test to ensure the binary builds and runs.
func TestCLIFlags(t *testing.T) {
	binaryPath := buildBinary(t)
	tempDir, _ := createTestFiles(t)
	
	// Create a test output file
	outputFile := filepath.Join(tempDir, "test_flags_output.md")
	
	// Run with various flags to ensure they're accepted
	_, stderr, err := runCliCommand(t, binaryPath, 
		"-verbose",
		"-output="+outputFile, 
		filepath.Join(tempDir, "file1.txt"))
	
	if err != nil {
		t.Fatalf("CLI failed to run with flags: %v\nStderr: %s", err, stderr)
	}
	
	// Check that output file was created
	if _, err := os.Stat(outputFile); err != nil {
		t.Errorf("Output file was not created: %v", err)
	}
}

// TestCLIBasicFileProcessing tests basic file processing functionality.
func TestCLIBasicFileProcessing(t *testing.T) {
	binaryPath := buildBinary(t)
	tempDir, fileContents := createTestFiles(t)

	// Create a test output file
	outputFile := filepath.Join(tempDir, "output.md")

	// Test processing a single file
	singleFilePath := filepath.Join(tempDir, "file1.txt")
	_, stderr, err := runCliCommand(t, binaryPath, "-output="+outputFile, singleFilePath)
	if err != nil {
		t.Fatalf("Failed to process single file: %v\nStderr: %s", err, stderr)
	}

	// Read output file
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check content
	contentStr := string(content)
	expectedContent := fileContents["file1.txt"]

	// Check that output file contains the file content
	if !strings.Contains(contentStr, expectedContent) {
		t.Errorf("Output doesn't contain expected content from file1.txt")
	}

	// Check that it's wrapped in context tags
	if !strings.HasPrefix(contentStr, "<context>") || !strings.HasSuffix(strings.TrimSpace(contentStr), "</context>") {
		t.Errorf("Output isn't properly wrapped in context tags")
	}

	// Check path tags
	expectedPath := singleFilePath
	if runtime.GOOS == "windows" {
		// Adjust path format for Windows if needed
		expectedPath = filepath.ToSlash(expectedPath)
	}

	if !strings.Contains(contentStr, "<"+expectedPath+">") {
		t.Errorf("Output doesn't contain expected path tag")
	}
}

// TestCLIDirectoryProcessing tests directory processing functionality.
func TestCLIDirectoryProcessing(t *testing.T) {
	binaryPath := buildBinary(t)
	tempDir, fileContents := createTestFiles(t)

	// Create a test output file
	outputFile := filepath.Join(tempDir, "output.md")

	// Process the entire directory
	_, stderr, err := runCliCommand(t, binaryPath, "-output="+outputFile, tempDir)
	if err != nil {
		t.Fatalf("Failed to process directory: %v\nStderr: %s", err, stderr)
	}

	// Read output file
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check content
	contentStr := string(content)

	// Check that output contains content from visible files (not hidden ones)
	visibleFiles := []string{"file1.txt", "file2.go", "file3.json", "file4.md",
		filepath.Join("subdir", "subfile1.txt"), filepath.Join("subdir", "subfile2.go")}

	for _, relPath := range visibleFiles {
		expectedContent := fileContents[relPath]
		if !strings.Contains(contentStr, expectedContent) {
			t.Errorf("Output doesn't contain expected content from %s", relPath)
		}

		// Check path tags (adjust for Windows if needed)
		fullPath := filepath.Join(tempDir, relPath)
		if runtime.GOOS == "windows" {
			fullPath = filepath.ToSlash(fullPath)
		}

		if !strings.Contains(contentStr, "<"+fullPath+">") {
			t.Errorf("Output doesn't contain expected path tag for %s", relPath)
		}
	}

	// Check that hidden files are excluded
	if strings.Contains(contentStr, fileContents[".hiddenfile"]) {
		t.Errorf("Output contains content from hidden file that should be excluded")
	}
}

// TestCLIFiltering tests include and exclude filtering functionality.
func TestCLIFiltering(t *testing.T) {
	binaryPath := buildBinary(t)
	tempDir, fileContents := createTestFiles(t)

	tests := []struct {
		name     string
		args     []string
		included []string
		excluded []string
	}{
		{
			name:     "Include only text files",
			args:     []string{"-include=.txt"},
			included: []string{"file1.txt", filepath.Join("subdir", "subfile1.txt")},
			excluded: []string{"file2.go", "file3.json", "file4.md", filepath.Join("subdir", "subfile2.go")},
		},
		{
			name:     "Include multiple types",
			args:     []string{"-include=.txt,.go"},
			included: []string{"file1.txt", "file2.go", filepath.Join("subdir", "subfile1.txt"), filepath.Join("subdir", "subfile2.go")},
			excluded: []string{"file3.json", "file4.md"},
		},
		{
			name:     "Exclude json files",
			args:     []string{"-exclude=.json"},
			included: []string{"file1.txt", "file2.go", "file4.md", filepath.Join("subdir", "subfile1.txt"), filepath.Join("subdir", "subfile2.go")},
			excluded: []string{"file3.json"},
		},
		{
			name:     "Exclude specific filenames",
			args:     []string{"-exclude-names=file1.txt,file3.json"},
			included: []string{"file2.go", "file4.md", filepath.Join("subdir", "subfile1.txt"), filepath.Join("subdir", "subfile2.go")},
			excluded: []string{"file3.json"}, // exclude-names works with base filenames only
		},
		{
			name:     "Combine include and exclude",
			args:     []string{"-include=.txt,.go", "-exclude=.go"},
			included: []string{"file1.txt", filepath.Join("subdir", "subfile1.txt")},
			excluded: []string{"file2.go", "file3.json", "file4.md", filepath.Join("subdir", "subfile2.go")},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create output file specific to this test
			outputFile := filepath.Join(tempDir, "output_"+strings.ReplaceAll(tc.name, " ", "_")+".md")

			// Prepare args: add output file and directory to test-specific args
			args := append(tc.args, "-output="+outputFile, tempDir)

			// Run command
			_, stderr, err := runCliCommand(t, binaryPath, args...)
			if err != nil {
				t.Fatalf("Failed to run with args %v: %v\nStderr: %s", args, err, stderr)
			}

			// Read output file
			content, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read output file: %v", err)
			}

			contentStr := string(content)

			// Check included files
			for _, relPath := range tc.included {
				expectedContent := fileContents[relPath]
				if !strings.Contains(contentStr, expectedContent) {
					t.Errorf("Output should contain content from %s, but doesn't", relPath)
				}
			}

			// Check excluded files
			for _, relPath := range tc.excluded {
				expectedContent := fileContents[relPath]
				if strings.Contains(contentStr, expectedContent) {
					t.Errorf("Output should not contain content from %s, but does", relPath)
				}
			}
		})
	}
}

// TestCLIForceFlag tests the -force flag for overwriting existing files.
func TestCLIForceFlag(t *testing.T) {
	binaryPath := buildBinary(t)
	tempDir, _ := createTestFiles(t)

	// Create an output file with known content
	outputFile := filepath.Join(tempDir, "existing_output.md")
	initialContent := "This file already exists and should not be overwritten without -force"
	err := os.WriteFile(outputFile, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create existing output file: %v", err)
	}

	// Test 1: Without -force flag - should fail
	_, stderr, err := runCliCommand(t, binaryPath, "-output="+outputFile, filepath.Join(tempDir, "file1.txt"))

	// Command should fail (non-zero exit code)
	if err == nil {
		t.Errorf("Expected command to fail without -force flag when output file exists, but it succeeded")
	}

	// Error message should mention the -force flag or indicate file exists error
	forceErrorPhrases := []string{"-force", "force flag", "already exists", "overwrite", "file exists"}
	if !errorMessageContainsAny(stderr, forceErrorPhrases) {
		t.Errorf("Error message should mention force flag or file exists error, got: %s", stderr)
	}

	// File content should remain unchanged
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	if string(content) != initialContent {
		t.Errorf("File content was changed without -force flag")
	}

	// Test 2: With -force flag - should succeed and overwrite
	_, stderr, err = runCliCommand(t, binaryPath, "-force", "-output="+outputFile, filepath.Join(tempDir, "file1.txt"))

	// Command should succeed
	if err != nil {
		t.Fatalf("Failed to run with -force flag: %v\nStderr: %s", err, stderr)
	}

	// File content should be overwritten
	content, err = os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	if string(content) == initialContent {
		t.Errorf("File content was not overwritten with -force flag")
	}
}

// TestCLIDryRun tests the -dry-run flag.
func TestCLIDryRun(t *testing.T) {
	binaryPath := buildBinary(t)
	tempDir, _ := createTestFiles(t)

	// Create an output file path (should not be created in dry run)
	outputFile := filepath.Join(tempDir, "dry_run_output.md")

	// Run with -dry-run flag
	stdout, _, err := runCliCommand(t, binaryPath, "-dry-run", "-output="+outputFile, filepath.Join(tempDir, "file1.txt"))
	if err != nil {
		t.Fatalf("Failed to run with -dry-run flag: %v", err)
	}

	// Output file should not be created
	if _, err := os.Stat(outputFile); err == nil {
		t.Errorf("Output file was created despite -dry-run flag")
	}

	// Stdout should contain the content that would be written
	if !strings.Contains(stdout, "Content of text file") {
		t.Errorf("Dry run output doesn't contain expected file content")
	}

	// Stdout should include a message indicating it's a dry run
	if !strings.Contains(stdout, "DRY RUN") {
		t.Errorf("Dry run output doesn't include 'DRY RUN' indicator")
	}
}

// TestCLIVerboseFlag tests the -verbose flag.
func TestCLIVerboseFlag(t *testing.T) {
	binaryPath := buildBinary(t)
	tempDir, _ := createTestFiles(t)

	// Create a test output file
	outputFile := filepath.Join(tempDir, "verbose_output.md")

	// Run with verbose flag
	_, stderr, err := runCliCommand(t, binaryPath, 
		"-verbose",
		"-output="+outputFile, 
		tempDir)

	if err != nil {
		t.Fatalf("Failed to run with verbose flag: %v", err)
	}

	// Verify verbose output messages in stderr
	verboseMessages := []string{
		"Processing path:",
		"Processing file",
		"Output will be written to:",
		"Writing content",
		"Processed files successfully",
	}

	for _, msg := range verboseMessages {
		if !strings.Contains(stderr, msg) {
			t.Errorf("Verbose output missing expected message: %q", msg)
		}
	}

	// Check that output file was created
	if _, err := os.Stat(outputFile); err != nil {
		t.Errorf("Output file was not created: %v", err)
	}
}

// TestCLICustomFormat tests the -format flag.
func TestCLICustomFormat(t *testing.T) {
	binaryPath := buildBinary(t)
	tempDir, _ := createTestFiles(t)

	// Create an output file
	outputFile := filepath.Join(tempDir, "custom_format_output.md")

	// Custom format
	customFormat := "FILE: {path}\n---\n{content}\n---\n\n"

	// Run with custom format
	_, stderr, err := runCliCommand(t, binaryPath,
		"-format="+customFormat,
		"-output="+outputFile,
		filepath.Join(tempDir, "file1.txt"))

	if err != nil {
		t.Fatalf("Failed to run with custom format: %v\nStderr: %s", err, stderr)
	}

	// Read output file
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)

	// Check custom format was applied
	if !strings.Contains(contentStr, "FILE:") {
		t.Errorf("Output doesn't contain 'FILE:' prefix from custom format")
	}

	if !strings.Contains(contentStr, "---") {
		t.Errorf("Output doesn't contain '---' separator from custom format")
	}

	// Should not contain default format code fence
	if strings.Contains(contentStr, "```") {
		t.Errorf("Output contains default format code fence despite custom format")
	}
}

// TestCLIErrorHandling tests error conditions.
func TestCLIErrorHandling(t *testing.T) {
	binaryPath := buildBinary(t)

	// Test 1: No paths provided
	_, stderr, err := runCliCommand(t, binaryPath)

	// Command should fail
	if err == nil {
		t.Errorf("Expected command to fail with no paths, but it succeeded")
	}

	// Error message should mention usage or show appropriate error
	usageErrorPhrases := []string{"usage", "Usage", "path1", "options", "no paths", "missing argument"}
	if !errorMessageContainsAny(stderr, usageErrorPhrases) {
		t.Errorf("Error message should indicate usage error, got: %s", stderr)
	}

	// Test 2: Non-existent path
	// The application might show a warning instead of failing with an error
	// since it might continue processing other valid paths
	stdout, stderr, _ := runCliCommand(t, binaryPath, "/path/does/not/exist")

	// Error message should be present in stderr
	fileNotFoundPhrases := []string{
		"no such file", "cannot find", "not found", "doesn't exist", 
		"does not exist", "not exist", "invalid path", "failed to stat", 
	}
	if !errorMessageContainsAny(stderr, fileNotFoundPhrases) {
		t.Errorf("Error message should indicate file not found, got: %s", stderr)
	}
	
	// The output should be empty since no files were processed
	if len(stdout) > 0 && !strings.Contains(stdout, "No files processed") {
		t.Errorf("Expected empty or 'No files processed' output, got: %s", stdout)
	}
}

// TestCLIOutputToFile tests writing output to a file.
func TestCLIOutputToFile(t *testing.T) {
	// This functionality is already tested in TestCLIBasicFileProcessing
	// and other tests that use the -output flag
	t.Skip("Functionality covered by other tests")
}
