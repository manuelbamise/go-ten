package prompts

import (
	"fmt"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/manuelbamise/go-ten/internal/generator"
)

// Stage represents the current stage in the multi-step flow
type Stage int

const (
	Stage1ProjectName Stage = iota
	Stage2AppType
	Stage3Package
	Stage4Summary
	Stage5Success
)

// Model represents the state of our multi-step selection UI
type Model struct {
	// Stage management
	currentStage Stage

	// Stage 1: Project Name Input
	projectName string
	inputValue  string
	inputCursor int

	// Stage 2: Application Type Selection
	appTypes        []string
	appTypeCursor   int
	selectedAppType string

	// Stage 3: Package Selection
	packages        []string
	packageCursor   int
	selectedPackage string

	// Stage 4: Summary
	quitting bool

	// Generation state
	generationError   error
	generationSuccess bool
}

// NewModel creates a new model with default values
func NewModel() Model {
	appTypes := []string{
		"Web API",
	}

	packages := []string{
		"stdlib",
	}

	return Model{
		currentStage:  Stage1ProjectName,
		appTypes:      appTypes,
		appTypeCursor: 0,
		packages:      packages,
		packageCursor: 0,
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
		// Quit keys (available in all stages except success)
		case "q", "ctrl+c":
			if m.currentStage != Stage5Success {
				m.quitting = true
				return m, tea.Quit
			}

		// Stage-specific key handling
		default:
			switch m.currentStage {
			case Stage1ProjectName:
				return m.updateStage1(msg)
			case Stage2AppType:
				return m.updateStage2(msg)
			case Stage3Package:
				return m.updateStage3(msg)
			case Stage4Summary:
				return m.updateStage4(msg)
			case Stage5Success:
				return m.updateStage5(msg)
			}
		}
	}

	return m, nil
}

// updateStage1 handles key input for project name input
func (m Model) updateStage1(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		// Validate input
		if err := m.validateProjectName(m.inputValue); err != nil {
			m.generationError = err
			return m, nil
		}

		// Set project name and move to next stage
		m.projectName = m.inputValue
		m.currentStage = Stage2AppType
		return m, nil

	case tea.KeyBackspace:
		if len(m.inputValue) > 0 {
			// Remove character at cursor position
			if m.inputCursor > 0 {
				m.inputValue = m.inputValue[:m.inputCursor-1] + m.inputValue[m.inputCursor:]
				m.inputCursor--
			}
		}

	case tea.KeyLeft:
		if m.inputCursor > 0 {
			m.inputCursor--
		}

	case tea.KeyRight:
		if m.inputCursor < len(m.inputValue) {
			m.inputCursor++
		}

	case tea.KeyRunes:
		// Add character at cursor position
		m.inputValue = m.inputValue[:m.inputCursor] + string(msg.Runes) + m.inputValue[m.inputCursor:]
		m.inputCursor += len(msg.Runes)
	}

	return m, nil
}

// updateStage2 handles key input for application type selection
func (m Model) updateStage2(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	// Navigation keys
	case "up", "k":
		if m.appTypeCursor > 0 {
			m.appTypeCursor--
		}

	case "down", "j":
		if m.appTypeCursor < len(m.appTypes)-1 {
			m.appTypeCursor++
		}

	// Selection key - advance to stage 3
	case "enter":
		m.selectedAppType = m.appTypes[m.appTypeCursor]
		m.currentStage = Stage3Package
	}

	return m, nil
}

// updateStage3 handles key input for package selection
func (m Model) updateStage3(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	// Navigation keys
	case "up", "k":
		if m.packageCursor > 0 {
			m.packageCursor--
		}

	case "down", "j":
		if m.packageCursor < len(m.packages)-1 {
			m.packageCursor++
		}

	// Selection key - advance to stage 4
	case "enter":
		m.selectedPackage = m.packages[m.packageCursor]
		m.currentStage = Stage4Summary
	}

	return m, nil
}

// updateStage4 handles key input for summary stage
func (m Model) updateStage4(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	// Confirm and create project
	case "enter":
		// Generate the project
		if err := m.generateProject(); err != nil {
			m.generationError = err
			return m, nil
		}

		m.generationSuccess = true
		m.currentStage = Stage5Success
	}

	return m, nil
}

// updateStage5 handles key input for success stage
func (m Model) updateStage5(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Any key exits
	m.quitting = true
	return m, tea.Quit
}

// View renders the UI
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	// Render based on current stage
	switch m.currentStage {
	case Stage1ProjectName:
		return m.renderStage1()
	case Stage2AppType:
		return m.renderStage2()
	case Stage3Package:
		return m.renderStage3()
	case Stage4Summary:
		return m.renderStage4()
	case Stage5Success:
		return m.renderStage5()
	default:
		return "Error: Unknown stage"
	}
}

