package biz

import (
	"encoding/json"
	"fmt"
	log "github.com/golang/glog"
	"gpex/gpexctl/api"
	"gpex/gpexctl/xhttp"
	"gpex/gpexctl/youtube"
	"math"
	"net/http"
	"time"
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
		cch := make(chan *youtube.Channels)
		go func() {
			var channelsResult *youtube.Channels
			url := api.YouTubeApi.MakeURL(api.YouTubeChannelsURL, video.FindChannelId())
			body := xhttp.Client.HandleRequest(http.MethodGet, url, nil)
			if err := json.Unmarshal(body, &channelsResult); err != nil {
				log.Errorf("%#+v", err)
				cch <- nil
			}
			cch <- channelsResult
		}()

		vch := make(chan *youtube.Videos)
		go func() {
			var videosResult *youtube.Videos
			url := api.YouTubeApi.MakeURL(api.YouTubeVideosURL, video.FindVideoId())
			body := xhttp.Client.HandleRequest(http.MethodGet, url, nil)
			if err := json.Unmarshal(body, &videosResult); err != nil {
				log.Errorf("%#+v", err)
				vch <- nil
			}
			vch <- videosResult
		}()

		viewCount := 0
		numSubscribers := 0
		for i := 0; i < 2; i++ {
			select {
			case channelsResult := <-cch:
				numSubscribers = video.FindNumSubscribers(*channelsResult)
			case videosResult := <-vch:
				viewCount = video.FindViewCount(*videosResult)
			case <-time.After(3 * time.Second):
				fmt.Println("fetch channel or video time out")
			}
		}

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
