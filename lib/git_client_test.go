package handoff

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestMockGitClient tests the basic functionality of the MockGitClient
func TestMockGitClient(t *testing.T) {
	// Test availability settings
	t.Run("Availability", func(t *testing.T) {
		availableClient := NewMockGitClient(true)
		if !availableClient.IsAvailable() {
			t.Error("MockGitClient with available=true should report IsAvailable() as true")
		}

		unavailableClient := NewMockGitClient(false)
		if unavailableClient.IsAvailable() {
			t.Error("MockGitClient with available=false should report IsAvailable() as false")
		}
	})

	t.Run("GitIgnore", func(t *testing.T) {
		client := NewMockGitClient(true)

		// Configure mock to consider certain files ignored
		client.SetIgnoredFiles(map[string]bool{
			"/path/to/ignored.txt": true,
			"/path/to/tracked.txt": false,
		})

		// Test configured ignored file
		if !client.IsGitIgnored("/path/to/ignored.txt") {
			t.Error("File explicitly set as ignored should be reported as ignored")
		}

		// Test configured non-ignored file
		if client.IsGitIgnored("/path/to/tracked.txt") {
			t.Error("File explicitly set as not ignored should be reported as not ignored")
		}

		// Test fall-back behavior for hidden files
		if !client.IsGitIgnored("/path/to/.hidden") {
			t.Error("Hidden file not explicitly configured should be treated as ignored")
		}

		// Test fall-back behavior for regular files
		if client.IsGitIgnored("/path/to/regular.go") {
			t.Error("Regular file not explicitly configured should be treated as not ignored")
		}
	})

	t.Run("GetGitFiles", func(t *testing.T) {
		// Test with git available
		availableClient := NewMockGitClient(true)
		expectedFiles := []string{"/dir/file1.go", "/dir/file2.go"}
		availableClient.SetFilesInDir("/dir", expectedFiles)

		files, err := availableClient.GetGitFiles("/dir")
		if err != nil {
			t.Errorf("GetGitFiles with configured directory should not return error: %v", err)
		}
		if !stringSlicesEqual(files, expectedFiles) {
			t.Errorf("GetGitFiles returned %v, expected %v", files, expectedFiles)
		}

		// Test with git available but directory not configured
		files, err = availableClient.GetGitFiles("/unknown")
		if err == nil {
			t.Error("GetGitFiles with unconfigured directory should return error")
		}
		if len(files) > 0 {
			t.Errorf("GetGitFiles with error should return empty list, got %v", files)
		}

		// Test with git unavailable
		unavailableClient := NewMockGitClient(false)
		unavailableClient.SetFilesInDir("/dir", expectedFiles)

		files, err = unavailableClient.GetGitFiles("/dir")
		if err == nil {
			t.Error("GetGitFiles with git unavailable should return error")
		}
		if len(files) > 0 {
			t.Errorf("GetGitFiles with git unavailable should return empty list, got %v", files)
		}
	})
}

// TestRealGitClientWithMocks tests the RealGitClient using MockGitClient for comparison
func TestGitClientIntegration(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "handoff-git-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if cleanErr := os.RemoveAll(tmpDir); cleanErr != nil {
			t.Logf("Failed to clean up test directory: %v", cleanErr)
		}
	}()

	// Create test files
	testFiles := []string{"file1.go", "file2.txt", ".hidden"}
	for _, file := range testFiles {
		filePath := filepath.Join(tmpDir, file)
		if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Test using both real client and mock client for comparison
	t.Run("Compare clients with same config", func(t *testing.T) {
		// Create both client types
		realClient := NewRealGitClient()
		mockClient := NewMockGitClient(realClient.IsAvailable())

		// Set up mock client to simulate the real behavior
		// Hidden files should be treated as ignored
		mockClient.SetIgnoredFiles(map[string]bool{
			filepath.Join(tmpDir, ".hidden"): true,
		})

		// Test IsGitIgnored behavior
		for _, file := range testFiles {
			filePath := filepath.Join(tmpDir, file)
			realResult := realClient.IsGitIgnored(filePath)
			mockResult := mockClient.IsGitIgnored(filePath)

			t.Logf("Testing IsGitIgnored on %s: real=%v, mock=%v", file, realResult, mockResult)

			// Hidden files should be treated as ignored by both implementations
			if strings.HasPrefix(file, ".") {
				if !realResult {
					t.Errorf("RealGitClient should treat hidden file %s as ignored", file)
				}
				if !mockResult {
					t.Errorf("MockGitClient should treat hidden file %s as ignored", file)
				}
			}
		}
	})

	// Test with different configurations to the client
	t.Run("Different client configurations", func(t *testing.T) {
		// Test with standard config using real git
		standardConfig := NewConfig()

		// Test with mock git that's unavailable
		mockConfig := NewConfig(WithGitClient(NewMockGitClient(false)))

		// ProcessProject with both configs
		standardContent, _, standardErr := ProcessProject([]string{tmpDir}, standardConfig)
		mockContent, _, mockErr := ProcessProject([]string{tmpDir}, mockConfig)

		// Both should work, just potentially with different file discovery approaches
		if standardErr != nil && !errors.Is(standardErr, ErrNoFilesProcessed) {
			t.Errorf("ProcessProject with standard config failed: %v", standardErr)
		}
		if mockErr != nil && !errors.Is(mockErr, ErrNoFilesProcessed) {
			t.Errorf("ProcessProject with mock config failed: %v", mockErr)
		}

		// Both should handle hidden files the same - they should be excluded
		if strings.Contains(standardContent, ".hidden") || strings.Contains(mockContent, ".hidden") {
			t.Error("Hidden files should be excluded in both real and mock implementations")
		}
	})
}

// stringSlicesEqual is a helper function to compare string slices
func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
