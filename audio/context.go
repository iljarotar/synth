package audio

import "github.com/gordonklaus/portaudio"

type ProcessCallback func([]float32)

type Context struct {
	*portaudio.Stream
}

func NewContext(sampleRate float64, callback ProcessCallback) (*Context, error) {
	ctx := &Context{}

	var err error
	ctx.Stream, err = portaudio.OpenDefaultStream(0, 1, sampleRate, 0, callback)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func (c *Context) Start() {
	err := c.Stream.Start()
	if err != nil {
		panic(err)
	}
}

func (c *Context) Stop() {
	err := c.Stream.Stop()
	if err != nil {
		panic(err)
	}
}

func Init() {
	portaudio.Initialize()
}

func Terminate() {
	portaudio.Terminate()
}
