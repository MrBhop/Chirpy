package handlers

import (
	"net/http"

	"github.com/MrBhop/Chirpy/internal/auth"
)

func (a *ApiConfig) HandlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	if err := a.Db.RevokeToken(r.Context(), token); err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't update token in database", err)
		return
	}

	respondWithJson(w, http.StatusNoContent, struct{}{})
}
