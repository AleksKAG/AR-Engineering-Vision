package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/AleksKAG/ar-backend/internal/s3"
	"github.com/AleksKAG/ar-backend/internal/store"
)

type aiAnalysisRequest struct {
	ImageURL     string `json:"image_url"`
	ExpectedType string `json:"expected_type"`
}

type aiAnalysisResponse struct {
	AIDetectedType string  `json:"ai_detected_type"`
	AIConfidence   float64 `json:"ai_confidence"`
	IsMatch        bool    `json:"is_match"`
	Message        string  `json:"message"`
}

// AnalyzeAndSaveIssue принимает фото, отправляет в AI, сохраняет результат
func AnalyzeAndSaveIssue(db *store.DB, s3c *s3.Client, aiServiceURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		elementID := vars["element_id"]

		// 1. Парсим multipart форму (фото + метаданные)
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			http.Error(w, "invalid form data", http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("photo")
		if err != nil {
			http.Error(w, "photo file is required", http.StatusBadRequest)
			return
		}
		defer file.Close()

		comment := r.FormValue("comment")
		deviationStr := r.FormValue("deviation_mm")
		deviation, _ := strconv.ParseFloat(deviationStr, 64)
		expectedType := r.FormValue("expected_type") // например, "duct"

		// 2. Загружаем фото в MinIO
		objectName := fmt.Sprintf("issues/%s/%d_%s", elementID, time.Now().Unix(), handler.Filename)
		_, err = s3c.Upload(r.Context(), objectName, file, handler.Size, handler.Header.Get("Content-Type"))
		if err != nil {
			http.Error(w, "S3 upload failed", http.StatusInternalServerError)
			return
		}

		// 3. Получаем Presigned URL для AI-сервиса
		photoURL, err := s3c.URL(objectName, time.Hour)
		if err != nil {
			http.Error(w, "Failed to generate URL", http.StatusInternalServerError)
			return
		}

		// 4. Отправляем запрос в Python AI-сервис
		aiReq := aiAnalysisRequest{
			ImageURL:     photoURL,
			ExpectedType: expectedType,
		}
		aiReqBody, _ := json.Marshal(aiReq)
		
		resp, err := http.Post(aiServiceURL+"/api/v1/ai/analyze", "application/json", bytes.NewBuffer(aiReqBody))
		if err != nil {
			// Если AI упал, всё равно сохраняем фото, но без AI-данных
			issue, _ := db.CreateIssue(r.Context(), elementID, photoURL, comment, deviation, "", 0.0, false)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(issue)
			return
		}
		defer resp.Body.Close()

		var aiResp aiAnalysisResponse
		json.NewDecoder(resp.Body).Decode(&aiResp)

		// 5. Сохраняем замечание с результатами AI в БД
		issue, err := db.CreateIssueWithAI(r.Context(), elementID, photoURL, comment, deviation, aiResp.AIDetectedType, aiResp.AIConfidence, aiResp.IsMatch)
		if err != nil {
			http.Error(w, "failed to save issue", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(issue)
	}
}
