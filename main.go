package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()

	r.HandleFunc("/", rssHandler)
	configRouter := r.PathPrefix("/config").Subrouter()
	configRouter.Use(authMiddleware)
	configRouter.HandleFunc("", configUIHandler)

	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.Use(authMiddleware)
	apiRouter.HandleFunc("/configs", listConfigsHandler).Methods("GET")
	apiRouter.HandleFunc("/configs/{name}", getConfigHandler).Methods("GET")
	apiRouter.HandleFunc("/configs", saveConfigHandler).Methods("POST")
	apiRouter.HandleFunc("/configs/{name}", deleteConfigHandler).Methods("DELETE")

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
