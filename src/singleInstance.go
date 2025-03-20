package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/fsnotify/fsnotify"
	"github.com/nightlyone/lockfile"
	"os"
	"path/filepath"
)

var (
	lockPath   = filepath.Join(os.TempDir(), AppName+".lock")
	notifyPath = filepath.Join(os.TempDir(), AppName+".notify")
)

func initLock() lockfile.Lockfile {
	lock, err := lockfile.New(lockPath)
	if err != nil {
		fmt.Printf("Cannot init lockfile: %v\n", err)
		os.Exit(1)
	}

	// Try to acquire the lock
	err = lock.TryLock()
	if err != nil {
		fmt.Println("Application is already running")

		// Touch the notification file to signal the first instance
		_ = os.WriteFile(notifyPath, []byte{}, 0644)

		os.Exit(0)
	}
	return lock
}

func newAppDetection() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error("Error creating watcher", "error", err)
		return
	}
	defer watcher.Close()

	notifyDir := filepath.Dir(notifyPath)

	err = watcher.Add(notifyDir)
	if err != nil {
		log.Error("Error watching directory", "error", err)
		return
	}

	log.Info("Watching for new instances", "path", notifyPath)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Name == notifyPath && (event.Op&fsnotify.Write != 0 || event.Op&fsnotify.Create != 0) {
				log.Printf("%+v", event)
				log.Info("New instance detected, bringing to front")
				window.Show()
				window.RequestFocus()
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Error("Watcher error", "error", err)
		}
	}
}
