package main

import (
	"fmt"
	"kortlink/api"

	"kortlink/internal/config"
	"kortlink/internal/database"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	env := os.Getenv("ENVIRONMENT")
	var connStr string
	if env == "debug" {
		log.Info().Msg("Starting Debug Kortlink project")
		dsn := "postgresql://%s:%s@%s/%s?sslmode=disable"
		connStr = fmt.Sprintf(dsn,
			config.Envs.DBUser,
			config.Envs.DBPassword,
			config.Envs.DBAddress,
			config.Envs.DBName,
		)
	} else {
		log.Info().Msg("Starting Production Kortlink project")
		connStr = os.Getenv("DB_URL")
	}

	// Attempt to connect to the database
	sqlStorage, err := database.NewPostgresStorage(connStr)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Initialize the database (e.g., create tables, seed data)
	if err := sqlStorage.InitializeDatabase(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}

	store := api.NewStore(sqlStorage.Pool())
	apiServer := api.NewAPIServer(":8080", store)
	log.Info().Msg("Starting API server on port 8080")
	apiServer.Serve()
}
