package config

import "os"

type Config struct {
	Port             string
	DatabaseURL      string
	JWTSecret        string
	UploadDir        string
	TelegramBotToken string
	TelegramChatID   string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-me"
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	return Config{
		Port:             port,
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		JWTSecret:        jwtSecret,
		UploadDir:        uploadDir,
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
	}
}
