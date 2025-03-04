package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
	processPath(tmpFile.Name(), &builder)

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
	processPath("non-existent-path", &builder)
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