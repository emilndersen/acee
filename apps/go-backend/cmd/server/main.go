package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/emilndersen/acee/apps/go-backend/internal/config"
	"github.com/emilndersen/acee/apps/go-backend/internal/db"
	httpapi "github.com/emilndersen/acee/apps/go-backend/internal/http"
)

func main() {
	// Игнорируем ошибку, так как в проде .env может не быть (переменные из окружения)
	_ = godotenv.Load()

	cfg := config.Load()

	// 1. Инициализация БД
	pool, err := db.NewPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	defer pool.Close()

	// 2. Инициализация Роутера
	router := httpapi.NewRouter(pool, cfg)

	// 3. Настройка HTTP сервера с таймаутами
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 4. Канал для перехвата сигналов остановки (Ctrl+C, SIGTERM)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 Server started on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Listen error: %v", err)
		}
	}()

	// Ждем сигнала остановки
	<-stop
	log.Println("⚠️ Shutting down server...")

	// 5. Даем серверу 10 секунд на завершение текущих дел
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}
