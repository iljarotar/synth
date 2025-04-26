package audio

import (
	"fmt"
)

type reader struct {
	readSample func() [2]float64
}

// Read func inspired by https://github.com/gopxl/beep
func (r *reader) Read(buf []byte) (int, error) {
	// TODO: make 4 a constant
	if len(buf)%4 != 0 {
		return 0, fmt.Errorf("buffer lenght must be divisible by 4")
	}

	for i := range buf[:len(buf)/4] {
		s := r.readSample()

		left := s[0]
		leftInt16 := int16(left * (1<<15 - 1))
		leftLow := byte(leftInt16)
		leftHigh := byte(leftInt16 >> 8)
		buf[i*4] = leftLow
		buf[i*4+1] = leftHigh

		right := s[1]
		rightInt16 := int16(right * (1<<15 - 1))
		rightLow := byte(rightInt16)
		rightHigh := byte(rightInt16 >> 8)
		buf[i*4+2] = rightLow
		buf[i*4+3] = rightHigh
	}

	return len(buf), nil
}
