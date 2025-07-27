package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MrBhop/Chirpy/internal/auth"
	"github.com/MrBhop/Chirpy/internal/database"
)

func (a *ApiConfig) HandlerUsersChange(w http.ResponseWriter, r *http.Request) {
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

	decoder := json.NewDecoder(r.Body)
	var params userAuthParams
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user, err := a.Db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID: userId,
		Email: params.Email,
		HashedPassword: hash,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	respondWithJson(w, http.StatusOK, userReturn{
		Id: user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed:user.IsChirpyRed,
	})
}
