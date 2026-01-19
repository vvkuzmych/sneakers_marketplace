package repository_gorm

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/vvkuzmych/sneakers_marketplace/internal/user/model_gorm"
)

// UserRepository handles database operations for users using GORM
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new GORM-based user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *model_gorm.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model_gorm.User, error) {
	var user model_gorm.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model_gorm.User, error) {
	var user model_gorm.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *model_gorm.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&model_gorm.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// List retrieves users with pagination
func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]model_gorm.User, int64, error) {
	var users []model_gorm.User
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&model_gorm.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

// FindActive finds active users
func (r *UserRepository) FindActive(ctx context.Context) ([]model_gorm.User, error) {
	var users []model_gorm.User
	if err := r.db.WithContext(ctx).
		Scopes(model_gorm.ActiveUsers).
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to find active users: %w", err)
	}
	return users, nil
}

// FindAdmins finds admin users
func (r *UserRepository) FindAdmins(ctx context.Context) ([]model_gorm.User, error) {
	var users []model_gorm.User
	if err := r.db.WithContext(ctx).
		Scopes(model_gorm.AdminUsers, model_gorm.ActiveUsers).
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to find admin users: %w", err)
	}
	return users, nil
}

// UpdateRole updates a user's role
func (r *UserRepository) UpdateRole(ctx context.Context, id int64, role string) error {
	if err := r.db.WithContext(ctx).
		Model(&model_gorm.User{}).
		Where("id = ?", id).
		Update("role", role).Error; err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}
	return nil
}

// Ban bans a user (sets is_active to false)
func (r *UserRepository) Ban(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).
		Model(&model_gorm.User{}).
		Where("id = ?", id).
		Update("is_active", false).Error; err != nil {
		return fmt.Errorf("failed to ban user: %w", err)
	}
	return nil
}

// Unban unbans a user (sets is_active to true)
func (r *UserRepository) Unban(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).
		Model(&model_gorm.User{}).
		Where("id = ?", id).
		Update("is_active", true).Error; err != nil {
		return fmt.Errorf("failed to unban user: %w", err)
	}
	return nil
}

// ============================================================================
// Address Repository Methods
// ============================================================================

type AddressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) *AddressRepository {
	return &AddressRepository{db: db}
}

// CreateAddress creates a new address
func (r *AddressRepository) CreateAddress(ctx context.Context, address *model_gorm.Address) error {
	return r.db.WithContext(ctx).Create(address).Error
}

// GetAddressByID retrieves an address by ID
func (r *AddressRepository) GetAddressByID(ctx context.Context, id int64) (*model_gorm.Address, error) {
	var address model_gorm.Address
	if err := r.db.WithContext(ctx).First(&address, id).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

// GetUserAddresses retrieves all addresses for a user
func (r *AddressRepository) GetUserAddresses(ctx context.Context, userID int64) ([]model_gorm.Address, error) {
	var addresses []model_gorm.Address
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}

// GetUserAddressesWithUser retrieves addresses with user preloaded
func (r *AddressRepository) GetUserAddressesWithUser(ctx context.Context, userID int64) ([]model_gorm.Address, error) {
	var addresses []model_gorm.Address
	if err := r.db.WithContext(ctx).
		Preload("User"). // Eager load user
		Where("user_id = ?", userID).
		Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}

// ============================================================================
// Session Repository Methods
// ============================================================================

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// CreateSession creates a new session
func (r *SessionRepository) CreateSession(ctx context.Context, session *model_gorm.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetSessionByToken retrieves a session by token hash
func (r *SessionRepository) GetSessionByToken(ctx context.Context, tokenHash string) (*model_gorm.Session, error) {
	var session model_gorm.Session
	if err := r.db.WithContext(ctx).
		Where("token_hash = ?", tokenHash).
		First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// GetSessionWithUser retrieves session with user preloaded
func (r *SessionRepository) GetSessionWithUser(ctx context.Context, tokenHash string) (*model_gorm.Session, error) {
	var session model_gorm.Session
	if err := r.db.WithContext(ctx).
		Preload("User"). // Eager load user
		Where("token_hash = ?", tokenHash).
		First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// DeleteSession deletes a session
func (r *SessionRepository) DeleteSession(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model_gorm.Session{}, id).Error
}

// DeleteUserSessions deletes all sessions for a user
func (r *SessionRepository) DeleteUserSessions(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&model_gorm.Session{}).Error
}
