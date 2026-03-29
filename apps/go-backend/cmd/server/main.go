package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/emilndersen/acee/apps/go-backend/internal/config"
	"github.com/emilndersen/acee/apps/go-backend/internal/db"
	httpapi "github.com/emilndersen/acee/apps/go-backend/internal/http"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	pool, err := db.NewPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	router := httpapi.NewRouter(pool)

	log.Println("Go backend running on :" + cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
