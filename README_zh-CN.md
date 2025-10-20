# dotm - 现代化的 Dotfile 管理器

`dotm` 是一个用于引导和管理您的开发环境的命令行工具。它用一个由 Go 驱动的、简单的配置驱动方法，替代了复杂的 Shell 脚本。

## 核心理念

- **配置即代码**: 所有的安装和设置逻辑都定义在一个清晰的 `config.yaml` 文件中。
- **模块化**: 每一款软件（如 zsh, git, fzf）都是一个独立的“模块”。
- **幂等性**: 多次运行本工具不会对您的系统造成任何损害。它在执行操作前会进行检查。
- **可扩展**: 通过在配置中定义一个新模块，您可以轻松地为您的环境添加新软件。
- **x-cmd 驱动**: 在可能的情况下，利用 `x-cmd` 通用包管理器来简化工具的安装。

## 快速开始

### 1. 安装

要使用 `dotm`，您的系统中需要先安装 Go 环境。

```bash
# 克隆包含本工具的仓库（例如您的 dotfiles 仓库）
# ...

# 编译二进制文件
cd scripts/dotm
go build
```

这会在当前目录下创建一个名为 `dotm` 的可执行文件。您可以将此文件移动到您系统 `$PATH` 的某个路径下（例如 `/usr/local/bin/`）以便全局访问。

### 2. 在新机器上引导环境

在新机器上进行设置的典型工作流程分为两步：

**第一步：同步您的 dotfiles 仓库**

此命令会将您的 dotfiles 裸仓库克隆到 `~/.dotfiles`，并检出（checkout）其中的文件到您的 Home 目录。

```bash
./dotm repo sync --url git@github.com:your-username/your-dotfiles.git
```

> **注意**：此初始版本不会自动处理与现有文件（例如系统默认的 `.bashrc`）的冲突。如果检出失败，您可能需要手动备份这些文件。

**第二步：按需安装您的工具**

`config.yaml` 文件扮演着您的个人软件仓库的角色。您可以按需安装其中的任何模块。

```bash
# 安装单个模块
./dotm install zsh

# 一次性安装多个模块
./dotm install fzf pyenv eza

# 工具会自动为您处理依赖关系。
```

### 3. 使用“Dry Run”安全预览

如果您想查看 `dotm` *将要* 执行哪些命令，而不想实际对系统做出任何更改，请使用 `--dry-run` 标志。强烈建议在新系统上运行时首先使用此功能。

```bash
./dotm install --dry-run eza
```

## 管理您的模块“仓库”

手动编辑 `config.yaml` 可能既繁琐又容易出错。`dotm` 提供了一套命令来帮助您高效地管理模块。

### 列出模块

查看 `config.yaml` 中所有可用的模块：

```bash
./dotm module list
```

### 添加新模块

使用 `module add` 命令和相关标志来添加一个新模块。

**示例：** 添加 `htop`

```bash
./dotm module add htop \
  --description "一个交互式的进程查看器" \
  --check "command -v htop" \
  --install-debian "sudo apt-get install -y htop" \
  --install-macos "brew install htop"
```
此命令会自动且正确地将 `htop` 模块追加到您的 `config.yaml` 文件中。

### 移除模块

移除一个您不再需要的模块：

```bash
./dotm module remove htop
```

## 配置文件 (`config.yaml`)

`dotm` 的核心是 `config.yaml` 文件。以下是其结构解析：

```yaml
modules:
  # 模块名
  fzf:
    # 描述信息
    description: "一个命令行的模糊查找工具"
    # 此模块依赖的其他模块
    dependencies: [x-cmd]
    # 用于检查此模块是否已安装的 Shell 命令。
    # 如果此命令以 0 状态码（成功）退出，则跳过安装。
    check: "command -v fzf"
    # 一个从操作系统到安装命令列表的映射。
    # 如果找不到特定系统的命令，则使用 'default'。
    install:
      default: ["x env use fzf"]
    # 安装后的配置步骤，例如向 dotfile 中注入内容。
    apply:
      - { strategy: "inject", target: "~/.zshrc", line: "source /path/to/fzf.zsh" }
```

这种声明式的方法让您可以从单一文件中轻松地查看、修改和扩展您的整个环境配置。