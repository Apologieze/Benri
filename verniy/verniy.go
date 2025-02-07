package verniy

import (
	"net/http"
	"time"

	"AnimeGUI/verniy/limiter"
)

// Client is anilist client.
type Client struct {
	Host        string
	Http        http.Client
	Limiter     limiter.Limiter
	AccessToken string
}

// New to create new anilist client.
func New() *Client {
	return &Client{
		Host: "https://graphql.anilist.co",
		Http: http.Client{
			Timeout: 10 * time.Second,
		},
		Limiter: limiter.New(90, time.Minute),
	}
}
