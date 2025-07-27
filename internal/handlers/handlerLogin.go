package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MrBhop/Chirpy/internal/auth"
)

func (a *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params userAuthParams
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := a.Db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithUnauthorized(w, err)
		return
	}

	if err := auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
		respondWithUnauthorized(w, err)
		return
	}

	expiresIn := 1 * time.Hour
	token, err := auth.MakeJWT(user.ID, a.Secret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT", err)
		return
	}

	refreshToken, err := auth.MakeRegisteredRefreshToken(a.Db, r, user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	respondWithJson(w, http.StatusOK, userReturn{
		Id: user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email: user.Email,
		Token: token,
		RefreshToken: refreshToken,
	})
}

func respondWithUnauthorized(w http.ResponseWriter, err error) {
	respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
}
