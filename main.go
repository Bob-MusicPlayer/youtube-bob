package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"youtube-bob/handler"
	"youtube-bob/repository"
	"youtube-bob/util"
)

var youtubeApiKey string

func FlagsAreValid() bool {
	youtubeApiKeyArg := flag.String("apiKey", "", "The Api Key from YouTube")

	flag.Parse()

	if youtubeApiKeyArg == nil {
		flag.PrintDefaults()
		return false
	}

	youtubeApiKey = *youtubeApiKeyArg

	return true
}

func main() {
	if !FlagsAreValid() {
		os.Exit(1)
	}

	player, err := util.NewPlayer()
	if err != nil {
		panic(err)
	}

	youtubeRepository := repository.NewYoutubeRepository(youtubeApiKey)

	youtubeHandler := handler.NewYoutubeHandler(player, youtubeRepository)

	http.HandleFunc("/api/v1/playback", youtubeHandler.HandleSetPlayback)
	http.HandleFunc("/api/v1/play", youtubeHandler.HandlePlay)
	http.HandleFunc("/api/v1/pause", youtubeHandler.HandlePause)
	http.HandleFunc("/api/v1/playback/info", youtubeHandler.HandlePlaybackInfo)
	http.HandleFunc("/api/v1/playback/seek", youtubeHandler.HandlePlaybackSeek)

	http.HandleFunc("/api/v1/playlist", youtubeHandler.HandlePlaylist)
	http.HandleFunc("/api/v1/search", youtubeHandler.HandleSearch)

	fmt.Println(http.ListenAndServe(":5001", nil))
}
