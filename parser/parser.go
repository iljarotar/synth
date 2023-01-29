package parser

import (
	"io/ioutil"

	s "github.com/iljarotar/synth/synth"
	"gopkg.in/yaml.v2"
)

type Parser struct {
	lastOpened, RootPath string
}

func NewParser() *Parser {
	return &Parser{RootPath: "examples"}
}

func (p *Parser) Load(file string, synth *s.Synth) error {
	data, err := ioutil.ReadFile(p.RootPath + "/" + file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, synth)
	if err != nil {
		return err
	}

	p.lastOpened = file

	return nil
}

func (p *Parser) LoadLastOpened(synth *s.Synth) error {
	return p.Load(p.lastOpened, synth)
}
