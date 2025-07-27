package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type webHooksParams struct {
	Event string `json:"event"`
	Data  struct {
		UserId uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (a *ApiConfig) HandlerWebHooks(w http.ResponseWriter, r *http.Request) {
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

	_, err := a.Db.EnableUserSubscription(r.Context(), params.Data.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user subscription", err)
		return
	}

	respondWithJson(w, http.StatusNoContent, struct{}{})
}
