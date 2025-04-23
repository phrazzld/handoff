// Package handoff provides functionality for collecting and formatting file contents
// from projects for sharing with AI assistants or other applications.
//
// The library simplifies the process of gathering code and documentation from
// multiple files and directories, filtering them by various criteria, and formatting
// the content in a standardized way, making it ideal for creating comprehensive
// context for AI assistants like Claude, ChatGPT, or Gemini.
//
// # Features
//
// The handoff library offers several key features:
//
//   - Processes multiple files and directories in a single operation
//   - Respects Git's ignore rules when Git is available
//   - Detects and skips binary files automatically
//   - Provides flexible file filtering by extension and name
//   - Supports custom output formatting with templates
//   - Calculates and reports statistics about processed content
//
// # Basic Usage
//
// The main API entry point is the ProcessProject function, which handles everything
// from file discovery to filtering and formatting:
//
//	// Create configuration with functional options
//	config := handoff.NewConfig(
//	    handoff.WithInclude(".go,.md"),         // Only include Go and Markdown files
//	    handoff.WithExcludeNames("go.sum"),     // Skip go.sum files
//	)
//
//	// Process files and get formatted content
//	content, stats, err := handoff.ProcessProject([]string{"./src", "README.md"}, config)
//	if err != nil {
//	    log.Fatalf("Error processing project: %v", err)
//	}
//
//	// Use the content with an AI assistant or save to a file
//	// The third parameter (true) allows overwriting existing files
//	if err := handoff.WriteToFile(content, "output.md", true); err != nil {
//	    log.Fatalf("Error writing to file: %v", err)
//	}
//
// # Customization Options
//
// The library offers several ways to customize the behavior through configuration options:
//
//	// Create default configuration with custom options
//	config := handoff.NewConfig(
//	    handoff.WithInclude(".go,.md,.txt"),       // Only include specific file types
//	    handoff.WithExclude(".min.js,.exe,.bin"),  // Exclude minified JS and binaries
//	    handoff.WithExcludeNames("go.sum,node_modules,package-lock.json"), // Skip specific files
//	    handoff.WithVerbose(true),                  // Enable verbose logging
//	    handoff.WithFormat("## {path}\n```\n{content}\n```\n\n"), // Custom output format
//	)
//
//	// Process files and get formatted content
//	content, stats, err := handoff.ProcessProject([]string{"./src", "README.md"}, config)
//	if err != nil {
//	    log.Fatalf("Error processing project: %v", err)
//	}
//
//	// Use stats to report information about the processed content
//	fmt.Printf("Processed %d/%d files\n", stats.FilesProcessed, stats.FilesTotal)
//	fmt.Printf("Total content: %d lines, %d chars, ~%d tokens\n", 
//	    stats.Lines, stats.Chars, stats.Tokens)
//
// # Additional API Functions
//
// Besides the main ProcessProject function, the package exports a few utilities:
//
//	// Calculate statistics for any arbitrary content
//	chars, lines, tokens := handoff.CalculateStatistics("Some content to analyze")
//	fmt.Printf("Analysis: %d chars, %d lines, ~%d tokens\n", chars, lines, tokens)
//
//	// Wrap content in context tags for better AI assistant compatibility
//	wrappedContent := handoff.WrapInContext("Content to wrap in context tags")
//	fmt.Println(wrappedContent) // Outputs: <context>Content to wrap in context tags</context>
//
// The package is designed to be simple to use with a focused API. ProcessProject
// is the main entry point for all file processing functionality.
package handoff
