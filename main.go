package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	handoff "github.com/phrazzld/handoff/lib"
)

// ErrClipboardFailed is returned when all clipboard commands fail
var ErrClipboardFailed = errors.New("clipboard commands failed")

// parseConfig defines and parses command-line flags, processes include/exclude extensions,
// and returns a populated Config struct from the library package.
// It also returns the CLI-specific options as separate values (output file path, force flag, and dry run flag).
func parseConfig() (*handoff.Config, string, bool, bool) {
	// Define flags for CLI use
	var (
		verbose       bool
		include       string
		exclude       string
		excludeNames  string
		format        string = "<{path}>\n```\n{content}\n```\n</{path}>\n\n"
		dryRun        bool
		outputFile    string
		force         bool
	)

	// Define flag bindings
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&dryRun, "dry-run", false, "Preview what would be copied without actually copying")
	flag.StringVar(&include, "include", "", "Comma-separated list of file extensions to include (e.g., .txt,.go)")
	flag.StringVar(&exclude, "exclude", "", "Comma-separated list of file extensions to exclude (e.g., .exe,.bin)")
	flag.StringVar(&excludeNames, "exclude-names", "", "Comma-separated list of file names to exclude (e.g., package-lock.json,yarn.lock)")
	flag.StringVar(&format, "format", format, "Custom format for output. Use {path} and {content} as placeholders")
	flag.StringVar(&outputFile, "output", "", "Write output to the specified file instead of clipboard (e.g., HANDOFF.md)")
	flag.BoolVar(&force, "force", false, "Allow overwriting existing files when using -output flag")

	// Parse command-line flags
	flag.Parse()

	// Create config with functional options based on CLI flags
	var options []handoff.Option
	
	if verbose {
		options = append(options, handoff.WithVerbose(verbose))
	}
	
	if include != "" {
		options = append(options, handoff.WithInclude(include))
	}
	
	if exclude != "" {
		options = append(options, handoff.WithExclude(exclude))
	}
	
	if excludeNames != "" {
		options = append(options, handoff.WithExcludeNames(excludeNames))
	}
	
	if format != "" {
		options = append(options, handoff.WithFormat(format))
	}
	
	config := handoff.NewConfig(options...)

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
	return fmt.Errorf("%w: %s", ErrClipboardFailed, strings.Join(errors, "; "))
}

// resolveOutputPath converts a relative path to an absolute path.
// It returns the absolute path and any error encountered.
func resolveOutputPath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("output path cannot be empty")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to determine absolute path for %q: %w", path, err)
	}

	return absPath, nil
}

// checkFileExists checks if a file exists at the specified path.
// It returns true if the file exists, false otherwise.
// An error is returned if there's a problem checking file existence (e.g., permission issues).
func checkFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		// File exists
		return true, nil
	}
	if os.IsNotExist(err) {
		// File does not exist
		return false, nil
	}
	// Some other error occurred (e.g., permission denied)
	return false, fmt.Errorf("cannot check if file %q exists: %w", path, err)
}

// Note: processPathUsingLib function was removed as it was unused after refactoring

// logStatisticsUsingLib logs statistics about the processed content
// using the Stats struct returned by ProcessProject
func logStatisticsUsingLib(stats handoff.Stats, config *handoff.Config, logger *handoff.Logger) {
	// Log statistics
	logger.Info("Handoff complete:")
	logger.Info("- Files: %d/%d", stats.FilesProcessed, stats.FilesTotal)
	logger.Info("- Lines: %d", stats.Lines)
	logger.Info("- Characters: %d", stats.Chars)
	logger.Info("- Estimated tokens: %d", stats.Tokens)

	if config.Verbose {
		logger.Verbose("Processed files successfully")
	}
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

		// Check if the file exists and handle according to force flag
		exists, err := checkFileExists(absOutputPath)
		if err != nil {
			logger.Error("Error checking output file: %v", err)
			os.Exit(1)
		}

		if exists && !force {
			logger.Error("Output file %s already exists. Use -force flag to overwrite.", absOutputPath)
			os.Exit(1)
		} else if exists && force {
			logger.Verbose("Output file %s exists, will be overwritten because -force flag is set", absOutputPath)
		}
	}

	// Check if we have any paths to process
	if flag.NArg() < 1 {
		logger.Error("usage: %s [options] path1 [path2 ...]", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Process paths and get content
	formattedContent, stats, err := handoff.ProcessProject(flag.Args(), config)
	if err != nil {
		logger.Error("Failed to process project: %v", err)
		os.Exit(1)
	}

	// Handle output based on precedence: dry-run > output file > clipboard
	if dryRun {
		// Highest precedence: dry-run mode
		fmt.Println("### DRY RUN: Content that would be generated ###")
		fmt.Println(formattedContent)
		logger.Info("Dry run complete. No file written or clipboard modified.")
	} else if outputFile != "" {
		// Medium precedence: write to file
		logger.Verbose("Writing content (%d bytes) to file: %s", len(formattedContent), absOutputPath)
		if err := handoff.WriteToFile(formattedContent, absOutputPath, force); err != nil {
			logger.Error("Failed to write to file %s: %v", absOutputPath, err)
			os.Exit(1)
		}
		logger.Info("Output successfully written to %s", absOutputPath)
	} else {
		// Lowest precedence: copy to clipboard (default behavior)
		if err := copyToClipboard(formattedContent); err != nil {
			logger.Error("Failed to copy to clipboard: %v", err)
			os.Exit(1)
		}
		logger.Info("Content successfully copied to clipboard.")
	}

	// Log statistics
	logStatisticsUsingLib(stats, config, logger)
}
