package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	
	handoff "github.com/phrazzld/handoff/lib"
)

// Using the exported GitAvailable from handoff library
var gitAvailable = handoff.GitAvailable

// isGitIgnored checks if a file is gitignored or hidden.
func isGitIgnored(file string) bool {
	if !gitAvailable {
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

// getGitFiles retrieves files from a directory using Git's ls-files command
func getGitFiles(dir string) ([]string, error) {
	if !gitAvailable {
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

// getFilesWithFilepathWalk retrieves files from a directory by walking the filesystem
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
// or the directory is not a Git repository.
func getFilesFromDir(dir string) ([]string, error) {
	if gitAvailable {
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
	binarySampleSize           = 512  // Number of bytes to sample for binary detection
	binaryNonPrintableThreshold = 0.3 // Threshold ratio of non-printable chars to consider a file binary
)

// isBinaryFile checks if a file is likely to be binary based on its content.
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

// shouldProcess decides if a file should be processed based on all filters
func shouldProcess(file string, config *handoff.Config) bool {
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

// processFile processes a single file with the given processor and config
func processFile(filePath string, logger *handoff.Logger, config *handoff.Config, processor handoff.ProcessorFunc) string {
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

// processDirectory processes all files in a directory with the given processor and config
func processDirectory(dirPath string, contentBuilder *strings.Builder, config *handoff.Config, logger *handoff.Logger, processor handoff.ProcessorFunc) {
	files, err := getFilesFromDir(dirPath)
	if err != nil {
		logger.Error("processing directory %s: %v", dirPath, err)
		return
	}
	
	for _, file := range files {
		output := processFile(file, logger, config, processor)
		if output != "" {
			contentBuilder.WriteString(output)
		}
	}
}

// processPathWithProcessor processes a single path (file or directory) with a custom processor function and config
func processPathWithProcessor(path string, contentBuilder *strings.Builder, config *handoff.Config, logger *handoff.Logger, processor handoff.ProcessorFunc) {
	info, err := os.Stat(path)
	if err != nil {
		// Just log the error and continue with other paths
		logger.Warn("%v", err)
		return
	}

	if info.IsDir() {
		processDirectory(path, contentBuilder, config, logger, processor)
	} else {
		output := processFile(path, logger, config, processor)
		if output != "" {
			contentBuilder.WriteString(output)
		}
	}
}

// processPath processes a single path (file or directory) with the default processor.
// This maintains backward compatibility with existing code.
func processPath(path string, builder *strings.Builder, config *handoff.Config, logger *handoff.Logger) {
	processor := func(file string, content []byte) string {
		return fmt.Sprintf("<%s>\n```\n%s\n```\n</%s>\n\n", file, string(content), file)
	}
	processPathWithProcessor(path, builder, config, logger, processor)
}

// This function has been moved to the library