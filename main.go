package main

import (
	"os"
	"os/exec"

	"github.com/Songmu/prompter"
	"github.com/iljarotar/synth/context"
	"github.com/iljarotar/synth/signal"
	"github.com/iljarotar/synth/wave"
)

var waves = []wave.Wave{
	{Type: wave.Sine, Freq: 220, Amplitude: 1},
	{Type: wave.Sine, Freq: 275, Amplitude: 1},
	{Type: wave.Sine, Freq: 330, Amplitude: 1},
	{Type: wave.Sine, Freq: 415, Amplitude: 1},
}

func main() {
	ctx := context.NewContext()
	ctx.Init()
	defer ctx.Terminate()

	w := wave.NewWaveTable(waves...)
	s := signal.NewSignal(w)
	defer s.Close()

	clear()
	go s.Play()

	for input := prompter.Prompt(">", ""); input != "exit"; {
		input = prompter.Prompt(">", "")
	}

	s.Stop()
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
