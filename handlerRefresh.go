package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MrBhop/Chirpy/internal/auth"
)

func (a *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	user, err := a.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Couldn't fetch token from db: %v", err), err)
		return
	}

	expiresIn := 1 * time.Hour
	newToken, err := auth.MakeJWT(user.ID, a.secret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create new token", err)
		return
	}

	respondWithJson(w, http.StatusOK, response{
		Token: newToken,
	})
}
