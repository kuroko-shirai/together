package player

import (
	"bytes"
	"os"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

type Player struct {
	player    *oto.Player
	ctx       *oto.Context
	readyChan chan struct{}
	isPlaying bool
}

func New() *Player {
	// Prepare an Oto context (this will use your default audio device) that will
	// play all our sounds. Its configuration can't be changed later.

	options := &oto.NewContextOptions{}

	// Usually 44100 or 48000. Other values might cause distortions in Oto
	options.SampleRate = 44100

	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
	options.ChannelCount = 2

	// Format of the source. go-mp3's format is signed 16bit integers.
	options.Format = oto.FormatSignedInt16LE

	otoCtx, readyChan, err := oto.NewContext(options)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}

	return &Player{
		player:    &oto.Player{},
		ctx:       otoCtx,
		readyChan: readyChan,
		isPlaying: false,
	}
}

func (p *Player) Play(filepaths ...string) {
	if p.isPlaying == false {
		filepath := filepaths[0]

		fileBytes, err := os.ReadFile(filepath)
		if err != nil {
			panic("reading my-file.mp3 failed: " + err.Error())
		}

		// Convert the pure bytes into a reader object that can be used with the mp3 decoder
		fileBytesReader := bytes.NewReader(fileBytes)

		// Decode file
		decodedMp3, err := mp3.NewDecoder(fileBytesReader)
		if err != nil {
			panic("mp3.NewDecoder failed: " + err.Error())
		}

		// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
		<-p.readyChan

		// Create a new 'player' that will handle our sound. Paused by default.
		p.player = p.ctx.NewPlayer(decodedMp3)
	}

	// Play starts playing the sound and returns without waiting for it (Play() is async).
	p.player.Play()
	p.isPlaying = true
}

func (p *Player) Pause() {
	p.player.Pause()
	p.isPlaying = false
}

func (p *Player) IsPlaying() bool {
	return p.isPlaying
}

func (p *Player) Close() error {
	if err := p.player.Close(); err != nil {
		return err
	}

	return nil
}
