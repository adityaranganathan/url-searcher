package repository

import (
	"fmt"
	"sort"
	"strconv"
	"team-cli/team"
)

type PlayerInfo struct {
	ID    string
	Name  string
	Age   string
	Teams []string
}

// PlayerRepository stores a map of player IDs to PlayerInfo.
type PlayerRepository map[int]PlayerInfo

// SavePlayers saves player information from a team result channel into the repository.
// Only players from target teams specified in the configuration will be saved.
// Once all target teams have been found, it sends to the done channel.
func (r PlayerRepository) SavePlayers(targetTeams map[string]bool, results chan team.Result, done chan bool) {

	for result := range results {
		if result.Err != nil {
			fmt.Println(fmt.Sprintf("processing team: %v", result.Err))
			continue
		}

		if _, ok := targetTeams[result.Team.Name]; !ok {
			continue
		}
		delete(targetTeams, result.Team.Name)

		for _, player := range result.Team.Players {
			err := r.Save(player, result.Team.Name)
			if err != nil {
				fmt.Println(fmt.Sprintf("saving player: %v", err))
				continue
			}
		}

		if len(targetTeams) == 0 {
			done <- true
		}
	}
}

// Save adds a player to the repository.
func (r PlayerRepository) Save(player team.Player, team string) error {
	playerID, err := strconv.Atoi(player.ID)
	if err != nil {
		return fmt.Errorf("converting player ID: %v", err)
	}

	var teams []string
	if info, ok := r[playerID]; ok {
		teams = info.Teams
	}
	teams = append(teams, team)

	r[playerID] = PlayerInfo{
		ID:    player.ID,
		Name:  player.Name,
		Age:   player.Age,
		Teams: teams,
	}

	return nil
}

// GetSortedPlayerIDs returns a slice of player IDs sorted in ascending order.
func (r PlayerRepository) GetSortedPlayerIDs() []int {
	var ids []int
	for id := range r {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	return ids
}
