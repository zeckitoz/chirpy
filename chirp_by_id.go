package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getChipByIDHandler(w http.ResponseWriter, req *http.Request) {
	path_value := req.PathValue("chirpID")
	id, err := uuid.Parse(path_value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert path value", err)
		return
	}

	raw_chirp, err := cfg.Queries.GetChirpByID(req.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	chirp := Chirp{
		raw_chirp.ID,
		raw_chirp.CreatedAt,
		raw_chirp.UpdatedAt,
		raw_chirp.Body,
		raw_chirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, chirp)
}
