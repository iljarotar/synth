package audio

import (
	"time"

	"github.com/ebitengine/oto/v3"
)

const (
	format         = oto.FormatFloat32LE
	bytesPerSample = 8
	bufferSize     = 512
)

type Context struct {
	ctx    *oto.Context
	player *oto.Player
}

func NewContext(sampleRate int, readSample func() [2]float64) (*Context, error) {
	ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: 2,
		Format:       format,
		BufferSize:   bufferDuration(bufferSize, float64(sampleRate)),
	})

	if err != nil {
		return nil, err
	}
	<-ready

	sampleReader := &reader{
		readSample: readSample,
	}

	player := ctx.NewPlayer(sampleReader)
	player.SetBufferSize(bufferSize * bytesPerSample)
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

func bufferDuration(bufferSize, sampleRate float64) time.Duration {
	return time.Duration(float64(time.Second) * bufferSize / sampleRate)
}
