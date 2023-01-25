package dsp

func LowpassFilter(signal []float32) []float32 {
	filtered := make([]float32, len(signal))
	for i := range signal {
		if i == 0 || i == len(signal)-1 {
			filtered[i] = signal[i]
			continue
		}
		filtered[i] = (signal[i-1] + signal[i] + signal[i+1]) / 3
	}
	return filtered
}
