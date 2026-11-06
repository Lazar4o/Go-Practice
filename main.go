package main

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
)

var ErrNotFound = errors.New("not found")
var ErrSeatTaken = errors.New("seat taken")

type Event struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Venue string `json:"venue"`
}

type Seat struct {
	ID string `json:"id"`
	EventID string `json:"eventId"`
	Row string `json:"row"`
	Number int `json:"number"`
	Status string `json:"status"`
}

type memStore struct {
	seats map[string]Seat
}

func newMemStore() memStore {
	m := make(map[string]Seat)

	for i := 0; i < 20; i++ {
		stringifiedI := strconv.Itoa(i + 1)
		m[stringifiedI] = Seat{
			ID: stringifiedI,
			Status: "free",
			EventID: "20",
			Row: "10",
			Number: i + 1,
		}
	}

	return memStore{seats: m}
}

func (s *memStore) SeatsByEventID(eventID string) []Seat {
	result := []Seat{}
	
	for _, seat := range s.seats {
		if eventID == seat.EventID {
			result = append(result, seat)
		}
	}

	slices.SortFunc(result, func(a, b Seat) int {
		return cmp.Compare(a.Number, b.Number)
	})

	return result
}

func main() {
	store := newMemStore()
	// fmt.Println(store)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /events/{id}/seats", func(w http.ResponseWriter, r *http.Request) {
		eventId := r.PathValue("id")

		status := r.URL.Query().Get("status")

		sortedEvents := store.SeatsByEventID(eventId)

		freeSeats := []Seat{}

		if status == "" {
			freeSeats = sortedEvents[:] // what is the difference with [:] and without it - only eventIDSeats ?? Ask?
		}

		if status != ""{
			for _, s := range sortedEvents {
				currentSeatStatus := s.Status
				if currentSeatStatus == status {
					freeSeats = append(freeSeats, s)
				}
			}
		}

		

		w.Header().Set("Content-type", "application/json")
		err := json.NewEncoder(w).Encode(freeSeats)
		if err != nil {
			log.Printf("encode seats: %v", err)
		}


	})

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
		w.Header().Set("Content-type", "text/plain")
		fmt.Fprintln(w, "ok")
	})
	log.Fatal(http.ListenAndServe(":8080", mux))
}