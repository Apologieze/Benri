package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var playMoreDialog *dialog.CustomDialog

func initPlayMoreDialog() {
	checkSkipOpening := widget.NewCheck("", func(b bool) {})
	checkSkipEnding := widget.NewCheck("", func(b bool) {})
	checkTrayIcon := widget.NewCheck("", func(b bool) {
	})

	rowSkipOpening := container.New(layout.NewFormLayout(),
		widget.NewLabelWithStyle("Automatically skip Opening", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), checkSkipOpening,
		widget.NewLabelWithStyle("Automatically skip Ending", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), checkSkipEnding,
		widget.NewLabelWithStyle("Add a tray icon for Benri", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), checkTrayIcon,
	)
	//form := container.New(layout.NewFormLayout(), rowSkipOpening)
	menuOption := container.NewBorder(nil, nil, nil, nil, rowSkipOpening)
	playMoreDialog = dialog.NewCustomWithoutButtons("More Actions", menuOption, window)
	playMoreDialog.Resize(fyne.NewSize(200, 300))
}

func openPlayMoreDialog() {
	playMoreDialog.Show()
}
