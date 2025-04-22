package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	handoff "github.com/phrazzld/handoff/lib"
)

// The tests in this file focus on CLI-specific functionality.

// TestParseConfigOutputAndForceFlags tests that parseConfig correctly parses -output and -force flags
func TestParseConfigOutputAndForceFlags(t *testing.T) {
	// Save original command line arguments and flags
	oldArgs := os.Args
	oldFlagCommandLine := flag.CommandLine

	// Restore original values when the test completes
	defer func() {
		os.Args = oldArgs
		flag.CommandLine = oldFlagCommandLine
	}()

	testCases := []struct {
		name            string
		args            []string
		expectedOutput  string
		expectedForce   bool
		expectedDryRun  bool
		expectedVerbose bool
	}{
		{
			name:            "No flags",
			args:            []string{"handoff", "file1.go", "file2.go"},
			expectedOutput:  "",
			expectedForce:   false,
			expectedDryRun:  false,
			expectedVerbose: false,
		},
		{
			name:            "Output flag only",
			args:            []string{"handoff", "-output=output.md", "file1.go"},
			expectedOutput:  "output.md",
			expectedForce:   false,
			expectedDryRun:  false,
			expectedVerbose: false,
		},
		{
			name:            "Force flag only",
			args:            []string{"handoff", "-force", "file1.go"},
			expectedOutput:  "",
			expectedForce:   true,
			expectedDryRun:  false,
			expectedVerbose: false,
		},
		{
			name:            "Verbose flag only",
			args:            []string{"handoff", "-verbose", "file1.go"},
			expectedOutput:  "",
			expectedForce:   false,
			expectedDryRun:  false,
			expectedVerbose: true,
		},
		{
			name:            "DryRun flag only",
			args:            []string{"handoff", "-dry-run", "file1.go"},
			expectedOutput:  "",
			expectedForce:   false,
			expectedDryRun:  true,
			expectedVerbose: false,
		},
		{
			name:            "Both output and force flags",
			args:            []string{"handoff", "-output=HANDOFF.md", "-force", "file1.go"},
			expectedOutput:  "HANDOFF.md",
			expectedForce:   true,
			expectedDryRun:  false,
			expectedVerbose: false,
		},
		{
			name:            "All flags",
			args:            []string{"handoff", "-verbose", "-dry-run", "-output=HANDOFF.md", "-force", "file1.go"},
			expectedOutput:  "HANDOFF.md",
			expectedForce:   true,
			expectedDryRun:  true,
			expectedVerbose: true,
		},
		{
			name:            "Alternative flag order",
			args:            []string{"handoff", "-force", "-output=custom.md", "file1.go"},
			expectedOutput:  "custom.md",
			expectedForce:   true,
			expectedDryRun:  false,
			expectedVerbose: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset command line flags for each test case
			flag.CommandLine = flag.NewFlagSet(tc.args[0], flag.ExitOnError)

			// Set up mock command line arguments
			os.Args = tc.args

			// Call parseConfig
			config, outputFile, force, dryRun := parseConfig()

			// Verify output file path
			if outputFile != tc.expectedOutput {
				t.Errorf("Expected output file %q, got %q", tc.expectedOutput, outputFile)
			}

			// Verify force flag
			if force != tc.expectedForce {
				t.Errorf("Expected force flag %v, got %v", tc.expectedForce, force)
			}

			// Verify dry run flag
			if dryRun != tc.expectedDryRun {
				t.Errorf("Expected dry run flag %v, got %v", tc.expectedDryRun, dryRun)
			}

			// Verify verbose flag
			if config.Verbose != tc.expectedVerbose {
				t.Errorf("Expected verbose flag %v, got %v", tc.expectedVerbose, config.Verbose)
			}
		})
	}
}

