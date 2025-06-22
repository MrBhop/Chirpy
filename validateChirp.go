package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type setOfString = map[string]struct{}

const maxChirpLength = 140

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	var params parameters
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	msg := params.Body
	if len(msg) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := setOfString{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}

	respondWithJson(w, http.StatusOK, returnVals{
		CleanedBody: cleanWords(msg, badWords),
	})
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
