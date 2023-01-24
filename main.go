package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/Songmu/prompter"
	"github.com/iljarotar/synth/audio"
	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/control"
)

func main() {
	if err := audio.Init(); err != nil {
		fmt.Println(err)
		return
	}
	defer audio.Terminate()
	clear()

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
