package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vvkuzmych/sneakers_marketplace/internal/user/model"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, phone, is_verified, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.IsVerified,
		user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone, 
		       COALESCE(role, 'user') as role, is_verified, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &model.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Role,
		&user.IsVerified,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone,
		       COALESCE(role, 'user') as role, 
		       is_verified, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &model.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Role,
		&user.IsVerified,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, phone = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.ID,
	).Scan(&user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// AddAddress adds a new address for a user
func (r *UserRepository) AddAddress(ctx context.Context, address *model.Address) error {
	query := `
		INSERT INTO addresses (user_id, address_type, street_line1, street_line2, 
		                      city, state, postal_code, country, is_default)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		address.UserID,
		address.AddressType,
		address.StreetLine1,
		address.StreetLine2,
		address.City,
		address.State,
		address.PostalCode,
		address.Country,
		address.IsDefault,
	).Scan(&address.ID, &address.CreatedAt, &address.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to add address: %w", err)
	}

	// If this is default address, unset other default addresses
	if address.IsDefault {
		if err := r.unsetOtherDefaultAddresses(ctx, address.UserID, address.ID, address.AddressType); err != nil {
			return err
		}
	}

	return nil
}

// GetAddresses retrieves all addresses for a user
func (r *UserRepository) GetAddresses(ctx context.Context, userID int64) ([]*model.Address, error) {
	query := `
		SELECT id, user_id, address_type, street_line1, street_line2, 
		       city, state, postal_code, country, is_default, created_at, updated_at
		FROM addresses
		WHERE user_id = $1
		ORDER BY is_default DESC, created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses: %w", err)
	}
	defer rows.Close()

	var addresses []*model.Address
	for rows.Next() {
		addr := &model.Address{}
		err := rows.Scan(
			&addr.ID,
			&addr.UserID,
			&addr.AddressType,
			&addr.StreetLine1,
			&addr.StreetLine2,
			&addr.City,
			&addr.State,
			&addr.PostalCode,
			&addr.Country,
			&addr.IsDefault,
			&addr.CreatedAt,
			&addr.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan address: %w", err)
		}
		addresses = append(addresses, addr)
	}

	return addresses, nil
}

// UpdateAddress updates an address
func (r *UserRepository) UpdateAddress(ctx context.Context, address *model.Address) error {
	query := `
		UPDATE addresses
		SET address_type = $1, street_line1 = $2, street_line2 = $3,
		    city = $4, state = $5, postal_code = $6, country = $7,
		    is_default = $8, updated_at = NOW()
		WHERE id = $9 AND user_id = $10
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		address.AddressType,
		address.StreetLine1,
		address.StreetLine2,
		address.City,
		address.State,
		address.PostalCode,
		address.Country,
		address.IsDefault,
		address.ID,
		address.UserID,
	).Scan(&address.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update address: %w", err)
	}

	// If this is default address, unset other default addresses
	if address.IsDefault {
		if err := r.unsetOtherDefaultAddresses(ctx, address.UserID, address.ID, address.AddressType); err != nil {
			return err
		}
	}

	return nil
}

// DeleteAddress deletes an address
func (r *UserRepository) DeleteAddress(ctx context.Context, addressID, userID int64) error {
	query := `DELETE FROM addresses WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(ctx, query, addressID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete address: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("address not found")
	}

	return nil
}

// CreateSession creates a new session
func (r *UserRepository) CreateSession(ctx context.Context, session *model.Session) error {
	query := `
		INSERT INTO sessions (user_id, token_hash, refresh_token_hash, ip_address, user_agent, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query,
		session.UserID,
		session.TokenHash,
		session.RefreshTokenHash,
		session.IPAddress,
		session.UserAgent,
		session.ExpiresAt,
	).Scan(&session.ID, &session.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// DeleteSession deletes a session by token hash
func (r *UserRepository) DeleteSession(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM sessions WHERE token_hash = $1`

	_, err := r.db.Exec(ctx, query, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

// DeleteExpiredSessions deletes all expired sessions
func (r *UserRepository) DeleteExpiredSessions(ctx context.Context) error {
	query := `DELETE FROM sessions WHERE expires_at < NOW()`

	_, err := r.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	return nil
}

// Helper: unset other default addresses of the same type
func (r *UserRepository) unsetOtherDefaultAddresses(ctx context.Context, userID, addressID int64, addressType string) error {
	query := `
		UPDATE addresses
		SET is_default = false
		WHERE user_id = $1 AND id != $2 AND address_type = $3 AND is_default = true
	`

	_, err := r.db.Exec(ctx, query, userID, addressID, addressType)
	if err != nil {
		return fmt.Errorf("failed to unset other default addresses: %w", err)
	}

	return nil
}
