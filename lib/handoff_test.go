package handoff

import (
	"testing"
)

// TestNewConfig tests the NewConfig function
func TestNewConfig(t *testing.T) {
	config := NewConfig()
	
	// Verify the default values are set correctly
	if config.Verbose != false {
		t.Errorf("Default Verbose value should be false, got %v", config.Verbose)
	}
	
	if config.Format != "<{path}>\n```\n{content}\n```\n</{path}>\n\n" {
		t.Errorf("Default Format value is incorrect, got %q", config.Format)
	}
	
	if config.Include != "" {
		t.Errorf("Default Include value should be empty, got %q", config.Include)
	}
	
	if config.Exclude != "" {
		t.Errorf("Default Exclude value should be empty, got %q", config.Exclude)
	}
	
	if config.ExcludeNamesStr != "" {
		t.Errorf("Default ExcludeNamesStr value should be empty, got %q", config.ExcludeNamesStr)
	}
	
	if len(config.IncludeExts) != 0 {
		t.Errorf("Default IncludeExts should be empty, got %v", config.IncludeExts)
	}
	
	if len(config.ExcludeExts) != 0 {
		t.Errorf("Default ExcludeExts should be empty, got %v", config.ExcludeExts)
	}
	
	if len(config.ExcludeNames) != 0 {
		t.Errorf("Default ExcludeNames should be empty, got %v", config.ExcludeNames)
	}
}