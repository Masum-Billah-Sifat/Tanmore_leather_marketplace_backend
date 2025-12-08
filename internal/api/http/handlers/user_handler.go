// user_handler.go — Handles user-related HTTP endpoints

package handlers

import (
	"encoding/json"
	"net/http"

	service "tanmore_backend/internal/services"

	"github.com/rs/zerolog/log"
)

// Input body struct
type CreateUserRequest struct {
	GoogleID        string `json:"google_id"`
	PrimaryEmail    string `json:"primary_email"`
	DisplayName     string `json:"display_name"`
	ProfileImageURL string `json:"profile_image_url"`
}

// POST /users — creates a new user
func CreateUserHandler(svc *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		user, err := svc.CreateUser(r.Context(), req.GoogleID, req.PrimaryEmail, req.DisplayName, req.ProfileImageURL)
		if err != nil {
			log.Error().Err(err).Msg("failed to create user")
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}
