package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/amir1226/go_mini_project/internal/auth"
	"github.com/amir1226/go_mini_project/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	if params.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Name is required")
		return
	}

	newUser, err := apiCfg.DB.CreateUser(
		r.Context(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      params.Name,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseUserToUser(newUser))
}

func (apiCfg *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request) {
	apikey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Auth error: %v", err))
		return
	}

	user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apikey)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error getting user: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
