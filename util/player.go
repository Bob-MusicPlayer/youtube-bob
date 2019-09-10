package util

import (
	bobModel "bob/model"
	"fmt"
	"github.com/DexterLB/mpvipc"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"
	"youtube-bob/repository"
)

type Player struct {
	mpvCommand      *exec.Cmd
	ipc             *mpvipc.Connection
	CurrentPlayback *bobModel.Playback
	bobRepository   *repository.BobRepository
}

func NewPlayer(bobRepository *repository.BobRepository) (*Player, error) {
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

	time.Sleep(time.Millisecond * 400)

	err = ipc.Open()

	if err != nil {
		return nil, err
	}

	return &Player{
		ipc:           ipc,
		mpvCommand:    mpvCommand,
		bobRepository: bobRepository,
	}, nil
}

func (p *Player) SetPlayback(id string) error {
	_, err := p.ipc.Call("loadfile", fmt.Sprintf("https://www.youtube.com/watch?v=%s", id))
	if err != nil {
		return err
	}

	p.CurrentPlayback = &bobModel.Playback{
		ID:     id,
		Source: "youtube",
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

func (p *Player) IsPlaying() (bool, error) {
	i, err := p.ipc.Get("pause")
	if err != nil {
		return false, err
	}

	return !i.(bool), nil
}

func (p *Player) Seek(seconds int) error {
	fmt.Println(seconds)
	err := p.ipc.Set("time-pos", seconds)
	if err != nil {
		return err
	}

	return nil
}

func (p *Player) ListenForCacheChanges() error {
	lastTime := float64(0)
	startTime := time.Now()

	for {
		t, err := p.GetCacheTime()
		if err != nil {
			continue
		}
		if lastTime != t {
			p.CurrentPlayback.CachePosition = t
			p.bobRepository.Sync()
			lastTime = t
		}
		duration, err := p.GetDuration()
		if err != nil {
			continue
		}

		if time.Now().Sub(startTime) > time.Second * 5 || (duration > 0 && duration - 2 < t)  {
			break
		}
		time.Sleep(time.Millisecond * 10)
	}

	return nil
}
