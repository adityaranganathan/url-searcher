package team_test

import (
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"team-cli/config"
	"team-cli/team"
	"testing"
)

func TestGetPlayers(t *testing.T) {
	g := NewGomegaWithT(t)

	var urlsSearched uint64
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&urlsSearched, 1)

		w.WriteHeader(http.StatusOK)
		resp := []byte(`
{
 "data": {
   "team": {
     "id": 1,
     "name": "TeamName",
     "players": [
       {
         "id": "100",
         "name": "Player1"
       }
     ]
   }
 }
}`)
		w.Write(resp)
	}))

	cfg := &config.Config{
		GetTeamURL: mockServer.URL + "/%d",
	}

	results := make(chan team.Result)
	done := make(chan bool)
	targetSearchRange := 100

	go func() {
		i := 0
		for range results {
			i++
			if i == targetSearchRange {
				done <- true
			}
		}
	}()

	team.GetPlayers(cfg, 5, results, done)

	// Assert that the number of URLs searched is greater than the target search range
	g.Expect(int(urlsSearched) > targetSearchRange).To(BeTrue())
}

func TestGetTeam(t *testing.T) {
	g := NewGomegaWithT(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		resp := []byte(`
{
  "data": {
    "team": {
      "id": 1,
      "name": "TeamName",
      "players": [
        { 
          "id": "100",
          "name": "Player1",
          "age": "30"
        }
      ]
    }
  }
}`)
		w.Write(resp)
	}))
	cfg := &config.Config{
		GetTeamURL: server.URL + "/%d",
	}

	results := make(chan team.Result)
	routines := make(chan int, 1)
	wg := &sync.WaitGroup{}

	routines <- 1
	wg.Add(1)
	go team.GetTeam(wg, cfg, 1, results, routines)

	result := <-results
	g.Expect(result.Err).NotTo(HaveOccurred())
	g.Expect(result.Team).To(Equal(team.Team{
		ID:   1,
		Name: "TeamName",
		Players: []team.Player{
			{
				ID:   "100",
				Name: "Player1",
				Age:  "30",
			},
		},
	}))

	wg.Wait()
	// GetTeam should remove itself from active routines channel
	g.Expect(len(routines)).To(Equal(0))
}
