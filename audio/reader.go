package audio

import (
	"fmt"
	"math"
)

type reader struct {
	readSample func() [2]float64
}

func (r *reader) Read(buf []byte) (int, error) {
	if len(buf)%int(bytesPerSample) != 0 {
		return 0, fmt.Errorf("buffer lenght must be divisible by %d", bytesPerSample)
	}
	numSamples := len(buf) / int(bytesPerSample)

	for i := range buf[:numSamples] {
		s := r.readSample()

		left := math.Float32bits(float32(s[0]))
		right := math.Float32bits(float32(s[1]))

		buf[i*bytesPerSample] = byte(left)
		buf[i*bytesPerSample+1] = byte(left >> 8)
		buf[i*bytesPerSample+2] = byte(left >> 16)
		buf[i*bytesPerSample+3] = byte(left >> 24)

		buf[i*bytesPerSample+4] = byte(right)
		buf[i*bytesPerSample+5] = byte(right >> 8)
		buf[i*bytesPerSample+6] = byte(right >> 16)
		buf[i*bytesPerSample+7] = byte(right >> 24)
	}

	return len(buf), nil
}
