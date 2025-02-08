package main

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

func openMenuOption() {
	dialogMenuOption := dialog.NewCustom("Menu", "Go back", container.NewCenter(), window)
	dialogMenuOption.Show()
}
