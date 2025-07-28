package models

import (
	"time"
)

// AgentManifest represents an agent.yaml file
type AgentManifest struct {
	Name         string            `yaml:"name"`
	Version      string            `yaml:"version"`
	Description  string            `yaml:"description"`
	Author       string            `yaml:"author"`
	License      string            `yaml:"license"`
	Runtime      string            `yaml:"runtime"`      // python, node, docker, etc.
	EntryPoint   string            `yaml:"entry_point"`  // main.py, index.js, etc.
	Dependencies map[string]string `yaml:"dependencies"` // tool:version mappings
	Environment  map[string]string `yaml:"environment"`  // env variables
	Config       map[string]interface{} `yaml:"config"`  // agent-specific config
	Tags         []string          `yaml:"tags"`
}

// ToolManifest represents a tool.yaml file
type ToolManifest struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	Description string            `yaml:"description"`
	Author      string            `yaml:"author"`
	License     string            `yaml:"license"`
	Runtime     string            `yaml:"runtime"`
	EntryPoint  string            `yaml:"entry_point"`
	Schema      ToolSchema        `yaml:"schema"`      // input/output schema
	Config      map[string]interface{} `yaml:"config"` // tool-specific config
	Tags        []string          `yaml:"tags"`
}

// ToolSchema defines the input/output schema for a tool
type ToolSchema struct {
	Input  map[string]interface{} `yaml:"input"`
	Output map[string]interface{} `yaml:"output"`
}

// ChainManifest represents a chain.yaml file
type ChainManifest struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	Description string            `yaml:"description"`
	Author      string            `yaml:"author"`
	License     string            `yaml:"license"`
	Steps       []ChainStep       `yaml:"steps"`
	Config      map[string]interface{} `yaml:"config"`
	Tags        []string          `yaml:"tags"`
}

// ChainStep represents a step in a chain
type ChainStep struct {
	Name     string                 `yaml:"name"`
	Type     string                 `yaml:"type"`     // agent, tool, prompt
	Package  string                 `yaml:"package"`  // package reference
	Config   map[string]interface{} `yaml:"config"`
	Inputs   map[string]string      `yaml:"inputs"`   // input mappings
	Outputs  map[string]string      `yaml:"outputs"`  // output mappings
	Condition string                `yaml:"condition,omitempty"` // conditional execution
}

// PromptManifest represents a prompt.yaml file
type PromptManifest struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	Description string            `yaml:"description"`
	Author      string            `yaml:"author"`
	License     string            `yaml:"license"`
	Template    string            `yaml:"template"`    // prompt template
	Variables   []PromptVariable  `yaml:"variables"`   // template variables
	Examples    []PromptExample   `yaml:"examples"`    // example inputs/outputs
	Config      map[string]interface{} `yaml:"config"`
	Tags        []string          `yaml:"tags"`
}

// PromptVariable represents a variable in a prompt template
type PromptVariable struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Default     string `yaml:"default,omitempty"`
}

// PromptExample represents an example for a prompt
type PromptExample struct {
	Name     string                 `yaml:"name"`
	Inputs   map[string]interface{} `yaml:"inputs"`
	Expected string                 `yaml:"expected"`
}

// DatasetManifest represents a dataset.yaml file
type DatasetManifest struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	Description string            `yaml:"description"`
	Author      string            `yaml:"author"`
	License     string            `yaml:"license"`
	Format      string            `yaml:"format"`      // csv, json, parquet, etc.
	Schema      DatasetSchema     `yaml:"schema"`      // data schema
	Files       []DatasetFile     `yaml:"files"`       // dataset files
	Config      map[string]interface{} `yaml:"config"`
	Tags        []string          `yaml:"tags"`
}

// DatasetSchema defines the schema for a dataset
type DatasetSchema struct {
	Columns []DatasetColumn `yaml:"columns"`
}

// DatasetColumn represents a column in a dataset
type DatasetColumn struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Nullable    bool   `yaml:"nullable"`
}

// DatasetFile represents a file in a dataset
type DatasetFile struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
	Size int64  `yaml:"size"`
	Hash string `yaml:"hash"`
}

// AgentPkgManifest represents an agentpkg.yaml file (multi-component package)
type AgentPkgManifest struct {
	Name         string            `yaml:"name"`
	Version      string            `yaml:"version"`
	Type         string            `yaml:"type"`         // multi, agent, tool, etc.
	Description  string            `yaml:"description"`
	Author       string            `yaml:"author"`
	License      string            `yaml:"license"`
	Components   Components        `yaml:"components"`
	Dependencies map[string]string `yaml:"dependencies"`
	Config       map[string]interface{} `yaml:"config"`
	Tags         []string          `yaml:"tags"`
}

// Components represents the components in a multi-package
type Components struct {
	Agents   []string `yaml:"agents,omitempty"`
	Tools    []string `yaml:"tools,omitempty"`
	Chains   []string `yaml:"chains,omitempty"`
	Prompts  []string `yaml:"prompts,omitempty"`
	Datasets []string `yaml:"datasets,omitempty"`
}

// AgentLock represents an agent.lock file (dependency lockfile)
type AgentLock struct {
	Version      string                 `yaml:"version"`       // lockfile format version
	Generated    time.Time              `yaml:"generated"`     // when lockfile was generated
	Dependencies map[string]LockedDep   `yaml:"dependencies"`  // locked dependencies
	Integrity    map[string]string      `yaml:"integrity"`     // integrity hashes
}

// LockedDep represents a locked dependency
type LockedDep struct {
	Version  string `yaml:"version"`
	Resolved string `yaml:"resolved"` // URL or path where it was resolved from
	Hash     string `yaml:"hash"`     // integrity hash
}

// PackageType represents the type of package
type PackageType string

const (
	PackageTypeAgent   PackageType = "agent"
	PackageTypeTool    PackageType = "tool"
	PackageTypeChain   PackageType = "chain"
	PackageTypePrompt  PackageType = "prompt"
	PackageTypeDataset PackageType = "dataset"
	PackageTypeMulti   PackageType = "multi"
)

// Manifest is a union interface for all manifest types
type Manifest interface {
	GetName() string
	GetVersion() string
	GetType() PackageType
}

// Implement the Manifest interface for each type
func (a *AgentManifest) GetName() string     { return a.Name }
func (a *AgentManifest) GetVersion() string  { return a.Version }
func (a *AgentManifest) GetType() PackageType { return PackageTypeAgent }

func (t *ToolManifest) GetName() string     { return t.Name }
func (t *ToolManifest) GetVersion() string  { return t.Version }
func (t *ToolManifest) GetType() PackageType { return PackageTypeTool }

func (c *ChainManifest) GetName() string     { return c.Name }
func (c *ChainManifest) GetVersion() string  { return c.Version }
func (c *ChainManifest) GetType() PackageType { return PackageTypeChain }

func (p *PromptManifest) GetName() string     { return p.Name }
func (p *PromptManifest) GetVersion() string  { return p.Version }
func (p *PromptManifest) GetType() PackageType { return PackageTypePrompt }

func (d *DatasetManifest) GetName() string     { return d.Name }
func (d *DatasetManifest) GetVersion() string  { return d.Version }
func (d *DatasetManifest) GetType() PackageType { return PackageTypeDataset }

func (a *AgentPkgManifest) GetName() string     { return a.Name }
func (a *AgentPkgManifest) GetVersion() string  { return a.Version }
func (a *AgentPkgManifest) GetType() PackageType { return PackageType(a.Type) } 