// Simple example demonstrating basic usage of the handoff library
//
// This example shows the minimal steps required to use the handoff library:
// 1. Create and configure a Config object
// 2. Process files/directories 
// 3. Handle the generated content
// 4. Get statistics about the processed content
//
// Usage examples:
//
//     # Process current directory, output to console
//     go run simple_usage.go
//
//     # Process specific directory, write to file
//     go run simple_usage.go --dir ./src --output codebase.md
//
//     # Include only certain file types
//     go run simple_usage.go --dir ./src --include .go,.md
//
//     # Exclude specific file types
//     go run simple_usage.go --dir ./src --exclude .exe,.bin,.jpg
//
package main

import (
	"flag"
	"fmt"
	"os"

	lib "github.com/phrazzld/handoff/lib"
)

func main() {
	// Parse command-line options
	inputDir := flag.String("dir", ".", "Directory or file to process")
	outputFile := flag.String("output", "", "Output file (if empty, prints to console)")
	includeExts := flag.String("include", "", "File extensions to include (e.g., '.go,.txt')")
	excludeExts := flag.String("exclude", ".exe,.bin,.obj,.jpg,.png,.gif", "File extensions to exclude")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	flag.Parse()

	// Create a new configuration with default settings
	config := lib.NewConfig()
	
	// Set configuration options from command-line flags
	config.Verbose = *verbose
	config.Include = *includeExts
	config.Exclude = *excludeExts
	
	// Process the string-based settings into slices
	// This step is required before using the configuration
	config.ProcessConfig()

	// Process the specified directory/file
	// The first parameter is a slice of paths to process
	content, err := lib.ProcessProject([]string{*inputDir}, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing files: %v\n", err)
		os.Exit(1)
	}

	// Calculate statistics about the processed content
	chars, lines, tokens := lib.CalculateStatistics(content)
	fmt.Printf("\nContent statistics:\n")
	fmt.Printf("- Characters: %d\n", chars)
	fmt.Printf("- Lines: %d\n", lines)
	fmt.Printf("- Estimated tokens: %d\n\n", tokens)

	// Either write to file or print to console
	if *outputFile != "" {
		// Write content to specified file
		if err := lib.WriteToFile(content, *outputFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Content successfully written to %s\n", *outputFile)
	} else {
		// If no output file specified, print to console
		fmt.Println("------- GENERATED CONTENT -------")
		fmt.Println(content)
		fmt.Println("--------------------------------")
	}
}