package main

import (
	"AnimeGUI/src/config"
	"AnimeGUI/src/richPresence"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var dialogMenuOption *dialog.CustomDialog

func initSettingDialog() {
	checkSkipOpening := widget.NewCheck("", func(b bool) { config.SetBool(config.SkipOpeningKey, b) })
	checkSkipEnding := widget.NewCheck("", func(b bool) { config.SetBool(config.SkipEndingKey, b) })
	checkTrayIcon := widget.NewCheck("", func(b bool) {
		config.SetBool(config.TrayIconKey, b)
		fmt.Println(config.Setting.TrayIcon)
	})
	checkDiscordPresence := widget.NewCheck("", func(b bool) {
		config.SetBool(config.DiscordPresence, b)
		richPresence.InitDiscordRichPresence()
	})

	checkSkipOpening.SetChecked(config.Setting.SkipOpening)
	checkSkipEnding.SetChecked(config.Setting.SkipEnding)
	checkTrayIcon.SetChecked(config.Setting.TrayIcon)
	checkDiscordPresence.SetChecked(config.Setting.DiscordPresence)
	toggleTrayFeature()

	rowSkipOpening := container.New(layout.NewFormLayout(),
		widget.NewLabelWithStyle("Automatically skip Opening", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), checkSkipOpening,
		widget.NewLabelWithStyle("Automatically skip Ending", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), checkSkipEnding,
		widget.NewLabelWithStyle("Add a tray icon for Benri", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), checkTrayIcon,
		widget.NewLabelWithStyle("Show Discord Activity", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), checkDiscordPresence,
	)
	//form := container.New(layout.NewFormLayout(), rowSkipOpening)
	menuOption := container.NewBorder(nil, nil, nil, nil, rowSkipOpening)
	dialogMenuOption = dialog.NewCustom("Settings", "Close Settings", menuOption, window)
	dialogMenuOption.Resize(fyne.NewSize(200, 300))
}

func openMenuOption() {
	dialogMenuOption.Show()
}

func toggleTrayFeature() {
	if config.Setting.TrayIcon {
		if desk, ok := appW.(desktop.App); ok {
			m := fyne.NewMenu("MyApp",
				fyne.NewMenuItem("Show", func() {
					window.Show()
					window.RequestFocus()
				}))
			desk.SetSystemTrayMenu(m)

			window.SetCloseIntercept(func() {
				window.Hide()
			})
		}
	}
}
