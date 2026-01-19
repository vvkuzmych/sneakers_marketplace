package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/vvkuzmych/sneakers_marketplace/internal/user/model"
	"github.com/vvkuzmych/sneakers_marketplace/internal/user/model_gorm"
	"github.com/vvkuzmych/sneakers_marketplace/internal/user/repository"
	"github.com/vvkuzmych/sneakers_marketplace/internal/user/repository_gorm"
)

const testDatabaseURL = "postgresql://postgres:postgres@localhost:5432/sneakers_marketplace?sslmode=disable"

var (
	rawRepo  *repository.UserRepository
	gormRepo *repository_gorm.UserRepository
)

// setupBenchmark initializes database connections
func setupBenchmark(b *testing.B) {
	ctx := context.Background()

	// Setup Raw SQL
	pgxPool, err := pgxpool.New(ctx, testDatabaseURL)
	if err != nil {
		b.Fatalf("Failed to connect to database (pgx): %v", err)
	}
	rawRepo = repository.NewUserRepository(pgxPool)

	// Setup GORM
	gormDB, err := gorm.Open(postgres.Open(testDatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		b.Fatalf("Failed to connect to database (GORM): %v", err)
	}
	gormRepo = repository_gorm.NewUserRepository(gormDB)
}

// BenchmarkCreate_RawSQL benchmarks user creation with Raw SQL
func BenchmarkCreate_RawSQL(b *testing.B) {
	setupBenchmark(b)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := &model.User{
			Email:        fmt.Sprintf("bench.raw.%d@example.com", time.Now().UnixNano()),
			PasswordHash: "hashed_password",
			FirstName:    "Bench",
			LastName:     "User",
			Phone:        "+1234567890",
			Role:         "user",
			IsVerified:   false,
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		err := rawRepo.Create(ctx, user)
		if err != nil {
			b.Fatalf("Failed to create user: %v", err)
		}

		// Note: Cleanup skipped - raw repo doesn't have Delete in current implementation
	}
}

// BenchmarkCreate_GORM benchmarks user creation with GORM
func BenchmarkCreate_GORM(b *testing.B) {
	setupBenchmark(b)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := &model_gorm.User{
			Email:        fmt.Sprintf("bench.gorm.%d@example.com", time.Now().UnixNano()),
			PasswordHash: "hashed_password",
			FirstName:    "Bench",
			LastName:     "User",
			Phone:        "+1234567890",
			Role:         "user",
			IsVerified:   false,
			IsActive:     true,
		}

		err := gormRepo.Create(ctx, user)
		if err != nil {
			b.Fatalf("Failed to create user: %v", err)
		}

		// Cleanup
		_ = gormRepo.Delete(ctx, user.ID)
	}
}

// BenchmarkGetByEmail_RawSQL benchmarks fetching user by email with Raw SQL
func BenchmarkGetByEmail_RawSQL(b *testing.B) {
	setupBenchmark(b)
	ctx := context.Background()

	// Create a test user
	user := &model.User{
		Email:        fmt.Sprintf("bench.raw.get.%d@example.com", time.Now().UnixNano()),
		PasswordHash: "hashed_password",
		FirstName:    "Bench",
		LastName:     "User",
		Phone:        "+1234567890",
		Role:         "user",
		IsVerified:   false,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	_ = rawRepo.Create(ctx, user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := rawRepo.GetByEmail(ctx, user.Email)
		if err != nil {
			b.Fatalf("Failed to get user: %v", err)
		}
	}

	// Note: Cleanup skipped - raw repo doesn't have Delete in current implementation
}

// BenchmarkGetByEmail_GORM benchmarks fetching user by email with GORM
func BenchmarkGetByEmail_GORM(b *testing.B) {
	setupBenchmark(b)
	ctx := context.Background()

	// Create a test user
	user := &model_gorm.User{
		Email:        fmt.Sprintf("bench.gorm.get.%d@example.com", time.Now().UnixNano()),
		PasswordHash: "hashed_password",
		FirstName:    "Bench",
		LastName:     "User",
		Phone:        "+1234567890",
		Role:         "user",
		IsVerified:   false,
		IsActive:     true,
	}
	_ = gormRepo.Create(ctx, user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := gormRepo.GetByEmail(ctx, user.Email)
		if err != nil {
			b.Fatalf("Failed to get user: %v", err)
		}
	}

	// Cleanup
	_ = gormRepo.Delete(ctx, user.ID)
}

// BenchmarkUpdate_RawSQL benchmarks user update with Raw SQL
func BenchmarkUpdate_RawSQL(b *testing.B) {
	setupBenchmark(b)
	ctx := context.Background()

	// Create a test user
	user := &model.User{
		Email:        fmt.Sprintf("bench.raw.update.%d@example.com", time.Now().UnixNano()),
		PasswordHash: "hashed_password",
		FirstName:    "Bench",
		LastName:     "User",
		Phone:        "+1234567890",
		Role:         "user",
		IsVerified:   false,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	_ = rawRepo.Create(ctx, user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user.FirstName = fmt.Sprintf("Updated-%d", i)
		user.UpdatedAt = time.Now()
		err := rawRepo.Update(ctx, user)
		if err != nil {
			b.Fatalf("Failed to update user: %v", err)
		}
	}

	// Note: Cleanup skipped - raw repo doesn't have Delete in current implementation
}

// BenchmarkUpdate_GORM benchmarks user update with GORM
func BenchmarkUpdate_GORM(b *testing.B) {
	setupBenchmark(b)
	ctx := context.Background()

	// Create a test user
	user := &model_gorm.User{
		Email:        fmt.Sprintf("bench.gorm.update.%d@example.com", time.Now().UnixNano()),
		PasswordHash: "hashed_password",
		FirstName:    "Bench",
		LastName:     "User",
		Phone:        "+1234567890",
		Role:         "user",
		IsVerified:   false,
		IsActive:     true,
	}
	_ = gormRepo.Create(ctx, user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user.FirstName = fmt.Sprintf("Updated-%d", i)
		err := gormRepo.Update(ctx, user)
		if err != nil {
			b.Fatalf("Failed to update user: %v", err)
		}
	}

	// Cleanup
	_ = gormRepo.Delete(ctx, user.ID)
}

// BenchmarkList_GORM benchmarks listing users with pagination
func BenchmarkList_GORM(b *testing.B) {
	setupBenchmark(b)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := gormRepo.List(ctx, 0, 10)
		if err != nil {
			b.Fatalf("Failed to list users: %v", err)
		}
	}
}

// BenchmarkFindActive_GORM benchmarks finding active users with scope
func BenchmarkFindActive_GORM(b *testing.B) {
	setupBenchmark(b)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := gormRepo.FindActive(ctx)
		if err != nil {
			b.Fatalf("Failed to find active users: %v", err)
		}
	}
}
