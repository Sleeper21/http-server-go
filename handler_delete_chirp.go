package main

import (
	"net/http"

	"github.com/Sleeper21/http-server/internal/auth"
	"github.com/google/uuid"
)

/*

	This endpoint requires an header with an access token jwt in the json format:
	{
		"Authorization": "Bearer <access token>"
	}

	the request also requires a {chirpID} in the path
	Only allow the deletion of a chirp if the user is the author of the chirp.
	the response will return:
		- a 403 status if the user verified in the jwt is not the author of the chirp
		- a 204 status if the chirp is deleted successfully
		- a 404 status if the chirp is not found

*/

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	// Get the Bearer token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		generateErrorJson(w, 401, "unauthorized", err)
		return
	}

	// Validate the jwt
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		generateErrorJson(w, 401, "invalid JWT", err)
		return
	}

	// Get the chirp id from the Path
	strChirpID := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(strChirpID)
	if err != nil {
		generateErrorJson(w, 400, "invalid chirp id", err)
		return
	}

	// Get chirp info from id (from chirp's id)
	chirp, err := cfg.dbQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		generateErrorJson(w, 404, "chirp not found", err)
		return
	}

	// Check if the user is the author of the chirp
	if userID != chirp.UserID {
		generateErrorJson(w, 403, "not allowed, user is not the author of the chirp", err)
		return
	}

	// Delete the chirp
	if err := cfg.dbQueries.DeleteChirpByID(r.Context(), chirpID); err != nil {
		generateErrorJson(w, http.StatusInternalServerError, "couldn't delete chirp", err)
	}

	w.WriteHeader(204)
}
