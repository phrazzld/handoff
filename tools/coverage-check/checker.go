package main

// CheckCoverageThreshold compares the actual coverage to the threshold
// Returns true if coverage meets or exceeds the threshold, false otherwise
func CheckCoverageThreshold(coverage, threshold float64) bool {
	return coverage >= threshold
}
