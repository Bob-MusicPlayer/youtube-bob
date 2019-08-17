package handler

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
)

type YoutubeHandler struct {
}

func NewYoutubeHandler() YoutubeHandler {
	return YoutubeHandler{
	}
}

func (yh *YoutubeHandler) HandlePlay(w http.ResponseWriter, req *http.Request) {
	url, err := url.Parse("https://www.youtube.com/watch?v=pZzSq8WfsKo")
	if err != nil {
		w.WriteHeader(901)
		w.Write([]byte(err.Error()))
		return
	}

	fmt.Println(url)

	// youtube-dl -o - https://www.youtube.com/watch?v=DKnIpsHe_YM

	args := "mpv --no-video -ytdl-format=bestaudio https://www.youtube.com/watch?v=pZzSq8WfsKo"
	mpv := exec.Command("mpv", strings.Split(args, " ")...)

	stderr, _ := mpv.StderrPipe()
	mpv.Start()

	scanner := bufio.NewScanner(stderr)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}