package file

import (
	"encoding/binary"
	"math"
	"os"
	"strings"
)

// file format
var (
	// header
	chunkID   = []byte{0x52, 0x49, 0x46, 0x46} // RIFF
	chunkSize = []byte{}                       // 36 + subchunk2Size
	format    = []byte{0x57, 0x41, 0x56, 0x45} // WAVE

	// subchunk1
	subchunk1ID   = []byte{0x66, 0x6d, 0x74, 0x20} // FMT
	subchunk1Size = []byte{0x10, 0x00, 0x00, 0x00} // 16
	audioFormat   = []byte{0x01, 0x00}             // 1
	numChannels   = []byte{0x02, 0x00}             // 2
	sampleRate    = []byte{}
	byteRate      = []byte{}           // sampleRate * numChannels * bitsPerSample/8
	blockAlign    = []byte{0x04, 0x00} // numChannels * bitsPerSample/8
	bitsPerSample = []byte{0x10, 0x00} // 16

	// subchunk2
	subchunk2ID   = []byte{0x64, 0x61, 0x74, 0x61} // DATA
	subchunk2Size = []byte{}                       // numSamples * numChannels * bitsPerSample/8
	data          = []byte{}
)

func writeWavFile(file string, sRate int, samples []sample) {
	if !strings.HasSuffix(file, ".wav") {
		file += ".wav"
	}

	f, _ := os.Create(file)
	b := make([]byte, 0)

	channels := 2
	bits := 16

	h := header(len(samples) * 4) // numSamples * numChannels * bitsPerSample/8
	ch1 := subchunk1(sRate, channels, bits)
	ch2 := subchunk2(channels, bits, samples)

	b = append(b, h...)
	b = append(b, ch1...)
	b = append(b, ch2...)

	f.Write(b)
	f.Close()
}

func header(size int) []byte {
	b := make([]byte, 0)
	s := 36 + size
	chunkSize = intToBytes(s, 4)

	b = append(b, chunkID...)
	b = append(b, chunkSize...)
	b = append(b, format...)

	return b
}

func subchunk1(sRate, channels, bits int) []byte {
	b := make([]byte, 0)
	sampleRate = intToBytes(sRate, 4)
	bRate := sRate * channels * bits / 8
	byteRate = intToBytes(bRate, 4)

	b = append(b, subchunk1ID...)
	b = append(b, subchunk1Size...)
	b = append(b, audioFormat...)
	b = append(b, numChannels...)
	b = append(b, sampleRate...)
	b = append(b, byteRate...)
	b = append(b, blockAlign...)
	b = append(b, bitsPerSample...)

	return b
}

func subchunk2(channels, bits int, samples []sample) []byte {
	b := make([]byte, 0)
	b = append(b, subchunk2ID...)

	s := len(samples) * channels * bits / 8
	subchunk2Size = intToBytes(s, 4)
	b = append(b, subchunk2Size...)
	b = append(b, getRawData(samples)...)

	return b
}

func getRawData(samples []sample) []byte {
	data := make([]byte, 0)

	for _, val := range samples {
		l := floatToBytes(val[0])
		r := floatToBytes(val[1])
		data = append(data, l...)
		data = append(data, r...)
	}

	return data
}

// returns little endian representation of n
func intToBytes(n, num int) []byte {
	b := make([]byte, 4)
	in := uint32(n)
	binary.LittleEndian.PutUint32(b, in)
	return b[:num]
}

// returns a 16 bit little endian representation of x
func floatToBytes(x float32) []byte {
	y := x * float32(math.MaxInt16)
	return intToBytes(int(y), 2)
}
