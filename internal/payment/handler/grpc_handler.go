package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vvkuzmych/sneakers_marketplace/internal/payment/model"
	"github.com/vvkuzmych/sneakers_marketplace/internal/payment/service"
	pb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/payment"
)

type PaymentHandler struct {
	pb.UnimplementedPaymentServiceServer
	service *service.PaymentService
}

func NewPaymentHandler(service *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service: service,
	}
}

// Helper: convert model.Payment to pb.Payment
func paymentToProto(p *model.Payment) *pb.Payment {
	if p == nil {
		return nil
	}

	payment := &pb.Payment{
		Id:             p.ID,
		PaymentId:      p.PaymentID,
		OrderId:        p.OrderID,
		UserId:         p.UserID,
		Amount:         p.Amount,
		Currency:       p.Currency,
		Status:         p.Status,
		RefundedAmount: p.RefundedAmount,
		CreatedAt:      timestamppb.New(p.CreatedAt),
		UpdatedAt:      timestamppb.New(p.UpdatedAt),
	}

	// Optional fields
	if p.StripePaymentIntentID.Valid {
		payment.StripePaymentIntentId = p.StripePaymentIntentID.String
	}
	if p.StripeChargeID.Valid {
		payment.StripeChargeId = p.StripeChargeID.String
	}
	if p.StripeCustomerID.Valid {
		payment.StripeCustomerId = p.StripeCustomerID.String
	}
	if p.PaymentMethod.Valid {
		payment.PaymentMethod = p.PaymentMethod.String
	}
	if p.CardLast4.Valid {
		payment.CardLast4 = p.CardLast4.String
	}
	if p.CardBrand.Valid {
		payment.CardBrand = p.CardBrand.String
	}
	if p.RefundReason.Valid {
		payment.RefundReason = p.RefundReason.String
	}
	if p.ProcessedAt.Valid {
		payment.ProcessedAt = timestamppb.New(p.ProcessedAt.Time)
	}
	if p.RefundedAt.Valid {
		payment.RefundedAt = timestamppb.New(p.RefundedAt.Time)
	}

	return payment
}

// Helper: convert model.Payout to pb.Payout
func payoutToProto(p *model.Payout) *pb.Payout {
	if p == nil {
		return nil
	}

	payout := &pb.Payout{
		Id:        p.ID,
		PayoutId:  p.PayoutID,
		OrderId:   p.OrderID,
		SellerId:  p.SellerID,
		PaymentId: p.PaymentID,
		Amount:    p.Amount,
		Currency:  p.Currency,
		Status:    p.Status,
		CreatedAt: timestamppb.New(p.CreatedAt),
		UpdatedAt: timestamppb.New(p.UpdatedAt),
	}

	// Optional fields
	if p.StripeTransferID.Valid {
		payout.StripeTransferId = p.StripeTransferID.String
	}
	if p.StripeAccountID.Valid {
		payout.StripeAccountId = p.StripeAccountID.String
	}
	if p.FailureReason.Valid {
		payout.FailureReason = p.FailureReason.String
	}
	if p.ProcessedAt.Valid {
		payout.ProcessedAt = timestamppb.New(p.ProcessedAt.Time)
	}

	return payout
}

// CreatePaymentIntent creates a Stripe PaymentIntent
func (h *PaymentHandler) CreatePaymentIntent(ctx context.Context, req *pb.CreatePaymentIntentRequest) (*pb.CreatePaymentIntentResponse, error) {
	if req.OrderId == 0 {
		return &pb.CreatePaymentIntentResponse{Error: "order_id is required"}, nil
	}
	if req.UserId == 0 {
		return &pb.CreatePaymentIntentResponse{Error: "user_id is required"}, nil
	}
	if req.Amount <= 0 {
		return &pb.CreatePaymentIntentResponse{Error: "amount must be greater than 0"}, nil
	}

	payment, clientSecret, err := h.service.CreatePaymentIntent(
		ctx, req.OrderId, req.UserId,
		req.Amount, req.Currency, req.StripeCustomerId,
	)
	if err != nil {
		return &pb.CreatePaymentIntentResponse{Error: err.Error()}, nil
	}

	return &pb.CreatePaymentIntentResponse{
		Payment:      paymentToProto(payment),
		ClientSecret: clientSecret,
	}, nil
}

// ConfirmPayment confirms a payment
func (h *PaymentHandler) ConfirmPayment(ctx context.Context, req *pb.ConfirmPaymentRequest) (*pb.ConfirmPaymentResponse, error) {
	if req.PaymentId == 0 {
		return &pb.ConfirmPaymentResponse{Error: "payment_id is required"}, nil
	}
	if req.StripePaymentIntentId == "" {
		return &pb.ConfirmPaymentResponse{Error: "stripe_payment_intent_id is required"}, nil
	}

	payment, err := h.service.ConfirmPayment(ctx, req.PaymentId, req.StripePaymentIntentId)
	if err != nil {
		return &pb.ConfirmPaymentResponse{Error: err.Error()}, nil
	}

	return &pb.ConfirmPaymentResponse{
		Payment: paymentToProto(payment),
	}, nil
}

// GetPayment retrieves a payment
func (h *PaymentHandler) GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.GetPaymentResponse, error) {
	if req.PaymentId == 0 {
		return &pb.GetPaymentResponse{Error: "payment_id is required"}, nil
	}

	payment, err := h.service.GetPayment(ctx, req.PaymentId)
	if err != nil {
		return &pb.GetPaymentResponse{Error: err.Error()}, nil
	}

	return &pb.GetPaymentResponse{Payment: paymentToProto(payment)}, nil
}

