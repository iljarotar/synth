package audio

import (
	"github.com/gordonklaus/portaudio"
)

type ProcessCallback func([]float32)

type Context struct {
	*portaudio.Stream
	Input chan struct{ Left, Right float32 }
}

func NewContext(input chan struct{ Left, Right float32 }, sampleRate float64) (*Context, error) {
	ctx := &Context{Input: input}

	var err error
	ctx.Stream, err = portaudio.OpenDefaultStream(0, 2, sampleRate, 0, ctx.Process)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func (c *Context) Process(out [][]float32) {
	for i := range out[0] {
		y := <-c.Input
		out[0][i] = y.Left
		out[1][i] = y.Right
	}
}

func Init() error {
	return portaudio.Initialize()
}

func Terminate() error {
	return portaudio.Terminate()
}
