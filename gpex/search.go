package gpex

import (
	"encoding/json"
	"fmt"
	log "github.com/golang/glog"
	"github.com/pkg/errors"
	"math"
	"net/http"
	"strconv"
	"time"
)

type Search struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	RegionCode    string `json:"regionCode"`
	PageInfo      PageInfo
	Items         []SearchItem
	Error         YouTubeDataApiError `json:"error"`
}

func (s Search) Meta() string {
	return fmt.Sprintf("kind: %v, etag: %v", s.Kind, s.Etag)
}

func (s Search) HasErrors() bool {
	return len(s.Error.Errors) > 0
}

// GetTitles just for test purpose
func (s Search) GetTitles() *[]string {
	r := make([]string, 0)
	for i, item := range s.Items {
		title := item.Snippet.Title
		r = append(r, fmt.Sprintf("%v: %v;", i, title))
	}
	return &r
}

type SearchItem struct {
	Kind string `json:"kind"`
	Etag string `json:"etag"`
	Id   struct {
		Kind    string `json:"kind"`
		VideoId string `json:"videoId"`
	} `json:"id"`
	Snippet Snippet
}

// FindTitle 获取视频的标题
func (item SearchItem) FindTitle() string {
	return item.Snippet.Title
}

// FindVideoUrl 获取视频的播放链接
func (item SearchItem) FindVideoUrl() string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%v", item.Id.VideoId)
}

// FindVideoId 获取视频ID
func (item SearchItem) FindVideoId() string {
	return item.Id.VideoId
}

// FindViewCount 获取视频的播放量
func (item SearchItem) FindViewCount(videos *Videos) (*int, error) {
	viewCount := videos.Items[0]. // FIXME panic: runtime error: index out of range [0] with length 0
					Statistics.
					ViewCount
	videoCount, err := strconv.Atoi(viewCount)
	if err != nil {
		return nil, err
	}
	return &videoCount, nil
}

// FindChannelId 获取视频所属频道的ID
func (item SearchItem) FindChannelId() string {
	return item.Snippet.ChannelId
}

// FindChannelUrl 获取视频所属频道的链接
func (item SearchItem) FindChannelUrl() string {
	return fmt.Sprintf("https://www.youtube.com/channel/%v", item.FindChannelId())
}

// FindChannelTitle 获取视频所属频道的名称
func (item SearchItem) FindChannelTitle(channels *Channels) (*string, error) {
	channel := channels.Items[0]
	title := channel.BrandingSettings.Channel.Title
	return &title, nil
}

// FindNumSubscribers 获取视频所属频道的订阅数
func (item SearchItem) FindNumSubscribers(channels *Channels) (*int, error) {
	channel := channels.Items[0]
	if channel.Statistics.HiddenSubscriberCount {
		s := 1000000
		return &s, nil
	} else {
		subscriberCount := channel.Statistics.SubscriberCount
		s, err := strconv.Atoi(subscriberCount)
		if err != nil {
			return nil, err
		}
		return &s, nil
	}
}

// HowOld 获取视频发布距现在几天
func (item SearchItem) HowOld() (*float64, error) {
	publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
	if err != nil {
		return nil, errors.Wrap(err, "parse video published at failed")
	}

	sub := time.Now().Sub(publishedAt)
	daysSincePublished := math.Floor(sub.Hours() / 24)
	if daysSincePublished == 0 {
		daysSincePublished = 1
	}
	return &daysSincePublished, nil
}

type Snippet struct {
	PublishedAt string `json:"publishedAt"`
	Title       string `json:"title"`
	ChannelId   string `json:"channelId"`
}

type SearchListEndpoint struct {
	apiKey string
	prefix string
}

func NewSearchListEndpoint(conf *Config) *SearchListEndpoint {
	return &SearchListEndpoint{
		apiKey: conf.ApiKey,
		prefix: apiPrefix,
	}
}

func (c SearchListEndpoint) Request(params Params) (*http.Request, error) {
	p := params.(*SearchParams)
	p.Validate()
	url := fmt.Sprintf("%v/search?key=%v&%v", c.prefix, c.apiKey, p.UrlParams())
	fmt.Printf("%s 🚀 url: %s\n", Indicator, YellowString(url))
	return http.NewRequest(http.MethodGet, url, nil)
}

func (c SearchListEndpoint) Response(body []byte) (YouTubeDataApi, error) {
	var result Search
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, errors.Wrap(err, "body of search list parse failed")
	}
	return result, nil
}

type SearchParams struct {
	Size     int
	Term     string
	Period   int
	prepared preparedSearchParams
}

type preparedSearchParams struct {
	key            string
	maxResults     int
	publishedAfter string
}

func (p SearchParams) Validate() bool {
	return true
}

func (p *SearchParams) Prepare() {
	var publishedAfter string
	if period, err := resolveSearchPeriod(p.Period); err != nil {
		log.Errorf("%#+v", errors.Wrap(err, "prepare published after failed"))
	} else {
		publishedAfter = period.Format(time.RFC3339)
	}
	p.prepared = preparedSearchParams{
		key:            p.Term,
		maxResults:     p.Size,
		publishedAfter: publishedAfter,
	}
}

func resolveSearchPeriod(period int) (*time.Time, error) {
	date := time.Now().AddDate(0, 0, -period)
	date, err := time.Parse("2006-01-02 15:04:05", date.Format("2006-01-02 15:04:05"))
	if err != nil {
		return nil, err
	}
	return &date, nil
}

func (p SearchParams) UrlParams() string {
	p.Prepare()
	pr := p.prepared
	// TODO
	// https://stackoverflow.com/questions/1760757/how-to-efficiently-concatenate-strings-in-go
	return fmt.Sprintf("part=id&part=snippet&q=%v&publishedAfter=%s&maxResults=%v&order=viewCount",
		pr.key, pr.publishedAfter, pr.maxResults)
}
