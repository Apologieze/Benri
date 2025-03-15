package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var playMorePopUp *widget.PopUp

func initPlayMorePopUp() {
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
	playMorePopUp = widget.NewPopUp(menuOption, window.Canvas())
	playMorePopUp.Resize(fyne.NewSize(290, 100))
}

func openPlayMorePopUp() {
	//sizeW := window.Canvas().Size()
	//sizePopup := playMorePopUp.Size()
	//playMorePopUp.ShowAtPosition(fyne.Position{sizeW.Width - sizePopup.Width - 50, sizeW.Height - sizePopup.Height - 50})
	playMorePopUp.ShowAtRelativePosition(fyne.Position{X: 0, Y: -110}, playButton)
}
