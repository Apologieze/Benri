package main

import (
	curd "animeFyne/curdInteg"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/rl404/verniy"
	"os"
	"path/filepath"
	"runtime"
)

var localAnime []curd.Anime
var userCurdConfig curd.CurdConfig
var databaseFile string
var user curd.User

func startCurdInteg() {
	//discordClientId := "1287457464148820089"

	//var anime curd.Anime

	var homeDir string
	if runtime.GOOS == "windows" {
		homeDir = os.Getenv("USERPROFILE")
	} else {
		homeDir = os.Getenv("HOME")
	}

	configFilePath := filepath.Join(homeDir, ".config", "curd", "curd.conf")

	// load curd userCurdConfig
	var err error
	userCurdConfig, err = curd.LoadConfig(configFilePath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	curd.SetGlobalConfig(&userCurdConfig)

	logFile := filepath.Join(os.ExpandEnv(userCurdConfig.StoragePath), "debug.log")
	//curd.ClearLogFile(logFile)

	// Get the token from the token file
	user.Token, err = curd.GetTokenFromFile(filepath.Join(os.ExpandEnv(userCurdConfig.StoragePath), "token"))
	if err != nil {
		log.Error("Error reading token", logFile)
	}
	if user.Token == "" {
		curd.ChangeToken(&userCurdConfig, &user)
	}

	if user.Id == 0 {
		user.Id, user.Username, err = curd.GetAnilistUserID(user.Token)
		if err != nil {
			log.Error(err)
		}
	}

	databaseFile = filepath.Join(os.ExpandEnv(userCurdConfig.StoragePath), "curd_history.txt")
	localAnime = curd.LocalGetAllAnime(databaseFile)
	for _, anime := range localAnime {
		fmt.Println(anime)
	}
	/*fmt.Println(databaseFile)*/

	fmt.Println(user.Username)
	fmt.Println(user.Id)

	/*allId := localAnime[0].AllanimeId
	url, _ := curd.GetEpisodeURL(userCurdConfig, allId, 1)
	fmt.Println(curd.PrioritizeLink(url))*/
}

func SearchFromLocalAniId(id int) *curd.Anime {
	for _, anime := range localAnime {
		if anime.AnilistId == id {
			return &anime
		}
	}
	return nil
}

type AllAnimeIdData struct {
	Id   string
	Name string
}

var keyValueArray []AllAnimeIdData

func OnPlayButtonClick(animeName string, animeData verniy.MediaList) {
	var allAnimeId string
	animeProgress := 0
	if animeData.Progress != nil {
		animeProgress = *animeData.Progress
	}
	animePointer := SearchFromLocalAniId(animeData.Media.ID)
	if animePointer == nil {
		allAnimeId = searchAllAnimeData(animeName, animeData.Media.Episodes)
		if allAnimeId == "" {
			log.Error("Failed to get allAnimeId")
			return
		}
		err := curd.LocalUpdateAnime(databaseFile, animeData.Media.ID, allAnimeId, animeProgress, 0, 0, animeName)
		if err != nil {
			log.Error("Can't update database file", err)
			return
		} else {
			log.Info("Successfully updated database file")
		}
	} else {
		fmt.Println(*animePointer)
		allAnimeId = animePointer.AllanimeId
	}
	log.Info("AllAnimeId!!!!!:", allAnimeId)

	if animeProgress == 0 {
		animeProgress = 1
	}
	log.Info("Anime Progress:", animeProgress)

	url, err := curd.GetEpisodeURL(userCurdConfig, allAnimeId, animeProgress)
	if err != nil {
		log.Error(err)
	}
	finalLink := curd.PrioritizeLink(url)
	fmt.Println(finalLink)

	mpvSocketPath, err := curd.StartVideo(finalLink, []string{}, fmt.Sprintf("%s - Episode %d", animeName, animeProgress))
	fmt.Println("MPV Socket Path:", mpvSocketPath)
}

func searchAllAnimeData(animeName string, epNumber *int) string {
	searchAnimeResult, err := curd.SearchAnime(animeName, "sub")
	if err != nil {
		log.Error(err)
	}

	var AllanimeId string

	if epNumber != nil {
		AllanimeId, err = curd.FindKeyByValue(searchAnimeResult, fmt.Sprintf("%v (%d episodes)", animeName, *epNumber))
		if err != nil {
			log.Error("Failed to find anime in animeList:", err)
		}

	}

	// If unable to get Allanime id automatically get manually
	if AllanimeId == "" {
		log.Error("Failed to link anime automatically")
		for key, value := range searchAnimeResult {
			keyValueArray = append(keyValueArray, AllAnimeIdData{Id: key, Name: value})
		}

		fmt.Println(keyValueArray)
	}
	fmt.Println(AllanimeId)
	return AllanimeId
}

/*func main() {
	startCurdInteg()
}*/
