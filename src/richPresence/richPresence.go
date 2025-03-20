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
}

type PresenceState int

const (
	MainMenu PresenceState = iota
	Watching
)

const advert = "Free, open-source, crossplatform Anime streaming app :D"

var timeNow = time.Now()

func ResetTime(waitingTime int) {
	go func() {
		time.Sleep(time.Duration(waitingTime) * time.Second)
		timeNow = time.Now()
	}()
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
			Start: &timeNow,
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

	err := client.SetActivity(client.Activity{
		Type:       client.ActivityTypeWatching,
		State:      fmt.Sprintf("%s remaining", numberToTime(anime.Duration-anime.PlaybackTime)),
		Details:    fmt.Sprintf("%s Episode %d/%d", anime.Name, anime.Ep, anime.TotalEp),
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
		Timestamps: &client.Timestamps{
			Start: &timeNow,
			//End:   &then,
		},
	})

	if err != nil {
		log.Error(err)
	}
}

func numberToTime(seconds int) string {
	minutes := seconds / 60
	seconds %= 60
	return fmt.Sprintf("%dmin %ds", minutes, seconds)
}
