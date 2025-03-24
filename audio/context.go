package audio

import (
	"io"
	"math"

	"github.com/ebitengine/oto/v3"
)

type AudioOutput struct {
	Left, Right float64
}

type Audio struct {
	ctx    *oto.Context
	player *oto.Player
}

type reader struct {
	output chan AudioOutput
}

func NewAudio(output chan AudioOutput, sampleRate int) (*Audio, error) {
	op := &oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: 2,
		Format:       oto.FormatFloat32LE,
		BufferSize:   0,
	}

	ctx, ready, err := oto.NewContext(op)
	if err != nil {
		return nil, err
	}
	<-ready

	r := &reader{
		output: output,
	}
	player := ctx.NewPlayer(r)
	player.Play()

	audio := &Audio{
		ctx:    ctx,
		player: player,
	}

	return audio, nil
}

func (a *Audio) Close() error {
	return a.player.Close()
}

func (r *reader) Read(buf []byte) (int, error) {
	var n int
	for i := 0; i < len(buf)/8; i++ {
		select {
		case y := <-r.output:
			leftBytes := math.Float32bits(float32(y.Left))
			rightBytes := math.Float32bits(float32(y.Right))
			sampleIdx := 8 * i

			buf[sampleIdx] = byte(leftBytes)
			buf[sampleIdx+1] = byte(leftBytes >> 4)
			buf[sampleIdx+2] = byte(leftBytes >> 8)
			buf[sampleIdx+3] = byte(leftBytes >> 12)

			buf[sampleIdx+4] = byte(rightBytes)
			buf[sampleIdx+5] = byte(rightBytes >> 4)
			buf[sampleIdx+6] = byte(rightBytes >> 8)
			buf[sampleIdx+7] = byte(rightBytes >> 12)

			n += 8
		default:
		}
	}

	return n, io.EOF
}
