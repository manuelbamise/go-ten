package prompts

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewModel(t *testing.T) {
	m := NewModel()

	// Test initial stage
	if m.currentStage != Stage1ProjectType {
		t.Errorf("Expected initial stage to be Stage1ProjectType, got %v", m.currentStage)
	}

	// Test project types initialization
	expectedProjectTypes := []string{"Web API", "CLI Tool", "gRPC Service", "Microservice"}
	if len(m.projectTypes) != len(expectedProjectTypes) {
		t.Errorf("Expected %d project types, got %d", len(expectedProjectTypes), len(m.projectTypes))
	}

	for i, projectType := range m.projectTypes {
		if projectType != expectedProjectTypes[i] {
			t.Errorf("Expected project type %s at index %d, got %s", expectedProjectTypes[i], i, projectType)
		}
	}

	// Test features initialization
	expectedFeatures := []string{
		"Docker support",
		"GitHub Actions CI/CD",
		"PostgreSQL integration",
		"Authentication (JWT)",
		"Logging (structured)",
	}
	if len(m.availableFeatures) != len(expectedFeatures) {
		t.Errorf("Expected %d available features, got %d", len(expectedFeatures), len(m.availableFeatures))
	}

	// Test that all features are initially unselected
	for _, feature := range expectedFeatures {
		if m.selectedFeatures[feature] {
			t.Errorf("Expected feature %s to be unselected initially", feature)
		}
	}

	// Test initial cursor positions
	if m.projectTypeCursor != 0 {
		t.Errorf("Expected project type cursor to be 0, got %d", m.projectTypeCursor)
	}

	if m.featureCursor != 0 {
		t.Errorf("Expected feature cursor to be 0, got %d", m.featureCursor)
	}

	// Test initial state
	if m.selectedProjectType != "" {
		t.Errorf("Expected selected project type to be empty, got %s", m.selectedProjectType)
	}

	if m.quitting {
		t.Error("Expected quitting to be false initially")
	}
}

func TestModelUpdateStage1Navigation(t *testing.T) {
	m := NewModel()

	// Test down arrow in stage 1
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if newModel.(Model).projectTypeCursor != 1 {
		t.Errorf("Expected project type cursor to be 1 after down arrow, got %d", newModel.(Model).projectTypeCursor)
	}

	// Test up arrow in stage 1
	m = newModel.(Model)
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if newModel.(Model).projectTypeCursor != 0 {
		t.Errorf("Expected project type cursor to be 0 after up arrow, got %d", newModel.(Model).projectTypeCursor)
	}

	// Test that we can't go below 0
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if newModel.(Model).projectTypeCursor != 0 {
		t.Errorf("Expected project type cursor to remain 0 when trying to go up from 0, got %d", newModel.(Model).projectTypeCursor)
	}

	// Test that we can't go above max
	m.projectTypeCursor = len(m.projectTypes) - 1
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if newModel.(Model).projectTypeCursor != len(m.projectTypes)-1 {
		t.Errorf("Expected project type cursor to remain at max when trying to go down from max, got %d", newModel.(Model).projectTypeCursor)
	}
}

func TestModelUpdateStage1Transition(t *testing.T) {
	m := NewModel()

	// Test enter key advances to stage 2
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	selectedModel := newModel.(Model)

	if selectedModel.currentStage != Stage2Features {
		t.Errorf("Expected stage to advance to Stage2Features, got %v", selectedModel.currentStage)
	}

	if selectedModel.selectedProjectType != "Web API" {
		t.Errorf("Expected selected project type to be 'Web API', got %s", selectedModel.selectedProjectType)
	}
}

func TestModelUpdateStage2Navigation(t *testing.T) {
	m := NewModel()
	// Advance to stage 2
	m.currentStage = Stage2Features

	// Test down arrow in stage 2
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if newModel.(Model).featureCursor != 1 {
		t.Errorf("Expected feature cursor to be 1 after down arrow, got %d", newModel.(Model).featureCursor)
	}

	// Test up arrow in stage 2
	m = newModel.(Model)
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if newModel.(Model).featureCursor != 0 {
		t.Errorf("Expected feature cursor to be 0 after up arrow, got %d", newModel.(Model).featureCursor)
	}
}

