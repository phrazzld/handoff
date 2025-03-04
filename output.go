package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

// Logger provides a simple logging interface with different log levels
type Logger struct {
	verbose bool
}

// newLogger creates a new Logger instance
func newLogger(verbose bool) *Logger {
	return &Logger{
		verbose: verbose,
	}
}

// Info logs an informational message to stderr
func (l *Logger) Info(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

// Warn logs a warning message to stderr
func (l *Logger) Warn(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "warning: "+format+"\n", args...)
}

// Error logs an error message to stderr
func (l *Logger) Error(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
}

// Verbose logs a message to stderr only if verbose mode is enabled
func (l *Logger) Verbose(format string, args ...interface{}) {
	if l.verbose {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}

// copyToClipboard copies text to the system clipboard with enhanced error reporting.
func copyToClipboard(text string) error {
	var errors []string
	
	// Try pbcopy (macOS)
	if _, err := exec.LookPath("pbcopy"); err == nil {
		cmd := exec.Command("pbcopy")
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err == nil {
			return nil // Success
		} else {
			errors = append(errors, fmt.Sprintf("pbcopy failed: %v", err))
		}
	} else {
		errors = append(errors, "pbcopy not found")
	}
	
	// Try xclip (X11/Linux)
	if _, err := exec.LookPath("xclip"); err == nil {
		cmd := exec.Command("xclip", "-selection", "clipboard")
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err == nil {
			return nil // Success
		} else {
			errors = append(errors, fmt.Sprintf("xclip failed: %v", err))
		}
	} else {
		errors = append(errors, "xclip not found")
	}
	
	// Try wl-copy (Wayland/Linux)
	if _, err := exec.LookPath("wl-copy"); err == nil {
		cmd := exec.Command("wl-copy")
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err == nil {
			return nil // Success
		} else {
			errors = append(errors, fmt.Sprintf("wl-copy failed: %v", err))
		}
	} else {
		errors = append(errors, "wl-copy not found")
	}
	
	// If we get here, all clipboard commands failed
	return fmt.Errorf("clipboard commands failed: %s", strings.Join(errors, "; "))
}

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
func logStatistics(content string, fileCount int, totalFiles int, config Config, logger *Logger) {
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
		if config.DryRun {
			logger.Verbose("Processed %d/%d files", fileCount, totalFiles)
		} else {
			logger.Verbose("Successfully copied content of %d/%d files to clipboard", fileCount, totalFiles)
		}
	}
}