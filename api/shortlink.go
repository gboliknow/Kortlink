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

// @Summary      Create Shortlink
// @Description  Create a new short URL
// @Tags         shortlinks
// @Accept       json
// @Produce      json
// @Param        body  body   models.ShortURLPayload  true  "Original URL payload"
// @Success      201   {object} models.ShortURL
// @Failure      400   {object} models.Response
// @Failure      500   {object} models.Response
// @Router       /api/v1/shortlink [post]
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

// @Summary      Redirect to the original URL
// @Description  Redirects to the original URL based on the provided short URL
// @Tags         shortlinks
// @Param        shortURL   path      string  true  "Short URL"
// @Success      302        {string}  string  "Redirected to the original URL"
// @Failure      400        {string}  string  "Short URL is required"
// @Failure      404        {string}  string  "Short URL not found"
// @Failure      500        {string}  string  "Failed to update access count"
// @Router       /api/v1/{shortURL} [get]
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

// @Summary      Update a short URL
// @Description  Update the original URL for a given short URL
// @Tags         shortlinks
// @Accept       json
// @Produce      json
// @Param        shortURL   path      string      true  "Short URL"
// @Param        body       body      models.ShortURL  true  "New original URL"
// @Success      200        {string}  string      "Short URL updated successfully"
// @Failure      400        {string}  string      "Invalid request payload or Short URL is required"
// @Failure      404        {string}  string      "Short URL not found"
// @Failure      500        {string}  string      "Failed to update short URL"
// @Router       /api/v1/{shortURL} [put]
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

// @Summary      Get short URL statistics
// @Description  Fetches the statistics (e.g., access count) for a given short URL
// @Tags         shortlinks
// @Param        shortURL   path      string  true  "Short URL"
// @Success      200        {object}  map[string]interface{}  "Statistics fetched successfully"
// @Failure      400        {string}  string  "Short URL is required"
// @Failure      404        {string}  string  "Short URL not found"
// @Router       /api/v1/{shortURL}/stats [get]
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

// @Summary      Delete a short URL
// @Description  Deletes a given short URL and its related data
// @Tags         shortlinks
// @Param        shortURL   path      string  true  "Short URL"
// @Success      200        {string}  string  "Short URL deleted successfully"
// @Failure      400        {string}  string  "Short URL is required"
// @Failure      404        {string}  string  "Short URL not found"
// @Failure      500        {string}  string  "Failed to delete short URL"
// @Router      /api/v1/{shortURL} [delete]
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

// @Summary      Get all short URLs
// @Description  Fetches a list of all short URLs stored in the system
// @Tags         shortlinks
// @Success      200        {array}   models.ShortURL  "Successfully fetched URLs"
// @Failure      500        {string}  string  "Failed to fetch URLs"
// @Router       /api/v1/shortlinks [get]
func (s *ShortlinkService) handleGetAllShortlinks(c *gin.Context) {
	urls, err := s.store.GetAllShortURLs()
	if err != nil {
		utility.WriteJSON(c.Writer, http.StatusInternalServerError, "Failed to fetch URLs", nil)
		return
	}

	utility.WriteJSON(c.Writer, http.StatusOK, "Successfully fetched URLs", urls)
}
