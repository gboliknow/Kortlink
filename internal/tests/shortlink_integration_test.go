package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"kortlink/api"
	"kortlink/internal/cache"
	"kortlink/internal/config"

	"kortlink/internal/database"
	"kortlink/internal/models"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)



func TestCreateAndRetrieveShortlink(t *testing.T) {
	// Directly use the actual database connection string

	dsn := "postgresql://%s:%s@%s/%s?sslmode=disable"
	connStr := fmt.Sprintf(dsn,
		config.Envs.DBUser,
		config.Envs.DBPassword,
		config.Envs.DBAddress,
		config.Envs.DBName,
	)
	sqlStorage, err := database.NewPostgresStorage(connStr)
	assert.NoError(t, err)
	redisCache := cache.NewRedisCache()

	// Setup Gin and API server
	store := api.NewStore(sqlStorage.Pool())
	shortlinkService := api.NewShortlinkService(store, redisCache)
	router := gin.Default()
	apiV1 := router.Group("/api/v1")
	shortlinkService.ShortlinkRoutes(apiV1)

	// 1. Test the creation of a shortlink
	shortURLPayload := models.ShortURLPayload{
		OriginalURL: "http://example.com",
	}
	payloadBytes, err := json.Marshal(shortURLPayload)
	assert.NoError(t, err) // Ensure marshalling was successful

	req, err := http.NewRequest("POST", "/api/v1/shortlink", bytes.NewBuffer(payloadBytes))

	log.Printf("Response Body: %s", payloadBytes) // Log the response body

	assert.NoError(t, err) // Check if the request was created successfully
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the status code
	assert.Equal(t, http.StatusCreated, w.Code)
	if w.Code != http.StatusCreated {
		log.Printf("Response Body: %s", w.Body.String()) // Log the response body
	}

	// Unmarshal response to get the created short URL
	var createdShortURL models.ShortURL
	err = json.Unmarshal(w.Body.Bytes(), &createdShortURL)
	assert.NoError(t, err)

	// Verify that the original URL matches what was sent
	assert.Equal(t, shortURLPayload.OriginalURL, createdShortURL.OriginalURL)
}
