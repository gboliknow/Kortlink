package api

import (
	"kortlink/internal/cache"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type APIServer struct {
	addr   string
	store  Store
	logger zerolog.Logger
	cache  *cache.RedisCache
}

func NewAPIServer(addr string, store Store) *APIServer {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	pass := os.Getenv("Password")
	redisCache := cache.NewRedisCache("localhost:6379", pass, 0)
	return &APIServer{addr: addr, store: store, logger: logger, cache: redisCache}
}

func (s *APIServer) Serve() {
	router := gin.Default()
	apiV1 := router.Group("/api/v1")

	//registering the routes
	shortlinkService := NewShortlinkService(s.store, s.cache)
	shortlinkService.ShortlinkRoutes(apiV1)

	s.logger.Info().Str("addr", s.addr).Msg("Starting API server")
	if err := http.ListenAndServe(s.addr, router); err != nil {
		s.logger.Fatal().Err(err).Msg("Server stopped")
	}
}
