package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validateChirpHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}

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

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: replaceProfaneWords(params.Body),
	})
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
