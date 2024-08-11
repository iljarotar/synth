package audio

import (
	"github.com/gordonklaus/portaudio"
)

type AudioOutput struct {
	Left, Right float64
}

type Context struct {
	*portaudio.Stream
	Input chan AudioOutput
}

func NewContext(input chan AudioOutput, sampleRate float64) (*Context, error) {
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
		out[0][i] = float32(y.Left)
		out[1][i] = float32(y.Right)
	}
}

func Init() error {
	return portaudio.Initialize()
}

func Terminate() error {
	return portaudio.Terminate()
}
