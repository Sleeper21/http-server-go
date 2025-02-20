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

var RefreshTokenExpiration time.Duration = time.Hour * 24 * 60 // --> 60 days

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	input := reqStruct{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error decoding inputs: %s", err)
		return
	}

	// // Check and set expiration time
	// if input.ExpiresIn > 3600 {
	// 	input.ExpiresIn = 3600 * time.Second
	// }
	// // If expiration is not set
	// if input.ExpiresIn == 0 {
	// 	input.ExpiresIn = 3600 * time.Second
	// }

	// Check email in the database to see if the user exists
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

	// Create a JWT token for this user to keep him logged in
	ExpiresIn := 3600 * time.Second // --> 1 hour
	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(ExpiresIn))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error creating session token. "))
		log.Printf("error creating JWT token: %s", err)
		return
	}

	// Create a refresh token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	// Store it in the database
	err = cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(RefreshTokenExpiration),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error saving refresh token in database: %s\n", err)
		return
	}

	// If all good, return a json with the user data, except the password
	type responseFields struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
		IsChirpyRed  bool      `json:"is_chirpy_red"`
	}

	loggedUser := responseFields{
		user.ID,
		user.CreatedAt,
		user.UpdatedAt,
		user.Email,
		token,
		refreshToken,
		user.IsChirpyRed,
	}

	respondWithJson(w, 200, loggedUser)
}
