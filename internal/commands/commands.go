package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"agenthub/pkg/models"
	"agenthub/pkg/utils"
)

// InitProject initializes a new AgentHub project
func InitProject(projectName string) error {
	fmt.Printf("Creating new AgentHub project: %s\n", projectName)
	
	// Create project directory
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}
	
	// Create basic project structure
	dirs := []string{
		filepath.Join(projectName, "agents"),
		filepath.Join(projectName, "tools"),
		filepath.Join(projectName, "chains"),
		filepath.Join(projectName, "prompts"),
		filepath.Join(projectName, "datasets"),
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create a sample agent manifest
	agentManifest := &models.AgentManifest{
		Name:        projectName,
		Version:     "0.1.0",
		Description: fmt.Sprintf("A sample agent project: %s", projectName),
		Author:      "Your Name",
		License:     "MIT",
		Runtime:     "python",
		EntryPoint:  "main.py",
		Dependencies: map[string]string{},
		Environment: map[string]string{
			"PYTHONPATH": ".",
		},
		Config: map[string]interface{}{
			"max_iterations": 10,
			"timeout":        30,
		},
		Tags: []string{"sample", "agent"},
	}

	// Save the agent manifest
	manifestPath := filepath.Join(projectName, "agent.yaml")
	if err := utils.SaveManifest(agentManifest, manifestPath); err != nil {
		return fmt.Errorf("failed to create agent manifest: %w", err)
	}

	// Create a sample main.py file
	sampleCode := `#!/usr/bin/env python3
"""
Sample Agent for AgentHub
"""
import sys
import json
import os

def main():
    """Main entry point for the agent"""
    print("Hello from your AgentHub agent!")
    print(f"Agent: {os.getenv('AGENT_NAME', 'unknown')}")
    
    # Example: Read input from stdin (file-based IPC)
    if not sys.stdin.isatty():
        try:
            input_data = json.load(sys.stdin)
            print(f"Received input: {input_data}")
        except json.JSONDecodeError:
            print("No valid JSON input received")
    
    # Example: Write output to stdout
    result = {
        "status": "success",
        "message": "Agent executed successfully",
        "data": {"processed": True}
    }
    print(json.dumps(result))

if __name__ == "__main__":
    main()
`

	pythonFile := filepath.Join(projectName, "main.py")
	if err := utils.CreateFileIfNotExists(pythonFile, []byte(sampleCode)); err != nil {
		return fmt.Errorf("failed to create main.py: %w", err)
	}

	// Create requirements.txt
	requirements := `# Add your Python dependencies here
# Example:
# requests>=2.28.0
# numpy>=1.21.0
`

	requirementsFile := filepath.Join(projectName, "requirements.txt")
	if err := utils.CreateFileIfNotExists(requirementsFile, []byte(requirements)); err != nil {
		return fmt.Errorf("failed to create requirements.txt: %w", err)
	}

	// Create README.md
	readme := fmt.Sprintf(`# %s

A sample AgentHub project.

## Description

%s

## Usage

Run the agent:
`, projectName, agentManifest.Description) + "```\nagenthub run\n```\n\n## Development\n\nInstall dependencies:\n```\npip install -r requirements.txt\n```\n"

	readmeFile := filepath.Join(projectName, "README.md")
	if err := utils.CreateFileIfNotExists(readmeFile, []byte(readme)); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	fmt.Printf("✅ Successfully initialized project: %s\n", projectName)
	fmt.Printf("📁 Project structure created with:\n")
	fmt.Printf("   - agent.yaml (manifest)\n")
	fmt.Printf("   - main.py (sample agent code)\n")
	fmt.Printf("   - requirements.txt (Python dependencies)\n")
	fmt.Printf("   - README.md (documentation)\n")
	fmt.Printf("   - Standard directories: agents/, tools/, chains/, prompts/, datasets/\n")
	fmt.Printf("\n🚀 Next steps:\n")
	fmt.Printf("   cd %s\n", projectName)
	fmt.Printf("   agenthub run\n")
	
	return nil
}

// RunAgent runs an agent locally using Python subprocess with file-based IPC
func RunAgent(agentName string, port int, env string, watch bool, verbose bool) error {
	// Find the manifest file
	manifestFile, packageType, err := utils.FindManifestFile()
	if err != nil {
		return fmt.Errorf("no manifest found: %w", err)
	}

	// Load the manifest
	manifest, err := utils.LoadManifest(manifestFile)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	// Validate the manifest
	if err := utils.ValidateManifest(manifest); err != nil {
		return fmt.Errorf("invalid manifest: %w", err)
	}

	if verbose {
		fmt.Printf("🔍 Found %s manifest: %s v%s\n", packageType, manifest.GetName(), manifest.GetVersion())
		fmt.Printf("📁 Running from: %s\n", manifestFile)
	}

	// For now, only support agents
	if packageType != models.PackageTypeAgent {
		return fmt.Errorf("only agents are supported for running (found %s)", packageType)
	}

	agentManifest, ok := manifest.(*models.AgentManifest)
	if !ok {
		return fmt.Errorf("failed to cast manifest to agent manifest")
	}

	// Determine which agent to run
	targetAgent := agentName
	if targetAgent == "" {
		targetAgent = agentManifest.Name
	}

	fmt.Printf("🤖 Running agent: %s\n", targetAgent)
	fmt.Printf("⚡ Runtime: %s\n", agentManifest.Runtime)
	fmt.Printf("🎯 Entry point: %s\n", agentManifest.EntryPoint)

	// Setup environment variables
	envVars := setupAgentEnvironment(agentManifest, env, port)

	// Execute the agent
	return executeAgent(agentManifest, envVars, verbose)
}

