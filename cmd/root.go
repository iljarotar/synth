package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/Songmu/prompter"
	"github.com/iljarotar/synth/audio"
	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/control"
	l "github.com/iljarotar/synth/loader"
	s "github.com/iljarotar/synth/synth"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "synth",
	Short: "command line synthesizer",
	Long:  "command line synthesizer",
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		s, _ := cmd.Flags().GetString("sample-rate")

		if s != "" {
			sRate, err := parseSampleRate(s)
			if err != nil {
				fmt.Println("could not parse sample rate. please provide an integer")
				return
			}
			config.Instance.SetSampleRate(sRate)
		}

		err := start(file)
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
	rootCmd.Flags().StringP("sample-rate", "s", "", "specify sample rate")
}

func start(file string) error {
	err := audio.Init()
	if err != nil {
		return err
	}
	defer audio.Terminate()
	clear()

	input := make(chan float32)
	ctx, err := audio.NewContext(input)
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

	loader, err := l.NewLoader(ctl)
	if err != nil {
		return err
	}
	defer loader.Close()

	var synth s.Synth
	err = loader.Load(file, &synth)
	if err != nil {
		return err
	}

	for {
		input := prompter.Prompt("type 'q' to quit", "")
		if input == "q" {
			ctl.Stop()
			time.Sleep(time.Millisecond * 50) // avoid clipping at the end
			break
		}
	}

	return nil
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func parseSampleRate(input string) (float64, error) {
	sampleRate, err := strconv.Atoi(input)
	if err != nil {
		return 0, errors.New("could not parse sample rate. please provide an integer")
	}

	return float64(sampleRate), nil
}
