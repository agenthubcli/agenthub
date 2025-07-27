package main

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "agenthub/internal/commands"
)

func TestCompleteWorkflow(t *testing.T) {
    // Initialize
    err := commands.ExecuteCommand("init")
    assert.NoError(t, err)

    // Build
    err = commands.ExecuteCommand("build")
    assert.NoError(t, err)

    // Publish
    err = commands.ExecuteCommand("publish")
    assert.NoError(t, err)

    // Install
    err = commands.ExecuteCommand("install")
    assert.NoError(t, err)
} 