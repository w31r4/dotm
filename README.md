# dotm - A Modern Dotfile Manager

`dotm` is a command-line tool for bootstrapping and managing your development environment. It replaces complex shell scripts with a simple, configuration-driven approach powered by Go.

## Core Concepts

- **Configuration as Code**: All setup logic is defined in a clear `config.yaml` file.
- **Modular**: Each piece of software (zsh, git, fzf) is a self-contained "module".
- **Idempotent**: Running the tool multiple times won't break your system. It checks before it acts.
- **Extensible**: Easily add new software to your setup by defining a new module in the config.
- **Powered by x-cmd**: Leverages the `x-cmd` universal package manager for simplified tool installation where possible.
- **Configuration Management**: Export, validate, and share your configurations with ease.

## Quick Start

### 1. Installation

To use `dotm`, you need to have Go installed on your system.

```bash
# Clone the repository (or your dotfiles repo containing this tool)
# ...

# Build the binary
cd scripts/dotm
go build
```

This will create a `dotm` executable in the current directory. You can move this to a location in your `$PATH` for global access, e.g., `sudo mv dotm /usr/local/bin/`.

### 2. Bootstrapping a New Machine

The typical workflow for setting up a new machine is a two-step process:

**Step 1: Sync your dotfiles repository**

This command clones your bare dotfiles repository to `~/.dotfiles` and checks out the files into your home directory.

```bash
./dotm repo sync --url git@github.com:your-username/your-dotfiles.git
```

> **Note**: If checkout/pull would overwrite existing files (e.g., a default `.bashrc`), `dotm` will back them up to `~/.dotfiles-backup/<timestamp>` and retry.

Useful flags:
- `--dir` to change the bare repo location (default `~/.dotfiles`)
- `--backup-dir` to change where conflicts are backed up
- `--pull=false` to skip pulling latest changes

**Step 2: On-Demand Installation**

`config.yaml` acts as your personal software repository. You can install any module from it on demand.

```bash
# Install a single module
./dotm install zsh

# Install multiple modules at once
./dotm install fzf pyenv eza

# The tool will automatically handle dependencies for you.
```

### 3. Safe Preview with Dry Run

To see what commands `dotm` *would* execute without actually changing anything, use the `--dry-run` flag. This is highly recommended before running on a new system.

```bash
./dotm install --dry-run eza
```

### Managing the Dotfiles Repo

If you use the bare-repo workflow, you can run git commands via `dotm` without setting up a separate shell alias:

```bash
./dotm repo git status
./dotm repo git add .gitconfig
./dotm repo git commit -m "Track gitconfig"
./dotm repo git push
```

## Managing Your Module "Repository"

Manually editing `config.yaml` can be tedious. `dotm` provides a suite of commands to manage your modules efficiently.

### List Modules

To see all the modules available in your `config.yaml`:

```bash
./dotm module list
```

### Add a New Module

To add a new module, use the `module add` command with flags.

**Example:** Adding `htop`

```bash
./dotm module add htop \
  --description "Interactive process viewer" \
  --check "command -v htop" \
  --install-debian "sudo apt-get install -y htop" \
  --install-macos "brew install htop"
```
This command will safely and correctly append the `htop` module to your `config.yaml`.

### Remove a Module

To remove a module you no longer need:

```bash
./dotm module remove htop
```

## Managing Configuration

The `config` command provides powerful utilities for managing, exporting, and validating your configuration files.

### Export/Download Configuration

Export your configuration to share with others or create a backup:

```bash
# Export to stdout
./dotm config export

# Export to a specific file
./dotm config export my-dotfiles-config.yaml

# You can also use the 'download' alias
./dotm config download backup.yaml
```

### View Configuration

Display detailed information about your modules:

```bash
# Show summary of all modules
./dotm config show

# Show details of a specific module
./dotm config show fzf
```

### Validate Configuration

Check your configuration file for errors:

```bash
./dotm config validate
```

This validates:
- YAML syntax
- Missing required fields
- Invalid module references in dependencies
- Circular dependencies

### Generate Templates

Generate configuration templates for new modules:

```bash
# Generate a template for a specific module
./dotm config template htop

# Generate a minimal config.yaml template
./dotm config template
```

### Copy Configuration

Create backups of your configuration:

```bash
# Copy with default names
./dotm config copy

# Specify source and destination
./dotm config copy config.yaml config.yaml.backup
```

## Shell Completion

Enable shell completions for faster workflows:

```bash
# Bash (current session)
source <(./dotm completion bash)

# Zsh (permanent)
./dotm completion zsh > "${fpath[1]}/_dotm"

# Fish (current session)
./dotm completion fish | source
```

## Configuration File (`config.yaml`)

The heart of `dotm` is the `config.yaml` file. Here's a breakdown of the structure:

```yaml
modules:
  # A module name
  fzf:
    # Description for humans
    description: "A command-line fuzzy finder"
    # Other modules that must be installed first
    dependencies: [x-cmd]
    # A shell command to check if the module is already installed.
    # If it exits with 0 (success), installation is skipped.
    check: "command -v fzf"
    # A map of OS to a list of installation commands.
    # 'default' is used if the specific OS is not found.
    install:
      default: ["x env use fzf"]
    # Post-installation steps, like configuring a dotfile.
    apply:
      - { strategy: "inject", target: "~/.zshrc", line: "source /path/to/fzf.zsh" }
```

This declarative approach makes it incredibly easy to see, modify, and extend your entire environment setup from a single file.
