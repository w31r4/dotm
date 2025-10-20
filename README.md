# dotm - A Modern Dotfile Manager

`dotm` is a command-line tool for bootstrapping and managing your development environment. It replaces complex shell scripts with a simple, configuration-driven approach powered by Go.

## Core Concepts

- **Configuration as Code**: All setup logic is defined in a clear `config.yaml` file.
- **Modular**: Each piece of software (zsh, git, fzf) is a self-contained "module".
- **Idempotent**: Running the tool multiple times won't break your system. It checks before it acts.
- **Extensible**: Easily add new software to your setup by defining a new module in the config.
- **Powered by x-cmd**: Leverages the `x-cmd` universal package manager for simplified tool installation where possible.

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

> **Note**: This initial version does not automatically handle conflicts with existing files (e.g., a default `.bashrc`). You may need to back them up manually if the checkout fails.

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