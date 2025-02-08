package main

import (
	curd "AnimeGUI/curdInteg"
	"AnimeGUI/src/anilist"
	"AnimeGUI/verniy"
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
	fynetooltip "github.com/dweymouth/fyne-tooltip"
	ttwidget "github.com/dweymouth/fyne-tooltip/widget"
	"image"
	"image/color"
	"io"
	"net/http"
	"net/url"
	"time"
)

var (
	animeList           *[]verniy.MediaList
	animeSelected       *verniy.MediaList
	window              fyne.Window
	appW                fyne.App
	episodeNumber       = widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	episodeLastPlayback = widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})
	changedToken        bool
	mpvPresent          bool
	grayScaleList       uint8 = 35
)

func main() {
	const AppName = "AnimeGUI"

	go dowloadMPV()

	appW = app.New()
	window = appW.NewWindow(AppName)
	window.Resize(fyne.NewSize(1000, 700))
	window.CenterOnScreen()
	window.Show()

	appW.Settings().SetTheme(&forcedVariant{
		Theme:   theme.DefaultTheme(),
		variant: theme.VariantDark,
	})
	//log.Info("Color", appW.Settings().Theme().Color(theme.ColorNameFocus, theme.VariantDark))

	startCurdInteg()
	if !changedToken {
		fmt.Println(window.Title(), AppName)
		initMainApp()
		appW.Run()
	}
}

func updateAnimeNames(data binding.ExternalStringList) (first bool) {
	if animeList == nil {
		log.Error("No list found")
		return
	}
	tempName := make([]string, 0, 25)
	for _, anime := range *animeList {
		if name := anilist.AnimeToName(anime.Media); name != nil {
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

func getAnimeImageFromImage(img image.Image, newWidth float32) *canvas.Image {
	size := img.Bounds().Size()
	ratio := float32(size.X) / float32(size.Y)

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
		fmt.Print("Selected:", allAnimeList[index])
		dialogC.Hide()
		err, tempAnime := curd.LocalUpdateAnime(databaseFile, animeSelected.Media.ID, allAnimeList[index].Id, animeProgress, 0, 0, animeName)
		if err != nil {
			log.Error("Can't update database file", err)
			return
		}
		if tempAnime != nil {
			localAnime = tempAnime
		}
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

func initMainApp() {
	secondCurdInit()
	anilist.Client.AccessToken = user.Token
	window.SetTitle("Benri")
	fmt.Println(localAnime)

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

	input := widget.NewEntry()
	input.SetPlaceHolder("Filter anime name")
	input.OnChanged = func(s string) {
		debounced(func() {
			fmt.Println(s)
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

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			setDialogAddAnime()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.MailAttachmentIcon(), func() {
			if animeSelected == nil {
				_ = appW.OpenURL(&url.URL{Scheme: "https", Host: "anilist.co", Path: "search/anime"})
				return
			}
			urlId := anilist.IdToUrl(animeSelected.Media.ID)
			if urlId != nil {
				err := appW.OpenURL(urlId)
				if err != nil {
					log.Error("Can't open url", err)
				}
			}
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.MenuIcon(), func() {
			openMenuOption()
		}),
	)

	inputContainer := container.NewBorder(nil, nil, nil, toolbar, input)

	vbox := container.NewVBox(
		inputContainer,
		radiobox,
	)

	/*if themeVariant == theme.VariantDark {
		grayScaleList = 220
	}*/
	listContainer := container.NewPadded(&canvas.Rectangle{FillColor: color.RGBA{R: grayScaleList, G: grayScaleList, B: grayScaleList, A: 255}, CornerRadius: 10}, listDisplay)

	leftSide := container.NewBorder(vbox, nil, nil, nil, listContainer)

	go anilist.GetData(radiobox, user.Username, deleteTokenFile)

	imageEx := &canvas.Image{}

	animeName := ttwidget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	animeName.Wrapping = fyne.TextWrapWord

	episodeMinus := widget.NewButton(" - ", func() { changeEpisodeInApp(-1) })
	episodePlus := widget.NewButton(" + ", func() { changeEpisodeInApp(1) })

	episodeContainer := container.NewHBox(layout.NewSpacer(), episodeMinus, episodeNumber, episodePlus, layout.NewSpacer())

	//nextEpisodeLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	nextEpisodeLabel := &canvas.Text{Text: "", Color: color.RGBA{156, 190, 93, 255}, Alignment: fyne.TextAlignCenter, TextStyle: fyne.TextStyle{Bold: true}, TextSize: theme.TextSize()}
	nextEpisodeLabel.Hide()

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

	imageContainer := container.NewVBox(imageEx, animeName, episodeContainer, nextEpisodeLabel, episodeLastPlayback, layout.NewSpacer(), playContainer)

	listDisplay.OnSelected = func(id int) {
		listName, err := data.GetValue(id)
		animeSelected = &(*animeList)[id]
		if err == nil {
			animeName.SetText(listName)
			animeName.SetToolTip(anilist.AnimeToRomaji(animeSelected.Media))
		}

		if animeSelected.Media.NextAiringEpisode != nil {
			nextEpisodeLabel.Text = fmt.Sprintf("Episode %d releasing in %s", animeSelected.Media.NextAiringEpisode.Episode, anilist.FormatDuration(animeSelected.Media.NextAiringEpisode.TimeUntilAiring))
			nextEpisodeLabel.Show()
			nextEpisodeLabel.Refresh()
		} else {
			nextEpisodeLabel.Hide()
		}

		if animeSelected.Progress != nil && animeSelected.Media.Episodes != nil {
			episodeNumber.SetText(fmt.Sprintf("Episode %d/%d", *animeSelected.Progress, *animeSelected.Media.Episodes))
		} else {
			episodeNumber.SetText("No episode data")
		}
		displayLocalProgress()

		imageLink := *animeSelected.Media.CoverImage.ExtraLarge

		imageFile := GetImageFromUrl(imageLink)
		if imageFile == nil {
			log.Error("No image found")
			return
		}

		*imageEx = *getAnimeImageFromImage(imageFile, 300)
		imageContainer.Refresh()
	}

	window.SetContent(fynetooltip.AddWindowToolTipLayer(container.NewBorder(nil, nil, nil, imageContainer, leftSide), window.Canvas()))
}
