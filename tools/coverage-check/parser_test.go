package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestParseCoverageOutput(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedOutput float64
		expectError    bool
	}{
		{
			name: "Valid coverage output",
			input: `
file1.go:	 FunctionOne		100.0%
file2.go:	 FunctionTwo		75.0%
total:	(statements)	85.7%
`,
			expectedOutput: 85.7,
			expectError:    false,
		},
		{
			name: "Valid coverage output with different formatting",
			input: `
file1.go:	 FunctionOne		100.0%
file2.go:	 FunctionTwo		75.0%
total:			(statements)			90.5%
`,
			expectedOutput: 90.5,
			expectError:    false,
		},
		{
			name: "No total line",
			input: `
file1.go:	 FunctionOne		100.0%
file2.go:	 FunctionTwo		75.0%
`,
			expectedOutput: 0,
			expectError:    true,
		},
		{
			name:           "Empty input",
			input:          "",
			expectedOutput: 0,
			expectError:    true,
		},
		{
			name:           "Invalid percentage format",
			input:          "total:	(statements)	invalid%",
			expectedOutput: 0,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseCoverageOutput(tc.input)

			// Check error expectation
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// If we don't expect an error, check the result
			if !tc.expectError {
				if result != tc.expectedOutput {
					t.Errorf("Expected coverage %.2f%%, got %.2f%%", tc.expectedOutput, result)
				}
			}
		})
	}
}

func TestParseCoverageFromFile(t *testing.T) {
	// Create a temporary file with sample coverage output
	tempDir, err := ioutil.TempDir("", "coverage-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a sample coverage.out file
	sampleCoverageContent := `mode: atomic
github.com/example/pkg/file1.go:10.40,12.2 1 1
github.com/example/pkg/file1.go:14.40,16.2 1 0
github.com/example/pkg/file2.go:20.40,22.2 1 1
`
	coverageFilePath := filepath.Join(tempDir, "coverage.out")
	if err := ioutil.WriteFile(coverageFilePath, []byte(sampleCoverageContent), 0644); err != nil {
		t.Fatalf("Failed to write sample coverage file: %v", err)
	}

	// Test for non-existent file
	_, err = ParseCoverageFromFile("non-existent-file.out")
	if err == nil {
		t.Error("Expected error for non-existent file but got nil")
	}

	// Note: We can't fully test the successful case without mocking go tool cover,
	// which would significantly complicate this test. We've covered the file existence
	// check, and the ParseCoverageOutput function is tested separately.
}
