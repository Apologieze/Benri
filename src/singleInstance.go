package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/nightlyone/lockfile"
	"os"
	"time"
)

var (
	lockPath   = os.TempDir() + "/" + AppName + ".lock"
	notifyPath = os.TempDir() + "/" + AppName + ".notify"
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
		currentTime := []byte(fmt.Sprintf("%d", time.Now().UnixNano()))
		_ = os.WriteFile(notifyPath, currentTime, 0644)

		os.Exit(0)
	}
	return lock
}

func pollingNewAppDetection() {
	var lastModTime time.Time

	info, err := os.Stat(notifyPath)
	if err == nil {
		lastModTime = info.ModTime()
	}

	for {
		// Check if notification file exists
		info, err := os.Stat(notifyPath)

		if err == nil && info.ModTime().After(lastModTime) {
			log.Info("New instance detected, bringing to front")
			window.Show()
			window.RequestFocus()
			lastModTime = info.ModTime()
		}

		time.Sleep(500 * time.Millisecond)
	}
}
