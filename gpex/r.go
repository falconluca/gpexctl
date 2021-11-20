package gpex

import (
	"fmt"
	log "github.com/golang/glog"
	"github.com/pkg/errors"
	"math"
)

type R struct {
	Title          string
	VideoUrl       string
	Views          int
	ChannelName    string
	NumSubscribers int
	ChannelUrl     string

	ViewSubscriberRatio float64
	CustomScore         float64

	daysSincePublished float64
}

type VideosHandler func(videoId string) (*Videos, error)

type ChannelsHandler func(videoId string) (*Channels, error)

func NewR(search Search, vh VideosHandler, ch ChannelsHandler) *[]R {
	var result []R
	for i, video := range search.Items {
		videoId := video.FindVideoId()
		fmt.Printf("%s 🎬 正在抓取第 %s 个视频的详情... %s\n", Indicator, CyanString(i+1), RedString(videoId))
		videos, err := vh(videoId)
		if err != nil {
			log.Errorf("%#+v", errors.Wrap(err, "videos handler fetch videos failed"))
			continue
		}
		fmt.Printf("%s 🍻 视频的详情抓取成功, 标题: %s\n", Indicator, CyanString(video.FindTitle()))
		viewCount, err := video.FindViewCount(videos)
		if err != nil {
			log.Errorf("%#+v", err)
			continue
		}

		channelId := video.FindChannelId()
		fmt.Printf("%s 📺 正在抓取所属频道... %s\n", Indicator, RedString(channelId))
		channels, err := ch(channelId)
		if err != nil {
			log.Errorf("%#+v", errors.Wrap(err, "channels handler fetch channels failed"))
			continue
		}
		fmt.Printf("%s 🎯 所属频道抓取成功, 频道名称: %s\n", Indicator, CyanString(channels.FindChannelTitle()))
		numSubscribers, err := video.FindNumSubscribers(channels)
		if err != nil {
			log.Errorf("%#+v", err)
			continue
		}
		howOld, err := video.HowOld()
		if err != nil {
			log.Errorf("%#+v", err)
			continue
		}
		channelTitle, err := video.FindChannelTitle(channels)
		if err != nil {
			log.Errorf("%#+v", err)
			continue
		}
		r := R{
			Title:              video.FindTitle(),
			VideoUrl:           video.FindVideoUrl(),
			Views:              *viewCount,
			ChannelName:        *channelTitle,
			NumSubscribers:     *numSubscribers,
			ChannelUrl:         video.FindChannelUrl(),
			daysSincePublished: *howOld,
		}
		r.exec()
		result = append(result, r)

	}
	return &result
}

// Exec 计算视频的总得分
func (r *R) exec() {
	r.execView2SubRatio()
	r.execCustomScore()
}

func (r *R) execView2SubRatio() {
	var ratio float64
	if r.NumSubscribers == 0 {
		ratio = 0
	} else {
		ratio = float64(r.Views / r.NumSubscribers)
	}
	r.ViewSubscriberRatio = ratio
}

func (r *R) execCustomScore() {
	ratio := math.Min(r.ViewSubscriberRatio, 5)
	r.CustomScore = (float64(r.Views) * ratio) / r.daysSincePublished
}
