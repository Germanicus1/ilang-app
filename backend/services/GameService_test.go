package services

import (
	"testing"
)

func TestFetchGames(t *testing.T) {
	// Mock the HTTP response from Supabase
	// server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte(`[{"id":"1","title":"Game 1"}]`))
	// }))
	// defer server.Close()

	games, err := FetchGames()
	if err != nil {
		t.Fatalf("FetchGames failed: %v", err)
	}

	if len(games) == 0 {
		t.Errorf("No games found: %v", err)
	}
}
