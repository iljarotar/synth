package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/iljarotar/synth/audio"
	"github.com/iljarotar/synth/ui"
	"gopkg.in/yaml.v2"

	c "github.com/iljarotar/synth/config"
	s "github.com/iljarotar/synth/synth"
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
		config, err := c.LoadConfig(cfg)
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
	defaultConfigPath, err := c.GetDefaultConfigPath()
	if err != nil {
		os.Exit(1)
	}
	rootCmd.Flags().Float64P("sample-rate", "s", c.DefaultSampleRate, "sample rate")
	rootCmd.Flags().Float64P("fade-in", "i", c.DefaultFadeIn, "fade-in in seconds")
	rootCmd.Flags().Float64P("fade-out", "o", c.DefaultFadeOut, "fade-out in seconds")
	rootCmd.Flags().StringP("config", "c", defaultConfigPath, "path to your config file")
	rootCmd.Flags().Float64P("duration", "d", c.DefaultDuration, "duration in seconds; if positive duration is specified, synth will stop playing after the defined time")
}

func parseFlags(cmd *cobra.Command, config *c.Config) error {
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

func start(file string, config *c.Config) error {
	err := audio.Init()
	if err != nil {
		return err
	}
	defer audio.Terminate()

	output := make(chan audio.AudioOutput)
	ctx, err := audio.NewContext(output, config.SampleRate)
	if err != nil {
		return err
	}
	defer ctx.Close()

	err = ctx.Start()
	if err != nil {
		return err
	}

	bytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	synth := s.Synth{}
	err = yaml.Unmarshal(bytes, &synth)
	if err != nil {
		return err
	}

	ctl, err := s.NewControl(&synth, *config, output)
	if err != nil {
		return err
	}

	p := tea.NewProgram(ui.NewModel(ctl), tea.WithAltScreen())

	callbacks := s.Callbacks{
		Quit: func() {
			p.Send(ui.QuitMsg(true))
		},
		UpdateTime: func(time float64) {
			p.Send(ui.TimeMsg(time))
		},
		SendVolumeWarning: func(output float64) {
			p.Send(ui.VolumeWarningMsg(output))
		},
		ShowVolume: func(volume float64) {
			p.Send(ui.VolumeMsg(volume))
		},
	}

	ctl.SetCallbacks(callbacks)
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("unable to start synth: %w", err)
	}

	return nil
}
