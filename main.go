package main

import (
	"fmt"
	"net/http"
	"youtube-bob/handler"
	"youtube-bob/util"
)

func main() {
	player, err := util.NewPlayer()
	if err != nil {
		panic(err)
	}

	youtubeHandler := handler.NewYoutubeHandler(player)

	http.HandleFunc("/api/v1/playback", youtubeHandler.HandleSetPlayback)
	http.HandleFunc("/api/v1/play", youtubeHandler.HandlePlay)
	http.HandleFunc("/api/v1/pause", youtubeHandler.HandlePause)
	http.HandleFunc("/api/v1/playback/info", youtubeHandler.HandlePlaybackInfo)

	fmt.Println(http.ListenAndServe(":5001", nil))
}
