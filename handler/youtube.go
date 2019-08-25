package handler

import (
	bobModel "bob/model"
	"net/http"
	shared "shared-bob"
	"youtube-bob/repository"
	"youtube-bob/util"
)

type YoutubeHandler struct {
	player            *util.Player
	youtubeRepository *repository.YoutubeRepository
}

func NewYoutubeHandler(player *util.Player, youtubeRepository *repository.YoutubeRepository) YoutubeHandler {
	return YoutubeHandler{
		player:            player,
		youtubeRepository: youtubeRepository,
	}
}

func (yh *YoutubeHandler) HandlePlay(w http.ResponseWriter, req *http.Request) {
	responseHelper := shared.NewResponseHelper(w, req)
	if responseHelper.ReturnOptionsOrNotAllowed(http.MethodPost) {
		return
	}

	err := yh.player.Play()
	if responseHelper.ReturnHasError(err) {
		return
	}

	responseHelper.ReturnOk(nil)
}

func (yh *YoutubeHandler) HandlePause(w http.ResponseWriter, req *http.Request) {
	responseHelper := shared.NewResponseHelper(w, req)
	if responseHelper.ReturnOptionsOrNotAllowed(http.MethodPost) {
		return
	}

	err := yh.player.Pause()
	if responseHelper.ReturnHasError(err) {
		return
	}

	responseHelper.ReturnOk(nil)
}

func (yh *YoutubeHandler) HandleSetPlayback(w http.ResponseWriter, req *http.Request) {
	responseHelper := shared.NewResponseHelper(w, req)
	if responseHelper.ReturnOptionsOrNotAllowed(http.MethodPost) {
		return
	}

	var playback bobModel.Playback

	err := responseHelper.DecodeBody(&playback)
	if responseHelper.ReturnHasError(err) {
		return
	}

	err = yh.player.SetPlayback(playback.ID)
	if responseHelper.ReturnHasError(err) {
		return
	}

	videoInfo, err := yh.youtubeRepository.GetVideoInfo(playback.ID)
	if responseHelper.ReturnHasError(err) {
		return
	}

	yh.player.CurrentPlayback.Title = videoInfo.Items[0].Snippet.Title
	yh.player.CurrentPlayback.Author = videoInfo.Items[0].Snippet.ChannelTitle
	yh.player.CurrentPlayback.ThumbnailUrl = videoInfo.Items[0].Snippet.Thumbnails.High.URL

	responseHelper.ReturnOk(nil)
}

func (yh *YoutubeHandler) HandlePlaybackInfo(w http.ResponseWriter, req *http.Request) {
	responseHelper := shared.NewResponseHelper(w, req)
	if responseHelper.ReturnOptionsOrNotAllowed(http.MethodGet) {
		return
	}

	cachedTime, _ := yh.player.GetCacheTime()
	position, _ := yh.player.GetPosition()
	duration, _ := yh.player.GetDuration()
	title, _ := yh.player.GetTitle()

	isPlaying, err := yh.player.IsPlaying()
	if responseHelper.ReturnHasError(err) {
		return
	}

	if yh.player.CurrentPlayback == nil {
		responseHelper.ReturnOk(bobModel.Playback{
			Source:    "youtube",
			IsPlaying: false,
		})
		return
	}

	responseHelper.ReturnOk(bobModel.Playback{
		ID:            yh.player.CurrentPlayback.ID,
		Title:         title,
		Author:        yh.player.CurrentPlayback.Author,
		Position:      position,
		Duration:      duration,
		CachePosition: cachedTime,
		Source:        "youtube",
		ThumbnailUrl:  yh.player.CurrentPlayback.ThumbnailUrl,
		IsPlaying:     isPlaying,
	})
}

func (yh *YoutubeHandler) HandlePlaylist(w http.ResponseWriter, req *http.Request) {
	responseHelper := shared.NewResponseHelper(w, req)
	if responseHelper.ReturnOptionsOrNotAllowed(http.MethodGet) {
		return
	}

	playlists, err := yh.youtubeRepository.GetPlaylist("PLRBp0Fe2GpgmgoscNFLxNyBVSFVdYmFkq")
	if responseHelper.ReturnHasError(err) {
		return
	}

	responseHelper.ReturnOk(playlists)
}

func (yh *YoutubeHandler) HandleSearch(w http.ResponseWriter, req *http.Request) {
	responseHelper := shared.NewResponseHelper(w, req)
	if responseHelper.ReturnOptionsOrNotAllowed(http.MethodGet) {
		return
	}

	query := req.URL.Query().Get("q")

	songs := make([]bobModel.Playback, 0)

	playlists, err := yh.youtubeRepository.Search(query)
	if responseHelper.ReturnHasError(err) {
		return
	}

	for _, playback := range playlists.Items {
		playback := bobModel.Playback{
			Title:        playback.Snippet.Title,
			Author:       playback.Snippet.ChannelTitle,
			Duration:     0,
			ID:           playback.ID.VideoID,
			ThumbnailUrl: playback.Snippet.Thumbnails.High.URL,
			Source:       "youtube",
		}

		songs = append(songs, playback)
	}

	responseHelper.ReturnOk(songs)
}
