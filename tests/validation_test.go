package tests

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/yourusername/agenthub/pkg"
)

func TestValidateYAMLSpecs(t *testing.T) {
    // Assuming that LoadAgentPkg loads a sample agent package
    agentPkg, err := pkg.LoadAgentPkg("sample_agent.yaml")
    assert.NoError(t, err)
    err = pkg.ValidateAgentPkg(agentPkg)
    assert.NoError(t, err)
}

func TestValidatePackageDependencies(t *testing.T) {
    // Example of validating package dependencies
    required := "^1.0.0"
    available := []string{"1.0.0", "1.0.1"}
    version, err := pkg.ResolveVersion(required, available)
    assert.NoError(t, err)
    assert.Equal(t, "1.0.0", version)
}