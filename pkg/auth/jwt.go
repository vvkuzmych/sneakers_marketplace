package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTManager manages JWT token generation and validation
type JWTManager struct {
	secretKey            []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

// Claims represents JWT claims
type Claims struct {
	jwt.RegisteredClaims
	Email  string `json:"email"`
	Role   string `json:"role"`
	UserID int64  `json:"user_id"`
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secretKey string, accessDuration, refreshDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:            []byte(secretKey),
		accessTokenDuration:  accessDuration,
		refreshTokenDuration: refreshDuration,
	}
}

// GenerateAccessToken generates a new access token (defaults to "user" role)
func (m *JWTManager) GenerateAccessToken(userID int64, email string) (string, error) {
	return m.GenerateAccessTokenWithRole(userID, email, "user")
}

// GenerateAccessTokenWithRole generates a new access token with specified role
func (m *JWTManager) GenerateAccessTokenWithRole(userID int64, email, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

// GenerateRefreshToken generates a new refresh token (defaults to "user" role)
func (m *JWTManager) GenerateRefreshToken(userID int64, email string) (string, error) {
	return m.GenerateRefreshTokenWithRole(userID, email, "user")
}

// GenerateRefreshTokenWithRole generates a new refresh token with specified role
func (m *JWTManager) GenerateRefreshTokenWithRole(userID int64, email, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

// ValidateToken validates a JWT token and returns the claims
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// ExtractUserID extracts user ID from token without full validation
func (m *JWTManager) ExtractUserID(tokenString string) (int64, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
