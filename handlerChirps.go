package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MrBhop/Chirpy/internal/database"
	"github.com/google/uuid"
)

type setOfString = map[string]struct{}

const maxChirpLength = 140

func (a *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	type returnVals struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Body       string    `json:"body"`
		UserId     uuid.UUID `json:"user_id"`
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

	user_id := params.UserId
	chirp, err := a.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleanedMessage,
		UserID: user_id,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJson(w, http.StatusCreated, returnVals{
		Id: chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Body: chirp.Body,
		UserId: chirp.UserID,
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
