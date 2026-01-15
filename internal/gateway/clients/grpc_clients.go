package clients

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	biddingPb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/bidding"
	orderPb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/order"
	paymentPb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/payment"
	productPb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/product"
	userPb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/user"
)

// GRPCClients holds all gRPC service clients
type GRPCClients struct {
	UserClient    userPb.UserServiceClient
	ProductClient productPb.ProductServiceClient
	BiddingClient biddingPb.BiddingServiceClient
	OrderClient   orderPb.OrderServiceClient
	PaymentClient paymentPb.PaymentServiceClient
}

// NewGRPCClients creates connections to all gRPC services
func NewGRPCClients(
	userAddr, productAddr, biddingAddr, orderAddr, paymentAddr string,
) (*GRPCClients, error) {
	// Connect to User Service
	userConn, err := grpc.NewClient(
		userAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	// Connect to Product Service
	productConn, err := grpc.NewClient(
		productAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}

	// Connect to Bidding Service
	biddingConn, err := grpc.NewClient(
		biddingAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to bidding service: %w", err)
	}

	// Connect to Order Service
	orderConn, err := grpc.NewClient(
		orderAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order service: %w", err)
	}

	// Connect to Payment Service
	paymentConn, err := grpc.NewClient(
		paymentAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service: %w", err)
	}

	log.Println("✅ Connected to all gRPC services")

	return &GRPCClients{
		UserClient:    userPb.NewUserServiceClient(userConn),
		ProductClient: productPb.NewProductServiceClient(productConn),
		BiddingClient: biddingPb.NewBiddingServiceClient(biddingConn),
		OrderClient:   orderPb.NewOrderServiceClient(orderConn),
		PaymentClient: paymentPb.NewPaymentServiceClient(paymentConn),
	}, nil
}
