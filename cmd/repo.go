package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

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
		dotfilesDir, _ := cmd.Flags().GetString("dir")
		backupBaseDir, _ := cmd.Flags().GetString("backup-dir")
		pullLatest, _ := cmd.Flags().GetBool("pull")

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("get home dir: %w", err)
		}

		dotfilesDir, err = expandHomePath(dotfilesDir, homeDir)
		if err != nil {
			return fmt.Errorf("invalid --dir: %w", err)
		}
		backupBaseDir, err = expandHomePath(backupBaseDir, homeDir)
		if err != nil {
			return fmt.Errorf("invalid --backup-dir: %w", err)
		}

		// 1) Clone the bare repository (if needed)
		if _, err := os.Stat(dotfilesDir); errors.Is(err, os.ErrNotExist) {
			fmt.Println("Cloning dotfiles bare repository...")
			if _, err := runCommand(dryRun, "git", "clone", "--bare", repoURL, dotfilesDir); err != nil {
				return fmt.Errorf("failed to clone repo: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("stat dotfiles dir: %w", err)
		} else {
			fmt.Println("Dotfiles repository already exists. Skipping clone.")
		}

		// 2) Keep the dotfiles repo readable by hiding home-directory untracked files.
		if _, err := runCommand(dryRun, "git", "--git-dir="+dotfilesDir, "config", "status.showUntrackedFiles", "no"); err != nil {
			return fmt.Errorf("failed to configure dotfiles repo: %w", err)
		}

		// 3) Checkout (with conflict backup)
		fmt.Println("Checking out dotfiles...")
		backupDir, err := checkoutWithBackup(dotfilesDir, homeDir, backupBaseDir, dryRun)
		if err != nil {
			return err
		}
		if backupDir != "" {
			fmt.Printf("Backed up conflicting files to: %s\n", backupDir)
		}

		// 4) Update repo (fast-forward only) and re-checkout if needed.
		if pullLatest {
			fmt.Println("Pulling latest changes...")
			pullBackupDir, err := pullWithBackup(dotfilesDir, homeDir, backupBaseDir, dryRun)
			if err != nil {
				return err
			}
			if pullBackupDir != "" {
				fmt.Printf("Backed up conflicting files to: %s\n", pullBackupDir)
			}
		}

		return nil
	},
}

