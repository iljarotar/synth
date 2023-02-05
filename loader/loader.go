package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/control"
	"github.com/iljarotar/synth/synth"
	s "github.com/iljarotar/synth/synth"
	"gopkg.in/yaml.v2"
)

type Loader struct {
	lastOpened string
	watcher    *fsnotify.Watcher
	watch      *bool
	control    *control.Control
}

func NewLoader(ctl *control.Control) (*Loader, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	watch := true
	l := Loader{watcher: watcher, watch: &watch, control: ctl}
	go l.StartWatching()
	return &l, nil
}

func (l *Loader) Close() error {
	*l.watch = false
	return l.watcher.Close()
}

func (l *Loader) SetRootPath(path string) error {
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

func (l *Loader) Load(file string, synth *s.Synth) error {
	fmt.Println("enter")
	c := config.Instance()
	data, err := ioutil.ReadFile(c.RootPath + "/" + file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, synth)
	if err != nil {
		return err
	}

	err = l.Watch(file)
	if err != nil {
		return err
	}

	l.control.LoadSynth(*synth)
	l.lastOpened = file
	return nil
}

func (l *Loader) Watch(file string) error {
	c := config.Instance()
	filePath := c.RootPath + "/" + file

	if l.lastOpened != "" {
		lastOpenedPath := c.RootPath + "/" + l.lastOpened

		err := l.watcher.Remove(lastOpenedPath)
		if err != nil {
			return err
		}
	}

	return l.watcher.Add(filePath)
}

func (l *Loader) StartWatching() {
	for *l.watch {
		select {
		case event, ok := <-l.watcher.Events:
			if !ok {
				return
			}

			if event.Has(fsnotify.Write) {
				var s synth.Synth

				err := l.Load(l.lastOpened, &s)
				if err != nil {
					fmt.Println("could not load file. error: " + err.Error())
				}
			}
		case err, ok := <-l.watcher.Errors:
			if !ok {
				return
			}

			fmt.Println("could not load file. error: " + err.Error())
		}
	}
}
