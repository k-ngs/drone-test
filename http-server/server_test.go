package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const (
	jsonContentType = "application/json"
)

func TestGetPlayers(t *testing.T) {
	store := StubPlayerStore{map[string]int{
		"Pepper": 20,
		"Floyd": 10,
	}, nil, nil}

	server := NewPlayerServer(&store)

	tests := []struct{
		name string
		player string
		expectedStatusCode int
		expectedScore string
	}{
		{
			name: "returns Pepper's score",
			player: "Pepper",
			expectedStatusCode: http.StatusOK,
			expectedScore: "20",
		},
		{
			name: "returns Floyd's score",
			player: "Floyd",
			expectedStatusCode: http.StatusOK,
			expectedScore: "10",
		},
		{
			name: "return 404 on missing players",
			player: "Apollo",
			expectedStatusCode: http.StatusNotFound,
			expectedScore: "0",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := newGetScoreRequest(test.player)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertStatusCode(t, response.Code, test.expectedStatusCode)
			assertResponseBody(t, response.Body.String(), test.expectedScore)
		})
	}
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{map[string]int{}, nil, nil}
	server := NewPlayerServer(&store)

	t.Run("it returns wins on POST", func(t *testing.T) {
		player := "Pepper"
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatusCode(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}
		if store.winCalls[0] != player {
			t.Errorf("did not store correct winner got %q, want %q", store.winCalls[0], player)
		}
	})
}

func TestLeague(t *testing.T) {
	t.Run("it return 200 on /league", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		store := StubPlayerStore{nil, nil, wantedLeague}
		server := NewPlayerServer(&store)

		response := httptest.NewRecorder()
		request := newLeagueRequest()

		server.ServeHTTP(response, request)
		got := getLeagueFromResponse(t, response.Body)

		assertStatusCode(t, response.Code, http.StatusOK)
		assertLeague(t, got, wantedLeague)
		assertContentType(t, response, jsonContentType)
	})
}

type StubPlayerStore struct {
	scores map[string]int
	winCalls []string
	league []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() []Player {
	return s.league
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	 return req
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func getLeagueFromResponse(t *testing.T, body io.Reader) (league []Player) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)
	if err != nil {
		t.Fatalf("failed to parse json %q, %v", body, err)
	}
	return
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertStatusCode(t *testing.T, got, want int ) {
	t.Helper()
	if got != want {
		t.Errorf("got status code %d, want status code %d", got, want)
	}
}

func assertLeague(t *testing.T, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Header().Get("content-type") != want {
		t.Errorf("response did't have application/json, got %v", response.Result().Header.Get("content-type"))
	}

}