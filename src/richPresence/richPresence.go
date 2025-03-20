package richPresence

import (
	"AnimeGUI/richgo/client"
	"AnimeGUI/src/config"
	"fmt"
	"github.com/charmbracelet/log"
	"time"
)

type PresenceAnime struct {
	Name         string
	Ep           int
	TotalEp      int
	ImageLink    string
	PlaybackTime int
	Duration     int
	Paused       bool
}

type PresenceState int

const (
	MainMenu PresenceState = iota
	Watching
)

const advert = "Free, open-source, crossplatform Anime streaming app :D"

var timeStartApp = time.Now()

var timeBegin = timeStartApp
var timeEnd = timeBegin.Add(time.Minute * 30)

func ResetTime(PlaybackTime, Duration int) {
	timeBegin = time.Now().Add(-time.Second * time.Duration(PlaybackTime))
	timeEnd = timeBegin.Add(time.Second * time.Duration(Duration))
}
func InitDiscordRichPresence() {
	if !config.Setting.DiscordPresence {
		client.Logout()
		return
	}
	err := client.Login("1046397185467621418")
	if err != nil {
		log.Error(err)
		return
	}

	SetMenuActivity()
}

func SetMenuActivity() {
	log.Info("Main Menu Activity presence")
	err := client.SetActivity(client.Activity{
		Type:       client.ActivityTypeWatching,
		Details:    "In Main Menu",
		State:      advert,
		LargeImage: "main-image",
		LargeText:  advert,
		/*SmallImage: "https://s4.anilist.co/file/anilistcdn/media/anime/cover/large/bx170942-B77wUSM1jQTu.jpg",
		SmallText:  "And this is the small image",*/
		/*Party: &client.Party{
			ID:         "-1",
			Players:    15,
			MaxPlayers: 24,
		},*/
		Buttons: []*client.Button{
			&client.Button{
				Label: "Try Benri", // Button label
				Url:   "https://uwu.apologize.fr/preview",
			},
			/*&client.Button{
				Label: "Try 2Benri", // Button label
				Url:   "https://uwu.apologize.fr/preview",
			},*/
		},
		Timestamps: &client.Timestamps{
			Start: &timeStartApp,
			//End:   &then,
		},
	})

	if err != nil {
		log.Error(err)
	}
}

func SetAnimeActivity(anime *PresenceAnime) {
	//zeroTime := time.Now()
	if anime == nil {
		SetMenuActivity()
		return
	}
	var timeStamps *client.Timestamps = nil

	if !anime.Paused {
		timeStamps = &client.Timestamps{
			Start: &timeBegin,
			End:   &timeEnd,
		}
	}

	err := client.SetActivity(client.Activity{
		Type:       client.ActivityTypeWatching,
		State:      fmt.Sprintf("Episode %d/%d %s", anime.Ep, anime.TotalEp, pausedToString(anime.Paused)),
		Details:    fmt.Sprintf("%s", anime.Name),
		LargeImage: anime.ImageLink,
		LargeText:  anime.Name,
		SmallImage: "main-image",
		SmallText:  advert,
		/*Party: &client.Party{
			ID:         "-1",
			Players:    15,
			MaxPlayers: 24,
		},*/
		Buttons: []*client.Button{
			&client.Button{
				Label: "Try Benri", // Button label
				Url:   "https://uwu.apologize.fr/preview",
			},
			/*&client.Button{
				Label: "Try 2Benri", // Button label
				Url:   "https://uwu.apologize.fr/preview",
			},*/
		},
		Timestamps: timeStamps,
	})

	if err != nil {
		log.Error(err)
	}
}

func pausedToString(paused bool) string {
	if paused {
		return "PAUSED"
	}
	return ""
}

func numberToTime(seconds int) string {
	minutes := seconds / 60
	seconds %= 60
	return fmt.Sprintf("%dmin %ds", minutes, seconds)
}
