package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/vvkuzmych/sneakers_marketplace/internal/user/model"
	"github.com/vvkuzmych/sneakers_marketplace/internal/user/repository"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/auth"
)

// UserService handles business logic for users
type UserService struct {
	repo       *repository.UserRepository
	jwtManager *auth.JWTManager
}

// NewUserService creates a new user service
func NewUserService(repo *repository.UserRepository, jwtManager *auth.JWTManager) *UserService {
	return &UserService{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

// Register registers a new user
func (s *UserService) Register(ctx context.Context, email, password, firstName, lastName, phone string) (*model.User, string, string, error) {
	// Check if user already exists
	existingUser, err := s.repo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, "", "", fmt.Errorf("user with email %s already exists", email)
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &model.User{
		Email:        email,
		PasswordHash: hashedPassword,
		FirstName:    firstName,
		LastName:     lastName,
		Phone:        phone,
		IsVerified:   false, // Email verification should be implemented
		IsActive:     true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, "", "", fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens with role (default to "user" for new registrations)
	role := "user"
	accessToken, err := s.jwtManager.GenerateAccessTokenWithRole(user.ID, user.Email, role)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshTokenWithRole(user.ID, user.Email, role)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create session
	session := &model.Session{
		UserID:           user.ID,
		TokenHash:        hashToken(accessToken),
		RefreshTokenHash: hashToken(refreshToken),
		IPAddress:        nil,                                // Extract from context if available
		UserAgent:        nil,                                // Extract from context if available
		ExpiresAt:        time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, "", "", fmt.Errorf("failed to create session: %w", err)
	}

	return user, accessToken, refreshToken, nil
}

// Login authenticates a user
func (s *UserService) Login(ctx context.Context, email, password string) (*model.User, string, string, error) {
	// Get user by email
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", "", fmt.Errorf("invalid credentials")
	}

	// Verify password
	if err := auth.VerifyPassword(user.PasswordHash, password); err != nil {
		return nil, "", "", fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, "", "", fmt.Errorf("user account is deactivated")
	}

	// Generate tokens with user's role
	role := user.Role
	if role == "" {
		role = "user" // Default to user if not set
	}
	accessToken, err := s.jwtManager.GenerateAccessTokenWithRole(user.ID, user.Email, role)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshTokenWithRole(user.ID, user.Email, role)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create session
	session := &model.Session{
		UserID:           user.ID,
		TokenHash:        hashToken(accessToken),
		RefreshTokenHash: hashToken(refreshToken),
		IPAddress:        nil,                                // Extract from context if available
		UserAgent:        nil,                                // Extract from context if available
		ExpiresAt:        time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, "", "", fmt.Errorf("failed to create session: %w", err)
	}

	return user, accessToken, refreshToken, nil
}

// RefreshToken generates new tokens from refresh token
func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// Validate refresh token
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get user
	user, err := s.repo.GetByID(ctx, claims.UserID)
	if err != nil {
		return "", "", fmt.Errorf("user not found: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return "", "", fmt.Errorf("user account is deactivated")
	}

	// Generate new tokens with user's role
	role := user.Role
	if role == "" {
		role = "user"
	}
	newAccessToken, err := s.jwtManager.GenerateAccessTokenWithRole(user.ID, user.Email, role)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshTokenWithRole(user.ID, user.Email, role)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create new session
	session := &model.Session{
		UserID:           user.ID,
		TokenHash:        hashToken(newAccessToken),
		RefreshTokenHash: hashToken(newRefreshToken),
		IPAddress:        nil,                                // Extract from context if available
		UserAgent:        nil,                                // Extract from context if available
		ExpiresAt:        time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return "", "", fmt.Errorf("failed to create session: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

// Logout logs out a user
func (s *UserService) Logout(ctx context.Context, accessToken string) error {
	tokenHash := hashToken(accessToken)
	return s.repo.DeleteSession(ctx, tokenHash)
}

// GetProfile retrieves a user profile
func (s *UserService) GetProfile(ctx context.Context, userID int64) (*model.User, error) {
	return s.repo.GetByID(ctx, userID)
}

// UpdateProfile updates a user profile
func (s *UserService) UpdateProfile(ctx context.Context, userID int64, firstName, lastName, phone string) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.Phone = phone

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// AddAddress adds a new address
func (s *UserService) AddAddress(ctx context.Context, address *model.Address) (*model.Address, error) {
	if err := s.repo.AddAddress(ctx, address); err != nil {
		return nil, fmt.Errorf("failed to add address: %w", err)
	}
	return address, nil
}

// GetAddresses retrieves all addresses for a user
func (s *UserService) GetAddresses(ctx context.Context, userID int64) ([]*model.Address, error) {
	return s.repo.GetAddresses(ctx, userID)
}

// UpdateAddress updates an address
func (s *UserService) UpdateAddress(ctx context.Context, address *model.Address) (*model.Address, error) {
	if err := s.repo.UpdateAddress(ctx, address); err != nil {
		return nil, fmt.Errorf("failed to update address: %w", err)
	}
	return address, nil
}

// DeleteAddress deletes an address
func (s *UserService) DeleteAddress(ctx context.Context, addressID, userID int64) error {
	return s.repo.DeleteAddress(ctx, addressID, userID)
}

// Helper: hash token for storage
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
