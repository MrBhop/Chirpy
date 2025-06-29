package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MrBhop/Chirpy/internal/auth"
)

func (a *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params userAuthParams
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := a.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithUnauthorized(w, err)
		return
	}

	if err := auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
		respondWithUnauthorized(w, err)
		return
	}

	var expiresIn time.Duration
	if params.ExpiresInSeconds == 0 {
		expiresIn = 1 * time.Hour
	} else {
		expiresIn = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	token, err := auth.MakeJWT(user.ID, a.secret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT", err)
		return
	}

	respondWithJson(w, http.StatusOK, userReturn{
		Id: user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email: user.Email,
		Token: token,
	})
}

func respondWithUnauthorized(w http.ResponseWriter, err error) {
	respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
}
