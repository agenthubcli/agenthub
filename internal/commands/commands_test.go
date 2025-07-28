package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"agenthub/pkg/models"
	"agenthub/pkg/utils"
)

func TestInitProject(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	os.Chdir(tempDir)
	
	err := InitProject("test-project")
	assert.NoError(t, err)
	
	// Verify project directory structure
	projectDir := filepath.Join(tempDir, "test-project")
	assert.DirExists(t, projectDir)
	
	expectedDirs := []string{
		"agents",
		"tools", 
		"chains",
		"prompts",
		"datasets",
	}
	
	for _, dir := range expectedDirs {
		dirPath := filepath.Join(projectDir, dir)
		assert.DirExists(t, dirPath, "Directory %s should exist", dir)
	}

	// Verify manifest and files were created
	assert.FileExists(t, filepath.Join(projectDir, "agent.yaml"))
	assert.FileExists(t, filepath.Join(projectDir, "main.py"))
	assert.FileExists(t, filepath.Join(projectDir, "requirements.txt"))
	assert.FileExists(t, filepath.Join(projectDir, "README.md"))
}

func TestInitProjectExistingDirectory(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	os.Chdir(tempDir)
	
	// Create project twice
	err := InitProject("duplicate-project")
	assert.NoError(t, err)
	
	err = InitProject("duplicate-project")
	assert.NoError(t, err) // Should not fail if directory already exists
}

func TestInstallAll(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	os.Chdir(tempDir)
	
	// Create a sample manifest for testing
	createSampleManifest(t, tempDir)
	
	err := InstallAll()
	assert.NoError(t, err)
	
	// Verify lock file was created
	assert.FileExists(t, "agent.lock")
}

func TestInstallPackage(t *testing.T) {
	testCases := []struct {
		name        string
		packageName string
	}{
		{"simple package", "simple-package"},
		{"scoped package", "@scope/package"},
		{"package with version", "package@1.0.0"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := InstallPackage(tc.packageName)
			assert.NoError(t, err)
		})
	}
}

func TestPublishPackage(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	os.Chdir(tempDir)
	
	// Create a sample manifest for testing
	createSampleManifest(t, tempDir)
	
	testCases := []struct {
		name    string
		dryRun  bool
		private bool
	}{
		{"normal publish", false, false},
		{"dry run publish", true, false},
		{"private publish", false, true},
		{"private dry run", true, true},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := PublishPackage(tc.dryRun, tc.private)
			assert.NoError(t, err)
		})
	}
}

func TestBuildPackage(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	os.Chdir(tempDir)
	
	// Create a sample manifest for testing
	createSampleManifest(t, tempDir)
	
	testCases := []struct {
		name      string
		verbose   bool
		outputDir string
	}{
		{"normal build", false, "dist"},
		{"verbose build", true, "build-output"},
		{"custom output", false, "custom"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := BuildPackage(tc.verbose, tc.outputDir)
			assert.NoError(t, err)
			
			// Verify output directory was created
			assert.DirExists(t, tc.outputDir)
		})
	}
}

func TestBuildPackageInvalidPath(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	os.Chdir(tempDir)
	
	// Create a sample manifest so we get past manifest validation
	createSampleManifest(t, tempDir)
	
	// Test with invalid output directory path that cannot be created
	// Use a path that would fail permission-wise on most systems
	invalidPath := "/root/invalid/path/that/cannot/be/created"
	if os.Getenv("CI") != "" {
		// On CI systems, use a different invalid path
		invalidPath = "/invalid/path/that/cannot/be/created"
	}
	
	err := BuildPackage(false, invalidPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create output directory")
}

func TestBuildPackageNoManifest(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	os.Chdir(tempDir)
	
	// Test without any manifest file
	err := BuildPackage(false, "dist")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no manifest found")
}

// Helper function to create a sample manifest for testing
func createSampleManifest(t *testing.T, dir string) {
	manifest := &models.AgentManifest{
		Name:        "test-agent",
		Version:     "1.0.0",
		Description: "Test agent for unit tests",
		Author:      "Test Author",
		License:     "MIT",
		Runtime:     "python",
		EntryPoint:  "main.py",
		Dependencies: map[string]string{},
		Environment: map[string]string{
			"PYTHONPATH": ".",
		},
		Config: map[string]interface{}{
			"test": true,
		},
		Tags: []string{"test"},
	}
	
	manifestPath := filepath.Join(dir, "agent.yaml")
	err := utils.SaveManifest(manifest, manifestPath)
	assert.NoError(t, err)
}
