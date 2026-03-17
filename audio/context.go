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

func Start(sampleRate int, readSample func() [2]float64) error {
	ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: 2,
		Format:       format,
		BufferSize:   bufferDuration(bufferSize, float64(sampleRate)),
	})

	if err != nil {
		return err
	}
	<-ready

	sampleReader := &reader{
		readSample: readSample,
	}

	player := ctx.NewPlayer(sampleReader)
	player.SetBufferSize(bufferSize * bytesPerSample)
	player.Play()
	return nil
}

func bufferDuration(bufferSize, sampleRate float64) time.Duration {
	return time.Duration(float64(time.Second) * bufferSize / sampleRate)
}
