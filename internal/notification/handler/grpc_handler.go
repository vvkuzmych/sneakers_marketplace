package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vvkuzmych/sneakers_marketplace/internal/notification/model"
	"github.com/vvkuzmych/sneakers_marketplace/internal/notification/service"
	pb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/notification"
)

type NotificationHandler struct {
	pb.UnimplementedNotificationServiceServer
	service *service.NotificationService
}

func NewNotificationHandler(service *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// SendNotification implements the gRPC method
func (h *NotificationHandler) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	data := make(map[string]interface{})
	if req.Data != "" {
		data["raw"] = req.Data
	}

	notif, err := h.service.SendNotification(
		ctx, req.UserId, req.Type, req.Title, req.Message,
		data, req.SendEmail, req.SendPush,
	)
	if err != nil {
		return &pb.SendNotificationResponse{Error: err.Error()}, nil
	}

	return &pb.SendNotificationResponse{
		Notification: modelToProto(notif),
	}, nil
}

// GetNotifications implements the gRPC method
func (h *NotificationHandler) GetNotifications(ctx context.Context, req *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	notifs, total, err := h.service.GetNotifications(ctx, req.UserId, req.UnreadOnly, req.Page, req.PageSize)
	if err != nil {
		return &pb.GetNotificationsResponse{Error: err.Error()}, nil
	}

	pbNotifs := make([]*pb.Notification, len(notifs))
	for i, n := range notifs {
		pbNotifs[i] = modelToProto(n)
	}

	return &pb.GetNotificationsResponse{
		Notifications: pbNotifs,
		Total:         total,
		Page:          req.Page,
		PageSize:      req.PageSize,
	}, nil
}

// GetUnreadCount implements the gRPC method
func (h *NotificationHandler) GetUnreadCount(ctx context.Context, req *pb.GetUnreadCountRequest) (*pb.GetUnreadCountResponse, error) {
	count, err := h.service.GetUnreadCount(ctx, req.UserId)
	if err != nil {
		return &pb.GetUnreadCountResponse{Error: err.Error()}, nil
	}

	return &pb.GetUnreadCountResponse{Count: count}, nil
}

// MarkAsRead implements the gRPC method
func (h *NotificationHandler) MarkAsRead(ctx context.Context, req *pb.MarkAsReadRequest) (*pb.MarkAsReadResponse, error) {
	err := h.service.MarkAsRead(ctx, req.NotificationId, req.UserId)
	if err != nil {
		return &pb.MarkAsReadResponse{Success: false, Error: err.Error()}, nil
	}

	return &pb.MarkAsReadResponse{Success: true}, nil
}

// MarkAllAsRead implements the gRPC method
func (h *NotificationHandler) MarkAllAsRead(ctx context.Context, req *pb.MarkAllAsReadRequest) (*pb.MarkAllAsReadResponse, error) {
	count, err := h.service.MarkAllAsRead(ctx, req.UserId)
	if err != nil {
		return &pb.MarkAllAsReadResponse{Error: err.Error()}, nil
	}

	return &pb.MarkAllAsReadResponse{UpdatedCount: count}, nil
}

// DeleteNotification implements the gRPC method
func (h *NotificationHandler) DeleteNotification(ctx context.Context, req *pb.DeleteNotificationRequest) (*pb.DeleteNotificationResponse, error) {
	err := h.service.DeleteNotification(ctx, req.NotificationId, req.UserId)
	if err != nil {
		return &pb.DeleteNotificationResponse{Success: false, Error: err.Error()}, nil
	}

	return &pb.DeleteNotificationResponse{Success: true}, nil
}

// GetPreferences implements the gRPC method
func (h *NotificationHandler) GetPreferences(ctx context.Context, req *pb.GetPreferencesRequest) (*pb.GetPreferencesResponse, error) {
	prefs, err := h.service.GetPreferences(ctx, req.UserId)
	if err != nil {
		return &pb.GetPreferencesResponse{Error: err.Error()}, nil
	}

	return &pb.GetPreferencesResponse{
		Preferences: prefsToProto(prefs),
	}, nil
}

