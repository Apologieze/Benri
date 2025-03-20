package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/charmbracelet/log"
	"time"
)

var loadingVideoPopup *dialog.CustomDialog
var step int

func InitLoadingVideoPopUp(window fyne.Window) {
	content := container.NewVBox(widget.NewLabel("Loading episode"))
	loadingVideoPopup = dialog.NewCustomWithoutButtons("Loading episode", content, window)
	loadingVideoPopup.Resize(fyne.NewSize(400, 400))

}

func ChangeLoadingStep(value int) {
	step = value
}

func ShowLoadingVideoPopUp() {
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
