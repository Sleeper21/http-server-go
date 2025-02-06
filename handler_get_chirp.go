package main

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {

	strID := r.PathValue("chirpID") // --> will return e string "12eae124-2345235-5235"

	// Convert id (string) to a uuid.UUID to match the db struct
	id, err := uuid.Parse(strID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid id format"))
		log.Printf("error: %s", err)
		return
	}

	// Get from database
	c, err := cfg.dbQueries.GetChirpByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(404)
		log.Printf("chirp not found : %s", err)
		return
	}

	chirp := Chirp{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserID:    c.UserID,
	}

	respondWithJson(w, 200, chirp)
}

// Get all chirps
func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {

	allChirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("couldn't retrieve chirps from database: %s", err)
		return
	}

	list := []Chirp{}

	for _, chirp := range allChirps {
		chirp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		list = append(list, chirp)
	}
	respondWithJson(w, 200, list)
}
