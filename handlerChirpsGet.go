package main

import (
	"net/http"

	"github.com/google/uuid"
)

const chirpIDParameter string = "chirpID"

func (a *apiConfig) handlerChirpsGetAll(w http.ResponseWriter, r *http.Request) {
	chirps, err := a.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't fetch chirps", err)
		return
	}

	output := []chirp{}
	for _, c := range chirps {
		output = append(output, chirp{
			Id: c.ID,
			Created_at: c.CreatedAt,
			Updated_at: c.UpdatedAt,
			Body: c.Body,
			UserId: c.UserID,
		})
	}

	respondWithJson(w, http.StatusOK, output)
}

func (a *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue(chirpIDParameter))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "id could not be parsed to uuid", err)
		return
	}

	responseChirp, err := a.db.GetChirpById(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't fetch chirp", err)
		return
	}

	respondWithJson(w, http.StatusOK, chirp{
		Id: responseChirp.ID,
		Created_at: responseChirp.CreatedAt,
		Updated_at: responseChirp.UpdatedAt,
		Body: responseChirp.Body,
		UserId: responseChirp.UserID,
	})
}
