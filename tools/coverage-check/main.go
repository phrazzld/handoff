// Package main provides a tool to check Go test coverage against a threshold
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Parse command line flags
	thresholdPtr := flag.Float64("threshold", 85.0, "Minimum coverage percentage required")
	filePtr := flag.String("file", "", "Coverage profile file (default reads from stdin)")
	verbosePtr := flag.Bool("verbose", false, "Show detailed output")
	flag.Parse()

	var coverage float64
	var err error

	// Parse coverage from file or stdin
	if *filePtr != "" {
		coverage, err = ParseCoverageFromFile(*filePtr)
	} else {
		coverage, err = ParseCoverageFromStdin()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing coverage: %v\n", err)
		os.Exit(2)
	}

	// Check coverage against threshold
	passed := CheckCoverageThreshold(coverage, *thresholdPtr)

	// Output results
	if *verbosePtr {
		fmt.Printf("Coverage: %.2f%%\n", coverage)
		fmt.Printf("Threshold: %.2f%%\n", *thresholdPtr)
		fmt.Printf("Status: %s\n", getStatusText(passed))
	} else {
		if passed {
			fmt.Printf("Coverage %.2f%% meets threshold of %.2f%%\n", coverage, *thresholdPtr)
		} else {
			fmt.Printf("Coverage %.2f%% is below threshold of %.2f%%\n", coverage, *thresholdPtr)
		}
	}

	if !passed {
		os.Exit(1)
	}
}

func getStatusText(passed bool) string {
	if passed {
		return "PASS"
	}
	return "FAIL"
}
