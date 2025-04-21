// Package handoff provides functionality for collecting and formatting file contents 
// from multiple files and directories for sharing with AI assistants or other applications.
// 
// The package supports file filtering by extension or name, respects Git's ignore rules,
// detects and skips binary files, and provides customizable output formatting.
//
// Basic usage:
//
//	config := handoff.NewConfig()
//	config.Include = ".go,.md"  // Only include Go and Markdown files
//	config.ProcessConfig()      // Process string-based config into slice-based filters
//	
//	content, err := handoff.ProcessProject([]string{"./src", "README.md"}, config)
//	if err != nil {
//	    // Handle error
//	}
//	
//	// Use the formatted content or write it to a file
//	handoff.WriteToFile(content, "output.md")
package handoff

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"unicode"
)

// Config holds all configuration options for file processing and output formatting.
// Users should create a Config with NewConfig() and set the desired options before
// calling ProcessConfig() to prepare the configuration for use.
type Config struct {
	// Verbose enables detailed logging output
	Verbose bool

	// Include is a comma-separated list of file extensions to include (e.g., ".go,.md")
	Include string

	// Exclude is a comma-separated list of file extensions to exclude (e.g., ".exe,.bin")
	Exclude string

	// ExcludeNamesStr is a comma-separated list of filenames to exclude (e.g., "package-lock.json,yarn.lock")
	ExcludeNamesStr string

	// Format is a template string for formatting output, using {path} and {content} placeholders
	Format string

	// IncludeExts contains the processed list of file extensions to include (populated by ProcessConfig)
	IncludeExts []string

	// ExcludeExts contains the processed list of file extensions to exclude (populated by ProcessConfig) 
	ExcludeExts []string

	// ExcludeNames contains the processed list of filenames to exclude (populated by ProcessConfig)
	ExcludeNames []string
}

// NewConfig creates a new Config with default values.
// By default, Verbose is false and Format uses a sensible default format
// with file path headers and code fences.
func NewConfig() *Config {
	return &Config{
		Verbose: false,
		Format:  "<{path}>\n```\n{content}\n```\n</{path}>\n\n",
	}
}

// ProcessorFunc is a function type that processes a file's content and returns formatted output.
// It receives the file path and raw content and should return the processed content as a string.
// This type is used for custom file processing in functions like ProcessFile and ProcessPathWithProcessor.
type ProcessorFunc func(filePath string, content []byte) string

// Logger provides a simple logging interface with different log levels.
// All messages are sent to stderr with appropriate prefixes for their level.
type Logger struct {
	// verbose determines whether Verbose-level messages are displayed
	verbose bool
}

// NewLogger creates a new Logger instance with the specified verbosity setting.
// When verbose is false, calls to the Verbose method will be suppressed.
func NewLogger(verbose bool) *Logger {
	return &Logger{
		verbose: verbose,
	}
}

// Info logs an informational message to stderr.
// These messages are always displayed regardless of the verbose setting.
func (l *Logger) Info(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

// Warn logs a warning message to stderr.
// Warning messages are prefixed with "warning: " and are always displayed.
func (l *Logger) Warn(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "warning: "+format+"\n", args...)
}

// Error logs an error message to stderr.
// Error messages are prefixed with "error: " and are always displayed.
func (l *Logger) Error(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
}

