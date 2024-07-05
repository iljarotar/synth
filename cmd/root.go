package cmd

import (
	"fmt"
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
		cfg, _ := cmd.Flags().GetString("config")

		if file == "" {
			cmd.Help()
			return
		}

		err := c.EnsureDefaultConfig()
		if err != nil {
			fmt.Printf("%v", err)
			return
		}

		defaultConfigPath, err := c.GetDefaultConfigPath()
		if err != nil {
			fmt.Printf("%v", err)
			return
		}

		if cfg == "" {
			cfg = defaultConfigPath
		}
		err = c.LoadConfig(cfg)
		if err != nil {
			fmt.Printf("could not load config file: %v", err)
			return
		}

		err = parseFlags(cmd)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}

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
	rootCmd.Flags().StringP("file", "f", "", "Path to your patch file")
	rootCmd.Flags().StringP("out", "o", "", "If provided, a .wav file with the given name will be recorded. A positive value for --duration must be given for the recording to work.")
	rootCmd.Flags().StringP("sample-rate", "s", "", "Sample rate")
	rootCmd.Flags().StringP("duration", "d", "", "Duration in seconds excluding fade-in and fade-out. a negative duration will cause the synth to play until stopped manually")
	rootCmd.Flags().String("fade-in", "", "Fade-in in seconds")
	rootCmd.Flags().String("fade-out", "", "Fade-out in seconds")
	rootCmd.Flags().StringP("config", "c", "", "Path to your config file.")
}

func parseFlags(cmd *cobra.Command) error {
	s, _ := cmd.Flags().GetString("sample-rate")
	in, _ := cmd.Flags().GetString("fade-in")
	out, _ := cmd.Flags().GetString("fade-out")
	d, _ := cmd.Flags().GetString("duration")
	record, _ := cmd.Flags().GetString("out")

	if s != "" {
		sampleRate, err := utils.ParseInt(s)
		if err != nil {
			return fmt.Errorf("could not parse sample rate: %w", err)
		}
		c.Config.SampleRate = float64(sampleRate)
	}

	if in != "" {
		fadeIn, err := utils.ParseFloat(in)
		if err != nil {
			return fmt.Errorf("could not parse fade-in: %w", err)
		}
		c.Config.FadeIn = fadeIn
	}

	if out != "" {
		fadeOut, err := utils.ParseFloat(out)
		if err != nil {
			return fmt.Errorf("could not parse fade-out: %w", err)
		}
		c.Config.FadeOut = fadeOut
	}

	if d != "" {
		duration, err := utils.ParseFloat(d)
		if err != nil {
			return fmt.Errorf("could not parse duration: %w", err)
		}
		c.Config.Duration = duration
	}

	if record != "" {
		c.Config.Out = record
	}

	return c.Config.Validate()
}

func start(file string) error {
	err := audio.Init()
	if err != nil {
		return err
	}
	defer audio.Terminate()

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
	rec := f.NewRecorder(recIn, speakerIn, c.Config.Out)
	go rec.StartRecording()

	exit := make(chan bool)
	quit := make(chan bool)
	u := ui.NewUI(file, quit)
	go u.Enter(exit)

	ctl := control.NewControl(recIn, exit)
	defer ctl.Close()

	loader, err := f.NewLoader(ctl, file)
	if err != nil {
		return err
	}
	defer loader.Close()

	err = loader.Load()
	if err != nil {
		return err
	}

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
	ui.LineBreaks(1)
	return err
}

func catchInterrupt(stop chan bool, sig chan os.Signal) {
	<-sig
	stop <- true
}
