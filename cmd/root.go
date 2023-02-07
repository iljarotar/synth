package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/Songmu/prompter"
	"github.com/iljarotar/synth/audio"
	"github.com/iljarotar/synth/control"
	l "github.com/iljarotar/synth/loader"
	s "github.com/iljarotar/synth/synth"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "synth",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		sampleRate, _ := cmd.Flags().GetString("sample-rate")

		err := start(file, sampleRate)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("file", "f", "", "specify which file to load")
	rootCmd.Flags().BoolP("help", "h", false, "print help")
	rootCmd.Flags().StringP("sample-rate", "s", "44100", "specify sample rate")
}

func start(file, sampleRate string) error {
	sRate, err := strconv.Atoi(sampleRate)
	if err != nil {
		return errors.New("could not parse sample rate. please provide an integer")
	}

	err = audio.Init()
	if err != nil {
		return err
	}
	defer audio.Terminate()

	input := make(chan float32)
	ctx, err := audio.NewContext(input, float64(sRate))
	if err != nil {
		return err
	}
	defer ctx.Close()

	err = ctx.Start()
	if err != nil {
		return err
	}

	ctl := control.NewControl(input)
	ctl.Start()

	loader, err := l.NewLoader()
	if err != nil {
		return err
	}
	defer loader.Close()

	var synth s.Synth
	err = loader.Load(file, &synth)
	if err != nil {
		return err
	}
	ctl.LoadSynth(synth)

	for {
		input := prompter.Prompt("type 'q' to quit", "")
		if input == "q" {
			break
		}
	}
	ctl.Stop()

	return nil
}

func loadFile(file string) error {
	return nil
}
