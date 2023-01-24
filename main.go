package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Songmu/prompter"
	"github.com/iljarotar/synth/audio"
	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/control"
	"github.com/iljarotar/synth/synth"
	"github.com/iljarotar/synth/wave"
)

var waves = []wave.Wave{
	{Type: wave.Sine, Freq: 220, Amplitude: 1},
	{Type: wave.Sine, Freq: 275, Amplitude: 1},
	{Type: wave.Sine, Freq: 330, Amplitude: 1},
	{Type: wave.Sine, Freq: 415, Amplitude: 1},
	{Type: wave.Noise, Amplitude: 0.1},
}

func main() {
	if err := audio.Init(); err != nil {
		fmt.Println(err)
		return
	}
	defer audio.Terminate()
	clear()

	w := wave.NewWaveTable(waves...)
	s := synth.NewSynth(w)

	c := config.Instance()
	input := make(chan float32)
	ctx, err := audio.NewContext(c.SampleRate, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()

	ctl := control.NewControl(*ctx, *s)
	err = ctl.Start()

	if err != nil {
		fmt.Println(err)
		return
	}

	for input := prompter.Prompt(">", ""); input != "exit"; {
		input = prompter.Prompt(">", "")
	}

	err = ctl.Stop()
	if err != nil {
		fmt.Println(err)
	}
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
