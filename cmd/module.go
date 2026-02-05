package cmd

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/w31r4/dotm/config"
	"gopkg.in/yaml.v3"
)

var moduleCmd = &cobra.Command{
	Use:   "module",
	Short: "Manage modules in the config.yaml file",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available modules",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("Error loading config from %s: %v", configPath, err)
		}
		fmt.Println("Available modules:")
		// Sort keys for consistent output
		keys := make([]string, 0, len(cfg.Modules))
		for k := range cfg.Modules {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("- %s: %s\n", k, cfg.Modules[k].Description)
		}
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove [module]",
	Short: "Remove a module from config.yaml",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		moduleName := args[0]
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("Error loading config from %s: %v", configPath, err)
		}

		if _, ok := cfg.Modules[moduleName]; !ok {
			log.Fatalf("Module '%s' not found.", moduleName)
		}

		delete(cfg.Modules, moduleName)

		data, err := yaml.Marshal(cfg)
		if err != nil {
			log.Fatalf("Error marshaling config: %v", err)
		}

		if err := os.WriteFile(configPath, data, 0644); err != nil {
			log.Fatalf("Error writing config file: %v", err)
		}
		fmt.Printf("Successfully removed module '%s'\n", moduleName)
	},
}

var addCmd = &cobra.Command{
	Use:   "add [module]",
	Short: "Add a new module to config.yaml",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		moduleName := args[0]
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			// If file doesn't exist, create a new config
			if os.IsNotExist(err) {
				cfg = &config.Config{Modules: make(map[string]config.Module)}
			} else {
				log.Fatalf("Error loading config from %s: %v", configPath, err)
			}
		}

		if _, ok := cfg.Modules[moduleName]; ok {
			log.Fatalf("Module '%s' already exists.", moduleName)
		}

		// Collect flags
		desc, _ := cmd.Flags().GetString("description")
		check, _ := cmd.Flags().GetString("check")
		deps, _ := cmd.Flags().GetStringSlice("dependencies")
		installDebian, _ := cmd.Flags().GetStringSlice("install-debian")
		installMacos, _ := cmd.Flags().GetStringSlice("install-macos")
		installArch, _ := cmd.Flags().GetStringSlice("install-arch")
		installDefault, _ := cmd.Flags().GetStringSlice("install-default")

		newModule := config.Module{
			Description:  desc,
			Check:        check,
			Dependencies: deps,
			Install:      make(map[string][]string),
		}

		if len(installDebian) > 0 {
			newModule.Install["debian"] = installDebian
		}
		if len(installMacos) > 0 {
			newModule.Install["macos"] = installMacos
		}
		if len(installArch) > 0 {
			newModule.Install["arch"] = installArch
		}
		if len(installDefault) > 0 {
			newModule.Install["default"] = installDefault
		}

		cfg.Modules[moduleName] = newModule

		data, err := yaml.Marshal(cfg)
		if err != nil {
			log.Fatalf("Error marshaling config: %v", err)
		}

		if err := os.WriteFile(configPath, data, 0644); err != nil {
			log.Fatalf("Error writing config file: %v", err)
		}
		fmt.Printf("Successfully added module '%s'\n", moduleName)
	},
}

func init() {
	rootCmd.AddCommand(moduleCmd)
	moduleCmd.AddCommand(listCmd)
	moduleCmd.AddCommand(removeCmd)
	moduleCmd.AddCommand(addCmd)

	// Flags for the 'add' command
	addCmd.Flags().String("description", "", "Module description")
	addCmd.Flags().String("check", "", "Command to check if the module is installed")
	addCmd.Flags().StringSlice("dependencies", []string{}, "Comma-separated list of dependencies")
	addCmd.Flags().StringSlice("install-debian", []string{}, "Install command(s) for Debian/Ubuntu")
	addCmd.Flags().StringSlice("install-macos", []string{}, "Install command(s) for macOS")
	addCmd.Flags().StringSlice("install-arch", []string{}, "Install command(s) for Arch Linux")
	addCmd.Flags().StringSlice("install-default", []string{}, "Default install command(s)")
}
