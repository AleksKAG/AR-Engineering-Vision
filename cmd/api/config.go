package api

import "os"

type Config struct {
    Port         string
    DatabaseURL  string
    JWTSecret    string
    S3Endpoint   string
    S3Bucket     string
    S3AccessKey  string
    S3SecretKey  string
}

func LoadConfig() *Config {
    return &Config{
        Port:         getEnv("PORT", "8080"),
        DatabaseURL:  getEnv("DATABASE_URL", "postgres://postgres:pass@postgres:5432/ar?sslmode=disable"),
        JWTSecret:    getEnv("JWT_SECRET", "supersecret"),
        S3Endpoint:   getEnv("S3_ENDPOINT", "minio:9000"),
        S3Bucket:     getEnv("S3_BUCKET", "ar-models"),
        S3AccessKey:  getEnv("S3_ACCESS_KEY", "minioadmin"),
        S3SecretKey:  getEnv("S3_SECRET_KEY", "minioadmin"),
    }
}

func getEnv(k, def string) string {
    if v := os.Getenv(k); v != "" {
        return v
    }
    return def
}
