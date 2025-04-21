package handoff

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNewConfig tests the NewConfig function
func TestNewConfig(t *testing.T) {
	config := NewConfig()

	// Verify the default values are set correctly
	if config.Verbose != false {
		t.Errorf("Default Verbose value should be false, got %v", config.Verbose)
	}

	if config.Format != "<{path}>\n```\n{content}\n```\n</{path}>\n\n" {
		t.Errorf("Default Format value is incorrect, got %q", config.Format)
	}

	if config.Include != "" {
		t.Errorf("Default Include value should be empty, got %q", config.Include)
	}

	if config.Exclude != "" {
		t.Errorf("Default Exclude value should be empty, got %q", config.Exclude)
	}

	if config.ExcludeNamesStr != "" {
		t.Errorf("Default ExcludeNamesStr value should be empty, got %q", config.ExcludeNamesStr)
	}

	if len(config.IncludeExts) != 0 {
		t.Errorf("Default IncludeExts should be empty, got %v", config.IncludeExts)
	}

	if len(config.ExcludeExts) != 0 {
		t.Errorf("Default ExcludeExts should be empty, got %v", config.ExcludeExts)
	}

	if len(config.ExcludeNames) != 0 {
		t.Errorf("Default ExcludeNames should be empty, got %v", config.ExcludeNames)
	}
}

// TestIsGitIgnored tests the isGitIgnored function (now internal to lib)
func TestIsGitIgnored(t *testing.T) {
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

	// Since isGitIgnored is now internal, we need to test it indirectly
	// We'll test it through ProcessFile which uses isGitIgnored internally

	// Create a hidden file that should be ignored
	hiddenFile := filepath.Join(tmpDir, ".hidden")
	if err := os.WriteFile(hiddenFile, []byte("hidden"), 0644); err != nil {
		t.Fatalf("Failed to create hidden file: %v", err)
	}

	// Create a visible file that should not be ignored
	visibleFile := filepath.Join(tmpDir, "visible")
	if err := os.WriteFile(visibleFile, []byte("visible"), 0644); err != nil {
		t.Fatalf("Failed to create visible file: %v", err)
	}

	// Setup test components
	config := NewConfig()
	logger := NewLogger(false)
	processor := func(file string, content []byte) string {
		return "processed"
	}

	// Test with hidden file - should be ignored and return empty string
	result := ProcessFile(hiddenFile, logger, config, processor)
	if result != "" {
		t.Errorf("Hidden file %s should have been ignored, but got result: %s", hiddenFile, result)
	}

	// Test with visible file - should be processed
	result = ProcessFile(visibleFile, logger, config, processor)
	if result != "processed" {
		t.Errorf("Visible file %s should have been processed, but got empty result", visibleFile)
	}
}

// TestGetFilesFromDir tests file collection functionality indirectly
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
	content, err := ProcessProject([]string{tmpDir}, config)
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
		t.Errorf("Content should not include .hidden file")
	}
}

