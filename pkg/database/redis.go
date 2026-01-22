package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/logger"
)

// RedisConfig holds Redis configuration
type RedisConfig struct {
	URL             string
	Host            string
	Password        string
	Port            int
	DB              int
	MaxRetries      int
	PoolSize        int
	MinIdleConns    int
	ConnMaxLifetime time.Duration
}

// NewRedisClient creates a new Redis client
func NewRedisClient(ctx context.Context, cfg RedisConfig, log *logger.Logger) (*redis.Client, error) {
	// Create Redis client
	var client *redis.Client

	if cfg.URL != "" {
		// Parse from URL
		opts, err := redis.ParseURL(cfg.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
		}
		client = redis.NewClient(opts)
	} else {
		// Build from individual fields
		client = redis.NewClient(&redis.Options{
			Addr:            fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Password:        cfg.Password,
			DB:              cfg.DB,
			MaxRetries:      cfg.MaxRetries,
			PoolSize:        cfg.PoolSize,
			MinIdleConns:    cfg.MinIdleConns,
			ConnMaxLifetime: cfg.ConnMaxLifetime,
		})
	}

	// Verify connection
	log.Info("Connecting to Redis...")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	log.Infof("Connected to Redis at %s:%d", cfg.Host, cfg.Port)

	return client, nil
}

// CloseRedisClient closes the Redis client gracefully
func CloseRedisClient(client *redis.Client, log *logger.Logger) {
	if client != nil {
		log.Info("Closing Redis client...")
		if err := client.Close(); err != nil {
			log.WithError(err).Logger().Error("Failed to close Redis client")
		} else {
			log.Info("Redis client closed")
		}
	}
}

// RedisHealthCheck checks if Redis is healthy
func RedisHealthCheck(ctx context.Context, client *redis.Client) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis health check failed: %w", err)
	}

	return nil
}

// GetRedisPoolStats returns Redis connection pool statistics
func GetRedisPoolStats(client *redis.Client) *redis.PoolStats {
	return client.PoolStats()
}
