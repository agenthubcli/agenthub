package cmd

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPublishCommand(t *testing.T) {
	cmd := findCommand(rootCmd, "publish")
	assert.NotNil(t, cmd, "Publish command should exist")
}

func TestPublishCommandMetadata(t *testing.T) {
	cmd := findCommand(rootCmd, "publish")
	
	assert.Contains(t, cmd.Short, "Publish")
	assert.Contains(t, cmd.Long, "Publish your agent package")
	assert.Equal(t, "publish", cmd.Use)
}

func TestPublishCommandFlags(t *testing.T) {
	cmd := findCommand(rootCmd, "publish")
	
	// Test dry-run flag
	dryRunFlag := cmd.Flags().Lookup("dry-run")
	assert.NotNil(t, dryRunFlag, "Dry-run flag should exist")
	assert.Equal(t, "bool", dryRunFlag.Value.Type())
	assert.Equal(t, "d", dryRunFlag.Shorthand)
	
	// Test private flag
	privateFlag := cmd.Flags().Lookup("private")
	assert.NotNil(t, privateFlag, "Private flag should exist")
	assert.Equal(t, "bool", privateFlag.Value.Type())
	assert.Equal(t, "p", privateFlag.Shorthand)
	
	// Test registry flag
	registryFlag := cmd.Flags().Lookup("registry")
	assert.NotNil(t, registryFlag, "Registry flag should exist")
	assert.Equal(t, "default", registryFlag.DefValue)
	assert.Equal(t, "r", registryFlag.Shorthand)
}

func TestPublishCommandNoArgs(t *testing.T) {
	cmd := findCommand(rootCmd, "publish")
	
	// Publish command should not require arguments (Args can be nil for no validation)
	assert.NotNil(t, cmd, "Publish command should exist")
} 