// TestProcessorFunc tests the ProcessorFunc functionality
func TestProcessorFunc(t *testing.T) {
	// Create a processor function similar to the default one
	processor := func(file string, content []byte) string {
		// Use isBinaryFile indirectly through ProcessFile
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
}

// TestIsBinaryFile tests the binary file detection functionality indirectly
func TestIsBinaryFile(t *testing.T) {
	// Create a temporary file for testing
	tmpDir, err := os.MkdirTemp("", "handoff-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if cleanErr := os.RemoveAll(tmpDir); cleanErr != nil {
			t.Logf("Failed to clean up temp directory: %v", cleanErr)
		}
	}()

	// Test with text content - create a text file
	textFile := filepath.Join(tmpDir, "text.txt")
	textContent := []byte("This is a text file with normal characters.")
	if err := os.WriteFile(textFile, textContent, 0644); err != nil {
		t.Fatalf("Failed to create text file: %v", err)
	}

	// Test with binary content - create a binary file
	binaryFile := filepath.Join(tmpDir, "binary.bin")
	binaryContent := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x00, 0x57, 0x6f, 0x72, 0x6c, 0x64} // Contains null byte
	if err := os.WriteFile(binaryFile, binaryContent, 0644); err != nil {
		t.Fatalf("Failed to create binary file: %v", err)
	}

	// Test with high percentage of non-printable characters
	nonPrintableFile := filepath.Join(tmpDir, "nonprintable.bin")
	nonPrintableContent := make([]byte, 100)
	for i := range nonPrintableContent {
		nonPrintableContent[i] = byte(i % 32)
	}
	if err := os.WriteFile(nonPrintableFile, nonPrintableContent, 0644); err != nil {
		t.Fatalf("Failed to create non-printable file: %v", err)
	}

	// Create test components
	config := NewConfig()
	logger := NewLogger(false)

	// Setup a processor that will be called only for non-binary files
	called := false
	processor := func(file string, content []byte) string {
		called = true
		return "processed"
	}

	// Test text file - should be processed
	called = false
	result := ProcessFile(textFile, logger, config, processor)
	if !called || result != "processed" {
		t.Errorf("Text file was not processed correctly")
	}

	// Test binary file - should be skipped
	called = false
	result = ProcessFile(binaryFile, logger, config, processor)
	if called || result != "" {
		t.Errorf("Binary file was processed when it should have been skipped")
	}

	// Test non-printable file - should be skipped
	called = false
	result = ProcessFile(nonPrintableFile, logger, config, processor)
	if called || result != "" {
		t.Errorf("Non-printable file was processed when it should have been skipped")
	}
}

// TestLogger tests the Logger functionality
func TestLogger(t *testing.T) {
	// Create a temporary capture of stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Create loggers
	verboseLogger := NewLogger(true)
	quietLogger := NewLogger(false)

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
	if err := w.Close(); err != nil {
		t.Errorf("Failed to close writer: %v", err)
	}
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

// TestWrapInContext tests the WrapInContext function
func TestWrapInContext(t *testing.T) {
	input := "test content"
	expected := "<context>\ntest content</context>"

	result := WrapInContext(input)
	if result != expected {
		t.Errorf("WrapInContext(%q) = %q, want %q", input, result, expected)
	}
}

// TestEstimateTokenCount tests the estimateTokenCount function indirectly through CalculateStatistics
func TestCalculateStatistics(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedChars  int
		expectedLines  int
		expectedTokens int
	}{
		{
			name:           "Empty string",
			input:          "",
			expectedChars:  0,
			expectedLines:  1,
			expectedTokens: 0,
		},
		{
			name:           "Single word",
			input:          "hello",
			expectedChars:  5,
			expectedLines:  1,
			expectedTokens: 1,
		},
		{
			name:           "Multiple words with spaces",
			input:          "hello world example",
			expectedChars:  19,
			expectedLines:  1,
			expectedTokens: 3,
		},
		{
			name:           "Words with mixed whitespace",
			input:          "hello\nworld\texample  test",
			expectedChars:  25, // Updated to fix test
			expectedLines:  2,
			expectedTokens: 4,
		},
		{
			name:           "Leading and trailing whitespace",
			input:          "  hello world  ",
			expectedChars:  15,
			expectedLines:  1,
			expectedTokens: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			chars, lines, tokens := CalculateStatistics(tc.input)

			if chars != tc.expectedChars {
				t.Errorf("Character count: got %d, want %d", chars, tc.expectedChars)
			}

			if lines != tc.expectedLines {
				t.Errorf("Line count: got %d, want %d", lines, tc.expectedLines)
			}

			if tokens != tc.expectedTokens {
				t.Errorf("Token count: got %d, want %d", tokens, tc.expectedTokens)
			}
		})
	}
}

