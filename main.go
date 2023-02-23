package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"team-cli/application"
	"team-cli/cmd"
)

func main() {
	err := Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Run() error {
	var maxRoutines int
	var runCLI bool
	flag.IntVar(&maxRoutines, "max_routines", 1, "number of goroutines to use")
	flag.BoolVar(&runCLI, "run_cli", true, "run as cli")
	flag.Parse()

	playerFetcher := application.NewPlayerFetcher(maxRoutines)

	var err error
	if runCLI {
		err = cmd.RunCLI(playerFetcher)
		return err
	} else {
		router := setupRoutes(playerFetcher)
		http.ListenAndServe(":9005", router)
	}

	return nil
}

func setupRoutes(fetcher *application.PlayerFetcher) *mux.Router {
	muxRouter := mux.NewRouter()
	muxRouter.Methods("GET").Path("/players").HandlerFunc(cmd.GetHandleTeamRequest(fetcher))
	muxRouter.Methods("GET").Path("/live").HandlerFunc(cmd.GetRespondOKHandler())
	muxRouter.Methods("GET").Path("/ready").HandlerFunc(cmd.GetRespondOKHandler())

	return muxRouter
}
