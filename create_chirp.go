package main

import (
	"encoding/json"
	"net/http"
	"time"

	"strings"

	"github.com/google/uuid"
	"github.com/zeckitoz/chirpy/internal/database"
)

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, req *http.Request) {
	type parameter struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameter{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}

	cleaned_body := replaceProfaneWords(params.Body)

	chirp_params := database.CreateChirpParams{}
	chirp_params.Body = cleaned_body
	chirp_params.UserID = params.UserID

	chirp, err := cfg.Queries.CreateChirp(req.Context(), chirp_params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create new chirp", err)
		return
	}

	new_chirp := Chirp{
		chirp.ID,
		chirp.CreatedAt,
		chirp.UpdatedAt,
		chirp.Body,
		chirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, new_chirp)
}

func replaceProfaneWords(text string) string {
	profane_words := [3]string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(text, " ")

	for i, word := range words {
		for _, p_word := range profane_words {
			if strings.ToLower(word) == p_word {
				words[i] = "****"
			}
		}
	}

	new_text := strings.Join(words, " ")
	return new_text
}
