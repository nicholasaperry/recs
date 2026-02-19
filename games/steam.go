package games

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"
)

const (
	ICON_BASE_URL  = "http://media.steampowered.com/steamcommunity/public/images/apps/"
	STEAM_API_KEY  = "STEAM_API_KEY"
	STEAM_BASE_URL = "STEAM_BASE_URL"
	STEAM_USER_ID  = "STEAM_USER_ID"
	STEAM          = "STEAM"
)

type SteamGame struct {
	Id              int    `json:"appid"`
	Name            string `json:"name"`
	IconSlug        string `json:"img_icon_url"`
	PlaytimeWindows int    `json:"playtime_windows_forever"`
	PlaytimeMac     int    `json:"playtime_mac_forever"`
	PlaytimeDeck    int    `json:"playtime_deck_forever"`
	PlaytimeLinux   int    `json:"playtime_linux_forever"`
	LastPlayed      int    `json:"rtime_last_played"`
}

type ApiResponse struct {
	Result struct {
		GameCount int         `json:"game_count"`
		Games     []SteamGame `json:"games"`
	} `json:"response"`
}

func (res ApiResponse) Format() []Game {
	if res.Result.Games == nil {
		return nil
	}
	var games []Game
	for _, steamGame := range res.Result.Games {
		minutesTotal := 0
		minutesTotal += steamGame.PlaytimeWindows + steamGame.PlaytimeMac + steamGame.PlaytimeDeck + steamGame.PlaytimeLinux
		lastPlayedValid := steamGame.LastPlayed > 0
		var lastPlayed sql.NullTime
		if lastPlayedValid {
			lastPlayed = sql.NullTime{
				Time:  time.Unix(int64(steamGame.LastPlayed), 0),
				Valid: true,
			}
		} else {
			lastPlayed = sql.NullTime{}
		}
		games = append(games, Game{
			name:       steamGame.Name,
			appId:      strconv.Itoa(steamGame.Id),
			iconUrl:    ICON_BASE_URL + strconv.Itoa(steamGame.Id) + "/" + steamGame.IconSlug + ".jpg",
			lastPlayed: lastPlayed,
			minutesPlayed: sql.NullInt32{
				Int32: int32(minutesTotal),
				Valid: true,
			},
		})
	}
	return games
}

type SteamReader struct {
	userId  string
	baseUrl string
	apiKey  string
}

func NewSteamReader(userId, baseUrl, apikey string) SteamReader {
	return SteamReader{userId: userId, baseUrl: baseUrl, apiKey: apikey}
}

func (r SteamReader) GetLibrary(userId string) []Game {
	getOwnedGamesUrl := r.baseUrl + "?format=json&include_appinfo=1&key=" + r.apiKey + "&steamid=" + userId
	response, err := http.Get(getOwnedGamesUrl)
	if err != nil {
		println("Error getting library for ", userId, ": ", err.Error())
		panic(err)
	}
	defer response.Body.Close()
	var responseBody ApiResponse
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		println("Error decoding response for ", userId, ": ", err.Error())
		panic(err)
	}
	return sortyByLastPlayed(responseBody.Format())
}

func sortyByLastPlayed(games []Game) []Game {
	sort.Slice(games, func(i, j int) bool {
		return games[i].lastPlayed.Time.Before(games[j].lastPlayed.Time)
	})
	return games
}
