package util

import (
	"fmt"
	"github.com/DexterLB/mpvipc"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"
)

type Player struct {
	mpvCommand *exec.Cmd
	ipc        *mpvipc.Connection
}

func NewPlayer() (*Player, error) {
	socketPath := ""

	if runtime.GOOS == "windows" {
		socketPath = "\\\\.\\pipe\\tmp\\mpv-socket"
		fmt.Println(socketPath)
	} else {
		socketPath = path.Join(os.TempDir(), "mpv-socket")
	}

	mpvCommand := exec.Command("mpv", "--idle=yes", fmt.Sprintf("--input-ipc-server=%s", socketPath), "-ytdl-format=bestaudio", "--keep-open=always", "--pause")
	err := mpvCommand.Start()
	if err != nil {
		return nil, err
	}

	ipc := mpvipc.NewConnection(socketPath)

	time.Sleep(time.Millisecond * 200)

	err = ipc.Open()

	if err != nil {
		return nil, err
	}

	return &Player{
		ipc:        ipc,
		mpvCommand: mpvCommand,
	}, nil
}

func (p *Player) SetPlayback(url string) error {
	_, err := p.ipc.Call("loadfile", url)
	if err != nil {
		return err
	}

	err = p.Play()

	return err
}

func (p *Player) Pause() error {
	err := p.ipc.Set("pause", true)

	return err
}

func (p *Player) Play() error {
	err := p.ipc.Set("pause", false)

	return err
}

func (p *Player) GetCacheTime() (float64, error) {
	i, err := p.ipc.Get("demuxer-cache-time")
	if err != nil {
		return 0, err
	}

	return i.(float64), nil
}

func (p *Player) GetPosition() (float64, error) {
	i, err := p.ipc.Get("time-pos")
	if err != nil {
		return 0, err
	}

	return i.(float64), nil
}

func (p *Player) GetDuration() (float64, error) {
	i, err := p.ipc.Get("duration")
	if err != nil {
		return 0, err
	}

	return i.(float64), nil
}

func (p *Player) GetTitle() (string, error) {
	i, err := p.ipc.Get("media-title")
	if err != nil {
		return "", err
	}

	return i.(string), nil
}

func (p *Player) IsPaused() (bool, error) {
	i, err := p.ipc.Get("pause")
	if err != nil {
		return false, err
	}

	return i.(bool), nil
}
