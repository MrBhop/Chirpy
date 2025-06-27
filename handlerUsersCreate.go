package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func handlerUsersCreate(config *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Email string `json:"email"`
		}
		type returnValues struct {
			Id uuid.UUID `json:"id"`
			Email string `json:"email"`
			Created_at time.Time `json:"created_at"`
			Updated_at time.Time `json:"updated_at"`
		}

		decoder := json.NewDecoder(r.Body)
		var params parameters
		if err := decoder.Decode(&params); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
			return
		}

		email := params.Email
		user, err := config.db.CreateUser(r.Context(), email)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
			return
		}

		respondWithJson(w, http.StatusCreated, returnValues{
			Id: user.ID,
			Email: user.Email,
			Created_at: user.CreatedAt,
			Updated_at: user.UpdatedAt,
		})
	}
}
