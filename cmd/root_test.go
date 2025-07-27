package cmd

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	// Test that the root command exists and has the correct structure
	assert.NotNil(t, rootCmd)
	assert.Equal(t, "agenthub", rootCmd.Use)
	assert.Contains(t, rootCmd.Short, "unified package manager")
	
	// Test that it has child commands
	commands := rootCmd.Commands()
	assert.True(t, len(commands) >= 4, "Should have at least 4 subcommands")
}

func TestRootCommandHasExpectedSubcommands(t *testing.T) {
	expectedCommands := []string{"init", "install", "publish", "build"}
	
	for _, expectedCmd := range expectedCommands {
		cmd := findCommand(rootCmd, expectedCmd)
		assert.NotNil(t, cmd, "Command %s should exist", expectedCmd)
	}
}

func TestRootCommandFlags(t *testing.T) {
	// Test persistent flags
	assert.NotNil(t, rootCmd.PersistentFlags().Lookup("config"))
	assert.NotNil(t, rootCmd.PersistentFlags().Lookup("verbose"))
} 