package main

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/cover"
)

func TestParseCoverageFromFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "coverage-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a sample coverage profile
	// This is a simplified coverage profile with known coverage values
	sampleCoverageContent := `mode: set
example.com/pkg/file1.go:10.20,15.3 3 1
example.com/pkg/file1.go:20.30,25.3 3 0
example.com/pkg/file1.go:30.40,35.3 3 1
example.com/pkg/file2.go:10.20,15.3 3 1
example.com/pkg/file2.go:20.30,25.3 3 1
`
	// In this sample:
	// - Total statements: 15 (5 blocks, 3 statements each)
	// - Covered statements: 12 (4 blocks with Count > 0, 3 statements each)
	// - Expected coverage: 12/15 * 100 = 80%

	coverageFilePath := filepath.Join(tempDir, "coverage.out")
	if err := os.WriteFile(coverageFilePath, []byte(sampleCoverageContent), 0644); err != nil {
		t.Fatalf("Failed to write sample coverage file: %v", err)
	}

	// Test parsing coverage from the file
	coverage, err := ParseCoverageFromFile(coverageFilePath)
	if err != nil {
		t.Fatalf("Failed to parse coverage from file: %v", err)
	}

	expectedCoverage := 80.0
	if coverage != expectedCoverage {
		t.Errorf("Expected coverage %.2f%%, got %.2f%%", expectedCoverage, coverage)
	}

	// Test for non-existent file
	_, err = ParseCoverageFromFile("non-existent-file.out")
	if err == nil {
		t.Error("Expected error for non-existent file but got nil")
	}
}

func TestCalculateCoverage(t *testing.T) {
	testCases := []struct {
		name             string
		profiles         []*cover.Profile
		expectedCoverage float64
	}{
		{
			name: "Multiple files with different coverage",
			profiles: []*cover.Profile{
				{
					FileName: "file1.go",
					Blocks: []cover.ProfileBlock{
						{StartLine: 10, StartCol: 20, EndLine: 15, EndCol: 30, NumStmt: 3, Count: 1},
						{StartLine: 20, StartCol: 20, EndLine: 25, EndCol: 30, NumStmt: 3, Count: 0},
					},
				},
				{
					FileName: "file2.go",
					Blocks: []cover.ProfileBlock{
						{StartLine: 10, StartCol: 20, EndLine: 15, EndCol: 30, NumStmt: 3, Count: 1},
						{StartLine: 20, StartCol: 20, EndLine: 25, EndCol: 30, NumStmt: 3, Count: 1},
					},
				},
			},
			expectedCoverage: 75.0, // 9/12 = 75%
		},
		{
			name: "All blocks covered",
			profiles: []*cover.Profile{
				{
					FileName: "file1.go",
					Blocks: []cover.ProfileBlock{
						{StartLine: 10, StartCol: 20, EndLine: 15, EndCol: 30, NumStmt: 3, Count: 1},
						{StartLine: 20, StartCol: 20, EndLine: 25, EndCol: 30, NumStmt: 3, Count: 1},
					},
				},
			},
			expectedCoverage: 100.0, // 6/6 = 100%
		},
		{
			name: "No blocks covered",
			profiles: []*cover.Profile{
				{
					FileName: "file1.go",
					Blocks: []cover.ProfileBlock{
						{StartLine: 10, StartCol: 20, EndLine: 15, EndCol: 30, NumStmt: 3, Count: 0},
						{StartLine: 20, StartCol: 20, EndLine: 25, EndCol: 30, NumStmt: 3, Count: 0},
					},
				},
			},
			expectedCoverage: 0.0, // 0/6 = 0%
		},
		{
			name:             "Empty profiles",
			profiles:         []*cover.Profile{},
			expectedCoverage: 0.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			coverage := calculateCoverage(tc.profiles)
			if coverage != tc.expectedCoverage {
				t.Errorf("Expected coverage %.2f%%, got %.2f%%", tc.expectedCoverage, coverage)
			}
		})
	}
}

func TestIsToolCoverOutput(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "cover-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testCases := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name: "Valid coverage profile",
			content: `mode: set
example.com/pkg/file1.go:10.20,15.3 3 1
example.com/pkg/file1.go:20.30,25.3 3 0`,
			expected: false,
		},
		{
			name: "go tool cover -func output",
			content: `file1.go:	 FunctionOne		100.0%
file2.go:	 FunctionTwo		75.0%
total:	(statements)	85.7%`,
			expected: true,
		},
		{
			name:     "Empty file",
			content:  "",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filePath := filepath.Join(tempDir, tc.name+".txt")
			if err := os.WriteFile(filePath, []byte(tc.content), 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			result := isToolCoverOutput(filePath)
			if result != tc.expected {
				t.Errorf("Expected isToolCoverOutput to return %v, got %v", tc.expected, result)
			}
		})
	}
}
