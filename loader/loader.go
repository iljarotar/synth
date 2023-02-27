package loader

import (
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/iljarotar/synth/control"
	"github.com/iljarotar/synth/screen"
	s "github.com/iljarotar/synth/synth"
	"gopkg.in/yaml.v2"
)

type Loader struct {
	watcher    *fsnotify.Watcher
	watch      *bool
	lastLoaded time.Time
	ctl        *control.Control
	logger     *screen.Logger
	file       string
}

func NewLoader(ctl *control.Control, log *screen.Logger, file string) (*Loader, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	watch := true
	l := Loader{watcher: watcher, watch: &watch, ctl: ctl, logger: log, file: file}
	go l.StartWatching()

	err = l.Watch(l.file)
	if err != nil {
		return nil, err
	}

	return &l, nil
}

func (l *Loader) Close() error {
	*l.watch = false
	return l.watcher.Close()
}

func (l *Loader) Load() error {
	data, err := os.ReadFile(l.file)
	if err != nil {
		return err
	}

	var synth s.Synth
	err = yaml.Unmarshal(data, &synth)
	if err != nil {
		return err
	}

	l.ctl.Stop(0.01)
	l.ctl.LoadSynth(synth)
	l.ctl.Start(0.01)

	l.lastLoaded = time.Now()
	return nil
}

func (l *Loader) Watch(file string) error {
	filePath, err := filepath.Abs(file)
	if err != nil {
		return err
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

			time.Sleep(time.Millisecond * 50) // to prevent occasional empty file loading

			// check last loaded time to prevent occasional double loading
			if !event.Has(fsnotify.Rename) && time.Now().Sub(l.lastLoaded) > 500*time.Millisecond {
				err := l.Load()
				if err != nil {
					l.logger.Log("could not load file. error: " + err.Error())
				}
			}
		case err, ok := <-l.watcher.Errors:
			if !ok {
				return
			}
			l.logger.Log("an error occurred. please restart synth. error: " + err.Error())
		}
	}
}
