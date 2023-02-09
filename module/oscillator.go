package module

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

type Oscillators map[string]*Oscillator

type Oscillator struct {
	Name    string         `yaml:"name"`
	Type    OscillatorType `yaml:"type"`
	Freq    float64        `yaml:"freq"`
	Amp     *Param         `yaml:"amp"`
	Phase   Param          `yaml:"phase"`
	Filters []string       `yaml:"filters"`
	Signal  SignalFunc
	Current float64
}

func (o *Oscillator) Initialize() {
	o.Signal = NewSignalFunc(o.Type)
	amp := 1.0

	if o.Amp != nil {
		amp = o.Amp.Val
	}

	o.Current = o.Signal(o.Phase.Val) * amp
}
