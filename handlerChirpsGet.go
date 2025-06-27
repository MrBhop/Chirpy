package main

import "net/http"

func (a *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirps, err := a.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't fetch chirps", err)
	}

	output := []chirp{}
	for _, c := range chirps {
		output = append(output, chirp{
			Id: c.ID,
			Created_at: c.CreatedAt,
			Updated_at: c.UpdatedAt,
			Body: c.Body,
			UserId: c.UserID,
		})
	}

	respondWithJson(w, http.StatusOK, output)
}
