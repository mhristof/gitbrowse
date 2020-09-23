package cmd

import (
	"fmt"
	"runtime"

	"github.com/mhristof/gitbrowse/log"
	"github.com/mhristof/go-update"
	"github.com/spf13/cobra"
)

var (
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update the binary with a new version",
		Run: func(cmd *cobra.Command, args []string) {
			url := fmt.Sprintf("https://github.com/mhristof/gitbrowse/releases/latest/download/gitbrowse.%s", runtime.GOOS)
			updates, updateFunc, err := update.Check(url)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("Cannot Check for an update")

			}

			if !updates {
				return
			}

			if silent, _ := cmd.Flags().GetBool("silent"); !silent {
				fmt.Println("New version downloaded!")
			}

			if dryrun, _ := cmd.Flags().GetBool("dryrun"); dryrun {
				return
			}

			updateFunc()
		},
	}
)

func init() {
	rootCmd.AddCommand(updateCmd)
}
