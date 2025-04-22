// Package handoff_test contains tests for the handoff library
// These tests verify functionality for collecting and formatting file contents.
package handoff

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsBinaryFile(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{
			name:     "Plain text",
			content:  []byte("This is plain text.\nIt has multiple lines."),
			expected: false,
		},
		{
			name:     "Binary with null bytes",
			content:  []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x00, 0x57, 0x6F, 0x72, 0x6C, 0x64},
			expected: true,
		},
		{
			name:     "Binary with high non-printable ratio",
			content:  []byte{0x7F, 0x7F, 0x7F, 0x7F, 0x41, 0x42, 0x43},
			expected: true,
		},
		{
			name:     "Text with some valid control chars",
			content:  []byte("Hello\nWorld\tWith\rSome\tControl\nChars"),
			expected: false,
		},
		{
			name:     "Empty content",
			content:  []byte{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isBinaryFile(tt.content)
			if result != tt.expected {
				t.Errorf("isBinaryFile() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestShouldProcess(t *testing.T) {
	testCases := []struct {
		name       string
		file       string
		config     *Config
		shouldProc bool
	}{
		{
			name: "Match include extension",
			file: "test.go",
			config: &Config{
				IncludeExts: []string{".go"},
			},
			shouldProc: true,
		},
		{
			name: "No match include extension",
			file: "test.txt",
			config: &Config{
				IncludeExts: []string{".go"},
			},
			shouldProc: false,
		},
		{
			name: "Match exclude extension",
			file: "test.bin",
			config: &Config{
				ExcludeExts: []string{".bin"},
			},
			shouldProc: false,
		},
		{
			name: "Match exclude name",
			file: "package-lock.json",
			config: &Config{
				ExcludeNames: []string{"package-lock.json"},
			},
			shouldProc: false,
		},
		{
			name:       "Default config (no filters)",
			file:       "anything.txt",
			config:     &Config{},
			shouldProc: true,
		},
		{
			name: "Match exclude has precedence over include",
			file: "something.exe",
			config: &Config{
				IncludeExts: []string{".exe"},
				ExcludeExts: []string{".exe"},
			},
			shouldProc: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := shouldProcess(tc.file, tc.config)
			if result != tc.shouldProc {
				t.Errorf("shouldProcess(%q) = %v, want %v", tc.file, result, tc.shouldProc)
			}
		})
	}
}

func TestGetFilesFromDir(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "handoff-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if cleanErr := os.RemoveAll(tmpDir); cleanErr != nil {
			t.Logf("Failed to clean up temp directory: %v", cleanErr)
		}
	}()

	// Create some files in the directory
	files := []string{"file1.txt", "file2.go", ".hidden"}
	for _, file := range files {
		path := filepath.Join(tmpDir, file)
		if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", path, err)
		}
	}

	// Create a subdirectory with files
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}
	subFiles := []string{"subfile1.txt", "subfile2.go"}
	for _, file := range subFiles {
		path := filepath.Join(subDir, file)
		if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", path, err)
		}
	}

	// Test through ProcessProject which will use getFilesFromDir internally
	config := NewConfig()
	// Process the directory
	content, _, err := ProcessProject([]string{tmpDir}, config)
	if err != nil {
		t.Fatalf("ProcessProject failed: %v", err)
	}

	// Verify content includes visible files but not hidden ones
	if !strings.Contains(content, "file1.txt") {
		t.Errorf("Content should include file1.txt")
	}
	if !strings.Contains(content, "file2.go") {
		t.Errorf("Content should include file2.go")
	}
	if !strings.Contains(content, "subfile1.txt") {
		t.Errorf("Content should include subfile1.txt")
	}
	if !strings.Contains(content, "subfile2.go") {
		t.Errorf("Content should include subfile2.go")
	}
	if strings.Contains(content, ".hidden") {
		t.Errorf("Content should not include .hidden")
	}
}

