package cmd

import (
	"fmt"
	"os"

	"github.com/iljarotar/synth/audio"
	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/file"
	"github.com/iljarotar/synth/log"
	"github.com/iljarotar/synth/player"
	"github.com/iljarotar/synth/ui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var version = "unknown"

var rootCmd = &cobra.Command{
	Use:     "synth",
	Version: version,
	Short:   "command line synthesizer",
	Long: `command line synthesizer
	
documentation and usage: https://github.com/iljarotar/synth`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := cmd.Flags().GetString("config")

		if len(args) == 0 {
			cmd.Help()
			return nil
		}

		if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
			return fmt.Errorf("too many arguments passed - at most one argument expected")
		}
		filename := args[0]

		err := config.EnsureDefaultConfig()
		if err != nil {
			return err
		}

		defaultConfigPath, err := config.GetDefaultConfigPath()
		if err != nil {
			return err
		}

		if cfg == "" {
			cfg = defaultConfigPath
		}

		c, err := config.LoadConfig(cfg)
		if err != nil {
			return fmt.Errorf("could not load config file: %v\n", err)
		}

		err = parseFlags(cmd, c)
		if err != nil {
			return err
		}

		return start(filename, c)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	defaultConfigPath, err := config.GetDefaultConfigPath()
	if err != nil {
		os.Exit(1)
	}
	rootCmd.Flags().IntP("sample-rate", "s", config.DefaultSampleRate, "sample rate")
	rootCmd.Flags().Float64P("fade-in", "i", config.DefaultFadeIn, "fade-in in seconds")
	rootCmd.Flags().Float64P("fade-out", "o", config.DefaultFadeOut, "fade-out in seconds")
	rootCmd.Flags().StringP("config", "c", defaultConfigPath, "path to your config file")
	rootCmd.Flags().Float64P("duration", "d", config.DefaultDuration, "duration in seconds; if positive duration is specified, synth will stop playing after the defined time")
}

func parseFlags(cmd *cobra.Command, config *config.Config) error {
	s, _ := cmd.Flags().GetInt("sample-rate")
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

func start(filename string, c *config.Config) error {
	logger := log.NewLogger(10)
	p, err := player.NewPlayer(logger, filename, int(c.SampleRate))
	if err != nil {
		return err
	}

	loader, err := file.NewLoader(logger, filename, p.UpdateSynth)
	if err != nil {
		return err
	}
	defer func() {
		err := loader.Close()
		if err != nil {
			fmt.Printf("failed to close loader:%v", err)
		}
	}()

	err = loader.LoadAndWatch()
	if err != nil {
		return err
	}

	audioCtx, err := audio.NewContext(int(c.SampleRate), p.ReadSample)
	if err != nil {
		return err
	}
	defer func() {
		err := audioCtx.Close()
		if err != nil {
			fmt.Printf("failed to close audio context:%v", err)
		}
	}()

	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to initialize raw terminal:%w", err)
	}
	defer func() {
		if err := term.Restore(int(os.Stdin.Fd()), state); err != nil {
			fmt.Printf("failed to restore terminal state:%v", err)
		}
	}()

	quitChan := make(chan bool)
	uiConfig := ui.Config{
		Logger:   logger,
		File:     filename,
		Duration: c.Duration,
		QuitChan: quitChan,
	}

	u := ui.NewUI(uiConfig)
	go u.Enter()

	<-quitChan

	return nil
}
