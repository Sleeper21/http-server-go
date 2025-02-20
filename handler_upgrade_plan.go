package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

/*

The endpoint will upgrade the user "is_chirpy_red" plan.
The endpoint should accept a request of this shape:
		{
		  "event": "user.upgraded",
		  "data": {
		    "user_id": "3311741c-680c-4546-99f3-fc9efac2036c"
		  }
		}

- If the event is anything other than user.upgraded, the endpoint should immediately respond with a 204 status code - we don’t care about any other events.

- If the event is user.upgraded, then it should update the user in the database, and mark that they are a Chirpy Red member.

- If the user is upgraded successfully, the endpoint should respond with a 204 status code and an empty response body. If the user can’t be found, the endpoint should respond with a 404 status code.

*/

type WebHookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

type resEmptyBody struct{}

func (cfg *apiConfig) handlerUpgradeUserPlan(w http.ResponseWriter, r *http.Request) {

	// Decode body to struct
	req := WebHookRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		generateErrorJson(w, http.StatusBadRequest, "couldn't decode request", err)
		return
	}

	// Validate the event and upgrade the user
	if req.Event == "user.upgraded" {

		// Check if user exists
		validUser, err := cfg.dbQueries.GetUserByID(r.Context(), req.Data.UserID)
		if err != nil {
			generateErrorJson(w, 404, "user not found", err)
			return
		}

		// Upgrade the user "is_chirpy_red" column in db
		err = cfg.dbQueries.UpgradeUserPlan(r.Context(), validUser.ID)
		if err != nil {
			generateErrorJson(w, http.StatusInternalServerError, "couldn't upgrade the user's plan", err)
			return
		}

		// responde with no content and empty body
		respondWithJson(w, 204, resEmptyBody{})
	}

	// If the event is not user.upgraded
	respondWithJson(w, 204, resEmptyBody{})

}
