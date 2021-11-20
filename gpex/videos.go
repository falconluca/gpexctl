package gpex

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

type Videos struct {
	Kind     string `json:"kind"`
	Etag     string `json:"etag"`
	Items    []VideoItem
	PageInfo PageInfo
	Error    YouTubeDataApiError `json:"error"`
}

func (v Videos) Meta() string {
	return fmt.Sprintf("kind: %v, etag: %v", v.Kind, v.Etag)
}

func (v Videos) HasErrors() bool {
	return len(v.Error.Errors) > 0
}

type VideoItem struct {
	Kind       string `json:"kind"`
	Etag       string `json:"etag"`
	Id         string `json:"id"`
	Statistics struct {
		ViewCount     string `json:"viewCount"`
		LikeCount     string `json:"likeCount"`
		DislikeCount  string `json:"dislikeCount"`
		FavoriteCount string `json:"favoriteCount"`
		CommentCount  string `json:"commentCount"`
	} `json:"statistics"`
}

type VideoListEndpoint struct {
	apiKey string
	prefix string
}

func NewVideoListEndpoint(conf *Config) *VideoListEndpoint {
	return &VideoListEndpoint{
		apiKey: conf.ApiKey,
		prefix: apiPrefix,
	}
}

func (v VideoListEndpoint) Request(params Params) (*http.Request, error) {
	p := params.(*VideoParams)
	p.Validate()
	url := fmt.Sprintf("%v/videos?key=%v&%v", v.prefix, v.apiKey, p.UrlParams())
	return http.NewRequest(http.MethodGet, url, nil)
}

func (v VideoListEndpoint) Response(body []byte) (YouTubeDataApi, error) {
	var result Videos
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, errors.Wrap(err, "body of video list parse failed")
	}
	return result, nil
}

type VideoParams struct {
	VideoId  string
	prepared preparedVideoParams
}

type preparedVideoParams struct {
	Id string
}

func (v VideoParams) Validate() bool {
	return true
}

func (v *VideoParams) Prepare() {
	v.prepared = preparedVideoParams{
		Id: v.VideoId,
	}
}

func (v VideoParams) UrlParams() string {
	v.Prepare()
	pr := v.prepared
	return fmt.Sprintf("part=statistics&id=%v", pr.Id)
}
