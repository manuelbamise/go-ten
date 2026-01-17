package commands

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// ExecuteEcho runs the echo command with the provided text in a secure, cross-platform manner
// It uses os/exec with separate arguments to prevent shell injection vulnerabilities
// The function captures command output and prints it to stdout
// Returns an error if the command fails
func executeEcho(text string) error {
	// Input validation to prevent shell injection and invalid inputs
	if err := validateInput(text); err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}

	// Cross-platform command execution
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// On Windows, use cmd /c echo
		cmd = exec.Command("cmd", "/c", "echo", text)
	} else {
		// On Unix-like systems, use echo directly
		cmd = exec.Command("echo", text)
	}

	// Execute command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("echo command failed: %w", err)
	}

	// Print output to stdout (trim trailing newline from echo command, then add newline for separation)
	fmt.Println(strings.TrimSuffix(string(output), "\n"))
	return nil
}

// EchoProjectConfig echoes the project configuration using ExecuteEcho
// It first echoes the project type, then echoes each selected feature on separate lines
// Provides a clear, readable format for the project configuration
// Returns an error if any echo command fails
func EchoProjectConfig(projectType string, features []string) error {
	// Validate project type input
	if projectType == "" {
		return fmt.Errorf("project type cannot be empty")
	}

	// Echo project type
	if err := executeEcho(projectType); err != nil {
		return fmt.Errorf("failed to echo project type: %w", err)
	}

	// Echo each feature on separate lines
	for _, feature := range features {
		// Validate individual feature input
		if err := validateInput(feature); err != nil {
			return fmt.Errorf("invalid feature input '%s': %w", feature, err)
		}

		if err := executeEcho(feature); err != nil {
			return fmt.Errorf("failed to echo feature '%s': %w", feature, err)
		}
	}

	return nil
}

// validateInput performs security validation on input strings
// Prevents shell injection by checking for dangerous characters and patterns
// Returns an error if validation fails
func validateInput(input string) error {
	if input == "" {
		return fmt.Errorf("input cannot be empty")
	}

	// Check for null bytes (can be used for injection attempts)
	if strings.Contains(input, string([]byte{0})) {
		return fmt.Errorf("input contains null bytes")
	}

	// Check for common shell injection patterns
	dangerousChars := []string{
		"`", "$", "|", "&", ";", "<", ">", "(", ")", "{", "}",
		"[", "]", "!", "*", "?", "~", "#", "%", "^", "=",
	}

	for _, char := range dangerousChars {
		if strings.Contains(input, char) {
			return fmt.Errorf("input contains potentially dangerous character: %s", char)
		}
	}

	// Check for command substitution patterns
	if strings.Contains(input, "$(") || strings.Contains(input, "`") {
		return fmt.Errorf("input contains command substitution patterns")
	}

	// Check length limits to prevent buffer overflow attempts
	if len(input) > 1000 {
		return fmt.Errorf("input too long (max 1000 characters)")
	}

	return nil
}
