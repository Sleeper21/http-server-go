package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Email string `json:"email"`
	}

	data := json.NewDecoder(r.Body)
	dataJSON := body{}
	err := data.Decode(&dataJSON)
	if err != nil {
		log.Printf("error decoding the data: %s", err)
		return
	}

	if dataJSON.Email == "" {
		log.Printf("email cannot be empty")
		return
	}

	// db query to add create user
	user, err := cfg.dbQueries.CreateUser(r.Context(), dataJSON.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error creating the user: %s", err)
		return
	}

	// parse the returned user to JSON and send it in response writer
	type userStruct struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	newUser := userStruct{
		user.ID,
		user.CreatedAt,
		user.UpdatedAt,
		user.Email,
	}

	respondWithJson(w, http.StatusCreated, newUser)
}
