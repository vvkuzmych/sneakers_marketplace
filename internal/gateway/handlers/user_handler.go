package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	userPb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/user"
)

type UserHandler struct {
	client userPb.UserServiceClient
}

func NewUserHandler(client userPb.UserServiceClient) *UserHandler {
	return &UserHandler{client: client}
}

// Register godoc
// @Summary Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body userPb.RegisterRequest true "Register Request"
// @Success 200 {object} userPb.RegisterResponse
// @Router /api/v1/auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req userPb.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.Register(c.Request.Context(), &req)
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

// Login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body userPb.LoginRequest true "Login Request"
// @Success 200 {object} userPb.LoginResponse
// @Router /api/v1/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req userPb.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if resp.Error != "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": resp.Error})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetProfile godoc
// @Summary Get user profile
// @Tags users
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} userPb.GetProfileResponse
// @Security BearerAuth
// @Router /api/v1/users/{user_id} [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	resp, err := h.client.GetProfile(c.Request.Context(), &userPb.GetProfileRequest{
		UserId: userID,
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
