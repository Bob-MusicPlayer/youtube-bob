package model

type PlaylistItems struct {
	Kind          string   `json:"kind"`
	Etag          string   `json:"etag"`
	NextPageToken string   `json:"nextPageToken"`
	PageInfo      PageInfo `json:"pageInfo"`
	Items         []Item   `json:"items"`
}

type Item struct {
	Kind    string  `json:"kind"`
	Etag    string  `json:"etag"`
	ID      string  `json:"id"`
	Snippet Snippet `json:"snippet"`
}

type Snippet struct {
	PublishedAt  string     `json:"publishedAt"`
	ChannelID    string     `json:"channelId"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Thumbnails   Thumbnails `json:"thumbnails"`
	ChannelTitle string     `json:"channelTitle"`
	PlaylistID   string     `json:"playlistId"`
	Position     int64      `json:"position"`
	ResourceID   ResourceID `json:"resourceId"`
}

type ResourceID struct {
	Kind    string `json:"kind"`
	VideoID string `json:"videoId"`
}

type Thumbnails struct {
	Default  Default `json:"default"`
	Medium   Default `json:"medium"`
	High     Default `json:"high"`
	Standard Default `json:"standard"`
	Maxres   Default `json:"maxres"`
}

type Default struct {
	URL    string `json:"url"`
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
}

type PageInfo struct {
	TotalResults   int64 `json:"totalResults"`
	ResultsPerPage int64 `json:"resultsPerPage"`
}
