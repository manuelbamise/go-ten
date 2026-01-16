package commands

import (
	"bufio"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestValidateInput(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid simple text",
			input:   "Hello World",
			wantErr: false,
		},
		{
			name:    "valid text with spaces",
			input:   "Project Type: Web API",
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "input with null bytes",
			input:   "test\x00",
			wantErr: true,
		},
		{
			name:    "input with backtick",
			input:   "test`malicious",
			wantErr: true,
		},
		{
			name:    "input with dollar sign",
			input:   "test$HOME",
			wantErr: true,
		},
		{
			name:    "input with pipe",
			input:   "test|rm -rf",
			wantErr: true,
		},
		{
			name:    "input with ampersand",
			input:   "test&evil",
			wantErr: true,
		},
		{
			name:    "input with semicolon",
			input:   "test;rm",
			wantErr: true,
		},
		{
			name:    "input with command substitution",
			input:   "test$(rm)",
			wantErr: true,
		},
		{
			name:    "input too long",
			input:   strings.Repeat("a", 1001),
			wantErr: true,
		},
		{
			name:    "input at max length",
			input:   strings.Repeat("a", 1000),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExecuteEcho(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		wantErr    bool
		expectText string
	}{
		{
			name:       "simple text",
			text:       "Hello World",
			wantErr:    false,
			expectText: "Hello World",
		},
		{
			name:       "empty text",
			text:       "",
			wantErr:    true,
			expectText: "",
		},
		{
			name:       "text with spaces",
			text:       "Project Type: Web API",
			wantErr:    false,
			expectText: "Project Type: Web API",
		},
		{
			name:       "text with special characters (allowed)",
			text:       "Feature-Name_v2",
			wantErr:    false,
			expectText: "Feature-Name_v2",
		},
		{
			name:    "text with dangerous characters",
			text:    "test;rm -rf",
			wantErr: true,
		},
		{
			name:    "text with null bytes",
			text:    "test\x00",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout to verify output
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := ExecuteEcho(tt.text)

			// Restore stdout and get captured output
			w.Close()
			os.Stdout = oldStdout
			scanner := bufio.NewScanner(r)
			var output strings.Builder
			for scanner.Scan() {
				if output.Len() > 0 {
					output.WriteString("\n")
				}
				output.WriteString(scanner.Text())
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteEcho() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// ExecuteEcho prints the text (without trailing newline due to TrimSuffix)
			// So the captured output should match exactly what we expect
			if !tt.wantErr && output.String() != tt.expectText {
				t.Errorf("ExecuteEcho() output = %q, want %q", output.String(), tt.expectText)
			}
		})
	}
}

func TestEchoProjectConfig(t *testing.T) {
	tests := []struct {
		name        string
		projectType string
		features    []string
		wantErr     bool
		expectLines []string
	}{
		{
			name:        "valid config with features",
			projectType: "Web API",
			features:    []string{"Docker support", "Authentication"},
			wantErr:     false,
			expectLines: []string{"Web API", "Docker support", "Authentication"},
		},
		{
			name:        "valid config without features",
			projectType: "CLI Tool",
			features:    []string{},
			wantErr:     false,
			expectLines: []string{"CLI Tool"},
		},
		{
			name:        "empty project type",
			projectType: "",
			features:    []string{"Docker support"},
			wantErr:     true,
		},
		{
			name:        "project type with dangerous characters",
			projectType: "test;rm",
			features:    []string{},
			wantErr:     true,
		},
		{
			name:        "feature with dangerous characters",
			projectType: "Web API",
			features:    []string{"Docker support", "test$(malicious)"},
			wantErr:     true,
		},
		{
			name:        "mixed valid and invalid features",
			projectType: "CLI Tool",
			features:    []string{"Valid Feature", "test`invalid"},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout to verify output
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := EchoProjectConfig(tt.projectType, tt.features)

			// Restore stdout and get captured output
			w.Close()
			os.Stdout = oldStdout
			scanner := bufio.NewScanner(r)
			var lines []string
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("EchoProjectConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !equalStringSlices(lines, tt.expectLines) {
				t.Errorf("EchoProjectConfig() output = %v, want %v", lines, tt.expectLines)
			}
		})
	}
}

// Helper function to compare string slices
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestCrossPlatformCommand(t *testing.T) {
	// Test that the function works on current platform
	t.Run("cross-platform echo", func(t *testing.T) {
		text := "Test Message"

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := ExecuteEcho(text)

		// Restore stdout and get captured output
		w.Close()
		os.Stdout = oldStdout
		scanner := bufio.NewScanner(r)
		var output strings.Builder
		for scanner.Scan() {
			output.WriteString(scanner.Text() + "\n")
		}
		outputStr := strings.TrimSpace(output.String())

		// Should work on any platform
		if err != nil {
			t.Errorf("ExecuteEcho() failed on %s: %v", runtime.GOOS, err)
		}

		// ExecuteEcho outputs the text, so we expect to find it
		if outputStr != text {
			t.Errorf("ExecuteEcho() output on %s = %q, want %q", runtime.GOOS, outputStr, text)
		}

		if !strings.Contains(outputStr, text) {
			t.Errorf("ExecuteEcho() output on %s = %q, want to contain %q", runtime.GOOS, outputStr, text)
		}
	})
}
