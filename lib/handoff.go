// Package handoff provides functionality for collecting and formatting file contents.
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

// Config holds application configuration settings
type Config struct {
	Verbose        bool
	Include        string
	Exclude        string
	ExcludeNamesStr string
	Format         string
	IncludeExts    []string
	ExcludeExts    []string
	ExcludeNames   []string
}

// NewConfig creates a new Config with default values.
func NewConfig() *Config {
	return &Config{
		Verbose: false,
		Format:  "<{path}>\n```\n{content}\n```\n</{path}>\n\n",
	}
}

// ProcessorFunc is a function type that processes a file's content and returns formatted output
type ProcessorFunc func(filePath string, content []byte) string

// Logger provides a simple logging interface with different log levels
type Logger struct {
	verbose bool
}

// NewLogger creates a new Logger instance
func NewLogger(verbose bool) *Logger {
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

// GitAvailable indicates whether the git command is available on the system.
var GitAvailable bool

func init() {
	_, err := exec.LookPath("git")
	GitAvailable = err == nil
}

// ProcessConfig processes the Include/Exclude strings in the Config and populates the extension slices
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

// IsGitIgnored checks if a file is gitignored or hidden.
func IsGitIgnored(file string) bool {
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

// GetGitFiles retrieves files from a directory using Git's ls-files command
func GetGitFiles(dir string) ([]string, error) {
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

// GetFilesWithFilepathWalk retrieves files from a directory by walking the filesystem
func GetFilesWithFilepathWalk(dir string) ([]string, error) {
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

// GetFilesFromDir retrieves all files to process from a directory.
// It tries to use Git first and falls back to filepath.Walk if Git is not available
// or the directory is not a Git repository.
func GetFilesFromDir(dir string) ([]string, error) {
	if GitAvailable {
		files, err := GetGitFiles(dir)
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
	return GetFilesWithFilepathWalk(dir)
}

// Constants for binary file detection
const (
	binarySampleSize           = 512  // Number of bytes to sample for binary detection
	binaryNonPrintableThreshold = 0.3 // Threshold ratio of non-printable chars to consider a file binary
)

// IsBinaryFile checks if a file is likely to be binary based on its content.
func IsBinaryFile(content []byte) bool {
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

// isWhitespace checks if a byte is a whitespace character
func isWhitespace(b byte) bool {
	return b == '\n' || b == '\r' || b == '\t' || b == ' '
}

// minInt returns the minimum of two integers
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ShouldProcess decides if a file should be processed based on all filters
func ShouldProcess(file string, config *Config) bool {
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

// ProcessFile processes a single file with the given processor and config
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
	if IsGitIgnored(filePath) {
		logger.Verbose("skipping gitignored file: %s", filePath)
		return ""
	}
	
	// Check if file should be processed based on filters
	if !ShouldProcess(filePath, config) {
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
	if IsBinaryFile(content) {
		logger.Verbose("skipping binary file: %s", filePath)
		return ""
	}
	
	// Process the content
	return processor(filePath, content)
}

// ProcessDirectory processes all files in a directory with the given processor and config
func ProcessDirectory(dirPath string, contentBuilder *strings.Builder, config *Config, logger *Logger, processor ProcessorFunc) {
	files, err := GetFilesFromDir(dirPath)
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

// ProcessPathWithProcessor processes a single path (file or directory) with a custom processor function and config
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

// ProcessPaths processes multiple paths and returns the number of processed files and total files
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
			if files, err := GetFilesFromDir(path); err == nil {
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

// WrapInContext wraps the content in a top-level context tag
func WrapInContext(content string) string {
	return "<context>\n" + content + "</context>"
}

// EstimateTokenCount counts tokens by tracking transitions between whitespace and non-whitespace characters
func EstimateTokenCount(text string) int {
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

// CalculateStatistics calculates statistics about the content
func CalculateStatistics(content string) (charCount, lineCount, tokenCount int) {
	charCount = len(content)
	lineCount = strings.Count(content, "\n") + 1
	tokenCount = EstimateTokenCount(content)
	return charCount, lineCount, tokenCount
}

// ProcessProject collects all files from the given paths according to filters,
// formats them, and returns the formatted content.
// This is the main function to use when integrating with other applications.
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

// WriteToFile writes the content to a file
func WriteToFile(content, filePath string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}