package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vvkuzmych/sneakers_marketplace/internal/user/model"
	"github.com/vvkuzmych/sneakers_marketplace/internal/user/service"
	pb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/user"
)

// UserHandler implements the gRPC UserService server
type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register handles user registration
func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Validate request
	if req.Email == "" || req.Password == "" {
		return &pb.RegisterResponse{
			Error: "email and password are required",
		}, nil
	}

	// Register user
	user, accessToken, refreshToken, err := h.userService.Register(
		ctx,
		req.Email,
		req.Password,
		req.FirstName,
		req.LastName,
		req.Phone,
	)

	if err != nil {
		return &pb.RegisterResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.RegisterResponse{
		User:         modelUserToProto(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Login handles user login
func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Validate request
	if req.Email == "" || req.Password == "" {
		return &pb.LoginResponse{
			Error: "email and password are required",
		}, nil
	}

	// Login user
	user, accessToken, refreshToken, err := h.userService.Login(
		ctx,
		req.Email,
		req.Password,
	)

	if err != nil {
		return &pb.LoginResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.LoginResponse{
		User:         modelUserToProto(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken handles token refresh
func (h *UserHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return &pb.RefreshTokenResponse{
			Error: "refresh token is required",
		}, nil
	}

	accessToken, refreshToken, err := h.userService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return &pb.RefreshTokenResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Logout handles user logout
func (h *UserHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	if req.AccessToken == "" {
		return &pb.LogoutResponse{
			Error: "access token is required",
		}, nil
	}

	if err := h.userService.Logout(ctx, req.AccessToken); err != nil {
		return &pb.LogoutResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &pb.LogoutResponse{
		Success: true,
	}, nil
}

// GetProfile retrieves user profile
func (h *UserHandler) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	user, err := h.userService.GetProfile(ctx, req.UserId)
	if err != nil {
		return &pb.GetProfileResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.GetProfileResponse{
		User: modelUserToProto(user),
	}, nil
}

// UpdateProfile updates user profile
func (h *UserHandler) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	user, err := h.userService.UpdateProfile(
		ctx,
		req.UserId,
		req.FirstName,
		req.LastName,
		req.Phone,
	)

	if err != nil {
		return &pb.UpdateProfileResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.UpdateProfileResponse{
		User: modelUserToProto(user),
	}, nil
}

// AddAddress adds a new address
func (h *UserHandler) AddAddress(ctx context.Context, req *pb.AddAddressRequest) (*pb.AddAddressResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	address := &model.Address{
		UserID:      req.UserId,
		AddressType: req.AddressType,
		StreetLine1: req.StreetLine1,
		StreetLine2: req.StreetLine2,
		City:        req.City,
		State:       req.State,
		PostalCode:  req.PostalCode,
		Country:     req.Country,
		IsDefault:   req.IsDefault,
	}

	result, err := h.userService.AddAddress(ctx, address)
	if err != nil {
		return &pb.AddAddressResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.AddAddressResponse{
		Address: modelAddressToProto(result),
	}, nil
}

// GetAddresses retrieves all addresses for a user
func (h *UserHandler) GetAddresses(ctx context.Context, req *pb.GetAddressesRequest) (*pb.GetAddressesResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	addresses, err := h.userService.GetAddresses(ctx, req.UserId)
	if err != nil {
		return &pb.GetAddressesResponse{
			Error: err.Error(),
		}, nil
	}

	protoAddresses := make([]*pb.Address, len(addresses))
	for i, addr := range addresses {
		protoAddresses[i] = modelAddressToProto(addr)
	}

	return &pb.GetAddressesResponse{
		Addresses: protoAddresses,
	}, nil
}

// UpdateAddress updates an address
func (h *UserHandler) UpdateAddress(ctx context.Context, req *pb.UpdateAddressRequest) (*pb.UpdateAddressResponse, error) {
	if req.AddressId == 0 || req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "address_id and user_id are required")
	}

	address := &model.Address{
		ID:          req.AddressId,
		UserID:      req.UserId,
		AddressType: req.AddressType,
		StreetLine1: req.StreetLine1,
		StreetLine2: req.StreetLine2,
		City:        req.City,
		State:       req.State,
		PostalCode:  req.PostalCode,
		Country:     req.Country,
		IsDefault:   req.IsDefault,
	}

	result, err := h.userService.UpdateAddress(ctx, address)
	if err != nil {
		return &pb.UpdateAddressResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.UpdateAddressResponse{
		Address: modelAddressToProto(result),
	}, nil
}

// DeleteAddress deletes an address
func (h *UserHandler) DeleteAddress(ctx context.Context, req *pb.DeleteAddressRequest) (*pb.DeleteAddressResponse, error) {
	if req.AddressId == 0 || req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "address_id and user_id are required")
	}

	if err := h.userService.DeleteAddress(ctx, req.AddressId, req.UserId); err != nil {
		return &pb.DeleteAddressResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &pb.DeleteAddressResponse{
		Success: true,
	}, nil
}

// Helper functions to convert between model and proto

func modelUserToProto(user *model.User) *pb.User {
	return &pb.User{
		Id:         user.ID,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Phone:      user.Phone,
		IsVerified: user.IsVerified,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func modelAddressToProto(address *model.Address) *pb.Address {
	return &pb.Address{
		Id:          address.ID,
		UserId:      address.UserID,
		AddressType: address.AddressType,
		StreetLine1: address.StreetLine1,
		StreetLine2: address.StreetLine2,
		City:        address.City,
		State:       address.State,
		PostalCode:  address.PostalCode,
		Country:     address.Country,
		IsDefault:   address.IsDefault,
		CreatedAt:   address.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
