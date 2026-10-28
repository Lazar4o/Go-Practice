package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
)

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


func main() {
	var seats []Seat // i can do like this the empty "slice" seats

	for i := 0; i < 20; i++{
		stringedI := strconv.Itoa(i + 1)

		status := []string{"free", "occupied"}[rand.IntN(2)]

		s := Seat{
			ID: stringedI,
			EventID: "20",
			Row: "A",
			Number: i + 1,
			Status: status,
		}
		seats = append(seats, s)
	}

	// LOOP - and PRINT the "seats" - 1st way to do it with for i := 0; i < len(seats); i++ {} 
	// for i := 0; i < len(seats); i++ {
	// 	fmt.Println(seats[i])
	// }

	// LOOP - and PRINT the "seats" - 2nd way to do it with range - and forr is the snippet
	// for _, seat := range seats {
	// 	fmt.Println(seat)
	// }

	//HIGHLIGHTED BLUE ERROR: unused write to field Status #unusedwritedefault
	// for _, seat := range seats {
	// 	seat.Status = "free"
	// }
	// fmt.Println(seats)

	for i := range seats {
		seats[i].Status = "kur"
	}
	fmt.Println(seats)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /events/{id}/seats", func(w http.ResponseWriter, r *http.Request) {
		eventId := r.PathValue("id")

		status := r.URL.Query().Get("status")

		// fmt.Fprintf(w, "We look for seats for event with ID: %s, filtered by status: %s", eventId, status)

		eventIDSeats := []Seat{}

		for _, s := range seats {
			if s.EventID == eventId {
				eventIDSeats = append(eventIDSeats, s)
			}
		}

		freeSeats := []Seat{}

		if status == "" {
			freeSeats = eventIDSeats[:] // what is the difference with [:] and without it - only eventIDSeats ?? Ask?
		}

		if status != ""{
			for _, s := range eventIDSeats {
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