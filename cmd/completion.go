package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// localCmd represents the local command
func completionCmd() *cobra.Command {
	// app := app.New()
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:

$ source <(binenv completion bash)

# To load completions for each session, execute once:
Linux:
  $ binenv completion bash > /etc/bash_completion.d/binenv
MacOS:
  $ binenv completion bash > /usr/local/etc/bash_completion.d/binenv

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ binenv completion zsh > "${fpath[1]}/_binenv"

# You will need to start a new shell for this setup to take effect.

Fish:

$ binenv completion fish | source

# To load completions for each session, execute once:
$ binenv completion fish > ~/.config/fish/completions/binenv.fish
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletion(os.Stdout)
			}
		},
	}
	// cmd.Flags().IntVarP(&a.Params.MinLength, "min-length", "m", 16, "Specify minimum password length, must not be less than 8")
	return cmd
}
