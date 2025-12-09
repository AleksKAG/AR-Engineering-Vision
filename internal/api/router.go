package api

import (
    "net/http"

    "github.com/gorilla/mux"
    "github.com/AleksKAG/ar-backend/internal/store"
)

func NewRouter(db *store.DB, cfg *Config) http.Handler {
    r := mux.NewRouter()

    // Auth
    r.HandleFunc("/api/v1/auth/login", Login(db)).Methods("POST")

    // Projects
    r.HandleFunc("/api/v1/projects", Auth(CreateProject(db))).Methods("POST")
    r.HandleFunc("/api/v1/projects", Auth(ListProjects(db))).Methods("GET")

    // Health
    r.HandleFunc("/api/v1/admin/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok"))
    })

    return r
}
