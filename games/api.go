package games

import (
	"encoding/json"
	"net/http"
	"os"
	"playtime/auth"
)

func HandleGetLibrary(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userId := auth.GetSteamID(r)
	if userId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"not logged in"}`))
		return
	}
	apiKey := os.Getenv(STEAM_API_KEY)
	baseURL := os.Getenv(STEAM_BASE_URL)
	if apiKey == "" || baseURL == "" {
		http.Error(w, "STEAM_API_KEY or STEAM_BASE_URL is not set", http.StatusInternalServerError)
		return
	}
	games := NewSteamReader(userId, baseURL, apiKey).GetLibrary(userId)
	gameMaps := make([]map[string]interface{}, len(games))
	for i, game := range games {
		gameMaps[i] = game.ToMap()
	}
	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"games": gameMaps,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
