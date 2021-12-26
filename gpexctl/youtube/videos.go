package youtube

type (
	Videos struct {
		DataApiResponse
		Items    []VideoItem
		PageInfo PageInfo
	}

	VideoItem struct {
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
)
