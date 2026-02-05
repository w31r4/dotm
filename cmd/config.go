package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/w31r4/dotm/config"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage and export configuration files",
	Long: `The config command provides utilities to manage, view, export, and validate
your dotm configuration. Use this to share your configuration with others,
backup your setup, or verify configuration correctness.`,
}

var exportCmd = &cobra.Command{
	Use:     "export [destination]",
	Aliases: []string{"download"},
	Short:   "Export configuration to a file or stdout",
	Long: `Export the current config.yaml to a specified destination.
If no destination is provided, the configuration is printed to stdout.
This is useful for sharing your configuration or creating backups.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("Error loading config from %s: %v", configPath, err)
		}

		data, err := yaml.Marshal(cfg)
		if err != nil {
			log.Fatalf("Error marshaling config: %v", err)
		}

		// If no destination provided, print to stdout
		if len(args) == 0 {
			fmt.Println(string(data))
			return
		}

		destination := args[0]
		if err := os.WriteFile(destination, data, 0644); err != nil {
			log.Fatalf("Error writing config to %s: %v", destination, err)
		}
		fmt.Printf("Configuration successfully exported to: %s\n", destination)
	},
}

var showCmd = &cobra.Command{
	Use:   "show [module]",
	Short: "Show detailed configuration information",
	Long: `Display detailed information about the configuration.
If a module name is provided, shows detailed information about that specific module.
Otherwise, shows a summary of all modules and the overall configuration structure.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("Error loading config from %s: %v", configPath, err)
		}

		// Show specific module
		if len(args) == 1 {
			moduleName := args[0]
			module, ok := cfg.Modules[moduleName]
			if !ok {
				log.Fatalf("Module '%s' not found in configuration", moduleName)
			}
			showModuleDetails(moduleName, module)
			return
		}

		// Show all modules summary
		showConfigSummary(cfg)
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the configuration file",
	Long: `Validate the config.yaml file for correctness.
This checks for:
- YAML syntax errors
- Missing required fields
- Invalid module references in dependencies
- Circular dependencies`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("❌ Configuration validation failed to load %s: %v", configPath, err)
		}

		errors := validateConfig(cfg)
		if len(errors) > 0 {
			fmt.Println("❌ Configuration validation failed with the following errors:")
			for i, err := range errors {
				fmt.Printf("%d. %s\n", i+1, err)
			}
			os.Exit(1)
		}

		fmt.Println("✅ Configuration is valid!")
		fmt.Printf("Total modules: %d\n", len(cfg.Modules))
	},
}

