package api

import (
	"kortlink/internal/models"
	"kortlink/internal/utility"
	"net/http"
	"time"

	"kortlink/internal/cache"

	"github.com/gin-gonic/gin"
)

type ShortlinkService struct {
	store Store
	cache *cache.RedisCache
}

func NewShortlinkService(s Store, c *cache.RedisCache) *ShortlinkService {
	return &ShortlinkService{store: s, cache: c}
}

func (s *ShortlinkService) ShortlinkRoutes(r *gin.RouterGroup) {
	r.POST("/shortlink", s.handleCreateShortlink)
	r.GET("/:shortURL", s.handleRedirect)
	r.PUT("/:shortURL", s.handleUpdateShortlink)
	r.DELETE("/:shortURL", s.handleDeleteShortlink)
	r.GET("/:shortURL/stats", s.handleGetStats)
	r.GET("/shortlinks", s.handleGetAllShortlinks)
	r.GET("/debug/healthCheck", s.handleHealthCheck)
}

func (s *ShortlinkService) handleHealthCheck(c *gin.Context) {
	utility.WriteJSON(c.Writer, http.StatusOK, "Health check successfully", nil)
}

func (s *ShortlinkService) handleCreateShortlink(c *gin.Context) {
	var payload models.ShortURL
	if err := c.ShouldBindJSON(&payload); err != nil {
		utility.WriteJSON(c.Writer, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}
	if err := utility.ValidateUrlRequest(payload.OriginalURL); err != nil {
		utility.WriteJSON(c.Writer, http.StatusBadRequest, err.Error(), nil)
		return
	}

	shortURL := utility.GenerateShortURL()
	now := time.Now()
	shortLink := &models.ShortURL{
		OriginalURL: payload.OriginalURL,
		ShortURL:    shortURL,
		AccessCount: 0,
		CreatedAt:   now,
	}

	err := s.store.CreateShortURL(shortLink)
	if err != nil {
		utility.WriteJSON(c.Writer, http.StatusInternalServerError, "Failed to create short link", nil)
		return
	}
	_ = s.cache.Set(shortLink.ShortURL, shortLink.OriginalURL, 24*time.Hour)
	utility.WriteJSON(c.Writer, http.StatusCreated, "Short link created successfully", shortLink)
}

func (s *ShortlinkService) handleRedirect(c *gin.Context) {
	shortURL := c.Param("shortURL")
	if shortURL == "" {
		utility.WriteJSON(c.Writer, http.StatusBadRequest, "Short URL is required", nil)
		return
	}

	originalURL, err := s.cache.Get(shortURL)
	if err == nil && originalURL != "" {
		err = s.store.IncrementAccessCount(shortURL)
		if err != nil {
			utility.WriteJSON(c.Writer, http.StatusInternalServerError, "Failed to update access count", nil)
			return
		}
		c.Redirect(http.StatusFound, originalURL)
		return
	}

	url, err := s.store.GetOriginalURL(shortURL)
	if err != nil {
		utility.WriteJSON(c.Writer, http.StatusNotFound, "Short URL not found", nil)
		return
	}

	_ = s.cache.Set(shortURL, originalURL, 24*time.Hour)
	err = s.store.IncrementAccessCount(shortURL)
	if err != nil {
		utility.WriteJSON(c.Writer, http.StatusInternalServerError, "Failed to update access count", nil)
		return
	}

	c.Redirect(http.StatusFound, url)
}

func (s *ShortlinkService) handleUpdateShortlink(c *gin.Context) {
	shortURL := c.Param("shortURL")
	if shortURL == "" {
		utility.WriteJSON(c.Writer, http.StatusBadRequest, "Short URL is required", nil)
		return
	}
	_, err := s.store.GetOriginalURL(shortURL)
	if err != nil {
		utility.WriteJSON(c.Writer, http.StatusNotFound, "Short URL not found", nil)
		return
	}

	var payload models.ShortURL
	if err := c.ShouldBindJSON(&payload); err != nil {
		utility.WriteJSON(c.Writer, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	if err := utility.ValidateUrlRequest(payload.OriginalURL); err != nil {
		utility.WriteJSON(c.Writer, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := s.store.UpdateShortURL(shortURL, payload.OriginalURL); err != nil {
		utility.WriteJSON(c.Writer, http.StatusInternalServerError, "Failed to update short URL", nil)
		return
	}
	_ = s.cache.Set(shortURL, payload.OriginalURL, 24*time.Hour)
	utility.WriteJSON(c.Writer, http.StatusOK, "Short URL updated successfully", nil)
}

func (s *ShortlinkService) handleGetStats(c *gin.Context) {
	shortURL := c.Param("shortURL")
	if shortURL == "" {
		utility.WriteJSON(c.Writer, http.StatusBadRequest, "Short URL is required", nil)
		return
	}

	stats, err := s.store.GetShortURLStats(shortURL)
	if err != nil {
		utility.WriteJSON(c.Writer, http.StatusNotFound, "Short URL not found", nil)
		return
	}

	utility.WriteJSON(c.Writer, http.StatusOK, "Statistics fetched successfully", stats)
}

func (s *ShortlinkService) handleDeleteShortlink(c *gin.Context) {
	shortURL := c.Param("shortURL")
	if shortURL == "" {
		utility.WriteJSON(c.Writer, http.StatusBadRequest, "Short URL is required", nil)
		return
	}
	_, err := s.store.GetOriginalURL(shortURL)
	if err != nil {
		utility.WriteJSON(c.Writer, http.StatusNotFound, "Short URL not found", nil)
		return
	}

	if err := s.store.DeleteShortURL(shortURL); err != nil {
		utility.WriteJSON(c.Writer, http.StatusInternalServerError, "Failed to delete short URL", nil)
		return
	}
	_ = s.cache.Delete(shortURL)
	utility.WriteJSON(c.Writer, http.StatusOK, "Short URL deleted successfully", nil)
}

func (s *ShortlinkService) handleGetAllShortlinks(c *gin.Context) {
	urls, err := s.store.GetAllShortURLs()
	if err != nil {
		utility.WriteJSON(c.Writer, http.StatusInternalServerError, "Failed to fetch URLs", nil)
		return
	}

	utility.WriteJSON(c.Writer, http.StatusOK, "Successfully fetched URLs", urls)
}
