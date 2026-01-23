package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vvkuzmych/sneakers_marketplace/internal/fees/repository"
	"github.com/vvkuzmych/sneakers_marketplace/internal/fees/service"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/logger"
)

type FeeHandler struct {
	service *service.FeeService
	log     *logger.Logger
}

// NewFeeHandler creates a new fee handler with optional subscription provider
// If subscriptionProvider is nil, uses default fees (Free tier: 1%)
func NewFeeHandler(feeRepo *repository.FeeRepository, log *logger.Logger, subscriptionProvider service.SubscriptionFeeProvider) *FeeHandler {
	feeService := service.NewFeeService(feeRepo, log, subscriptionProvider)
	return &FeeHandler{
		service: feeService,
		log:     log,
	}
}

// CalculateFees calculates fee breakdown for a given price based on seller's subscription
// GET /api/v1/fees/calculate?vertical=sneakers&price=200&seller_user_id=123
// If seller_user_id is not provided, uses default Free tier fees (1%)
func (h *FeeHandler) CalculateFees(c *gin.Context) {
	vertical := c.Query("vertical")
	if vertical == "" {
		vertical = "sneakers" // default
	}

	priceStr := c.Query("price")
	if priceStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price parameter is required"})
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price, must be a positive number"})
		return
	}

	// Get seller_user_id (optional - if not provided, use 0 which triggers default fees)
	sellerUserIDStr := c.Query("seller_user_id")
	var sellerUserID int64
	if sellerUserIDStr != "" {
		sellerUserID, err = strconv.ParseInt(sellerUserIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid seller_user_id format"})
			return
		}
	}

	// If sellerUserID is 0, use special value to indicate "use default fees"
	if sellerUserID == 0 {
		sellerUserID = -1 // Special value handled by fee service
	}

	// Calculate fees
	breakdown, err := h.service.CalculateFees(c.Request.Context(), vertical, price, sellerUserID)
	if err != nil {
		h.log.Errorf("Failed to calculate fees: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate fees"})
		return
	}

	c.JSON(http.StatusOK, breakdown)
}

// GetFeeConfig retrieves fee configuration for a vertical
// GET /api/v1/fees/config/:vertical
func (h *FeeHandler) GetFeeConfig(c *gin.Context) {
	vertical := c.Param("vertical")
	if vertical == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vertical parameter is required"})
		return
	}

	config, err := h.service.GetFeeConfig(c.Request.Context(), vertical)
	if err != nil {
		h.log.Errorf("Failed to get fee config: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "fee config not found"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// GetAllFeeConfigs retrieves all fee configurations
// GET /api/v1/fees/configs
func (h *FeeHandler) GetAllFeeConfigs(c *gin.Context) {
	configs, err := h.service.GetAllFeeConfigs(c.Request.Context())
	if err != nil {
		h.log.Errorf("Failed to get fee configs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get fee configs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"configs": configs,
	})
}

// GetRevenue retrieves platform revenue for a date range
// GET /api/v1/fees/revenue?start_date=2026-01-01&end_date=2026-01-31
func (h *FeeHandler) GetRevenue(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, use YYYY-MM-DD"})
			return
		}
	} else {
		// Default to start of current month
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, use YYYY-MM-DD"})
			return
		}
		// Add one day to include the end date
		endDate = endDate.Add(24 * time.Hour)
	} else {
		// Default to end of current month
		endDate = time.Now()
	}

	// Get total revenue
	totalRevenue, err := h.service.GetFeeConfig(c.Request.Context(), "sneakers") // TODO: Fix this to use actual revenue method
	if err != nil {
		h.log.Errorf("Failed to get revenue: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get revenue"})
		return
	}

	// For now, return mock data structure
	// TODO: Implement actual revenue calculation from repository
	c.JSON(http.StatusOK, gin.H{
		"start_date":    startDate.Format("2006-01-02"),
		"end_date":      endDate.Format("2006-01-02"),
		"total_revenue": 0.0, // TODO: Get from service
		"by_vertical": map[string]float64{
			"sneakers": 0.0,
			"tickets":  0.0,
		},
		"transaction_count": 0,
		"config":            totalRevenue, // Temp, remove later
	})
}

// GetTransactionFee retrieves transaction fee by match ID
// GET /api/v1/fees/transaction/:match_id
func (h *FeeHandler) GetTransactionFee(c *gin.Context) {
	matchIDStr := c.Param("match_id")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match_id"})
		return
	}

	fee, err := h.service.GetTransactionFee(context.Background(), matchID)
	if err != nil {
		h.log.Errorf("Failed to get transaction fee: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction fee not found"})
		return
	}

	c.JSON(http.StatusOK, fee)
}
