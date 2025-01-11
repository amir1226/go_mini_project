package main

import (
	"fmt"
	"net/http"

	"github.com/amir1226/go_mini_project/internal/auth"
	"github.com/amir1226/go_mini_project/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfd *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apikey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := cfd.DB.GetUserByAPIKey(r.Context(), apikey)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error getting user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
