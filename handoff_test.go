package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	
	handoff "github.com/phrazzld/handoff/lib"
)

// TestIsGitIgnored tests the isGitIgnored function
func TestIsGitIgnored(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "handoff-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test with a hidden file (should be considered ignored)
	hiddenFile := filepath.Join(tmpDir, ".hidden")
	if err := os.WriteFile(hiddenFile, []byte("hidden"), 0644); err != nil {
		t.Fatalf("Failed to create hidden file: %v", err)
	}

	if !isGitIgnored(hiddenFile) {
		t.Errorf("Hidden file %s should be considered ignored", hiddenFile)
	}

	// Test with a visible file (should not be considered ignored when not in a git repo)
	visibleFile := filepath.Join(tmpDir, "visible")
	if err := os.WriteFile(visibleFile, []byte("visible"), 0644); err != nil {
		t.Fatalf("Failed to create visible file: %v", err)
	}

	if isGitIgnored(visibleFile) {
		t.Errorf("Visible file %s should not be considered ignored", visibleFile)
	}
}

// TestGetFilesFromDir tests the getFilesFromDir function
func TestGetFilesFromDir(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "handoff-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

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

	// Get files from the directory
	foundFiles, err := getFilesFromDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to get files from directory: %v", err)
	}

	// We should find the visible files, but not the hidden one
	if len(foundFiles) < 4 { // file1.txt, file2.go, subdir/subfile1.txt, subdir/subfile2.go
		t.Errorf("Expected at least 4 files, but found %d: %v", len(foundFiles), foundFiles)
	}

	// The hidden file should not be included
	for _, f := range foundFiles {
		if strings.HasSuffix(f, ".hidden") {
			t.Errorf("Hidden file %s should not be included", f)
		}
	}
}

// TestProcessorFunc tests the processor function in isolation
func TestProcessorFunc(t *testing.T) {
	// Create a processor function similar to the default one
	processor := func(file string, content []byte) string {
		if isBinaryFile(content) {
			return ""
		}
		return fmt.Sprintf("<%s>\n```\n%s\n```\n</%s>\n\n", file, string(content), file)
	}

	// Test with text content
	filePath := "/path/to/file.txt"
	fileContent := []byte("Test content")
	
	result := processor(filePath, fileContent)

	// Check the result
	expectedPrefix := "<" + filePath + ">\n```\n"
	expectedSuffix := "\n```\n</" + filePath + ">\n\n"
	if !strings.HasPrefix(result, expectedPrefix) {
		t.Errorf("Expected result to start with %q, but got %q", expectedPrefix, result)
	}
	if !strings.HasSuffix(result, expectedSuffix) {
		t.Errorf("Expected result to end with %q, but got %q", expectedSuffix, result)
	}
	if !strings.Contains(result, string(fileContent)) {
		t.Errorf("Expected result to contain %q, but got %q", string(fileContent), result)
	}

	// Test with binary content
	binaryContent := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x00, 0x57, 0x6f, 0x72, 0x6c, 0x64} // Contains null byte
	result = processor(filePath, binaryContent)
	if result != "" {
		t.Errorf("Expected empty result for binary content, but got %q", result)
	}
}

// TestProcessPath tests the processPath function
func TestProcessPath(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "handoff-test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write content to the file
	content := "Test content"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Process the path (file)
	var builder strings.Builder
	logger := handoff.NewLogger(false) // Non-verbose logger for testing
	config := handoff.NewConfig()
	processPath(tmpFile.Name(), &builder, config, logger)

	// Check the result
	result := builder.String()
	expectedPrefix := "<" + tmpFile.Name() + ">\n```\n"
	expectedSuffix := "\n```\n</" + tmpFile.Name() + ">\n\n"
	if !strings.HasPrefix(result, expectedPrefix) {
		t.Errorf("Expected result to start with %q, but got %q", expectedPrefix, result)
	}
	if !strings.HasSuffix(result, expectedSuffix) {
		t.Errorf("Expected result to end with %q, but got %q", expectedSuffix, result)
	}
	if !strings.Contains(result, content) {
		t.Errorf("Expected result to contain %q, but got %q", content, result)
	}

	// Test with a non-existent path
	builder.Reset()
	processPath("non-existent-path", &builder, config, logger)
	if builder.String() != "" {
		t.Errorf("Expected empty result for non-existent path, but got %q", builder.String())
	}
	
	// Note: The top-level context tags are added in main(), not in processPath()
	// so we don't test for them here
}

// TestIsBinaryFile tests the isBinaryFile function
func TestIsBinaryFile(t *testing.T) {
	// Test with text content
	textContent := []byte("This is a text file with normal characters.")
	if isBinaryFile(textContent) {
		t.Errorf("Text content incorrectly identified as binary")
	}

	// Test with binary content (content with null bytes)
	binaryContent := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x00, 0x57, 0x6f, 0x72, 0x6c, 0x64}
	if !isBinaryFile(binaryContent) {
		t.Errorf("Binary content not identified as binary")
	}

	// Test with high percentage of non-printable characters
	nonPrintableContent := make([]byte, 100)
	for i := range nonPrintableContent {
		nonPrintableContent[i] = byte(i % 32)
	}
	if !isBinaryFile(nonPrintableContent) {
		t.Errorf("Content with high percentage of non-printable characters not identified as binary")
	}

	// Test with whitespace characters (they should be considered printable)
	whitespaceContent := []byte("\n\r\t ")
	if isBinaryFile(whitespaceContent) {
		t.Errorf("Whitespace content incorrectly identified as binary")
	}
}

