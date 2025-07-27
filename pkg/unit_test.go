package pkg

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestValidateAgentPkg(t *testing.T) {
    agentPkg := &AgentPkg{
        Name:    "test-agent",
        Version: "1.0.0",
    }
    err := ValidateAgentPkg(agentPkg)
    assert.NoError(t, err)
}

func TestLoadAgentPkg(t *testing.T) {
    agentPkg, err := LoadAgentPkg("valid_agent.yaml")
    assert.NoError(t, err)
    assert.NotEmpty(t, agentPkg)
} 