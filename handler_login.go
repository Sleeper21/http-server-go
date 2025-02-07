package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Sleeper21/http-server/internal/auth"
	"github.com/google/uuid"
)

// User tries to login
// request body should contain:
//      {
//      	"email": "somename@example.com"
//      	"password": "04234",
//      }

type reqStruct struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	input := reqStruct{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error decoding inputs: %s", err)
		return
	}

	// Check email in the database
	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), input.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("email not found"))
		log.Printf("couldn't find email in database: %s", err)
		return
	}

	// Compare the password input with the hashed password stored
	err = auth.CheckPasswordHash(user.HashedPassword, input.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("incorrect password"))
		log.Printf("passwords don't match: %s", err)
		return
	}

	// If all good, return a json with the user data, except the password
	type responseFields struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Password  string    `json:"hashed_password"`
	}

	loggedUser := responseFields{
		user.ID,
		user.CreatedAt,
		user.UpdatedAt,
		user.Email,
		user.HashedPassword,
	}

	respondWithJson(w, 200, loggedUser)
}
