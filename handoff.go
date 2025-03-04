package main

import (
    "bytes"
    "flag"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

// gitAvailable indicates whether the git command is available on the system.
var gitAvailable bool

func init() {
    _, err := exec.LookPath("git")
    gitAvailable = err == nil
}

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

// getFilesFromDir retrieves all files to process from a directory.
func getFilesFromDir(dir string) ([]string, error) {
    if gitAvailable {
        cmd := exec.Command("git", "-C", dir, "ls-files", "--cached", "--others", "--exclude-standard")
        output, err := cmd.Output()
        if err == nil {
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
        if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 128 {
            // Not a git repo, fall back to filepath.Walk
        } else if err != nil {
            return nil, fmt.Errorf("error running git ls-files: %v", err)
        }
    }

    // Fallback to walking the directory, excluding hidden files and dirs
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

// isBinaryFile checks if a file is likely to be binary based on its content.
func isBinaryFile(content []byte) bool {
    // Check for null bytes, which are common in binary files
    if len(content) > 0 && bytes.IndexByte(content, 0) != -1 {
        return true
    }

    // Check for a high percentage of non-printable, non-whitespace characters
    // which suggest binary content
    nonPrintable := 0
    sampleSize := minInt(len(content), 512) // Sample the first 512 bytes
    for i := 0; i < sampleSize; i++ {
        if content[i] < 32 && !isWhitespace(content[i]) {
            nonPrintable++
        }
    }

    // If more than 30% of the sampled bytes are non-printable, consider it binary
    return nonPrintable > sampleSize*3/10
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

// processFile reads a file and formats its contents.
func processFile(file string) string {
    // First check if file exists
    if _, statErr := os.Stat(file); statErr != nil {
        if os.IsNotExist(statErr) {
            // Skip without warning if the file simply doesn't exist
            return ""
        }
    }
    
    content, err := os.ReadFile(file)
    if err != nil {
        // Log warning for other errors
        fmt.Fprintf(os.Stderr, "warning: cannot read %s: %v\n", file, err)
        return ""
    }

    // Skip binary files
    if isBinaryFile(content) {
        fmt.Fprintf(os.Stderr, "skipping binary file: %s\n", file)
        return ""
    }

    return fmt.Sprintf("%s\n```\n%s\n```\n\n", file, string(content))
}

// ProcessorFunc is a function type that processes a file's content and returns formatted output
type ProcessorFunc func(filePath string, content []byte) string

// processPathWithProcessor processes a single path with a custom processor function
func processPathWithProcessor(path string, builder *strings.Builder, processor ProcessorFunc) {
    info, err := os.Stat(path)
    if err != nil {
        // Just log the error and continue with other paths
        fmt.Fprintf(os.Stderr, "warning: %v\n", err)
        return
    }

    if info.IsDir() {
        files, err := getFilesFromDir(path)
        if err != nil {
            fmt.Fprintf(os.Stderr, "error processing directory %s: %v\n", path, err)
            return
        }
        for _, file := range files {
            // First check if file exists
            if _, statErr := os.Stat(file); statErr != nil {
                if os.IsNotExist(statErr) {
                    // Skip without warning if the file simply doesn't exist
                    continue
                }
            }
            
            content, err := os.ReadFile(file)
            if err != nil {
                // Log warning for other errors
                fmt.Fprintf(os.Stderr, "warning: cannot read %s: %v\n", file, err)
                continue
            }
            if output := processor(file, content); output != "" {
                builder.WriteString(output)
            }
        }
    } else if !isGitIgnored(path) {
        // First check if file exists
        if _, statErr := os.Stat(path); statErr != nil {
            if os.IsNotExist(statErr) {
                // Skip without warning if the file simply doesn't exist
                return
            }
        }
        
        content, err := os.ReadFile(path)
        if err != nil {
            // Log warning for other errors
            fmt.Fprintf(os.Stderr, "warning: cannot read %s: %v\n", path, err)
            return
        }
        if output := processor(path, content); output != "" {
            builder.WriteString(output)
        }
    } else {
        fmt.Fprintf(os.Stderr, "skipping gitignored file: %s\n", path)
    }
}

// processPath processes a single path (file or directory) with the default processor.
// This maintains backward compatibility with existing code.
func processPath(path string, builder *strings.Builder) {
    processor := func(file string, content []byte) string {
        if isBinaryFile(content) {
            fmt.Fprintf(os.Stderr, "skipping binary file: %s\n", file)
            return ""
        }
        return fmt.Sprintf("<%s>\n```\n%s\n```\n</%s>\n\n", file, string(content), file)
    }
    processPathWithProcessor(path, builder, processor)
}

// copyToClipboard copies text to the system clipboard.
func copyToClipboard(text string) error {
    if _, err := exec.LookPath("pbcopy"); err == nil {
        cmd := exec.Command("pbcopy")
        cmd.Stdin = strings.NewReader(text)
        return cmd.Run()
    } else if _, err := exec.LookPath("xclip"); err == nil {
        cmd := exec.Command("xclip", "-selection", "clipboard")
        cmd.Stdin = strings.NewReader(text)
        return cmd.Run()
    } else if _, err := exec.LookPath("wl-copy"); err == nil {
        cmd := exec.Command("wl-copy")
        cmd.Stdin = strings.NewReader(text)
        return cmd.Run()
    }
    return fmt.Errorf("no supported clipboard command found (pbcopy, xclip, wl-copy)")
}

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

func main() {
    // Define command-line flags
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

    // Check if we have any paths to process
    if flag.NArg() < 1 {
        fmt.Fprintf(os.Stderr, "usage: %s [options] path1 [path2 ...]\n", os.Args[0])
        flag.PrintDefaults()
        os.Exit(1)
    }

    // Process paths
    var builder strings.Builder
    totalFiles := 0
    processedFiles := 0

    for _, path := range flag.Args() {
        if config.Verbose {
            fmt.Fprintf(os.Stderr, "Processing path: %s\n", path)
        }

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
            // Skip files based on extension filters
            ext := strings.ToLower(filepath.Ext(file))
            if len(config.IncludeExts) > 0 {
                included := false
                for _, includeExt := range config.IncludeExts {
                    if ext == includeExt {
                        included = true
                        break
                    }
                }
                if !included {
                    if config.Verbose {
                        fmt.Fprintf(os.Stderr, "Skipping file (not in include list): %s\n", file)
                    }
                    return ""
                }
            }
            if len(config.ExcludeExts) > 0 {
                for _, excludeExt := range config.ExcludeExts {
                    if ext == excludeExt {
                        if config.Verbose {
                            fmt.Fprintf(os.Stderr, "Skipping file (in exclude list): %s\n", file)
                        }
                        return ""
                    }
                }
            }

            // Skip binary files
            if isBinaryFile(fileContent) {
                if config.Verbose {
                    fmt.Fprintf(os.Stderr, "Skipping binary file: %s\n", file)
                }
                return ""
            }

            processedFiles++
            if config.Verbose {
                fmt.Fprintf(os.Stderr, "Processing file (%d/%d): %s\n", processedFiles, totalFiles, file)
            }

            // Format the output using the custom format
            output := config.Format
            output = strings.ReplaceAll(output, "{path}", file)
            output = strings.ReplaceAll(output, "{content}", string(fileContent))
            return output
        }

        // Process the path with our custom processor
        processPathWithProcessor(path, &builder, pathProcessor)
    }

    // Wrap everything in a top-level context tag
    text := "<context>\n" + builder.String() + "</context>"
    
    // In dry-run mode, just print what would be copied
    if config.DryRun {
        fmt.Println("### DRY RUN: Content that would be copied to clipboard ###")
        fmt.Println(text)
        if config.Verbose {
            fmt.Fprintf(os.Stderr, "Processed %d/%d files\n", processedFiles, totalFiles)
        }
        return
    }

    // Copy to clipboard
    if err := copyToClipboard(text); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    if config.Verbose {
        fmt.Fprintf(os.Stderr, "Successfully copied content of %d/%d files to clipboard\n", processedFiles, totalFiles)
    }
}
