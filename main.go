package main

import (
	"animeFyne/anilist"
	curd "animeFyne/curdInteg"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/bep/debounce"
	"github.com/charmbracelet/log"
	"github.com/rl404/verniy"
	"image"
	"image/color"
	"io"
	"net/http"
	"time"
)

var animeList *[]verniy.MediaList
var window fyne.Window
var animeSelected *verniy.MediaList
var episodeNumber = widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

func main() {
	startCurdInteg()
	fmt.Println(localAnime)
	a := app.New()
	//a.Settings().SetTheme(&myTheme{})

	window = a.NewWindow("AnimeGui")
	window.Resize(fyne.NewSize(1000, 700))

	debounced := debounce.New(400 * time.Millisecond)

	data := binding.BindStringList(
		&[]string{},
	)

	listDisplay := widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return &widget.Label{Text: "template"}
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	hello := widget.NewLabel("Hello Fyne!")

	input := widget.NewEntry()
	input.SetPlaceHolder("Anime name")
	input.OnChanged = func(s string) {
		debounced(func() {
			hello.SetText(s)
		})
	}

	radiobox := widget.NewRadioGroup([]string{"Watching", "Planning", "Completed", "Dropped"}, func(s string) {
		animeList = anilist.FindList(s)
		if updateAnimeNames(data) {
			listDisplay.Unselect(0)
			listDisplay.Select(0)
			listDisplay.ScrollToTop()
		}
	})
	radiobox.Required = true
	radiobox.Horizontal = true

	vbox := container.NewVBox(
		input,
		radiobox,
	)

	var grayScaleList uint8 = 35
	listContainer := container.NewStack(canvas.NewRectangle(color.RGBA{R: grayScaleList, G: grayScaleList, B: grayScaleList, A: 255}), listDisplay)

	leftSide := container.NewBorder(vbox, nil, nil, nil, listContainer)

	go anilist.GetData(radiobox, user.Username)

	imageEx := &canvas.Image{}

	animeName := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})

	episodeMinus := widget.NewButton(" - ", func() { changeEpisodeInApp(-1) })
	episodePlus := widget.NewButton(" + ", func() { changeEpisodeInApp(1) })

	episodeContainer := container.NewHBox(layout.NewSpacer(), episodeMinus, episodeNumber, episodePlus, layout.NewSpacer())

	button := widget.NewButtonWithIcon("Play!", theme.MediaPlayIcon(), func() {
		//fmt.Println(anilist.Search())
		fmt.Println(animeSelected.Media.ID)
		if animeName.Text == "" {
			return
		}
		OnPlayButtonClick(animeName.Text, animeSelected)
	})

	button.IconPlacement = widget.ButtonIconTrailingText
	button.Importance = widget.HighImportance

	playContainer := container.NewHBox(layout.NewSpacer(), button, layout.NewSpacer())

	imageContainer := container.NewVBox(imageEx, animeName, episodeContainer, layout.NewSpacer(), playContainer)

	listDisplay.OnSelected = func(id int) {
		listName, err := data.GetValue(id)
		animeSelected = &(*animeList)[id]
		if err == nil {
			animeName.SetText(listName)
		}

		if animeSelected.Progress != nil && animeSelected.Media.Episodes != nil {
			episodeNumber.SetText(fmt.Sprintf("Episode %d/%d", *animeSelected.Progress, *animeSelected.Media.Episodes))
		} else {
			episodeNumber.SetText("No episode data")
		}

		imageLink := *animeSelected.Media.CoverImage.ExtraLarge

		imageFile := GetImageFromUrl(imageLink)
		if imageFile == nil {
			log.Error("No image found")
			return
		}

		*imageEx = *getAnimeImageFromImage(imageFile)
		imageContainer.Refresh()
	}

	window.SetContent(container.NewBorder(nil, nil, nil, imageContainer, leftSide))

	window.ShowAndRun()
}

func updateAnimeNames(data binding.ExternalStringList) (first bool) {
	if animeList == nil {
		log.Error("No list found")
		return
	}
	tempName := make([]string, 0, 25)
	for _, anime := range *animeList {
		if name := anilist.AnimeToName(anime); name != nil {
			tempName = append(tempName, *name)
		} else {
			tempName = append(tempName, "Error name")
		}
	}
	if len(tempName) != 0 {
		first = true
	}
	err := data.Set(tempName)
	if err != nil {
		return
	}
	return first
}

func GetImageFromUrl(url string) image.Image {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading image:", err)
		return nil
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return nil
	}
	return img
}

func getAnimeImageFromImage(img image.Image) *canvas.Image {
	size := img.Bounds().Size()
	ratio := float32(size.X) / float32(size.Y)
	var newWidth float32 = 300
	newHeight := newWidth / ratio

	imageEx := canvas.NewImageFromImage(img)
	imageEx.FillMode = canvas.ImageFillContain

	imageEx.SetMinSize(fyne.NewSize(newWidth, newHeight))

	return imageEx
}

func selectCorrectLinking(allAnimeList []AllAnimeIdData, animeName string, animeProgress int) {
	linkingList := widget.NewList(func() int {
		return len(allAnimeList)
	},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(allAnimeList[i].Name)
		})

	linkingContainer := container.NewBorder(nil, nil, nil, nil, linkingList)
	dialogC := dialog.NewCustom("Select the correct anime", "Cancel", linkingContainer, window)

	linkingList.OnSelected = func(index widget.ListItemID) {
		fmt.Println("Selected:", allAnimeList[index])
		dialogC.Hide()
		err := curd.LocalUpdateAnime(databaseFile, animeSelected.Media.ID, allAnimeList[index].Id, animeProgress, 0, 0, animeName)
		if err != nil {
			log.Error("Can't update database file", err)
			return
		}
		localAnime = curd.LocalGetAllAnime(databaseFile)
		OnPlayButtonClick(animeName, animeSelected)
	}

	dialogC.Resize(fyne.NewSize(600, 900))
	dialogC.Show()
	log.Info("Select the correct anime", allAnimeList)

}

func changeEpisodeInApp(variation int) {
	var currentSelected = animeSelected
	newNumber := *currentSelected.Progress + variation
	fmt.Println("New number:", newNumber, *currentSelected.Progress)
	if newNumber >= 0 && newNumber <= *currentSelected.Media.Episodes {
		go UpdateAnimeProgress(currentSelected.Media.ID, newNumber)
		currentSelected.Progress = &newNumber
		episodeNumber.SetText(fmt.Sprintf("Episode %d/%d", newNumber, *currentSelected.Media.Episodes))
	}
}
