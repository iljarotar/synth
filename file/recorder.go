package file

import c "github.com/iljarotar/synth/config"

type sample [2]float32

type Recorder struct {
	in, out chan struct{ Left, Right float32 }
	buffer  []sample
	file    string
}

func NewRecorder(in, out chan struct{ Left, Right float32 }, file string) Recorder {
	r := Recorder{in: in, out: out, file: file}
	if c.Config.Duration >= 0 {
		r.buffer = make([]sample, int(c.Config.SampleRate*(c.Config.Duration+c.Config.FadeIn+c.Config.FadeOut)))
	}
	return r
}

func (r *Recorder) StartRecording() {
	defer close(r.out)
	var i int
	for y := range r.in {
		if len(r.buffer) > i {
			r.buffer[i] = [2]float32{y.Left, y.Right}
		}
		r.out <- y
		i++
	}
}

func (r *Recorder) StopRecording() error {
	if r.file == "" {
		return nil
	}

	writer := newWaveWriter(r.buffer, int(c.Config.SampleRate))
	return writer.write(r.file)
}
