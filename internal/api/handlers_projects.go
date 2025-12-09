package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/AleksKAG/ar-backend/internal/store"
)

type createProjectReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CreateProject(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createProjectReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		// For MVP, we assume user ID is in context (from Auth middleware)
		uid := r.Context().Value(ctxUserIDKey{}).(string)

		p, err := db.CreateProject(r.Context(), uid, req.Name, req.Description)
		if err != nil {
			http.Error(w, "failed to create project", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
	}
}

func ListProjects(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value(ctxUserIDKey{}).(string)
		list, err := db.ListProjectsByOwner(r.Context(), uid)
		if err != nil {
			http.Error(w, "failed to list", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(list)
	}
}

func GetProject(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		p, err := db.GetProjectByID(r.Context(), id)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(p)
	}
}
