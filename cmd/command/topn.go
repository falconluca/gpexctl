package command

import "github.com/spf13/cobra"

func TopNCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "top",
		Aliases: []string{"t"},
		Short:   "Top N for Weibo, Baidu...",
	}
	cmd.AddCommand(listAllTopN())
	cmd.AddCommand(getWeiboTopN())
	cmd.AddCommand(getBaiduTopN())
	return cmd
}

func listAllTopN() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all",
		Short: "get all top n",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	return cmd
}

func getWeiboTopN() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "weibo",
		Short: "get weibo top n",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	return cmd
}

func getBaiduTopN() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "baidu",
		Short: "get baidu top n",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	return cmd
}
