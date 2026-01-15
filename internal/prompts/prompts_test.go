package prompts

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewModel(t *testing.T) {
	m := NewModel()

	expectedChoices := []string{"Web API", "CLI Tool", "gRPC Service", "Microservice"}

	if len(m.choices) != len(expectedChoices) {
		t.Errorf("Expected %d choices, got %d", len(expectedChoices), len(m.choices))
	}

	for i, choice := range m.choices {
		if choice != expectedChoices[i] {
			t.Errorf("Expected choice %s at index %d, got %s", expectedChoices[i], i, choice)
		}
	}

	if m.cursor != 0 {
		t.Errorf("Expected cursor to be 0, got %d", m.cursor)
	}

	if m.chosen != "" {
		t.Errorf("Expected chosen to be empty, got %s", m.chosen)
	}
}

func TestModelUpdateNavigation(t *testing.T) {
	m := NewModel()

	// Test down arrow
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if newModel.(Model).cursor != 1 {
		t.Errorf("Expected cursor to be 1 after down arrow, got %d", newModel.(Model).cursor)
	}

	// Test up arrow
	m = newModel.(Model)
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if newModel.(Model).cursor != 0 {
		t.Errorf("Expected cursor to be 0 after up arrow, got %d", newModel.(Model).cursor)
	}
}

func TestModelUpdateSelection(t *testing.T) {
	m := NewModel()

	// Select the first option
	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	selectedModel := newModel.(Model)
	if selectedModel.chosen != "Web API" {
		t.Errorf("Expected chosen to be 'Web API', got %s", selectedModel.chosen)
	}

	if cmd == nil {
		t.Error("Expected tea.Quit command, got nil")
	}
}

func TestModelUpdateQuit(t *testing.T) {
	m := NewModel()

	// Test quit with 'q'
	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	selectedModel := newModel.(Model)
	if !selectedModel.quitting {
		t.Error("Expected quitting to be true")
	}

	if cmd == nil {
		t.Error("Expected tea.Quit command, got nil")
	}
}

func TestModelView(t *testing.T) {
	m := NewModel()

	view := m.View()

	// Check that view is not empty
	if view == "" {
		t.Error("Expected non-empty view")
	}

	// Note: We can't easily test the full view content due to ANSI escape sequences
	// but we can verify it's not empty
}

func TestGetChosen(t *testing.T) {
	m := NewModel()

	// Initially no choice
	if m.GetChosen() != "" {
		t.Errorf("Expected empty chosen, got %s", m.GetChosen())
	}

	// After selection
	m.chosen = "CLI Tool"
	if m.GetChosen() != "CLI Tool" {
		t.Errorf("Expected 'CLI Tool', got %s", m.GetChosen())
	}
}
