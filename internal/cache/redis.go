package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisCache struct {
	Client *redis.Client
}

func NewRedisCache() *RedisCache {
	opt, _ := redis.ParseURL("rediss://default:AVPEAAIjcDEwMzE0ZGM5ZjkyNDE0MmI5YWQyOTdjZTFhNTFkNWYyNHAxMA@real-owl-21444.upstash.io:6379")
	client := redis.NewClient(opt)
	log.Info().
		Msg("Redis client initialized")

	return &RedisCache{
		Client: client,
	}
}

func (r *RedisCache) Get(key string) (string, error) {
	ctx := context.Background()
	val, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Warn().
			Str("key", key).
			Msg("Cache miss")
		return "", nil
	} else if err != nil {
		log.Error().
			Err(err).
			Str("key", key).
			Msg("Error retrieving key from Redis")
		return "", err
	}
	log.Info().
		Str("key", key).
		Msg("Cache hit")
	return val, nil
}

func (r *RedisCache) Set(key string, value string, expiration time.Duration) error {
	ctx := context.Background()
	err := r.Client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Error().
			Err(err).
			Str("key", key).
			Msg("Error setting key in Redis")
		return err
	}
	log.Info().
		Str("key", key).
		Dur("expiration", expiration).
		Msg("Cache set successfully")
	return nil
}

func (r *RedisCache) Delete(key string) error {
	ctx := context.Background()
	err := r.Client.Del(ctx, key).Err()
	if err != nil {
		log.Error().
			Err(err).
			Str("key", key).
			Msg("Error deleting key from Redis")
		return err
	}
	log.Info().
		Str("key", key).
		Msg("Cache deleted successfully")
	return nil
}

func (r *RedisCache) CacheStats(key string, value string, expiration time.Duration) error {
	ctx := context.Background()
	err := r.Client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Error().
			Err(err).
			Str("key", key).
			Msg("Error caching stats in Redis")
		return err
	}
	log.Info().
		Str("key", key).
		Dur("expiration", expiration).
		Msg("Stats cached successfully")
	return nil
}

func (r *RedisCache) GetCachedStats(key string) (string, error) {
	ctx := context.Background()
	val, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Warn().
			Str("key", key).
			Msg("Cache miss for stats")
		return "", nil // Key does not exist
	} else if err != nil {
		log.Error().
			Err(err).
			Str("key", key).
			Msg("Error retrieving stats key from Redis")
		return "", err
	}
	log.Info().
		Str("key", key).
		Msg("Cache hit for stats")
	return val, nil
}
