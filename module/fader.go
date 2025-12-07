package module

type (
	fader struct {
		current float64
		target  float64
		step    float64
	}
)

func (f *fader) initialize(duration, sampleRate float64) {
	delta := f.target - f.current
	if duration == 0 || sampleRate == 0 {
		f.step = delta
	}
	f.step = delta / (duration * sampleRate)
}

func (f *fader) fade() float64 {
	if f.current == f.target {
		return f.current
	}

	new := f.current + f.step
	if (f.target-f.current >= 0) != (f.target-new >= 0) {
		f.current = f.target
	}
	f.current = new
	return f.current
}
