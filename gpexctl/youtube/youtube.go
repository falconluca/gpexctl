package youtube

import "fmt"

type (
	DataApiResponse struct {
		Kind  string `json:"kind"`
		Etag  string `json:"etag"`
		Error struct {
			Errors []Error `json:"errors"`
		} `json:"error"`
	}

	Error struct {
		Message string `json:"message"`
		Domain  string `json:"domain"`
		Reason  string `json:"reason"`
	}

	PageInfo struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	}
)

func (r DataApiResponse) HasErrors() bool {
	return len(r.Error.Errors) > 0
}

func (r DataApiResponse) Meta() string {
	return fmt.Sprintf("kind: %v, etag: %v", r.Kind, r.Etag)
}
