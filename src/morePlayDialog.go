package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var playMoreDialog *widget.PopUp

func initPlayMoreDialog() {
	/*checkSkipOpening := widget.NewCheck("", func(b bool) {})
	checkSkipEnding := widget.NewCheck("", func(b bool) {})
	checkTrayIcon := widget.NewCheck("", func(b bool) {
	})*/

	playPreviousButton := widget.NewButtonWithIcon("Play previous episode", theme.Icon(theme.IconNameMediaPlay), func() {

	})
	playPreviousButton.IconPlacement = widget.ButtonIconTrailingText

	/*rowSkipOpening := container.New(layout.NewFormLayout(),
		widget.NewLabelWithStyle("Automatically skip Opening", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), checkSkipOpening,
		widget.NewLabelWithStyle("Automatically skip Ending", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), checkSkipEnding,
		widget.NewLabelWithStyle("Add a tray icon for Benri", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), checkTrayIcon,
	)*/
	vbox := container.NewVBox(layout.NewSpacer(), playPreviousButton, layout.NewSpacer())

	menuOption := container.NewBorder(nil, nil, nil, nil, container.NewPadded(vbox))
	playMoreDialog = widget.NewPopUp(menuOption, window.Canvas())
	//playMoreDialog.Resize(fyne.NewSize(200, 100))
}

func openPlayMoreDialog() {
	playMoreDialog.ShowAtRelativePosition(fyne.Position{X: 100}, episodeLastPlayback)
}
