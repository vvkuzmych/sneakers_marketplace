package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	Env string

	// Database
	Database DatabaseConfig

	// Redis
	Redis RedisConfig

	// Kafka
	Kafka KafkaConfig

	// JWT
	JWT JWTConfig

	// Server
	Server ServerConfig

	// Stripe
	Stripe StripeConfig

	// Email
	Email EmailConfig

	// Services (for service discovery)
	Services ServicesConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	URL      string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	URL      string
}

type KafkaConfig struct {
	Brokers []string
	GroupID string
}

type JWTConfig struct {
	Secret            string
	Expiration        time.Duration
	RefreshExpiration time.Duration
}

type ServerConfig struct {
	Port            int
	Host            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type StripeConfig struct {
	SecretKey      string
	WebhookSecret  string
	PublishableKey string
}

type EmailConfig struct {
	SendGridAPIKey string
	FromEmail      string
	FromName       string
}

type ServicesConfig struct {
	UserService         string
	ProductService      string
	BiddingService      string
	OrderService        string
	PaymentService      string
	NotificationService string
	SearchService       string
	AnalyticsService    string
	AuthService         string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Env: getEnv("ENV", "development"),
	}

	// Database
	cfg.Database = DatabaseConfig{
		Host:     getEnv("DATABASE_HOST", "localhost"),
		Port:     getEnvAsInt("DATABASE_PORT", 5435), // Changed to 5435
		User:     getEnv("DATABASE_USER", "postgres"),
		Password: getEnv("DATABASE_PASSWORD", "postgres"),
		DBName:   getEnv("DATABASE_NAME", "sneakers_marketplace"),
		SSLMode:  getEnv("DATABASE_SSL_MODE", "disable"),
		URL:      getEnv("DATABASE_URL", ""),
	}

	// Build DATABASE_URL if not provided
	if cfg.Database.URL == "" {
		cfg.Database.URL = fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.DBName,
			cfg.Database.SSLMode,
		)
	}

	// Redis
	cfg.Redis = RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnvAsInt("REDIS_PORT", 6380), // Changed to 6380
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       getEnvAsInt("REDIS_DB", 0),
		URL:      getEnv("REDIS_URL", ""),
	}

	if cfg.Redis.URL == "" {
		if cfg.Redis.Password != "" {
			cfg.Redis.URL = fmt.Sprintf("redis://:%s@%s:%d/%d",
				cfg.Redis.Password, cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.DB)
		} else {
			cfg.Redis.URL = fmt.Sprintf("redis://%s:%d/%d",
				cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.DB)
		}
	}

	// Kafka
	cfg.Kafka = KafkaConfig{
		Brokers: getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9094"}), // Changed to 9094
		GroupID: getEnv("KAFKA_GROUP_ID", "sneakers-marketplace"),
	}

	// JWT
	cfg.JWT = JWTConfig{
		Secret:            getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		Expiration:        getEnvAsDuration("JWT_EXPIRATION", 24*time.Hour),
		RefreshExpiration: getEnvAsDuration("REFRESH_TOKEN_EXPIRATION", 168*time.Hour),
	}

	// Server
	cfg.Server = ServerConfig{
		Port:            getEnvAsInt("SERVER_PORT", 8080),
		Host:            getEnv("SERVER_HOST", "0.0.0.0"),
		ReadTimeout:     getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
		WriteTimeout:    getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
		ShutdownTimeout: getEnvAsDuration("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
	}

	// Stripe
	cfg.Stripe = StripeConfig{
		SecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
		WebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", ""),
		PublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", ""),
	}

	// Email
	cfg.Email = EmailConfig{
		SendGridAPIKey: getEnv("SENDGRID_API_KEY", ""),
		FromEmail:      getEnv("SENDGRID_FROM_EMAIL", "noreply@sneakersmarketplace.com"),
		FromName:       getEnv("SENDGRID_FROM_NAME", "Sneakers Marketplace"),
	}

	// Services (for gRPC connections)
	cfg.Services = ServicesConfig{
		UserService:         getEnv("USER_SERVICE_ADDR", "localhost:50051"),
		ProductService:      getEnv("PRODUCT_SERVICE_ADDR", "localhost:50052"),
		BiddingService:      getEnv("BIDDING_SERVICE_ADDR", "localhost:50053"),
		OrderService:        getEnv("ORDER_SERVICE_ADDR", "localhost:50054"),
		PaymentService:      getEnv("PAYMENT_SERVICE_ADDR", "localhost:50055"),
		NotificationService: getEnv("NOTIFICATION_SERVICE_ADDR", "localhost:50056"),
		SearchService:       getEnv("SEARCH_SERVICE_ADDR", "localhost:50057"),
		AnalyticsService:    getEnv("ANALYTICS_SERVICE_ADDR", "localhost:50058"),
		AuthService:         getEnv("AUTH_SERVICE_ADDR", "localhost:50059"),
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.JWT.Secret == "your-secret-key-change-in-production" && c.Env == "production" {
		return fmt.Errorf("JWT_SECRET must be changed in production")
	}

	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	return nil
}

// IsDevelopment returns true if environment is development
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// IsProduction returns true if environment is production
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	// Simple comma-separated values
	// For production, consider using a proper CSV parser
	return []string{valueStr}
}
