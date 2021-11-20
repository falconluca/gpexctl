package gpex

import (
	"fmt"
	u "net/url"
	"sort"
	"strings"
)

type Arranger struct {
	Conf   *Config
	Client RestClient
	vh     VideosHandler
	ch     ChannelsHandler
}

func NewArranger(conf *Config) (*Arranger, error) {
	client, err := NewRestClient(conf)
	if err != nil {
		return nil, err
	}

	a := &Arranger{
		Conf:   conf,
		Client: *client,
	}

	a.vh = a.newVideosHandler()
	a.ch = a.newChannelsHandler()

	return a, nil
}

func (a Arranger) newChannelsHandler() ChannelsHandler {
	return func(channelId string) (*Channels, error) {
		channels := NewChannelListEndpoint(a.Conf)
		channelsReq, err := channels.Request(&ChannelParams{
			ChannelId: channelId,
		})
		if err != nil {
			return nil, err
		}
		raw, err := a.Client.Do(channelsReq)
		if err != nil {
			return nil, err
		}
		resp, err := channels.Response(raw)
		if err != nil {
			return nil, err
		}
		channelsResp := resp.(Channels)
		if channelsResp.HasErrors() {
			// TODO
		}
		return &channelsResp, nil
	}
}

func (a Arranger) newVideosHandler() VideosHandler {
	return func(videoId string) (*Videos, error) {
		video := NewVideoListEndpoint(a.Conf)
		videoReq, err := video.Request(&VideoParams{
			VideoId: videoId,
		})
		if err != nil {
			return nil, err
		}
		videoRaw, err := a.Client.Do(videoReq)
		if err != nil {
			return nil, err
		}
		resp, err := video.Response(videoRaw)
		if err != nil {
			return nil, err
		}
		videosResp := resp.(Videos)
		if videosResp.HasErrors() {
			// TODO
		}
		return &videosResp, nil
	}
}

func (a Arranger) Arrange(terms []string, period int) error {
	var result []R
	for _, term := range terms {
		searchResp, err := a.DoSearch(term, period)
		if err != nil {
			return err
		}
		r := *NewR(*searchResp, a.vh, a.ch)
		result = append(result, r...)
	}
	fmt.Printf("%v 🍾 共抓取视频 %v 个\n", Indicator, GreenString(len(result)))
	sort.Slice(result, func(i, j int) bool {
		return result[i].CustomScore > result[j].CustomScore
	})

	// TODO statistics

	printResults(&result, terms)
	return nil
}

func (a Arranger) DoSearch(term string, period int) (*Search, error) {
	size := 50
	fmt.Printf("%v 🔍 正在根据关键词: %v 搜索相关视频, 预期结果: %v 个...\n", Indicator, RedString(term), size)
	search := NewSearchListEndpoint(a.Conf)
	searchReq, err := search.Request(&SearchParams{
		Size:   size,
		Term:   urlEncode(term),
		Period: period,
	})
	if err != nil {
		return nil, err
	}
	searchRaw, err := a.Client.Do(searchReq)
	if err != nil {
		return nil, err
	}
	resp, err := search.Response(searchRaw)
	if err != nil {
		return nil, err
	}
	searchResp := resp.(Search)
	if searchResp.HasErrors() {
		// TODO
	}
	fmt.Printf("%v 🍾 搜索成功, 实际结果: %v 个...\n", Indicator, RedString(len(searchResp.Items)))
	return &searchResp, nil
}

func printResults(r *[]R, terms []string) {
	keywords := strings.Join(terms, ",")
	fmt.Println("===============================")
	fmt.Printf("关键期 '%s' 最值得播放视频\n", keywords)
	fmt.Println("===============================")
	for i, rr := range *r {
		fmt.Printf("编号 #%v:\n", i+1)
		if rr.CustomScore > 0 {
			fmt.Printf("得分: %v\n", RedString(rr.CustomScore))
		}
		fmt.Printf("'%s' \n%s 有 %v 播放量, 所属频道有 %v 订阅量, 传送门 🚪: %v\n",
			CyanString(rr.Title), GreenString("> "), RedString(rr.Views), RedString(rr.NumSubscribers), rr.VideoUrl)
		fmt.Println("===============================")
	}
}

func urlEncode(term string) string {
	return u.QueryEscape(term)
}
