package module

type (
	Sampler struct {
		Module
		In           string `yaml:"in"`
		Trigger      string `yaml:"trigger"`
		triggerValue float64
	}

	SamplerMap map[string]*Sampler
)

func (s *Sampler) Update(new *Sampler) {
	if new == nil {
		return
	}

	s.In = new.In
	s.Trigger = new.Trigger
}

func (s *Sampler) Step(modules *ModuleMap) {
	triggerValue := getMono(modules, s.Trigger)

	if triggerValue > 0 && s.triggerValue <= 0 {
		val := getMono(modules, s.In)

		s.current = Output{
			Mono:  val,
			Left:  val / 2,
			Right: val / 2,
		}
	}

	s.triggerValue = triggerValue
}
