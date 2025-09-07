package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

func rssHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Printf("name: %s\n", name)
	if name != "" {
		res, err := GetByName(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/rss+xml")
		fmt.Fprint(w, string(res))
	} else {
		url := r.URL.Query().Get("url")
		res, err := GetByURL(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprint(w, string(res))
	}
}

func configUIHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(".", "config.html"))
}

func listConfigsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dbClient, err := NewDBClient(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dbClient.Close()

	names, err := dbClient.ListConfigs(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(names)
}

func getConfigHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	ctx := r.Context()
	dbClient, err := NewDBClient(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dbClient.Close()

	config, err := dbClient.GetConfig(ctx, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

func saveConfigHandler(w http.ResponseWriter, r *http.Request) {
	var config Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	dbClient, err := NewDBClient(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dbClient.Close()

	if err := dbClient.SaveConfig(ctx, &config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteConfigHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	ctx := r.Context()
	dbClient, err := NewDBClient(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dbClient.Close()

	if err := dbClient.DeleteConfig(ctx, name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
