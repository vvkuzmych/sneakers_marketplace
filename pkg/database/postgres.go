package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/logger"
)

// PostgresConfig holds PostgreSQL configuration
type PostgresConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// NewPostgresPool creates a new PostgreSQL connection pool
func NewPostgresPool(ctx context.Context, cfg PostgresConfig, log *logger.Logger) (*pgxpool.Pool, error) {
	// Parse config
	config, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Set connection pool settings
	if cfg.MaxOpenConns > 0 {
		config.MaxConns = int32(cfg.MaxOpenConns)
	} else {
		config.MaxConns = 25 // Default
	}

	if cfg.ConnMaxLifetime > 0 {
		config.MaxConnLifetime = cfg.ConnMaxLifetime
	} else {
		config.MaxConnLifetime = time.Hour // Default
	}

	if cfg.ConnMaxIdleTime > 0 {
		config.MaxConnIdleTime = cfg.ConnMaxIdleTime
	} else {
		config.MaxConnIdleTime = 30 * time.Minute // Default
	}

	// Create connection pool
	log.Info("Connecting to PostgreSQL...")
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Infof("Connected to PostgreSQL (max connections: %d)", config.MaxConns)

	return pool, nil
}

// ClosePostgresPool closes the PostgreSQL connection pool gracefully
func ClosePostgresPool(pool *pgxpool.Pool, log *logger.Logger) {
	if pool != nil {
		log.Info("Closing PostgreSQL connection pool...")
		pool.Close()
		log.Info("PostgreSQL connection pool closed")
	}
}

// HealthCheck checks if database is healthy
func HealthCheck(ctx context.Context, pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// GetPoolStats returns connection pool statistics
func GetPoolStats(pool *pgxpool.Pool) *PoolStats {
	stats := pool.Stat()
	return &PoolStats{
		AcquireCount:         stats.AcquireCount(),
		AcquireDuration:      stats.AcquireDuration(),
		AcquiredConns:        stats.AcquiredConns(),
		CanceledAcquireCount: stats.CanceledAcquireCount(),
		ConstructingConns:    stats.ConstructingConns(),
		EmptyAcquireCount:    stats.EmptyAcquireCount(),
		IdleConns:            stats.IdleConns(),
		MaxConns:             stats.MaxConns(),
		TotalConns:           stats.TotalConns(),
	}
}

// PoolStats holds connection pool statistics
type PoolStats struct {
	AcquireCount         int64
	AcquireDuration      time.Duration
	CanceledAcquireCount int64
	EmptyAcquireCount    int64
	AcquiredConns        int32
	ConstructingConns    int32
	IdleConns            int32
	MaxConns             int32
	TotalConns           int32
}

// String returns a string representation of pool stats
func (s *PoolStats) String() string {
	return fmt.Sprintf(
		"acquired=%d idle=%d total=%d max=%d acquiring=%d",
		s.AcquiredConns,
		s.IdleConns,
		s.TotalConns,
		s.MaxConns,
		s.ConstructingConns,
	)
}