func TestModelUpdateStage2FeatureToggle(t *testing.T) {
	m := NewModel()
	// Advance to stage 2
	m.currentStage = Stage2Features

	// Test spacebar toggles feature
	feature := m.availableFeatures[0]
	if m.selectedFeatures[feature] {
		t.Errorf("Expected feature %s to be unselected initially", feature)
	}

	// Toggle to selected
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	selectedModel := newModel.(Model)
	if !selectedModel.selectedFeatures[feature] {
		t.Errorf("Expected feature %s to be selected after spacebar, got %v", feature, selectedModel.selectedFeatures[feature])
	}

	// Toggle back to unselected
	newModel, _ = selectedModel.Update(tea.KeyMsg{Type: tea.KeySpace})
	selectedModel = newModel.(Model)
	if selectedModel.selectedFeatures[feature] {
		t.Errorf("Expected feature %s to be unselected after second spacebar, got %v", feature, selectedModel.selectedFeatures[feature])
	}
}

func TestModelUpdateStage2Transition(t *testing.T) {
	m := NewModel()
	// Advance to stage 2
	m.currentStage = Stage2Features

	// Test enter key advances to stage 3
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	selectedModel := newModel.(Model)

	if selectedModel.currentStage != Stage3Summary {
		t.Errorf("Expected stage to advance to Stage3Summary, got %v", selectedModel.currentStage)
	}
}

func TestModelUpdateStage3Confirmation(t *testing.T) {
	m := NewModel()
	// Advance to stage 3
	m.currentStage = Stage3Summary
	m.selectedProjectType = "Web API"

	// Test enter key quits
	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	selectedModel := newModel.(Model)

	if !selectedModel.quitting {
		t.Error("Expected quitting to be true after enter in stage 3")
	}

	if cmd == nil {
		t.Error("Expected tea.Quit command, got nil")
	}
}

func TestModelUpdateQuit(t *testing.T) {
	m := NewModel()

	// Test quit with 'q' in any stage
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

	// Test stage 1 view
	view := m.View()
	if view == "" {
		t.Error("Expected non-empty view for stage 1")
	}

	// Test stage 2 view
	m.currentStage = Stage2Features
	view = m.View()
	if view == "" {
		t.Error("Expected non-empty view for stage 2")
	}

	// Test stage 3 view
	m.currentStage = Stage3Summary
	m.selectedProjectType = "Web API"
	view = m.View()
	if view == "" {
		t.Error("Expected non-empty view for stage 3")
	}
}

func TestGetConfiguration(t *testing.T) {
	m := NewModel()

	// Initially empty configuration
	config := m.GetConfiguration()
	if config.ProjectType != "" {
		t.Errorf("Expected empty project type, got %s", config.ProjectType)
	}

	if len(config.Features) != 0 {
		t.Errorf("Expected no features, got %v", config.Features)
	}

	// After selection
	m.selectedProjectType = "CLI Tool"
	m.selectedFeatures["Docker support"] = true
	m.selectedFeatures["Authentication (JWT)"] = true

	config = m.GetConfiguration()
	if config.ProjectType != "CLI Tool" {
		t.Errorf("Expected project type 'CLI Tool', got %s", config.ProjectType)
	}

	expectedFeatures := []string{"Docker support", "Authentication (JWT)"}
	if len(config.Features) != len(expectedFeatures) {
		t.Errorf("Expected %d features, got %d", len(expectedFeatures), len(config.Features))
	}
}

func TestGetSelectedFeatures(t *testing.T) {
	m := NewModel()

	// Initially no features selected
	features := m.getSelectedFeatures()
	if len(features) != 0 {
		t.Errorf("Expected no selected features, got %v", features)
	}

	// Select some features
	m.selectedFeatures["Docker support"] = true
	m.selectedFeatures["Authentication (JWT)"] = true

	features = m.getSelectedFeatures()
	if len(features) != 2 {
		t.Errorf("Expected 2 selected features, got %d", len(features))
	}

	// Check that the correct features are selected
	expectedFeatures := []string{"Docker support", "Authentication (JWT)"}
	for i, feature := range features {
		if feature != expectedFeatures[i] {
			t.Errorf("Expected feature %s at index %d, got %s", expectedFeatures[i], i, feature)
		}
	}
}
