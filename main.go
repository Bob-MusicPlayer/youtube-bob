package main

import (
	"fmt"
	"net/http"
	"youtube-bob/handler"
)

func main() {
	youtubeHandler := handler.NewYoutubeHandler()

	http.HandleFunc("/play", youtubeHandler.HandlePlay)

	fmt.Println(http.ListenAndServe(":5001", nil))
}