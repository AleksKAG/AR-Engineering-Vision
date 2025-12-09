package main

import (
    "log"
    "net/http"
    "os"

    "github.com/AleksKAG/ar-backend/internal/api"
    "github.com/AleksKAG/ar-backend/internal/store"
)

func main() {
    cfg := api.LoadConfig()

    db, err := store.NewPostgres(cfg.DatabaseURL)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    router := api.NewRouter(db, cfg)

    log.Printf("Server listening on :%s", cfg.Port)
    log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
