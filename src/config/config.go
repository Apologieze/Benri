package config

import "fyne.io/fyne/v2"

type Config struct {
	preference      fyne.Preferences
	SkipOpening     bool `json:"skip_opening"`
	SkipEnding      bool `json:"skip_ending"`
	TrayIcon        bool `json:"tray_icon"`
	DiscordPresence bool
}

const (
	SkipOpeningKey  = "skip_opening"
	SkipEndingKey   = "skip_ending"
	TrayIconKey     = "tray_icon"
	DiscordPresence = "discord_presence"
)

var Setting Config

// NewConfig returns a new Config struct
func CreateConfig(app fyne.Preferences) {
	Setting = Config{
		preference:      app,
		SkipOpening:     app.Bool(SkipOpeningKey),
		SkipEnding:      app.Bool(SkipEndingKey),
		TrayIcon:        app.Bool(TrayIconKey),
		DiscordPresence: app.BoolWithFallback(DiscordPresence, true),
	}
}

func SetBool(key string, value bool) {
	switch key {
	case SkipOpeningKey:
		Setting.SkipOpening = value
	case SkipEndingKey:
		Setting.SkipEnding = value
	case TrayIconKey:
		Setting.TrayIcon = value
	case DiscordPresence:
		Setting.DiscordPresence = value
	}
	Setting.preference.SetBool(key, value)
}
