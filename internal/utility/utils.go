package utility

import (
	"encoding/json"
	"kortlink/internal/models"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func WriteJSON(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := models.Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func WriteJSONGin(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"message": message,
		"data":    data,
	})
}

func IsValidURL(url string) bool {
	re := regexp.MustCompile(`^(http|https)://[^\s/$.?#].[^\s]*$`)
	return re.MatchString(url)
}

func GenerateShortURL() string {
	return uuid.New().String()[:8] // Example: Generate an 8-character short URL from a UUID
}
