package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	paymentPb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/payment"
)

type PaymentHandler struct {
	client paymentPb.PaymentServiceClient
}

func NewPaymentHandler(client paymentPb.PaymentServiceClient) *PaymentHandler {
	return &PaymentHandler{client: client}
}

// CreatePaymentIntent godoc
// @Summary Create payment intent for order
// @Tags payments
// @Accept json
// @Produce json
// @Param request body paymentPb.CreatePaymentIntentRequest true "Payment Intent Request"
// @Success 200 {object} paymentPb.CreatePaymentIntentResponse
// @Security BearerAuth
// @Router /api/v1/payments/intent [post]
func (h *PaymentHandler) CreatePaymentIntent(c *gin.Context) {
	var req paymentPb.CreatePaymentIntentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.CreatePaymentIntent(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if resp.Error != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp.Error})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetPayment godoc
// @Summary Get payment by ID
// @Tags payments
// @Produce json
// @Param id path int true "Payment ID"
// @Success 200 {object} paymentPb.GetPaymentResponse
// @Security BearerAuth
// @Router /api/v1/payments/{id} [get]
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment id"})
		return
	}

	resp, err := h.client.GetPayment(c.Request.Context(), &paymentPb.GetPaymentRequest{
		PaymentId: id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if resp.Error != "" {
		c.JSON(http.StatusNotFound, gin.H{"error": resp.Error})
		return
	}

	c.JSON(http.StatusOK, resp)
}
