package cmd

import (
    "fmt"
    "os"
    
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "agenthub",
    Short: "A unified package manager for Agents, Tools, Chains, Prompts, and Datasets",
    Long: `A unified package manager, for Agents, Tools, Chains, Prompts, and Datasets.`,
    Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from AgentHub 🎉")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
    return rootCmd.Execute()
}

func init() {
    cobra.OnInitialize(initConfig)
    
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.agenthub.yaml)")
    rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
    
    viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, err := os.UserHomeDir()
        cobra.CheckErr(err)
        
        viper.AddConfigPath(home)
        viper.SetConfigType("yaml")
        viper.SetConfigName(".agenthub")
    }
    
    viper.AutomaticEnv()
    
    if err := viper.ReadInConfig(); err == nil {
        fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
    }
}
