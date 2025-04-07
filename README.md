# Handoff

A tool to collect and format code from multiple files for sharing with AI assistants or other programs.

## Features

- Collects files from specified paths (files or directories)
- Formats file contents with customizable templates
- Copies aggregated content to clipboard or writes to file
- Can be used as a library in other Go programs
- Provides content statistics (files, lines, characters, tokens)
- Supports git-aware file collection
- Filters files by extension or name
- Detects and skips binary files
- Includes file overwrite protection

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

### Command Line

```bash
./handoff [options] [path1] [path2] ...
```

#### Options

- `-verbose`: Enable verbose output
- `-dry-run`: Preview what would be copied without actually copying
- `-output`: Write output to the specified file instead of clipboard (e.g., `HANDOFF.md`)
- `-force`: Allow overwriting existing files when using `-output` flag
- `-include`: Comma-separated list of file extensions to include (e.g., `.txt,.go`)
- `-exclude`: Comma-separated list of file extensions to exclude (e.g., `.exe,.bin`)
- `-exclude-names`: Comma-separated list of file names to exclude (e.g., `package-lock.json,yarn.lock`)
- `-format`: Custom format for output. Use `{path}` and `{content}` as placeholders

#### Examples

```bash
# Copy all files in the current directory
./handoff .

# Copy only .go files from the src directory
./handoff -include=.go src/

# Copy specific files
./handoff main.go utils.go config.go

# Use a custom format
./handoff -format="File: {path}\n```go\n{content}\n```\n\n" .

# Write output to a file instead of clipboard
./handoff -output=HANDOFF.md .

# Write output to a file, overwriting if it exists
./handoff -output=HANDOFF.md -force .

# Preview content that would be written to file
./handoff -output=HANDOFF.md -dry-run .
```

### As a Library

Handoff can also be used as a library in other Go programs:

```go
import "github.com/phrazzld/handoff/lib"

// Create configuration
config := handoff.NewConfig()
config.Verbose = true
config.Exclude = ".exe,.bin,.jpg,.png"
config.ProcessConfig()

// Process project files
content, err := handoff.ProcessProject([]string{"./my-project"}, config)
if err != nil {
    panic(err)
}

// Use the formatted content
fmt.Println(content)

// Or write to a file
handoff.WriteToFile(content, "output.txt")
```

Check the `examples` directory for more detailed examples.

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
- Monitoring the size of your clipboard or file content

### Output Mode Precedence

When multiple output options are specified, Handoff follows this precedence:
1. `-dry-run`: Highest priority - outputs to screen only, no clipboard/file modifications
2. `-output`: Medium priority - writes to the specified file
3. Clipboard: Default behavior - copies to clipboard when no other output option is specified

### File Overwrite Protection

When using the `-output` flag, Handoff includes built-in protection against accidental file overwrites:
- If the specified output file already exists, Handoff will refuse to overwrite it
- To allow overwriting an existing file, use the `-force` flag
- This ensures you don't accidentally lose content in existing files

## Git Integration

When processing directories, Handoff respects `.gitignore` rules:
- In Git repositories, files ignored by Git will not be included
- In non-Git directories, hidden files (starting with `.`) will be skipped
- This ensures that binary files, build artifacts, and other irrelevant files are not copied

## Examples

### Gemini Planner

The `examples/gemini_planner.go` demonstrates how to use Handoff to extract code from a project and send it to Gemini to generate a technical plan.

```bash
go run examples/gemini_planner.go --project ./my-project --prompt "Add authentication to the application" --output PLAN.md
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.