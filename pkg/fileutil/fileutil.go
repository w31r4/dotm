package fileutil

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// expandHome resolves the '~' character to the user's home directory.
func expandHome(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return strings.Replace(path, "~", home, 1), nil
	}
	return path, nil
}

// InjectLine ensures a specific line is present in a file.
// If the line already exists, it does nothing. Otherwise, it appends the line.
func InjectLine(filePath, lineToInject string) error {
	expandedPath, err := expandHome(filePath)
	if err != nil {
		return fmt.Errorf("could not expand home directory in path '%s': %w", filePath, err)
	}

	// Ensure the file exists, create if it doesn't
	file, err := os.OpenFile(expandedPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("could not open or create file '%s': %w", expandedPath, err)
	}
	defer file.Close()

	// Check if the line already exists
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == strings.TrimSpace(lineToInject) {
			fmt.Printf("Line already exists in %s. Skipping.\n", expandedPath)
			return nil
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning file '%s': %w", expandedPath, err)
	}

	// If we reached here, the line doesn't exist, so append it.
	// We add a newline before our line for better formatting if the file is not empty.
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	if stat.Size() > 0 {
		if _, err := file.WriteString("\n"); err != nil {
			return fmt.Errorf("failed to write newline to '%s': %w", expandedPath, err)
		}
	}

	if _, err := file.WriteString(lineToInject + "\n"); err != nil {
		return fmt.Errorf("failed to write to '%s': %w", expandedPath, err)
	}

	fmt.Printf("Successfully injected line into %s\n", expandedPath)
	return nil
}