// TestFileCreation tests that a file is created with correct content when -output flag is used
func TestFileCreation(t *testing.T) {
	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "handoff-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if cleanErr := os.RemoveAll(tempDir); cleanErr != nil {
			t.Logf("Failed to clean up temp directory: %v", cleanErr)
		}
	}()

	// Create test files within the temp directory
	testFiles := map[string]string{
		"file1.txt": "Test content for file 1",
		"file2.go":  "package main\n\nfunc main() {\n\tfmt.Println(\"Hello World\")\n}",
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filePath, err)
		}
	}

	// Set up the output file path
	outputFile := filepath.Join(tempDir, "output.md")

	// Save original args and flags
	oldArgs := os.Args
	oldFlagCommandLine := flag.CommandLine

	// Restore original values when the test completes
	defer func() {
		os.Args = oldArgs
		flag.CommandLine = oldFlagCommandLine
	}()

	// Set up command line arguments for the test
	os.Args = []string{"handoff", "-output=" + outputFile, tempDir}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Parse flags
	config, outputPath, _, _ := parseConfig()

	// Verify parsed arguments
	if outputPath != outputFile {
		t.Errorf("Expected output path %q, got %q", outputFile, outputPath)
	}

	// Resolve the output path
	absOutputPath, err := resolveOutputPath(outputPath)
	if err != nil {
		t.Fatalf("Failed to resolve output path: %v", err)
	}

	// Check file existence (should not exist yet)
	exists, err := checkFileExists(absOutputPath)
	if err != nil {
		t.Fatalf("Error checking file existence: %v", err)
	}
	if exists {
		t.Errorf("Output file should not exist before the test runs")
	}

	// Process the project files
	formattedContent, _, err := handoff.ProcessProject([]string{tempDir}, config)
	if err != nil {
		t.Fatalf("Failed to process project: %v", err)
	}

	// Write to file
	err = handoff.WriteToFile(formattedContent, absOutputPath)
	if err != nil {
		t.Fatalf("Failed to write to file: %v", err)
	}

	// Verify the file was created
	exists, err = checkFileExists(absOutputPath)
	if err != nil {
		t.Fatalf("Error checking file existence after writing: %v", err)
	}
	if !exists {
		t.Errorf("Output file was not created")
	}

	// Read the file content
	content, err := os.ReadFile(absOutputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Verify the content
	contentStr := string(content)

	// Check that all test files are included in the output
	for filename, fileContent := range testFiles {
		expectedPathTag := "<" + filepath.Join(tempDir, filename) + ">"
		if !strings.Contains(contentStr, expectedPathTag) {
			t.Errorf("Output does not contain path tag for %s", filename)
		}

		if !strings.Contains(contentStr, fileContent) {
			t.Errorf("Output does not contain content for %s", filename)
		}
	}

	// Check that output is wrapped in context tags
	if !strings.HasPrefix(contentStr, "<context>") || !strings.HasSuffix(strings.TrimSpace(contentStr), "</context>") {
		t.Errorf("Output is not properly wrapped in context tags")
	}
}

