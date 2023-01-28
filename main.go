package main

import (
	"fmt"

	"github.com/iljarotar/synth/audio"
	"github.com/iljarotar/synth/control"
	"github.com/iljarotar/synth/ui"
)

func main() {
	if err := audio.Init(); err != nil {
		fmt.Println(err)
		return
	}
	defer audio.Terminate()

	ctx, err := audio.NewContext()
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

	exit := make(chan bool)
	UI := ui.NewUI(ctl, exit)

	UI.ClearScreen()
	go UI.AcceptInput()
	<-exit
	UI.ClearScreen()
}
