package main

import (
	"os"
	"os/exec"

	"github.com/Songmu/prompter"
	"github.com/iljarotar/synth/context"
	"github.com/iljarotar/synth/signal"
	"github.com/iljarotar/synth/wave"
)

func main() {
	ctx := context.NewContext()
	ctx.Init()
	defer ctx.Terminate()

	w := wave.NewWaveTable(wave.Wave{Type: wave.Sine})
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
