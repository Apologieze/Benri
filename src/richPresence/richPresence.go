package richPresence

import (
	"AnimeGUI/src/config"
	"github.com/charmbracelet/log"
	"github.com/hugolgst/rich-go/client"
	"time"
)

var timeNow = time.Now()

func ResetTime() {
	timeNow = time.Now()
}
func InitDiscordRichPresence() {
	if !config.Setting.DiscordPresence {
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
	err := client.SetActivity(client.Activity{
		Details:    "In Main Menu",
		State:      "Free, open-source, crossplatform Anime streaming app :D",
		LargeImage: "https://apologize.fr/benri/icon.jpg",
		LargeText:  "Free, open-source, crossplatform Anime streaming app :D",
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

func SetActivity() {
	err := client.SetActivity(client.Activity{
		State:      "Episode 12/24 - 12min 3s",
		Details:    "Blue Box",
		LargeImage: "https://s4.anilist.co/file/anilistcdn/media/anime/cover/large/bx170942-B77wUSM1jQTu.jpg",
		LargeText:  "This is the large image :D",
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
