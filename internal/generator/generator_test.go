package generator

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetCurrentDirName(t *testing.T) {
	dirName, err := GetCurrentDirName()
	if err != nil {
		t.Fatalf("GetCurrentDirName failed: %v", err)
	}

	if dirName == "" {
		t.Error("GetCurrentDirName returned empty string")
	}

	// Should not be a path, just a name
	if strings.Contains(dirName, string(filepath.Separator)) {
		t.Errorf("GetCurrentDirName returned path instead of name: %s", dirName)
	}
}

func TestCreateDirectory(t *testing.T) {
	testDir := "test_create_dir"
	defer os.RemoveAll(testDir)

	// Test creating new directory
	err := createDirectory(testDir)
	if err != nil {
		t.Fatalf("createDirectory failed: %v", err)
	}

	// Verify directory exists
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Error("Directory was not created")
	}

	// Test creating existing directory (should not error)
	err = createDirectory(testDir)
	if err != nil {
		t.Errorf("createDirectory failed on existing directory: %v", err)
	}
}

func TestProcessTemplate(t *testing.T) {
	templateContent := "Project: {{.ProjectName}}, Type: {{.AppType}}"
	config := ProjectConfig{
		ProjectName: "test-project",
		AppType:     "web-api",
	}

	result, err := processTemplate(templateContent, config)
	if err != nil {
		t.Fatalf("processTemplate failed: %v", err)
	}

	expected := "Project: test-project, Type: web-api"
	if result != expected {
		t.Errorf("processTemplate result mismatch. Got: %s, Expected: %s", result, expected)
	}
}

func TestGetTemplateFS(t *testing.T) {
	// Test with valid template
	templateFS, err := getTemplateFS("web-api", "stdlib")
	if err != nil {
		t.Fatalf("getTemplateFS failed: %v", err)
	}

	// Verify we can read files
	files, err := fs.ReadDir(templateFS, ".")
	if err != nil {
		t.Fatalf("Failed to read template directory: %v", err)
	}

	if len(files) == 0 {
		t.Error("No files found in template filesystem")
	}

	// Test with invalid template
	_, err = getTemplateFS("invalid", "invalid")
	if err == nil {
		t.Error("getTemplateFS should have failed with invalid template")
	}
}

func TestGenerate(t *testing.T) {
	testDir := "test_generate_output"
	defer os.RemoveAll(testDir)

	config := ProjectConfig{
		ProjectName:   "test-project",
		ModuleName:    "test-project",
		AppType:       "web-api",
		Package:       "stdlib",
		TargetDir:     testDir,
		UseCurrentDir: false,
	}

	err := Generate(config)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Verify files were created
	goModPath := filepath.Join(testDir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		t.Error("go.mod file was not created")
	}

	testTxtPath := filepath.Join(testDir, "test.txt")
	if _, err := os.Stat(testTxtPath); os.IsNotExist(err) {
		t.Error("test.txt file was not created")
	}

	// Verify content
	content, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}

	expectedContent := "module test-project"
	if !strings.Contains(string(content), expectedContent) {
		t.Errorf("go.mod content mismatch. Got: %s, Expected to contain: %s", string(content), expectedContent)
	}
}
