package handlers

import (
	"net/http"

	"github.com/MrBhop/Chirpy/internal/database"
	"github.com/google/uuid"
)

const ChirpIDParameter string = "chirpID"
const ChirpAuthorIDParameter string = "author_id"

func (a *ApiConfig) HandlerChirpsGetAll(w http.ResponseWriter, r *http.Request) {
	authorIdString := r.URL.Query().Get(ChirpAuthorIDParameter)

	var chirps []database.Chirp
	if authorIdString == "" {
		dbChirps, err := a.Db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't fetch chirps", err)
			return
		}
		chirps = dbChirps
	} else {
		authorId, err := uuid.Parse(authorIdString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "could not parse author_id to uuid", err)
		}
		dbChirps, err := a.Db.GetChirpByAuthorId(r.Context(), authorId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't fetch chirps", err)
			return
		}
		chirps = dbChirps
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

func (a *ApiConfig) HandlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue(ChirpIDParameter))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "id could not be parsed to uuid", err)
		return
	}

	responseChirp, err := a.Db.GetChirpById(r.Context(), id)
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
