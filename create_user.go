package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zeckitoz/chirpy/internal/auth"
	"github.com/zeckitoz/chirpy/internal/database"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type User struct {
		ID             uuid.UUID `json:"id"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
		Email          string    `json:"email"`
		HashedPassword string    `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user_params := database.CreateUserParams{}
	user_params.Email = params.Email
	user_params.HashedPassword = hash

	user, err := cfg.Queries.CreateUser(req.Context(), user_params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create new user", err)
		return
	}

	new_user := User{user.ID, user.CreatedAt, user.UpdatedAt, user.Email, "****"}
	respondWithJSON(w, http.StatusCreated, new_user)
}
