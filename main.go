package main

import (
	"animeFyne/anilist"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/bep/debounce"
	"github.com/charmbracelet/log"
	"github.com/rl404/verniy"
	"image"
	"os"
	"time"
)

var animeList []verniy.MediaList

func main() {
	a := app.New()
	a.Settings().SetTheme(&myTheme{})
	w := a.NewWindow("Hello")
	w.Resize(fyne.NewSize(900, 600))

	debounced := debounce.New(400 * time.Millisecond)

	data := binding.BindStringList(
		&[]string{},
	)

	listDisplay := widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	listDisplay.OnSelected = func(id int) {
		//log.Infof("Selected: %d", id)
		fmt.Println(*animeList[id].Media.CoverImage.ExtraLarge)
	}

	hello := widget.NewLabel("Hello Fyne!")
	button := widget.NewButton("Hi!", func() {
		//fmt.Println(anilist.Search())
		err := data.Set([]string{"feur"})
		if err != nil {
			return
		}
	})

	input := widget.NewEntry()
	input.SetPlaceHolder("Anime name")
	input.OnChanged = func(s string) {
		debounced(func() {
			hello.SetText(s)
		})
	}

	radiobox := widget.NewRadioGroup([]string{"Watching", "Planning", "Completed", "Dropped"}, func(s string) {
		animeList = anilist.FindList(s)
		updateName(data)
	})
	radiobox.Required = true
	radiobox.Horizontal = true

	vbox := container.NewVBox(
		button,
		input,
		radiobox,
	)

	leftSide := container.NewBorder(vbox, nil, nil, nil, listDisplay)

	go anilist.GetData(radiobox)

	// Load image from file
	imgFile, err := os.Open("asset/img.png")
	anilist.Fatal(err)
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	anilist.Fatal(err)

	size := img.Bounds().Size()
	ratio := float32(size.X) / float32(size.Y)
	newWidth := 400
	newHeight := float32(newWidth) / ratio

	imageEx := canvas.NewImageFromImage(img)
	imageEx.FillMode = canvas.ImageFillContain

	imageEx.SetMinSize(fyne.NewSize(400, newHeight))

	imageContainer := container.NewVBox(imageEx, layout.NewSpacer())

	w.SetContent(container.NewBorder(nil, nil, nil, imageContainer, leftSide))

	w.ShowAndRun()
}

func updateName(data binding.ExternalStringList) {
	if animeList == nil {
		log.Error("No list found")
		return
	}
	tempName := make([]string, 0, 25)
	for _, anime := range animeList {
		if media := anime.Media; media != nil {
			if title := media.Title; title != nil {
				if english := title.English; english != nil {
					//fmt.Println(*english)
					tempName = append(tempName, *english)
				} else {
					//log.Error(*title.Romaji)
					tempName = append(tempName, *title.Romaji)
				}
			}
		}
	}
	err := data.Set(tempName)
	if err != nil {
		return
	}
}
