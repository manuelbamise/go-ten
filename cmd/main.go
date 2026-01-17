package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/manuelbamise/go-ten/internal/prompts"
)

func main() {
	// Create and run the bubbletea program
	p := prompts.NewProgram()

	// Run the program and get the result
	model, err := p.Run()
	if err != nil {
		if err == tea.ErrInterrupted {
			fmt.Println("\nOperation cancelled by user")
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}

	// Type assert to get our model
	m, ok := model.(prompts.Model)
	if !ok {
		fmt.Fprintf(os.Stderr, "Unexpected model type\n")
		os.Exit(1)
	}

	// Check if user completed the process successfully
	if m.GenerationSuccess() {
		fmt.Println("Project generation completed successfully!")
		os.Exit(0)
	} else {
		// User quit without completing or there was an error
		fmt.Println("Project generation cancelled or failed")
		os.Exit(0)
	}
}