func TestProcessFile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "handoff-process-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if cleanErr := os.RemoveAll(tmpDir); cleanErr != nil {
			t.Logf("Failed to clean up temp directory: %v", cleanErr)
		}
	}()

	// Create a temporary file for testing
	filePath := filepath.Join(tmpDir, "test.txt")
	fileContent := "This is a test file.\nIt has multiple lines."
	if err := os.WriteFile(filePath, []byte(fileContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a custom processor function
	processor := func(file string, content []byte) string {
		return fmt.Sprintf("PROCESSED: %s\n%s", file, string(content))
	}

	// Create a test config
	config := &Config{}

	// Create a logger
	logger := NewLogger(false)

	// Test processing a valid file
	result := ProcessFile(filePath, logger, config, processor)
	expected := "PROCESSED: " + filePath + "\n" + fileContent
	if result != expected {
		t.Errorf("ProcessFile() = %q, want %q", result, expected)
	}

	// Test file that doesn't exist
	nonExistentPath := filepath.Join(tmpDir, "non-existent")
	result = ProcessFile(nonExistentPath, logger, config, processor)
	if result != "" {
		t.Errorf("ProcessFile() for non-existent file returned %q, want empty string", result)
	}

	// Test binary file (simulated using isBinaryFile mock)
	// Create a temporary binary file
	binaryFilePath := filepath.Join(tmpDir, "binary.bin")
	binaryContent := []byte{0x00, 0x01, 0x02, 0x03}
	if err := os.WriteFile(binaryFilePath, binaryContent, 0644); err != nil {
		t.Fatalf("Failed to create binary test file: %v", err)
	}

	// Test processing a binary file
	result = ProcessFile(binaryFilePath, logger, config, processor)
	if result != "" {
		t.Errorf("ProcessFile() for binary file returned %q, want empty string", result)
	}

	// Test with exclusion config
	configWithExclude := &Config{
		ExcludeExts: []string{".txt"},
	}
	result = ProcessFile(filePath, logger, configWithExclude, processor)
	if result != "" {
		t.Errorf("ProcessFile() for excluded extension returned %q, want empty string", result)
	}
}

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	// Check default values
	if config.Verbose {
		t.Errorf("Default config.Verbose = true, want false")
	}
	if config.Format != "<{path}>\n```\n{content}\n```\n</{path}>\n\n" {
		t.Errorf("Default config.Format = %q, want default format", config.Format)
	}
	if len(config.IncludeExts) != 0 {
		t.Errorf("Default config.IncludeExts = %v, want empty slice", config.IncludeExts)
	}
	if len(config.ExcludeExts) != 0 {
		t.Errorf("Default config.ExcludeExts = %v, want empty slice", config.ExcludeExts)
	}
	if len(config.ExcludeNames) != 0 {
		t.Errorf("Default config.ExcludeNames = %v, want empty slice", config.ExcludeNames)
	}
}

func TestProcessConfig(t *testing.T) {
	testCases := []struct {
		name             string
		includeStr       string
		excludeStr       string
		excludeNamesStr  string
		wantIncludeExts  []string
		wantExcludeExts  []string
		wantExcludeNames []string
	}{
		{
			name:             "Empty strings",
			includeStr:       "",
			excludeStr:       "",
			excludeNamesStr:  "",
			wantIncludeExts:  nil,
			wantExcludeExts:  nil,
			wantExcludeNames: nil,
		},
		{
			name:             "Simple extensions",
			includeStr:       ".go,.md",
			excludeStr:       ".exe,.bin",
			excludeNamesStr:  "package-lock.json,yarn.lock",
			wantIncludeExts:  []string{".go", ".md"},
			wantExcludeExts:  []string{".exe", ".bin"},
			wantExcludeNames: []string{"package-lock.json", "yarn.lock"},
		},
		{
			name:             "Extensions without dots",
			includeStr:       "go,md",
			excludeStr:       "exe,bin",
			excludeNamesStr:  "",
			wantIncludeExts:  []string{".go", ".md"},
			wantExcludeExts:  []string{".exe", ".bin"},
			wantExcludeNames: nil,
		},
		{
			name:             "Whitespace handling",
			includeStr:       " go , md ",
			excludeStr:       " exe , bin ",
			excludeNamesStr:  " package-lock.json , yarn.lock ",
			wantIncludeExts:  []string{".go", ".md"},
			wantExcludeExts:  []string{".exe", ".bin"},
			wantExcludeNames: []string{"package-lock.json", "yarn.lock"},
		},
		{
			name:             "Mixed dot presence",
			includeStr:       ".go,md,  .txt ",
			excludeStr:       "exe,.bin",
			excludeNamesStr:  "",
			wantIncludeExts:  []string{".go", ".md", ".txt"},
			wantExcludeExts:  []string{".exe", ".bin"},
			wantExcludeNames: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &Config{
				Include:        tc.includeStr,
				Exclude:        tc.excludeStr,
				ExcludeNamesStr: tc.excludeNamesStr,
			}

			config.ProcessConfig()

			// Check include extensions
			if !equalSlices(config.IncludeExts, tc.wantIncludeExts) {
				t.Errorf("IncludeExts = %v, want %v", config.IncludeExts, tc.wantIncludeExts)
			}

			// Check exclude extensions
			if !equalSlices(config.ExcludeExts, tc.wantExcludeExts) {
				t.Errorf("ExcludeExts = %v, want %v", config.ExcludeExts, tc.wantExcludeExts)
			}

			// Check exclude names
			if !equalSlices(config.ExcludeNames, tc.wantExcludeNames) {
				t.Errorf("ExcludeNames = %v, want %v", config.ExcludeNames, tc.wantExcludeNames)
			}
		})
	}
}

