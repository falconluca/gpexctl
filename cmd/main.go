package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gpex/cmd/command"
	"gpex/gpexctl/api"
	"gpex/gpexctl/common"
	"gpex/gpexctl/config"
	"gpex/gpexctl/xhttp"
	"os"
)

func main() {
	config.InitConfig()
	api.InitAPI()
	xhttp.InitClient()

	rootCmd := &cobra.Command{
		Use:   "gpexctl",
		Short: "Gpex is valuable YouTube videos finder",
		Long: `A tool to intentionally discover valuable videos.
More details is available at https://github.com/shaohsiung/gpex`,
	}

	rootCmd.AddCommand(command.YouTubeCmd())

	rootCmd.PersistentFlags().BoolVar(&common.Flags.Debug, "debug", false,
		"enable debug mode.")
	rootCmd.PersistentFlags().BoolVar(&common.Flags.UITable, "table", false,
		"output as UI table format.")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
