package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/manuelbamise/go-ten/internal/prompts"
)

func main() {
	// Create and run the bubbletea program
	p := prompts.NewProgram()

	// Run the program and get the result
	model, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}

	// Type assert to get our model
	m, ok := model.(prompts.Model)
	if !ok {
		fmt.Fprintf(os.Stderr, "Unexpected model type\n")
		os.Exit(1)
	}

	// Get the complete configuration
	config := m.GetConfiguration()

	// Check if user made a selection (project type will be empty if they quit)
	if config.ProjectType != "" {
		// Display the final configuration in structured format
		fmt.Println("Project Configuration:")
		fmt.Printf("- Type: %s\n", config.ProjectType)

		if len(config.Features) > 0 {
			fmt.Printf("- Features: %s\n", strings.Join(config.Features, ", "))
		} else {
			fmt.Println("- Features: (none)")
		}

		os.Exit(0)
	} else {
		// User quit without making a selection
		fmt.Println("No selection made")
		os.Exit(0)
	}
}
