package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/AleksKAG/ar-backend/internal/s3"
	"github.com/AleksKAG/ar-backend/internal/store"
)

func CreateIssue(db *store.DB, s3c *s3.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		elementID := vars["element_id"]

		r.ParseMultipartForm(10 << 20) // 10MB limit for photos
		file, handler, _ := r.FormFile("photo")
		defer file.Close()

		comment := r.FormValue("comment")
		deviationStr := r.FormValue("deviation_mm")
		deviation, _ := strconv.ParseFloat(deviationStr, 64)

		// 1. Загрузка фото в MinIO
		objectName := "issues/" + elementID + "/" + handler.Filename
		s3c.Upload(r.Context(), objectName, file, handler.Size, "image/jpeg")

		// 2. Формируем Presigned URL, чтобы фронтенд мог сразу показать фото
		photoURL, _ := s3c.URL(objectName, 24*7*time.Hour) // Ссылка на 7 дней

		// 3. Сохраняем замечание в БД
		issue, err := db.CreateIssue(r.Context(), elementID, photoURL, comment, deviation)
		if err != nil {
			http.Error(w, "failed to create issue", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(issue)
	}
}
