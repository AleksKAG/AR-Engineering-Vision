package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/AleksKAG/ar-backend/internal/s3"
	"github.com/AleksKAG/ar-backend/internal/store"
)

// UploadModel принимает BIM-модель (IFC/glTF) и сохраняет в MinIO + Postgres
func UploadModel(db *store.DB, s3c *s3.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		projectID := vars["project_id"]

		// Лимит 100MB для тяжелых BIM-моделей
		if err := r.ParseMultipartForm(100 << 20); err != nil {
			http.Error(w, "file too large", http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "invalid file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 1. Загружаем файл в MinIO
		objectName := "projects/" + projectID + "/models/" + handler.Filename
		_, err = s3c.Upload(r.Context(), objectName, file, handler.Size, handler.Header.Get("Content-Type"))
		if err != nil {
			http.Error(w, "S3 upload failed", http.StatusInternalServerError)
			return
		}

		// 2. Сохраняем метаданные в Postgres
		model, err := db.CreateModel(r.Context(), projectID, handler.Filename, handler.Header.Get("Content-Type"), objectName)
		if err != nil {
			http.Error(w, "DB save failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(model)
	}
}
