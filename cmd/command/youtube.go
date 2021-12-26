package command

import (
	"encoding/json"
	"fmt"
	log "github.com/golang/glog"
	"github.com/gosuri/uitable"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gpex/gpexctl/api"
	"gpex/gpexctl/biz"
	"gpex/gpexctl/common"
	"gpex/gpexctl/config"
	"gpex/gpexctl/xhttp"
	"gpex/gpexctl/youtube"
	"net/http"
	"sort"
	"strings"
	"time"
)

type YouTubeFlags struct {
	Terms      []string
	Period     int
	MaxResults int
}

var (
	flags YouTubeFlags
)

func YouTubeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "youtube",
		Aliases: []string{"y", "yt"},
		Short:   "YouTube abilities",
	}

	cmd.AddCommand(searchYouTubeCmd())
	return cmd
}

func searchYouTubeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "YouTube finder",
		Run: func(cmd *cobra.Command, args []string) {
			if len(flags.Terms) == 0 {
				common.ExitWithError(errors.New("term at least one"))
			}

			var publishedAfter string
			if period, err := resolveSearchPeriod(flags.Period); err != nil {
				common.ExitWithError(err)
			} else {
				publishedAfter = period.Format(time.RFC3339)
			}

			var result []biz.CustomScore
			for _, q := range flags.Terms {
				url := api.YouTubeApi.MakeURL(api.YouTubeSearchURL, q, publishedAfter, flags.MaxResults)
				body := xhttp.Client.HandleRequest(http.MethodGet, url, nil)

				var searchResult youtube.Search
				if err := json.Unmarshal(body, &searchResult); err != nil {
					log.Errorf("%#+v", err)
				}

				res := biz.NewYouTubeBizWithSearch(searchResult)
				customScoreList := res.CustomScoreList()
				result = append(result, customScoreList...)
			}
			fmt.Printf("%v 🍾 共抓取视频 %v 个\n", config.Indicator, config.GreenString(len(result)))
			sort.Slice(result, func(i, j int) bool {
				return result[i].CustomScore > result[j].CustomScore
			})

			if common.Flags.UITable {
				printUITable(result)
			} else {
				printBody(result)
			}
		},
	}

	cmd.PersistentFlags().StringArrayVarP(&flags.Terms, "term", "q", []string{},
		"search keywords for videos.")
	cmd.PersistentFlags().IntVarP(&flags.Period, "period", "p", 10,
		"how long are videos from today.")
	cmd.PersistentFlags().IntVarP(&flags.MaxResults, "size", "s", 10,
		"the size of videos.")

	return cmd
}

func resolveSearchPeriod(period int) (*time.Time, error) {
	date := time.Now().AddDate(0, 0, -period)
	date, err := time.Parse("2006-01-02 15:04:05", date.Format("2006-01-02 15:04:05"))
	if err != nil {
		return nil, err
	}
	return &date, nil
}

func printUITable(cs []biz.CustomScore) {
	table := uitable.New()
	table.MaxColWidth = 50

	table.AddRow("编号", "得分", "标题", "播放量", "订阅量", "🚪 传送门")

	keywords := strings.Join(flags.Terms, ",")
	fmt.Println("===============================")
	fmt.Printf("关键期 '%s' 最值得播放视频\n", keywords)
	fmt.Println("===============================")
	for i, rr := range cs {
		var customScoreString string
		if rr.CustomScore > 0 {
			customScoreString = config.RedString(rr.CustomScore)
		} else {
			customScoreString = "-"
		}
		table.AddRow(fmt.Sprintf("No.%d", i+1), customScoreString, config.CyanString(rr.Title),
			config.RedString(rr.Views), config.RedString(rr.NumSubscribers), rr.VideoURL)
	}

	fmt.Println(table)
}

func printBody(cs []biz.CustomScore) {
	keywords := strings.Join(flags.Terms, ",")
	fmt.Println("===============================")
	fmt.Printf("关键期 '%s' 最值得播放视频\n", keywords)
	fmt.Println("===============================")
	for i, rr := range cs {
		fmt.Printf("编号 #%v:\n", i+1)
		if rr.CustomScore > 0 {
			fmt.Printf("得分: %v\n", config.RedString(rr.CustomScore))
		}
		fmt.Printf("'%s' \n%s 有 %v 播放量, 所属频道有 %v 订阅量, 传送门 🚪: %v\n",
			config.CyanString(rr.Title), config.GreenString("> "),
			config.RedString(rr.Views), config.RedString(rr.NumSubscribers), rr.VideoURL)
		fmt.Println("===============================")
	}
}