// Helper function to check if two string slices are equal
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Helper to create a test directory with various files
func createTestDir(t *testing.T) (string, []string) {
	tmpDir, err := os.MkdirTemp("", "handoff-process-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create regular files
	textFiles := []string{"file1.txt", "file2.md", "readme.md"}
	codeFiles := []string{"main.go", "util.go"}
	binaryFiles := []string{"binary.bin", "executable.exe"}

	// Add content to regular files
	for _, file := range append(textFiles, codeFiles...) {
		path := filepath.Join(tmpDir, file)
		content := fmt.Sprintf("This is the content of %s\nIt has multiple lines.\n", file)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Add binary content to binary files
	for _, file := range binaryFiles {
		path := filepath.Join(tmpDir, file)
		content := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05} // Some binary-like content
		if err := os.WriteFile(path, content, 0644); err != nil {
			t.Fatalf("Failed to create binary file %s: %v", file, err)
		}
	}

	// Create subdirectory with files
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	subFiles := []string{"subfile1.txt", "subfile2.go"}
	for _, file := range subFiles {
		path := filepath.Join(subDir, file)
		content := fmt.Sprintf("This is a file in the subdirectory: %s\n", file)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file in subdirectory %s: %v", file, err)
		}
	}

	allFiles := append(append(textFiles, codeFiles...), binaryFiles...)
	allFiles = append(allFiles, filepath.Join("subdir", subFiles[0]), filepath.Join("subdir", subFiles[1]))

	return tmpDir, allFiles
}

