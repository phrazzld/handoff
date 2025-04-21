// Example program demonstrating how to use the handoff library
//
// This program demonstrates integrating Handoff with AI models by:
// 1. Collecting code from a project using the handoff library
// 2. Creating a prompt that includes the code context
// 3. Sending the prompt to Gemini (simulated in this example)
// 4. Writing the response to a PLAN.md file
//
// Usage:
//
//	go run gemini_planner.go --project ./my-project --prompt "Add user authentication" --output PLAN.md
//	go run gemini_planner.go --project ./my-project --prompt-file my-prompt.txt --verbose
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/phrazzld/handoff/lib"
)

func main() {
	// Parse command line arguments
	projectPath := flag.String("project", ".", "Path to the project directory")
	userPrompt := flag.String("prompt", "", "User's description of the work to be done")
	promptFile := flag.String("prompt-file", "", "File containing the user's prompt (alternative to --prompt)")
	outputFile := flag.String("output", "PLAN.md", "Output file for the generated plan")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	exclude := flag.String("exclude", ".exe,.bin,.jpg,.png,.gif,.mp3,.mp4,.avi,.mov", "Comma-separated list of file extensions to exclude")
	excludeNames := flag.String("exclude-names", "node_modules,package-lock.json,yarn.lock", "Comma-separated list of file names to exclude")

	flag.Parse()

	// Validate inputs
	if *userPrompt == "" && *promptFile == "" {
		fmt.Fprintf(os.Stderr, "Error: Either --prompt or --prompt-file must be provided\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Load prompt from file if specified
	finalPrompt := *userPrompt
	if *promptFile != "" {
		content, err := os.ReadFile(*promptFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading prompt file: %v\n", err)
			os.Exit(1)
		}
		finalPrompt = string(content)
	}

	// Configure handoff library
	// Create a new configuration with default settings
	config := handoff.NewConfig()
	
	// Apply command-line options
	config.Verbose = *verbose
	config.Exclude = *exclude          // Extensions to exclude (e.g., ".exe,.bin")
	config.ExcludeNamesStr = *excludeNames  // Files/dirs to exclude (e.g., "node_modules")
	
	// Very important: process the configuration strings into slices
	// This must be called after setting string-based config options
	config.ProcessConfig()

	// Process project files using the library
	// This will collect all the files that match our filter criteria
	// and format them according to the configuration
	content, err := handoff.ProcessProject([]string{*projectPath}, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing project: %v\n", err)
		os.Exit(1)
	}
	
	// Get statistics about the content (useful for LLM context limits)
	chars, lines, tokens := handoff.CalculateStatistics(content)
	if *verbose {
		fmt.Printf("Content statistics:\n")
		fmt.Printf("- Characters: %d\n", chars)
		fmt.Printf("- Lines: %d\n", lines)
		fmt.Printf("- Estimated tokens: %d\n", tokens)
	}

	// Create the final prompt for Gemini
	geminiPrompt := fmt.Sprintf(`I'm working on a coding project and need to create a technical plan for the following work:

%s

Here's the current codebase for context:

%s

Please create a detailed technical plan (PLAN.md) that includes:
1. A clear breakdown of the tasks needed
2. Implementation details for each task
3. Any potential challenges or considerations
4. Testing strategy

Format your response as a markdown document that I can use as my implementation guide.`, finalPrompt, content)

	// This is where you would send the prompt to Gemini API
	// For this example, we'll just output that we'd send this to Gemini
	fmt.Println("Generated Gemini prompt with codebase context.")

	if *verbose {
		fmt.Printf("Prompt length: %d characters\n", len(geminiPrompt))
		fmt.Printf("User request: %s\n", finalPrompt)
	}

	// Here you would make the API call to Gemini and receive the response
	// For demonstration, we'll just create a placeholder response
	geminiResponse := "# Technical Plan\n\n*This would be the response from Gemini with a detailed plan.*"

	// Determine the output file path
	outputPath := *outputFile
	if !filepath.IsAbs(outputPath) {
		// If a relative path is provided, make it relative to the project directory
		outputPath = filepath.Join(*projectPath, outputPath)
	}

	// Use the library's WriteToFile function to write the response
	// This function handles directory creation and error checking
	if err := handoff.WriteToFile(geminiResponse, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing plan to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Plan generated successfully and written to %s\n", outputPath)
}
