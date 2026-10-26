package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Event struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Venue string `json:"venue"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /events/{id}", func(w http.ResponseWriter, r *http.Request){
		id := r.PathValue("id")
		if id == "0" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-type", "application/json")

		err := json.NewEncoder(w).Encode(Event{ID: id, Name: "Concert", Venue: "Venue"})
		if err != nil {
			log.Printf("encode event: %v", err)
		}
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-type", "plain/text")
		fmt.Fprintln(w, "ok")
	})
	log.Fatal(http.ListenAndServe(":8080", mux))
}