// TestProcessProject tests the ProcessProject function with various configurations
func TestProcessProject(t *testing.T) {
	// Create test data
	tmpDir, _ := createTestDir(t)
	defer func() {
		if cleanErr := os.RemoveAll(tmpDir); cleanErr != nil {
			t.Logf("Failed to clean up test directory: %v", cleanErr)
		}
	}()

	testCases := []struct {
		name          string
		paths         []string
		config        *Config
		wantErr       bool
		expectContent []string     // Strings that should be in the result
		rejectContent []string     // Strings that should NOT be in the result
		customFormat  bool         // Whether to test custom format
		noProcess     bool         // Skip processing (e.g., for error cases)
	}{
		{
			name:          "Basic test with default config",
			paths:         []string{tmpDir},
			config:        NewConfig(),
			wantErr:       false,
			expectContent: []string{"file1.txt", "main.go", "subfile1.txt"},
			rejectContent: []string{"binary.bin", "executable.exe"}, // Binary files should be excluded
		},
		{
			name:          "Filter by include extension",
			paths:         []string{tmpDir},
			config:        &Config{Include: ".go"},
			wantErr:       false,
			expectContent: []string{"main.go", "util.go", "subfile2.go"},
			rejectContent: []string{"file1.txt", "readme.md", "binary.bin"},
		},
		{
			name:          "Filter by exclude extension",
			paths:         []string{tmpDir},
			config:        &Config{Exclude: ".txt,.md"},
			wantErr:       false,
			expectContent: []string{"main.go", "util.go"},
			rejectContent: []string{"file1.txt", "readme.md"},
		},
		{
			name:          "Process specific paths only",
			paths:         []string{filepath.Join(tmpDir, "main.go"), filepath.Join(tmpDir, "file1.txt")},
			config:        NewConfig(),
			wantErr:       false,
			expectContent: []string{"main.go", "file1.txt"},
			rejectContent: []string{"util.go", "readme.md", "subfile1.txt"},
		},
		{
			name:          "Custom format",
			paths:         []string{filepath.Join(tmpDir, "main.go")},
			config:        &Config{Format: "FILE: {path}\n---\n{content}\n---\n"},
			wantErr:       false,
			expectContent: []string{"FILE:", "---"},
			customFormat:  true,
		},
		{
			name:    "Empty paths",
			paths:   []string{},
			config:  NewConfig(),
			wantErr: true,
			noProcess: true,
		},
		{
			name:    "Invalid path",
			paths:   []string{"/path/that/does/not/exist"},
			config:  NewConfig(),
			wantErr: false, // We shouldn't return an error for non-existent paths now
			expectContent: []string{},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Process config if needed
			if tt.config != nil && !tt.noProcess {
				tt.config.ProcessConfig()
			}

			// Call ProcessProject
			content, _, err := ProcessProject(tt.paths, tt.config)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expected an error, no need to check content
			if tt.wantErr {
				if content != "" {
					t.Errorf("ProcessProject() content should be empty, got: %s", content)
				}
				return
			}

			// Check expected content
			for _, expected := range tt.expectContent {
				if !strings.Contains(content, expected) {
					t.Errorf("ProcessProject() content should contain %q, but doesn't", expected)
				}
			}

			// Check rejected content
			for _, unexpected := range tt.rejectContent {
				if strings.Contains(content, unexpected) {
					t.Errorf("ProcessProject() content should not contain %q, but does", unexpected)
				}
			}

			// Check custom format if needed
			if tt.customFormat {
				if !strings.Contains(content, "FILE:") {
					t.Errorf("ProcessProject() with custom format should contain 'FILE:' prefix")
				}
				if !strings.Contains(content, "---") {
					t.Errorf("ProcessProject() with custom format should contain '---' separator")
				}
			}
		})
	}
}

// TestProcessProjectWithVerbose tests the verbose output mode of ProcessProject
func TestProcessProjectWithVerbose(t *testing.T) {
	tmpDir, _ := createTestDir(t)
	defer func() {
		if cleanErr := os.RemoveAll(tmpDir); cleanErr != nil {
			t.Logf("Failed to clean up test directory: %v", cleanErr)
		}
	}()

	// Create a Config with Verbose enabled
	config := NewConfig()
	config.Verbose = true

	// Capture stderr output
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Run ProcessProject
	content, stats, err := ProcessProject([]string{tmpDir}, config)
	if err != nil {
		t.Fatalf("ProcessProject failed: %v", err)
	}

	// Close the pipe writer and restore stderr
	if err := w.Close(); err != nil {
		t.Errorf("Failed to close writer: %v", err)
	}
	os.Stderr = oldStderr

	// Read the captured output to verify verbose logging messages
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Errorf("Failed to read stderr output: %v", err)
	}
	
	// Check for expected verbose output messages
	verboseOutput := buf.String()
	expectedMessages := []string{
		"Processing path:",
		"Processing file",
	}
	
	for _, msg := range expectedMessages {
		if !strings.Contains(verboseOutput, msg) {
			t.Errorf("Verbose output missing expected message: %q", msg)
		}
	}

	// Verify we have valid content
	if len(content) == 0 {
		t.Error("ProcessProject() returned empty content")
	}

	// Verify that stats contain valid data
	if stats.FilesProcessed == 0 {
		t.Error("ProcessProject() returned stats with zero files processed")
	}
	if stats.Chars == 0 {
		t.Error("ProcessProject() returned stats with zero characters")
	}
	
	// ProcessProject no longer logs statistics directly as that's the caller's responsibility
	// This test now verifies that the Stats struct is properly populated instead
}

