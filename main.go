package main

import (
	"time"

	"github.com/iljarotar/synth/context"
	"github.com/iljarotar/synth/signal"
	"github.com/iljarotar/synth/wave"
)

func main() {
	ctx := context.NewContext()
	ctx.Init()
	defer ctx.Terminate()

	w := wave.NewWaveTable(wave.SineFunc(440), wave.NoiseFunc())
	s := signal.NewSignal(w)
	defer s.Close()

	s.Start()
	time.Sleep(time.Second * 10)
	s.Stop()
}
