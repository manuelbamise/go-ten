package generator

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/*
var templateFS embed.FS

// ProjectConfig holds the configuration for project generation
type ProjectConfig struct {
	ProjectName   string // e.g., "my-api" or extracted from pwd
	ModuleName    string // same as ProjectName for now
	AppType       string // "web-api"
	Package       string // "stdlib"
	TargetDir     string // "./my-api/" or "./"
	UseCurrentDir bool   // true if user entered "."
}

// Generate is the main orchestration function for project generation
func Generate(config ProjectConfig) error {
	// Create target directory if not using current dir
	if !config.UseCurrentDir {
		if err := createDirectory(config.TargetDir); err != nil {
			return fmt.Errorf("failed to create target directory: %w", err)
		}
	}

	// Get the embedded template filesystem for the config
	templateFS, err := getTemplateFS(config.AppType, config.Package)
	if err != nil {
		return fmt.Errorf("failed to get template filesystem: %w", err)
	}

	// Walk through template files and copy them
	if err := copyTemplateFiles(templateFS, config.TargetDir, config); err != nil {
		return fmt.Errorf("failed to copy template files: %w", err)
	}

	return nil
}

// getTemplateFS returns the embedded filesystem for specific template
func getTemplateFS(appType, packageName string) (fs.FS, error) {
	// Template path format: "templates/{appType}-{packageName}"
	// Example: "templates/web-api-stdlib"
	templatePath := fmt.Sprintf("templates/%s-%s", appType, packageName)

	templateSubFS, err := fs.Sub(templateFS, templatePath)
	if err != nil {
		return nil, fmt.Errorf("template not found: %s (available templates: check templates/ directory)", templatePath)
	}

	// Verify the template exists by checking if we can read at least one file
	_, err = fs.ReadDir(templateSubFS, ".")
	if err != nil {
		return nil, fmt.Errorf("template directory is empty or invalid: %s", templatePath)
	}

	return templateSubFS, nil
}

// createDirectory creates directory and all parent directories
func createDirectory(path string) error {
	// Check if directory already exists
	if _, err := os.Stat(path); err == nil {
		return nil // Directory already exists
	}

	// Create directory with all parent directories
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}

	return nil
}

// copyTemplateFiles walks through all files in templateFS and copies them
func copyTemplateFiles(templateFS fs.FS, targetDir string, config ProjectConfig) error {
	return fs.WalkDir(templateFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory
		if path == "." {
			return nil
		}

		// Construct the full target path
		targetPath := filepath.Join(targetDir, path)

		// If it's a directory, create it
		if d.IsDir() {
			return createDirectory(targetPath)
		}

		// Handle files
		return copyFile(templateFS, path, targetPath, config)
	})
}

// copyFile copies a single file from template to target, processing templates if needed
func copyFile(templateFS fs.FS, sourcePath, targetPath string, config ProjectConfig) error {
	// Read the source file
	sourceContent, err := fs.ReadFile(templateFS, sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source file %s: %w", sourcePath, err)
	}

	var finalContent string
	var finalPath string

	// Check if it's a template file (.tmpl extension)
	if strings.HasSuffix(sourcePath, ".tmpl") {
		// Process the template
		processedContent, err := processTemplate(string(sourceContent), config)
		if err != nil {
			return fmt.Errorf("failed to process template %s: %w", sourcePath, err)
		}
		finalContent = processedContent
		// Remove .tmpl extension from target path
		finalPath = strings.TrimSuffix(targetPath, ".tmpl")
	} else {
		// Copy file as-is
		finalContent = string(sourceContent)
		finalPath = targetPath
	}

	// Write the target file
	if err := os.WriteFile(finalPath, []byte(finalContent), 0644); err != nil {
		return fmt.Errorf("failed to write target file %s: %w", finalPath, err)
	}

	return nil
}

// processTemplate uses text/template to replace variables
func processTemplate(templateContent string, config ProjectConfig) (string, error) {
	// Create and parse the template
	tmpl, err := template.New("project").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute the template with config data
	var result strings.Builder
	if err := tmpl.Execute(&result, config); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return result.String(), nil
}

// GetCurrentDirName gets the current working directory name (not full path)
func GetCurrentDirName() (string, error) {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Extract just the directory name from the path
	dirName := filepath.Base(cwd)

	// Handle edge case where we're at root
	if dirName == "." || dirName == "/" {
		return "", fmt.Errorf("cannot determine project name from root directory")
	}

	return dirName, nil
}
