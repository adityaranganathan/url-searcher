package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"team-cli/repository"
)

//go:generate mockery --name=PlayerFetcher
type PlayerFetcher interface {
	GetTeamPlayers(map[string]bool) (repository.PlayerRepository, error)
}

func GetHandleTeamRequest(fetcher PlayerFetcher) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		rawTeams := request.URL.Query()

		targetTeams := make(map[string]bool)
		for _, t := range rawTeams["team"] {
			targetTeams[t] = true
		}

		playerInfo, err := fetcher.GetTeamPlayers(targetTeams)
		if err != nil {
			fmt.Println(fmt.Sprintf("encoding player response: %v", err))

			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		var response HandleTeamResponse
		sortedIDs := playerInfo.GetSortedPlayerIDs()
		for _, id := range sortedIDs {
			response = append(response, PlayerInfoResponse{
				ID:    playerInfo[id].ID,
				Name:  playerInfo[id].Name,
				Age:   playerInfo[id].Age,
				Teams: playerInfo[id].Teams,
			})
		}

		writer.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(writer).Encode(response)
		if err != nil {
			fmt.Println(fmt.Sprintf("encoding player response: %v", err))
			writer.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}
}

type HandleTeamResponse []PlayerInfoResponse

type PlayerInfoResponse struct {
	ID    string
	Name  string
	Age   string
	Teams []string
}

func GetRespondOKHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("OK"))
	}
}
