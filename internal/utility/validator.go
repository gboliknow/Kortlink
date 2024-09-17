package utility

import "errors"

func ValidateUrlRequest(OriginalURL string) error {
	if OriginalURL == "" {
		return errors.New("original URL is required")
	}

	if !IsValidURL(OriginalURL) {
		return errors.New("invalid URL format")
	}

	return nil
}
