package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"agenthub/internal/commands"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [agent-name]",
	Short: "Run an installed agent locally",
	Long: `Run an agent locally using the default Python runtime with subprocess execution.
The agent will be executed with file-based IPC for communication.

If no agent name is provided, it will look for a default agent in the current project.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		agentName := ""
		if len(args) > 0 {
			agentName = args[0]
		}
		
		port, _ := cmd.Flags().GetInt("port")
		env, _ := cmd.Flags().GetString("env")
		watch, _ := cmd.Flags().GetBool("watch")
		verbose, _ := cmd.Flags().GetBool("verbose")
		
		fmt.Printf("Running agent: %s\n", agentName)
		return commands.RunAgent(agentName, port, env, watch, verbose)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().IntP("port", "p", 8080, "port to run the agent on")
	runCmd.Flags().StringP("env", "e", "development", "environment to run in")
	runCmd.Flags().BoolP("watch", "w", false, "watch for changes and restart")
	runCmd.Flags().BoolP("verbose", "v", false, "verbose output")
} 