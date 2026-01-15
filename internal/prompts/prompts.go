package prompts

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the state of our selection UI
type Model struct {
	choices  []string // Available project template options
	cursor   int      // Current cursor position
	chosen   string   // Selected option (empty until selection)
	quitting bool     // Flag to indicate quit state
}

// NewModel creates a new model with default values
func NewModel() Model {
	return Model{
		choices: []string{
			"Web API",
			"CLI Tool",
			"gRPC Service",
			"Microservice",
		},
		cursor: 0,
	}
}

// Init initializes the bubbletea program
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages and updates the model state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Quit keys
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		// Navigation keys
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// Selection key
		case "enter":
			m.chosen = m.choices[m.cursor]
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	s := "Create a new Go project\n\n"

	// Render the list of choices
	for i, choice := range m.choices {
		// Cursor indicator
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		// Highlight the currently selected option
		if m.cursor == i {
			s += fmt.Sprintf("%s \x1b[1m%s\x1b[0m\n", cursor, choice)
		} else {
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}
	}

	s += "\n(press enter to select, q to quit)"
	return s
}

// GetChosen returns the selected option after the program exits
func (m Model) GetChosen() string {
	return m.chosen
}

// NewProgram creates and returns a new bubbletea program for project selection
func NewProgram() *tea.Program {
	return tea.NewProgram(NewModel())
}
