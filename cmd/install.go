package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/w31r4/dotm/config"
	"github.com/w31r4/dotm/pkg/executor"
	"github.com/w31r4/dotm/pkg/fileutil"
)

var dryRun bool
var installedModules = make(map[string]bool)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [module...]",
	Short: "Install and configure one or more modules",
	Long: `Install modules defined in the config.yaml file.
This command will check for dependencies, install the required software,
and apply the necessary dotfile configurations.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("Error loading config from %s: %v", configPath, err)
		}

		for _, moduleName := range args {
			if err := installModule(moduleName, cfg, dryRun); err != nil {
				log.Fatalf("Failed to install module %s: %v", moduleName, err)
			}
		}
		fmt.Println("\nAll requested modules installed successfully!")
	},
}

func installModule(name string, cfg *config.Config, dryRun bool) error {
	if installedModules[name] {
		return nil // Already installed in this run
	}

	fmt.Printf("--- Installing module: %s ---\n", name)

	module, ok := cfg.Modules[name]
	if !ok {
		return fmt.Errorf("module '%s' not found in config.yaml", name)
	}

	// 1. Handle dependencies first
	if len(module.Dependencies) > 0 {
		fmt.Println("Checking dependencies...")
		for _, depName := range module.Dependencies {
			if err := installModule(depName, cfg, dryRun); err != nil {
				return fmt.Errorf("dependency '%s' for module '%s' failed to install: %w", depName, name, err)
			}
		}
	}

	// 2. Check if the module is already installed
	if module.Check != "" {
		fmt.Printf("Running check: %s\n", module.Check)
		if err := executor.Execute(module.Check, dryRun); err == nil {
			fmt.Println("Module is already installed. Skipping installation.")
			installedModules[name] = true
			// Even if installed, we might want to re-apply configs
			return applyConfiguration(module, dryRun)
		}
		fmt.Println("Module not found, proceeding with installation.")
	}

	// 3. Install the software
	os := executor.GetOS()
	installCmds, ok := module.Install[os]
	if !ok {
		// If no specific command for the OS, check for a "default"
		installCmds, ok = module.Install["default"]
		if !ok {
			return fmt.Errorf("no install command found for OS '%s' or default in module '%s'", os, name)
		}
	}

	fmt.Printf("Running install commands for %s...\n", name)
	for _, cmd := range installCmds {
		if err := executor.Execute(cmd, dryRun); err != nil {
			return fmt.Errorf("installation command '%s' failed: %w", cmd, err)
		}
	}

	// 4. Apply dotfile configurations
	if err := applyConfiguration(module, dryRun); err != nil {
		return err
	}

	installedModules[name] = true
	fmt.Printf("--- Successfully installed module: %s ---\n", name)
	return nil
}

func applyConfiguration(module config.Module, dryRun bool) error {
	if len(module.Apply) > 0 {
		fmt.Println("Applying configurations...")
		for _, step := range module.Apply {
			switch step.Strategy {
			case "inject":
				if err := fileutil.InjectLine(step.Target, step.Line, dryRun); err != nil {
					return fmt.Errorf("failed to apply inject strategy on '%s': %w", step.Target, err)
				}
			default:
				return fmt.Errorf("unknown apply strategy: '%s'", step.Strategy)
			}
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate the installation without making any changes")
}
