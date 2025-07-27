package pkg

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

// Tests for AgentPkg validation
func TestValidateAgentPkg(t *testing.T) {
    agentPkg := &AgentPkg{
        Name:    "test-agent",
        Version: "1.0.0",
    }
    err := ValidateAgentPkg(agentPkg)
    assert.NoError(t, err)
}

func TestValidateAgentPkgNil(t *testing.T) {
    err := ValidateAgentPkg(nil)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "cannot be nil")
}

func TestValidateAgentPkgMissingName(t *testing.T) {
    agentPkg := &AgentPkg{
        Version: "1.0.0",
    }
    err := ValidateAgentPkg(agentPkg)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "name is required")
}

// Tests for AgentPkg loading
func TestLoadAgentPkg(t *testing.T) {
    agentPkg, err := LoadAgentPkg("valid_agent.yaml")
    assert.NoError(t, err)
    assert.NotEmpty(t, agentPkg)
}

func TestLoadAgentPkgEmptyFilename(t *testing.T) {
    agentPkg, err := LoadAgentPkg("")
    assert.Error(t, err)
    assert.Nil(t, agentPkg)
    assert.Contains(t, err.Error(), "filename cannot be empty")
}

// Tests for YAML specs validation (moved from validation_test.go)
func TestValidateYAMLSpecs(t *testing.T) {
    // Load and validate a sample agent package
    agentPkg, err := LoadAgentPkg("sample_agent.yaml")
    assert.NoError(t, err)
    assert.NotNil(t, agentPkg)
    err = ValidateAgentPkg(agentPkg)
    assert.NoError(t, err)
}

// Tests for version resolution
func TestResolveVersion(t *testing.T) {
    required := "^1.0.0"
    available := []string{"1.0.0", "1.0.1"}
    version, err := ResolveVersion(required, available)
    assert.NoError(t, err)
    assert.Equal(t, "1.0.0", version)
}

func TestResolveVersionEmptyRequired(t *testing.T) {
    available := []string{"1.0.0", "1.0.1"}
    version, err := ResolveVersion("", available)
    assert.Error(t, err)
    assert.Empty(t, version)
    assert.Contains(t, err.Error(), "required version cannot be empty")
}

func TestResolveVersionNoAvailable(t *testing.T) {
    version, err := ResolveVersion("^1.0.0", []string{})
    assert.Error(t, err)
    assert.Empty(t, version)
    assert.Contains(t, err.Error(), "no available versions")
}

// Tests for package dependencies validation (moved from validation_test.go)
func TestValidatePackageDependencies(t *testing.T) {
    // Example of validating package dependencies
    required := "^1.0.0"
    available := []string{"1.0.0", "1.0.1"}
    version, err := ResolveVersion(required, available)
    assert.NoError(t, err)
    assert.Equal(t, "1.0.0", version)
} 