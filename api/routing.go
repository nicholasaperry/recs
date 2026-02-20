package api

import (
	"net/http"
	"sommelier/auth"
	"sommelier/games"
)

func RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/api/auth/steam/login", auth.LoginHandler)
	router.HandleFunc("/api/auth/steam/callback", auth.CallbackHandler)
	router.HandleFunc("/api/auth/me", auth.MeHandler)
	router.HandleFunc("/api/auth/logout", auth.LogoutHandler)
	router.HandleFunc("/api/games/library", games.HandleGetLibrary)
	router.Handle("/", http.FileServer(http.Dir("frontend/dist")))
}
