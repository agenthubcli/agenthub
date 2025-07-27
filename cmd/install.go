package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"agenthub/internal/commands"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [package-name]",
	Short: "Install packages from the registry",
	Long: `Install agent packages, tools, chains, prompts, or datasets from the AgentHub registry.
If no package name is provided, it will install all dependencies from the project file.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println("Installing all project dependencies...")
			return commands.InstallAll()
		}
		
		packageName := args[0]
		fmt.Printf("Installing package: %s\n", packageName)
		return commands.InstallPackage(packageName)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolP("dev", "d", false, "install development dependencies")
	installCmd.Flags().String("version", "latest", "specify version to install")
	installCmd.Flags().BoolP("global", "g", false, "install package globally")
} 