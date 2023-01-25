package main

import (
	"fmt"
	"io/ioutil"

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
	ui.ClearScreen()

	c := config.Instance()
	input := make(chan float32)
	ctx, err := audio.NewContext(c.SampleRate, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()

	ctl := control.NewControl(ctx)

	data, err := ioutil.ReadFile("patches/example.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ctl.Parse(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ctl.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	done := make(chan bool)
	go ui.AcceptInput(done)
	<-done

	err = ctl.Stop()
	if err != nil {
		fmt.Println(err)
	}
}
