package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	biddingPb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/bidding"
)

type BiddingHandler struct {
	client biddingPb.BiddingServiceClient
}

func NewBiddingHandler(client biddingPb.BiddingServiceClient) *BiddingHandler {
	return &BiddingHandler{client: client}
}

// PlaceBid godoc
// @Summary Place a bid (buy order)
// @Tags bidding
// @Accept json
// @Produce json
// @Param request body biddingPb.PlaceBidRequest true "Bid Request"
// @Success 200 {object} biddingPb.PlaceBidResponse
// @Security BearerAuth
// @Router /api/v1/bids [post]
func (h *BiddingHandler) PlaceBid(c *gin.Context) {
	// Parse JSON body manually to handle camelCase
	var body struct {
		ProductID      int64   `json:"productId"`
		SizeID         int64   `json:"sizeId"`
		Price          float64 `json:"price"`
		Quantity       int32   `json:"quantity"`
		ExpiresInHours int32   `json:"expiresInHours"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract user_id from JWT claims (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}

	// Convert user_id to int64
	var userIDInt64 int64
	switch v := userID.(type) {
	case float64:
		userIDInt64 = int64(v)
	case int64:
		userIDInt64 = v
	case int:
		userIDInt64 = int64(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id type"})
		return
	}

	// Build gRPC request with snake_case fields
	req := &biddingPb.PlaceBidRequest{
		UserId:         userIDInt64,
		ProductId:      body.ProductID,
		SizeId:         body.SizeID,
		Price:          body.Price,
		Quantity:       body.Quantity,
		ExpiresInHours: body.ExpiresInHours,
	}

	resp, err := h.client.PlaceBid(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PlaceAsk godoc
// @Summary Place an ask (sell order)
// @Tags bidding
// @Accept json
// @Produce json
// @Param request body biddingPb.PlaceAskRequest true "Ask Request"
// @Success 200 {object} biddingPb.PlaceAskResponse
// @Security BearerAuth
// @Router /api/v1/asks [post]
func (h *BiddingHandler) PlaceAsk(c *gin.Context) {
	// Parse JSON body manually to handle camelCase
	var body struct {
		ProductID      int64   `json:"productId"`
		SizeID         int64   `json:"sizeId"`
		Price          float64 `json:"price"`
		Quantity       int32   `json:"quantity"`
		ExpiresInHours int32   `json:"expiresInHours"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract user_id from JWT claims (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}

	// Convert user_id to int64
	var userIDInt64 int64
	switch v := userID.(type) {
	case float64:
		userIDInt64 = int64(v)
	case int64:
		userIDInt64 = v
	case int:
		userIDInt64 = int64(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id type"})
		return
	}

	// Build gRPC request with snake_case fields
	req := &biddingPb.PlaceAskRequest{
		UserId:         userIDInt64,
		ProductId:      body.ProductID,
		SizeId:         body.SizeID,
		Price:          body.Price,
		Quantity:       body.Quantity,
		ExpiresInHours: body.ExpiresInHours,
	}

	resp, err := h.client.PlaceAsk(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetMarketPrice godoc
// @Summary Get market price for product/size
// @Tags bidding
// @Produce json
// @Param product_id path int true "Product ID"
// @Param size_id path int true "Size ID"
// @Success 200 {object} biddingPb.GetMarketPriceResponse
// @Router /api/v1/market/{product_id}/{size_id} [get]
func (h *BiddingHandler) GetMarketPrice(c *gin.Context) {
	productID, _ := strconv.ParseInt(c.Param("product_id"), 10, 64)
	sizeID, _ := strconv.ParseInt(c.Param("size_id"), 10, 64)

	resp, err := h.client.GetMarketPrice(c.Request.Context(), &biddingPb.GetMarketPriceRequest{
		ProductId: productID,
		SizeId:    sizeID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetProductBids godoc
// @Summary Get all bids for a product
// @Tags bidding
// @Produce json
// @Param product_id path int true "Product ID"
// @Param size_id query int false "Size ID"
// @Success 200 {object} biddingPb.GetProductBidsResponse
// @Router /api/v1/bids/product/{product_id} [get]
func (h *BiddingHandler) GetProductBids(c *gin.Context) {
	productID, _ := strconv.ParseInt(c.Param("product_id"), 10, 64)
	sizeID, _ := strconv.ParseInt(c.Query("size_id"), 10, 64)

	resp, err := h.client.GetProductBids(c.Request.Context(), &biddingPb.GetProductBidsRequest{
		ProductId: productID,
		SizeId:    sizeID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetProductAsks godoc
// @Summary Get all asks for a product
// @Tags bidding
// @Produce json
// @Param product_id path int true "Product ID"
// @Param size_id query int false "Size ID"
// @Success 200 {object} biddingPb.GetProductAsksResponse
// @Router /api/v1/asks/product/{product_id} [get]
func (h *BiddingHandler) GetProductAsks(c *gin.Context) {
	productID, _ := strconv.ParseInt(c.Param("product_id"), 10, 64)
	sizeID, _ := strconv.ParseInt(c.Query("size_id"), 10, 64)

	resp, err := h.client.GetProductAsks(c.Request.Context(), &biddingPb.GetProductAsksRequest{
		ProductId: productID,
		SizeId:    sizeID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
