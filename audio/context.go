package audio

import (
	"github.com/gordonklaus/portaudio"
)

type ProcessCallback func([]float32)

type Context struct {
	*portaudio.Stream
	Input chan float32
}

func NewContext(sampleRate float64, input chan float32) (*Context, error) {
	ctx := &Context{Input: input}

	var err error
	ctx.Stream, err = portaudio.OpenDefaultStream(0, 1, sampleRate, 0, ctx.Process)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func (c *Context) Process(out []float32) {
	for i := range out {
		out[i] = <-c.Input
	}
}

func Init() error {
	return portaudio.Initialize()
}

func Terminate() error {
	return portaudio.Terminate()
}
