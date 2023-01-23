package signal

import (
	"fmt"

	"github.com/gordonklaus/portaudio"
	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/wave"
)

type Signal struct {
	*portaudio.Stream
}

func NewSignal(w wave.WaveTable) *Signal {
	s := &Signal{}
	c := config.Instance()

	var err error
	s.Stream, err = portaudio.OpenDefaultStream(0, 1, c.SampleRate, 0, w.Process)
	if err != nil {
		fmt.Print(err)
	}

	return s
}

func (s *Signal) Play() {
	err := s.Stream.Start()
	if err != nil {
		panic(err)
	}
}

func (s *Signal) Stop() {
	err := s.Stream.Stop()
	if err != nil {
		panic(err)
	}
}
