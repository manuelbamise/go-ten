package main

import (
	"fmt"
	"os"

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

	// Get the chosen option
	chosen := m.GetChosen()

	if chosen != "" {
		fmt.Printf("You chose: %s\n", chosen)
		os.Exit(0)
	} else {
		// User quit without making a selection
		fmt.Println("No selection made")
		os.Exit(0)
	}
}
