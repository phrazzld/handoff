# Handoff

A command-line utility to copy file contents to the clipboard in a formatted way. Perfect for quickly sharing code snippets, documentation, or any text files with others.

## Purpose

Handoff simplifies the process of sharing file contents by:
- Reading specified files and directories
- Formatting the content with filenames
- Automatically copying everything to your clipboard
- Respecting `.gitignore` rules when processing directories
- Providing useful statistics about the copied content

## Installation

### Prerequisites
- Go 1.16 or higher

### Building from Source
```bash
# Clone the repository
git clone https://github.com/phrazzld/handoff.git
cd handoff

# Build the binary
go build

# Optionally, install it to your GOPATH
go install
```

## Usage

### Basic Usage

```bash
# Copy a single file to clipboard
./handoff path/to/file.txt

# Copy multiple files to clipboard
./handoff file1.go file2.go

# Copy all files from a directory (respecting .gitignore)
./handoff path/to/directory
```

### Advanced Options

Handoff supports several command-line options for advanced usage:

```bash
# Show verbose output while processing files
./handoff --verbose path/to/directory

# Preview what would be copied without actually copying to clipboard
./handoff --dry-run path/to/directory

# Only include specific file extensions
./handoff --include=.go,.md path/to/directory

# Exclude specific file extensions
./handoff --exclude=.exe,.bin,.o path/to/directory

# Exclude specific files by name (regardless of directory)
./handoff --exclude-names=package-lock.json,yarn.lock path/to/directory

# Use a custom format for the output
./handoff --format="File: {path}\n{content}\n---\n" path/to/directory
```

### Command-line Flags

| Flag | Description |
|------|-------------|
| `--verbose` | Enable detailed output showing which files are processed |
| `--dry-run` | Preview what would be copied without actually copying to clipboard |
| `--include=.ext1,.ext2` | Only include files with specified extensions |
| `--exclude=.ext1,.ext2` | Exclude files with specified extensions |
| `--exclude-names=file1,file2` | Exclude specific files by name (e.g., package-lock.json,yarn.lock) |
| `--format="..."` | Customize the output format using {path} and {content} placeholders |

### Output Format

The copied content will be formatted as:

````
<context>
<path/to/file.txt>
```
content of file.txt
```
</path/to/file.txt>

<path/to/another/file.go>
```
content of file.go
```
</path/to/another/file.go>
</context>
````

You can customize this format using the `--format` flag with `{path}` and `{content}` placeholders.

### Output Statistics

After processing files, Handoff displays useful statistics about the copied content:

```
Handoff complete:
- Files: 3
- Lines: 45
- Characters: 1024
- Estimated tokens: 256
```

These statistics are particularly helpful for:
- Understanding how much content you're sharing
- Estimating LLM token usage when pasting into AI tools
- Monitoring the size of your clipboard content

## Git Integration

When processing directories, Handoff respects `.gitignore` rules:
- In Git repositories, files ignored by Git will not be included
- In non-Git directories, hidden files (starting with `.`) will be skipped
- This ensures that binary files, build artifacts, and other irrelevant files are not copied

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.