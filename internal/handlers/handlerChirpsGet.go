package handlers

import (
	"net/http"
	"sort"
	"strings"

	"github.com/MrBhop/Chirpy/internal/database"
	"github.com/google/uuid"
)

const ChirpIDParameter string = "chirpID"
const ChirpAuthorIDParameter string = "author_id"
const ChirpSortParameter string = "sort"

func (a *ApiConfig) HandlerChirpsGetAll(w http.ResponseWriter, r *http.Request) {
	authorIdString := r.URL.Query().Get(ChirpAuthorIDParameter)
	sortString := r.URL.Query().Get(ChirpSortParameter)

	sortAscending := strings.ToLower(sortString) != "desc"

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
			return
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

	sort.Slice(output, func(i, j int) bool {
		if sortAscending {
			return output[i].Created_at.Before(output[j].Created_at)
		} else {
			return output[i].Created_at.After(output[j].Created_at)
		}
	})

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
