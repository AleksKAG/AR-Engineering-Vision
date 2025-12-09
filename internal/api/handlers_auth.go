package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AleksKAG/ar-backend/internal/store"
	"github.com/AleksKAG/ar-backend/internal/auth"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func Login(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		user, err := db.GetUserByEmail(r.Context(), req.Email)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		// create token
		token, err := auth.GenerateJWT(user.ID.String(), time.Hour*24)
		if err != nil {
			http.Error(w, "failed to create token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(loginResponse{Token: token})
	}
}
