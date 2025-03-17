package main

import (
	"AnimeGUI/src/anilist"
	"AnimeGUI/verniy"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/charmbracelet/log"
	"image/color"
)

var displayToCategories = map[string]string{
	"Watching":      "CURRENT",
	"Completed":     "COMPLETED",
	"Paused":        "PAUSED",
	"Dropped":       "DROPPED",
	"Plan to Watch": "PLANNING",
}

var categoriesToDisplay = map[string]string{
	"CURRENT":   "Watching",
	"COMPLETED": "Completed",
	"PAUSED":    "Paused",
	"DROPPED":   "Dropped",
	"PLANNING":  "Plan to Watch",
}

var displayCategories = []string{
	"Watching",
	"Plan to Watch",
	"Completed",
	"Paused",
	"Dropped",
}

func setDialogAddAnime() {
	var searchResult []verniy.Media
	var selectedAnime *verniy.Media
	var selectedAnimeStatus string
	isAnimeSelected := binding.NewBool()

	animeImageHolder := &canvas.Image{}
	selectCategory := widget.NewSelect(displayCategories, func(s string) {})
	selectCategory.Alignment = fyne.TextAlignCenter
	selectCategory.SetSelected(displayCategories[0])
	labelInfo := &widget.Label{Text: "Change to category:", Alignment: fyne.TextAlignCenter}
	imageContainer := container.NewVBox(animeImageHolder, labelInfo, selectCategory, layout.NewSpacer())
	//imageContainer.Hide()

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
		isAnimeSelected.Set(false)
		selectedAnime = nil
		listAnimeDisplay.UnselectAll()
		if s == "" {
			return
		}
		fmt.Println(s)
		result := anilist.SearchFromQuery(inputSearch.Text)
		if result == nil {
			return
		}
		searchResult = result
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

	dialogAdd := dialog.NewCustomWithoutButtons("Add new anime from Anilist", container.NewBorder(searchBar, nil, nil, imageContainer, listContainer), window)

	addButton := &widget.Button{Text: "Add", OnTapped: dialogAdd.Hide, Icon: theme.ConfirmIcon(), Importance: widget.HighImportance}

	addButton.OnTapped = func() {
		if selectedAnime == nil {
			return
		}
		err := anilist.UpdateAnimeStatus(user.Token, selectedAnime.ID, displayToCategories[selectCategory.Selected])
		if err != nil {
			log.Error("Error updating anime status:", err)
			return
		}
		log.Info("Anime added successfully")
		go anilist.GetData(categoryRadiobox, user.Username, func() { log.Error("Invalid token") })
		dialogAdd.Hide()
	}

	selectCategory.OnChanged = func(s string) {
		if status, ok := displayToCategories[s]; ok {
			if status == selectedAnimeStatus {
				addButton.Disable()
				return
			}
		}
		addButton.Enable()
	}

	dialogAdd.SetButtons([]fyne.CanvasObject{
		&widget.Button{Text: "Cancel", OnTapped: dialogAdd.Hide, Icon: theme.CancelIcon()},
		addButton,
	})
	dialogAdd.Resize(fyne.NewSize(850, 580))

	listAnimeDisplay.OnSelected = func(id int) {

		imageLink := searchResult[id].CoverImage.Large
		isAnimeSelected.Set(true)
		selectedAnime = &searchResult[id]

		if imageLink != nil {
			imageFile := GetImageFromUrl(*imageLink)
			if imageFile == nil {
				log.Error("No image found")
				return
			}

			*animeImageHolder = *getAnimeImageFromImage(imageFile, 220)
			imageContainer.Refresh()
		}

		status := anilist.FindStatusFromId(selectedAnime.ID)
		if status != nil {
			selectedAnimeStatus = string(*status)
			if display, ok := categoriesToDisplay[selectedAnimeStatus]; ok {
				selectCategory.SetSelected(display)
			}
		} else {
			selectedAnimeStatus = ""
			selectCategory.SetSelected(displayCategories[0])
		}
	}

	isAnimeSelected.AddListener(binding.NewDataListener(func() {
		get, err := isAnimeSelected.Get()
		if err != nil {
			return
		}
		log.Info("Is anime selected:", get)
		if get {
			addButton.Show()
			imageContainer.Show()
		} else {
			addButton.Hide()
			imageContainer.Hide()
		}
	}))

	dialogAdd.Show()
	window.Canvas().Focus(inputSearch)
}
