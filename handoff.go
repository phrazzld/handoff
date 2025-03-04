package main

import (
    "bytes"
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
                    files = append(files, filepath.Join(dir, line))
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
    content, err := os.ReadFile(file)
    if err != nil {
        fmt.Fprintf(os.Stderr, "error reading file %s: %v\n", file, err)
        return ""
    }

    // Skip binary files
    if isBinaryFile(content) {
        fmt.Fprintf(os.Stderr, "skipping binary file: %s\n", file)
        return ""
    }

    return fmt.Sprintf("%s\n```\n%s\n```\n\n", file, string(content))
}

// processPath processes a single path (file or directory).
func processPath(path string, builder *strings.Builder) {
    info, err := os.Stat(path)
    if err != nil {
        fmt.Fprintf(os.Stderr, "path not found: %s\n", path)
        return
    }

    if info.IsDir() {
        files, err := getFilesFromDir(path)
        if err != nil {
            fmt.Fprintf(os.Stderr, "error processing directory %s: %v\n", path, err)
            return
        }
        for _, file := range files {
            builder.WriteString(processFile(file))
        }
    } else if !isGitIgnored(path) {
        builder.WriteString(processFile(path))
    } else {
        fmt.Fprintf(os.Stderr, "skipping gitignored file: %s\n", path)
    }
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

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintf(os.Stderr, "usage: %s path1 [path2 ...]\n", os.Args[0])
        os.Exit(1)
    }

    var builder strings.Builder
    for _, path := range os.Args[1:] {
        processPath(path, &builder)
    }

    text := builder.String()
    if err := copyToClipboard(text); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
