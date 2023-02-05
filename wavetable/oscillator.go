package wavetable

type OscillatorType string

func (t OscillatorType) String() string {
	return string(t)
}

const (
	Sawtooth         OscillatorType = "Sawtooth"
	InvertedSawtooth OscillatorType = "InvertedSawtooth"
	Sine             OscillatorType = "Sine"
	Square           OscillatorType = "Square"
	Triangle         OscillatorType = "Triangle"
	Noise            OscillatorType = "Noise"
)

type WaveTable struct {
	SignalFunc  SignalFunc
	Amplitude   *Param        `yaml:"amplitude"`
	Oscillators []*Oscillator `yaml:"oscillators"`
	Filters     []Filter      `yaml:"filters"`
}

type Oscillator struct {
	Type       OscillatorType `yaml:"type"`
	Freq       float64        `yaml:"freq"`
	Amplitude  *Param         `yaml:"amplitude"`
	Phase      *Param         `yaml:"phase"`
	signalFunc SignalFunc
}

type Param struct {
	Value      float64    `yaml:"value"`
	Modulation *WaveTable `yaml:"modulation"`
}

func (w *WaveTable) Initialize() {
	if w.Amplitude != nil && w.Amplitude.Modulation != nil {
		w.Amplitude.Modulation.Initialize()
	}

	for i := range w.Filters {
		w.Filters[i].Initialize()
	}

	for i := range w.Oscillators {
		osc := w.Oscillators[i]
		osc.signalFunc = NewSignalFunc(osc.Type)

		if osc.Phase != nil && osc.Phase.Modulation != nil {
			osc.Phase.Modulation.Initialize()
		}

		if osc.Amplitude != nil && osc.Amplitude.Modulation != nil {
			osc.Amplitude.Modulation.Initialize()
		}
	}

	w.SignalFunc = w.makeSignalFunc()
}

func (w *WaveTable) makeSignalFunc() SignalFunc {
	wAmp := 1.0

	signalFunc := func(x float64) float64 {
		var y float64

		for i := range w.Oscillators {
			osc := w.Oscillators[i]
			freq := osc.Freq
			arg := x
			amp := 1.0

			if osc.Amplitude != nil {
				amp = osc.Amplitude.Value

				if osc.Amplitude.Modulation != nil {
					amp += osc.Amplitude.Modulation.SignalFunc(x)
				}
			}

			if osc.Phase != nil {
				arg += osc.Phase.Value

				if osc.Phase.Modulation != nil {
					arg += osc.Phase.Modulation.SignalFunc(x)
				}
			}

			for j := range w.Filters {
				filter := w.Filters[j]
				amp *= filter.Apply(freq, x)
			}

			if amp < 0 {
				amp = 0
			}

			y += osc.signalFunc(arg*freq) * amp
		}

		if w.Amplitude != nil {
			wAmp = w.Amplitude.Value

			if w.Amplitude.Modulation != nil {
				wAmp += w.Amplitude.Modulation.SignalFunc(x)
			}
		}

		return y * wAmp / float64(len(w.Oscillators))
	}

	return signalFunc
}
