package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/iljarotar/synth/audio"
	c "github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/control"
	l "github.com/iljarotar/synth/loader"
	s "github.com/iljarotar/synth/screen"
	"github.com/iljarotar/synth/utils"
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
		in, _ := cmd.Flags().GetString("fade-in")
		out, _ := cmd.Flags().GetString("fade-out")

		if file == "" {
			cmd.Help()
			return
		}

		sampleRate, err := utils.ParseInt(s)
		if err != nil {
			fmt.Println("could not parse sample rate: %w", err)
			return
		}
		c.Config.SampleRate = sampleRate

		fadeIn, err := utils.ParseFloat(in)
		if err != nil {
			fmt.Println("could not parse fade-in: %w", err)
			return
		}
		c.Config.FadeIn = fadeIn

		fadeOut, err := utils.ParseFloat(out)
		if err != nil {
			fmt.Println("could not parse fade-out: %w", err)
			return
		}
		c.Config.FadeOut = fadeOut

		err = start(file)
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
	sampleRate := fmt.Sprintf("%v", c.Default.SampleRate)
	fadeIn := fmt.Sprintf("%v", c.Default.FadeIn)
	fadeOut := fmt.Sprintf("%v", c.Default.FadeOut)

	rootCmd.Flags().StringP("file", "f", "", "specify which file to load")
	rootCmd.Flags().BoolP("help", "h", false, "print help")
	rootCmd.Flags().StringP("sample-rate", "s", sampleRate, "specify sample rate")
	rootCmd.Flags().String("fade-in", fadeIn, "length of the fade-in in seconds")
	rootCmd.Flags().String("fade-out", fadeOut, "length of the fade-out in seconds")
}

func start(file string) error {
	err := audio.Init()
	if err != nil {
		return err
	}
	defer audio.Terminate()
	s.Clear()

	input := make(chan struct{ Left, Right float32 })
	ctx, err := audio.NewContext(input, c.Config.SampleRate)
	if err != nil {
		return err
	}
	defer ctx.Close()

	err = ctx.Start()
	if err != nil {
		return err
	}

	ctl := control.NewControl(input)

	log := make(chan string)
	logger := s.NewLogger(log)

	loader, err := l.NewLoader(ctl, logger, file)
	if err != nil {
		return err
	}
	defer loader.Close()

	err = loader.Load()
	if err != nil {
		return err
	}

	quit := make(chan bool)
	exit := make(chan bool)
	screen := s.NewScreen(logger, quit)
	go screen.Enter(exit)
	ctl.Start(c.Config.FadeIn)

	<-quit

	ctl.Stop(c.Config.FadeOut)
	time.Sleep(time.Millisecond * 200) // avoid clipping at the end
	return nil
}
