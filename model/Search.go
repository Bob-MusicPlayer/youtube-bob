package model

type Search struct {
	Kind          string       `json:"kind"`
	Etag          string       `json:"etag"`
	NextPageToken string       `json:"nextPageToken"`
	RegionCode    string       `json:"regionCode"`
	PageInfo      PageInfo     `json:"pageInfo"`
	Items         []SearchItem `json:"items"`
}

type SearchItem struct {
	Kind    string  `json:"kind"`
	Etag    string  `json:"etag"`
	ID      ID      `json:"id"`
	Snippet Snippet `json:"snippet"`
}

type ID struct {
	Kind    string `json:"kind"`
	VideoID string `json:"videoId"`
}
