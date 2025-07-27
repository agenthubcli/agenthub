package pkg

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestValidateYAMLSpecs(t *testing.T) {
    // Assuming that LoadAgentPkg loads a sample agent package
    agentPkg, err := LoadAgentPkg("sample_agent.yaml")
    assert.NoError(t, err)
    assert.NotNil(t, agentPkg)
    err = ValidateAgentPkg(agentPkg)
    assert.NoError(t, err)
}

func TestValidatePackageDependencies(t *testing.T) {
    // Example of validating package dependencies
    required := "^1.0.0"
    available := []string{"1.0.0", "1.0.1"}
    version, err := ResolveVersion(required, available)
    assert.NoError(t, err)
    assert.Equal(t, "1.0.0", version)
} 