// UpdatePreferences implements the gRPC method
func (h *NotificationHandler) UpdatePreferences(ctx context.Context, req *pb.UpdatePreferencesRequest) (*pb.UpdatePreferencesResponse, error) {
	// Get current preferences
	prefs, err := h.service.GetPreferences(ctx, req.UserId)
	if err != nil {
		return &pb.UpdatePreferencesResponse{Error: err.Error()}, nil
	}

	// Update only provided fields
	if req.EmailEnabled != nil {
		prefs.EmailEnabled = *req.EmailEnabled
	}
	if req.EmailMatchCreated != nil {
		prefs.EmailMatchCreated = *req.EmailMatchCreated
	}
	if req.PushEnabled != nil {
		prefs.PushEnabled = *req.PushEnabled
	}
	// ... (other fields can be added similarly)

	updated, err := h.service.UpdatePreferences(ctx, prefs)
	if err != nil {
		return &pb.UpdatePreferencesResponse{Error: err.Error()}, nil
	}

	return &pb.UpdatePreferencesResponse{
		Preferences: prefsToProto(updated),
	}, nil
}

// NotifyMatchCreated implements the gRPC method
func (h *NotificationHandler) NotifyMatchCreated(ctx context.Context, req *pb.NotifyMatchCreatedRequest) (*pb.NotifyMatchCreatedResponse, error) {
	buyerNotified, sellerNotified, err := h.service.NotifyMatchCreated(
		ctx, req.MatchId, req.BuyerId, req.SellerId, req.ProductId,
		req.ProductName, req.Size, req.Price,
	)

	if err != nil {
		return &pb.NotifyMatchCreatedResponse{Error: err.Error()}, nil
	}

	return &pb.NotifyMatchCreatedResponse{
		BuyerNotified:  buyerNotified,
		SellerNotified: sellerNotified,
	}, nil
}

// NotifyOrderUpdate implements the gRPC method
func (h *NotificationHandler) NotifyOrderUpdate(ctx context.Context, req *pb.NotifyOrderUpdateRequest) (*pb.NotifyOrderUpdateResponse, error) {
	err := h.service.NotifyOrderUpdate(
		ctx, req.OrderId, req.OrderNumber, req.BuyerId, req.SellerId,
		req.OldStatus, req.NewStatus, req.TrackingNumber,
	)

	if err != nil {
		return &pb.NotifyOrderUpdateResponse{Success: false, Error: err.Error()}, nil
	}

	return &pb.NotifyOrderUpdateResponse{Success: true}, nil
}

// NotifyPaymentEvent implements the gRPC method
func (h *NotificationHandler) NotifyPaymentEvent(ctx context.Context, req *pb.NotifyPaymentEventRequest) (*pb.NotifyPaymentEventResponse, error) {
	err := h.service.NotifyPaymentEvent(
		ctx, req.PaymentId, req.UserId, req.OrderId,
		req.EventType, req.Amount, req.Currency,
	)

	if err != nil {
		return &pb.NotifyPaymentEventResponse{Success: false, Error: err.Error()}, nil
	}

	return &pb.NotifyPaymentEventResponse{Success: true}, nil
}

// Helper functions
func modelToProto(n *model.Notification) *pb.Notification {
	pbNotif := &pb.Notification{
		Id:        n.ID,
		UserId:    n.UserID,
		Type:      n.Type,
		Title:     n.Title,
		Message:   n.Message,
		Data:      string(n.Data),
		EmailSent: n.EmailSent,
		PushSent:  n.PushSent,
		IsRead:    n.IsRead,
		CreatedAt: timestamppb.New(n.CreatedAt),
	}

	if n.EmailSentAt.Valid {
		pbNotif.EmailSentAt = timestamppb.New(n.EmailSentAt.Time)
	}
	if n.PushSentAt.Valid {
		pbNotif.PushSentAt = timestamppb.New(n.PushSentAt.Time)
	}
	if n.ReadAt.Valid {
		pbNotif.ReadAt = timestamppb.New(n.ReadAt.Time)
	}

	return pbNotif
}

func prefsToProto(p *model.NotificationPreferences) *pb.NotificationPreferences {
	return &pb.NotificationPreferences{
		UserId:               p.UserID,
		EmailEnabled:         p.EmailEnabled,
		EmailMatchCreated:    p.EmailMatchCreated,
		EmailOrderCreated:    p.EmailOrderCreated,
		EmailOrderShipped:    p.EmailOrderShipped,
		EmailPaymentReceived: p.EmailPaymentReceived,
		EmailPayoutCompleted: p.EmailPayoutCompleted,
		PushEnabled:          p.PushEnabled,
		PushMatchCreated:     p.PushMatchCreated,
		PushOrderUpdates:     p.PushOrderUpdates,
		PushPaymentUpdates:   p.PushPaymentUpdates,
		InappEnabled:         p.InAppEnabled,
		CreatedAt:            timestamppb.New(p.CreatedAt),
		UpdatedAt:            timestamppb.New(p.UpdatedAt),
	}
}
