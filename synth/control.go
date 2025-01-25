package synth

import (
	"math"

	"github.com/iljarotar/synth/audio"
	cfg "github.com/iljarotar/synth/config"
)

type QuitFunc func()
type TimeIsUpFunc func()
type UpdateTimeFunc func(time float64)
type ShowVolumeWarningFunc func(output float64)
type ShowVolumeFunc func(volume float64)

type Callbacks struct {
	Quit              QuitFunc
	TimeIsUp          TimeIsUpFunc
	UpdateTime        UpdateTimeFunc
	SendVolumeWarning ShowVolumeWarningFunc
	ShowVolume        ShowVolumeFunc
}

type Control struct {
	config                   cfg.Config
	synth                    *Synth
	output                   chan audio.AudioOutput
	synthDone                chan bool
	autoStop                 chan bool
	maxOutput                float64
	volumeWarningTriggeredAt float64
	callbacks                Callbacks
	quitting                 bool
	secondsPassed            float64
}

func NewControl(synth *Synth, config cfg.Config, output chan audio.AudioOutput) (*Control, error) {
	err := synth.initialize(config.SampleRate)
	if err != nil {
		return nil, err
	}

	callbacks := Callbacks{
		Quit:              func() {},
		TimeIsUp:          func() {},
		UpdateTime:        func(time float64) {},
		SendVolumeWarning: func(output float64) {},
		ShowVolume:        func(volume float64) {},
	}

	ctl := &Control{
		config:    config,
		synth:     synth,
		output:    output,
		synthDone: make(chan bool),
		autoStop:  make(chan bool),
		callbacks: callbacks,
	}

	return ctl, nil
}

func (c *Control) GetVolume() float64 {
	return c.synth.volumeMemory
}

func (c *Control) IncreaseVolume() {
	vol := c.synth.Volume + 0.02
	if vol > maxVolume {
		vol = maxVolume
	}
	c.synth.volumeMemory = vol
	c.synth.Volume = vol
}

func (c *Control) DecreaseVolume() {
	vol := c.synth.Volume - 0.02
	if vol < 0 {
		vol = 0
	}
	c.synth.volumeMemory = vol
	c.synth.Volume = vol
	c.maxOutput = 0
}

func (c *Control) SetCallbacks(callbacks Callbacks) {
	c.callbacks = callbacks
}

func (c *Control) Start() {
	outputChan := make(chan Output)
	c.synth.active = true
	go c.synth.play(outputChan)
	go c.receiveOutput(outputChan)
	c.synth.startFading(FadeDirectionIn, c.config.FadeIn)
}

func (c *Control) Stop() {
	c.synth.notifyFadeOutDone(c.synthDone)
	c.synth.startFading(FadeDirectionOut, c.config.FadeOut)
	go func() {
		<-c.synthDone
		c.synth.active = false
		c.callbacks.Quit()
	}()
}

func (c *Control) checkVolume(output, time float64) {
	// only consider up to three decimals
	abs := math.Round(math.Abs(output)*1000) / 1000

	if abs > 1 && abs > c.maxOutput && time-c.volumeWarningTriggeredAt >= 0.5 {
		c.maxOutput = abs
		c.volumeWarningTriggeredAt = time
		c.callbacks.SendVolumeWarning(abs)
		return
	}

	if c.maxOutput <= 1 && c.volumeWarningTriggeredAt > 0 {
		c.volumeWarningTriggeredAt = 0
		c.callbacks.SendVolumeWarning(0)
	}
}

func (c *Control) receiveOutput(outputChan <-chan Output) {
	defer close(c.output)

	for out := range outputChan {
		c.checkDuration(out.Time)
		c.checkVolume(out.Mono, out.Time)
		c.sendTime(out.Time)

		c.output <- audio.AudioOutput{
			Left:  out.Left,
			Right: out.Right,
		}
	}
}

func (c *Control) sendTime(time float64) {
	if time-c.secondsPassed >= 1 {
		c.secondsPassed += 1
		c.callbacks.UpdateTime(c.secondsPassed)
	}
}

func (c *Control) checkDuration(time float64) {
	if c.config.Duration < 0 {
		return
	}
	if c.quitting {
		return
	}
	duration := c.config.Duration - c.config.FadeOut
	if time < duration {
		return
	}
	c.quitting = true
	c.callbacks.TimeIsUp()
}
