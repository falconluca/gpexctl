package cmd

import (
	"fmt"
	log "github.com/golang/glog"
	"gpex/gpex"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		terms, err := cmd.Flags().GetStringArray("terms")
		if err != nil {
			log.Errorf("%#+v", err)
			return
		}
		if len(terms) == 0 {
			fmt.Println("terms flag is required")
			return
		}
		period, err := cmd.Flags().GetInt("period")
		if err != nil {
			log.Errorf("%#+v", err)
			return
		}

		conf, err := gpex.LoadConfig("config.yaml")
		if err != nil {
			log.Errorf("%#+v", err)
		}
		arranger, err := gpex.NewArranger(conf)
		if err != nil {
			log.Errorf("%#+v", err)
		}
		if err = arranger.Arrange(terms, period); err != nil {
			log.Errorf("%#+v", err)
		}
	},
}

func init() {
	searchCmd.Flags().IntP("period", "p", 3, "published at")
	searchCmd.Flags().StringArrayP("terms", "t", []string{""}, "What video your want to find")
	rootCmd.AddCommand(searchCmd)
}