// TestProcessConfig tests the ProcessConfig method of Config
func TestProcessConfig(t *testing.T) {
	tests := []struct {
		name          string
		config        Config
		wantInclude   []string
		wantExclude   []string
		wantExclNames []string
	}{
		{
			name: "Empty config",
			config: Config{
				Include:         "",
				Exclude:         "",
				ExcludeNamesStr: "",
			},
			wantInclude:   nil,
			wantExclude:   nil,
			wantExclNames: nil,
		},
		{
			name: "Include extensions with dots",
			config: Config{
				Include: ".go,.txt,.md",
			},
			wantInclude:   []string{".go", ".txt", ".md"},
			wantExclude:   nil,
			wantExclNames: nil,
		},
		{
			name: "Include extensions without dots",
			config: Config{
				Include: "go,txt,md",
			},
			wantInclude:   []string{".go", ".txt", ".md"},
			wantExclude:   nil,
			wantExclNames: nil,
		},
		{
			name: "Include extensions mixed format",
			config: Config{
				Include: ".go,txt,.md",
			},
			wantInclude:   []string{".go", ".txt", ".md"},
			wantExclude:   nil,
			wantExclNames: nil,
		},
		{
			name: "Exclude extensions with dots",
			config: Config{
				Exclude: ".exe,.bin,.obj",
			},
			wantInclude:   nil,
			wantExclude:   []string{".exe", ".bin", ".obj"},
			wantExclNames: nil,
		},
		{
			name: "Exclude extensions without dots",
			config: Config{
				Exclude: "exe,bin,obj",
			},
			wantInclude:   nil,
			wantExclude:   []string{".exe", ".bin", ".obj"},
			wantExclNames: nil,
		},
		{
			name: "Exclude names",
			config: Config{
				ExcludeNamesStr: "package-lock.json,yarn.lock,.DS_Store",
			},
			wantInclude:   nil,
			wantExclude:   nil,
			wantExclNames: []string{"package-lock.json", "yarn.lock", ".DS_Store"},
		},
		{
			name: "All filters with whitespace",
			config: Config{
				Include:         " .go , txt , .md ",
				Exclude:         " .exe , bin , .obj ",
				ExcludeNamesStr: " package-lock.json , yarn.lock , .DS_Store ",
			},
			wantInclude:   []string{".go", ".txt", ".md"},
			wantExclude:   []string{".exe", ".bin", ".obj"},
			wantExclNames: []string{"package-lock.json", "yarn.lock", ".DS_Store"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Process the config
			tt.config.ProcessConfig()

			// Check include extensions
			if !sliceEqual(tt.config.IncludeExts, tt.wantInclude) {
				t.Errorf("ProcessConfig() IncludeExts = %v, want %v", tt.config.IncludeExts, tt.wantInclude)
			}

			// Check exclude extensions
			if !sliceEqual(tt.config.ExcludeExts, tt.wantExclude) {
				t.Errorf("ProcessConfig() ExcludeExts = %v, want %v", tt.config.ExcludeExts, tt.wantExclude)
			}

			// Check exclude names
			if !sliceEqual(tt.config.ExcludeNames, tt.wantExclNames) {
				t.Errorf("ProcessConfig() ExcludeNames = %v, want %v", tt.config.ExcludeNames, tt.wantExclNames)
			}
		})
	}
}

