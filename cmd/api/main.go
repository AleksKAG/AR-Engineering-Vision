package main

import (
	"log"
	"net/http"

	"github.com/AleksKAG/ar-backend/internal/api"
	"github.com/AleksKAG/ar-backend/internal/s3"
	"github.com/AleksKAG/ar-backend/internal/store"
)

func main() {
	cfg := api.LoadConfig()
	
	// 1. Инициализация БД
	db, err := store.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 2. Инициализация MinIO (S3)
	s3c, err := s3.NewClient(cfg.S3Endpoint, cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3Bucket, false)
	if err != nil {
		log.Fatal("Failed to connect to MinIO:", err)
	}

	// 3. Запуск сервера
	router := api.NewRouter(db, s3c, cfg)
	log.Printf("🚀 AR Engineering Vision Backend listening on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
