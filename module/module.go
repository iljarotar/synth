package module

type IModule interface {
	Current() Output
	Integral() float64
}

type ModulesMap map[string]IModule

type Module struct {
	current  Output
	integral float64
}

func (m *Module) Current() Output {
	return m.current
}

func (m *Module) Integral() float64 {
	return m.integral
}

type limits [2]float64

var (
	gainLimits   = limits{0, 1}
	outputLimits = limits{-1, 1}
	freqLimits   = limits{0, 20000}
)

type Output struct {
	Mono, Left, Right float64
}
