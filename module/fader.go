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
		return
	}
	f.step = delta / (duration * sampleRate)
}

func (f *fader) fade() float64 {
	// technically, this case is covered below, but for efficiency we make an early return here
	if f.current == f.target {
		return f.current
	}

	f.current += f.step
	if (f.current > f.target) == (f.step > 0) {
		f.current = f.target
	}

	return f.current
}
