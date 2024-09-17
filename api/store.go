package api

import (
	"context"
	"fmt"
	"kortlink/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	CreateShortURL(shortURL *models.ShortURL) error
	GetOriginalURL(shortURL string) (string, error)
	IncrementAccessCount(shortURL string) error
	UpdateShortURL(shortURL string, newOriginalURL string) error
	DeleteShortURL(shortURL string) error
	GetShortURLStats(shortURL string) (*models.ShortURL, error)
	GetAllShortURLs() ([]models.ShortURL, error)
}

type Storage struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Storage {
	return &Storage{
		pool: pool,
	}
}

func (s *Storage) CreateShortURL(shortURL *models.ShortURL) error {
	query := `
		INSERT INTO urls (original_url, short_url, access_count, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	err := s.pool.QueryRow(context.Background(), query,
		shortURL.OriginalURL,
		shortURL.ShortURL,
		shortURL.AccessCount,
		shortURL.CreatedAt,
	).Scan(&shortURL.ID)

	if err != nil {
		return fmt.Errorf("could not insert short URL: %w", err)
	}

	return nil
}
func (s *Storage) GetOriginalURL(shortURL string) (string, error) {
	var originalURL string
	query := `SELECT original_url FROM urls WHERE short_url = $1`
	err := s.pool.QueryRow(context.Background(), query, shortURL).Scan(&originalURL)
	if err != nil {
		return "", err
	}
	return originalURL, nil
}
func (s *Storage) IncrementAccessCount(shortURL string) error {
	query := `UPDATE urls SET access_count = access_count + 1, updated_at = NOW() WHERE short_url = $1`
	_, err := s.pool.Exec(context.Background(), query, shortURL)
	return err
}
func (s *Storage) UpdateShortURL(shortURL string, newOriginalURL string) error {
	query := `
		UPDATE urls
		SET original_url = $1, updated_at = NOW()
		WHERE short_url = $2
	`
	_, err := s.pool.Exec(context.Background(), query, newOriginalURL, shortURL)
	return err
}
func (s *Storage) DeleteShortURL(shortURL string) error {
	query := `DELETE FROM urls WHERE short_url = $1`
	_, err := s.pool.Exec(context.Background(), query, shortURL)
	return err
}
func (s *Storage) GetShortURLStats(shortURL string) (*models.ShortURL, error) {
	var url models.ShortURL
	query := `SELECT original_url, short_url, access_count, created_at, updated_at FROM urls WHERE short_url = $1`
	err := s.pool.QueryRow(context.Background(), query, shortURL).Scan(&url.OriginalURL, &url.ShortURL, &url.AccessCount, &url.CreatedAt, &url.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &url, nil
}
func (s *Storage) GetAllShortURLs() ([]models.ShortURL, error) {
	query := `
		SELECT short_url, original_url, access_count, created_at, updated_at
		FROM urls
	`
	rows, err := s.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []models.ShortURL
	for rows.Next() {
		var url models.ShortURL
		if err := rows.Scan(&url.ShortURL, &url.OriginalURL, &url.AccessCount, &url.CreatedAt, &url.UpdatedAt); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}
