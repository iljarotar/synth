package loader

import (
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/iljarotar/synth/control"
	"github.com/iljarotar/synth/screen"
	"github.com/iljarotar/synth/synth"
	s "github.com/iljarotar/synth/synth"
	"gopkg.in/yaml.v2"
)

type Loader struct {
	currentFile string
	watcher     *fsnotify.Watcher
	watch       *bool
	lastLoaded  time.Time
	ctl         *control.Control
	logger      *screen.Logger
}

func NewLoader(ctl *control.Control, log *screen.Logger) (*Loader, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	watch := true
	l := Loader{watcher: watcher, watch: &watch, ctl: ctl, logger: log}
	go l.StartWatching()

	return &l, nil
}

func (l *Loader) Close() error {
	*l.watch = false
	return l.watcher.Close()
}

func (l *Loader) Load(file string, synth *s.Synth) error {
	// to prevent clipping when write event is sent twice for the same change
	if time.Now().Sub(l.lastLoaded) < 500*time.Millisecond {
		return nil
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = l.Watch(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, synth)
	if err != nil {
		return err
	}

	l.lastLoaded = time.Now()
	l.currentFile = file
	l.ctl.LoadSynth(*synth)

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

			if !event.Has(fsnotify.Rename) {
				var s synth.Synth

				err := l.Load(l.currentFile, &s)
				if err != nil {
					l.logger.Log("could not load file. error: " + err.Error())
				}
			}
		case err, ok := <-l.watcher.Errors:
			if !ok {
				return
			}

			l.logger.Log("could not load file. error: " + err.Error())
		}
	}
}
