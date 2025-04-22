package main

import "testing"

func TestCheckCoverageThreshold(t *testing.T) {
	testCases := []struct {
		name            string
		coverage        float64
		threshold       float64
		expectedOutcome bool
	}{
		{
			name:            "Coverage exceeds threshold",
			coverage:        90.0,
			threshold:       85.0,
			expectedOutcome: true,
		},
		{
			name:            "Coverage equals threshold",
			coverage:        85.0,
			threshold:       85.0,
			expectedOutcome: true,
		},
		{
			name:            "Coverage below threshold",
			coverage:        80.0,
			threshold:       85.0,
			expectedOutcome: false,
		},
		{
			name:            "Zero coverage",
			coverage:        0.0,
			threshold:       85.0,
			expectedOutcome: false,
		},
		{
			name:            "Zero threshold",
			coverage:        85.0,
			threshold:       0.0,
			expectedOutcome: true,
		},
		{
			name:            "Both zero",
			coverage:        0.0,
			threshold:       0.0,
			expectedOutcome: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CheckCoverageThreshold(tc.coverage, tc.threshold)
			if result != tc.expectedOutcome {
				t.Errorf("Expected %v but got %v for coverage=%.2f%%, threshold=%.2f%%",
					tc.expectedOutcome, result, tc.coverage, tc.threshold)
			}
		})
	}
}