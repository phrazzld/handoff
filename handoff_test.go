package main

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"testing"
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

// The following tests have been removed as they are redundant with existing tests:
// - TestFileCreation: Functionality covered by TestCLIBasicFileProcessing in cli_integration_test.go
// - TestFileOverwriteProtection: Functionality covered by TestCLIForceFlag in cli_integration_test.go
// - TestInvalidPathErrorHandling: Functionality covered by TestCLIErrorHandling in cli_integration_test.go 
//   and TestWriteToFileWithDirectoryCreation in lib/handoff_test.go

// Note: determineOutputMode function was removed as it duplicated the output mode logic in main.go

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

			// Directly determine and verify the expected output mode based on the flags
			var actualMode string
			if dryRun {
				actualMode = "dry-run"
			} else if outputPath != "" {
				actualMode = "file"
			} else {
				actualMode = "clipboard"
			}

			// Verify the correct output mode based on precedence rules
			if actualMode != tc.expectedOutputMode {
				t.Errorf("Expected output mode %q, got %q", tc.expectedOutputMode, actualMode)
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
		// Use errors.Is for more robust error type checking
		if !errors.Is(err, ErrClipboardFailed) {
			t.Errorf("Expected error type ErrClipboardFailed, got: %v", err)
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
