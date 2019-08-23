package repository

import (
	"fmt"
	shared "shared-bob"
	"youtube-bob/model"
)

type YoutubeRepository struct {
	requestHelper *shared.RequestHelper
	apiKey        string
}

func NewYoutubeRepository(apiKey string) *YoutubeRepository {
	return &YoutubeRepository{
		requestHelper: shared.NewRequestHelper("https://www.googleapis.com/youtube"),
		apiKey:        apiKey,
	}
}

func (y *YoutubeRepository) GetPlaylist(id string) (*model.PlaylistItems, error) {
	var playlistItems model.PlaylistItems

	response, err := y.requestHelper.Get(fmt.Sprintf("/v3/playlistItems?part=snippet&maxResults=50&playlistId=%s&key=%s", id, y.apiKey), nil)
	if err != nil {
		return nil, err
	}

	err = response.DecodeBody(&playlistItems)
	if err != nil {
		return nil, err
	}

	return &playlistItems, err
}

func (y *YoutubeRepository) Search(query string) (*model.Search, error) {
	var playlistItems model.Search

	response, err := y.requestHelper.Get(fmt.Sprintf("/v3/search?part=snippet&maxResults=25&q=%s&type=video,playlist&key=%s", query, y.apiKey), nil)
	if err != nil {
		return nil, err
	}

	err = response.DecodeBody(&playlistItems)
	if err != nil {
		return nil, err
	}

	return &playlistItems, err
}

func (y *YoutubeRepository) GetVideoInfo(id string) (*model.VideoInfo, error) {
	var videoInfo model.VideoInfo

	response, err := y.requestHelper.Get(fmt.Sprintf("/v3/videos?part=snippet&id=%s&key=%s", id, y.apiKey), nil)
	if err != nil {
		return nil, err
	}

	err = response.DecodeBody(&videoInfo)
	if err != nil {
		return nil, err
	}

	return &videoInfo, err
}
