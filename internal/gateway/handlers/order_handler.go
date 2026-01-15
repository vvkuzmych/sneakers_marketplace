package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	orderPb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/order"
)

type OrderHandler struct {
	client orderPb.OrderServiceClient
}

func NewOrderHandler(client orderPb.OrderServiceClient) *OrderHandler {
	return &OrderHandler{client: client}
}

// GetOrder godoc
// @Summary Get order by ID
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} orderPb.GetOrderResponse
// @Security BearerAuth
// @Router /api/v1/orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	resp, err := h.client.GetOrder(c.Request.Context(), &orderPb.GetOrderRequest{
		OrderId: id,
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

// ListBuyerOrders godoc
// @Summary List orders for buyer
// @Tags orders
// @Produce json
// @Param buyer_id path int true "Buyer ID"
// @Success 200 {object} orderPb.GetBuyerOrdersResponse
// @Security BearerAuth
// @Router /api/v1/orders/buyer/{buyer_id} [get]
func (h *OrderHandler) ListBuyerOrders(c *gin.Context) {
	buyerID, _ := strconv.ParseInt(c.Param("buyer_id"), 10, 64)

	resp, err := h.client.GetBuyerOrders(c.Request.Context(), &orderPb.GetBuyerOrdersRequest{
		BuyerId: buyerID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
