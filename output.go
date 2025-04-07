// Package main implements a command-line utility for copying file contents in a formatted way.
package main

import (
	"strings"
	"unicode"
	
	handoff "github.com/phrazzld/handoff/lib"
)

// estimateTokenCount counts tokens by tracking transitions between whitespace and non-whitespace characters
func estimateTokenCount(text string) int {
	count := 0
	inToken := false
	for _, r := range text {
		if unicode.IsSpace(r) {
			if inToken {
				count++
				inToken = false
			}
		} else {
			inToken = true
		}
	}
	if inToken {
		count++ // Count the last token if text ends with non-whitespace
	}
	return count
}

// wrapInContext wraps the content in a top-level context tag
func wrapInContext(content string) string {
	return "<context>\n" + content + "</context>"
}

// logStatistics calculates and logs statistics about the copied content
func logStatistics(content string, fileCount int, totalFiles int, config *handoff.Config, logger *handoff.Logger) {
	charCount := len(content)
	lineCount := strings.Count(content, "\n") + 1
	tokenCount := estimateTokenCount(content)
	
	// Log statistics
	logger.Info("Handoff complete:")
	logger.Info("- Files: %d", fileCount)
	logger.Info("- Lines: %d", lineCount)
	logger.Info("- Characters: %d", charCount)
	logger.Info("- Estimated tokens: %d", tokenCount)
	
	if config.Verbose {
		logger.Verbose("Processed %d/%d files", fileCount, totalFiles)
	}
}