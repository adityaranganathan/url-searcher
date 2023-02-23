package cmd_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"team-cli/cmd"
	"team-cli/cmd/mocks"
	"testing"
)

func TestName(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost/players", nil)
	fetcher := mocks.NewPlayerFetcher(t)
	fetcher.On("GetTeamPlayers", map[string]bool{}).Return(nil, nil)

	rr := httptest.NewRecorder()
	cmd.GetHandleTeamRequest(fetcher).ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatal(fmt.Sprintf("expected %v, got %v", 200, rr.Result().StatusCode))
	}
}