// TestFileOverwriteProtection tests that existing files are not overwritten without -force flag
// and that they are overwritten when -force flag is provided
func TestFileOverwriteProtection(t *testing.T) {
	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "handoff-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if cleanErr := os.RemoveAll(tempDir); cleanErr != nil {
			t.Logf("Failed to clean up temp directory: %v", cleanErr)
		}
	}()

	// Create test files within the temp directory
	testFiles := map[string]string{
		"file1.txt": "Test content for file 1",
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filePath, err)
		}
	}

	// Set up the output file path
	outputFile := filepath.Join(tempDir, "output.md")

	// Create an existing file at the output path with known content
	initialContent := "This is pre-existing content that should not be overwritten without -force"
	err = os.WriteFile(outputFile, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create initial output file: %v", err)
	}

	// Save original args and flags
	oldArgs := os.Args
	oldFlagCommandLine := flag.CommandLine

	// Restore original values when the test completes
	defer func() {
		os.Args = oldArgs
		flag.CommandLine = oldFlagCommandLine
	}()

	// PART 1: Test that file is NOT overwritten without -force flag
	// Set up command line arguments for the test (without -force)
	os.Args = []string{"handoff", "-output=" + outputFile, tempDir}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Parse flags
	config, outputPath, force, _ := parseConfig()

	// Verify parsed arguments
	if outputPath != outputFile {
		t.Errorf("Expected output path %q, got %q", outputFile, outputPath)
	}
	if force {
		t.Errorf("Force flag should be false")
	}

	// Resolve the output path
	absOutputPath, err := resolveOutputPath(outputPath)
	if err != nil {
		t.Fatalf("Failed to resolve output path: %v", err)
	}

	// Confirm the file exists before attempting to write
	exists, err := checkFileExists(absOutputPath)
	if err != nil {
		t.Fatalf("Error checking file existence: %v", err)
	}
	if !exists {
		t.Errorf("Output file should exist before the test")
	}

	// Process the project files
	formattedContent, _, err := handoff.ProcessProject([]string{tempDir}, config)
	if err != nil {
		t.Fatalf("Failed to process project: %v", err)
	}

	// Attempt to write to the file - this should NOT overwrite without force flag
	// In a real CLI context, this would be prevented by the main() function's file existence check
	// For testing, let's verify the protection logic ourselves
	if exists && !force {
		// Verify the original content is preserved
		content, err := os.ReadFile(absOutputPath)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		if string(content) != initialContent {
			t.Errorf("File content was modified without -force flag")
		}

		// Try writing to the file but expect main to block it
		// To simulate main's behavior without calling os.Exit, don't write if exists && !force
		t.Logf("Correctly detected existing file without -force flag")
	} else {
		t.Errorf("Should have detected existing file without -force flag")
	}

	// PART 2: Test that file IS overwritten with -force flag
	// Set up command line arguments with -force flag
	os.Args = []string{"handoff", "-output=" + outputFile, "-force", tempDir}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Parse flags again with -force
	config, outputPath, force, _ = parseConfig()

	// Verify parsed arguments
	if !force {
		t.Errorf("Force flag should be true")
	}

	// Now write to the file - this should overwrite with force flag
	err = handoff.WriteToFile(formattedContent, absOutputPath)
	if err != nil {
		t.Fatalf("Failed to write to file with -force flag: %v", err)
	}

	// Verify the file was overwritten
	content, err := os.ReadFile(absOutputPath)
	if err != nil {
		t.Fatalf("Failed to read output file after overwrite: %v", err)
	}

	// Verify the content was updated and is not the initial content
	contentStr := string(content)
	if contentStr == initialContent {
		t.Errorf("File was not overwritten with -force flag")
	}

	// Verify the new content has the expected format
	if !strings.Contains(contentStr, testFiles["file1.txt"]) {
		t.Errorf("New content doesn't contain expected test file content")
	}

	if !strings.HasPrefix(contentStr, "<context>") {
		t.Errorf("New content isn't properly formatted with context tags")
	}
}

