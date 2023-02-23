package application

import (
	"fmt"
	"team-cli/config"
	"team-cli/repository"
	"team-cli/team"
)

type PlayerFetcher struct {
	maxRoutines int
}

func NewPlayerFetcher(maxRoutines int) *PlayerFetcher {
	return &PlayerFetcher{maxRoutines: maxRoutines}
}

func (f *PlayerFetcher) GetTeamPlayers(targetTeams map[string]bool) (repository.PlayerRepository, error) {
	results := make(chan team.Result)
	done := make(chan bool)

	cfg, err := config.Get()
	if err != nil {
		return repository.PlayerRepository{}, fmt.Errorf("getting config: %v", err)
	}

	repo := repository.PlayerRepository{}
	go repo.SavePlayers(targetTeams, results, done)

	team.GetPlayers(cfg, f.maxRoutines, results, done)

	return repo, nil
}
