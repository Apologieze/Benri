package main

import (
	"fyne.io/fyne/v2/dialog"
	"unicode"
)

func showErrorPopUp(err error) {
	errorText := capitalizeFirstChar(err.Error())
	dialog.ShowInformation("Error", errorText, window)
}

func capitalizeFirstChar(s string) string {
	if len(s) == 0 {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
