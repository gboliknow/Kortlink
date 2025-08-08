// internal/tests/shortlink_handler_test.go
package tests

import (
	"bytes"
	"encoding/json"
	"kortlink/api"
	"kortlink/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleCreateShortlink(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockStore := new(MockStore)
	mockCache := new(MockCache)

	service := api.NewShortlinkService(mockStore, mockCache)

	// Expectation: store.CreateShortURL should be called and succeed
	mockStore.On("CreateShortURL", mock.AnythingOfType("*models.ShortURL")).Return(nil)
	mockCache.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("string"), 24*time.Hour).Return(nil)

	router := gin.Default()
	apiV1 := router.Group("/api/v1")
	service.ShortlinkRoutes(apiV1)

	payload := models.ShortURL{
		OriginalURL: "https://google.com",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/shortlink", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	mockStore.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}
