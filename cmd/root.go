package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iljarotar/synth/audio"
	"github.com/iljarotar/synth/config"
	c "github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/control"
	f "github.com/iljarotar/synth/file"
	"github.com/iljarotar/synth/ui"
	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "synth",
	Version: version,
	Short:   "command line synthesizer",
	Long: `command line synthesizer
	
documentation and usage: https://github.com/iljarotar/synth`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := cmd.Flags().GetString("config")

		if len(args) == 0 {
			cmd.Help()
			return
		}

		if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
			fmt.Println("too many arguments passed - at most one argument expected")
		}
		file := args[0]

		err := c.EnsureDefaultConfig()
		if err != nil {
			fmt.Println(err)
			return
		}

		defaultConfigPath, err := c.GetDefaultConfigPath()
		if err != nil {
			fmt.Println(err)
			return
		}

		if cfg == "" {
			cfg = defaultConfigPath
		}
		err = c.LoadConfig(cfg)
		if err != nil {
			fmt.Printf("could not load config file: %v\n", err)
			return
		}

		err = parseFlags(cmd)
		if err != nil {
			fmt.Println(err)
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
	defaultConfigPath, err := config.GetDefaultConfigPath()
	if err != nil {
		os.Exit(1)
	}
	rootCmd.Flags().Float64P("sample-rate", "s", config.Default.SampleRate, "sample rate")
	rootCmd.Flags().Float64("fade-in", config.Default.FadeIn, "fade-in in seconds")
	rootCmd.Flags().Float64("fade-out", config.Default.FadeOut, "fade-out in seconds")
	rootCmd.Flags().StringP("config", "c", defaultConfigPath, "path to your config file")
	rootCmd.Flags().Float64P("duration", "d", config.Default.Duration, "duration in seconds; if positive duration is specified, synth will stop playing after the defined time")
}

func parseFlags(cmd *cobra.Command) error {
	s, _ := cmd.Flags().GetFloat64("sample-rate")
	in, _ := cmd.Flags().GetFloat64("fade-in")
	out, _ := cmd.Flags().GetFloat64("fade-out")
	duration, _ := cmd.Flags().GetFloat64("duration")

	c.Config.SampleRate = s
	c.Config.FadeIn = in
	c.Config.FadeOut = out
	c.Config.Duration = duration

	return c.Config.Validate()
}

func start(file string) error {
	err := audio.Init()
	if err != nil {
		return err
	}
	defer audio.Terminate()

	quit := make(chan bool)
	autoStop := make(chan bool)
	u := ui.NewUI(file, quit, autoStop)
	go u.Enter()

	output := make(chan struct{ Left, Right float32 })
	ctx, err := audio.NewContext(output, c.Config.SampleRate)
	if err != nil {
		return err
	}
	defer ctx.Close()

	err = ctx.Start()
	if err != nil {
		return err
	}

	ctl := control.NewControl(output, autoStop)
	ctl.Start()
	defer ctl.StopSynth()

	loader, err := f.NewLoader(ctl, file)
	if err != nil {
		return err
	}
	defer loader.Close()

	err = loader.Load()
	if err != nil {
		ui.Clear()
		return fmt.Errorf("unable to load file %s: %w", file, err)
	}

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	interrupt := make(chan bool)
	go catchInterrupt(interrupt, sig)

	ctl.FadeIn(c.Config.FadeIn)
	var fadingOut bool

Loop:
	for {
		select {
		case <-quit:
			if fadingOut {
				ui.Logger.Info("already received quit signal")
				continue
			}
			fadingOut = true
			ui.Logger.Info(fmt.Sprintf("fading out in %fs", c.Config.FadeOut))
			ctl.Stop(c.Config.FadeOut)
		case <-interrupt:
			ctl.Stop(0.05)
		case <-ctl.SynthDone:
			break Loop
		}
	}

	time.Sleep(time.Millisecond * 200) // avoid clipping at the end
	ui.LineBreaks(2)
	return err
}

func catchInterrupt(stop chan bool, sig chan os.Signal) {
	<-sig
	stop <- true
}
