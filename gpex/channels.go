package gpex

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

type Channels struct {
	Kind  string `json:"kind"`
	Etag  string `json:"etag"`
	Items []ChannelsItem
	Error YouTubeDataApiError `json:"error"`
}

func (c Channels) Meta() string {
	return fmt.Sprintf("kind: %v, etag: %v", c.Kind, c.Etag)
}

func (c Channels) HasErrors() bool {
	return len(c.Error.Errors) > 0
}

func (c Channels) FindChannelTitle() string {
	return c.Items[0].BrandingSettings.Channel.Title
}

type ChannelsItem struct {
	Kind       string `json:"kind"`
	Etag       string `json:"etag"`
	Id         string `json:"id"`
	Statistics struct {
		HiddenSubscriberCount bool   `json:"hiddenSubscriberCount"`
		SubscriberCount       string `json:"subscriberCount"`
	} `json:"statistics"`
	BrandingSettings struct {
		Channel struct {
			Title string `json:"title"`
		} `json:"channel"`
	} `json:"brandingSettings"`
}

type ChannelListEndpoint struct {
	apiKey string
	prefix string
}

func NewChannelListEndpoint(conf *Config) *ChannelListEndpoint {
	return &ChannelListEndpoint{
		apiKey: conf.ApiKey,
		prefix: apiPrefix,
	}
}

func (c ChannelListEndpoint) Request(params Params) (*http.Request, error) {
	p := params.(*ChannelParams)
	p.Validate()
	url := fmt.Sprintf("%v/channels?key=%v&%v", c.prefix, c.apiKey, p.UrlParams())
	return http.NewRequest(http.MethodGet, url, nil)
}

func (c ChannelListEndpoint) Response(body []byte) (YouTubeDataApi, error) {
	var result Channels
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, errors.Wrap(err, "body of channel list parse failed")
	}
	return result, nil
}

type ChannelParams struct {
	ChannelId string
	prepared  preparedChannelParams
}

type preparedChannelParams struct {
	Id string
}

func (v ChannelParams) Validate() bool {
	return true
}

func (v *ChannelParams) Prepare() {
	v.prepared = preparedChannelParams{
		Id: v.ChannelId,
	}
}

func (v ChannelParams) UrlParams() string {
	v.Prepare()
	pr := v.prepared
	// FIXME
	return "part=brandingSettings%2Cstatistics&&id=" + pr.Id
}
