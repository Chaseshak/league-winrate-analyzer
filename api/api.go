package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Chaseshak/league-winrate-analyzer/models"
)

const (
	host            = "api.riotgames.com"
	region          = "americas"
	puuid_path      = "/riot/account/v1/accounts/by-riot-id"
	match_ids_path  = "/lol/match/v5/matches/by-puuid"
	match_path      = "/lol/match/v5/matches"
	ranked_queue_id = "420"
	rate_limit      = 15
)

func FetchMatches(summoner *models.Summoner) []models.Match {
	fetchPUUID(summoner)
	matchIDs := fetchMatchIDs(summoner)
	fmt.Println("Found", len(matchIDs), "matches")

	sem := make(chan struct{}, 20)

	var matches []models.Match
	var mu sync.Mutex

	for _, matchID := range matchIDs {
		sem <- struct{}{}

		go func(id string) {
			defer func() { <-sem }()
			url := buildUrl(fmt.Sprintf("%s/%s", match_path, id))
			body := apiGet(url, make(map[string]string))

			var match models.Match
			err := json.Unmarshal(body, &match)
			if err != nil {
				log.Println("Error parsing ", id, ": ", err)
			}

			mu.Lock()
			matches = append(matches, match)
			mu.Unlock()
		}(matchID)

		time.Sleep(time.Second / rate_limit)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}

	return []models.Match{}
}

func fetchMatchIDs(summoner *models.Summoner) []string {
	url := buildUrl(fmt.Sprintf("%s/%s/ids", match_ids_path, summoner.Puuid))
	perPage := 100
	start, _ := time.Parse("2006-01-02 15:04:05", "2023-08-10 00:00:00")
	queryParams := map[string]string{
		"startTime": strconv.FormatInt(start.Unix(), 10),
		"count":     strconv.Itoa(perPage),
		"queue":     ranked_queue_id,
	}
	var matchIDs []string
	page := 0

	for {
		queryParams["start"] = fmt.Sprintf("%d", page*perPage)
		body := apiGet(url, queryParams)

		var newMatchIDs []string
		err := json.Unmarshal(body, &newMatchIDs)
		if err != nil {
			log.Fatal("Error parsing Match IDs:", err)
		}

		matchIDs = append(matchIDs, newMatchIDs...)

		if len(newMatchIDs) < perPage {
			break
		}

		page++
	}

	return matchIDs
}

func fetchPUUID(summoner *models.Summoner) {
	url := buildUrl(fmt.Sprintf("%s/%s/%s", puuid_path, summoner.GameName, summoner.TagLine))

	body := apiGet(url, make(map[string]string))

	err := json.Unmarshal(body, &summoner)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
	}
}

func apiGet(url url.URL, params map[string]string) []byte {
	q := url.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	url.RawQuery = q.Encode()

	log.Println("GET: ", url.String())

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url.String(), nil)
	req.Header.Add("X-Riot-Token", apiKey())

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body: ", err)
	}
	if resp.StatusCode != 200 {
		log.Printf("Error %d from %s: %s", resp.StatusCode, url.String(), responseBody)
		log.Fatal("Exiting...")
	}

	return responseBody
}

func buildUrl(path string) url.URL {
	return url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s.%s", region, host),
		Path:   path,
	}
}

func apiKey() string {
	return os.Getenv("RIOT_API_KEY")
}
