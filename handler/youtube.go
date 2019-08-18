package handler

import (
	"encoding/json"
	"net/http"
	"youtube-bob/model"
	"youtube-bob/util"
)

type YoutubeHandler struct {
	player *util.Player
}

func NewYoutubeHandler(player *util.Player) YoutubeHandler {
	return YoutubeHandler{
		player: player,
	}
}

func (yh *YoutubeHandler) HandlePlay(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost && req.Method != http.MethodOptions {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if req.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(200)
		return
	}

	err := yh.player.Play()

	if err != nil {
		w.WriteHeader(901)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(200)
}

func (yh *YoutubeHandler) HandlePause(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost && req.Method != http.MethodOptions {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if req.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(200)
		return
	}

	err := yh.player.Pause()

	if err != nil {
		w.WriteHeader(901)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(200)
}

func (yh *YoutubeHandler) HandleSetPlayback(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost && req.Method != http.MethodOptions {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if req.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(200)
		return
	}

	var playback model.SetPlayback

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&playback)
	if err != nil {
		w.WriteHeader(901)
		w.Write([]byte(err.Error()))
		return
	}

	err = yh.player.SetPlayback(playback.Url)
	if err != nil {
		w.WriteHeader(901)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
}

func (yh *YoutubeHandler) HandlePlaybackInfo(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost && req.Method != http.MethodOptions {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if req.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(200)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var playbackInfo model.PlaybackInfo

	cachedTime, err := yh.player.GetCacheTime()
	if err != nil {
		w.WriteHeader(901)
		w.Write([]byte(err.Error()))
		return
	}

	position, err := yh.player.GetPosition()
	if err != nil {
		w.WriteHeader(901)
		w.Write([]byte(err.Error()))
		return
	}

	duration, err := yh.player.GetDuration()
	if err != nil {
		w.WriteHeader(901)
		w.Write([]byte(err.Error()))
		return
	}

	title, err := yh.player.GetTitle()
	if err != nil {
		w.WriteHeader(901)
		w.Write([]byte(err.Error()))
	}

	playbackInfo.CachedTime = cachedTime
	playbackInfo.Position = position
	playbackInfo.Duration = duration
	playbackInfo.Title = title

	json, err := json.Marshal(&playbackInfo)
	if err != nil {
		w.WriteHeader(901)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(json)
}
