package youtube

type (
	Channels struct {
		DataApiResponse
		Items []ChannelsItem
	}

	ChannelsItem struct {
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
)

func (c Channels) FindChannelTitle() string {
	return c.Items[0].BrandingSettings.Channel.Title
}
