package handler

import (
	bobModel "bob/model"
	"fmt"
	"net/http"
	shared "shared-bob"
	"youtube-bob/model"
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

	err = yh.player.SetPlayback(fmt.Sprintf("https://www.youtube.com/watch?v=%s", playback.ID))
	if responseHelper.ReturnHasError(err) {
		return
	}

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

	paused, err := yh.player.IsPaused()
	if responseHelper.ReturnHasError(err) {
		return
	}

	responseHelper.ReturnOk(model.PlaybackInfo{
		CachedTime: cachedTime,
		Position:   position,
		Duration:   duration,
		Title:      title,
		Paused:     paused,
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

	playlists, err := yh.youtubeRepository.Search(query)
	if responseHelper.ReturnHasError(err) {
		return
	}

	responseHelper.ReturnOk(playlists)
}
