package cmd

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/iljarotar/synth/audio"
	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/control"
	"github.com/iljarotar/synth/file"
	"github.com/iljarotar/synth/log"
	"github.com/iljarotar/synth/ui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var version = "unknown"

var rootCmd = &cobra.Command{
	Use:     "synth",
	Version: version,
	Short:   "A modular synthesizer for the command line",
	Long: `A modular synthesizer for the command line.
	
Documentation and usage: https://github.com/iljarotar/synth`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := cmd.Flags().GetString("config")

		err := config.EnsureDefaultConfig()
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return cmd.Help()
		}

		if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
			return fmt.Errorf("too many arguments passed - at most one argument expected")
		}
		filename := args[0]

		defaultConfigPath, err := config.GetDefaultConfigPath()
		if err != nil {
			return err
		}

		if cfg == "" {
			cfg = defaultConfigPath
		}

		c, err := config.LoadConfig(cfg)
		if err != nil {
			return fmt.Errorf("could not load config file: %w", err)
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
}

func parseFlags(cmd *cobra.Command, config *config.Config) error {
	s, _ := cmd.Flags().GetInt("sample-rate")
	in, _ := cmd.Flags().GetFloat64("fade-in")
	out, _ := cmd.Flags().GetFloat64("fade-out")

	if cmd.Flag("sample-rate").Changed {
		config.SampleRate = s
	}
	if cmd.Flag("fade-in").Changed {
		config.FadeIn = in
	}
	if cmd.Flag("fade-out").Changed {
		config.FadeOut = out
	}

	return config.Validate()
}

func start(filename string, c *config.Config) error {
	logger := log.NewLogger(5)
	ctl, err := control.NewControl(logger, c)
	if err != nil {
		return err
	}

	loader, err := file.NewLoader(logger, filename, ctl.LoadSynth)
	if err != nil {
		return err
	}
	defer func() {
		err := loader.Close()
		if err != nil {
			fmt.Printf("failed to close loader: %v", err)
		}
	}()

	err = loader.LoadAndWatch()
	if err != nil {
		return err
	}

	audioCtx, err := audio.NewContext(int(c.SampleRate), ctl.ReadSample)
	if err != nil {
		return err
	}
	defer func() {
		runtime.KeepAlive(audioCtx)
	}()

	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to initialize raw terminal: %w", err)
	}
	defer func() {
		if err := term.Restore(int(os.Stdin.Fd()), state); err != nil {
			fmt.Printf("failed to restore terminal state: %v", err)
		}
	}()

	signalChan := make(chan ui.Signal)
	uiConfig := ui.Config{
		Logger:     logger,
		File:       filename,
		SignalChan: signalChan,
	}

	u := ui.NewUI(uiConfig)
	go u.Enter()

	done := make(chan bool)
	var fadingOut bool

Loop:
	for {
		select {
		case signal := <-signalChan:
			if signal == ui.SignalQuit {
				if fadingOut {
					logger.Info("already received quit signal")
					continue
				}

				fadingOut = true
				logger.Info(fmt.Sprintf("fading out in %fs", c.FadeOut))
				loader.Stop()
				go ctl.Stop(done, false)
			}

			if signal == ui.SignalInterrupt {
				loader.Stop()
				go ctl.Stop(done, true)
			}

		case <-done:
			break Loop
		}
	}

	time.Sleep(200 * time.Millisecond) // avoid clipping at the end
	ui.LineBreaks(2)
	return nil
}
