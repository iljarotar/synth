package audio

import (
	"github.com/ebitengine/oto/v3"
)

type Context struct {
	ctx    *oto.Context
	player *oto.Player
}

func NewContext(sampleRate int, readSample func() [2]float64) (*Context, error) {
	ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: 2,
		Format:       oto.FormatSignedInt16LE,
		BufferSize:   0,
	})

	if err != nil {
		return nil, err
	}
	<-ready

	sampleReader := &reader{
		readSample: readSample,
	}

	player := ctx.NewPlayer(sampleReader)
	// TODO: set buffer size?
	player.Play()

	context := &Context{
		ctx:    ctx,
		player: player,
	}

	return context, nil
}

func (a *Context) Close() error {
	return a.player.Close()
}
