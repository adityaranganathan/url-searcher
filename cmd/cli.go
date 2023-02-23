package cmd

import (
	"fmt"
	"strings"
	"team-cli/application"
	"team-cli/config"
	"team-cli/repository"
)

func RunCLI(fetcher *application.PlayerFetcher) error {
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("getting config: %v", err)
	}

	playerInfo, err := fetcher.GetTeamPlayers(cfg.GetTargetTeams())
	if err != nil {
		return fmt.Errorf("getting player info: %v", err)
	}

	DisplayPlayers(playerInfo)

	return nil
}

// DisplayPlayers renders to stdout the information about players ordered by their ID.
func DisplayPlayers(repo repository.PlayerRepository) {
	for _, id := range repo.GetSortedPlayerIDs() {
		player := repo[id]
		fmt.Println(fmt.Sprintf("%s; %s; %s; %s", player.ID, player.Name, player.Age, strings.Join(player.Teams, ",")))
	}
}
