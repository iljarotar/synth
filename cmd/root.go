package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/iljarotar/synth/audio"
	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/control"
	l "github.com/iljarotar/synth/loader"
	"github.com/iljarotar/synth/screen"
	s "github.com/iljarotar/synth/synth"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "synth",
	Short: "command line synthesizer",
	Long: `command line synthesizer
	
documentation and usage: https://github.com/iljarotar/synth`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		s, _ := cmd.Flags().GetString("sample-rate")

		if file == "" {
			cmd.Help()
			return
		}

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
	screen.Clear()

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

	log := make(chan string)
	logger := screen.NewLogger(log)

	loader, err := l.NewLoader(ctl, logger)
	if err != nil {
		return err
	}
	defer loader.Close()

	var synth s.Synth
	err = loader.Load(file, &synth)
	if err != nil {
		return err
	}

	done := make(chan bool)
	s := screen.NewScreen(logger, done)
	go s.Enter()
	<-done

	ctl.Stop()
	time.Sleep(time.Millisecond * 200) // avoid clipping at the end
	return nil
}

func parseSampleRate(input string) (float64, error) {
	sampleRate, err := strconv.Atoi(input)
	if err != nil {
		return 0, errors.New("could not parse sample rate. please provide an integer")
	}

	return float64(sampleRate), nil
}