// ListPayments lists payments
func (h *PaymentHandler) ListPayments(ctx context.Context, req *pb.ListPaymentsRequest) (*pb.ListPaymentsResponse, error) {
	page := req.Page
	if page == 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 10
	}

	payments, total, err := h.service.ListPayments(ctx, req.UserId, req.Status, page, pageSize)
	if err != nil {
		return &pb.ListPaymentsResponse{Error: err.Error()}, nil
	}

	pbPayments := make([]*pb.Payment, len(payments))
	for i, p := range payments {
		pbPayments[i] = paymentToProto(p)
	}

	return &pb.ListPaymentsResponse{
		Payments: pbPayments,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// CreateRefund creates a refund
func (h *PaymentHandler) CreateRefund(ctx context.Context, req *pb.CreateRefundRequest) (*pb.CreateRefundResponse, error) {
	if req.PaymentId == 0 {
		return &pb.CreateRefundResponse{Error: "payment_id is required"}, nil
	}
	if req.Reason == "" {
		return &pb.CreateRefundResponse{Error: "reason is required"}, nil
	}

	stripeRefundID, err := h.service.CreateRefund(ctx, req.PaymentId, req.Amount, req.Reason)
	if err != nil {
		return &pb.CreateRefundResponse{Error: err.Error()}, nil
	}

	payment, _ := h.service.GetPayment(ctx, req.PaymentId)

	return &pb.CreateRefundResponse{
		Payment:        paymentToProto(payment),
		StripeRefundId: stripeRefundID,
	}, nil
}

// GetRefundStatus retrieves refund status
func (h *PaymentHandler) GetRefundStatus(ctx context.Context, req *pb.GetRefundStatusRequest) (*pb.GetRefundStatusResponse, error) {
	if req.PaymentId == 0 {
		return &pb.GetRefundStatusResponse{Error: "payment_id is required"}, nil
	}

	payment, err := h.service.GetRefundStatus(ctx, req.PaymentId)
	if err != nil {
		return &pb.GetRefundStatusResponse{Error: err.Error()}, nil
	}

	response := &pb.GetRefundStatusResponse{
		RefundedAmount: payment.RefundedAmount,
	}
	if payment.RefundReason.Valid {
		response.RefundReason = payment.RefundReason.String
	}
	if payment.RefundedAt.Valid {
		response.RefundedAt = timestamppb.New(payment.RefundedAt.Time)
	}

	return response, nil
}

// HandleStripeWebhook handles Stripe webhook
func (h *PaymentHandler) HandleStripeWebhook(ctx context.Context, req *pb.HandleStripeWebhookRequest) (*pb.HandleStripeWebhookResponse, error) {
	eventType, err := h.service.HandleStripeWebhook(ctx, req.Payload, req.Signature)
	if err != nil {
		return &pb.HandleStripeWebhookResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &pb.HandleStripeWebhookResponse{
		Success:   true,
		EventType: eventType,
	}, nil
}

// CreatePayout creates a payout
func (h *PaymentHandler) CreatePayout(ctx context.Context, req *pb.CreatePayoutRequest) (*pb.CreatePayoutResponse, error) {
	if req.OrderId == 0 {
		return &pb.CreatePayoutResponse{Error: "order_id is required"}, nil
	}
	if req.SellerId == 0 {
		return &pb.CreatePayoutResponse{Error: "seller_id is required"}, nil
	}
	if req.PaymentId == 0 {
		return &pb.CreatePayoutResponse{Error: "payment_id is required"}, nil
	}
	if req.Amount <= 0 {
		return &pb.CreatePayoutResponse{Error: "amount must be greater than 0"}, nil
	}
	if req.StripeAccountId == "" {
		return &pb.CreatePayoutResponse{Error: "stripe_account_id is required"}, nil
	}

	payout, err := h.service.CreatePayout(
		ctx, req.OrderId, req.SellerId, req.PaymentId,
		req.Amount, req.StripeAccountId,
	)
	if err != nil {
		return &pb.CreatePayoutResponse{Error: err.Error()}, nil
	}

	return &pb.CreatePayoutResponse{
		Payout: payoutToProto(payout),
	}, nil
}

// GetPayout retrieves a payout
func (h *PaymentHandler) GetPayout(ctx context.Context, req *pb.GetPayoutRequest) (*pb.GetPayoutResponse, error) {
	if req.PayoutId == 0 {
		return &pb.GetPayoutResponse{Error: "payout_id is required"}, nil
	}

	payout, err := h.service.GetPayout(ctx, req.PayoutId)
	if err != nil {
		return &pb.GetPayoutResponse{Error: err.Error()}, nil
	}

	return &pb.GetPayoutResponse{Payout: payoutToProto(payout)}, nil
}

// ListPayouts lists payouts
func (h *PaymentHandler) ListPayouts(ctx context.Context, req *pb.ListPayoutsRequest) (*pb.ListPayoutsResponse, error) {
	page := req.Page
	if page == 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 10
	}

	payouts, total, err := h.service.ListPayouts(ctx, req.SellerId, req.Status, page, pageSize)
	if err != nil {
		return &pb.ListPayoutsResponse{Error: err.Error()}, nil
	}

	pbPayouts := make([]*pb.Payout, len(payouts))
	for i, p := range payouts {
		pbPayouts[i] = payoutToProto(p)
	}

	return &pb.ListPayoutsResponse{
		Payouts:  pbPayouts,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
