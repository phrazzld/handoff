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
//	handoff.WriteToFile(content, "output.md")
//
// # Advanced Usage
//
// For more control over the processing, the library provides lower-level functions
// that can be combined as needed:
//
//	config := handoff.NewConfig()
//	config.ProcessConfig()
//	logger := handoff.NewLogger(true) // Enable verbose logging
//
//	// Create a custom processor function
//	processor := func(file string, content []byte) string {
//	    return fmt.Sprintf("File: %s\n%s\n---\n", file, string(content))
//	}
//
//	builder := &strings.Builder{}
//	for _, path := range paths {
//	    handoff.ProcessPathWithProcessor(path, builder, config, logger, processor)
//	}
//
//	// Get and use the result
//	result := builder.String()
//
// The package is designed to be flexible yet simple to use, making it suitable
// for a variety of scenarios where code needs to be collected and shared.
package handoff
