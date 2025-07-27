package tests

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/yourusername/agenthub/pkg"
)

func TestValidateAgentPkg(t *testing.T) {
    agentPkg := pkg.AgentPkg{/* Initialize with test data */}
    err := pkg.ValidateAgentPkg(agentPkg)
    assert.NoError(t, err)
}

func TestLoadAgentPkg(t *testing.T) {
    agentPkg, err := pkg.LoadAgentPkg("valid_agent.yaml")
    assert.NoError(t, err)
    assert.NotEmpty(t, agentPkg)
}