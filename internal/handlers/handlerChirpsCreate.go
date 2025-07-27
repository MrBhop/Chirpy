package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MrBhop/Chirpy/internal/auth"
	"github.com/MrBhop/Chirpy/internal/database"
	"github.com/google/uuid"
)

type setOfString = map[string]struct{}

type chirp struct {
	Id         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Body       string    `json:"body"`
	UserId     uuid.UUID `json:"user_id"`
}

const maxChirpLength = 140

func (a *ApiConfig) HandlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

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
	var params parameters
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	msg := params.Body
	cleanedMessage, err := validateChirp(msg)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	newChirp, err := a.Db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleanedMessage,
		UserID: userId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJson(w, http.StatusCreated, chirp{
		Id: newChirp.ID,
		Created_at: newChirp.CreatedAt,
		Updated_at: newChirp.UpdatedAt,
		Body: newChirp.Body,
		UserId: newChirp.UserID,
	})
}

func validateChirp(msg string) (string, error) {
	if len(msg) > maxChirpLength {
		return "", fmt.Errorf("Chirp is too long")
	}

	badWords := setOfString{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}

	return cleanWords(msg, badWords), nil
}

func cleanWords(msg string, wordsToExclude setOfString) string {
	const delimiter = " "
	words := strings.Split(msg, delimiter)

	for i, word := range words {
		if _, exists := wordsToExclude[strings.ToLower(word)]; exists {
			words[i] = "****"
		}
	}

	return strings.Join(words, delimiter)
}
