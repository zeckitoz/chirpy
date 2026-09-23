package main

import (
	"errors"
	"net/http"
)

func (cfg *apiConfig) resetServer(w http.ResponseWriter, req *http.Request) {

	if cfg.Platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden", errors.New("Not dev environment"))
		return
	}

	err := cfg.Queries.ResetUsers(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete all users from the database.", err)
		return
	}

	cfg.FileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server has been reset."))
}
