package cmd

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCliCommands(t *testing.T) {
    // Test init command
    output, err := ExecuteCliCommand("init")
    assert.NoError(t, err)
    assert.Contains(t, output, "Initializing a new AgentHub project...")

    // Test install command
    output, err = ExecuteCliCommand("install")
    assert.NoError(t, err)
    assert.Contains(t, output, "Installing packages from the registry...")

    // Test publish command
    output, err = ExecuteCliCommand("publish")
    assert.NoError(t, err)
    assert.Contains(t, output, "Publishing package to the GitHub registry...")
} 