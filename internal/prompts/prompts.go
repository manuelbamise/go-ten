package prompts

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Stage represents the current stage in the multi-step flow
type Stage int

const (
	Stage1ProjectType Stage = iota
	Stage2Features
	Stage3Summary
)

// Model represents the state of our multi-step selection UI
type Model struct {
	// Stage management
	currentStage Stage

	// Stage 1: Project Type Selection
	projectTypes        []string
	projectTypeCursor   int
	selectedProjectType string

	// Stage 2: Features Selection
	availableFeatures []string
	selectedFeatures  map[string]bool
	featureCursor     int

	// Stage 3: Summary
	quitting bool
}

// NewModel creates a new model with default values
func NewModel() Model {
	projectTypes := []string{
		"Web API",
		"CLI Tool",
		"gRPC Service",
		"Microservice",
	}

	availableFeatures := []string{
		"Docker support",
		"GitHub Actions CI/CD",
		"PostgreSQL integration",
		"Authentication (JWT)",
		"Logging (structured)",
	}

	// Initialize selected features map with all features unselected
	selectedFeatures := make(map[string]bool)
	for _, feature := range availableFeatures {
		selectedFeatures[feature] = false
	}

	return Model{
		currentStage:      Stage1ProjectType,
		projectTypes:      projectTypes,
		projectTypeCursor: 0,
		availableFeatures: availableFeatures,
		selectedFeatures:  selectedFeatures,
		featureCursor:     0,
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
		// Quit keys (available in all stages)
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		// Stage-specific key handling
		default:
			switch m.currentStage {
			case Stage1ProjectType:
				return m.updateStage1(msg)
			case Stage2Features:
				return m.updateStage2(msg)
			case Stage3Summary:
				return m.updateStage3(msg)
			}
		}
	}

	return m, nil
}

// updateStage1 handles key input for project type selection
func (m Model) updateStage1(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	// Navigation keys
	case "up", "k":
		if m.projectTypeCursor > 0 {
			m.projectTypeCursor--
		}

	case "down", "j":
		if m.projectTypeCursor < len(m.projectTypes)-1 {
			m.projectTypeCursor++
		}

	// Selection key - advance to stage 2
	case "enter":
		m.selectedProjectType = m.projectTypes[m.projectTypeCursor]
		m.currentStage = Stage2Features
	}

	return m, nil
}

// updateStage2 handles key input for features selection
func (m Model) updateStage2(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	// Navigation keys
	case "up", "k":
		if m.featureCursor > 0 {
			m.featureCursor--
		}

	case "down", "j":
		if m.featureCursor < len(m.availableFeatures)-1 {
			m.featureCursor++
		}

	// Toggle feature selection
	case " ":
		feature := m.availableFeatures[m.featureCursor]
		m.selectedFeatures[feature] = !m.selectedFeatures[feature]

	// Advance to stage 3
	case "enter":
		m.currentStage = Stage3Summary
	}

	return m, nil
}

// updateStage3 handles key input for summary stage
func (m Model) updateStage3(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	// Confirm and create project
	case "enter":
		m.quitting = true
		return m, tea.Quit
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	// Render based on current stage
	switch m.currentStage {
	case Stage1ProjectType:
		return m.renderStage1()
	case Stage2Features:
		return m.renderStage2()
	case Stage3Summary:
		return m.renderStage3()
	default:
		return "Error: Unknown stage"
	}
}

// renderStage1 renders the project type selection screen
func (m Model) renderStage1() string {
	s := "What type of project do you want to create?\n\n"

	// Render the list of project types
	for i, projectType := range m.projectTypes {
		// Cursor indicator
		cursor := " "
		if m.projectTypeCursor == i {
			cursor = ">"
		}

		// Highlight the currently selected option
		if m.projectTypeCursor == i {
			s += fmt.Sprintf("%s \x1b[1m%s\x1b[0m\n", cursor, projectType)
		} else {
			s += fmt.Sprintf("%s %s\n", cursor, projectType)
		}
	}

	s += "\n(Use arrow keys to navigate, press Enter to continue, q to quit)"
	return s
}

// renderStage2 renders the features selection screen
func (m Model) renderStage2() string {
	s := "Select additional features: (space to toggle, enter to continue)\n\n"

	// Render the list of features with checkboxes
	for i, feature := range m.availableFeatures {
		// Cursor indicator
		cursor := " "
		if m.featureCursor == i {
			cursor = ">"
		}

		// Checkbox state
		checkbox := "[ ]"
		if m.selectedFeatures[feature] {
			checkbox = "[x]"
		}

		// Highlight the currently selected option
		if m.featureCursor == i {
			s += fmt.Sprintf("%s %s \x1b[1m%s\x1b[0m\n", cursor, checkbox, feature)
		} else {
			s += fmt.Sprintf("%s %s %s\n", cursor, checkbox, feature)
		}
	}

	s += "\n(Use arrow keys to navigate, Space to toggle, Enter to continue, q to quit)"
	return s
}

// renderStage3 renders the summary screen
func (m Model) renderStage3() string {
	s := "Project Configuration Summary\n\n"

	// Display selected project type
	s += fmt.Sprintf("Project Type: \x1b[1m%s\x1b[0m\n", m.selectedProjectType)

	// Display selected features
	s += "\nFeatures:\n"
	selectedCount := 0
	for _, feature := range m.availableFeatures {
		if m.selectedFeatures[feature] {
			s += fmt.Sprintf("  • %s\n", feature)
			selectedCount++
		}
	}

	if selectedCount == 0 {
		s += "  (No additional features selected)\n"
	}

	s += "\nPress Enter to create project or 'q' to quit"
	return s
}

// GetConfiguration returns the complete project configuration
func (m Model) GetConfiguration() ProjectConfiguration {
	return ProjectConfiguration{
		ProjectType: m.selectedProjectType,
		Features:    m.getSelectedFeatures(),
	}
}

// getSelectedFeatures returns a slice of selected feature names
func (m Model) getSelectedFeatures() []string {
	var features []string
	for _, feature := range m.availableFeatures {
		if m.selectedFeatures[feature] {
			features = append(features, feature)
		}
	}
	return features
}

// ProjectConfiguration represents the final project configuration
type ProjectConfiguration struct {
	ProjectType string
	Features    []string
}

// NewProgram creates and returns a new bubbletea program for project selection
func NewProgram() *tea.Program {
	return tea.NewProgram(NewModel())
}
