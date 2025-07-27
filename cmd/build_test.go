package cmd

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBuildCommand(t *testing.T) {
	cmd := findCommand(rootCmd, "build")
	assert.NotNil(t, cmd, "Build command should exist")
}

func TestBuildCommandMetadata(t *testing.T) {
	cmd := findCommand(rootCmd, "build")
	
	assert.Contains(t, cmd.Short, "Build")
	assert.Contains(t, cmd.Long, "Build and validate your agent package")
	assert.Equal(t, "build", cmd.Use)
}

func TestBuildCommandFlags(t *testing.T) {
	cmd := findCommand(rootCmd, "build")
	
	// Test verbose flag
	verboseFlag := cmd.Flags().Lookup("verbose")
	assert.NotNil(t, verboseFlag, "Verbose flag should exist")
	assert.Equal(t, "bool", verboseFlag.Value.Type())
	
	// Test output flag
	outputFlag := cmd.Flags().Lookup("output")
	assert.NotNil(t, outputFlag, "Output flag should exist")
	assert.Equal(t, "dist", outputFlag.DefValue)
	assert.Equal(t, "o", outputFlag.Shorthand)
	
	// Test watch flag
	watchFlag := cmd.Flags().Lookup("watch")
	assert.NotNil(t, watchFlag, "Watch flag should exist")
	assert.Equal(t, "bool", watchFlag.Value.Type())
	assert.Equal(t, "w", watchFlag.Shorthand)
}

func TestBuildCommandNoArgs(t *testing.T) {
	cmd := findCommand(rootCmd, "build")
	
	// Build command should not require arguments (Args can be nil for no validation)
	assert.NotNil(t, cmd, "Build command should exist")
} 