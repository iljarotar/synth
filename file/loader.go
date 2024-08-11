package file

import (
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/iljarotar/synth/control"
	s "github.com/iljarotar/synth/synth"
	"github.com/iljarotar/synth/ui"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v2"
)

type Loader struct {
	logger     *ui.Logger
	watcher    *fsnotify.Watcher
	watch      *bool
	lastLoaded time.Time
	ctl        *control.Control
	file       string
}

func NewLoader(logger *ui.Logger, ctl *control.Control, file string) (*Loader, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	watch := true
	l := &Loader{
		logger:  logger,
		watcher: watcher,
		watch:   &watch,
		ctl:     ctl,
		file:    file,
	}
	go l.StartWatching()

	return l, nil
}

func (l *Loader) Close() error {
	*l.watch = false
	return l.watcher.Close()
}

func (l *Loader) Load() error {
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

	l.ctl.LoadSynth(synth)

	l.lastLoaded = time.Now()
	return nil
}

func (l *Loader) Watch(file string) error {
	filePath, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	if slices.Contains(l.watcher.WatchList(), file) {
		l.watcher.Remove(file)
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
			if ui.State.Closed {
				return
			}

			time.Sleep(time.Millisecond * 50) // to prevent occasional empty file loading

			// check last loaded time to prevent occasional double loading
			if !event.Has(fsnotify.Rename) && time.Now().Sub(l.lastLoaded) > 500*time.Millisecond {
				waitForFadeOut := make(chan bool)

				l.ctl.FadeOut(0.01, waitForFadeOut)
				<-waitForFadeOut

				err := l.Load()
				if err != nil {
					l.logger.Error("could not load file. error: " + err.Error())
				} else {
					l.logger.Info("reloaded patch file")
					l.logger.ShowOverdriveWarning(false)
				}

				l.ctl.FadeIn(0.01)
			}
		case err, ok := <-l.watcher.Errors:
			if !ok {
				return
			}
			l.logger.Error("an error occurred. please restart synth. error: " + err.Error())
		}
	}
}
