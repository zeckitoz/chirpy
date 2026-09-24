package main

import (
	"net/http"
)

func (cfg *apiConfig) listChirpsHandler(w http.ResponseWriter, req *http.Request) {

	raw_chirps, err := cfg.Queries.ListChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve the chirps from database", err)
		return
	}

	chirps := []Chirp{}
	for _, chirp := range raw_chirps {
		new_chirp := Chirp{
			chirp.ID,
			chirp.CreatedAt,
			chirp.UpdatedAt,
			chirp.Body,
			chirp.UserID,
		}
		chirps = append(chirps, new_chirp)
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
