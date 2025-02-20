package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Sleeper21/http-server/internal/auth"
	"github.com/Sleeper21/http-server/internal/database"
	"github.com/google/uuid"
)

type BodyUserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	data := json.NewDecoder(r.Body)
	dataJSON := BodyUserCredentials{}
	err := data.Decode(&dataJSON)
	if err != nil {
		log.Printf("error decoding the data: %s", err)
		return
	}

	if dataJSON.Email == "" {
		log.Println("email cannot be empty")
		return
	}
	if dataJSON.Password == "" {
		log.Println("password cannot be empty")
	}

	//Hash password before storing
	hashedPassword, err := auth.HashPassword(dataJSON.Password)
	if err != nil {
		log.Printf("error hashing password: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// db query to create user
	userToSave := database.CreateUserParams{
		Email:          dataJSON.Email,
		HashedPassword: hashedPassword,
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), userToSave)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error creating the user: %s", err)
		return
	}
	// parse the returned user to JSON and send it in response writer
	// excluding the hashed password
	type responseFields struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}

	newUser := responseFields{
		user.ID,
		user.CreatedAt,
		user.UpdatedAt,
		user.Email,
		user.IsChirpyRed,
	}

	respondWithJson(w, http.StatusCreated, newUser)
}
