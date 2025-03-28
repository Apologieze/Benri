package curdInteg

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
)

// skipTimesResponse struct to hold the response from the AniSkip API
type skipTimesResponse struct {
	Found   bool         `json:"found"`
	Results []skipResult `json:"results"`
}

// skipResult struct to hold individual skip result data
type skipResult struct {
	Interval skipInterval `json:"interval"`
}

// skipInterval struct to hold the start and end times for skip intervals
type skipInterval struct {
	StartTime float64 `json:"start_time"`
	EndTime   float64 `json:"end_time"`
}

// GetAniSkipData fetches skip times data for a given anime ID and episode
func GetAniSkipData(animeMalId int, episode int) (string, error) {
	baseURL := "https://api.aniskip.com/v1/skip-times"
	url := fmt.Sprintf("%s/%d/%d?types=op&types=ed", baseURL, animeMalId, episode)

	resp, err := http.Get(url)
	if err != nil {
		Log(fmt.Sprintf("error fetching data from AniSkip API: %w", err), logFile)
		return "", fmt.Errorf("error fetching data from AniSkip API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		Log(fmt.Sprintf("failed with status %d", resp.StatusCode), logFile)
		return "", fmt.Errorf("failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		Log(fmt.Sprintf("failed to read response body: %w", err), logFile)
		return "", fmt.Errorf("failed to read response body %w", err)
	}

	return string(body), nil
}

// RoundTime rounds a time value to the specified precision
func RoundTime(timeValue float64, precision int) float64 {
	multiplier := math.Pow(10, float64(precision))
	return math.Floor(timeValue*multiplier+0.5) / multiplier
}

// ParseAniSkipResponse parses the response text from the AniSkip API and updates the Anime struct
func ParseAniSkipResponse(responseText string, anime *Anime, timePrecision int) error {
	if responseText == "" {
		return fmt.Errorf("response text is empty")
	}

	var data skipTimesResponse
	err := json.Unmarshal([]byte(responseText), &data)
	if err != nil {
		return fmt.Errorf("error unmarshalling response: %w", err)
	}

	if !data.Found {
		return fmt.Errorf("no skip times found")
	}

	// Populate skip times for the anime's episode
	if len(data.Results) > 0 {
		op := data.Results[0].Interval
		anime.Ep.SkipTimes.Op = Skip{
			Start: int(RoundTime(op.StartTime, timePrecision)),
			End:   int(RoundTime(op.EndTime, timePrecision)),
		}
	}

	if len(data.Results) > 1 {
		ed := data.Results[len(data.Results)-1].Interval
		anime.Ep.SkipTimes.Ed = Skip{
			Start: int(RoundTime(ed.StartTime, timePrecision)),
			End:   int(RoundTime(ed.EndTime, timePrecision)),
		}
	}

	return nil
}

// GetAndParseAniSkipData fetches and parses skip times for a given anime ID and episode
func GetAndParseAniSkipData(animeMalId int, episode int, timePrecision int, anime *Anime) error {
	responseText, err := GetAniSkipData(animeMalId, episode)
	if err != nil {
		return err
	}
	return ParseAniSkipResponse(responseText, anime, timePrecision)
}
