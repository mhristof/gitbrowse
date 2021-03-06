package cmd

import (
	"fmt"
	"os"

	"github.com/mhristof/gitbrowse/git"
	"github.com/mhristof/gitbrowse/log"
	"github.com/spf13/cobra"
)

var version = "devel"

var rootCmd = &cobra.Command{
	Use:     "gitbrowse",
	Short:   "Translate local git repositories to URLs",
	Args:    cobra.ExactArgs(1),
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		Verbose(cmd)

		if _, err := os.Stat(args[0]); os.IsNotExist(err) {
			log.WithFields(log.Fields{
				"args[0]": args[0],
			}).Error("Does not exist")
		}

		repo, err := git.New(args[0])
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Fatal("Cant create a repo")
		}

		line, err := cmd.Flags().GetInt("line")
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Panic("Cannot retrieve line arg")

		}

		url, err := repo.URL(args[0], line)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Fatal("Cannot calculate url")

		}
		fmt.Println(url)

	},
}

// Verbose Increase verbosity
func Verbose(cmd *cobra.Command) {
	verbose, err := cmd.Flags().GetBool("verbose")

	if err != nil {
		log.Panic(err)
	}

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Increase verbosity")
	rootCmd.PersistentFlags().BoolP("dryrun", "n", false, "Dry run")
	rootCmd.Flags().IntP("line", "l", -1, "Line number")
}

// Execute The main function for the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
