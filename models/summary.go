package models

type ChampionStats map[string]Result

type OpponentStats map[string]ChampionStats

type Result struct {
	Wins  int
	Games int
}
