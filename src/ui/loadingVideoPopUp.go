package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/charmbracelet/log"
	"time"
)

var loadingVideoPopup *dialog.CustomDialog
var step int
var mainWindow fyne.Window

func InitLoadingVideoPopUp(window fyne.Window) {
	//content := container.NewVBox(widget.NewLabel("Loading episode"))
	/*content := widget.NewProgressBarInfinite()
	loadingVideoPopup = dialog.NewCustomWithoutButtons("Loading Episode", content, window)
	loadingVideoPopup.Resize(fyne.NewSize(400, 50))*/
	mainWindow = window
}

func ChangeLoadingStep(value int) {
	step = value
}

func ShowLoadingVideoPopUp(title string) {
	content := widget.NewProgressBarInfinite()
	loadingVideoPopup = dialog.NewCustomWithoutButtons(fmt.Sprint("Loading ", title), content, mainWindow)
	loadingVideoPopup.Resize(fyne.NewSize(400, 50))
	step = 1
	if loadingVideoPopup == nil {
		log.Error("Couldn't open loading popup")
		return
	}
	fmt.Println("popup")
	loadingVideoPopup.Show()
}

func CloseLoadingPopup(duration time.Duration) {
	if loadingVideoPopup == nil {
		log.Error("Couldn't close loading popup")
		return
	}
	go func() {
		time.Sleep(time.Second * duration)
		loadingVideoPopup.Hide()
	}()
}
