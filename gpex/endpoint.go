package gpex

import (
	"net/http"
)

var (
	apiPrefix = "https://www.googleapis.com/youtube/v3"
)

type YouTubeDataApi interface {
	Meta() string
	HasErrors() bool
}

// YouTubeDataApiError api err
type YouTubeDataApiError struct {
	Errors []struct {
		Message string `json:"message"`
		Domain  string `json:"domain"`
		Reason  string `json:"reason"`
	} `json:"errors"`
}

type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

type ApiEndpoint interface {
	Request(params Params) (*http.Request, error)
	Response([]byte) (YouTubeDataApi, error)
}

type Params interface {
	Validate() bool
	Prepare()
}
