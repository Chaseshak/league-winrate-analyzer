package analyze

import (
	"errors"

	"github.com/Chaseshak/league-winrate-analyzer/models"
)

func Summarize(mySummoner models.Summoner, matches []models.Match) map[string]models.ChampionStats {
	results := make(map[string]models.ChampionStats)

	for _, match := range matches {
		myParticipant, err := findMyParticipant(mySummoner, match.Info.Participants)
		if err != nil {
			continue
		}

		opponentParticipant, err := findOpponentParticipant(myParticipant, match.Info.Participants)
		if err != nil {
			continue
		}

		if results[opponentParticipant.ChampionName] == nil {
			results[opponentParticipant.ChampionName] = make(models.ChampionStats)
		}

		tmpResults := results[opponentParticipant.ChampionName][myParticipant.ChampionName]
		tmpResults.Games++

		if myParticipant.Win {
			tmpResults.Wins++
		}

		results[opponentParticipant.ChampionName][myParticipant.ChampionName] = tmpResults
	}

	return results
}

func findMyParticipant(mySummoner models.Summoner, participants []models.Participant) (models.Participant, error) {
	for _, participant := range participants {
		if participant.Puuid == mySummoner.Puuid {
			return participant, nil
		}
	}

	return models.Participant{}, errors.New("Participant not found for summoner")
}

func findOpponentParticipant(myParticipant models.Participant, participants []models.Participant) (models.Participant, error) {
	for _, participant := range participants {
		if participant.Puuid != myParticipant.Puuid && participant.TeamPosition == myParticipant.TeamPosition {
			return participant, nil
		}
	}

	return models.Participant{}, errors.New("Opponent participant not found")
}
