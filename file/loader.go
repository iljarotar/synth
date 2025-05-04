package file

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/iljarotar/synth/log"
	s "github.com/iljarotar/synth/synth"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v2"
)

type callbackFunc func(*s.Synth) error

type Loader struct {
	logger   *log.Logger
	file     string
	callback callbackFunc

	watcher    *fsnotify.Watcher
	lastLoaded time.Time
	active     bool
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

// FIX: saving file quickly multiple times breaks watcher (also occurs on main branch)
func (l *Loader) StartWatching() {
	for l.active {
		select {
		case event := <-l.watcher.Events:
			time.Sleep(time.Millisecond * 50) // to prevent occasional empty file loading

			// check last loaded time to prevent occasional double loading
			if !event.Has(fsnotify.Rename) && time.Now().Sub(l.lastLoaded) > 500*time.Millisecond {
				err := l.LoadAndWatch()
				if err != nil {
					l.logger.Error("could not load file. error: " + err.Error())
				} else {
					l.logger.Info("reloaded patch file")
				}
			}
		case err := <-l.watcher.Errors:
			l.logger.Error(fmt.Sprintf("failed to watch file:%v", err))
		}
	}
}