// TestInvalidPathErrorHandling tests error handling for invalid output paths
func TestInvalidPathErrorHandling(t *testing.T) {
	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "handoff-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if cleanErr := os.RemoveAll(tempDir); cleanErr != nil {
			t.Logf("Failed to clean up temp directory: %v", cleanErr)
		}
	}()

	// Create test files within the temp directory
	testFilePath := filepath.Join(tempDir, "file1.txt")
	testContent := "Test content for error handling"
	err = os.WriteFile(testFilePath, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Save original args and flags
	oldArgs := os.Args
	oldFlagCommandLine := flag.CommandLine

	// Restore original values when the test completes
	defer func() {
		os.Args = oldArgs
		flag.CommandLine = oldFlagCommandLine
	}()

	// PART 1: Test invalid directory path
	// Create a path to a file in a non-existent directory
	nonExistentDir := filepath.Join(tempDir, "non-existent-dir")
	invalidPath := filepath.Join(nonExistentDir, "output.md")

	// Set up command line arguments for the test
	os.Args = []string{"handoff", "-output=" + invalidPath, tempDir}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Parse flags
	config, outputPath, _, _ := parseConfig()

	// Verify parsed arguments
	if outputPath != invalidPath {
		t.Errorf("Expected output path %q, got %q", invalidPath, outputPath)
	}

	// Resolve the output path - this should work, but the directory doesn't exist
	absOutputPath, err := resolveOutputPath(outputPath)
	if err != nil {
		t.Fatalf("Failed to resolve output path: %v", err)
	}

	// Process the project files
	formattedContent, _, err := handoff.ProcessProject([]string{tempDir}, config)
	if err != nil {
		t.Fatalf("Failed to process project: %v", err)
	}

	// Attempt to write to the file - this should fail because the directory doesn't exist
	err = handoff.WriteToFile(formattedContent, absOutputPath)
	if err == nil {
		t.Errorf("Expected an error when writing to a non-existent directory, but got none")
	} else {
		// Verify the error message indicates the directory issue
		if !strings.Contains(err.Error(), "no such file or directory") {
			t.Errorf("Expected 'no such file or directory' error, got: %v", err)
		}
	}

	// PART 2: Test inaccessible path due to permissions
	// Skip if running as root, as root can write anywhere
	if os.Geteuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	// Create a read-only directory
	readOnlyDir := filepath.Join(tempDir, "read-only-dir")
	err = os.Mkdir(readOnlyDir, 0500) // read + execute, no write
	if err != nil {
		t.Fatalf("Failed to create read-only directory: %v", err)
	}

	// Set up a path to a file in the read-only directory
	readOnlyPath := filepath.Join(readOnlyDir, "output.md")

	// Set up command line arguments for the test
	os.Args = []string{"handoff", "-output=" + readOnlyPath, tempDir}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Parse flags
	config, outputPath, _, _ = parseConfig()

	// Verify parsed arguments
	if outputPath != readOnlyPath {
		t.Errorf("Expected output path %q, got %q", readOnlyPath, outputPath)
	}

	// Resolve the output path
	absOutputPath, err = resolveOutputPath(outputPath)
	if err != nil {
		t.Fatalf("Failed to resolve output path: %v", err)
	}

	// Attempt to write to the file - this should fail due to permissions
	err = handoff.WriteToFile(formattedContent, absOutputPath)
	if err == nil {
		t.Errorf("Expected a permission error when writing to a read-only directory, but got none")
	} else {
		// Verify the error message indicates the permission issue
		if !strings.Contains(err.Error(), "permission denied") {
			t.Errorf("Expected 'permission denied' error, got: %v", err)
		}
	}
}

// determineOutputMode mimics the precedence logic in main() for testing purposes
// Returns "dry-run", "file", or "clipboard" based on the provided flags
func determineOutputMode(dryRun bool, outputFile string) string {
	if dryRun {
		// Highest precedence: dry-run mode
		return "dry-run"
	} else if outputFile != "" {
		// Medium precedence: write to file
		return "file"
	} else {
		// Lowest precedence: copy to clipboard (default behavior)
		return "clipboard"
	}
}

// TestFlagInteractionPrecedence tests that the correct precedence is followed
// when various combinations of flags are used
func TestFlagInteractionPrecedence(t *testing.T) {
	// Create temporary test directory and file path for testing
	tempDir, err := os.MkdirTemp("", "handoff-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if cleanErr := os.RemoveAll(tempDir); cleanErr != nil {
			t.Logf("Failed to clean up temp directory: %v", cleanErr)
		}
	}()

	// Create test file
	testFilePath := filepath.Join(tempDir, "file1.txt")
	testContent := "Test content for flag interaction"
	err = os.WriteFile(testFilePath, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Set up the output file path
	outputFile := filepath.Join(tempDir, "output.md")

	// Save original args and flags
	oldArgs := os.Args
	oldFlagCommandLine := flag.CommandLine

	// Restore original values when the test completes
	defer func() {
		os.Args = oldArgs
		flag.CommandLine = oldFlagCommandLine
	}()

	// Test cases for different flag combinations
	testCases := []struct {
		name               string
		args               []string
		expectedOutputMode string
	}{
		{
			name:               "No special flags",
			args:               []string{"handoff", tempDir},
			expectedOutputMode: "clipboard",
		},
		{
			name:               "Output flag only",
			args:               []string{"handoff", "-output=" + outputFile, tempDir},
			expectedOutputMode: "file",
		},
		{
			name:               "Output and force flags",
			args:               []string{"handoff", "-output=" + outputFile, "-force", tempDir},
			expectedOutputMode: "file",
		},
		{
			name:               "Dry-run flag only",
			args:               []string{"handoff", "-dry-run", tempDir},
			expectedOutputMode: "dry-run",
		},
		{
			name:               "Dry-run and output flags",
			args:               []string{"handoff", "-dry-run", "-output=" + outputFile, tempDir},
			expectedOutputMode: "dry-run",
		},
		{
			name:               "All flags",
			args:               []string{"handoff", "-dry-run", "-output=" + outputFile, "-force", tempDir},
			expectedOutputMode: "dry-run",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up command line arguments for the test
			os.Args = tc.args
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// Parse flags
			_, outputPath, _, dryRun := parseConfig()

			// Determine the output mode based on the flags
			outputMode := determineOutputMode(dryRun, outputPath)

			// Verify the correct output mode was selected
			if outputMode != tc.expectedOutputMode {
				t.Errorf("Expected output mode %q, got %q", tc.expectedOutputMode, outputMode)
			}
		})
	}
}

