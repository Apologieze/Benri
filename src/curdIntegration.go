package main

import (
	curd "AnimeGUI/curdInteg"
	"AnimeGUI/src/anilist"
	"AnimeGUI/src/config"
	"AnimeGUI/src/richPresence"
	"AnimeGUI/verniy"
	"fmt"
	"github.com/charmbracelet/log"
	"os"
	"path/filepath"
	"runtime"
	"time"
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

	//var logFile = filepath.Join(os.ExpandEnv(userCurdConfig.StoragePath), "debug.log")
	//curd.ClearLogFile(logFile)

	// Get the token from the token file
	user.Token, err = curd.GetTokenFromFile(filepath.Join(os.ExpandEnv(userCurdConfig.StoragePath), "token"))
	if err != nil {
		log.Error("Error reading token")
	}
	if user.Token == "" {
		setTokenGraphicaly(filepath.Join(os.ExpandEnv(userCurdConfig.StoragePath), "token"), &user)
	}
}

func secondCurdInit() {
	if user.Token == "" {
		curd.ChangeToken(&userCurdConfig, &user)
	}

	var err error
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

func OnPlayButtonClick(animeName string, animeData *verniy.MediaList, savingWatch bool) {
	if mpvPresent == false {
		log.Error("MPV is not yet dl")
		return
	}
	if animeData == nil {
		log.Error("Anime data is nil")
		return
	}
	var allAnimeId string
	animeProgress := 0
	if animeData.Progress != nil && animeData.Media.Episodes != nil {
		animeProgress = min(*animeData.Progress, *animeData.Media.Episodes-1)
	}
	animePointer := SearchFromLocalAniId(animeData.Media.ID)
	if animePointer == nil {
		allAnimeId = searchAllAnimeData(anilist.AnimeToRomaji(animeData.Media), animeData.Media.Episodes, animeProgress)
		if allAnimeId == "" {
			log.Error("Failed to get allAnimeId")
			return
		}
		err, _ := curd.LocalUpdateAnime(databaseFile, animeData.Media.ID, allAnimeId, animeProgress, 0, 0, animeName)
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

	animeProgress++
	if !savingWatch {
		animeProgress = max(animeProgress-1, 1)
	}

	log.Info("Anime Progress:", animeProgress)

	fmt.Println("Start getting url")
	url, err := curd.GetEpisodeURL(userCurdConfig, allAnimeId, animeProgress)
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println("--Got all urls")
	finalLink := curd.PrioritizeLink(url)
	if len(finalLink) < 5 {
		log.Error("No valid link found")
		return
	}
	fmt.Println("Final Link:", finalLink)

	mpvSocketPath, err := curd.StartVideo(finalLink, []string{}, fmt.Sprintf("%s - Episode %d", animeName, animeProgress))
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println("MPV Socket Path:", mpvSocketPath)
	playingAnime := curd.Anime{AnilistId: animeData.Media.ID, AllanimeId: allAnimeId}
	playingAnime.Ep.Player.SocketPath = mpvSocketPath
	playingAnime.Title.English = animeName
	playingAnime.Ep.Number = animeProgress - 1
	playingAnime.TotalEpisodes = *animeData.Media.Episodes
	if animePointer != nil {
		fmt.Println("AnimePointer:", animePointer.Ep.Number, playingAnime.Ep.Number)
		if animePointer.Ep.Number == playingAnime.Ep.Number {
			playingAnime.Ep.Player.PlaybackTime = animePointer.Ep.Player.PlaybackTime
		}
	}
	playingAnimeLoop(playingAnime, animeData, savingWatch)
}

func searchAllAnimeData(animeName string, epNumber *int, animeProgress int) string {
	fmt.Println(animeName)
	searchAnimeResult, err := curd.SearchAnime(animeName, "sub")
	fmt.Println(searchAnimeResult)
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
		var keyValueArray []AllAnimeIdData
		log.Error("Failed to link anime automatically")
		for key, value := range searchAnimeResult {
			keyValueArray = append(keyValueArray, AllAnimeIdData{Id: key, Name: value})
		}

		selectCorrectLinking(keyValueArray, animeName, animeProgress)
		return ""
	}
	fmt.Println(AllanimeId)
	return AllanimeId
}

func playingAnimeLoop(playingAnime curd.Anime, animeData *verniy.MediaList, savingWatch bool) {
	fmt.Println(playingAnime.Ep.Player.PlaybackTime, "ah oue")
	// Get video duration
	go func() {
		defer richPresence.SetMenuActivity()
		for {
			time.Sleep(1 * time.Second)
			if playingAnime.Ep.Duration == 0 {
				// Get video duration
				durationPos, err := curd.MPVSendCommand(playingAnime.Ep.Player.SocketPath, []interface{}{"get_property", "duration"})
				if err != nil {
					log.Error("Error getting video duration: " + err.Error())
				} else if durationPos != nil {
					if duration, ok := durationPos.(float64); ok {
						playingAnime.Ep.Duration = int(duration + 0.5) // Round to nearest integer
						log.Infof("Video duration: %d seconds", playingAnime.Ep.Duration)

						if !savingWatch {
							playingAnime.Ep.Player.PlaybackTime = int(float64(playingAnime.Ep.Duration) * 0.80)
						}

						if playingAnime.Ep.Player.PlaybackTime > 10 {
							_, err := curd.SeekMPV(playingAnime.Ep.Player.SocketPath, max(0, playingAnime.Ep.Player.PlaybackTime-5))
							if err != nil {
								log.Error("Error seeking video: " + err.Error())
							}
						} else {
							log.Error("Error seeking video: playback time is", playingAnime.Ep.Player.PlaybackTime)
						}
						//richPresence.ResetTime()
						break
					} else {
						log.Error("Error: duration is not a float64")
					}
				}
			}
		}

		presenceAnime := richPresence.PresenceAnime{Name: playingAnime.Title.English, Ep: playingAnime.Ep.Number + 1, ImageLink: *animeData.Media.CoverImage.Large, PlaybackTime: 0, Duration: playingAnime.Ep.Duration, TotalEp: playingAnime.TotalEpisodes}
		for {
			time.Sleep(1 * time.Second)
			timePos, err := curd.MPVSendCommand(playingAnime.Ep.Player.SocketPath, []interface{}{"get_property", "time-pos"})
			if err != nil {
				if savingWatch {
					log.Error("Error getting video position: " + err.Error())
					fmt.Println("EH en vrai", playingAnime.Ep.Player.PlaybackTime, playingAnime.Ep.Duration)
					percentageWatched := curd.PercentageWatched(playingAnime.Ep.Player.PlaybackTime, playingAnime.Ep.Duration)

					if int(percentageWatched) >= userCurdConfig.PercentageToMarkComplete {
						playingAnime.Ep.Number++
						playingAnime.Ep.Player.PlaybackTime = 0
						var newProgress int = playingAnime.Ep.Number
						animeData.Progress = &newProgress
						go UpdateProgressHandler(playingAnime)
						episodeNumber.SetText(fmt.Sprintf("Episode %d/%d", playingAnime.Ep.Number, playingAnime.TotalEpisodes))
					}

					err, tempAnime := curd.LocalUpdateAnime(databaseFile, playingAnime.AnilistId, playingAnime.AllanimeId, playingAnime.Ep.Number, playingAnime.Ep.Player.PlaybackTime, 0, playingAnime.Title.English)
					if err == nil && tempAnime != nil {
						log.Info("Successfully updated database file")
						localAnime = curd.LocalGetAllAnime(databaseFile)
					}
					displayLocalProgress()
				}
				break
			}
			if timePos != nil && playingAnime.Ep.Duration != 0 {
				if timing, ok := timePos.(float64); ok {
					playingAnime.Ep.Player.PlaybackTime = int(timing + 0.5)
					log.Infof("Video position: %d seconds", playingAnime.Ep.Player.PlaybackTime)
					if config.Setting.DiscordPresence {
						presenceAnime.PlaybackTime = playingAnime.Ep.Player.PlaybackTime
						richPresence.SetAnimeActivity(&presenceAnime)
					}
				} else {
					log.Error("Error: time-pos is not a float64")
				}

			}

		}
	}()
}

func UpdateProgressHandler(anime curd.Anime) {
	UpdateAnimeProgress(anime.AnilistId, anime.Ep.Number)
	if anime.TotalEpisodes == anime.Ep.Number {
		// Mark as complete
		go anilist.GetData(categoryRadiobox, user.Username, func() { log.Error("Invalid token") })
	}
}

func UpdateAnimeProgress(animeId int, episode int) {
	err := curd.UpdateAnimeProgress(user.Token, animeId, episode)
	if err != nil {
		log.Error(err)
	}
}

func deleteTokenFile() {
	err := os.Remove(filepath.Join(os.ExpandEnv(userCurdConfig.StoragePath), "token"))
	if err != nil {
		log.Error(err)
	}
}

func displayLocalProgress() {
	localDbAnime := SearchFromLocalAniId(animeSelected.Media.ID)
	AnimeProgress := IntPointerFallback(animeSelected.Progress, 0)
	AnimeEpisode := IntPointerFallback(animeSelected.Media.Episodes, 0)

	currentEp := min(AnimeProgress+1, AnimeEpisode)
	playButton.Text = fmt.Sprint("Play Ep", currentEp)
	fmt.Println("Current Ep:", currentEp)

	defer setPlayButtonVisibility()
	if localDbAnime != nil {
		if localDbAnime.Ep.Number == AnimeProgress {
			if localDbAnime.Ep.Player.PlaybackTime == 0 {
				episodeLastPlayback.SetText(fmt.Sprintf("Just finished Episode %d", localDbAnime.Ep.Number))
				episodeLastPlayback.Show()
				return
			} else {
				episodeLastPlayback.SetText(fmt.Sprintf("Last saved at EP%d: [%s]", localDbAnime.Ep.Number+1, time.Second*time.Duration(localDbAnime.Ep.Player.PlaybackTime)))
				playButton.Text = fmt.Sprint("Resuming Ep", currentEp)
				episodeLastPlayback.Show()
				return
			}
		}
	}
	episodeLastPlayback.Hide()
}

func setPlayButtonVisibility() {
	defer playButton.Refresh()
	if animeSelected.Media.NextAiringEpisode != nil {
		if *animeSelected.Progress+1 == animeSelected.Media.NextAiringEpisode.Episode {
			playButton.Hide()
			return
		}
	}
	if animeSelected.Media.Episodes == nil {
		playButton.Hide()
		return
	}
	playButton.Show()
}

func IntPointerFallback(ptr *int, value int) int {
	if ptr == nil {
		return value
	}
	return *ptr
}
