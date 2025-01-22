package main

import (
	curd "animeFyne/curdInteg"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/charmbracelet/log"
	"net/url"
	"strings"
)

func setTokenGraphicaly(tokenPath string, user *curd.User) {
	changedToken = true
	fmt.Println(tokenPath, "Token path")

	window.SetTitle("Set Anilist Token")
	//window.Resize(fyne.NewSize(700, 400))

	/*window.SetOnClosed(func() {

	})*/

	urlLink, _ := url.Parse("https://anilist.co/api/v2/oauth/authorize?client_id=20686&response_type=token")
	err := appW.OpenURL(urlLink)
	if err != nil {
		log.Error("Error opening url", err)
	}
	hyperlink := widget.NewHyperlink("Website to generate token", urlLink)
	hyperlink.Alignment = fyne.TextAlignCenter

	labelTitle := widget.NewLabelWithStyle("Generate a token and paste it below", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	input := widget.NewEntry()
	input.SetPlaceHolder("Anilist token")
	input.OnSubmitted = func(s string) {
		fmt.Println(s)
		if changeTokenCustom(s, tokenPath, user) {
			initMainApp()
		}
	}

	validateButton := widget.NewButton("Validate", func() {
		if changeTokenCustom(input.Text, tokenPath, user) {
			initMainApp()
		}
	})
	validateButton.Importance = widget.HighImportance
	centerBtnContainer := container.NewHBox(layout.NewSpacer(), validateButton, layout.NewSpacer())

	window.SetContent(container.NewVBox(labelTitle, hyperlink, input, centerBtnContainer))
	window.Show()
	appW.Run()
}

func processToken(t string) string {
	token := strings.TrimSpace(t)
	if len(token) < 20 {
		return ""
	}
	return token
}

func changeTokenCustom(t, tokenPath string, user *curd.User) bool {
	token := processToken(t)
	if token == "" {
		return false
	}
	user.Token = token
	err := curd.WriteTokenToFile(token, tokenPath)
	if err != nil {
		return false
	}
	return true
}
