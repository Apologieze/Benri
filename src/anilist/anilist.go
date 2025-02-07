package anilist

import (
	"AnimeGUI/verniy"
	verniy2 "AnimeGUI/verniy"
	"fmt"
	"fyne.io/fyne/v2/widget"
	"github.com/charmbracelet/log"
	"net/url"
	"strings"
)

var fields = []verniy2.MediaListGroupField{
	verniy2.MediaListGroupFieldName,
	verniy2.MediaListGroupFieldEntries(
		verniy2.MediaListFieldID,
		verniy2.MediaListFieldStatus,
		verniy2.MediaListFieldScore,
		verniy2.MediaListFieldProgress,
		verniy2.MediaListFieldMedia(
			verniy2.MediaFieldID,
			verniy2.MediaFieldNextAiringEpisode(
				verniy2.AiringScheduleFieldEpisode,
				verniy2.AiringScheduleFieldAiringAt,
				verniy2.AiringScheduleFieldTimeUntilAiring,
			),
			verniy2.MediaFieldTitle(
				verniy2.MediaTitleFieldRomaji,
				verniy2.MediaTitleFieldEnglish,
				verniy2.MediaTitleFieldNative),
			verniy2.MediaFieldType,
			verniy2.MediaFieldFormat,
			verniy2.MediaFieldStatusV2,
			verniy2.MediaFieldCoverImage(verniy2.MediaCoverImageFieldLarge, verniy2.MediaCoverImageFieldExtraLarge),
			verniy2.MediaFieldAverageScore,
			verniy2.MediaFieldPopularity,
			verniy2.MediaFieldIsAdult,
			verniy2.MediaFieldEpisodes)),
}

var UserData []verniy2.MediaListGroup

var categoriesToInt = make(map[string]int)

/*var categoriesToInt = map[string]int{
	"Completed": 0,
	"Dropped":   1,
	"Watching":  2,
	"Planning":  3,
}*/

var Client *verniy2.Client = verniy2.New()

func GetData(radio *widget.RadioGroup, username string, delete func()) {
	typeAnime, err := Client.GetUserAnimeListSort(username, verniy.MediaListSortUpdatedTimeDesc, fields...)
	if err != nil {
		typeAnime = make([]verniy2.MediaListGroup, 4)
		log.Error("Invalid token")
		log.Error(err)
		delete()
	}

	categoriesToInt = make(map[string]int)
	for i := 0; i < len(typeAnime); i++ {
		if typeAnime[i].Name != nil {
			categoriesToInt[*typeAnime[i].Name] = i
			//typeAnime[i].Entries =
		}
	}

	UserData = typeAnime
	if radio != nil {
		if radio.Selected == "" {
			radio.SetSelected("Watching")
		}
	}
}

func FindList(categoryName string) *[]verniy2.MediaList {
	if UserData == nil {
		log.Error("No data found")
		return nil
	}
	categoryIndex, exists := categoriesToInt[categoryName]
	if !exists {
		log.Error("Category not found in user")
		return &[]verniy2.MediaList{}
	}
	return &UserData[categoryIndex].Entries
}

func AnimeToName(anime *verniy2.Media) *string {
	if anime == nil {
		return nil
	}
	if anime.Title == nil {
		return nil
	}
	if anime.Title.English != nil {
		return anime.Title.English
	}
	return anime.Title.Romaji
}

func FormatDuration(seconds int) string {
	days := seconds / 86400
	remaining := seconds % 86400
	hours := remaining / 3600
	remaining %= 3600
	minutes := remaining / 60

	var parts []string

	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d days", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	// Include minutes if it's non-zero or no parts exist (e.g., 0d 0h 0m → "0m")
	if minutes > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}

	return strings.Join(parts, " ")
}

func IdToUrl(id int) *url.URL {
	u, err := url.Parse(fmt.Sprintf("https://anilist.co/anime/%d", id))
	if err != nil {
		return nil
	}
	return u
}

func SearchFromQuery(strQuery string) []verniy2.Media {
	if strQuery == "" {
		return nil
	}
	query := verniy.PageParamMedia{Search: strQuery}
	result, err := Client.SearchAnime(query, 1, 15)
	if err != nil {
		log.Error("Error searching anime:", err)
		return nil
	}
	return result.Media
}
