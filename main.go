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
	r.HandleFunc("/config", configUIHandler)
	r.HandleFunc("/api/configs", listConfigsHandler).Methods("GET")
	r.HandleFunc("/api/configs/{name}", getConfigHandler).Methods("GET")
	r.HandleFunc("/api/configs", saveConfigHandler).Methods("POST")
	r.HandleFunc("/api/configs/{name}", deleteConfigHandler).Methods("DELETE")

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
