package handlers

import (
	"errors"
	"net/http"

	"github.com/MrBhop/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (a *ApiConfig) HandlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	userId, err := auth.ValidateJWT(token, a.Secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Token", err)
		return
	}

	chirpId, err := uuid.Parse(r.PathValue(ChirpIDParameter))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "id could not be parsed to uuid", err)
		return
	}

	chirp, err := a.Db.GetChirpById(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldn't find chirp matching id", err)
		return
	}

	if chirp.UserID != userId {
		err := errors.New("you can't delete this chirp")
		respondWithError(w, http.StatusForbidden, err.Error(), err)
		return
	}

	if err := a.Db.DeleteChirpById(r.Context(), chirpId); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't delete chirp", err)
		return
	}

	respondWithJson(w, http.StatusNoContent, struct{}{})
}