// TestLogger tests the Logger functionality
func TestLogger(t *testing.T) {
	// Create a temporary capture of stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	
	// Create loggers
	verboseLogger := handoff.NewLogger(true)
	quietLogger := handoff.NewLogger(false)
	
	// Log some messages
	verboseLogger.Info("Info message")
	verboseLogger.Warn("Warning message")
	verboseLogger.Error("Error message")
	verboseLogger.Verbose("Verbose message")
	
	quietLogger.Info("Info message from quiet logger")
	quietLogger.Warn("Warning message from quiet logger")
	quietLogger.Error("Error message from quiet logger")
	quietLogger.Verbose("This verbose message should not appear")
	
	// Close the writer and restore stderr
	w.Close()
	os.Stderr = oldStderr
	
	// Read the output
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()
	
	// Check that messages were logged correctly
	if !strings.Contains(output, "Info message") {
		t.Error("Info message not found in logger output")
	}
	if !strings.Contains(output, "warning: Warning message") {
		t.Error("Warning message not found in logger output")
	}
	if !strings.Contains(output, "error: Error message") {
		t.Error("Error message not found in logger output")
	}
	if !strings.Contains(output, "Verbose message") {
		t.Error("Verbose message not found in logger output")
	}
	
	// Check that quiet logger suppresses verbose messages
	if strings.Contains(output, "This verbose message should not appear") {
		t.Error("Verbose message from quiet logger should not appear")
	}
}

// TestWrapInContext tests the wrapInContext function
func TestWrapInContext(t *testing.T) {
	input := "test content"
	expected := "<context>\ntest content</context>"
	
	result := wrapInContext(input)
	if result != expected {
		t.Errorf("wrapInContext(%q) = %q, want %q", input, result, expected)
	}
}

// TestLogStatistics tests the logStatistics function
func TestLogStatistics(t *testing.T) {
	// Create a temporary capture of stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	
	// Create a logger
	logger := handoff.NewLogger(true)
	
	// Mock content and config
	content := "Line 1\nLine 2\nLine 3\nThis is a test of the statistics function.\n"
	config := handoff.NewConfig()
	config.Verbose = true
	
	// Call logStatistics
	logStatistics(content, 3, 5, config, logger)
	
	// Close the writer and restore stderr
	w.Close()
	os.Stderr = oldStderr
	
	// Read the output
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()
	
	// Check that statistics were logged correctly
	statsToCheck := []string{
		"Handoff complete:",
		"Files: 3",
		"Lines: 5", // 4 newlines + 1 = 5 lines
		"Characters: " + fmt.Sprintf("%d", len(content)),
		"Estimated tokens: " + fmt.Sprintf("%d", estimateTokenCount(content)),
		"Processed 3/5 files",
	}
	
	for _, stat := range statsToCheck {
		if !strings.Contains(output, stat) {
			t.Errorf("Expected output to contain %q, but got %q", stat, output)
		}
	}
	
	// Test with verbose config (since DryRun was moved)
	oldStderr = os.Stderr
	r, w, _ = os.Pipe()
	os.Stderr = w
	
	verboseConfig := handoff.NewConfig()
	verboseConfig.Verbose = true
	
	logStatistics(content, 3, 5, verboseConfig, logger)
	
	w.Close()
	os.Stderr = oldStderr
	
	buf.Reset()
	_, _ = io.Copy(&buf, r)
	dryRunOutput := buf.String()
	
	if !strings.Contains(dryRunOutput, "Processed 3/5 files") {
		t.Errorf("Expected output to mention processed files, but got %q", dryRunOutput)
	}
}

// TestEstimateTokenCount tests the estimateTokenCount function
func TestEstimateTokenCount(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "Single word",
			input:    "hello",
			expected: 1,
		},
		{
			name:     "Multiple words with spaces",
			input:    "hello world example",
			expected: 3,
		},
		{
			name:     "Words with mixed whitespace",
			input:    "hello\nworld\texample  test",
			expected: 4,
		},
		{
			name:     "Leading and trailing whitespace",
			input:    "  hello world  ",
			expected: 2,
		},
		{
			name:     "Symbols and punctuation",
			input:    "hello, world! This is a test.",
			expected: 6, // Punctuation attaches to words in our simple tokenizer
		},
		{
			name:     "Code-like text",
			input:    "func main() { fmt.Println(\"Hello\") }",
			expected: 5, // Punctuation and symbols stay attached to words in our tokenizer
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := estimateTokenCount(tc.input)
			if result != tc.expected {
				t.Errorf("estimateTokenCount(%q) = %d, want %d", tc.input, result, tc.expected)
			}
		})
	}
}