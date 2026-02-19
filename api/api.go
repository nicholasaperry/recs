package api

import (
	"log"
	"net/http"
)

func Run(port string) {
	if port == "" {
		port = "8080"
	}
	router := http.NewServeMux()
	RegisterRoutes(router)
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal("Error starting API: ", err)
	}
	log.Println("API started on port ", port)
}
