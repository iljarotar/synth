package module

import "math/rand"

type (
	Noise struct {
		Module
	}

	NoiseMap map[string]*Noise
)

func (n *Noise) Step() {
	val := rand.Float64()*2 - 1

	n.current = Output{
		Mono:  val,
		Left:  val / 2,
		Right: val / 2,
	}
}
