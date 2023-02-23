package repository_test

import (
	. "github.com/onsi/gomega"
	"team-cli/repository"
	"team-cli/team"
	"testing"
)

func TestSavePlayers(t *testing.T) {
	g := NewGomegaWithT(t)

	targetTeams := map[string]bool{
		"Team1": true,
		"Team3": true,
	}

	results := make(chan team.Result)
	done := make(chan bool)

	repo := repository.PlayerRepository{}
	go repo.SavePlayers(targetTeams, results, done)

	team1 := team.Team{
		Name: "Team1",
		Players: []team.Player{
			{ID: "1", Name: "Name1", Age: "20"},
			{ID: "2", Name: "Name2", Age: "30"},
		},
	}
	team2 := team.Team{
		Name: "Team2",
		Players: []team.Player{
			{ID: "3", Name: "Name3", Age: "40"},
			{ID: "4", Name: "Name4", Age: "50"},
		},
	}
	team3 := team.Team{
		Name: "Team3",
		Players: []team.Player{
			{ID: "1", Name: "Name1", Age: "20"},
			{ID: "5", Name: "Name5", Age: "50"},
		},
	}

	results <- team.Result{Team: team1}
	results <- team.Result{Team: team2}
	results <- team.Result{Team: team3}

	<-done
	close(results)

	expectedPlayers := []repository.PlayerInfo{
		{ID: "1", Name: "Name1", Age: "20", Teams: []string{"Team1", "Team3"}},
		{ID: "2", Name: "Name2", Age: "30", Teams: []string{"Team1"}},
		{ID: "5", Name: "Name5", Age: "50", Teams: []string{"Team3"}},
	}

	g.Expect(len(repo)).To(Equal(3))
	g.Expect(repo).To(HaveKeyWithValue(1, expectedPlayers[0]))
	g.Expect(repo).To(HaveKeyWithValue(2, expectedPlayers[1]))
	g.Expect(repo).To(HaveKeyWithValue(5, expectedPlayers[2]))
}

func TestSave(t *testing.T) {
	g := NewGomegaWithT(t)

	repo := repository.PlayerRepository{}
	player := team.Player{ID: "1", Name: "Name1", Age: "25"}

	err := repo.Save(player, "Team1")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(repo)).To(Equal(1))
	g.Expect(repo[1]).To(Equal(repository.PlayerInfo{
		ID:    player.ID,
		Name:  player.Name,
		Age:   player.Age,
		Teams: []string{"Team1"},
	}))
}

func TestGetSortedPlayerIDs(t *testing.T) {
	g := NewGomegaWithT(t)

	players := []team.Player{{ID: "3"}, {ID: "2"}, {ID: "4"}}

	repo := repository.PlayerRepository{}
	for _, player := range players {
		err := repo.Save(player, "TeamName")
		g.Expect(err).NotTo(HaveOccurred())
	}

	g.Expect(repo.GetSortedPlayerIDs()).To(Equal([]int{2, 3, 4}))
}
