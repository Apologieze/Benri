package anilist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func makePostRequest(url, query string, variables map[string]interface{}, headers map[string]string) (map[string]interface{}, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json") // <-- Important!
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status %d: %s", resp.StatusCode, body)
	}

	var responseData map[string]interface{}
	// Unmarshal the response into a map
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return responseData, nil
}

func UpdateAnimeStatus(token string, mediaID int, status string) error {
	url := "https://graphql.anilist.co"
	query := `
	mutation($mediaId: Int, $status: MediaListStatus) {
		SaveMediaListEntry(mediaId: $mediaId, status: $status) {
			id
			status
		}
	}`

	variables := map[string]interface{}{
		"mediaId": mediaID,
		"status":  status,
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	_, err := makePostRequest(url, query, variables, headers)
	if err != nil {
		return fmt.Errorf("failed to update anime status: %w", err)
	}
	return nil
}
