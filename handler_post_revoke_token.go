package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Sleeper21/http-server/internal/auth"
	"github.com/Sleeper21/http-server/internal/database"
)

/*
This endpoint resets the user's refresh token and stores it in the db
This new endpoint does not accept a request body, but does require a refresh token to be present in the headers, in the same Authorization: Bearer <token> format.
*/

func (cfg *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {

	strToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("error getting authorization header: %s", err)
		return
	}

	// Check if the token retrieved from req header exists in database
	res, err := cfg.dbQueries.GetRefreshTokenByToken(r.Context(), strToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("refresh token not found: %s", err)
		return
	}

	// Generate new refresh token
	newToken, err := auth.MakeRefreshToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update the refresh token
	newTokenParams := database.UpdateRefreshTokenParams{
		Token:     newToken,
		ExpiresAt: time.Now().Add(RefreshTokenExpiration),
		Token_2:   res.Token,
	}

	err = cfg.dbQueries.UpdateRefreshToken(r.Context(), newTokenParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("couldn't update the refresh token: %s", err)
		return
	}

	w.WriteHeader(204) // --> A 204 status means the request was successful but no body is returned.
}