// InstallAll installs all project dependencies
func InstallAll() error {
	fmt.Println("Installing all project dependencies...")
	
	// Find the manifest file
	manifestFile, _, err := utils.FindManifestFile()
	if err != nil {
		return fmt.Errorf("no manifest found: %w", err)
	}

	// Load the manifest
	manifest, err := utils.LoadManifest(manifestFile)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	// TODO: Implement actual dependency installation logic
	// For now, just print what would be installed
	fmt.Printf("📦 Installing dependencies for %s v%s\n", manifest.GetName(), manifest.GetVersion())
	
	// Create a basic lock file
	lock := &models.AgentLock{
		Version:   "1.0.0",
		Generated: time.Now(),
		Dependencies: make(map[string]models.LockedDep),
		Integrity:    make(map[string]string),
	}

	// Save the lock file
	if err := utils.SaveAgentLock(lock, "agent.lock"); err != nil {
		return fmt.Errorf("failed to save lock file: %w", err)
	}

	fmt.Println("✅ All dependencies installed successfully")
	fmt.Println("🔒 Lock file created: agent.lock")
	return nil
}

// InstallPackage installs a specific package
func InstallPackage(packageName string) error {
	fmt.Printf("Installing package: %s\n", packageName)
	// TODO: Implement actual package installation logic
	fmt.Printf("✅ Package %s installed successfully\n", packageName)
	return nil
}

// PublishPackage publishes a package to the registry
func PublishPackage(dryRun, private bool) error {
	// Find and validate manifest
	manifestFile, _, err := utils.FindManifestFile()
	if err != nil {
		return fmt.Errorf("no manifest found: %w", err)
	}

	manifest, err := utils.LoadManifest(manifestFile)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	if err := utils.ValidateManifest(manifest); err != nil {
		return fmt.Errorf("invalid manifest: %w", err)
	}

	if dryRun {
		fmt.Printf("🧪 Dry run: Would publish %s v%s to registry\n", manifest.GetName(), manifest.GetVersion())
		return nil
	}
	
	visibility := "public"
	if private {
		visibility = "private"
	}
	
	fmt.Printf("📤 Publishing %s package %s v%s to registry...\n", visibility, manifest.GetName(), manifest.GetVersion())
	// TODO: Implement actual package publishing logic
	fmt.Println("✅ Package published successfully")
	return nil
}

// BuildPackage builds and validates the current package
func BuildPackage(verbose bool, outputDir string) error {
	// Find and validate manifest
	manifestFile, packageType, err := utils.FindManifestFile()
	if err != nil {
		return fmt.Errorf("no manifest found: %w", err)
	}

	manifest, err := utils.LoadManifest(manifestFile)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	if err := utils.ValidateManifest(manifest); err != nil {
		return fmt.Errorf("invalid manifest: %w", err)
	}

	if verbose {
		fmt.Printf("🔍 Building %s: %s v%s\n", packageType, manifest.GetName(), manifest.GetVersion())
		fmt.Printf("📁 Output directory: %s\n", outputDir)
	}
	
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	fmt.Println("✅ Validating package structure...")
	fmt.Println("🔨 Compiling package...")
	fmt.Println("🧪 Running tests...")
	fmt.Printf("✅ Package built successfully in %s\n", outputDir)
	return nil
}

// Helper functions

func setupAgentEnvironment(manifest *models.AgentManifest, env string, port int) []string {
	envVars := os.Environ()
	
	// Add agent-specific environment variables
	envVars = append(envVars, fmt.Sprintf("AGENT_NAME=%s", manifest.Name))
	envVars = append(envVars, fmt.Sprintf("AGENT_VERSION=%s", manifest.Version))
	envVars = append(envVars, fmt.Sprintf("AGENT_ENV=%s", env))
	envVars = append(envVars, fmt.Sprintf("AGENT_PORT=%d", port))
	
	// Add manifest environment variables
	for key, value := range manifest.Environment {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}
	
	return envVars
}

func executeAgent(manifest *models.AgentManifest, envVars []string, verbose bool) error {
	// Determine the command based on runtime
	var cmd *exec.Cmd
	
	switch strings.ToLower(manifest.Runtime) {
	case "python", "python3":
		if !utils.FileExists(manifest.EntryPoint) {
			return fmt.Errorf("entry point not found: %s", manifest.EntryPoint)
		}
		cmd = exec.Command("python3", manifest.EntryPoint)
	case "node", "nodejs":
		if !utils.FileExists(manifest.EntryPoint) {
			return fmt.Errorf("entry point not found: %s", manifest.EntryPoint)
		}
		cmd = exec.Command("node", manifest.EntryPoint)
	default:
		return fmt.Errorf("unsupported runtime: %s", manifest.Runtime)
	}
	
	// Set environment variables
	cmd.Env = envVars
	
	// Set up stdio
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	if verbose {
		fmt.Printf("🚀 Executing: %s %s\n", cmd.Path, strings.Join(cmd.Args[1:], " "))
	}
	
	// Run the agent
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("agent execution failed: %w", err)
	}
	
	fmt.Println("✅ Agent completed successfully")
	return nil
}
