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
	var req biddingPb.PlaceBidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.PlaceBid(c.Request.Context(), &req)
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
	var req biddingPb.PlaceAskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.PlaceAsk(c.Request.Context(), &req)
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
