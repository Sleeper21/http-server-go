package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Sleeper21/http-server/internal/database"
	"github.com/google/uuid"
)

type parameters struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	const maxLength = 140

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Printf("error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError) // --> 500
		return
	}

	if len(params.Body) > maxLength {
		generateErrorJson(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	// Check if user exists
	user, err := cfg.dbQueries.GetUserByID(r.Context(), params.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		log.Printf("user not found: %s", err)
		return
	}

	// Censor profane words
	filteredMsg := hideProfaneWords(params.Body)

	// Create chirp
	returnedChirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   filteredMsg,
		UserID: user.ID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("couldn't create new chirp"))
		log.Printf("error creating chirp: %s", err)
		return
	}

	newChirp := chirp{
		ID:        returnedChirp.ID,
		CreatedAt: returnedChirp.CreatedAt,
		UpdatedAt: returnedChirp.UpdatedAt,
		Body:      returnedChirp.Body,
		UserID:    returnedChirp.UserID,
	}

	respondWithJson(w, http.StatusCreated, newChirp)
}

func hideProfaneWords(msg string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	message := strings.Fields(msg)
	for _, profane := range profaneWords {
		for i, word := range message {
			if strings.ToLower(word) == profane {
				message[i] = "****"
			}
		}
	}
	return strings.Join(message, " ")
}
