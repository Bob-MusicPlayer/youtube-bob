package handler

import (
	bobModel "bob/model"
	"fmt"
	"net/http"
	shared "shared-bob"
	"strconv"
	"time"
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

	for {
		if title, err := yh.player.GetTitle(); err != nil || (title != fmt.Sprintf("watch?v=%s", playback.ID) && title != "") {
			yh.player.CurrentPlayback.Title = title
			break
		}
		time.Sleep(time.Millisecond * 200)
	}
	if responseHelper.ReturnHasError(err) {
		return
	}

	go yh.player.ListenForCacheChanges()

	if len(videoInfo.Items) == 0 {
		fmt.Println("Videoinfo cant be loaded for id " + playback.ID + ". This is strange.")
		return
	}

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

func (yh *YoutubeHandler) HandlePlaybackSeek(w http.ResponseWriter, req *http.Request) {
	responseHelper := shared.NewResponseHelper(w, req)
	if responseHelper.ReturnOptionsOrNotAllowed(http.MethodPost) {
		return
	}

	seconds := req.URL.Query().Get("seconds")

	fmt.Println(seconds)

	if seconds == "" {
		responseHelper.ReturnError(fmt.Errorf("seconds must be specified"))
		return
	}

	sec, err := strconv.Atoi(seconds)
	if responseHelper.ReturnHasError(err) {
		return
	}

	err = yh.player.Seek(sec)
	if responseHelper.ReturnHasError(err) {
		return
	}

	responseHelper.ReturnOk(nil)
}
