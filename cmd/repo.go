package cmd

import (
	"dotm/pkg/executor"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage the dotfiles repository",
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Clone dotfiles repo and checkout files",
	RunE: func(cmd *cobra.Command, args []string) error {
		repoURL, _ := cmd.Flags().GetString("url")
		dotfilesDir := os.Getenv("HOME") + "/.dotfiles"

		// 1. Clone the bare repository
		if _, err := os.Stat(dotfilesDir); os.IsNotExist(err) {
			fmt.Println("Cloning dotfiles bare repository...")
			cloneCmd := fmt.Sprintf("git clone --bare %s %s", repoURL, dotfilesDir)
			if err := executor.Execute(cloneCmd, dryRun); err != nil {
				return fmt.Errorf("failed to clone repo: %w", err)
			}
		} else {
			fmt.Println("Dotfiles repository already exists. Skipping clone.")
		}

		// For simplicity, we will shell out to git for the checkout logic
		// A more robust solution would use a Go git library.
		// This part is complex due to potential conflicts.
		// The original shell script logic for backing up is non-trivial.
		fmt.Println("Checking out dotfiles...")
		fmt.Println("NOTE: Conflict handling is not yet implemented in this Go version.")
		fmt.Println("You may need to manually back up conflicting files.")

		checkoutCmd := fmt.Sprintf("/usr/bin/git --git-dir=%s --work-tree=%s checkout", dotfilesDir, os.Getenv("HOME"))
		if err := executor.Execute(checkoutCmd, dryRun); err != nil {
			fmt.Printf("Checkout command failed. This might be due to conflicts with existing files.\n")
			// We don't return error here to allow the process to continue
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(repoCmd)
	repoCmd.AddCommand(syncCmd)
	syncCmd.Flags().String("url", "git@github.com:w31r4/dotfiles.git", "The URL of the dotfiles repository")
}
