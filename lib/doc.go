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
//   - Offers both a high-level API for common use cases and more granular control when needed
//
// # Basic Usage
//
// The primary entry point is the ProcessProject function, which handles everything
// from file discovery to filtering and formatting:
//
//	// Create default configuration and customize as needed
//	config := handoff.NewConfig()
//	config.Include = ".go,.md"          // Only include Go and Markdown files
//	config.ExcludeNamesStr = "go.sum"   // Skip go.sum files
//	config.ProcessConfig()              // Process string-based config into slice-based filters
//
//	// Process files and get formatted content
//	content, err := handoff.ProcessProject([]string{"./src", "README.md"}, config)
//	if err != nil {
//	    log.Fatalf("Error processing project: %v", err)
//	}
//
//	// Use the content with an AI assistant or save to a file
//	handoff.WriteToFile(content, "output.md", true)
//
// # Customization Options
//
// The library offers several ways to customize the behavior through configuration options:
//
//	// Create default configuration with custom options
//	config := handoff.NewConfig(
//	    handoff.WithInclude(".go,.md"),         // Only include Go and Markdown files
//	    handoff.WithExclude(".min.js"),         // Exclude minified JavaScript
//	    handoff.WithExcludeNames("go.sum"),     // Skip specific files by name
//	    handoff.WithVerbose(true),              // Enable verbose logging
//	    handoff.WithFormat("## {path}\n{content}\n\n"), // Custom output format
//	)
//
//	// Process files and get formatted content
//	content, stats, err := handoff.ProcessProject([]string{"./src", "README.md"}, config)
//
//	// Use stats to report information about the processed content
//	fmt.Printf("Processed %d/%d files\n", stats.FilesProcessed, stats.FilesTotal)
//	fmt.Printf("Total content: %d lines, %d chars, ~%d tokens\n", 
//	    stats.Lines, stats.Chars, stats.Tokens)
//
// The package is designed to be flexible yet simple to use, making it suitable
// for a variety of scenarios where code needs to be collected and shared.
package handoff
