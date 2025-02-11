package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Sleeper21/http-server/internal/auth"
)

// The header should contain a refresh token in the same Authorization: Bearer <token> format:
/* return a new access token in this format:
{
 	"token": "eyJhbGciOiJIUzII6IkpXVCJ9.eyJzdWIiOiIxwIiwibmFtZSI6I"
}
*/

/*
Each user has a refresh token stored in the database
this endpoint will check the refresh token, if valid the respective user will receive a new access-token valid for one hour.
*/

func (cfg *apiConfig) handlerResetAccessToken(w http.ResponseWriter, r *http.Request) {

	// Check if the refresh token sent in the header is valid for this user
	refToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("error: ", err)
		return
	}

	// Get user from refresh token from database
	user, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), string(refToken))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("error finding user in the database: ", err)
		return
	}
	// Check if the refresh token is expired
	if user.ExpiresAt.Before(time.Now()) {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("error, refresh token expired")
		return
	}

	// if the refresh token is valid and not expired,
	// Create a new access token jwt that expires in 1 hour
	token, err := auth.MakeJWT(user.UserID, cfg.secret, 3600*time.Second)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error creating jwt token: %s ", err)
		return
	}

	type responseStruct struct {
		Token string `json:"token"`
	}

	response := responseStruct{
		token,
	}

	respondWithJson(w, 200, response)
}
