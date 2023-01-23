package main

import (
	"fmt"
	"time"

	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/context"
	"github.com/iljarotar/synth/signal"
	"github.com/iljarotar/synth/wave"
)

const sampleRate = 44100

func main() {
	ctx := context.NewContext()
	ctx.Init()
	defer ctx.Terminate()

	c := config.NewConfig(sampleRate)

	s := signal.NewSignal(wave.Sine(c, 440))
	defer s.Close()

	err := s.Start()
	if err != nil {
		fmt.Print(err)
	}

	time.Sleep(time.Second * 10)

	err = s.Stop()
	if err != nil {
		fmt.Print(err)
	}
}
