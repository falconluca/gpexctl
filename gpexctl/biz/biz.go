package biz

import (
	"encoding/json"
	log "github.com/golang/glog"
	"gpex/gpexctl/api"
	"gpex/gpexctl/xhttp"
	"gpex/gpexctl/youtube"
	"math"
	"net/http"
)

type (
	YouTubeBiz struct {
		Search youtube.Search
	}

	CustomScore struct {
		CustomScore    float64
		Title          string
		VideoURL       string
		Views          int
		NumSubscribers int
	}
)

func NewYouTubeBizWithSearch(search youtube.Search) *YouTubeBiz {
	return &YouTubeBiz{
		Search: search,
	}
}

func (b YouTubeBiz) CustomScoreList() []CustomScore {
	var cs []CustomScore
	for _, video := range b.Search.Items {
		url := api.YouTubeApi.MakeURL(api.YouTubeChannelsURL, video.FindChannelId())
		body := xhttp.Client.HandleRequest(http.MethodGet, url, nil)
		var channelsResult *youtube.Channels
		if err := json.Unmarshal(body, &channelsResult); err != nil {
			log.Errorf("%#+v", err)
		}
		numSubscribers := video.FindNumSubscribers(*channelsResult)

		url = api.YouTubeApi.MakeURL(api.YouTubeVideosURL, video.FindVideoId())
		body = xhttp.Client.HandleRequest(http.MethodGet, url, nil)
		var videosResult *youtube.Videos
		if err := json.Unmarshal(body, &videosResult); err != nil {
			log.Errorf("%#+v", err)
		}
		viewCount := video.FindViewCount(*videosResult)

		ratio := viewSubscriberRatio(viewCount, numSubscribers)
		howOld, err := video.HowOld()
		if err != nil {
			log.Errorf("%#+v", err)
		}

		cs = append(cs, CustomScore{
			Title:          video.FindTitle(),
			VideoURL:       video.FindVideoUrl(),
			Views:          viewCount,
			NumSubscribers: numSubscribers,
			CustomScore:    customScore(ratio, howOld, viewCount),
		})
	}
	return cs
}

func viewSubscriberRatio(viewCount int, numSubscribers int) float64 {
	if numSubscribers == 0 {
		return 0
	}
	return float64(viewCount / numSubscribers)
}

func customScore(viewSubscriberRatio float64, howOld float64, viewCount int) float64 {
	ratio := math.Min(viewSubscriberRatio, 5)
	return (float64(viewCount) * ratio) / howOld
}
