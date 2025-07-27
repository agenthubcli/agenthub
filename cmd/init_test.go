package cmd

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestInitCommand(t *testing.T) {
	cmd := findCommand(rootCmd, "init")
	assert.NotNil(t, cmd, "Init command should exist")
}

func TestInitCommandMetadata(t *testing.T) {
	cmd := findCommand(rootCmd, "init")
	
	assert.Contains(t, cmd.Short, "Initialize")
	assert.Contains(t, cmd.Long, "Initialize a new AgentHub project")
	assert.Equal(t, "init [project-name]", cmd.Use)
}

func TestInitCommandFlags(t *testing.T) {
	cmd := findCommand(rootCmd, "init")
	
	// Test that required flags exist
	templateFlag := cmd.Flags().Lookup("template")
	assert.NotNil(t, templateFlag, "Template flag should exist")
	assert.Equal(t, "bool", templateFlag.Value.Type())
	
	registryFlag := cmd.Flags().Lookup("registry")
	assert.NotNil(t, registryFlag, "Registry flag should exist")
	assert.Equal(t, "default", registryFlag.DefValue)
}

func TestInitCommandArgs(t *testing.T) {
	cmd := findCommand(rootCmd, "init")
	
	// Test argument validation - should accept 0 or 1 args
	assert.NotNil(t, cmd.Args)
} 