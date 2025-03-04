# Handoff

A command-line utility to copy file contents to the clipboard in a formatted way. Perfect for quickly sharing code snippets, documentation, or any text files with others.

## Purpose

Handoff simplifies the process of sharing file contents by:
- Reading specified files and directories
- Formatting the content with filenames
- Automatically copying everything to your clipboard
- Respecting `.gitignore` rules when processing directories

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

```bash
# Copy a single file to clipboard
./handoff path/to/file.txt

# Copy multiple files to clipboard
./handoff file1.go file2.go

# Copy all files from a directory (respecting .gitignore)
./handoff path/to/directory
```

### Output Format

The copied content will be formatted as:

```
/path/to/file.txt
```
content of file.txt
```

/path/to/another/file.go
```
content of file.go
```
```

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