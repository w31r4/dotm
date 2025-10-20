package executor

import (
	"bufio"
	"fmt"
	"os/exec"
	"runtime"
)

// GetOS returns the current operating system (e.g., "linux", "darwin").
// We'll use "debian" for linux as a simplification for this project.
func GetOS() string {
	switch runtime.GOOS {
	case "linux":
		// For this project, we'll assume Debian-based Linux.
		// A more robust solution would check /etc/os-release.
		return "debian"
	case "darwin":
		return "macos"
	default:
		return runtime.GOOS
	}
}

// Execute runs a command and streams its output to stdout.
func Execute(command string) error {
	fmt.Printf("Executing: %s\n", command)

	// Use sh -c to properly handle commands with pipes or multiple parts.
	cmd := exec.Command("sh", "-c", command)

	// Get pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("could not get stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("could not get stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("could not start command '%s': %w", command, err)
	}

	// Create scanners to read output line by line
	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)

	// Concurrently read from both pipes
	go func() {
		for stdoutScanner.Scan() {
			fmt.Println(stdoutScanner.Text())
		}
	}()
	go func() {
		for stderrScanner.Scan() {
			fmt.Println(stderrScanner.Text())
		}
	}()

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command '%s' failed: %w", command, err)
	}

	return nil
}
