package main

import (
	"AnimeGUI/src/anilist"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

func setDialogAddAnime() {
	animesNames := binding.BindStringList(
		&[]string{},
	)

	listAnimeDisplay := widget.NewListWithData(animesNames,
		func() fyne.CanvasObject {
			return &widget.Label{Text: "template"}
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})
	listContainer := container.NewPadded(&canvas.Rectangle{FillColor: color.RGBA{R: grayScaleList, G: grayScaleList, B: grayScaleList, A: 255}, CornerRadius: 10}, listAnimeDisplay)

	inputSearch := widget.NewEntry()
	inputSearch.SetPlaceHolder("Search")
	inputSearch.OnSubmitted = func(s string) {
		if s == "" {
			return
		}
		fmt.Println(s)
		result := anilist.SearchFromQuery(inputSearch.Text)
		if result == nil {
			return
		}
		animesNames.Set([]string{})
		for i := 0; i < len(result); i++ {
			name := anilist.AnimeToName(&result[i])
			if name != nil {
				animesNames.Append(*name)
			}
		}

		fmt.Printf("Result: %+v\n", result)
	}
	buttonSearch := widget.NewButtonWithIcon("", theme.SearchIcon(), func() { inputSearch.OnSubmitted(inputSearch.Text) })
	buttonSearch.Importance = widget.WarningImportance
	searchBar := container.NewBorder(nil, nil, nil, buttonSearch, inputSearch)

	dialogAdd := dialog.NewCustom("Add new anime from Anilist", "Cancel", container.NewBorder(searchBar, nil, nil, nil, listContainer), window)
	dialogAdd.Resize(fyne.NewSize(800, 550))
	dialogAdd.Show()
}
