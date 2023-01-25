package control

type buffer struct {
	Input, Ouput chan float32
	buffer, send []float32
}

func NewBuffer(output chan float32) buffer {
	input := make(chan float32)
	return buffer{Input: input, Ouput: output}
}

func (b *buffer) Pipe() {
	for {
		v := <-b.Input
		b.buffer = append(b.buffer, v)
		if len(b.buffer) > 1024 {
			b.send = b.buffer // filter b.buffer if necessary
			b.Ouput <- b.send[0]
			b.send = b.send[1:]
			b.buffer = b.buffer[1:]
		}
	}
}
