package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/iljarotar/synth/log"
	s "github.com/iljarotar/synth/synth"
	"gopkg.in/yaml.v2"
)

type callbackFunc func(*s.Synth) error

type Loader struct {
	logger   *log.Logger
	file     string
	callback callbackFunc

	watcher *fsnotify.Watcher
	active  bool
}

func NewLoader(logger *log.Logger, filename string, callback callbackFunc) (*Loader, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	l := &Loader{
		logger:   logger,
		file:     filename,
		callback: callback,
		watcher:  watcher,
		active:   true,
	}
	go l.StartWatching()

	return l, nil
}

func (l *Loader) Close() error {
	return l.watcher.Close()
}

func (l *Loader) Stop() {
	l.active = false
}

func (l *Loader) LoadAndWatch() error {
	err := l.Watch(l.file)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(l.file)
	if err != nil {
		return err
	}

	var synth s.Synth
	err = yaml.Unmarshal(data, &synth)
	if err != nil {
		return err
	}

	err = l.callback(&synth)
	if err != nil {
		return err
	}

	return nil
}

func (l *Loader) Watch(file string) error {
	absolutePath, err := filepath.Abs(file)
	if err != nil {
		return err
	}
	parentDir, _ := filepath.Split(absolutePath)

	return l.watcher.Add(parentDir)
}

func (l *Loader) StartWatching() {
	for l.active {
		select {
		case event := <-l.watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				err := l.LoadAndWatch()
				if err != nil {
					l.logger.Error(fmt.Sprintf("failed to load file:%v", err))
					continue
				}

				l.logger.Info("reloaded patch file")
			}
		case err := <-l.watcher.Errors:
			l.logger.Error(fmt.Sprintf("failed to watch file:%v", err))
		}
	}
}
