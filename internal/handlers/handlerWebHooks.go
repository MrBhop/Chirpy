package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MrBhop/Chirpy/internal/auth"
	"github.com/google/uuid"
)

type webHooksParams struct {
	Event string `json:"event"`
	Data  struct {
		UserId uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (a *ApiConfig) HandlerWebHooks(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find api key", err)
		return
	}

	if apiKey != a.PolkaApiKey {
		err := errors.New("Invalid ApiKey")
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var params webHooksParams
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJson(w, http.StatusNoContent, struct{}{})
		return
	}

	if _, err := a.Db.EnableUserSubscription(r.Context(), params.Data.UserId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user subscription", err)
		return
	}

	respondWithJson(w, http.StatusNoContent, struct{}{})
}
