package module

import (
	"math"

	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/utils"
)

type TextProcessor struct {
	Module
	Name string `yaml:"name"`
	Text string `yaml:"text"`
	BPM  Param  `yaml:"bpm"`
	Amp  Param  `yaml:"amp"`
	Pan  Param  `yaml:"pan"`
	data []float64
}

func (p *TextProcessor) Initialize() {
	p.limitParams()
	p.data = utils.Normalize(p.textToSignal(p.Text), -1, 1)

	y := p.signalValue(0, p.Amp.Val, p.BPM.Val)
	p.current = stereo(y, p.Pan.Val)
}

func (p *TextProcessor) Next(t float64, modMap ModulesMap) {
	pan := utils.Limit(p.Pan.Val+modulate(p.Pan.Mod, modMap)*p.Pan.ModAmp, panLimits.min, panLimits.max)
	amp := utils.Limit(p.Amp.Val+modulate(p.Amp.Mod, modMap)*p.Amp.ModAmp, ampLimits.min, ampLimits.max)
	bpm := utils.Limit(p.BPM.Val+modulate(p.BPM.Mod, modMap)*p.BPM.ModAmp, bpmLimits.min, bpmLimits.max)

	y := p.signalValue(t, amp, bpm)
	p.current = stereo(y, pan)
}

func (p *TextProcessor) limitParams() {
	p.Amp.ModAmp = utils.Limit(p.Amp.ModAmp, ampLimits.min, ampLimits.max)
	p.Amp.Val = utils.Limit(p.Amp.Val, ampLimits.min, ampLimits.max)

	p.Pan.ModAmp = utils.Limit(p.Pan.ModAmp, panLimits.min, panLimits.max)
	p.Pan.Val = utils.Limit(p.Pan.Val, panLimits.min, panLimits.max)

	p.BPM.ModAmp = utils.Limit(p.BPM.ModAmp, bpmLimits.min, bpmLimits.max)
	p.BPM.Val = utils.Limit(p.BPM.Val, bpmLimits.min, bpmLimits.max)
}

func (p *TextProcessor) signalValue(t, amp, bpm float64) float64 {
	freq := bpm / 60 / float64(len(p.data))

	idx := int(math.Floor(t * float64(len(p.data)) * freq))
	var val float64

	if l := len(p.data); l > 0 {
		val = p.data[idx%l]
	}

	y := amp * val
	p.integral += y / config.Config.SampleRate

	return y
}

func (p *TextProcessor) textToSignal(input string) []float64 {
	signal := make([]float64, len(p.Text))

	for i, char := range p.Text {
		signal[i] = float64(int(char))
	}

	return signal
}
