package utils

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"agenthub/pkg/models"
)

// LoadManifest loads a manifest file and returns the appropriate type
func LoadManifest(filePath string) (models.Manifest, error) {
	if !FileExists(filePath) {
		return nil, fmt.Errorf("manifest file not found: %s", filePath)
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest file: %w", err)
	}

	// Determine manifest type based on filename
	fileName := filepath.Base(filePath)
	switch {
	case strings.HasPrefix(fileName, "agent"):
		return loadAgentManifest(data)
	case strings.HasPrefix(fileName, "tool"):
		return loadToolManifest(data)
	case strings.HasPrefix(fileName, "chain"):
		return loadChainManifest(data)
	case strings.HasPrefix(fileName, "prompt"):
		return loadPromptManifest(data)
	case strings.HasPrefix(fileName, "dataset"):
		return loadDatasetManifest(data)
	case strings.HasPrefix(fileName, "agentpkg"):
		return loadAgentPkgManifest(data)
	default:
		return nil, fmt.Errorf("unknown manifest type: %s", fileName)
	}
}

// LoadAgentLock loads an agent.lock file
func LoadAgentLock(filePath string) (*models.AgentLock, error) {
	if !FileExists(filePath) {
		return nil, fmt.Errorf("lock file not found: %s", filePath)
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read lock file: %w", err)
	}

	var lock models.AgentLock
	if err := yaml.Unmarshal(data, &lock); err != nil {
		return nil, fmt.Errorf("failed to parse lock file: %w", err)
	}

	return &lock, nil
}

// SaveManifest saves a manifest to a file
func SaveManifest(manifest models.Manifest, filePath string) error {
	data, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	if err := EnsureDir(filepath.Dir(filePath)); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write manifest file: %w", err)
	}

	return nil
}

// SaveAgentLock saves an agent.lock file
func SaveAgentLock(lock *models.AgentLock, filePath string) error {
	data, err := yaml.Marshal(lock)
	if err != nil {
		return fmt.Errorf("failed to marshal lock file: %w", err)
	}

	if err := EnsureDir(filepath.Dir(filePath)); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write lock file: %w", err)
	}

	return nil
}

// ValidateManifest validates a manifest
func ValidateManifest(manifest models.Manifest) error {
	if manifest == nil {
		return fmt.Errorf("manifest cannot be nil")
	}

	name := manifest.GetName()
	if name == "" {
		return fmt.Errorf("manifest name is required")
	}

	version := manifest.GetVersion()
	if version == "" {
		return fmt.Errorf("manifest version is required")
	}

	// Additional type-specific validation
	switch m := manifest.(type) {
	case *models.AgentManifest:
		return validateAgentManifest(m)
	case *models.ToolManifest:
		return validateToolManifest(m)
	case *models.ChainManifest:
		return validateChainManifest(m)
	case *models.PromptManifest:
		return validatePromptManifest(m)
	case *models.DatasetManifest:
		return validateDatasetManifest(m)
	case *models.AgentPkgManifest:
		return validateAgentPkgManifest(m)
	}

	return nil
}

// FindManifestFile finds the manifest file in the current directory
func FindManifestFile() (string, models.PackageType, error) {
	manifestFiles := []struct {
		name string
		typ  models.PackageType
	}{
		{"agent.yaml", models.PackageTypeAgent},
		{"agent.yml", models.PackageTypeAgent},
		{"tool.yaml", models.PackageTypeTool},
		{"tool.yml", models.PackageTypeTool},
		{"chain.yaml", models.PackageTypeChain},
		{"chain.yml", models.PackageTypeChain},
		{"prompt.yaml", models.PackageTypePrompt},
		{"prompt.yml", models.PackageTypePrompt},
		{"dataset.yaml", models.PackageTypeDataset},
		{"dataset.yml", models.PackageTypeDataset},
		{"agentpkg.yaml", models.PackageTypeMulti},
		{"agentpkg.yml", models.PackageTypeMulti},
	}

	for _, mf := range manifestFiles {
		if FileExists(mf.name) {
			return mf.name, mf.typ, nil
		}
	}

	return "", "", fmt.Errorf("no manifest file found in current directory")
}

// Private helper functions

func loadAgentManifest(data []byte) (*models.AgentManifest, error) {
	var manifest models.AgentManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse agent manifest: %w", err)
	}
	return &manifest, nil
}

func loadToolManifest(data []byte) (*models.ToolManifest, error) {
	var manifest models.ToolManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse tool manifest: %w", err)
	}
	return &manifest, nil
}

func loadChainManifest(data []byte) (*models.ChainManifest, error) {
	var manifest models.ChainManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse chain manifest: %w", err)
	}
	return &manifest, nil
}

func loadPromptManifest(data []byte) (*models.PromptManifest, error) {
	var manifest models.PromptManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse prompt manifest: %w", err)
	}
	return &manifest, nil
}

func loadDatasetManifest(data []byte) (*models.DatasetManifest, error) {
	var manifest models.DatasetManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse dataset manifest: %w", err)
	}
	return &manifest, nil
}

func loadAgentPkgManifest(data []byte) (*models.AgentPkgManifest, error) {
	var manifest models.AgentPkgManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse agentpkg manifest: %w", err)
	}
	return &manifest, nil
}

func validateAgentManifest(manifest *models.AgentManifest) error {
	if manifest.Runtime == "" {
		return fmt.Errorf("agent runtime is required")
	}
	if manifest.EntryPoint == "" {
		return fmt.Errorf("agent entry_point is required")
	}
	return nil
}

func validateToolManifest(manifest *models.ToolManifest) error {
	if manifest.Runtime == "" {
		return fmt.Errorf("tool runtime is required")
	}
	if manifest.EntryPoint == "" {
		return fmt.Errorf("tool entry_point is required")
	}
	return nil
}

func validateChainManifest(manifest *models.ChainManifest) error {
	if len(manifest.Steps) == 0 {
		return fmt.Errorf("chain must have at least one step")
	}
	for i, step := range manifest.Steps {
		if step.Name == "" {
			return fmt.Errorf("step %d: name is required", i)
		}
		if step.Type == "" {
			return fmt.Errorf("step %d: type is required", i)
		}
		if step.Package == "" {
			return fmt.Errorf("step %d: package is required", i)
		}
	}
	return nil
}

func validatePromptManifest(manifest *models.PromptManifest) error {
	if manifest.Template == "" {
		return fmt.Errorf("prompt template is required")
	}
	return nil
}

func validateDatasetManifest(manifest *models.DatasetManifest) error {
	if manifest.Format == "" {
		return fmt.Errorf("dataset format is required")
	}
	if len(manifest.Files) == 0 {
		return fmt.Errorf("dataset must have at least one file")
	}
	return nil
}

func validateAgentPkgManifest(manifest *models.AgentPkgManifest) error {
	if manifest.Type == "" {
		return fmt.Errorf("agentpkg type is required")
	}
	
	// For multi-type packages, check that components are specified
	if manifest.Type == "multi" {
		components := manifest.Components
		if len(components.Agents) == 0 && len(components.Tools) == 0 && 
		   len(components.Chains) == 0 && len(components.Prompts) == 0 && 
		   len(components.Datasets) == 0 {
			return fmt.Errorf("multi-type package must specify at least one component")
		}
	}
	
	return nil
} 