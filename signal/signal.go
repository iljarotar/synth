package signal

import (
	"fmt"

	"github.com/gordonklaus/portaudio"
	"github.com/iljarotar/synth/wave"
)

type Signal struct {
	*portaudio.Stream
}

func NewSignal(w *wave.WaveTable) *Signal {
	s := &Signal{}

	var err error
	s.Stream, err = portaudio.OpenDefaultStream(0, 1, w.Config.SampleRate, 0, w.Process)
	if err != nil {
		fmt.Print(err)
	}

	return s
}
