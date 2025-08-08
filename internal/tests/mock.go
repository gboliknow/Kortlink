

// internal/tests/mocks.go
package tests

import (
	"kortlink/internal/models"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) CreateShortURL(shortURL *models.ShortURL) error {
	args := m.Called(shortURL)
	return args.Error(0)
}

func (m *MockStore) GetOriginalURL(shortURL string) (string, error) {
	args := m.Called(shortURL)
	return args.String(0), args.Error(1)
}

func (m *MockStore) IncrementAccessCount(shortURL string) error {
	args := m.Called(shortURL)
	return args.Error(0)
}

func (m *MockStore) UpdateShortURL(shortURL, newOriginalURL string) error {
	args := m.Called(shortURL, newOriginalURL)
	return args.Error(0)
}

func (m *MockStore) DeleteShortURL(shortURL string) error {
	args := m.Called(shortURL)
	return args.Error(0)
}

func (m *MockStore) GetShortURLStats(shortURL string) (*models.ShortURL, error) {
	args := m.Called(shortURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ShortURL), args.Error(1)
}

func (m *MockStore) GetAllShortURLs() ([]models.ShortURL, error) {
	args := m.Called()
	return args.Get(0).([]models.ShortURL), args.Error(1)
}

// ---- Mock RedisCache ----
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Set(key, value string, duration time.Duration) error {
	args := m.Called(key, value, duration)
	return args.Error(0)
}

func (m *MockCache) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockCache) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}