// Helper function to compare slices
func sliceEqual(a, b []string) bool {
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

// TestShouldProcess tests the shouldProcess function
func TestShouldProcess(t *testing.T) {
	tests := []struct {
		name       string
		file       string
		config     *Config
		wantResult bool
	}{
		{
			name:       "No filters",
			file:       "/path/to/file.txt",
			config:     &Config{},
			wantResult: true,
		},
		{
			name:       "Include match",
			file:       "/path/to/file.go",
			config:     &Config{IncludeExts: []string{".go", ".md"}},
			wantResult: true,
		},
		{
			name:       "Include no match",
			file:       "/path/to/file.txt",
			config:     &Config{IncludeExts: []string{".go", ".md"}},
			wantResult: false,
		},
		{
			name:       "Exclude match",
			file:       "/path/to/file.exe",
			config:     &Config{ExcludeExts: []string{".exe", ".bin"}},
			wantResult: false,
		},
		{
			name:       "Exclude no match",
			file:       "/path/to/file.txt",
			config:     &Config{ExcludeExts: []string{".exe", ".bin"}},
			wantResult: true,
		},
		{
			name:       "Exclude name match",
			file:       "/path/to/package-lock.json",
			config:     &Config{ExcludeNames: []string{"package-lock.json", "yarn.lock"}},
			wantResult: false,
		},
		{
			name:       "Exclude name no match",
			file:       "/path/to/file.json",
			config:     &Config{ExcludeNames: []string{"package-lock.json", "yarn.lock"}},
			wantResult: true,
		},
		{
			name: "Include and exclude - included file",
			file: "/path/to/file.go",
			config: &Config{
				IncludeExts: []string{".go", ".md"},
				ExcludeExts: []string{".exe", ".bin"},
			},
			wantResult: true,
		},
		{
			name: "Include and exclude - excluded file",
			file: "/path/to/file.exe",
			config: &Config{
				IncludeExts: []string{".go", ".md", ".exe"},
				ExcludeExts: []string{".exe", ".bin"},
			},
			wantResult: false,
		},
		{
			name: "Exclude name takes precedence over include ext",
			file: "/path/to/special.go",
			config: &Config{
				IncludeExts:  []string{".go", ".md"},
				ExcludeNames: []string{"special.go"},
			},
			wantResult: false,
		},
		{
			name: "Case sensitivity in extensions",
			file: "/path/to/file.GO",
			config: &Config{
				IncludeExts: []string{".go"},
			},
			wantResult: true, // Extension comparison is case-insensitive
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldProcess(tt.file, tt.config)
			if result != tt.wantResult {
				t.Errorf("shouldProcess(%q, %+v) = %v, want %v", tt.file, tt.config, result, tt.wantResult)
			}
		})
	}
}

// TestWriteToFile tests the WriteToFile function
func TestWriteToFile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "handoff-write-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if cleanErr := os.RemoveAll(tmpDir); cleanErr != nil {
			t.Logf("Failed to clean up temp directory: %v", cleanErr)
		}
	}()

	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "Write simple content",
			content: "Test content",
			wantErr: false,
		},
		{
			name:    "Write empty content",
			content: "",
			wantErr: false,
		},
		{
			name:    "Write multi-line content",
			content: "Line 1\nLine 2\nLine 3",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a simple filename without problematic characters
			filename := strings.ReplaceAll(tt.name, " ", "-")
			filename = strings.ReplaceAll(filename, "/", "-")
			filePath := filepath.Join(tmpDir, fmt.Sprintf("test-file-%s.txt", filename))

			// Write the content
			err := WriteToFile(tt.content, filePath)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteToFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify content was written correctly
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Errorf("Failed to read written file: %v", err)
					return
				}

				if string(content) != tt.content {
					t.Errorf("Written content = %q, want %q", string(content), tt.content)
				}

				// Verify file permissions (0644)
				info, err := os.Stat(filePath)
				if err != nil {
					t.Errorf("Failed to stat written file: %v", err)
					return
				}

				// Check permission bits (accounting for umask)
				// We only check if it's readable, as the exact permission bits
				// may be affected by the system's umask
				if info.Mode().Perm()&0444 == 0 {
					t.Errorf("File permissions %v do not include read permission", info.Mode().Perm())
				}
			}
		})
	}

	// Test writing to a non-existent directory
	t.Run("Write to non-existent directory", func(t *testing.T) {
		nonExistentPath := filepath.Join(tmpDir, "non/existent/dir/file.txt")
		err := WriteToFile("content", nonExistentPath)
		if err == nil {
			t.Errorf("WriteToFile() to non-existent directory succeeded, want error")
		}
	})
}
