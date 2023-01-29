package parser

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/iljarotar/synth/config"
	s "github.com/iljarotar/synth/synth"
	"gopkg.in/yaml.v2"
)

type Parser struct {
	lastOpened string
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) SetRootPath(path string) error {
	c := config.Instance()

	if strings.Split(path, "/")[0] == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		c.SetRootPath(home + path[1:])
	} else {
		c.SetRootPath(path)
	}

	return nil
}

func (p *Parser) Load(file string, synth *s.Synth) error {
	c := config.Instance()
	data, err := ioutil.ReadFile(c.RootPath + "/" + file)
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
	if p.lastOpened == "" {
		return errors.New("don't know what to apply. please load a file first")
	}

	return p.Load(p.lastOpened, synth)
}
