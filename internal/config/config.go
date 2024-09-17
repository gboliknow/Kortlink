package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port       string
	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string
}

var Envs = InitializeConfig()

func InitializeConfig() Config {

	err := godotenv.Load()
	if err != nil {
		log.Error().Err(err).Msg("Error loading .env file")
	}

	return Config{
		Port:       getEnv("PORT", "5432"),
		DBUser:     getEnv("DB_USER", "user_3"),
		DBPassword: getEnv("DB_PASSWORD", "pass_3"),
		DBName:     getEnv("DB_NAME", "kortlink"),
		DBAddress:  fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "5432")),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
