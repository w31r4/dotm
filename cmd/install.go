package cmd

import (
	"dotm/config"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [module]",
	Short: "Install and configure a module",
	Long: `Install a module defined in the config.yaml file.
This command will check for dependencies, install the required software,
and apply the necessary dotfile configurations.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		moduleName := args[0]
		fmt.Printf("Attempting to install module: %s\n", moduleName)

		cfg, err := config.LoadConfig("config.yaml")
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}

		module, ok := cfg.Modules[moduleName]
		if !ok {
			log.Fatalf("Error: module '%s' not found in config.yaml", moduleName)
		}

		fmt.Printf("Found module: %s\n", module.Description)
		// TODO: Implement the actual installation logic here
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}