package cmd

import (
    "bytes"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/spf13/cobra"
)

// Helper to capture command output
func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
    buf := new(bytes.Buffer)
    root.SetOut(buf)
    root.SetErr(buf)
    root.SetArgs(args)

    err = root.Execute()
    return buf.String(), err
}

func TestRootCommand(t *testing.T) {
    // Just test that the root command exists and has the correct structure
    assert.NotNil(t, rootCmd)
    assert.Equal(t, "agenthub", rootCmd.Use)
    assert.Contains(t, rootCmd.Short, "unified package manager")
    
    // Test that it has child commands
    commands := rootCmd.Commands()
    assert.True(t, len(commands) >= 4, "Should have at least 4 subcommands")
}

func TestInitCommand(t *testing.T) {
    // Find init command in the command list
    var initCmd *cobra.Command
    for _, cmd := range rootCmd.Commands() {
        if cmd.Name() == "init" {
            initCmd = cmd
            break
        }
    }
    
    assert.NotNil(t, initCmd, "Init command should exist")
    assert.Contains(t, initCmd.Short, "Initialize")
    
    // Test flags exist
    assert.NotNil(t, initCmd.Flags().Lookup("template"))
    assert.NotNil(t, initCmd.Flags().Lookup("registry"))
}

func TestInstallCommand(t *testing.T) {
    var installCmd *cobra.Command
    for _, cmd := range rootCmd.Commands() {
        if cmd.Name() == "install" {
            installCmd = cmd
            break
        }
    }
    
    assert.NotNil(t, installCmd, "Install command should exist")
    assert.Contains(t, installCmd.Short, "Install")
    
    // Test flags exist
    assert.NotNil(t, installCmd.Flags().Lookup("dev"))
    assert.NotNil(t, installCmd.Flags().Lookup("version"))
    assert.NotNil(t, installCmd.Flags().Lookup("global"))
}

func TestPublishCommand(t *testing.T) {
    var publishCmd *cobra.Command
    for _, cmd := range rootCmd.Commands() {
        if cmd.Name() == "publish" {
            publishCmd = cmd
            break
        }
    }
    
    assert.NotNil(t, publishCmd, "Publish command should exist")
    assert.Contains(t, publishCmd.Short, "Publish")
    
    // Test flags exist
    assert.NotNil(t, publishCmd.Flags().Lookup("dry-run"))
    assert.NotNil(t, publishCmd.Flags().Lookup("private"))
    assert.NotNil(t, publishCmd.Flags().Lookup("registry"))
}

func TestBuildCommand(t *testing.T) {
    var buildCmd *cobra.Command
    for _, cmd := range rootCmd.Commands() {
        if cmd.Name() == "build" {
            buildCmd = cmd
            break
        }
    }
    
    assert.NotNil(t, buildCmd, "Build command should exist")
    assert.Contains(t, buildCmd.Short, "Build")
    
    // Test flags exist
    assert.NotNil(t, buildCmd.Flags().Lookup("verbose"))
    assert.NotNil(t, buildCmd.Flags().Lookup("output"))
    assert.NotNil(t, buildCmd.Flags().Lookup("watch"))
} 