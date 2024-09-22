package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	// Create the PostgreSQL connection pool
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Error().Err(err).Msg("Unable to connect to database")
		return nil, err
	}

	log.Info().Msg("Connected to PostgreSQL!")
	return &PostgresStorage{pool: pool}, nil
}

// Pool returns the underlying connection pool
func (s *PostgresStorage) Pool() *pgxpool.Pool {
	return s.pool
}

func (s *PostgresStorage) InitializeDatabase() error {
	// Create the users table
	if err := s.createUrlsTable(); err != nil {
		log.Error().Err(err).Msg("Failed to create urls table")
		return err
	}
	log.Info().Msg("url table created successfully")

	return nil
}

func (s *PostgresStorage) createUrlsTable() error {
	sql := `
    CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		original_url TEXT NOT NULL,
		short_url TEXT NOT NULL UNIQUE,
		access_count INT DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
    `
	_, err := s.pool.Exec(context.Background(), sql)
	return err
}
