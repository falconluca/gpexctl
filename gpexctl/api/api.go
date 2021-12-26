package api

import (
	"fmt"
	"gpex/gpexctl/common"
	"gpex/gpexctl/config"
	"log"
)

const (
	YouTubeURL = "www.googleapis.com/youtube/v3"

	YouTubeSearchURL = "/search?part=id&part=snippet&q=%v&publishedAfter=%v&maxResults=%v&order=viewCount"

	YouTubeChannelsURL = "/channels?part=brandingSettings,statistics&id=%v"

	YouTubeVideosURL = "/videos?part=statistics&id=%v"
)

type (
	YouTube struct {
		ApiKey string
	}
)

var (
	YouTubeApi *YouTube
)

func InitAPI() {
	YouTubeApi = &YouTube{
		ApiKey: config.Conf.ApiKey,
	}
}

func (y YouTube) MakeURL(urlTemplate string, a ...interface{}) string {
	url := "https://" + YouTubeURL + fmt.Sprintf(urlTemplate, a...) + "&key=" + y.ApiKey
	if common.Flags.Debug {
		log.Println(url)
	}
	return url
}
