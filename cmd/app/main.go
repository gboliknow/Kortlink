package main

import (
	"kortlink/api"

	"kortlink/internal/database"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("Starting Kortlink project")
	dsn := os.Getenv("DB_URL")
	// connStr := fmt.Sprintf(dsn,
	// 	// config.Envs.DBUser,
	// 	// config.Envs.DBPassword,
	// 	// config.Envs.DBAddress,
	// 	// config.Envs.DBName,
	// )
	sqlStorage, err := database.NewPostgresStorage(dsn)
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
