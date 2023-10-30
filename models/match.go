package models

type TeamPosition string

const (
	Top     TeamPosition = "TOP"
	Jungle  TeamPosition = "JUNGLE"
	Mid     TeamPosition = "MID"
	Bot     TeamPosition = "BOT"
	Support TeamPosition = "SUPPORT"
)

type Match struct {
	Info MatchInfo `json:"info"`
}

type MatchInfo struct {
	Participants []Participant `json:"participants"`
}

type Participant struct {
	ChampionName string       `json:"championName"`
	Puuid        string       `json:"puuid"`
	TeamPosition TeamPosition `json:"teamPosition"`
	Win          bool         `json:"win"`
}
