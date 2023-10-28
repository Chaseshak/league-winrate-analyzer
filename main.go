package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/Chaseshak/league-winrate-analyzer/api"
	"github.com/Chaseshak/league-winrate-analyzer/models"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

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

	api.FetchMatches(&summoner)
}

