package api

import (
	"encoding/json"
	"net/http"

	"github.com/AleksKAG/ar-backend/internal/store"
)

// BimProcessedRequest - структура, которую присылает Python-скрипт
type BimProcessedRequest struct {
	ProjectID string               `json:"project_id"`
	GlbS3Key  string               `json:"glb_s3_key"`
	Rooms     []store.Room         `json:"rooms"`
	Elements  []store.ElementInput `json:"elements"`
}

// ProcessBimCallback принимает результат парсинга от Python-сервиса
func ProcessBimCallback(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BimProcessedRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		// Транзакция для атомарного сохранения комнат и элементов
		err := db.SaveBimData(r.Context(), req.ProjectID, req.GlbS3Key, req.Rooms, req.Elements)
		if err != nil {
			http.Error(w, "failed to save bim data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}
