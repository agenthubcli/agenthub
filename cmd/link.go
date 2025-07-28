package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"agenthub/internal/commands"
)

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link [package-path]",
	Short: "Link local packages for development and testing",
	Long: `Link local packages for development and testing purposes.
This creates symbolic links to local packages, allowing you to test 
changes without publishing to the registry.

Usage examples:
  agenthub link                    # Link current package globally
  agenthub link ../my-tool         # Link specific package from path
  agenthub link --list             # List all linked packages
  agenthub link --unlink my-tool   # Unlink a specific package`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packagePath := ""
		if len(args) > 0 {
			packagePath = args[0]
		}
		
		list, _ := cmd.Flags().GetBool("list")
		unlink, _ := cmd.Flags().GetString("unlink")
		global, _ := cmd.Flags().GetBool("global")
		verbose, _ := cmd.Flags().GetBool("verbose")
		
		if list {
			fmt.Println("Listing all linked packages...")
			return commands.ListLinkedPackages()
		}
		
		if unlink != "" {
			fmt.Printf("Unlinking package: %s\n", unlink)
			return commands.UnlinkPackage(unlink, verbose)
		}
		
		fmt.Printf("Linking package: %s\n", packagePath)
		return commands.LinkPackage(packagePath, global, verbose)
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
	linkCmd.Flags().BoolP("list", "l", false, "list all linked packages")
	linkCmd.Flags().StringP("unlink", "u", "", "unlink a specific package")
	linkCmd.Flags().BoolP("global", "g", false, "link package globally")
	linkCmd.Flags().BoolP("verbose", "v", false, "verbose output")
} 