package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "eic",
	Short: "Ensure import comment",
	Run:   rootRun,
}

var (
	_dirpath  string
	_filepath string
	_dryrun   bool
)

func init() {
	rootCmd.Flags().StringVarP(&_dirpath, "dir", "d", "", "transfer directory")
	rootCmd.Flags().StringVarP(&_filepath, "file", "f", "", "transfer a file")
	rootCmd.Flags().BoolVarP(&_dryrun, "dryrun", "n", false, "show what would have been transferred")
}

func rootRun(cmd *cobra.Command, args []string) {
	w := &Worker{
		DryRun: _dryrun,
	}

	switch {
	case _dirpath != "":
		if err := w.WorkDir(_dirpath); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case _filepath != "":
		if err := w.WorkFile(_filepath); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		cmd.Help()
		os.Exit(1)

	}

}
