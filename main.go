package main

import (
	"animeFyne/anilist"
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

var animeList []verniy.MediaList

func main() {
	startCurdInteg()
	fmt.Println(localAnime)
	a := app.New()
	//a.Settings().SetTheme(&myTheme{})
	var animeSelected verniy.MediaList

	w := a.NewWindow("AnimeGui")
	w.Resize(fyne.NewSize(1000, 700))

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

	dialogC := dialog.NewCustom("Select the correct anime", "Cancel", container.NewCenter(widget.NewLabel("Salut")), w)

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

	// Load image from file
	/*imgFile, err := os.Open("asset/img.png")
	anilist.Fatal(err)
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	anilist.Fatal(err)*/

	imageEx := &canvas.Image{}

	animeName := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})

	episodeNumber := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	episodeMinus := widget.NewButton(" - ", func() {})
	episodePlus := widget.NewButton(" + ", func() {})

	episodeContainer := container.NewHBox(layout.NewSpacer(), episodeMinus, episodeNumber, episodePlus, layout.NewSpacer())

	button := widget.NewButtonWithIcon("Play!", theme.MediaPlayIcon(), func() {
		//fmt.Println(anilist.Search())
		fmt.Println(animeSelected.Media.ID)
		if animeName.Text == "" {
			return
		}
		OnPlayButtonClick(animeName.Text, animeSelected)
		dialogC.Show()
	})

	button.IconPlacement = widget.ButtonIconTrailingText
	button.Importance = widget.HighImportance

	playContainer := container.NewHBox(layout.NewSpacer(), button, layout.NewSpacer())

	imageContainer := container.NewVBox(imageEx, animeName, episodeContainer, layout.NewSpacer(), playContainer)

	listDisplay.OnSelected = func(id int) {
		listName, err := data.GetValue(id)
		animeSelected = animeList[id]
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

	w.SetContent(container.NewBorder(nil, nil, nil, imageContainer, leftSide))

	w.ShowAndRun()
}

func updateAnimeNames(data binding.ExternalStringList) (first bool) {
	if animeList == nil {
		log.Error("No list found")
		return
	}
	tempName := make([]string, 0, 25)
	for _, anime := range animeList {
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
