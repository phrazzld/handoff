package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	handoff "github.com/phrazzld/handoff/lib"
)

// parseConfig defines and parses command-line flags, processes include/exclude extensions,
// and returns a populated Config struct, output file path, force flag, and dry run flag.
func parseConfig() (*handoff.Config, string, bool, bool) {
	config := handoff.NewConfig()
	
	// Define flags
	var dryRun bool
	var outputFile string
	var force bool
	flag.BoolVar(&config.Verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&dryRun, "dry-run", false, "Preview what would be copied without actually copying")
	flag.StringVar(&config.Include, "include", "", "Comma-separated list of file extensions to include (e.g., .txt,.go)")
	flag.StringVar(&config.Exclude, "exclude", "", "Comma-separated list of file extensions to exclude (e.g., .exe,.bin)")
	flag.StringVar(&config.ExcludeNamesStr, "exclude-names", "", "Comma-separated list of file names to exclude (e.g., package-lock.json,yarn.lock)")
	flag.StringVar(&config.Format, "format", "<{path}>\n```\n{content}\n```\n</{path}>\n\n", "Custom format for output. Use {path} and {content} as placeholders")
	flag.StringVar(&outputFile, "output", "", "Write output to the specified file instead of clipboard (e.g., HANDOFF.md)")
	flag.BoolVar(&force, "force", false, "Allow overwriting existing files when using -output flag")

	// Parse command-line flags
	flag.Parse()

	// Process config (converts include/exclude strings to slices)
	config.ProcessConfig()
	
	return config, outputFile, force, dryRun
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

// resolveOutputPath converts a relative path to an absolute path.
// It returns the absolute path and any error encountered.
func resolveOutputPath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("output path is empty")
	}
	
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to determine absolute path: %w", err)
	}
	
	return absPath, nil
}

func main() {
	// Parse command-line flags and get configuration
	config, outputFile, force, dryRun := parseConfig()
	logger := handoff.NewLogger(config.Verbose)
	
	// Resolve output path if specified
	var absOutputPath string
	if outputFile != "" {
		var err error
		absOutputPath, err = resolveOutputPath(outputFile)
		if err != nil {
			logger.Error("Invalid output path: %v", err)
			os.Exit(1)
		}
		logger.Verbose("Output will be written to: %s", absOutputPath)
	}

	// Check if we have any paths to process
	if flag.NArg() < 1 {
		logger.Error("usage: %s [options] path1 [path2 ...]", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Process paths and get content
	formattedContent, err := handoff.ProcessProject(flag.Args(), config)
	if err != nil {
		logger.Error("Failed to process project: %v", err)
		os.Exit(1)
	}
	
	// Handle dry-run or copy to clipboard
	if dryRun {
		fmt.Println("### DRY RUN: Content that would be copied to clipboard ###")
		fmt.Println(formattedContent)
	} else {
		// Copy to clipboard
		if err := copyToClipboard(formattedContent); err != nil {
			logger.Error("Failed to copy to clipboard: %v", err)
			os.Exit(1)
		}
	}
	
	// Calculate and log statistics
	charCount, lineCount, tokenCount := handoff.CalculateStatistics(formattedContent)
	// Count processed files from the content
	processedFiles := strings.Count(formattedContent, "</")
	
	logger.Info("Handoff complete:")
	logger.Info("- Files: %d", processedFiles)
	logger.Info("- Lines: %d", lineCount)
	logger.Info("- Characters: %d", charCount)
	logger.Info("- Estimated tokens: %d", tokenCount)
}