package module

type Sampler struct {
	Module
	In           string `yaml:"in"`
	Trigger      string `yaml:"trigger"`
	triggerValue float64
}

type SamplerMap map[string]*Sampler

func (s *Sampler) Step(modules ModulesMap) {
	triggerValue := getMono(modules[s.Trigger])

	if triggerValue > 0 && s.triggerValue <= 0 {
		val := getMono(modules[s.In])

		s.current = Output{
			Mono:  val,
			Left:  val / 2,
			Right: val / 2,
		}
	}

	s.triggerValue = triggerValue
}
