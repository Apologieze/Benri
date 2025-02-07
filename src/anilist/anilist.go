package anilist

import (
	"AnimeGUI/verniy"
	"fmt"
	"fyne.io/fyne/v2/widget"
	"github.com/charmbracelet/log"
	"net/url"
	"strings"
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
			verniy.MediaFieldNextAiringEpisode(
				verniy.AiringScheduleFieldEpisode,
				verniy.AiringScheduleFieldAiringAt,
				verniy.AiringScheduleFieldTimeUntilAiring,
			),
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

var categoriesToInt = make(map[string]int)

/*var categoriesToInt = map[string]int{
	"Completed": 0,
	"Dropped":   1,
	"Watching":  2,
	"Planning":  3,
}*/

var Client *verniy.Client = verniy.New()

func GetData(radio *widget.RadioGroup, username string, delete func()) {
	typeAnime, err := Client.GetUserAnimeListSort(username, verniy.MediaListSortUpdatedTimeDesc, fields...)
	if err != nil {
		typeAnime = make([]verniy.MediaListGroup, 4)
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

func FindList(categoryName string) *[]verniy.MediaList {
	if UserData == nil {
		log.Error("No data found")
		return nil
	}
	categoryIndex, exists := categoriesToInt[categoryName]
	if !exists {
		log.Error("Category not found in user")
		return &[]verniy.MediaList{}
	}
	return &UserData[categoryIndex].Entries
}

func AnimeToName(anime *verniy.Media) *string {
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

func AnimeToRomaji(anime *verniy.Media) string {
	if anime == nil {
		return ""
	}
	if anime.Title == nil {
		return ""
	}
	if anime.Title.Romaji != nil {
		return *anime.Title.Romaji
	}
	return ""
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

func SearchFromQuery(strQuery string) []verniy.Media {
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
