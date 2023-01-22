package audio

import (
	"fmt"

	"github.com/gordonklaus/portaudio"
	w "github.com/iljarotar/synth/wave"
)

type Audio struct {
}

type Signal struct {
	*portaudio.Stream
}

func NewAudio() *Audio {
	return &Audio{}
}

func NewSignal(wave *w.WaveTable) *Signal {
	s := &Signal{}

	var err error
	s.Stream, err = portaudio.OpenDefaultStream(0, 1, wave.SampleRate, 0, wave.Process)
	if err != nil {
		fmt.Print(err)
	}

	return s
}

func (a *Audio) Init() {
	portaudio.Initialize()
}

func (a *Audio) Terminate() {
	portaudio.Terminate()
}
