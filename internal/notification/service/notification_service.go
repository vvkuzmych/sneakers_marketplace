package service

import (
	"context"
	"fmt"
	"log"

	"github.com/vvkuzmych/sneakers_marketplace/internal/notification/email"
	"github.com/vvkuzmych/sneakers_marketplace/internal/notification/model"
	"github.com/vvkuzmych/sneakers_marketplace/internal/notification/repository"
)

type NotificationService struct {
	repo         *repository.NotificationRepository
	emailService *email.EmailService
}

func NewNotificationService(repo *repository.NotificationRepository, emailService *email.EmailService) *NotificationService {
	return &NotificationService{
		repo:         repo,
		emailService: emailService,
	}
}

// SendNotification creates and sends a notification
func (s *NotificationService) SendNotification(
	ctx context.Context,
	userID int64,
	notifType, title, message string,
	data map[string]interface{},
	sendEmail, sendPush bool,
) (*model.Notification, error) {
	// Create notification record
	notif, err := model.NewNotification(userID, notifType, title, message, data)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	// Save to database
	created, err := s.repo.CreateNotification(ctx, notif)
	if err != nil {
		return nil, fmt.Errorf("failed to save notification: %w", err)
	}

	// Send email if requested (async, non-blocking)
	if sendEmail {
		go func() {
			// Get user preferences
			prefs, err := s.repo.GetPreferences(context.Background(), userID)
			if err != nil || !prefs.ShouldSendEmail(notifType) {
				return
			}

			// Send email (simplified - using user_id as email for now)
			// In real app, you'd fetch user's email from User Service
			userEmail := fmt.Sprintf("user%d@example.com", userID)
			err = s.emailService.SendEmail(userEmail, title, message)
			if err != nil {
				log.Printf("Failed to send email: %v", err)
				return
			}

			// Mark as sent
			_ = s.repo.UpdateEmailSent(context.Background(), created.ID)
		}()
	}

	// TODO: Send push/WebSocket if requested

	return created, nil
}

// GetNotifications retrieves user's notifications
func (s *NotificationService) GetNotifications(
	ctx context.Context,
	userID int64,
	unreadOnly bool,
	page, pageSize int32,
) ([]*model.Notification, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.repo.ListNotifications(ctx, userID, unreadOnly, page, pageSize)
}

// GetUnreadCount returns count of unread notifications
func (s *NotificationService) GetUnreadCount(ctx context.Context, userID int64) (int64, error) {
	return s.repo.GetUnreadCount(ctx, userID)
}

// MarkAsRead marks a notification as read
func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID, userID int64) error {
	return s.repo.MarkAsRead(ctx, notificationID, userID)
}

// MarkAllAsRead marks all user's notifications as read
func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID int64) (int64, error) {
	return s.repo.MarkAllAsRead(ctx, userID)
}

// DeleteNotification deletes a notification
func (s *NotificationService) DeleteNotification(ctx context.Context, notificationID, userID int64) error {
	return s.repo.DeleteNotification(ctx, notificationID, userID)
}

// GetPreferences retrieves user's preferences
func (s *NotificationService) GetPreferences(ctx context.Context, userID int64) (*model.NotificationPreferences, error) {
	return s.repo.GetPreferences(ctx, userID)
}

// UpdatePreferences updates user's preferences
func (s *NotificationService) UpdatePreferences(ctx context.Context, prefs *model.NotificationPreferences) (*model.NotificationPreferences, error) {
	return s.repo.UpdatePreferences(ctx, prefs)
}

// NotifyMatchCreated sends notifications when a match is created
func (s *NotificationService) NotifyMatchCreated(
	ctx context.Context,
	matchID, buyerID, sellerID, productID int64,
	productName, size string,
	price float64,
) (bool, bool, error) {
	orderNumber := fmt.Sprintf("MATCH-%d", matchID)

	// Notify buyer
	buyerData := map[string]interface{}{
		"match_id":     matchID,
		"product_id":   productID,
		"product_name": productName,
		"size":         size,
		"price":        price,
		"role":         "buyer",
	}

	_, err := s.SendNotification(
		ctx, buyerID,
		model.TypeMatchCreated,
		"🎯 Your bid has been matched!",
		fmt.Sprintf("Your bid for %s (%s) has been matched at $%.2f", productName, size, price),
		buyerData,
		true, true,
	)
	buyerNotified := err == nil

	// Send buyer email
	go func() {
		buyerEmail := fmt.Sprintf("user%d@example.com", buyerID)
		_ = s.emailService.SendMatchCreatedEmail(buyerEmail, "bid", productName, size, price, orderNumber)
	}()

	// Notify seller
	sellerData := map[string]interface{}{
		"match_id":     matchID,
		"product_id":   productID,
		"product_name": productName,
		"size":         size,
		"price":        price,
		"role":         "seller",
	}

	_, err = s.SendNotification(
		ctx, sellerID,
		model.TypeMatchCreated,
		"🎯 Your ask has been matched!",
		fmt.Sprintf("Your ask for %s (%s) has been matched at $%.2f", productName, size, price),
		sellerData,
		true, true,
	)
	sellerNotified := err == nil

	// Send seller email
	go func() {
		sellerEmail := fmt.Sprintf("user%d@example.com", sellerID)
		_ = s.emailService.SendMatchCreatedEmail(sellerEmail, "ask", productName, size, price, orderNumber)
	}()

	return buyerNotified, sellerNotified, nil
}

