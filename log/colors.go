package log

type Color string

const (
	ColorWhite        Color = "\033[0m"
	ColorOrangeStrong Color = "\033[1;33m"
	ColorBlueStrong   Color = "\033[1;34m"
	ColorGreenStrong  Color = "\033[1;32m"
	ColorRedStrong    Color = "\033[1;31m"
)
