package gramophone

import (
	"bytes"
	"os"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

type Gramophone struct {
	player    *oto.Player
	ctx       *oto.Context
	readyChan chan struct{}
	isPlaying bool
}

func New() *Gramophone {
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

	return &Gramophone{
		player:    &oto.Player{},
		ctx:       otoCtx,
		readyChan: readyChan,
		isPlaying: false,
	}
}

func (this *Gramophone) Play(filepaths ...string) {
	if this.isPlaying == false {
		filepath := filepaths[0]

		fileBytes, err := os.ReadFile(filepath)
		if err != nil {
			panic("reading " + filepath + " failed: " + err.Error())
		}

		// Convert the pure bytes into a reader object that can be used with the mp3 decoder
		fileBytesReader := bytes.NewReader(fileBytes)

		// Decode file
		decodedMp3, err := mp3.NewDecoder(fileBytesReader)
		if err != nil {
			panic("mp3.NewDecoder failed: " + err.Error())
		}

		// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
		<-this.readyChan

		// Create a new 'player' that will handle our sound. Paused by default.
		this.player = this.ctx.NewPlayer(decodedMp3)
	}

	// Play starts playing the sound and returns without waiting for it (Play() is async).
	this.player.Play()
	this.isPlaying = true
}

func (this *Gramophone) Pause() {
	this.player.Pause()
	this.isPlaying = false
}

func (this *Gramophone) IsPlaying() bool {
	return this.isPlaying
}

func (this *Gramophone) Close() error {
	if err := this.player.Close(); err != nil {
		return err
	}

	return nil
}
