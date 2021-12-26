package youtube

import (
	"fmt"
	"github.com/pkg/errors"
	"math"
	"strconv"
	"time"
)

type (
	Search struct {
		DataApiResponse
		NextPageToken string `json:"nextPageToken"`
		RegionCode    string `json:"regionCode"`
		PageInfo      PageInfo
		Items         []SearchItem
	}

	SearchItem struct {
		Kind string `json:"kind"`
		Etag string `json:"etag"`
		Id   struct {
			Kind    string `json:"kind"`
			VideoId string `json:"videoId"`
		} `json:"id"`
		Snippet Snippet
	}

	Snippet struct {
		PublishedAt string `json:"publishedAt"`
		Title       string `json:"title"`
		ChannelId   string `json:"channelId"`
	}
)

func (s Search) GetTitles() []string {
	r := make([]string, 0)
	for i, item := range s.Items {
		title := item.Snippet.Title
		r = append(r, fmt.Sprintf("%v: %v;", i, title))
	}
	return r
}

func (si SearchItem) FindTitle() string {
	return si.Snippet.Title
}

func (si SearchItem) FindVideoUrl() string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%v", si.Id.VideoId)
}

func (si SearchItem) FindVideoId() string {
	return si.Id.VideoId
}

func (si SearchItem) FindViewCount(videos Videos) int {
	if len(videos.Items) == 0 {
		return 0
	}
	viewCount := videos.Items[0].
		Statistics.
		ViewCount
	videoCount, err := strconv.Atoi(viewCount)
	if err != nil {
		return 0
	}
	return videoCount
}

func (si SearchItem) FindChannelId() string {
	return si.Snippet.ChannelId
}

func (si SearchItem) FindChannelUrl() string {
	return fmt.Sprintf("https://www.youtube.com/channel/%v", si.FindChannelId())
}

func (si SearchItem) FindChannelTitle(channels *Channels) (string, error) {
	channel := channels.Items[0]
	title := channel.BrandingSettings.Channel.Title
	return title, nil
}

func (si SearchItem) FindNumSubscribers(channels Channels) int {
	channel := channels.Items[0]
	if channel.Statistics.HiddenSubscriberCount {
		return 1000000
	}

	subscriberCount := channel.Statistics.SubscriberCount
	s, err := strconv.Atoi(subscriberCount)
	if err != nil {
		return 0
	}
	return s
}

func (si SearchItem) HowOld() (float64, error) {
	publishedAt, err := time.Parse(time.RFC3339, si.Snippet.PublishedAt)
	if err != nil {
		return 0, errors.Wrap(err, "parse video published at failed")
	}

	sub := time.Now().Sub(publishedAt)
	daysSincePublished := math.Floor(sub.Hours() / 24)
	if daysSincePublished == 0 {
		daysSincePublished = 1
	}
	return daysSincePublished, nil
}
