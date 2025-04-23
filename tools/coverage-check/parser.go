package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/tools/cover"
)

// ParseCoverageFromFile parses a coverage profile file and calculates the coverage percentage
func ParseCoverageFromFile(filepath string) (float64, error) {
	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return 0, fmt.Errorf("coverage file %s does not exist", filepath)
	}

	// Parse the coverage profile using the cover package
	profiles, err := cover.ParseProfiles(filepath)
	if err != nil {
		return 0, fmt.Errorf("failed to parse coverage profile: %v", err)
	}

	// Calculate total coverage percentage
	coverage := calculateCoverage(profiles)
	return coverage, nil
}

// ParseCoverageFromStdin parses coverage profile data from stdin
func ParseCoverageFromStdin() (float64, error) {
	// Create a temporary file to store the stdin data
	tempFile, err := os.CreateTemp("", "coverage-*.out")
	if err != nil {
		return 0, fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(tempFile)

	// Read from stdin and write to the temporary file
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return 0, fmt.Errorf("error reading from stdin: %v", err)
		}

		if line != "" {
			if _, err := writer.WriteString(line); err != nil {
				return 0, fmt.Errorf("error writing to temporary file: %v", err)
			}
		}

		if err == io.EOF {
			break
		}
	}

	if err := writer.Flush(); err != nil {
		return 0, fmt.Errorf("error flushing data to temporary file: %v", err)
	}

	// Make sure we're at the beginning of the file for reading
	if _, err := tempFile.Seek(0, 0); err != nil {
		return 0, fmt.Errorf("error seeking in temporary file: %v", err)
	}

	// Parse the coverage profile
	profiles, err := cover.ParseProfiles(tempFile.Name())
	if err != nil {
		// Check if the input might be the output of go tool cover -func
		if isToolCoverOutput(tempFile.Name()) {
			return 0, fmt.Errorf("input appears to be the output of 'go tool cover -func'. Please provide a coverage profile file instead")
		}
		return 0, fmt.Errorf("failed to parse coverage profile from stdin: %v", err)
	}

	// Calculate total coverage percentage
	coverage := calculateCoverage(profiles)
	return coverage, nil
}

// isToolCoverOutput checks if the file contains the output of go tool cover -func
// instead of an actual coverage profile
func isToolCoverOutput(filepath string) bool {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return false
	}

	content := string(data)
	return len(content) > 0 && !strings.HasPrefix(content, "mode:") && 
		(strings.Contains(content, "(statements)") || strings.Contains(content, "total:"))
}

// calculateCoverage computes the coverage percentage from profile data
func calculateCoverage(profiles []*cover.Profile) float64 {
	var totalStmts, coveredStmts int

	for _, profile := range profiles {
		for _, block := range profile.Blocks {
			totalStmts += block.NumStmt
			if block.Count > 0 {
				coveredStmts += block.NumStmt
			}
		}
	}

	if totalStmts == 0 {
		return 0.0
	}

	return float64(coveredStmts) * 100.0 / float64(totalStmts)
}