var repoGitCmd = &cobra.Command{
	Use:                "git [git-args...]",
	Short:              "Run git against the dotfiles bare repo",
	DisableFlagParsing: true,
	Args:               cobra.ArbitraryArgs,
	Long: `Run a git command with the dotfiles bare repository and your home directory work tree.

Examples:
  dotm repo git status
  dotm repo git add .zshrc
  dotm repo git commit -m "Update zsh config"

By default, the repo is assumed to be at ~/.dotfiles. Override with DOTM_DOTFILES_DIR.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("get home dir: %w", err)
		}

		dotfilesDir := os.Getenv("DOTM_DOTFILES_DIR")
		if dotfilesDir == "" {
			dotfilesDir = "~/.dotfiles"
		}

		dotfilesDir, err = expandHomePath(dotfilesDir, homeDir)
		if err != nil {
			return fmt.Errorf("invalid DOTM_DOTFILES_DIR: %w", err)
		}

		if _, err := os.Stat(dotfilesDir); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("dotfiles repo not found at %s (run `dotm repo sync` first)", dotfilesDir)
		} else if err != nil {
			return fmt.Errorf("stat dotfiles dir: %w", err)
		}

		gitArgs := append([]string{"--git-dir=" + dotfilesDir, "--work-tree=" + homeDir}, args...)
		_, err = runCommand(false, "git", gitArgs...)
		return err
	},
}

func expandHomePath(path, homeDir string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("empty path")
	}
	if path == "~" {
		return homeDir, nil
	}
	if strings.HasPrefix(path, "~"+string(os.PathSeparator)) {
		return filepath.Join(homeDir, strings.TrimPrefix(path, "~"+string(os.PathSeparator))), nil
	}
	return path, nil
}

func runCommand(dryRun bool, name string, args ...string) (string, error) {
	if dryRun {
		fmt.Printf("[DRY RUN] Would execute: %s\n", formatCommand(name, args))
		return "", nil
	}

	fmt.Printf("Executing: %s\n", formatCommand(name, args))
	command := exec.Command(name, args...)
	out, err := command.CombinedOutput()
	if len(out) > 0 {
		fmt.Print(string(out))
	}
	if err != nil {
		return string(out), fmt.Errorf("%s failed: %w", name, err)
	}
	return string(out), nil
}

func formatCommand(name string, args []string) string {
	if len(args) == 0 {
		return name
	}
	return name + " " + strings.Join(args, " ")
}

func checkoutWithBackup(dotfilesDir, homeDir, backupBaseDir string, dryRun bool) (string, error) {
	return backupAndRetry(dotfilesDir, homeDir, backupBaseDir, dryRun, []string{"checkout"})
}

func pullWithBackup(dotfilesDir, homeDir, backupBaseDir string, dryRun bool) (string, error) {
	return backupAndRetry(dotfilesDir, homeDir, backupBaseDir, dryRun, []string{"pull", "--ff-only"})
}

func backupAndRetry(dotfilesDir, homeDir, backupBaseDir string, dryRun bool, gitArgs []string) (string, error) {
	const maxAttempts = 5
	backupDir := ""
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		args := append([]string{"--git-dir=" + dotfilesDir, "--work-tree=" + homeDir}, gitArgs...)
		out, err := runCommand(dryRun, "git", args...)
		if err == nil {
			return backupDir, nil
		}
		lastErr = err

		conflicts := parseOverwrittenFilePaths(out)
		if len(conflicts) == 0 {
			return backupDir, fmt.Errorf("git %s failed: %w", strings.Join(gitArgs, " "), err)
		}

		if backupDir == "" {
			backupDir = filepath.Join(backupBaseDir, time.Now().Format("20060102-150405.000000000"))
		}

		if err := backupPaths(homeDir, backupDir, conflicts, dryRun); err != nil {
			return backupDir, err
		}
	}

	return backupDir, fmt.Errorf("git %s failed after %d attempts: %w", strings.Join(gitArgs, " "), maxAttempts, lastErr)
}

func backupPaths(homeDir, backupDir string, paths []string, dryRun bool) error {
	if dryRun {
		for _, p := range paths {
			fmt.Printf("[DRY RUN] Would back up: %s -> %s\n", filepath.Join(homeDir, p), filepath.Join(backupDir, p))
		}
		return nil
	}

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("create backup dir: %w", err)
	}

	for _, rawPath := range paths {
		relPath, err := safeRelPath(rawPath)
		if err != nil {
			return err
		}

		srcPath := filepath.Join(homeDir, relPath)
		dstPath := filepath.Join(backupDir, relPath)

		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return fmt.Errorf("create backup parent dir: %w", err)
		}

		if err := os.Rename(srcPath, dstPath); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return fmt.Errorf("backup %s: %w", relPath, err)
		}
	}

	return nil
}

func safeRelPath(p string) (string, error) {
	clean := filepath.Clean(strings.TrimSpace(p))
	if clean == "." || clean == "" {
		return "", fmt.Errorf("invalid path %q", p)
	}
	if filepath.IsAbs(clean) {
		return "", fmt.Errorf("refusing absolute path %q", p)
	}
	if clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("refusing path traversal %q", p)
	}
	return clean, nil
}

func parseOverwrittenFilePaths(gitOutput string) []string {
	var paths []string
	inList := false

	for _, line := range strings.Split(gitOutput, "\n") {
		trimmed := strings.TrimSpace(line)

		if strings.Contains(line, "would be overwritten by checkout:") || strings.Contains(line, "would be overwritten by merge:") {
			inList = true
			continue
		}

		if !inList {
			continue
		}

		if trimmed == "" {
			continue
		}

		// Common terminators after the file list
		if strings.HasPrefix(trimmed, "Please ") || strings.HasPrefix(trimmed, "Aborting") || strings.HasPrefix(trimmed, "hint:") {
			inList = false
			continue
		}

		// If we hit another header-like line, end the current list.
		if strings.HasPrefix(trimmed, "error:") || strings.HasPrefix(trimmed, "fatal:") {
			inList = false
			continue
		}

		paths = append(paths, trimmed)
	}

	paths = slices.Compact(paths)
	return paths
}

func init() {
	rootCmd.AddCommand(repoCmd)
	repoCmd.AddCommand(syncCmd)
	repoCmd.AddCommand(repoGitCmd)
	syncCmd.Flags().String("url", "git@github.com:w31r4/dotfiles.git", "The URL of the dotfiles repository")
	syncCmd.Flags().String("dir", "~/.dotfiles", "The directory where the bare dotfiles repo is stored")
	syncCmd.Flags().String("backup-dir", "~/.dotfiles-backup", "Where to back up conflicting files during checkout/pull")
	syncCmd.Flags().Bool("pull", true, "Pull latest changes (fast-forward only) after checkout")
	syncCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate actions without making changes")
}
