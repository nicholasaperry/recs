package main

import (
	"os"
	"playtime/api"

	"github.com/joho/godotenv"
)

func main() {
	// load env
	err := godotenv.Load()
	if err != nil {
		println("Error loading .env file ", err.Error())
	}
	api.Run(os.Getenv("PORT"))
	// stores := map[string]games.Reader{
	// 	games.STEAM: games.NewSteamReader(os.Getenv(games.STEAM_USER_ID), os.Getenv(games.STEAM_BASE_URL), os.Getenv(games.STEAM_API_KEY)),
	// }
	// // step 1: which game stores are we checking?
	// featuresEnabled := map[string]bool{}
	// // start with only checking steam
	// steamApiKey := os.Getenv(games.STEAM_API_KEY)
	// if steamApiKey != "" {
	// 	featuresEnabled[games.STEAM] = true
	// }
	// // call the relevant method for each registered playtime provider
	// for library, reader := range stores {
	// 	if featuresEnabled[library] {
	// 		println("Getting library for ", library)
	// 		userId := os.Getenv(games.STEAM_USER_ID)
	// 		games := reader.GetLibrary(userId)
	// 		for _, game := range games {
	// 			fmt.Println(game.String())
	// 		}
	// 	}
	// }
}
