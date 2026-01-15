package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vvkuzmych/sneakers_marketplace/internal/product/model"
	"github.com/vvkuzmych/sneakers_marketplace/internal/product/service"
	pb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/product"
)

// ProductHandler implements the gRPC ProductService server
type ProductHandler struct {
	pb.UnimplementedProductServiceServer
	productService *service.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct handles product creation
func (h *ProductHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	if req.Sku == "" || req.Name == "" {
		return &pb.CreateProductResponse{
			Error: "SKU and name are required",
		}, nil
	}

	product := &model.Product{
		SKU:         req.Sku,
		Name:        req.Name,
		Brand:       req.Brand,
		Model:       req.Model,
		Color:       req.Color,
		Description: req.Description,
		Category:    req.Category,
		ReleaseYear: int(req.ReleaseYear),
		RetailPrice: req.RetailPrice,
		IsActive:    true,
	}

	result, err := h.productService.CreateProduct(ctx, product)
	if err != nil {
		return &pb.CreateProductResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.CreateProductResponse{
		Product: modelProductToProto(result, nil, nil),
	}, nil
}

// GetProduct retrieves a product
func (h *ProductHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	if req.ProductId == 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	product, err := h.productService.GetProduct(ctx, req.ProductId)
	if err != nil {
		return &pb.GetProductResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.GetProductResponse{
		Product: modelProductWithDetailsToProto(product),
	}, nil
}

// ListProducts retrieves products with pagination
func (h *ProductHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	page := int(req.Page)
	pageSize := int(req.PageSize)

	products, total, err := h.productService.ListProducts(
		ctx,
		page,
		pageSize,
		req.Category,
		req.Brand,
		req.ActiveOnly,
	)

	if err != nil {
		return &pb.ListProductsResponse{
			Error: err.Error(),
		}, nil
	}

	protoProducts := make([]*pb.Product, len(products))
	for i, p := range products {
		protoProducts[i] = modelProductToProto(p, nil, nil)
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// UpdateProduct updates a product
func (h *ProductHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	if req.ProductId == 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	product := &model.Product{
		ID:          req.ProductId,
		Name:        req.Name,
		Brand:       req.Brand,
		Model:       req.Model,
		Color:       req.Color,
		Description: req.Description,
		Category:    req.Category,
		ReleaseYear: int(req.ReleaseYear),
		RetailPrice: req.RetailPrice,
		IsActive:    req.IsActive,
	}

	result, err := h.productService.UpdateProduct(ctx, product)
	if err != nil {
		return &pb.UpdateProductResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.UpdateProductResponse{
		Product: modelProductToProto(result, nil, nil),
	}, nil
}

// DeleteProduct deletes a product
func (h *ProductHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	if req.ProductId == 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	if err := h.productService.DeleteProduct(ctx, req.ProductId); err != nil {
		return &pb.DeleteProductResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &pb.DeleteProductResponse{
		Success: true,
	}, nil
}

// SearchProducts searches products
func (h *ProductHandler) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
	if req.Query == "" {
		return &pb.SearchProductsResponse{
			Error: "query is required",
		}, nil
	}

	products, total, err := h.productService.SearchProducts(
		ctx,
		req.Query,
		int(req.Page),
		int(req.PageSize),
	)

	if err != nil {
		return &pb.SearchProductsResponse{
			Error: err.Error(),
		}, nil
	}

	protoProducts := make([]*pb.Product, len(products))
	for i, p := range products {
		protoProducts[i] = modelProductToProto(p, nil, nil)
	}

	return &pb.SearchProductsResponse{
		Products: protoProducts,
		Total:    total,
	}, nil
}

// AddProductImage adds an image to a product
func (h *ProductHandler) AddProductImage(ctx context.Context, req *pb.AddProductImageRequest) (*pb.AddProductImageResponse, error) {
	if req.ProductId == 0 || req.ImageUrl == "" {
		return &pb.AddProductImageResponse{
			Error: "product_id and image_url are required",
		}, nil
	}

	image := &model.ProductImage{
		ProductID:    req.ProductId,
		ImageURL:     req.ImageUrl,
		AltText:      req.AltText,
		DisplayOrder: int(req.DisplayOrder),
		IsPrimary:    req.IsPrimary,
	}

	result, err := h.productService.AddProductImage(ctx, image)
	if err != nil {
		return &pb.AddProductImageResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.AddProductImageResponse{
		Image: modelImageToProto(result),
	}, nil
}

// DeleteProductImage deletes an image
func (h *ProductHandler) DeleteProductImage(ctx context.Context, req *pb.DeleteProductImageRequest) (*pb.DeleteProductImageResponse, error) {
	if req.ImageId == 0 {
		return nil, status.Error(codes.InvalidArgument, "image_id is required")
	}

	if err := h.productService.DeleteProductImage(ctx, req.ImageId); err != nil {
		return &pb.DeleteProductImageResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &pb.DeleteProductImageResponse{
		Success: true,
	}, nil
}

// AddSize adds a size with inventory
func (h *ProductHandler) AddSize(ctx context.Context, req *pb.AddSizeRequest) (*pb.AddSizeResponse, error) {
	if req.ProductId == 0 || req.Size == "" {
		return &pb.AddSizeResponse{
			Error: "product_id and size are required",
		}, nil
	}

	size := &model.Size{
		ProductID: req.ProductId,
		Size:      req.Size,
		Quantity:  int(req.Quantity),
	}

	result, err := h.productService.AddSize(ctx, size)
	if err != nil {
		return &pb.AddSizeResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.AddSizeResponse{
		Size: modelSizeToProto(result),
	}, nil
}

// GetAvailableSizes retrieves available sizes
func (h *ProductHandler) GetAvailableSizes(ctx context.Context, req *pb.GetAvailableSizesRequest) (*pb.GetAvailableSizesResponse, error) {
	if req.ProductId == 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	sizes, err := h.productService.GetAvailableSizes(ctx, req.ProductId)
	if err != nil {
		return &pb.GetAvailableSizesResponse{
			Error: err.Error(),
		}, nil
	}

	protoSizes := make([]*pb.Size, len(sizes))
	for i, s := range sizes {
		protoSizes[i] = modelSizeToProto(s)
	}

	return &pb.GetAvailableSizesResponse{
		Sizes: protoSizes,
	}, nil
}

// UpdateInventory updates inventory quantity
func (h *ProductHandler) UpdateInventory(ctx context.Context, req *pb.UpdateInventoryRequest) (*pb.UpdateInventoryResponse, error) {
	if req.SizeId == 0 {
		return nil, status.Error(codes.InvalidArgument, "size_id is required")
	}

	size, err := h.productService.UpdateInventory(ctx, req.SizeId, int(req.Quantity), req.Notes)
	if err != nil {
		return &pb.UpdateInventoryResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.UpdateInventoryResponse{
		Size: modelSizeToProto(size),
	}, nil
}

// ReserveInventory reserves inventory for an order
func (h *ProductHandler) ReserveInventory(ctx context.Context, req *pb.ReserveInventoryRequest) (*pb.ReserveInventoryResponse, error) {
	if req.SizeId == 0 || req.OrderId == "" {
		return &pb.ReserveInventoryResponse{
			Success: false,
			Error:   "size_id and order_id are required",
		}, nil
	}

	if err := h.productService.ReserveInventory(ctx, req.SizeId, int(req.Quantity), req.OrderId); err != nil {
		return &pb.ReserveInventoryResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &pb.ReserveInventoryResponse{
		Success: true,
	}, nil
}

// ReleaseInventory releases reserved inventory
func (h *ProductHandler) ReleaseInventory(ctx context.Context, req *pb.ReleaseInventoryRequest) (*pb.ReleaseInventoryResponse, error) {
	if req.SizeId == 0 || req.OrderId == "" {
		return &pb.ReleaseInventoryResponse{
			Success: false,
			Error:   "size_id and order_id are required",
		}, nil
	}

	if err := h.productService.ReleaseInventory(ctx, req.SizeId, int(req.Quantity), req.OrderId); err != nil {
		return &pb.ReleaseInventoryResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &pb.ReleaseInventoryResponse{
		Success: true,
	}, nil
}

// Helper functions to convert between model and proto

func modelProductToProto(product *model.Product, images []model.ProductImage, sizes []model.Size) *pb.Product {
	protoImages := make([]*pb.ProductImage, len(images))
	for i, img := range images {
		protoImages[i] = &pb.ProductImage{
			Id:           img.ID,
			ProductId:    img.ProductID,
			ImageUrl:     img.ImageURL,
			AltText:      img.AltText,
			DisplayOrder: int32(img.DisplayOrder),
			IsPrimary:    img.IsPrimary,
			CreatedAt:    img.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	protoSizes := make([]*pb.Size, len(sizes))
	for i, s := range sizes {
		protoSizes[i] = &pb.Size{
			Id:        s.ID,
			ProductId: s.ProductID,
			Size:      s.Size,
			Quantity:  int32(s.Quantity),
			Reserved:  int32(s.Reserved),
			CreatedAt: s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &pb.Product{
		Id:          product.ID,
		Sku:         product.SKU,
		Name:        product.Name,
		Brand:       product.Brand,
		Model:       product.Model,
		Color:       product.Color,
		Description: product.Description,
		Category:    product.Category,
		ReleaseYear: int64(product.ReleaseYear),
		RetailPrice: product.RetailPrice,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Images:      protoImages,
		Sizes:       protoSizes,
	}
}

func modelProductWithDetailsToProto(product *model.ProductWithDetails) *pb.Product {
	return modelProductToProto(&product.Product, product.Images, product.Sizes)
}

func modelImageToProto(image *model.ProductImage) *pb.ProductImage {
	return &pb.ProductImage{
		Id:           image.ID,
		ProductId:    image.ProductID,
		ImageUrl:     image.ImageURL,
		AltText:      image.AltText,
		DisplayOrder: int32(image.DisplayOrder),
		IsPrimary:    image.IsPrimary,
		CreatedAt:    image.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func modelSizeToProto(size *model.Size) *pb.Size {
	return &pb.Size{
		Id:        size.ID,
		ProductId: size.ProductID,
		Size:      size.Size,
		Quantity:  int32(size.Quantity),
		Reserved:  int32(size.Reserved),
		CreatedAt: size.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: size.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
