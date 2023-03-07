package file

import "os"

var (
	// header
	chunkID   = []byte{0x52, 0x49, 0x46, 0x46} // RIFF
	chunkSize = []byte{0x00, 0x00, 0x00, 0x28} // just for testing, normally: 36 + subchunk2Size
	format    = []byte{0x57, 0x41, 0x56, 0x45} // WAVE

	// subchunk1
	subchunk1ID   = []byte{0x66, 0x6d, 0x74, 0x20} // FMT
	subchunk1Size = []byte{0x10, 0x00, 0x00, 0x00} // 16
	audioFormat   = []byte{0x01, 0x00}             // 1
	numChannels   = []byte{0x02, 0x00}             // 2
	sampleRate    = []byte{0x22, 0x56, 0x00, 0x00} // just for testing, should come from outside
	bitsPerSample = []byte{0x10, 0x00}             // 16
	byteRate      = []byte{0x88, 0x58, 0x01, 0x00} // just for testing, normally: sampleRate * numChannels * bitsPerSample/8
	blockAlign    = []byte{0x04, 0x00}             // numChannels * bitsPerSample/8

	// subchunk2
	subchunk2ID   = []byte{0x64, 0x61, 0x74, 0x61} // DATA
	subchunk2Size = []byte{0x00, 0x08, 0x00, 0x00} // numSamples * numChannels * bitsPerSample/8
	data          = []byte{
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	} // just testing
)

func test() {
	f, _ := os.Create("test.wav")
	b := make([]byte, 0)

	b = append(b, chunkID...)
	b = append(b, chunkSize...)
	b = append(b, format...)
	b = append(b, subchunk1ID...)
	b = append(b, subchunk1Size...)
	b = append(b, audioFormat...)
	b = append(b, numChannels...)
	b = append(b, sampleRate...)
	b = append(b, bitsPerSample...)
	b = append(b, byteRate...)
	b = append(b, blockAlign...)
	b = append(b, subchunk2ID...)
	b = append(b, subchunk2Size...)
	b = append(b, data...)

	f.Write(b)
	f.Close()
}
