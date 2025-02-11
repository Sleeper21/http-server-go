package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Sleeper21/http-server/internal/auth"
	"github.com/Sleeper21/http-server/internal/database"
	"github.com/google/uuid"
)

/*
PUT /api/users endpoint so that users can update their own (but not others’) email and password. It requires:
  - An access token in the header
  - A new password and email in the request body (both new)
*/
func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {

	// Get Bearer access token
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		generateErrorJson(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	// Validate access token
	userID, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		generateErrorJson(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	// Treat the body request data received
	var reqBody BodyUserCredentials
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		generateErrorJson(w, http.StatusInternalServerError, "couldn't decode data", err)
		return
	}

	// Hash the new password
	hashedPassword, err := auth.HashPassword(reqBody.Password)
	if err != nil {
		generateErrorJson(w, 500, "couldn't process request", err)
		return
	}

	// update users credentials
	updatedUser, err := cfg.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          reqBody.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	})
	if err != nil {
		generateErrorJson(w, 500, "couldn't update credentials", err)
		return
	}

	// Respond with returned data - updated user. this does not return the hashed password
	type updatedUserJson struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	respondWithJson(w, 200, updatedUserJson{
		ID:        updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
	})
}
