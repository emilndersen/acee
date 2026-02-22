package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	return Config{
		Port:        port,
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}
