package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
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
	err := InstallAll()
	assert.NoError(t, err)
	// Since this is a stub implementation, we just verify it doesn't error
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
	// Test with invalid output directory path
	err := BuildPackage(false, "/invalid/path/that/cannot/be/created")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create output directory")
}
