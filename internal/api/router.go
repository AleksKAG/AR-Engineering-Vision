package api

import (
	"net/http"
	"time" // Добавьте этот импорт

	"github.com/gorilla/mux"
	"github.com/AleksKAG/ar-backend/internal/s3"
	"github.com/AleksKAG/ar-backend/internal/store"
)

// Добавляем s3c в аргументы
func NewRouter(db *store.DB, s3c *s3.Client, cfg *Config) http.Handler {
	r := mux.NewRouter()

	// --- AUTH ---
	r.HandleFunc("/api/v1/auth/login", Login(db)).Methods("POST")

	// --- PROJECTS ---
	r.HandleFunc("/api/v1/projects", Auth(CreateProject(db))).Methods("POST")
	r.HandleFunc("/api/v1/projects", Auth(ListProjects(db))).Methods("GET")

	// --- BIM MODELS (Загрузка IFC/glTF) ---
	r.HandleFunc("/api/v1/projects/{project_id}/models", Auth(UploadModel(db, s3c))).Methods("POST")

	// --- AR DATA (Для AR-движка) ---
	r.HandleFunc("/api/v1/rooms/{room_id}/elements", Auth(GetElementsForAR(db))).Methods("GET")

	// --- ISSUES (Фотофиксация и замечания) ---
	r.HandleFunc("/api/v1/elements/{element_id}/issues", Auth(CreateIssue(db, s3c))).Methods("POST")
	
// Внутренний API для Python-конвертера (в продакшене стоит закрыть middleware'ем или сетью)
r.HandleFunc("/api/v1/internal/bim-processed", ProcessBimCallback(db)).Methods("POST")
	
	// --- HEALTH ---
	r.HandleFunc("/api/v1/admin/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	return r
}
