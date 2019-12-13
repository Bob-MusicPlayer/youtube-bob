package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"youtube-bob/handler"
	"youtube-bob/repository"
	"youtube-bob/util"
)

var youtubeApiKey string
var bobUrl string

const (
	version = "0.0.1"
)

func FlagsAreValid() bool {
	youtubeApiKeyArg := flag.String("apiKey", "", "The Api Key from YouTube")
	bobUrlArg := flag.String("bobUrl", "http://localhost:5002", "The url of bob")

	flag.Parse()

	if *youtubeApiKeyArg == "" {
		flag.PrintDefaults()
		return false
	}

	youtubeApiKey = *youtubeApiKeyArg
	bobUrl = *bobUrlArg

	return true
}

func main() {
	if !FlagsAreValid() {
		os.Exit(1)
	}

	printBanner()

	logrus.Info("Start Youtube Bob")

	bobRepository := repository.NewBobRepository(bobUrl)

	player, err := util.NewPlayer(bobRepository)
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

	logrus.WithFields(logrus.Fields{
		"hostname": "0.0.0.0",
		"port":     "5001",
	}).Info("Start Webserver")
	fmt.Println(http.ListenAndServe("0.0.0.0:5001", nil))
}

func printBanner() {
	fmt.Println("+------------------------------------------------+")
	fmt.Println("|                                                |")
	fmt.Println("|    \u001b[31m▄██████████▄\033[0m   ██████\033[1;34m╗\033[0m  ██████\033[1;34m╗\033[0m ██████\033[1;34m╗\033[0m     |")
	fmt.Println("|   \u001b[31m██████\033[0;41m█▄\u001b[31m██████\033[0m  ██\033[1;34m╔══\033[0m██\033[1;34m╗\033[0m██\033[1;34m╔═══\033[0m██\033[1;34m╗\033[0m██\033[1;34m╔══\033[0m██\033[1;34m╗\033[0m    |")
	fmt.Println("|   \u001b[31m██████\033[0;41m███■\u001b[31m████\033[0m  █████\033[1;34m╔╝\033[0m███\033[1;34m║\033[0m   ██\033[1;34m║\033[0m██████\033[1;34m╔╝\033[0m    |")
	fmt.Println("|   \u001b[31m██████\033[0;41m█▀\u001b[31m██████\033[0m  ██\033[1;34m╔══\033[0m██\033[1;34m╗\033[0m██\033[1;34m║\033[0m   ██\033[1;34m║\033[0m██\033[1;34m╔══\033[0m██\033[1;34m╗\033[0m    |")
	fmt.Println("|    \u001b[31m▀██████████▀\033[0m   ██████\033[1;34m╔╝╚\033[0m██████\033[1;34m╔╝\033[0m██████\033[1;34m╔╝\033[0m    |")
	fmt.Println("|                   \033[1;34m╚═════╝  ╚═════╝ ╚═════╝\033[0m     |")
	fmt.Println("|                                                |")
	fmt.Println("+------------------------------------------------+")
	fmt.Printf("|  Version: %7s                              |\n", version)
	fmt.Println("+------------------------------------------------+")
}