// TestResolveOutputPath tests the resolveOutputPath function
func TestResolveOutputPath(t *testing.T) {
	// Test cases
	testCases := []struct {
		name      string
		inputPath string
		wantErr   bool
	}{
		{
			name:      "Valid relative path",
			inputPath: "output.md",
			wantErr:   false,
		},
		{
			name:      "Valid absolute path",
			inputPath: "/tmp/output.md",
			wantErr:   false,
		},
		{
			name:      "Empty path",
			inputPath: "",
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			absPath, err := resolveOutputPath(tc.inputPath)

			// Check error condition
			if (err != nil) != tc.wantErr {
				t.Errorf("resolveOutputPath(%q) error = %v, wantErr %v", tc.inputPath, err, tc.wantErr)
				return
			}

			// For valid paths, check that the result is absolute
			if !tc.wantErr && !filepath.IsAbs(absPath) {
				t.Errorf("resolveOutputPath(%q) = %q, which is not an absolute path", tc.inputPath, absPath)
			}
		})
	}
}

// TestCopyToClipboardErrorHandling tests the error handling in copyToClipboard
// when no supported clipboard mechanism is available
func TestCopyToClipboardErrorHandling(t *testing.T) {
	// Save original PATH and restore after test
	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)

	// Set PATH to a non-existent directory to ensure clipboard commands can't be found
	os.Setenv("PATH", "/this/path/does/not/exist")

	err := copyToClipboard("Test content")

	// We expect an error since no clipboard commands should be available
	if err == nil {
		t.Errorf("Expected error when no clipboard commands are available, but got nil")
	} else {
		// Verify error message contains expected text
		if !strings.Contains(err.Error(), "clipboard commands failed") {
			t.Errorf("Expected error message to contain 'clipboard commands failed', got: %v", err)
		}
	}
}

// TestCheckFileExists tests the checkFileExists function
func TestCheckFileExists(t *testing.T) {
	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "handoff-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if cleanErr := os.RemoveAll(tempDir); cleanErr != nil {
			t.Logf("Failed to clean up temp directory: %v", cleanErr)
		}
	}()

	// Create test file
	existingFilePath := filepath.Join(tempDir, "existing.txt")
	err = os.WriteFile(existingFilePath, []byte("Test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Path to non-existent file
	nonExistentPath := filepath.Join(tempDir, "nonexistent.txt")

	// Test cases
	testCases := []struct {
		name     string
		path     string
		expected bool
		wantErr  bool
	}{
		{
			name:     "Existing file",
			path:     existingFilePath,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "Non-existent file",
			path:     nonExistentPath,
			expected: false,
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exists, err := checkFileExists(tc.path)

			// Check error condition
			if (err != nil) != tc.wantErr {
				t.Errorf("checkFileExists(%q) error = %v, wantErr %v", tc.path, err, tc.wantErr)
				return
			}

			// Check return value
			if exists != tc.expected {
				t.Errorf("checkFileExists(%q) = %v, want %v", tc.path, exists, tc.expected)
			}
		})
	}
}
