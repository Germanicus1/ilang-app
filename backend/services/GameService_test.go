package services

import (
	"log"
	"testing"
)

func TestFetchGames(t *testing.T) {
	userID := "5803acaf-821a-4463-b8b4-15ac6e0e466a"
	// Mock the HTTP response from Supabase
	// server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte(`[{"id":"1","title":"Game 1"}]`))
	// }))
	// defer server.Close()

	games, err := FetchGames(userID)
	if err != nil {
		t.Fatalf("FetchGames failed: %v", err)
	}
	log.Println(games)
	if len(games) != 2 {
		t.Errorf("FetchGames failed: %v", err)
	}
}
