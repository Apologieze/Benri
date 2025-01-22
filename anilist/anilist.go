package anilist

import (
	"fyne.io/fyne/v2/widget"
	"github.com/charmbracelet/log"
	"github.com/rl404/verniy"
)

var fields = []verniy.MediaListGroupField{
	verniy.MediaListGroupFieldName,
	verniy.MediaListGroupFieldEntries(
		verniy.MediaListFieldID,
		verniy.MediaListFieldStatus,
		verniy.MediaListFieldScore,
		verniy.MediaListFieldProgress,
		verniy.MediaListFieldMedia(
			verniy.MediaFieldID,
			verniy.MediaFieldTitle(
				verniy.MediaTitleFieldRomaji,
				verniy.MediaTitleFieldEnglish,
				verniy.MediaTitleFieldNative),
			verniy.MediaFieldType,
			verniy.MediaFieldFormat,
			verniy.MediaFieldStatusV2,
			verniy.MediaFieldCoverImage(verniy.MediaCoverImageFieldLarge, verniy.MediaCoverImageFieldExtraLarge),
			verniy.MediaFieldAverageScore,
			verniy.MediaFieldPopularity,
			verniy.MediaFieldIsAdult,
			verniy.MediaFieldEpisodes)),
}

var UserData []verniy.MediaListGroup

var categoriesToInt = map[string]int{
	"Completed": 0,
	"Dropped":   1,
	"Watching":  2,
	"Planning":  3,
}

func GetData(radio *widget.RadioGroup, username string, delete func()) {
	v := verniy.New()

	typeAnime, err := v.GetUserAnimeList(username, fields...)
	if err != nil {
		typeAnime = make([]verniy.MediaListGroup, 4)
		log.Error("Invalid token")
		delete()
	}
	UserData = typeAnime
	if radio != nil {
		if radio.Selected == "" {
			radio.SetSelected("Watching")
		}
	}
}

func FindList(categoryName string) *[]verniy.MediaList {
	if UserData == nil {
		log.Error("No data found")
		return nil
	}
	categoryIndex := categoriesToInt[categoryName]
	return &UserData[categoryIndex].Entries
}

func AnimeToName(anime verniy.MediaList) *string {
	if anime.Media == nil {
		return nil
	}
	if anime.Media.Title == nil {
		return nil
	}
	if anime.Media.Title.English != nil {
		return anime.Media.Title.English
	}
	return anime.Media.Title.Romaji
}

/*func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}*/
