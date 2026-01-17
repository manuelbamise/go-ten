package prompts

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewModel(t *testing.T) {
	model := NewModel()

	if model.currentStage != Stage1ProjectName {
		t.Errorf("Expected stage %d, got %d", Stage1ProjectName, model.currentStage)
	}

	if len(model.appTypes) == 0 {
		t.Error("appTypes should not be empty")
	}

	if len(model.packages) == 0 {
		t.Error("packages should not be empty")
	}
}

func TestValidateProjectName(t *testing.T) {
	model := NewModel()

	// Test valid names
	validNames := []string{"my-project", "my_project", "myproject", "my-project-123", "."}
	for _, name := range validNames {
		err := model.validateProjectName(name)
		if err != nil {
			t.Errorf("Valid name %s failed validation: %v", name, err)
		}
	}

	// Test invalid names
	invalidNames := []string{"", "my project", "my@project", "my#project", "my/project"}
	for _, name := range invalidNames {
		err := model.validateProjectName(name)
		if err == nil {
			t.Errorf("Invalid name %s should have failed validation", name)
		}
	}
}

func TestGetTargetDir(t *testing.T) {
	model := NewModel()

	// Test with regular project name
	model.projectName = "test-project"
	targetDir := model.getTargetDir()
	expected := "./test-project/"
	if targetDir != expected {
		t.Errorf("Expected %s, got %s", expected, targetDir)
	}

	// Test with current directory
	model.projectName = "."
	targetDir = model.getTargetDir()
	expected = "./"
	if targetDir != expected {
		t.Errorf("Expected %s, got %s", expected, targetDir)
	}
}

func TestGenerationSuccess(t *testing.T) {
	model := NewModel()

	// Initially should be false
	if model.GenerationSuccess() {
		t.Error("GenerationSuccess should initially be false")
	}

	// Set to true
	model.generationSuccess = true
	if !model.GenerationSuccess() {
		t.Error("GenerationSuccess should be true when set")
	}
}

func TestUpdateStage1ProjectName(t *testing.T) {
	model := NewModel()

	// Test character input
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t', 'e', 's', 't'}}
	updatedModel, _ := model.updateStage1(msg)

	if updatedModel.(Model).inputValue != "test" {
		t.Errorf("Expected input value 'test', got '%s'", updatedModel.(Model).inputValue)
	}

	// Test backspace
	msg = tea.KeyMsg{Type: tea.KeyBackspace}
	updatedModel, _ = updatedModel.(Model).updateStage1(msg)

	if updatedModel.(Model).inputValue != "tes" {
		t.Errorf("Expected input value 'tes', got '%s'", updatedModel.(Model).inputValue)
	}
}

func TestUpdateStage2AppType(t *testing.T) {
	model := NewModel()
	model.currentStage = Stage2AppType

	// Test Enter key
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(msg)

	um := updatedModel.(Model)
	if um.currentStage != Stage3Package {
		t.Errorf("Expected stage %d, got %d", Stage3Package, um.currentStage)
	}

	if um.selectedAppType == "" {
		t.Error("selectedAppType should not be empty")
	}
}

func TestUpdateStage3Package(t *testing.T) {
	model := NewModel()
	model.currentStage = Stage3Package

	// Test Enter key
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(msg)

	um := updatedModel.(Model)
	if um.currentStage != Stage4Summary {
		t.Errorf("Expected stage %d, got %d", Stage4Summary, um.currentStage)
	}

	if um.selectedPackage == "" {
		t.Error("selectedPackage should not be empty")
	}
}
