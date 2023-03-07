package file

import (
	"encoding/binary"
	"math"
	"os"
	"strings"
)

type waveWriter struct {
	file waveFile
}

type waveFile struct {
	header waveHeader
	format waveFormat
	data   waveData
	bytes  []byte
}

type waveHeader struct {
	chunkID   []byte
	chunkSize []byte
	format    []byte
	bytes     []byte
}

type waveFormat struct {
	subchunkID    []byte
	subchunkSize  []byte
	audioFormat   []byte
	numChannels   int
	sampleRate    int
	byteRate      []byte
	blockAlign    []byte
	bitsPerSample int
	bytes         []byte
}

type waveData struct {
	subchunkID   []byte
	subchunkSize int
	data         []byte
	bytes        []byte
}

func newWaveWriter(samples []sample, sampleRate int) waveWriter {
	writer := waveWriter{}
	f := waveFile{}

	// the order of these calls is important
	f.getFormat(sampleRate)
	f.getData(samples)
	f.getHeader()
	f.getBytes()

	writer.file = f
	return writer
}

func (w *waveWriter) write(file string) error {
	if !strings.HasSuffix(file, ".wav") {
		file += ".wav"
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	b := make([]byte, 0)
	b = append(b, w.file.bytes...)

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func (f *waveFile) getBytes() {
	b := make([]byte, 0)

	b = append(b, f.header.bytes...)
	b = append(b, f.format.bytes...)
	b = append(b, f.data.bytes...)

	f.bytes = b
}

func (f *waveFile) getHeader() {
	b := make([]byte, 0)

	f.header.chunkID = []byte{0x52, 0x49, 0x46, 0x46} // RIFF
	f.header.format = []byte{0x57, 0x41, 0x56, 0x45}  // WAVE
	size := 36 + f.data.subchunkSize
	f.header.chunkSize = intToBytes(size, 4)

	b = append(b, f.header.chunkID...)
	b = append(b, f.header.chunkSize...)
	b = append(b, f.header.format...)

	f.header.bytes = b
}

func (f *waveFile) getFormat(sampleRate int) {
	b := make([]byte, 0)

	f.format.subchunkID = []byte{0x66, 0x6d, 0x74, 0x20}   // FMT
	f.format.subchunkSize = []byte{0x10, 0x00, 0x00, 0x00} // 16
	f.format.audioFormat = []byte{0x01, 0x00}              // 1
	f.format.numChannels = 2
	f.format.sampleRate = sampleRate

	bytes := f.format.sampleRate * f.format.numChannels * f.format.bitsPerSample / 8
	f.format.byteRate = intToBytes(bytes, 4)
	f.format.blockAlign = []byte{0x04, 0x00} // 4
	f.format.bitsPerSample = 16

	sampleRateBytes := intToBytes(f.format.sampleRate, 4)
	bitsBytes := intToBytes(f.format.bitsPerSample, 2)
	channelsBytes := intToBytes(f.format.numChannels, 2)

	b = append(b, f.format.subchunkID...)
	b = append(b, f.format.subchunkSize...)
	b = append(b, f.format.audioFormat...)
	b = append(b, channelsBytes...)
	b = append(b, sampleRateBytes...)
	b = append(b, f.format.byteRate...)
	b = append(b, f.format.blockAlign...)
	b = append(b, bitsBytes...)

	f.format.bytes = b
}

func (f *waveFile) getData(samples []sample) {
	b := make([]byte, 0)

	f.data.subchunkID = []byte{0x64, 0x61, 0x74, 0x61} // DATA
	f.data.subchunkSize = len(samples) * f.format.numChannels * f.format.bitsPerSample / 8
	f.data.data = getRawData(samples)
	size := intToBytes(f.data.subchunkSize, 4)

	b = append(b, f.data.subchunkID...)
	b = append(b, size...)
	b = append(b, getRawData(samples)...)

	f.data.bytes = b
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