// NotifyOrderUpdate sends notification when order status changes
func (s *NotificationService) NotifyOrderUpdate(
	ctx context.Context,
	orderID int64,
	orderNumber string,
	buyerID, sellerID int64,
	oldStatus, newStatus string,
	trackingNumber string,
) error {
	data := map[string]interface{}{
		"order_id":        orderID,
		"order_number":    orderNumber,
		"old_status":      oldStatus,
		"new_status":      newStatus,
		"tracking_number": trackingNumber,
	}

	// Determine who to notify and message
	var userID int64
	var title, message string

	switch newStatus {
	case "paid":
		userID = sellerID
		title = "💳 Payment received for your sale"
		message = fmt.Sprintf("Order %s has been paid. Please prepare for shipment.", orderNumber)
	case "shipped":
		userID = buyerID
		title = "📦 Your order has shipped!"
		message = fmt.Sprintf("Order %s is on its way. Tracking: %s", orderNumber, trackingNumber)

		// Send email
		go func() {
			buyerEmail := fmt.Sprintf("user%d@example.com", buyerID)
			_ = s.emailService.SendOrderShippedEmail(buyerEmail, orderNumber, trackingNumber, "USPS")
		}()
	case "delivered":
		// Notify both
		_, _ = s.SendNotification(ctx, buyerID, model.TypeOrderDelivered,
			"✅ Order delivered",
			fmt.Sprintf("Order %s has been delivered!", orderNumber),
			data, true, true)

		userID = sellerID
		title = "✅ Order delivered to buyer"
		message = fmt.Sprintf("Order %s has been delivered. Payout will be processed soon.", orderNumber)
	default:
		userID = buyerID
		title = "📦 Order update"
		message = fmt.Sprintf("Order %s status: %s", orderNumber, newStatus)
	}

	_, err := s.SendNotification(ctx, userID, model.TypeOrderShipped, title, message, data, true, true)
	return err
}

// NotifyPaymentEvent sends notification for payment events
func (s *NotificationService) NotifyPaymentEvent(
	ctx context.Context,
	paymentID, userID, orderID int64,
	eventType string,
	amount float64,
	currency string,
) error {
	data := map[string]interface{}{
		"payment_id": paymentID,
		"order_id":   orderID,
		"amount":     amount,
		"currency":   currency,
	}

	var title, message string
	var notifType string

	switch eventType {
	case "payment_succeeded":
		title = "✅ Payment successful"
		message = fmt.Sprintf("Your payment of $%.2f has been confirmed", amount)
		notifType = model.TypePaymentSucceeded

		// Send email
		go func() {
			userEmail := fmt.Sprintf("user%d@example.com", userID)
			orderNumber := fmt.Sprintf("ORD-%d", orderID)
			_ = s.emailService.SendPaymentReceivedEmail(userEmail, orderNumber, amount)
		}()
	case "payment_failed":
		title = "❌ Payment failed"
		message = "Your payment could not be processed. Please try again."
		notifType = model.TypePaymentFailed
	case "refund_issued":
		title = "💰 Refund issued"
		message = fmt.Sprintf("A refund of $%.2f has been issued to your account", amount)
		notifType = model.TypeRefundIssued
	case "payout_completed":
		title = "💰 Payout completed"
		message = fmt.Sprintf("$%.2f has been transferred to your account", amount)
		notifType = model.TypePayoutCompleted

		// Send email
		go func() {
			userEmail := fmt.Sprintf("user%d@example.com", userID)
			orderNumber := fmt.Sprintf("ORD-%d", orderID)
			_ = s.emailService.SendPayoutCompletedEmail(userEmail, orderNumber, amount)
		}()
	default:
		title = "Payment notification"
		message = fmt.Sprintf("Payment event: %s", eventType)
		notifType = model.TypePaymentSucceeded
	}

	_, err := s.SendNotification(ctx, userID, notifType, title, message, data, true, true)
	return err
}
