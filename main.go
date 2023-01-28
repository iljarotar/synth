package main

import (
	"fmt"

	"github.com/iljarotar/synth/audio"
	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/control"
	"github.com/iljarotar/synth/ui"
)

func main() {
	if err := audio.Init(); err != nil {
		fmt.Println(err)
		return
	}
	defer audio.Terminate()

	c := config.Instance()
	input := make(chan float32)
	ctx, err := audio.NewContext(c.SampleRate, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()

	ctl := control.NewControl(ctx)
	err = ctl.LoadSynth()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ctl.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	UI := ui.NewUI()
	UI.ClearScreen()

	done := make(chan bool)
	go UI.AcceptInput(done)
	<-done
	UI.ClearScreen()

	err = ctl.Stop()
	if err != nil {
		fmt.Println(err)
	}
}
