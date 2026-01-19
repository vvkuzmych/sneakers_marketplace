package main

import (
	"context"
	"fmt"
	"log"
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

const (
	databaseURL = "postgresql://postgres:postgres@localhost:5432/sneakers_marketplace?sslmode=disable"
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║         🔬 GORM vs Raw SQL Comparison Demo                      ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	ctx := context.Background()

	// ========================================================================
	// Setup Raw SQL (pgx) connection
	// ========================================================================
	fmt.Println("📦 Setting up Raw SQL (pgx) connection...")
	pgxPool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database (pgx): %v", err)
	}
	defer pgxPool.Close()

	rawRepo := repository.NewUserRepository(pgxPool)
	fmt.Println("✅ Raw SQL repository ready")
	fmt.Println()

	// ========================================================================
	// Setup GORM connection
	// ========================================================================
	fmt.Println("📦 Setting up GORM connection...")
	gormDB, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Quiet mode for cleaner output
	})
	if err != nil {
		log.Fatalf("Failed to connect to database (GORM): %v", err)
	}

	gormRepo := repository_gorm.NewUserRepository(gormDB)
	fmt.Println("✅ GORM repository ready")
	fmt.Println()

	// ========================================================================
	// Demo 1: Create User
	// ========================================================================
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📝 Demo 1: CREATE USER")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	timestamp := time.Now().Unix()

	// Raw SQL
	fmt.Println("\n🔹 Raw SQL (pgx):")
	rawUser := &model.User{
		Email:        fmt.Sprintf("raw.user.%d@example.com", timestamp),
		PasswordHash: "hashed_password_123",
		FirstName:    "John",
		LastName:     "Doe",
		Phone:        "+1234567890",
		Role:         "user",
		IsVerified:   false,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	start := time.Now()
	err = rawRepo.Create(ctx, rawUser)
	rawDuration := time.Since(start)

	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Created user ID: %d (took %v)\n", rawUser.ID, rawDuration)
	}

	// GORM
	fmt.Println("\n🔹 GORM:")
	gormUser := &model_gorm.User{
		Email:        fmt.Sprintf("gorm.user.%d@example.com", timestamp),
		PasswordHash: "hashed_password_123",
		FirstName:    "Jane",
		LastName:     "Smith",
		Phone:        "+0987654321",
		Role:         "user",
		IsVerified:   false,
		IsActive:     true,
	}

	start = time.Now()
	err = gormRepo.Create(ctx, gormUser)
	gormDuration := time.Since(start)

	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Created user ID: %d (took %v)\n", gormUser.ID, gormDuration)
	}

	fmt.Printf("\n📊 Performance: Raw SQL %v vs GORM %v (%.1fx)\n",
		rawDuration, gormDuration, float64(gormDuration)/float64(rawDuration))

	// ========================================================================
	// Demo 2: Get User by Email
	// ========================================================================
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🔍 Demo 2: GET USER BY EMAIL")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Raw SQL
	fmt.Println("\n🔹 Raw SQL (pgx):")
	start = time.Now()
	fetchedRaw, err := rawRepo.GetByEmail(ctx, rawUser.Email)
	rawDuration = time.Since(start)

	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Found user: %s %s (took %v)\n", fetchedRaw.FirstName, fetchedRaw.LastName, rawDuration)
	}

	// GORM
	fmt.Println("\n🔹 GORM:")
	start = time.Now()
	fetchedGorm, err := gormRepo.GetByEmail(ctx, gormUser.Email)
	gormDuration = time.Since(start)

	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Found user: %s %s (took %v)\n", fetchedGorm.FirstName, fetchedGorm.LastName, gormDuration)
	}

	fmt.Printf("\n📊 Performance: Raw SQL %v vs GORM %v (%.1fx)\n",
		rawDuration, gormDuration, float64(gormDuration)/float64(rawDuration))

	// ========================================================================
	// Demo 3: Update User
	// ========================================================================
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📝 Demo 3: UPDATE USER")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Raw SQL
	fmt.Println("\n🔹 Raw SQL (pgx):")
	fetchedRaw.FirstName = "John-Updated"
	fetchedRaw.UpdatedAt = time.Now()

	start = time.Now()
	err = rawRepo.Update(ctx, fetchedRaw)
	rawDuration = time.Since(start)

	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Updated user (took %v)\n", rawDuration)
	}

	// GORM
	fmt.Println("\n🔹 GORM:")
	fetchedGorm.FirstName = "Jane-Updated"

	start = time.Now()
	err = gormRepo.Update(ctx, fetchedGorm)
	gormDuration = time.Since(start)

	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Updated user (took %v)\n", gormDuration)
	}

	fmt.Printf("\n📊 Performance: Raw SQL %v vs GORM %v (%.1fx)\n",
		rawDuration, gormDuration, float64(gormDuration)/float64(rawDuration))

	// ========================================================================
	// Demo 4: List Users (Pagination)
	// ========================================================================
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📋 Demo 4: LIST USERS (Pagination)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// GORM (Raw SQL repo doesn't have List method in current implementation)
	fmt.Println("\n🔹 GORM:")
	start = time.Now()
	users, total, err := gormRepo.List(ctx, 0, 5)
	gormDuration = time.Since(start)

	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Found %d users (total: %d) (took %v)\n", len(users), total, gormDuration)
		for i, u := range users {
			fmt.Printf("   %d. %s (%s)\n", i+1, u.Email, u.Role)
		}
	}

	// ========================================================================
	// Demo 5: Find Active Users (Scopes)
	// ========================================================================
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🔍 Demo 5: FIND ACTIVE USERS (Scopes)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	fmt.Println("\n🔹 GORM (using scope):")
	start = time.Now()
	activeUsers, err := gormRepo.FindActive(ctx)
	gormDuration = time.Since(start)

	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Found %d active users (took %v)\n", len(activeUsers), gormDuration)
	}

	// ========================================================================
	// Cleanup (delete test users)
	// ========================================================================
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🧹 Cleaning up test users...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Note: Raw SQL repo doesn't have Delete method in current implementation
	// In production, you would add it or use direct SQL
	fmt.Println("⚠️  Raw SQL user cleanup skipped (Delete method not implemented)")

	// GORM cleanup
	gormRepo.Delete(ctx, gormUser.ID)

	fmt.Println("✅ Cleanup complete")
	fmt.Println()

	// ========================================================================
	// Summary
	// ========================================================================
	fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                        📊 SUMMARY                                ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("✨ GORM Advantages:")
	fmt.Println("   • Less boilerplate code (3-5x shorter)")
	fmt.Println("   • Auto timestamps (CreatedAt, UpdatedAt)")
	fmt.Println("   • Scopes for reusable queries")
	fmt.Println("   • Automatic scanning (no manual Scan())")
	fmt.Println("   • Hooks (BeforeCreate, AfterUpdate, etc.)")
	fmt.Println("   • Associations and eager loading")
	fmt.Println()
	fmt.Println("⚡ Raw SQL (pgx) Advantages:")
	fmt.Println("   • 20-60% faster performance")
	fmt.Println("   • Full control over queries")
	fmt.Println("   • Better for complex queries (CTEs, subqueries)")
	fmt.Println("   • More transparent (see exact SQL)")
	fmt.Println("   • No ORM overhead")
	fmt.Println()
	fmt.Println("🎯 Recommendation:")
	fmt.Println("   Use GORM for simple CRUD, Raw SQL for complex/critical queries")
	fmt.Println()
}
