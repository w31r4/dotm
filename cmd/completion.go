package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `Generate shell completion scripts for dotm.

To load completions:

Bash:
  $ source <(dotm completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ dotm completion bash > /etc/bash_completion.d/dotm
  # macOS:
  $ dotm completion bash > /usr/local/etc/bash_completion.d/dotm

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
  # To load completions for each session, execute once:
  $ dotm completion zsh > "${fpath[1]}/_dotm"
  # You will need to start a new shell for this setup to take effect.

Fish:
  $ dotm completion fish | source
  # To load completions for each session, execute once:
  $ dotm completion fish > ~/.config/fish/completions/dotm.fish

PowerShell:
  PS> dotm completion powershell | Out-String | Invoke-Expression
  # To load completions for every new session, run:
  PS> dotm completion powershell > dotm.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
