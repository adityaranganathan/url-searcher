package team

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"team-cli/config"
	"time"
)

type Team struct {
	ID      int
	Name    string
	Players []Player
}

type Player struct {
	ID   string
	Name string
	Age  string
}

type Result struct {
	Team Team
	Err  error
}

// GetPlayers concurrently retrieves data for multiple teams using a specified max number of goroutines.
// It stops when it receives a signal from the done channel and waits for existing goroutines to finish.
func GetPlayers(cfg *config.Config, maxRoutines int, results chan Result, done chan bool) {
	var wg sync.WaitGroup

	routines := make(chan int, maxRoutines)
	teamID := 0
loop:
	for {
		select {
		case routines <- 1:
			wg.Add(1)
			go GetTeam(&wg, cfg, teamID, results, routines)
			teamID++
		case <-done:
			break loop

		}
	}
	wg.Wait()
}

// GetTeam retrieves data for one team and sends the result to the results channel
// It also removes 1 value from the routines channel that tracks running goroutines.
func GetTeam(wg *sync.WaitGroup, cfg *config.Config, teamID int, results chan Result, routines chan int) {
	defer func() {
		<-routines
		wg.Done()
	}()

	client := http.Client{
		Timeout: time.Duration(cfg.TimeoutInMilliSecond) * time.Millisecond,
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(cfg.GetTeamURL, teamID), nil)
	if err != nil {
		results <- Result{Err: fmt.Errorf("creating request: %v", err)}
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		results <- Result{Err: fmt.Errorf("getting team data: %v", err)}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		results <- Result{}
		return
	}

	var responseData GetTeamResponse
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		results <- Result{Err: fmt.Errorf("decoding response body: %v", err)}
		return
	}

	results <- Result{Team: responseData.Data.Team}
}

type GetTeamResponse struct {
	Data Data
}

type Data struct {
	Team Team
}
