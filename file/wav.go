package file

import (
	"os"
	"strings"
)

type Recorder struct {
	in, out chan struct{ Left, Right float32 }
	record  bool
	buffer  [2][]float32
	file    string
}

func NewRecorder(in, out chan struct{ Left, Right float32 }, file string) Recorder {
	return Recorder{in: in, out: out, record: true, file: file}
}

func (r *Recorder) StartRecording() {
	defer close(r.out)
	for y := range r.in {
		r.buffer[0] = append(r.buffer[0], y.Left)
		r.buffer[1] = append(r.buffer[1], y.Right)
		r.out <- y
	}
}

func (r *Recorder) StopRecording() error {
	if r.file == "" {
		return nil
	}

	if !strings.HasSuffix(r.file, ".wav") {
		r.file = r.file + ".wav"
	}

	f, err := os.Create(r.file)
	if err != nil {
		return err
	}
	defer f.Close()

	// TODO: write wav data here

	return nil
}
