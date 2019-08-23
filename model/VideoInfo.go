package model

type VideoInfo struct {
	Kind     string      `json:"kind"`
	Etag     string      `json:"etag"`
	PageInfo PageInfo    `json:"pageInfo"`
	Items    []VideoItem `json:"items"`
}

type VideoItem struct {
	Kind    string       `json:"kind"`
	Etag    string       `json:"etag"`
	ID      string       `json:"id"`
	Snippet VideoSnippet `json:"snippet"`
}

type VideoSnippet struct {
	PublishedAt          string     `json:"publishedAt"`
	ChannelID            string     `json:"channelId"`
	Title                string     `json:"title"`
	Description          string     `json:"description"`
	Thumbnails           Thumbnails `json:"thumbnails"`
	ChannelTitle         string     `json:"channelTitle"`
	Tags                 []string   `json:"tags"`
	CategoryID           string     `json:"categoryId"`
	LiveBroadcastContent string     `json:"liveBroadcastContent"`
	Localized            Localized  `json:"localized"`
}

type Localized struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
