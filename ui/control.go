package ui

type control interface {
	DecreaseVolume()
	IncreaseVolume()
	Volume() float64
}
