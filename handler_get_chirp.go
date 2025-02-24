package main

import (
	"errors"
	"log"
	"net/http"
	"sort"
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

// Get chirp by its Id
func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {

	strID := r.PathValue("chirpID") // --> will return string "12eae124-2345235-5235"

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
// Optional query params:
// ?author_id=<userid> --> retrieves all chirps from that user
// ?sort=asc --> sort the chirps by 'created_at' in ascending order
// ?sort=desc --> sort the chirps by 'created_at' in descending order
// Default sort will be asc by created_at

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {

	// Check if optional query params were provided
	stringID := r.URL.Query().Get("author_id") // --> returns a string id if existent
	sorting := r.URL.Query().Get("sort")

	if stringID != "" {
		// Covert id string into uuid.UUID to match the db struct
		userID, err := uuid.Parse(stringID)
		if err != nil {
			generateErrorJson(w, http.StatusBadRequest, "invalid id format", err)
			return
		}

		// Check if its a valid user id
		user, err := cfg.dbQueries.GetUserByID(r.Context(), userID)
		if err != nil {
			generateErrorJson(w, http.StatusNotFound, "user not found", err)
			return
		}

		// Get all chirps from that user
		chirpsDefault, err := cfg.dbQueries.GetChirpsByUserID(r.Context(), user.ID)
		if err != nil {
			generateErrorJson(w, http.StatusInternalServerError, "couldn't retrieve any chirp from user", err)
			return
		}

		// Convert database.Chirp into Chirp type because of Json capability
		chirpsDefaultJSON := []Chirp{}

		for _, chirp := range chirpsDefault {
			chirp := Chirp{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			}
			chirpsDefaultJSON = append(chirpsDefaultJSON, chirp)
		}

		// handle sort param
		if sorting != "" {
			err = validateSortQuery(sorting) // Validate order of sorting
			if err != nil {
				generateErrorJson(w, http.StatusBadRequest, "couldn't sort", err)
				return
			}

			if sorting == "desc" {
				sortedChirps := sortDataDesc(chirpsDefaultJSON)
				respondWithJson(w, 200, sortedChirps)
				return
			}
		}

		// return User's chirps
		respondWithJson(w, 200, chirpsDefault)
		return
	}

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

	// Handle sort param
	if sorting != "" {
		err = validateSortQuery(sorting)
		if err != nil {
			generateErrorJson(w, http.StatusBadRequest, "couldn't sort", err)
			return
		}

		if sorting == "desc" {
			sortedChirps := sortDataDesc(list)
			respondWithJson(w, 200, sortedChirps)
			return
		}
	}

	respondWithJson(w, 200, list)
}

func sortDataDesc(chirps []Chirp) []Chirp {
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
	})
	return chirps
}

func validateSortQuery(order string) error {
	if order != "desc" && order != "asc" {
		return errors.New("invalid sort parameter. Use 'asc' or 'desc' only")
	}
	return nil
}
