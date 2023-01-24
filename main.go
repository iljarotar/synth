package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Songmu/prompter"
	"github.com/iljarotar/synth/audio"
	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/wave"
)

var waves = []wave.Wave{
	{Type: wave.Sine, Freq: 220, Amplitude: 1},
	{Type: wave.Sine, Freq: 275, Amplitude: 1},
	{Type: wave.Sine, Freq: 330, Amplitude: 1},
	{Type: wave.Sine, Freq: 415, Amplitude: 1},
}

var noise = []wave.Wave{
	{Type: wave.Noise, Amplitude: 1},
}

func main() {
	audio.Init()
	defer audio.Terminate()

	c := config.Instance()
	w := wave.NewWaveTable(noise...)

	ctx, err := audio.NewContext(c.SampleRate, w.Process)
	if err != nil {
		fmt.Println(err)
	}

	defer ctx.Close()

	clear()
	go ctx.Start()

	for input := prompter.Prompt(">", ""); input != "exit"; {
		input = prompter.Prompt(">", "")
	}

	ctx.Stop()
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