// renderStage1 renders the project name input screen
func (m Model) renderStage1() string {
	s := "Enter your project name (or '.' for current directory):\n\n"

	// Show input field with cursor
	s += "> "

	// Display input with cursor indicator
	for i, char := range m.inputValue {
		if i == m.inputCursor {
			s += fmt.Sprintf("%c|", char)
		} else {
			s += string(char)
		}
	}

	// Show cursor at end if at the end of input
	if m.inputCursor == len(m.inputValue) {
		s += "|"
	}

	// Show error if validation failed
	if m.generationError != nil {
		s += fmt.Sprintf("\n\nError: %v", m.generationError)
		m.generationError = nil // Clear error after displaying
	}

	s += "\n\n(Enter to submit, q to quit)"
	return s
}

// renderStage2 renders the application type selection screen
func (m Model) renderStage2() string {
	s := "Select application type:\n\n"

	// Render the list of application types
	for i, appType := range m.appTypes {
		// Cursor indicator
		cursor := " "
		if m.appTypeCursor == i {
			cursor = ">"
		}

		// Highlight the currently selected option
		if m.appTypeCursor == i {
			s += fmt.Sprintf("%s \x1b[1m%s\x1b[0m\n", cursor, appType)
		} else {
			s += fmt.Sprintf("%s %s\n", cursor, appType)
		}
	}

	s += "\n(Use arrow keys to navigate, press Enter to continue, q to quit)"
	return s
}

// renderStage3 renders the package selection screen
func (m Model) renderStage3() string {
	s := "Select package:\n\n"

	// Render the list of packages
	for i, pkg := range m.packages {
		// Cursor indicator
		cursor := " "
		if m.packageCursor == i {
			cursor = ">"
		}

		// Highlight the currently selected option
		if m.packageCursor == i {
			s += fmt.Sprintf("%s \x1b[1m%s\x1b[0m\n", cursor, pkg)
		} else {
			s += fmt.Sprintf("%s %s\n", cursor, pkg)
		}
	}

	s += "\n(Use arrow keys to navigate, press Enter to continue, q to quit)"
	return s
}

// renderStage4 renders the summary screen
func (m Model) renderStage4() string {
	s := "Project Configuration Summary\n\n"

	// Display project name
	targetDir := m.getTargetDir()
	s += fmt.Sprintf("Name: \x1b[1m%s\x1b[0m\n", m.projectName)

	// Display selected application type
	s += fmt.Sprintf("Type: \x1b[1m%s\x1b[0m\n", m.selectedAppType)

	// Display selected package
	s += fmt.Sprintf("Package: \x1b[1m%s\x1b[0m\n", m.selectedPackage)

	// Display target location
	s += fmt.Sprintf("Location: \x1b[1m%s\x1b[0m\n", targetDir)

	// Show error if generation failed
	if m.generationError != nil {
		s += fmt.Sprintf("\n\x1b[31mError: %v\x1b[0m\n", m.generationError)
		s += "\nPress Enter to retry or 'q' to quit"
	} else {
		s += "\nPress Enter to generate or 'q' to quit"
	}

	return s
}

// renderStage5 renders the success screen
func (m Model) renderStage5() string {
	s := "\x1b[32m✓ Project created successfully!\x1b[0m\n\n"
	s += "Next steps:\n"

	targetDir := m.getTargetDir()
	if m.projectName != "." {
		s += fmt.Sprintf("cd %s\n", targetDir)
	}

	s += "go mod tidy\n"
	s += "go run ./cmd/api\n\n"
	s += fmt.Sprintf("Your Web API is ready at: %s\n", targetDir)
	s += "\nPress any key to exit"

	return s
}

// validateProjectName validates the project name input
func (m Model) validateProjectName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Allow "." as special case for current directory
	if name == "." {
		return nil
	}

	// Validate project name format (alphanumeric, hyphens, underscores)
	validName := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("project name must contain only letters, numbers, hyphens, and underscores")
	}

	return nil
}

// getTargetDir returns the target directory path
func (m Model) getTargetDir() string {
	if m.projectName == "." {
		return "./"
	}
	return fmt.Sprintf("./%s/", m.projectName)
}

// generateProject creates the project using the generator
func (m Model) generateProject() error {
	// Determine project name and target directory
	var projectName string
	var targetDir string
	var useCurrentDir bool

	if m.projectName == "." {
		// Use current directory name as project name
		currentDir, err := generator.GetCurrentDirName()
		if err != nil {
			return fmt.Errorf("failed to get current directory name: %w", err)
		}
		projectName = currentDir
		targetDir = "./"
		useCurrentDir = true
	} else {
		projectName = m.projectName
		targetDir = fmt.Sprintf("./%s/", m.projectName)
		useCurrentDir = false
	}

	// Create project configuration
	config := generator.ProjectConfig{
		ProjectName:   projectName,
		ModuleName:    projectName,
		AppType:       "web-api",
		Package:       m.selectedPackage,
		TargetDir:     targetDir,
		UseCurrentDir: useCurrentDir,
	}

	// Generate the project
	return generator.Generate(config)
}

// GenerationSuccess returns true if the project was generated successfully
func (m Model) GenerationSuccess() bool {
	return m.generationSuccess
}

// NewProgram creates and returns a new bubbletea program for project selection
func NewProgram() *tea.Program {
	return tea.NewProgram(NewModel())
}
