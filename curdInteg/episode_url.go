package curdInteg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"
)

type allanimeResponse struct {
	Data struct {
		Episode struct {
			SourceUrls []struct {
				SourceUrl string `json:"sourceUrl"`
			} `json:"sourceUrls"`
		} `json:"episode"`
	} `json:"data"`
}

func decodeProviderID(encoded string) string {
	// Split the string into pairs of characters (.. equivalent of 'sed s/../&\n/g')
	re := regexp.MustCompile("..")
	pairs := re.FindAllString(encoded, -1)

	// Mapping for the replacements
	replacements := map[string]string{
		"01": "9", "08": "0", "05": "=", "0a": "2", "0b": "3", "0c": "4", "07": "?",
		"00": "8", "5c": "d", "0f": "7", "5e": "f", "17": "/", "54": "l", "09": "1",
		"48": "p", "4f": "w", "0e": "6", "5b": "c", "5d": "e", "0d": "5", "53": "k",
		"1e": "&", "5a": "b", "59": "a", "4a": "r", "4c": "t", "4e": "v", "57": "o",
		"51": "i",
	}

	// Perform the replacement equivalent to sed 's/^../.../'
	for i, pair := range pairs {
		if val, exists := replacements[pair]; exists {
			pairs[i] = val
		}
	}

	// Join the modified pairs back into a single string
	result := strings.Join(pairs, "")

	// Replace "/clock" with "/clock.json" equivalent of sed "s/\/clock/\/clock\.json/"
	result = strings.ReplaceAll(result, "/clock", "/clock.json")

	// Print the final result
	return result
}

func extractLinks(provider_id string) map[string]interface{} {
	allanime_base := "https://allanime.day"
	url := allanime_base + provider_id
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	var videoData map[string]interface{}
	if err != nil {
		Log(fmt.Sprint("Error creating request:", err), logFile)
		return videoData
	}

	// Add the headers
	req.Header.Set("Referer", "https://allanime.to")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/121.0")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		Log(fmt.Sprint("Error sending request:", err), logFile)
		return videoData
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(fmt.Sprint("Error reading response:", err))
		return videoData
	}

	// Parse the JSON response
	err = json.Unmarshal(body, &videoData)
	if err != nil {
		log.Error(fmt.Sprint("Error parsing JSON:", err))
		return videoData
	}

	// Process the data as needed
	return videoData
}

// Get anime episode url respective to given config
// If the link is found, it returns a list of links. Otherwise, it returns an error.
//
// Parameters:
// - config: Configuration of the anime search.
// - id: Allanime id of the anime to search for.
// - epNo: Anime episode number to get links for.
//
// Returns:
// - []string: a list of links for specified episode.
// - error: an error if the episode is not found or if there is an issue during the search.
func GetEpisodeURL(config CurdConfig, id string, epNo int) ([]string, error) {
	query := `query($showId:String!,$translationType:VaildTranslationTypeEnumType!,$episodeString:String!){episode(showId:$showId,translationType:$translationType,episodeString:$episodeString){episodeString sourceUrls}}`

	variables := map[string]string{
		"showId":          id,
		"translationType": config.SubOrDub,
		"episodeString":   fmt.Sprintf("%d", epNo),
	}

	variablesJSON, err := json.Marshal(variables)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Set("query", query)
	values.Set("variables", string(variablesJSON))

	reqURL := fmt.Sprintf("%s/api?%s", "https://api.allanime.day", values.Encode())

	client := &http.Client{}
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/121.0")
	req.Header.Set("Referer", "https://allanime.to")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response allanimeResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		mu      sync.Mutex
		found   bool
		allinks []string
	)

	var wg sync.WaitGroup
	resultChan := make(chan []string, 1) // Buffer for immediate return

	for _, url := range response.Data.Episode.SourceUrls {
		if len(url.SourceUrl) > 2 && unicode.IsDigit(rune(url.SourceUrl[2])) {
			decodedProviderID := decodeProviderID(url.SourceUrl[2:])

			wg.Add(1)
			go func(id string) {
				defer wg.Done()

				// Early exit if context cancelled
				select {
				case <-ctx.Done():
					return
				default:
				}

				extractedLinks := extractLinks(id)

				// Check links for sharepoint
				if links, ok := extractedLinks["links"].([]interface{}); ok {
					for _, linkInterface := range links {
						if linkMap, ok := linkInterface.(map[string]interface{}); ok {
							if link, ok := linkMap["link"].(string); ok {
								fmt.Println(link)
								if strings.Contains(link, LinkPriorities[0]) {
									fmt.Println("found")
									mu.Lock()
									if !found {
										found = true
										resultChan <- []string{link} // Immediate result
										cancel()                     // Cancel all other operations
									}
									mu.Unlock()
									return
								}

								// Add normal link if no sharepoint found
								mu.Lock()
								allinks = append(allinks, link)
								mu.Unlock()
							}
						}
					}
				}
			}(decodedProviderID)
		}
	}

	// Wait for either the first sharepoint link or all goroutines
	go func() {
		wg.Wait()
		mu.Lock()
		defer mu.Unlock()
		if !found {
			resultChan <- allinks // Send final results
		}
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case <-time.After(10 * time.Second): // Safety timeout
		return nil, fmt.Errorf("timeout waiting for results")
	}
}