var templateCmd = &cobra.Command{
	Use:   "template [module-name]",
	Short: "Generate a configuration template",
	Long: `Generate a template configuration for a new module.
If a module name is provided, generates a template for that specific module.
Otherwise, generates a minimal config.yaml template.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			generateModuleTemplate(args[0])
		} else {
			generateConfigTemplate()
		}
	},
}

var copyCmd = &cobra.Command{
	Use:   "copy [source] [destination]",
	Short: "Copy configuration file to another location",
	Long: `Copy the config.yaml file to another location.
This is useful for creating backups or setting up configuration in a different directory.`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		source := configPath
		destination := fmt.Sprintf("%s.backup", configPath)

		if len(args) >= 1 {
			source = args[0]
		}
		if len(args) >= 2 {
			destination = args[1]
		}

		if err := copyFile(source, destination); err != nil {
			log.Fatalf("Error copying file: %v", err)
		}
		fmt.Printf("✅ Configuration copied from %s to %s\n", source, destination)
	},
}

func showModuleDetails(name string, module config.Module) {
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("Module: %s\n", name)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("Description: %s\n", module.Description)

	if len(module.Dependencies) > 0 {
		fmt.Printf("\nDependencies:\n")
		for _, dep := range module.Dependencies {
			fmt.Printf("  - %s\n", dep)
		}
	}

	if module.Check != "" {
		fmt.Printf("\nCheck Command: %s\n", module.Check)
	}

	if len(module.Install) > 0 {
		fmt.Printf("\nInstall Commands:\n")
		for os, cmds := range module.Install {
			fmt.Printf("  %s:\n", os)
			for _, cmd := range cmds {
				fmt.Printf("    - %s\n", cmd)
			}
		}
	}

	if len(module.Apply) > 0 {
		fmt.Printf("\nPost-Install Configuration:\n")
		for _, step := range module.Apply {
			fmt.Printf("  Strategy: %s\n", step.Strategy)
			fmt.Printf("    Target: %s\n", step.Target)
			fmt.Printf("    Line: %s\n", step.Line)
		}
	}
}

func showConfigSummary(cfg *config.Config) {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Configuration Summary")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Total Modules: %d\n\n", len(cfg.Modules))

	// Sort modules by name
	keys := make([]string, 0, len(cfg.Modules))
	for k := range cfg.Modules {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Println("Available Modules:")
	fmt.Println("─────────────────────────────────────────")
	for _, name := range keys {
		module := cfg.Modules[name]
		deps := ""
		if len(module.Dependencies) > 0 {
			deps = fmt.Sprintf(" [deps: %s]", strings.Join(module.Dependencies, ", "))
		}
		fmt.Printf("  %-25s %s%s\n", name, module.Description, deps)
	}
}

func validateConfig(cfg *config.Config) []string {
	var errors []string

	// Check for empty modules
	if len(cfg.Modules) == 0 {
		errors = append(errors, "No modules defined in configuration")
		return errors
	}

	// Validate each module
	for name, module := range cfg.Modules {
		// Check for missing description
		if module.Description == "" {
			errors = append(errors, fmt.Sprintf("Module '%s' is missing a description", name))
		}

		// Check for missing install commands
		if len(module.Install) == 0 {
			errors = append(errors, fmt.Sprintf("Module '%s' has no install commands defined", name))
		}

		// Validate dependencies exist
		for _, dep := range module.Dependencies {
			if _, ok := cfg.Modules[dep]; !ok {
				errors = append(errors, fmt.Sprintf("Module '%s' depends on non-existent module '%s'", name, dep))
			}
		}

		// Validate apply steps
		for i, step := range module.Apply {
			if step.Strategy == "" {
				errors = append(errors, fmt.Sprintf("Module '%s' apply step %d is missing strategy", name, i))
			}
			if step.Target == "" {
				errors = append(errors, fmt.Sprintf("Module '%s' apply step %d is missing target", name, i))
			}
		}
	}

	// Check for circular dependencies
	for name := range cfg.Modules {
		if hasCyclicDependency(name, cfg.Modules, make(map[string]bool), make(map[string]bool)) {
			errors = append(errors, fmt.Sprintf("Circular dependency detected for module '%s'", name))
		}
	}

	return errors
}

func hasCyclicDependency(moduleName string, modules map[string]config.Module, visiting, visited map[string]bool) bool {
	if visiting[moduleName] {
		return true
	}
	if visited[moduleName] {
		return false
	}

	visiting[moduleName] = true

	module, ok := modules[moduleName]
	if ok {
		for _, dep := range module.Dependencies {
			if hasCyclicDependency(dep, modules, visiting, visited) {
				return true
			}
		}
	}

	visiting[moduleName] = false
	visited[moduleName] = true
	return false
}

func generateModuleTemplate(moduleName string) {
	template := fmt.Sprintf(`modules:
  %s:
    description: "Description of %s"
    dependencies: []
    check: "command -v %s"
    install:
      debian: ["sudo apt-get install -y %s"]
      macos: ["brew install %s"]
      default: []
    apply:
      - { strategy: "inject", target: "~/.zshrc", line: "# Configure %s here" }
`, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName)

	fmt.Println("Module Template:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Print(template)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("\nCopy this template and customize it for your needs.")
	fmt.Println("You can add it to your config.yaml using:")
	fmt.Printf("  ./dotm module add %s --description \"Your description\" ...\n", moduleName)
}

func generateConfigTemplate() {
	template := `modules:
  example-tool:
    description: "An example tool to demonstrate the configuration structure"
    dependencies: []
    check: "command -v example-tool"
    install:
      debian: ["sudo apt-get update", "sudo apt-get install -y example-tool"]
      macos: ["brew install example-tool"]
      arch: ["sudo pacman -S --noconfirm example-tool"]
      default: []
    apply:
      - { strategy: "inject", target: "~/.zshrc", line: "# Configure example-tool" }
`

	fmt.Println("Configuration Template:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Print(template)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("\nSave this template as config.yaml and customize it with your tools.")
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(exportCmd)
	configCmd.AddCommand(showCmd)
	configCmd.AddCommand(validateCmd)
	configCmd.AddCommand(templateCmd)
	configCmd.AddCommand(copyCmd)
}
