package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	"github.com/mhristof/gitbrowse/log"
	"github.com/riywo/loginshell"
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "autocomplete",
	Short: "Generate completion script",
	Long: heredoc.Doc(fmt.Sprintf(`
		For zsh, you need to run this command

		%s completion zsh > "${fpath[1]}/_%s"
		echo 'source ${fpath[1]}/_%s && compdef _%s %s' >> ~/.zshrc
	`, rootCmd.Use, rootCmd.Use, rootCmd.Use, rootCmd.Use, rootCmd.Use)),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		shell, err := loginshell.Shell()
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Cannot retrieve login shell")

		}

		switch filepath.Base(shell) {
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

func init() {
	rootCmd.AddCommand(completionCmd)
}
