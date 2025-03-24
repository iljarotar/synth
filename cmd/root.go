package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/iljarotar/synth/audio"
	"github.com/iljarotar/synth/log"
	"golang.org/x/term"

	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/control"
	"github.com/iljarotar/synth/file"
	"github.com/iljarotar/synth/ui"
	"github.com/spf13/cobra"
)

var version = "unknown"

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

		err := config.EnsureDefaultConfig()
		if err != nil {
			fmt.Println(err)
			return
		}

		defaultConfigPath, err := config.GetDefaultConfigPath()
		if err != nil {
			fmt.Println(err)
			return
		}

		if cfg == "" {
			cfg = defaultConfigPath
		}
		config, err := config.LoadConfig(cfg)
		if err != nil {
			fmt.Printf("could not load config file: %v\n", err)
			return
		}

		err = parseFlags(cmd, config)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = start(file, config)
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
	rootCmd.Flags().Float64P("sample-rate", "s", config.DefaultSampleRate, "sample rate")
	rootCmd.Flags().Float64P("fade-in", "i", config.DefaultFadeIn, "fade-in in seconds")
	rootCmd.Flags().Float64P("fade-out", "o", config.DefaultFadeOut, "fade-out in seconds")
	rootCmd.Flags().StringP("config", "c", defaultConfigPath, "path to your config file")
	rootCmd.Flags().Float64P("duration", "d", config.DefaultDuration, "duration in seconds; if positive duration is specified, synth will stop playing after the defined time")
}

func parseFlags(cmd *cobra.Command, config *config.Config) error {
	s, _ := cmd.Flags().GetFloat64("sample-rate")
	in, _ := cmd.Flags().GetFloat64("fade-in")
	out, _ := cmd.Flags().GetFloat64("fade-out")
	duration, _ := cmd.Flags().GetFloat64("duration")

	if cmd.Flag("sample-rate").Changed {
		config.SampleRate = s
	}
	if cmd.Flag("fade-in").Changed {
		config.FadeIn = in
	}
	if cmd.Flag("fade-out").Changed {
		config.FadeOut = out
	}
	if cmd.Flag("duration").Changed {
		config.Duration = duration
	}

	return config.Validate()
}

func start(fileName string, cfg *config.Config) error {
	logger := log.NewLogger(10)
	quit := make(chan bool)
	autoStop := make(chan bool)
	var closing bool
	interrupt := make(chan bool)

	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to read input %v", err))
	}
	defer func() {
		if err := term.Restore(int(os.Stdin.Fd()), state); err != nil {
			logger.Error(fmt.Sprintf("failed to restore terminal %v", err))
		}
	}()

	output := make(chan audio.AudioOutput)
	a, err := audio.NewAudio(output, int(cfg.SampleRate))
	if err != nil {
		return err
	}
	defer func() {
		err := a.Close()
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}()

	ctl, err := control.NewControl(logger, *cfg, output, autoStop, &closing)
	if err != nil {
		return err
	}
	ctl.Start()
	defer ctl.StopSynth()

	u := ui.NewUI(logger, fileName, quit, autoStop, cfg.Duration, &closing, interrupt, ctl)
	go u.Enter()

	loader, err := file.NewLoader(logger, ctl, fileName, &closing)
	if err != nil {
		return err
	}
	defer loader.Close()

	err = loader.Load()
	if err != nil {
		ui.Clear()
		return fmt.Errorf("unable to load file %s: %w", fileName, err)
	}

	ctl.FadeIn(cfg.FadeIn)
	var fadingOut bool

Loop:
	for {
		select {
		case <-quit:
			if fadingOut {
				logger.Info("already received quit signal")
				continue
			}
			fadingOut = true
			logger.Info(fmt.Sprintf("fading out in %fs", cfg.FadeOut))
			ctl.Stop(cfg.FadeOut)
		case <-interrupt:
			logger.Info("interrupt received")
			ctl.Stop(0.05)
		case <-ctl.SynthDone:
			break Loop
		}
	}

	time.Sleep(time.Millisecond * 200) // avoid clipping at the end
	ui.LineBreaks(2)
	return err
}
