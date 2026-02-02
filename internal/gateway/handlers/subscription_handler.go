package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	subscriptionService "github.com/vvkuzmych/sneakers_marketplace/internal/subscription/service"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/logger"
)

type SubscriptionHandler struct {
	service *subscriptionService.SubscriptionService
	log     *logger.Logger
}

func NewSubscriptionHandler(service *subscriptionService.SubscriptionService, log *logger.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
		log:     log,
	}
}

// GetSubscriptionPlans returns all available subscription plans
// GET /api/v1/subscriptions/plans
func (h *SubscriptionHandler) GetSubscriptionPlans(c *gin.Context) {
	plans, err := h.service.GetAllPlans(c.Request.Context())
	if err != nil {
		h.log.Errorf("Failed to get subscription plans: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get subscription plans"})
		return
	}

	c.JSON(http.StatusOK, plans)
}

// GetCurrentSubscription returns user's current active subscription
// GET /api/v1/subscriptions/current
func (h *SubscriptionHandler) GetCurrentSubscription(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	subscription, err := h.service.GetUserSubscription(c.Request.Context(), userID.(int64))
	if err != nil {
		// User might not have a subscription yet (Free tier)
		h.log.Infof("User %d has no subscription, returning Free tier", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "no active subscription found"})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// Subscribe creates a new subscription for the user
// POST /api/v1/subscriptions/subscribe
func (h *SubscriptionHandler) Subscribe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		PlanID          int64  `json:"planId" binding:"required"`
		BillingCycle    string `json:"billingCycle" binding:"required,oneof=monthly yearly"`
		PaymentMethodID string `json:"paymentMethodId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement actual Stripe subscription creation
	// For now, return a mock response
	h.log.Infof("User %d subscribing to plan %d (%s) with payment method %s",
		userID, req.PlanID, req.BillingCycle, req.PaymentMethodID)

	c.JSON(http.StatusOK, gin.H{
		"message":        "Subscription created successfully",
		"requiresAction": false,
	})
}

// CancelSubscription cancels user's subscription
// POST /api/v1/subscriptions/cancel
func (h *SubscriptionHandler) CancelSubscription(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		SubscriptionID    int64 `json:"subscriptionId" binding:"required"`
		CancelAtPeriodEnd bool  `json:"cancelAtPeriodEnd"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Infof("User %d canceling subscription %d (cancelAtPeriodEnd: %v)",
		userID, req.SubscriptionID, req.CancelAtPeriodEnd)

	// TODO: Implement actual subscription cancellation
	c.JSON(http.StatusOK, gin.H{"message": "Subscription canceled successfully"})
}

// UpdateSubscription updates user's subscription (upgrade/downgrade)
// PUT /api/v1/subscriptions/update
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		SubscriptionID  int64  `json:"subscriptionId" binding:"required"`
		NewPlanID       int64  `json:"newPlanId" binding:"required"`
		NewBillingCycle string `json:"newBillingCycle,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Infof("User %d updating subscription %d to plan %d",
		userID, req.SubscriptionID, req.NewPlanID)

	// TODO: Implement actual subscription update
	c.JSON(http.StatusOK, gin.H{"message": "Subscription updated successfully"})
}

// GetTransactions returns user's subscription transaction history
// GET /api/v1/subscriptions/transactions
func (h *SubscriptionHandler) GetTransactions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	h.log.Infof("Fetching transactions for user %d", userID)

	// TODO: Implement actual transaction history retrieval
	// For now, return empty array
	c.JSON(http.StatusOK, []interface{}{})
}

// CalculateSavings calculates potential fee savings for a plan
// GET /api/v1/subscriptions/savings?plan_id=2&sale_price=1000
func (h *SubscriptionHandler) CalculateSavings(c *gin.Context) {
	// TODO: Implement savings calculation
	c.JSON(http.StatusOK, gin.H{
		"currentPlan":       "Free",
		"currentFeePercent": 1.0,
		"targetPlan":        "Pro",
		"targetFeePercent":  0.75,
		"salePrice":         1000.0,
		"currentFee":        10.0,
		"targetFee":         7.5,
		"savings":           2.5,
		"savingsPercent":    25.0,
	})
}
