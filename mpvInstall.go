package main

import (
	"github.com/charmbracelet/log"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func dowloadMPV() {
	log.Info(runtime.GOOS)
	if runtime.GOOS != "windows" {
		return
	}
	exePath, err := os.Executable()
	if err != nil {
		log.Error(err)
		return
	}
	exeDir := filepath.Dir(exePath)
	mpvPath := filepath.Join(exeDir, "bin", "mpv.exe")
	log.Info(mpvPath)
	if _, err := os.Stat(mpvPath); os.IsNotExist(err) {
		log.Info("mpv.exe does not exist")
	} else if err != nil {
		log.Error(err)
	} else {
		return
	}
	log.Info("Downloading mpv")
	resp, err := http.Get("http://apologize.fr/mpv.exe")
	if err != nil {
		log.Error("Error downloading mpv.exe:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	mpvDir := filepath.Join(exeDir, "bin")
	if err := os.MkdirAll(mpvDir, os.ModePerm); err != nil {
		log.Error("Error creating directories:", err)
		return
	}

	out, err := os.Create(mpvPath)
	if err != nil {
		log.Error("Error creating mpv.exe file:", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Error("Error saving mpv.exe:", err)
		return
	}
	log.Info("mpv.exe downloaded successfully")
}
