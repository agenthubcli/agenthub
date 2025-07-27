package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"agenthub/internal/commands"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build agent package",
	Long: `Build and validate your agent package, tool, chain, prompt, or dataset.
This will compile, validate, and prepare your package for distribution.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Building agent package...")
		
		verbose, _ := cmd.Flags().GetBool("verbose")
		output, _ := cmd.Flags().GetString("output")
		
		return commands.BuildPackage(verbose, output)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().Bool("verbose", false, "verbose build output")
	buildCmd.Flags().StringP("output", "o", "dist", "output directory for built package")
	buildCmd.Flags().BoolP("watch", "w", false, "watch for changes and rebuild")
} 