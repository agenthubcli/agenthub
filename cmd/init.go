package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"agenthub/internal/commands"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new AgentHub project",
	Long: `Initialize a new AgentHub project with the default structure and configuration.
This will create the necessary files and directories for an AgentHub project.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := "my-agent-project"
		if len(args) > 0 {
			projectName = args[0]
		}
		
		fmt.Printf("Initializing AgentHub project: %s\n", projectName)
		return commands.InitProject(projectName)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolP("template", "t", false, "use a template for initialization")
	initCmd.Flags().StringP("registry", "r", "default", "specify the registry to use")
} 