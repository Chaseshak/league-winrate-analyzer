package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/Chaseshak/league-winrate-analyzer/models"
)

type Match struct {
}

const (
	host       = "api.riotgames.com"
	region     = "americas"
	puuid_path = "/riot/account/v1/accounts/by-riot-id/"
)

func FetchMatches(summoner *models.Summoner) []Match {
	fetchPUUID(summoner)
	fmt.Printf("Summoner: %+v\n", summoner)

	return []Match{}
}

func fetchPUUID(summoner *models.Summoner) {
	url := buildUrl(fmt.Sprintf("%s%s/%s", puuid_path, summoner.GameName, summoner.TagLine))

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-Riot-Token", apiKey())

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}
	if resp.StatusCode != 200 {
		fmt.Println("Error fetching summoner:", string(body))
	}

	err = json.Unmarshal(body, &summoner)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
	}
}

func buildUrl(path string) string {
	u := &url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s.%s", region, host),
		Path:   path,
	}

	return u.String()
}

func apiKey() string {
	return os.Getenv("RIOT_API_KEY")
}