// Verbose logs a message to stderr only if verbose mode is enabled.
// These messages are useful for detailed progress information.
func (l *Logger) Verbose(format string, args ...interface{}) {
	if l.verbose {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}

// GitAvailable indicates whether the git command is available on the system.
// This variable is set during package initialization and can be used
// to determine if Git functionality (like respecting .gitignore rules) can be used.
var GitAvailable bool

func init() {
	_, err := exec.LookPath("git")
	GitAvailable = err == nil
}

// ProcessConfig processes the string-based Include, Exclude, and ExcludeNamesStr fields 
// in the Config struct and populates the corresponding slice fields (IncludeExts, ExcludeExts, ExcludeNames).
//
// This method should be called after setting the string fields and before using the Config
// with ProcessProject or other processing functions. It ensures file extensions start with a dot
// and normalizes strings by trimming spaces.
func (c *Config) ProcessConfig() {
	// Process include/exclude extensions
	if c.Include != "" {
		c.IncludeExts = strings.Split(c.Include, ",")
		for i, ext := range c.IncludeExts {
			c.IncludeExts[i] = strings.TrimSpace(ext)
			if !strings.HasPrefix(c.IncludeExts[i], ".") {
				c.IncludeExts[i] = "." + c.IncludeExts[i]
			}
		}
	}
	if c.Exclude != "" {
		c.ExcludeExts = strings.Split(c.Exclude, ",")
		for i, ext := range c.ExcludeExts {
			c.ExcludeExts[i] = strings.TrimSpace(ext)
			if !strings.HasPrefix(c.ExcludeExts[i], ".") {
				c.ExcludeExts[i] = "." + c.ExcludeExts[i]
			}
		}
	}
	// Process exclude names
	if c.ExcludeNamesStr != "" {
		c.ExcludeNames = strings.Split(c.ExcludeNamesStr, ",")
		for i, name := range c.ExcludeNames {
			c.ExcludeNames[i] = strings.TrimSpace(name)
		}
	}
}

// isGitIgnored checks if a file is gitignored or hidden (internal helper).
func isGitIgnored(file string) bool {
	if !GitAvailable {
		return strings.HasPrefix(filepath.Base(file), ".")
	}
	dir := filepath.Dir(file)
	filename := filepath.Base(file)
	cmd := exec.Command("git", "-C", dir, "check-ignore", "-q", filename)
	err := cmd.Run()
	if err == nil { // Exit code 0: file is ignored
		return true
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() == 1 { // Exit code 1: file is not ignored
			return false
		}
	}
	// Other errors (e.g., not a git repo): fall back to checking if hidden
	return strings.HasPrefix(filename, ".")
}

// getGitFiles retrieves files from a directory using Git's ls-files command (internal helper)
func getGitFiles(dir string) ([]string, error) {
	if !GitAvailable {
		return nil, fmt.Errorf("git not available")
	}

	cmd := exec.Command("git", "-C", dir, "ls-files", "--cached", "--others", "--exclude-standard")
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 128 {
			return nil, fmt.Errorf("not a git repository")
		}
		return nil, fmt.Errorf("error running git ls-files: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var files []string
	for _, line := range lines {
		if line != "" {
			filePath := filepath.Join(dir, line)
			// Check if file still exists before adding it
			if _, err := os.Stat(filePath); err == nil {
				files = append(files, filePath)
			}
		}
	}
	return files, nil
}

// getFilesWithFilepathWalk retrieves files from a directory by walking the filesystem (internal helper)
func getFilesWithFilepathWalk(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") && path != dir {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasPrefix(info.Name(), ".") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// getFilesFromDir retrieves all files to process from a directory.
// It tries to use Git first and falls back to filepath.Walk if Git is not available
// or the directory is not a Git repository. (internal helper)
func getFilesFromDir(dir string) ([]string, error) {
	if GitAvailable {
		files, err := getGitFiles(dir)
		if err == nil {
			return files, nil
		}
		// If there's an error running git ls-files and it's not "not a git repository", return the error
		if !strings.Contains(err.Error(), "not a git repository") {
			return nil, err
		}
		// Otherwise fall back to filepath.Walk
	}

	// Fallback to walking the directory, excluding hidden files and dirs
	return getFilesWithFilepathWalk(dir)
}

// Constants for binary file detection
const (
	binarySampleSize            = 512 // Number of bytes to sample for binary detection
	binaryNonPrintableThreshold = 0.3 // Threshold ratio of non-printable chars to consider a file binary
)

// isBinaryFile checks if a file is likely to be binary based on its content (internal helper).
func isBinaryFile(content []byte) bool {
	// Check for null bytes, which are common in binary files
	if len(content) > 0 && bytes.IndexByte(content, 0) != -1 {
		return true
	}

	// Check for a high percentage of non-printable, non-whitespace characters
	// which suggest binary content
	nonPrintable := 0
	sampleSize := minInt(len(content), binarySampleSize) // Sample the first bytes
	for i := 0; i < sampleSize; i++ {
		if content[i] < 32 && !isWhitespace(content[i]) {
			nonPrintable++
		}
	}

	// If more than the threshold percentage of sampled bytes are non-printable, consider it binary
	return float64(nonPrintable) > float64(sampleSize)*binaryNonPrintableThreshold
}

// isWhitespace checks if a byte is a whitespace character (unexported, internal helper)
func isWhitespace(b byte) bool {
	return b == '\n' || b == '\r' || b == '\t' || b == ' '
}

// minInt returns the minimum of two integers (unexported, internal helper)
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// shouldProcess decides if a file should be processed based on all filters (internal helper)
func shouldProcess(file string, config *Config) bool {
	base := filepath.Base(file)
	ext := strings.ToLower(filepath.Ext(file))

	// Check exclude names filter
	if len(config.ExcludeNames) > 0 && slices.Contains(config.ExcludeNames, base) {
		return false
	}

	// Check include extensions filter
	if len(config.IncludeExts) > 0 {
		included := false
		for _, includeExt := range config.IncludeExts {
			if ext == includeExt {
				included = true
				break
			}
		}
		if !included {
			return false
		}
	}

	// Check exclude extensions filter
	if len(config.ExcludeExts) > 0 {
		for _, excludeExt := range config.ExcludeExts {
			if ext == excludeExt {
				return false
			}
		}
	}

	return true
}

// ProcessFile processes a single file with the given processor function and configuration.
// It applies various filters (gitignore, extension/name filters, binary detection) and
// passes the valid file's content to the processor function to generate formatted output.
//
// If the file should be skipped (doesn't exist, is gitignored, doesn't match filters, 
// or is binary), an empty string is returned and appropriate messages are logged.
//
// Parameters:
//   - filePath: The path to the file to process
//   - logger: Logger for status and error messages
//   - config: Configuration options controlling filtering
//   - processor: Function to process the file content
//
// Returns a formatted string for valid files or an empty string for skipped files.
func ProcessFile(filePath string, logger *Logger, config *Config, processor ProcessorFunc) string {
	// First check if file exists
	if _, statErr := os.Stat(filePath); statErr != nil {
		if os.IsNotExist(statErr) {
			// Skip without warning if the file simply doesn't exist
			return ""
		}
		// Log warning for other errors
		logger.Warn("stat %s: %v", filePath, statErr)
		return ""
	}

	// Check if file is gitignored
	if isGitIgnored(filePath) {
		logger.Verbose("skipping gitignored file: %s", filePath)
		return ""
	}

	// Check if file should be processed based on filters
	if !shouldProcess(filePath, config) {
		if len(config.ExcludeNames) > 0 && slices.Contains(config.ExcludeNames, filepath.Base(filePath)) {
			logger.Verbose("skipping file (in exclude-names list): %s", filePath)
		}
		return ""
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		logger.Warn("cannot read %s: %v", filePath, err)
		return ""
	}

	// Skip binary files
	if isBinaryFile(content) {
		logger.Verbose("skipping binary file: %s", filePath)
		return ""
	}

	// Process the content
	return processor(filePath, content)
}

// ProcessDirectory processes all files in a directory with the given processor and config.
// It discovers files in the directory (respecting Git's ignore rules if available),
// filters them according to the provided configuration, and processes each valid file,
// appending the formatted output to the contentBuilder.
//
// Parameters:
//   - dirPath: The directory path to process
//   - contentBuilder: Builder to append formatted content to
//   - config: Configuration for file filtering and processing
//   - logger: Logger for status and error messages
//   - processor: Function to process each file's content
func ProcessDirectory(dirPath string, contentBuilder *strings.Builder, config *Config, logger *Logger, processor ProcessorFunc) {
	files, err := getFilesFromDir(dirPath)
	if err != nil {
		logger.Error("processing directory %s: %v", dirPath, err)
		return
	}

	for _, file := range files {
		output := ProcessFile(file, logger, config, processor)
		if output != "" {
			contentBuilder.WriteString(output)
		}
	}
}

// ProcessPathWithProcessor processes a single path (file or directory) with a custom processor function.
// This is a unified entry point that handles both files and directories, dispatching to the
// appropriate handler based on the path type. For directories, it processes all contained files
// recursively. For files, it processes the file directly.
//
// Parameters:
//   - path: Path to a file or directory
//   - contentBuilder: Builder to append formatted content to
//   - config: Configuration for file filtering and processing
//   - logger: Logger for status and error messages
//   - processor: Function to process each file's content
func ProcessPathWithProcessor(path string, contentBuilder *strings.Builder, config *Config, logger *Logger, processor ProcessorFunc) {
	info, err := os.Stat(path)
	if err != nil {
		// Just log the error and continue with other paths
		logger.Warn("%v", err)
		return
	}

	if info.IsDir() {
		ProcessDirectory(path, contentBuilder, config, logger, processor)
	} else {
		output := ProcessFile(path, logger, config, processor)
		if output != "" {
			contentBuilder.WriteString(output)
		}
	}
}

// ProcessPaths processes multiple file or directory paths according to the configuration.
// It creates a customized processor function that tracks progress and formats output
// using the config's Format template, then processes each path with this processor.
//
// Parameters:
//   - paths: List of file or directory paths to process
//   - config: Configuration for file filtering and processing
//   - logger: Logger for status and error messages
//
// Returns:
//   - A string containing the combined formatted content
//   - The number of files successfully processed
//   - The total number of candidate files found (before filtering)
func ProcessPaths(paths []string, config *Config, logger *Logger) (string, int, int) {
	contentBuilder := &strings.Builder{}
	processedFiles := 0
	totalFiles := 0

	for _, path := range paths {
		logger.Verbose("Processing path: %s", path)

		// Count total files before processing
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			totalFiles++
		} else if err == nil && info.IsDir() {
			if files, err := getFilesFromDir(path); err == nil {
				totalFiles += len(files)
			}
		}

		// Custom process function with config and progress tracking
		pathProcessor := func(file string, fileContent []byte) string {
			processedFiles++
			logger.Verbose("Processing file (%d/%d): %s", processedFiles, totalFiles, file)

			// Format the output using the custom format
			output := config.Format
			output = strings.ReplaceAll(output, "{path}", file)
			output = strings.ReplaceAll(output, "{content}", string(fileContent))
			return output
		}

		// Process the path with our custom processor
		ProcessPathWithProcessor(path, contentBuilder, config, logger, pathProcessor)
	}

	return contentBuilder.String(), processedFiles, totalFiles
}

// WrapInContext wraps the content in top-level context tags.
// This provides consistent formatting for the final output, making it easier
// to identify the boundaries of the collected content.
//
// Parameter:
//   - content: The raw content to wrap
//
// Returns:
//   - The content wrapped with <context> tags
func WrapInContext(content string) string {
	return "<context>\n" + content + "</context>"
}

// estimateTokenCount counts tokens by tracking transitions between whitespace and non-whitespace characters (internal helper)
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

// CalculateStatistics calculates useful statistics about the content.
// This function analyzes the provided content and returns counts of characters,
// lines, and tokens (words or code-like tokens) it contains.
//
// Parameter:
//   - content: The content to analyze
//
// Returns:
//   - charCount: Total number of characters in the content
//   - lineCount: Total number of lines in the content (based on newlines)
//   - tokenCount: Estimated number of tokens/words in the content
func CalculateStatistics(content string) (charCount, lineCount, tokenCount int) {
	charCount = len(content)
	lineCount = strings.Count(content, "\n") + 1
	tokenCount = estimateTokenCount(content)
	return charCount, lineCount, tokenCount
}

// ProcessProject collects and formats content from files in the specified paths.
// This is the main entry point for the library and the primary function that external
// applications should use. It handles all aspects of file collection, filtering, and
// formatting according to the provided configuration.
//
// If config is nil, default configuration is used. The function automatically calls
// ProcessConfig() on the configuration to ensure string-based filters are processed.
//
// Parameters:
//   - paths: List of file or directory paths to process
//   - config: Configuration for file filtering and processing (can be nil for defaults)
//
// Returns:
//   - The formatted content wrapped in context tags
//   - An error if no paths are provided or if processing fails
func ProcessProject(paths []string, config *Config) (string, error) {
	if config == nil {
		config = NewConfig()
	}

	config.ProcessConfig()
	logger := NewLogger(config.Verbose)

	if len(paths) == 0 {
		return "", fmt.Errorf("no paths provided")
	}

	// Process paths
	content, processedFiles, totalFiles := ProcessPaths(paths, config, logger)

	// Wrap content in context tag
	formattedContent := WrapInContext(content)

	// Log statistics
	if config.Verbose {
		charCount, lineCount, tokenCount := CalculateStatistics(formattedContent)
		logger.Info("Handoff complete:")
		logger.Info("- Files: %d", processedFiles)
		logger.Info("- Lines: %d", lineCount)
		logger.Info("- Characters: %d", charCount)
		logger.Info("- Estimated tokens: %d", tokenCount)
		logger.Verbose("Processed %d/%d files", processedFiles, totalFiles)
	}

	return formattedContent, nil
}

// WriteToFile writes the content to a file at the specified path.
// This is a convenience function for saving the output of ProcessProject
// directly to a file. The file is created with 0644 permissions.
//
// Parameters:
//   - content: The content to write to the file
//   - filePath: The path where the file should be created
//
// Returns:
//   - An error if the file cannot be written (e.g., due to permissions or path issues)
func WriteToFile(content, filePath string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}
