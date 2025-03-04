package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// ProcessorFunc is a function type that processes a file's content and returns formatted output
type ProcessorFunc func(filePath string, content []byte) string

// Config holds application configuration settings
type Config struct {
	Verbose     bool
	DryRun      bool
	Include     string
	Exclude     string
	Format      string
	IncludeExts []string
	ExcludeExts []string
}

// parseConfig defines and parses command-line flags, processes include/exclude extensions,
// and returns a populated Config struct.
func parseConfig() Config {
	var config Config
	flag.BoolVar(&config.Verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&config.DryRun, "dry-run", false, "Preview what would be copied without actually copying")
	flag.StringVar(&config.Include, "include", "", "Comma-separated list of file extensions to include (e.g., .txt,.go)")
	flag.StringVar(&config.Exclude, "exclude", "", "Comma-separated list of file extensions to exclude (e.g., .exe,.bin)")
	flag.StringVar(&config.Format, "format", "<{path}>\n```\n{content}\n```\n</{path}>\n\n", "Custom format for output. Use {path} and {content} as placeholders")

	// Parse command-line flags
	flag.Parse()

	// Process include/exclude extensions
	if config.Include != "" {
		config.IncludeExts = strings.Split(config.Include, ",")
		for i, ext := range config.IncludeExts {
			config.IncludeExts[i] = strings.TrimSpace(ext)
			if !strings.HasPrefix(config.IncludeExts[i], ".") {
				config.IncludeExts[i] = "." + config.IncludeExts[i]
			}
		}
	}
	if config.Exclude != "" {
		config.ExcludeExts = strings.Split(config.Exclude, ",")
		for i, ext := range config.ExcludeExts {
			config.ExcludeExts[i] = strings.TrimSpace(ext)
			if !strings.HasPrefix(config.ExcludeExts[i], ".") {
				config.ExcludeExts[i] = "." + config.ExcludeExts[i]
			}
		}
	}
	
	return config
}

func main() {
	// Parse command-line flags and get configuration
	config := parseConfig()
	logger := newLogger(config.Verbose)

	// Check if we have any paths to process
	if flag.NArg() < 1 {
		logger.Error("usage: %s [options] path1 [path2 ...]", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Process paths
	contentBuilder := &strings.Builder{}
	processedFiles, totalFiles := processPaths(flag.Args(), contentBuilder, config, logger)
	
	// Wrap content in context tag
	formattedContent := wrapInContext(contentBuilder.String())
	
	// Handle dry-run or copy to clipboard
	if config.DryRun {
		fmt.Println("### DRY RUN: Content that would be copied to clipboard ###")
		fmt.Println(formattedContent)
	} else {
		// Copy to clipboard
		if err := copyToClipboard(formattedContent); err != nil {
			logger.Error("Failed to copy to clipboard: %v", err)
			os.Exit(1)
		}
	}
	
	// Log statistics
	logStatistics(formattedContent, processedFiles, totalFiles, config, logger)
}