package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vvkuzmych/sneakers_marketplace/internal/notification/model"
)

type NotificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// CreateNotification creates a new notification
func (r *NotificationRepository) CreateNotification(ctx context.Context, notif *model.Notification) (*model.Notification, error) {
	query := `
		INSERT INTO notifications (user_id, type, title, message, data, email_sent, push_sent, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query,
		notif.UserID, notif.Type, notif.Title, notif.Message, notif.Data,
		notif.EmailSent, notif.PushSent, notif.IsRead, notif.CreatedAt,
	).Scan(&notif.ID, &notif.CreatedAt)

	return notif, err
}

// GetNotificationByID retrieves a notification by ID
func (r *NotificationRepository) GetNotificationByID(ctx context.Context, id int64) (*model.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, data,
		       email_sent, email_sent_at, push_sent, push_sent_at,
		       is_read, read_at, created_at
		FROM notifications
		WHERE id = $1`

	var n model.Notification
	err := r.db.QueryRow(ctx, query, id).Scan(
		&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.Data,
		&n.EmailSent, &n.EmailSentAt, &n.PushSent, &n.PushSentAt,
		&n.IsRead, &n.ReadAt, &n.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("notification not found")
	}
	return &n, err
}

// ListNotifications retrieves notifications for a user with pagination
func (r *NotificationRepository) ListNotifications(
	ctx context.Context, userID int64, unreadOnly bool, page, pageSize int32,
) ([]*model.Notification, int64, error) {
	offset := (page - 1) * pageSize

	whereClause := "WHERE user_id = $1"
	if unreadOnly {
		whereClause += " AND is_read = FALSE"
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM notifications %s", whereClause)
	var total int64
	err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get notifications
	query := fmt.Sprintf(`
		SELECT id, user_id, type, title, message, data,
		       email_sent, email_sent_at, push_sent, push_sent_at,
		       is_read, read_at, created_at
		FROM notifications
		%s
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`, whereClause)

	rows, err := r.db.Query(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	notifications := make([]*model.Notification, 0)
	for rows.Next() {
		var n model.Notification
		err := rows.Scan(
			&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.Data,
			&n.EmailSent, &n.EmailSentAt, &n.PushSent, &n.PushSentAt,
			&n.IsRead, &n.ReadAt, &n.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		notifications = append(notifications, &n)
	}

	return notifications, total, rows.Err()
}

// GetUnreadCount returns the count of unread notifications for a user
func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID int64) (int64, error) {
	var count int64
	query := "SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = FALSE"
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

// MarkAsRead marks a notification as read
func (r *NotificationRepository) MarkAsRead(ctx context.Context, id, userID int64) error {
	query := `
		UPDATE notifications
		SET is_read = TRUE, read_at = $1
		WHERE id = $2 AND user_id = $3 AND is_read = FALSE`

	tag, err := r.db.Exec(ctx, query, time.Now(), id, userID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("notification not found or already read")
	}

	return nil
}

// MarkAllAsRead marks all user's notifications as read
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID int64) (int64, error) {
	query := `
		UPDATE notifications
		SET is_read = TRUE, read_at = $1
		WHERE user_id = $2 AND is_read = FALSE`

	tag, err := r.db.Exec(ctx, query, time.Now(), userID)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}

// DeleteNotification deletes a notification
func (r *NotificationRepository) DeleteNotification(ctx context.Context, id, userID int64) error {
	query := "DELETE FROM notifications WHERE id = $1 AND user_id = $2"

	tag, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

// UpdateEmailSent marks email as sent
func (r *NotificationRepository) UpdateEmailSent(ctx context.Context, id int64) error {
	query := "UPDATE notifications SET email_sent = TRUE, email_sent_at = $1 WHERE id = $2"
	_, err := r.db.Exec(ctx, query, time.Now(), id)
	return err
}

// UpdatePushSent marks push as sent
func (r *NotificationRepository) UpdatePushSent(ctx context.Context, id int64) error {
	query := "UPDATE notifications SET push_sent = TRUE, push_sent_at = $1 WHERE id = $2"
	_, err := r.db.Exec(ctx, query, time.Now(), id)
	return err
}

// ===== Preferences =====

// GetPreferences retrieves user's notification preferences
func (r *NotificationRepository) GetPreferences(ctx context.Context, userID int64) (*model.NotificationPreferences, error) {
	query := `
		SELECT user_id, email_enabled, email_match_created, email_order_created,
		       email_order_shipped, email_payment_received, email_payout_completed,
		       push_enabled, push_match_created, push_order_updates, push_payment_updates,
		       inapp_enabled, created_at, updated_at
		FROM notification_preferences
		WHERE user_id = $1`

	var p model.NotificationPreferences
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&p.UserID, &p.EmailEnabled, &p.EmailMatchCreated, &p.EmailOrderCreated,
		&p.EmailOrderShipped, &p.EmailPaymentReceived, &p.EmailPayoutCompleted,
		&p.PushEnabled, &p.PushMatchCreated, &p.PushOrderUpdates, &p.PushPaymentUpdates,
		&p.InAppEnabled, &p.CreatedAt, &p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Create default preferences
		return r.CreateDefaultPreferences(ctx, userID)
	}

	return &p, err
}

// CreateDefaultPreferences creates default notification preferences for a user
func (r *NotificationRepository) CreateDefaultPreferences(ctx context.Context, userID int64) (*model.NotificationPreferences, error) {
	query := `
		INSERT INTO notification_preferences (user_id)
		VALUES ($1)
		RETURNING user_id, email_enabled, email_match_created, email_order_created,
		          email_order_shipped, email_payment_received, email_payout_completed,
		          push_enabled, push_match_created, push_order_updates, push_payment_updates,
		          inapp_enabled, created_at, updated_at`

	var p model.NotificationPreferences
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&p.UserID, &p.EmailEnabled, &p.EmailMatchCreated, &p.EmailOrderCreated,
		&p.EmailOrderShipped, &p.EmailPaymentReceived, &p.EmailPayoutCompleted,
		&p.PushEnabled, &p.PushMatchCreated, &p.PushOrderUpdates, &p.PushPaymentUpdates,
		&p.InAppEnabled, &p.CreatedAt, &p.UpdatedAt,
	)

	return &p, err
}

// UpdatePreferences updates user's notification preferences
func (r *NotificationRepository) UpdatePreferences(ctx context.Context, prefs *model.NotificationPreferences) (*model.NotificationPreferences, error) {
	query := `
		UPDATE notification_preferences
		SET email_enabled = $2, email_match_created = $3, email_order_created = $4,
		    email_order_shipped = $5, email_payment_received = $6, email_payout_completed = $7,
		    push_enabled = $8, push_match_created = $9, push_order_updates = $10,
		    push_payment_updates = $11, inapp_enabled = $12, updated_at = $13
		WHERE user_id = $1
		RETURNING user_id, email_enabled, email_match_created, email_order_created,
		          email_order_shipped, email_payment_received, email_payout_completed,
		          push_enabled, push_match_created, push_order_updates, push_payment_updates,
		          inapp_enabled, created_at, updated_at`

	prefs.UpdatedAt = time.Now()

	var p model.NotificationPreferences
	err := r.db.QueryRow(ctx, query,
		prefs.UserID, prefs.EmailEnabled, prefs.EmailMatchCreated, prefs.EmailOrderCreated,
		prefs.EmailOrderShipped, prefs.EmailPaymentReceived, prefs.EmailPayoutCompleted,
		prefs.PushEnabled, prefs.PushMatchCreated, prefs.PushOrderUpdates,
		prefs.PushPaymentUpdates, prefs.InAppEnabled, prefs.UpdatedAt,
	).Scan(
		&p.UserID, &p.EmailEnabled, &p.EmailMatchCreated, &p.EmailOrderCreated,
		&p.EmailOrderShipped, &p.EmailPaymentReceived, &p.EmailPayoutCompleted,
		&p.PushEnabled, &p.PushMatchCreated, &p.PushOrderUpdates, &p.PushPaymentUpdates,
		&p.InAppEnabled, &p.CreatedAt, &p.UpdatedAt,
	)

	return &p, err
}
