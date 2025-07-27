package cmd

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestInstallCommand(t *testing.T) {
	cmd := findCommand(rootCmd, "install")
	assert.NotNil(t, cmd, "Install command should exist")
}

func TestInstallCommandMetadata(t *testing.T) {
	cmd := findCommand(rootCmd, "install")
	
	assert.Contains(t, cmd.Short, "Install")
	assert.Contains(t, cmd.Long, "Install agent packages")
	assert.Equal(t, "install [package-name]", cmd.Use)
}

func TestInstallCommandFlags(t *testing.T) {
	cmd := findCommand(rootCmd, "install")
	
	// Test development dependencies flag
	devFlag := cmd.Flags().Lookup("dev")
	assert.NotNil(t, devFlag, "Dev flag should exist")
	assert.Equal(t, "bool", devFlag.Value.Type())
	assert.Equal(t, "d", devFlag.Shorthand)
	
	// Test version flag
	versionFlag := cmd.Flags().Lookup("version")
	assert.NotNil(t, versionFlag, "Version flag should exist")
	assert.Equal(t, "latest", versionFlag.DefValue)
	
	// Test global flag
	globalFlag := cmd.Flags().Lookup("global")
	assert.NotNil(t, globalFlag, "Global flag should exist")
	assert.Equal(t, "bool", globalFlag.Value.Type())
	assert.Equal(t, "g", globalFlag.Shorthand)
}

func TestInstallCommandArgs(t *testing.T) {
	cmd := findCommand(rootCmd, "install")
	
	// Should accept 0 or 1 args (package name is optional)
	assert.NotNil(t, cmd.Args)
} 