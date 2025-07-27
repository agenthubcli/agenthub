package main

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCompleteWorkflow(t *testing.T) {
    // Initialize
    err := ExecuteCommand("init")
    assert.NoError(t, err)

    // Build
    err = ExecuteCommand("build")
    assert.NoError(t, err)

    // Publish
    err = ExecuteCommand("publish")
    assert.NoError(t, err)

    // Install
    err = ExecuteCommand("install")
    assert.NoError(t, err)
} 