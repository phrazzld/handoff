package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// ParseCoverageFromFile runs go tool cover on the given profile file and parses the output
func ParseCoverageFromFile(filepath string) (float64, error) {
	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return 0, fmt.Errorf("coverage file %s does not exist", filepath)
	}

	// Run go tool cover to get the percentage
	cmd := exec.Command("go", "tool", "cover", "-func", filepath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("failed to run go tool cover: %v - output: %s", err, output)
	}

	// Parse the output
	coverage, err := ParseCoverageOutput(string(output))
	if err != nil {
		return 0, err
	}

	return coverage, nil
}

// ParseCoverageFromStdin parses coverage data from stdin
func ParseCoverageFromStdin() (float64, error) {
	reader := bufio.NewReader(os.Stdin)
	var output strings.Builder

	// Read all input
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return 0, fmt.Errorf("error reading from stdin: %v", err)
		}

		output.WriteString(line)

		if err == io.EOF {
			break
		}
	}

	// Parse the coverage from the accumulated output
	coverage, err := ParseCoverageOutput(output.String())
	if err != nil {
		return 0, err
	}

	return coverage, nil
}

// ParseCoverageOutput extracts the coverage percentage from go tool cover output
func ParseCoverageOutput(output string) (float64, error) {
	// Regular expression to match the total coverage line
	// Example line: "total:	(statements)	85.7%"
	re := regexp.MustCompile(`total:\s+\(statements\)\s+(\d+\.\d+)%`)
	matches := re.FindStringSubmatch(output)

	if len(matches) < 2 {
		return 0, fmt.Errorf("could not find coverage percentage in output: %s", output)
	}

	// Parse the percentage
	coverageStr := matches[1]
	coverage, err := strconv.ParseFloat(coverageStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse coverage percentage '%s': %v", coverageStr, err)
	}

	return coverage, nil
}