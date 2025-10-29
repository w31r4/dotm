# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2025-10-29

### Added

- **Configuration Management Commands** (`config` command)
  - `config export` - Export configuration to a file or stdout (supports `download` alias)
  - `config show` - Display detailed configuration information for all modules or specific module
  - `config validate` - Validate configuration file for syntax errors, missing fields, invalid dependencies, and circular dependencies
  - `config template` - Generate configuration templates for new modules or entire config.yaml
  - `config copy` - Copy configuration file to another location for backups

- **Version Command**
  - `version` - Display the current version of dotm

- **Shell Completion Support**
  - `completion` - Generate shell completion scripts for bash, zsh, fish, and powershell
  - Enables tab completion for commands and flags across all supported shells

- **Global Configuration Path Flag**
  - `--config` flag now available on all commands to specify custom config file location
  - Default path remains `./config.yaml`

### Changed

- Enhanced root command description with clearer messaging about dotm's purpose
- All commands now use the configurable `configPath` variable instead of hardcoded "config.yaml"
- Improved error messages to show which config file failed to load

### Documentation

- Updated README.md with comprehensive documentation for new `config` command
- Updated README_zh-CN.md with Chinese documentation for new features
- Added new "Managing Configuration" section to both README files
- Updated "Core Concepts" section to highlight configuration management capabilities
- Created CHANGELOG.md to track version history

## [0.0.1] - Initial Release

### Added

- Core module installation system with dependency resolution
- `install` command for installing modules
- `module list` command for listing available modules
- `module add` command for adding new modules
- `module remove` command for removing modules
- `repo sync` command for bootstrapping dotfiles repository
- Dry-run support for safe preview of installation commands
- OS-specific command selection (debian, macos, arch, default)
- Idempotent installation with check commands
- Post-install configuration with inject strategy
- Sample configuration with common development tools
