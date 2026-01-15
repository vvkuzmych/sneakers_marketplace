package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	productPb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/product"
)

type ProductHandler struct {
	client productPb.ProductServiceClient
}

func NewProductHandler(client productPb.ProductServiceClient) *ProductHandler {
	return &ProductHandler{client: client}
}

// ListProducts godoc
// @Summary List all products
// @Tags products
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} productPb.ListProductsResponse
// @Router /api/v1/products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 32)

	resp, err := h.client.ListProducts(c.Request.Context(), &productPb.ListProductsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetProduct godoc
// @Summary Get product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} productPb.GetProductResponse
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	resp, err := h.client.GetProduct(c.Request.Context(), &productPb.GetProductRequest{
		ProductId: id,
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

// SearchProducts godoc
// @Summary Search products
// @Tags products
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {object} productPb.SearchProductsResponse
// @Router /api/v1/products/search [get]
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search query required"})
		return
	}

	resp, err := h.client.SearchProducts(c.Request.Context(), &productPb.SearchProductsRequest{
		Query: query,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
