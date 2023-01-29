package audio

import (
	"github.com/gordonklaus/portaudio"
	"github.com/iljarotar/synth/config"
)

type ProcessCallback func([]float32)

type Context struct {
	*portaudio.Stream
	Input chan float32
}

func NewContext() (*Context, error) {
	ctx := &Context{Input: make(chan float32)}
	c := config.Instance()

	var err error
	ctx.Stream, err = portaudio.OpenDefaultStream(0, 1, c.SampleRate, 0, ctx.Process)
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
