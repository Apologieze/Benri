package main

import (
	curd "animeFyne/curdInteg"
	"fmt"
	"github.com/charmbracelet/log"
	"os"
	"path/filepath"
	"runtime"
)

var localAnime []curd.Anime

func startCurdInteg() {
	//discordClientId := "1287457464148820089"

	//var anime curd.Anime
	var user curd.User

	var homeDir string
	if runtime.GOOS == "windows" {
		homeDir = os.Getenv("USERPROFILE")
	} else {
		homeDir = os.Getenv("HOME")
	}

	configFilePath := filepath.Join(homeDir, ".config", "curd", "curd.conf")

	// load curd userCurdConfig
	userCurdConfig, err := curd.LoadConfig(configFilePath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	curd.SetGlobalConfig(&userCurdConfig)

	logFile := filepath.Join(os.ExpandEnv(userCurdConfig.StoragePath), "debug.log")
	curd.ClearLogFile(logFile)

	// Get the token from the token file
	user.Token, err = curd.GetTokenFromFile(filepath.Join(os.ExpandEnv(userCurdConfig.StoragePath), "token"))
	if err != nil {
		curd.Log("Error reading token", logFile)
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

	databaseFile := filepath.Join(os.ExpandEnv(userCurdConfig.StoragePath), "curd_history.txt")
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

func SearchFromAniId(id int) *curd.Anime {
	for _, anime := range localAnime {
		if anime.AnilistId == id {
			return &anime
		}
	}
	return nil
}

type KeyValue struct {
	Name  string
	Value string
}

var keyValueArray []KeyValue

func OnPlayButtonClick(id int) {
	if id == -1 {
		return
	}
	animePointer := SearchFromAniId(id)
	if animePointer == nil {
		return
	}
	fmt.Println(*animePointer)
	searchAnimeResult, err := curd.SearchAnime("The Eminence in Shadow", "sub")
	if err != nil {
		log.Error(err)
	}

	for key, value := range searchAnimeResult {
		keyValueArray = append(keyValueArray, KeyValue{Name: key, Value: value})
	}

	fmt.Println(keyValueArray)

}

/*func main() {
	startCurdInteg()
}*/
