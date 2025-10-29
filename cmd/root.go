/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const version = "0.1.0"

var configPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dotm",
	Short: "dotm is a declarative environment bootstrapper.",
	Long: `dotm is a declarative environment bootstrapper designed to simplify
and codify your dotfiles workflow. Define your toolchain once in YAML and
install or update on any machine with confidence.

Key commands include:
- dotm install <modules...> to install specific tools
- dotm module <subcommand> to manage entries in config.yaml
- dotm config download/export to share or back up config
- dotm repo sync to bootstrap your dotfiles repository`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if configPath == "" {
			configPath = "config.yaml"
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file (default is ./config.yaml)")
}
