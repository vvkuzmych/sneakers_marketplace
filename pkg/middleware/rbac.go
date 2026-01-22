package middleware

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/golang-jwt/jwt/v5"
)

// Role represents a user role
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// UserContext holds user information from JWT
type UserContext struct {
	Email  string
	Role   Role
	UserID int64
}

// contextKey is the key for storing UserContext in context
type contextKey string

const userContextKey contextKey = "user_context"

// RequireRole creates a gRPC interceptor that requires specific role
func RequireRole(jwtSecret string, requiredRole Role) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Extract JWT from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization token")
		}

		tokenString := tokens[0]
		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// Validate JWT and extract user info
		userCtx, err := validateJWTAndExtractUser(tokenString, jwtSecret)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token: "+err.Error())
		}

		// Check if user has required role
		if !hasRole(userCtx.Role, requiredRole) {
			return nil, status.Errorf(codes.PermissionDenied,
				"insufficient permissions: requires %s role", requiredRole)
		}

		// Add user context to request context
		ctx = context.WithValue(ctx, userContextKey, userCtx)

		// Call the handler
		return handler(ctx, req)
	}
}

// RequireAdmin creates an interceptor that requires admin role
func RequireAdmin(jwtSecret string) grpc.UnaryServerInterceptor {
	return RequireRole(jwtSecret, RoleAdmin)
}

// RequireAuthentication creates an interceptor that requires any authenticated user
func RequireAuthentication(jwtSecret string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Extract JWT from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization token")
		}

		tokenString := tokens[0]
		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// Validate JWT and extract user info
		userCtx, err := validateJWTAndExtractUser(tokenString, jwtSecret)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token: "+err.Error())
		}

		// Add user context to request context
		ctx = context.WithValue(ctx, userContextKey, userCtx)

		// Call the handler
		return handler(ctx, req)
	}
}

// GetUserFromContext extracts UserContext from context
func GetUserFromContext(ctx context.Context) (*UserContext, error) {
	userCtx, ok := ctx.Value(userContextKey).(*UserContext)
	if !ok {
		return nil, fmt.Errorf("user context not found")
	}
	return userCtx, nil
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (int64, error) {
	userCtx, err := GetUserFromContext(ctx)
	if err != nil {
		return 0, err
	}
	return userCtx.UserID, nil
}

// IsAdmin checks if user in context has admin role
func IsAdmin(ctx context.Context) bool {
	userCtx, err := GetUserFromContext(ctx)
	if err != nil {
		return false
	}
	return userCtx.Role == RoleAdmin
}

// validateJWTAndExtractUser validates JWT token and extracts user information
func validateJWTAndExtractUser(tokenString, jwtSecret string) (*UserContext, error) {
	// Parse JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Extract user_id
	userID, err := extractUserID(claims)
	if err != nil {
		return nil, err
	}

	// Extract email
	email, ok := claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("missing email in token")
	}

	// Extract role (defaults to "user" if not present)
	role := RoleUser
	if roleStr, ok := claims["role"].(string); ok {
		role = Role(roleStr)
	}

	return &UserContext{
		UserID: userID,
		Email:  email,
		Role:   role,
	}, nil
}

// extractUserID extracts user ID from JWT claims
func extractUserID(claims jwt.MapClaims) (int64, error) {
	userIDClaim, ok := claims["user_id"]
	if !ok {
		return 0, fmt.Errorf("missing user_id in token")
	}

	// Handle different numeric types
	switch v := userIDClaim.(type) {
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("invalid user_id type: %T", v)
	}
}

// hasRole checks if user has the required role
func hasRole(userRole, requiredRole Role) bool {
	// Admin has access to everything
	if userRole == RoleAdmin {
		return true
	}
	// Otherwise, exact match required
	return userRole == requiredRole
}

// ChainInterceptors chains multiple unary interceptors
func ChainInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Build chain of handlers
		wrappedHandler := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			currentInterceptor := interceptors[i]
			currentHandler := wrappedHandler
			wrappedHandler = func(ctx context.Context, req interface{}) (interface{}, error) {
				return currentInterceptor(ctx, req, info, currentHandler)
			}
		}
		return wrappedHandler(ctx, req)
	}
}

// MethodMatcher allows different middleware for different methods
type MethodMatcher struct {
	adminMethods  map[string]bool
	publicMethods map[string]bool
}

// NewMethodMatcher creates a new method matcher
func NewMethodMatcher() *MethodMatcher {
	return &MethodMatcher{
		adminMethods:  make(map[string]bool),
		publicMethods: make(map[string]bool),
	}
}

// AddAdminMethod adds a method that requires admin role
func (m *MethodMatcher) AddAdminMethod(methodName string) {
	m.adminMethods[methodName] = true
}

// AddPublicMethod adds a method that doesn't require authentication
func (m *MethodMatcher) AddPublicMethod(methodName string) {
	m.publicMethods[methodName] = true
}

// CreateInterceptor creates an interceptor with method-specific rules
func (m *MethodMatcher) CreateInterceptor(jwtSecret string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Check if method is public
		if m.publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		// Check if method requires admin
		if m.adminMethods[info.FullMethod] {
			return RequireAdmin(jwtSecret)(ctx, req, info, handler)
		}

		// Default: require authentication
		return RequireAuthentication(jwtSecret)(ctx, req, info, handler)
	}
}

// LoggingInterceptor logs all requests with user information
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Try to get user context
	userCtx, err := GetUserFromContext(ctx)
	if err == nil {
		// User is authenticated, log with user info
		fmt.Printf("[RBAC] Method: %s, User: %d (%s), Role: %s\n",
			info.FullMethod, userCtx.UserID, userCtx.Email, userCtx.Role)
	} else {
		// Unauthenticated request
		fmt.Printf("[RBAC] Method: %s, User: unauthenticated\n", info.FullMethod)
	}

	return handler(ctx, req)
}
