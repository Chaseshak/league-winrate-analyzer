package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/Chaseshak/league-winrate-analyzer/analyze"
	"github.com/Chaseshak/league-winrate-analyzer/api"
	"github.com/Chaseshak/league-winrate-analyzer/models"
)

func main() {
	loadEnv()

	summonerName := flag.String("summoner", "", "Summoner name to fetch matches for")
	tagline := flag.String("tagline", "", "Riot ID tagline (e.g. NA1)")
	flag.Parse()

	if *summonerName == "" {
		fmt.Println("Please provide a summoner name with the -summoner flag")
		os.Exit(1)
	}
	if *tagline == "" {
		fmt.Println("Please provide a tagline with the -tagline flag")
		os.Exit(1)
	}

	summoner := models.Summoner{
		GameName: *summonerName,
		TagLine:  *tagline,
	}

	matches := api.FetchMatches(&summoner)
	results := analyze.Summarize(summoner, matches)

	for champion, championStats := range results {
		fmt.Println("---When playing against---", champion)

		for myChampion, result := range championStats {
			winrate := float64(result.Wins) / float64(result.Games) * 100
			winrate = float64(int(winrate*100)) / 100
			fmt.Println("\t", myChampion, ":", winrate, "%", "(", result.Wins, "/", result.Games, ")")
		}
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		fmt.Println("Please create a .env following the .env.example file")
		os.Exit(1)
	}
}
