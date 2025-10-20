package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the entire configuration loaded from config.yaml
type Config struct {
	Modules map[string]Module `yaml:"modules"`
}

// Module represents a single installable unit (e.g., zsh, fzf).
type Module struct {
	Description  string              `yaml:"description"`
	Dependencies []string            `yaml:"dependencies"`
	Check        string              `yaml:"check"`
	Install      map[string][]string `yaml:"install"`
	Apply        []ApplyStep         `yaml:"apply"`
}

// ApplyStep defines a single action to configure a dotfile.
type ApplyStep struct {
	Strategy string `yaml:"strategy"`
	Target   string `yaml:"target"`
	Line     string `yaml:"line"`
}

// LoadConfig reads and parses the config.yaml file from the given path.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
