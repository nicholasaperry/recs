package games

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Game struct {
	name          string
	lastPlayed    sql.NullTime
	minutesPlayed sql.NullInt32
	appId         string
	iconUrl       string
}

func (g Game) String() string {
	return fmt.Sprintf("Game: %s, Last Played: %s, Minutes Played: %d, Icon Url: %s", g.name, g.lastPlayed.Time.Format("01/02/2006"), g.minutesPlayed.Int32, g.iconUrl)
}

func (g Game) ToJson() string {
	jsonBytes, err := json.Marshal(g.ToMap())
	if err != nil {
		println("Error marshalling game to JSON: ", err.Error())
		return ""
	}
	return string(jsonBytes)
}

// ToMap returns the game as a map for JSON encoding (e.g. API response).
func (g Game) ToMap() map[string]any {
	lastPlayed := ""
	if g.lastPlayed.Valid {
		lastPlayed = g.lastPlayed.Time.Format(time.RFC3339)
	}
	minutesPlayed := int32(0)
	if g.minutesPlayed.Valid {
		minutesPlayed = g.minutesPlayed.Int32
	}
	return map[string]any{
		"name":          g.name,
		"lastPlayed":    lastPlayed,
		"minutesPlayed": minutesPlayed,
		"appId":         g.appId,
		"iconUrl":       g.iconUrl,
	}
}

type Reader interface {
	GetLibrary(userId string) []Game
}
