package cmd

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iljarotar/synth/audio"
	c "github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/control"
	f "github.com/iljarotar/synth/file"
	"github.com/iljarotar/synth/ui"
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
		d, _ := cmd.Flags().GetString("duration")
		record, _ := cmd.Flags().GetString("out")

		if file == "" {
			cmd.Help()
			return
		}

		sampleRate, err := utils.ParseInt(s)
		if err != nil {
			fmt.Println("could not parse sample rate:", err)
			return
		}
		c.Config.SampleRate = float64(sampleRate)

		fadeIn, err := utils.ParseFloat(in)
		if err != nil {
			fmt.Println("could not parse fade-in:", err)
			return
		}
		c.Config.FadeIn = fadeIn

		fadeOut, err := utils.ParseFloat(out)
		if err != nil {
			fmt.Println("could not parse fade-out:", err)
			return
		}
		c.Config.FadeOut = fadeOut

		duration, err := utils.ParseInt(d)
		if err != nil {
			fmt.Println("could not parse duration:", err)
			return
		}

		if duration*sampleRate > math.MaxInt32 {
			fmt.Printf("duration too long. maximum duration is floor(%v / samplerate)\n", math.MaxInt32)
			return
		}
		c.Config.Duration = duration

		err = start(file, record)
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
	duration := fmt.Sprintf("%v", c.Config.Duration)

	rootCmd.Flags().StringP("file", "f", "", "path to your patch file")
	rootCmd.Flags().StringP("out", "o", "", "if provided recording will be written to the given file")
	rootCmd.Flags().BoolP("help", "h", false, "print help")
	rootCmd.Flags().StringP("sample-rate", "s", sampleRate, "sample rate")
	rootCmd.Flags().StringP("duration", "d", duration, "duration in seconds. if omitted playback will continue until stopped manually")
	rootCmd.Flags().String("fade-in", fadeIn, "length of the fade-in in seconds")
	rootCmd.Flags().String("fade-out", fadeOut, "length of the fade-out in seconds")
}

func start(file, record string) error {
	err := audio.Init()
	if err != nil {
		return err
	}
	defer audio.Terminate()
	ui.Clear()

	speakerIn := make(chan struct{ Left, Right float32 })
	ctx, err := audio.NewContext(speakerIn, c.Config.SampleRate)
	if err != nil {
		return err
	}
	defer ctx.Close()

	err = ctx.Start()
	if err != nil {
		return err
	}

	recIn := make(chan struct{ Left, Right float32 })
	rec := f.NewRecorder(recIn, speakerIn, record)
	go rec.StartRecording()

	exit := make(chan bool)
	ctl := control.NewControl(recIn, exit)
	defer ctl.Close()

	log := make(chan string)
	logger := ui.NewLogger(log)

	loader, err := f.NewLoader(ctl, logger, file)
	if err != nil {
		return err
	}
	defer loader.Close()

	err = loader.Load()
	if err != nil {
		return err
	}

	quit := make(chan bool)
	u := ui.NewUI(logger, quit)
	go u.Enter(exit)

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	interrupt := make(chan bool)
	go catchInterrupt(interrupt, sig)
	ctl.Start(c.Config.FadeIn)

	select {
	case <-quit:
		ctl.Stop(c.Config.FadeOut)
	case <-interrupt:
		ctl.Stop(0.05)
	}

	err = rec.StopRecording()
	time.Sleep(time.Millisecond * 200) // avoid clipping at the end
	ui.Clear()
	return err
}

func catchInterrupt(stop chan bool, sig chan os.Signal) {
	<-sig
	stop <- true
}
