package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "gpex",
	Short: "Gpex is valuable YouTube videos finder",
	Long: `A tool to intentionally discover valuable videos.
More details is available at https://github.com/shaohsiung/gpex`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