// TestProcessPaths_ErrNoFilesProcessed tests the error handling in ProcessPaths when no files are processed
func TestProcessPaths_ErrNoFilesProcessed(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "handoff-nofiles-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a file that we can exclude with config
	filePath := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create config that will exclude the file
	config := NewConfig()
	config.Exclude = ".txt" // Exclude the .txt file we created
	config.ProcessConfig()
	logger := NewLogger(false)

	// Call ProcessPaths
	_, stats, err := ProcessPaths([]string{tmpDir}, config, logger)

	// Assert that we get the expected error
	if !errors.Is(err, ErrNoFilesProcessed) {
		t.Errorf("Expected ErrNoFilesProcessed, got %v", err)
	}

	// Assert that FilesProcessed is 0
	if stats.FilesProcessed != 0 {
		t.Errorf("Expected FilesProcessed to be 0, got %d", stats.FilesProcessed)
	}

	// Assert that FilesTotal is greater than 0
	if stats.FilesTotal <= 0 {
		t.Errorf("Expected FilesTotal to be > 0, got %d", stats.FilesTotal)
	}
}

// TestProcessProject_NoFilesProcessed tests that ProcessProject returns ErrNoFilesProcessed when appropriate
func TestProcessProject_NoFilesProcessed(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "handoff-project-nofiles-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a file that we can exclude with configuration
	filePath := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Manually set GitAvailable to false to ensure we use filepath.Walk
	// which will count the files properly for this test
	originalGitAvailable := GitAvailable
	GitAvailable = false
	defer func() {
		// Restore the original value
		GitAvailable = originalGitAvailable
	}()

	// Create config that will exclude the file
	config := NewConfig()
	config.Exclude = ".txt" // Exclude the .txt file we created
	config.ProcessConfig()

	// Call ProcessProject
	content, stats, err := ProcessProject([]string{tmpDir}, config)

	// Assert that we get the expected error
	if !errors.Is(err, ErrNoFilesProcessed) {
		t.Errorf("Expected ErrNoFilesProcessed, got %v", err)
	}

	// Assert that the returned content is empty
	if content != "" {
		t.Errorf("Expected empty content, got %q", content)
	}

	// Assert that FilesProcessed is 0
	if stats.FilesProcessed != 0 {
		t.Errorf("Expected FilesProcessed to be 0, got %d", stats.FilesProcessed)
	}
}

// TestWriteToFileWithDirectoryCreation tests that WriteToFile creates parent directories
func TestWriteToFileWithDirectoryCreation(t *testing.T) {
	// Create a temporary base directory for testing
	tmpDir, err := os.MkdirTemp("", "handoff-writetofile-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if cleanErr := os.RemoveAll(tmpDir); cleanErr != nil {
			t.Logf("Failed to clean up temp directory: %v", cleanErr)
		}
	}()

	// Create a nested path that doesn't exist yet
	nestedDirPath := filepath.Join(tmpDir, "level1", "level2", "level3")
	filePath := filepath.Join(nestedDirPath, "test.txt")
	
	// Content to write
	content := "This is test content for directory creation"
	
	// Write to the file with non-existent directories
	err = WriteToFile(content, filePath)
	if err != nil {
		t.Errorf("WriteToFile() failed with nested directories: %v", err)
	}
	
	// Verify the file was created with the correct content
	readContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to read written file: %v", err)
	}
	
	if string(readContent) != content {
		t.Errorf("Written content doesn't match expected. Got %q, want %q", string(readContent), content)
	}
	
	// Verify all parent directories were created
	for _, dir := range []string{
		filepath.Join(tmpDir, "level1"),
		filepath.Join(tmpDir, "level1", "level2"),
		filepath.Join(tmpDir, "level1", "level2", "level3"),
	} {
		if info, err := os.Stat(dir); err != nil || !info.IsDir() {
			t.Errorf("Parent directory %s was not created or is not a directory", dir)
		}
	}
}