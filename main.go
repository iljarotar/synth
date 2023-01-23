package main

import (
	"fmt"
	"time"

	"github.com/iljarotar/synth/context"
	"github.com/iljarotar/synth/signal"
	"github.com/iljarotar/synth/wave"
)

func main() {
	ctx := context.NewContext()
	ctx.Init()
	defer ctx.Terminate()

	s := signal.NewSignal(wave.Sine(440))
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
