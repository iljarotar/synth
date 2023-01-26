package parser

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func Parse(out interface{}) error {
	data, err := ioutil.ReadFile("patches/fm.yaml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, out)
	if err != nil {
		return err
	}

	return nil
}
