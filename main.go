package main

import (
	"fmt"
	"math"
	"time"

	a "github.com/iljarotar/synth/audio"
	w "github.com/iljarotar/synth/wave"
)

const sampleRate = 44100

func root(x float64) float64 {
	return math.Sin(2 * math.Pi * 440 * x)
}

func fifth(x float64) float64 {
	return math.Sin(2 * math.Pi * 660 * x)
}

func third(x float64) float64 {
	return math.Sin(2 * math.Pi * 550 * x)
}

var functions = []func(float64) float64{root, third}

func main() {
	audio := a.NewAudio()
	audio.Init()
	defer audio.Terminate()

	signal := a.NewSignal(w.Custom(sampleRate, functions))
	defer signal.Close()

	err := signal.Start()
	if err != nil {
		fmt.Print(err)
	}

	time.Sleep(time.Second * 10)

	err = signal.Stop()
	if err != nil {
		fmt.Print(err)
	}
}
