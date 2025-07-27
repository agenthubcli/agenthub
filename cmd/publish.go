package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"agenthub/internal/commands"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish package to the registry",
	Long: `Publish your agent package, tool, chain, prompt, or dataset to the AgentHub registry.
This will make your package available for others to install and use.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Publishing package to registry...")
		
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		private, _ := cmd.Flags().GetBool("private")
		
		return commands.PublishPackage(dryRun, private)
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)
	publishCmd.Flags().BoolP("dry-run", "d", false, "perform a dry run without actually publishing")
	publishCmd.Flags().BoolP("private", "p", false, "publish as private package")
	publishCmd.Flags().StringP("registry", "r", "default", "specify the registry to publish to")
} 