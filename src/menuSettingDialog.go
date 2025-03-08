package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var dialogMenuOption *dialog.CustomDialog

func initMenuOption() {
	rowSkipOpening := container.New(layout.NewFormLayout(),
		widget.NewLabelWithStyle("Automatically skip Opening", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewCheck("", func(b bool) { fmt.Println(b) }),
		widget.NewLabelWithStyle("Automatically skip Ending", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewCheck("", func(b bool) { fmt.Println(b) }),
		widget.NewLabelWithStyle("Add a tray icon for Benri", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewCheck("", func(b bool) { fmt.Println(b) }),
	)
	//form := container.New(layout.NewFormLayout(), rowSkipOpening)
	menuOption := container.NewBorder(nil, nil, nil, nil, rowSkipOpening)
	dialogMenuOption = dialog.NewCustom("Settings", "Close menu", menuOption, window)
	dialogMenuOption.Resize(fyne.NewSize(200, 300))
}

func openMenuOption() {
	dialogMenuOption.Show()
}
