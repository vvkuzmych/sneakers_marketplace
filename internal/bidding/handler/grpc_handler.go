package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vvkuzmych/sneakers_marketplace/internal/bidding/model"
	"github.com/vvkuzmych/sneakers_marketplace/internal/bidding/service"
	pb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/bidding"
)

// BiddingHandler implements the gRPC BiddingService server
type BiddingHandler struct {
	pb.UnimplementedBiddingServiceServer
	biddingService *service.BiddingService
}

// NewBiddingHandler creates a new bidding handler
func NewBiddingHandler(biddingService *service.BiddingService) *BiddingHandler {
	return &BiddingHandler{
		biddingService: biddingService,
	}
}

// PlaceBid handles bid placement
func (h *BiddingHandler) PlaceBid(ctx context.Context, req *pb.PlaceBidRequest) (*pb.PlaceBidResponse, error) {
	if req.UserId == 0 || req.ProductId == 0 || req.SizeId == 0 || req.Price <= 0 {
		return &pb.PlaceBidResponse{
			Error: "user_id, product_id, size_id, and price are required",
		}, nil
	}

	bid, match, err := h.biddingService.PlaceBid(
		ctx,
		req.UserId,
		req.ProductId,
		req.SizeId,
		req.Price,
		int(req.Quantity),
		int(req.ExpiresInHours),
	)

	if err != nil {
		return &pb.PlaceBidResponse{
			Error: err.Error(),
		}, nil
	}

	response := &pb.PlaceBidResponse{
		Bid: modelBidToProto(bid),
	}

	if match != nil {
		response.Match = modelMatchToProto(match)
	}

	return response, nil
}

// GetBid retrieves a bid
func (h *BiddingHandler) GetBid(ctx context.Context, req *pb.GetBidRequest) (*pb.GetBidResponse, error) {
	if req.BidId == 0 {
		return nil, status.Error(codes.InvalidArgument, "bid_id is required")
	}

	bid, err := h.biddingService.GetBid(ctx, req.BidId)
	if err != nil {
		return &pb.GetBidResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.GetBidResponse{
		Bid: modelBidToProto(bid),
	}, nil
}

// GetUserBids retrieves user's bids
func (h *BiddingHandler) GetUserBids(ctx context.Context, req *pb.GetUserBidsRequest) (*pb.GetUserBidsResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	bids, total, err := h.biddingService.GetUserBids(
		ctx,
		req.UserId,
		req.Status,
		int(req.Page),
		int(req.PageSize),
	)

	if err != nil {
		return &pb.GetUserBidsResponse{
			Error: err.Error(),
		}, nil
	}

	protoBids := make([]*pb.Bid, len(bids))
	for i, bid := range bids {
		protoBids[i] = modelBidToProto(bid)
	}

	return &pb.GetUserBidsResponse{
		Bids:  protoBids,
		Total: total,
	}, nil
}

// GetProductBids retrieves bids for a product
func (h *BiddingHandler) GetProductBids(ctx context.Context, req *pb.GetProductBidsRequest) (*pb.GetProductBidsResponse, error) {
	if req.ProductId == 0 || req.SizeId == 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id and size_id are required")
	}

	bids, total, err := h.biddingService.GetProductBids(
		ctx,
		req.ProductId,
		req.SizeId,
		req.Status,
		int(req.Page),
		int(req.PageSize),
	)

	if err != nil {
		return &pb.GetProductBidsResponse{
			Error: err.Error(),
		}, nil
	}

	protoBids := make([]*pb.Bid, len(bids))
	for i, bid := range bids {
		protoBids[i] = modelBidToProto(bid)
	}

	return &pb.GetProductBidsResponse{
		Bids:  protoBids,
		Total: total,
	}, nil
}

// CancelBid cancels a bid
func (h *BiddingHandler) CancelBid(ctx context.Context, req *pb.CancelBidRequest) (*pb.CancelBidResponse, error) {
	if req.BidId == 0 || req.UserId == 0 {
		return &pb.CancelBidResponse{
			Success: false,
			Error:   "bid_id and user_id are required",
		}, nil
	}

	if err := h.biddingService.CancelBid(ctx, req.BidId, req.UserId); err != nil {
		return &pb.CancelBidResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &pb.CancelBidResponse{
		Success: true,
	}, nil
}

// PlaceAsk handles ask placement
func (h *BiddingHandler) PlaceAsk(ctx context.Context, req *pb.PlaceAskRequest) (*pb.PlaceAskResponse, error) {
	if req.UserId == 0 || req.ProductId == 0 || req.SizeId == 0 || req.Price <= 0 {
		return &pb.PlaceAskResponse{
			Error: "user_id, product_id, size_id, and price are required",
		}, nil
	}

	ask, match, err := h.biddingService.PlaceAsk(
		ctx,
		req.UserId,
		req.ProductId,
		req.SizeId,
		req.Price,
		int(req.Quantity),
		int(req.ExpiresInHours),
	)

	if err != nil {
		return &pb.PlaceAskResponse{
			Error: err.Error(),
		}, nil
	}

	response := &pb.PlaceAskResponse{
		Ask: modelAskToProto(ask),
	}

	if match != nil {
		response.Match = modelMatchToProto(match)
	}

	return response, nil
}

// GetAsk retrieves an ask
func (h *BiddingHandler) GetAsk(ctx context.Context, req *pb.GetAskRequest) (*pb.GetAskResponse, error) {
	if req.AskId == 0 {
		return nil, status.Error(codes.InvalidArgument, "ask_id is required")
	}

	ask, err := h.biddingService.GetAsk(ctx, req.AskId)
	if err != nil {
		return &pb.GetAskResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.GetAskResponse{
		Ask: modelAskToProto(ask),
	}, nil
}

// GetUserAsks retrieves user's asks
func (h *BiddingHandler) GetUserAsks(ctx context.Context, req *pb.GetUserAsksRequest) (*pb.GetUserAsksResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	asks, total, err := h.biddingService.GetUserAsks(
		ctx,
		req.UserId,
		req.Status,
		int(req.Page),
		int(req.PageSize),
	)

	if err != nil {
		return &pb.GetUserAsksResponse{
			Error: err.Error(),
		}, nil
	}

	protoAsks := make([]*pb.Ask, len(asks))
	for i, ask := range asks {
		protoAsks[i] = modelAskToProto(ask)
	}

	return &pb.GetUserAsksResponse{
		Asks:  protoAsks,
		Total: total,
	}, nil
}

// GetProductAsks retrieves asks for a product
func (h *BiddingHandler) GetProductAsks(ctx context.Context, req *pb.GetProductAsksRequest) (*pb.GetProductAsksResponse, error) {
	if req.ProductId == 0 || req.SizeId == 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id and size_id are required")
	}

	asks, total, err := h.biddingService.GetProductAsks(
		ctx,
		req.ProductId,
		req.SizeId,
		req.Status,
		int(req.Page),
		int(req.PageSize),
	)

	if err != nil {
		return &pb.GetProductAsksResponse{
			Error: err.Error(),
		}, nil
	}

	protoAsks := make([]*pb.Ask, len(asks))
	for i, ask := range asks {
		protoAsks[i] = modelAskToProto(ask)
	}

	return &pb.GetProductAsksResponse{
		Asks:  protoAsks,
		Total: total,
	}, nil
}

// CancelAsk cancels an ask
func (h *BiddingHandler) CancelAsk(ctx context.Context, req *pb.CancelAskRequest) (*pb.CancelAskResponse, error) {
	if req.AskId == 0 || req.UserId == 0 {
		return &pb.CancelAskResponse{
			Success: false,
			Error:   "ask_id and user_id are required",
		}, nil
	}

	if err := h.biddingService.CancelAsk(ctx, req.AskId, req.UserId); err != nil {
		return &pb.CancelAskResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &pb.CancelAskResponse{
		Success: true,
	}, nil
}

// GetHighestBid retrieves highest bid
func (h *BiddingHandler) GetHighestBid(ctx context.Context, req *pb.GetHighestBidRequest) (*pb.GetHighestBidResponse, error) {
	if req.ProductId == 0 || req.SizeId == 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id and size_id are required")
	}

	bid, err := h.biddingService.GetHighestBid(ctx, req.ProductId, req.SizeId)
	if err != nil {
		return &pb.GetHighestBidResponse{
			Error: err.Error(),
		}, nil
	}

	if bid == nil {
		return &pb.GetHighestBidResponse{}, nil
	}

	return &pb.GetHighestBidResponse{
		Bid: modelBidToProto(bid),
	}, nil
}

// GetLowestAsk retrieves lowest ask
func (h *BiddingHandler) GetLowestAsk(ctx context.Context, req *pb.GetLowestAskRequest) (*pb.GetLowestAskResponse, error) {
	if req.ProductId == 0 || req.SizeId == 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id and size_id are required")
	}

	ask, err := h.biddingService.GetLowestAsk(ctx, req.ProductId, req.SizeId)
	if err != nil {
		return &pb.GetLowestAskResponse{
			Error: err.Error(),
		}, nil
	}

	if ask == nil {
		return &pb.GetLowestAskResponse{}, nil
	}

	return &pb.GetLowestAskResponse{
		Ask: modelAskToProto(ask),
	}, nil
}

// GetMarketPrice retrieves market price data
func (h *BiddingHandler) GetMarketPrice(ctx context.Context, req *pb.GetMarketPriceRequest) (*pb.GetMarketPriceResponse, error) {
	if req.ProductId == 0 || req.SizeId == 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id and size_id are required")
	}

	marketPrice, err := h.biddingService.GetMarketPrice(ctx, req.ProductId, req.SizeId)
	if err != nil {
		return &pb.GetMarketPriceResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.GetMarketPriceResponse{
		HighestBid: marketPrice.HighestBid,
		LowestAsk:  marketPrice.LowestAsk,
		LastSale:   marketPrice.LastSale,
		TotalBids:  marketPrice.TotalBids,
		TotalAsks:  marketPrice.TotalAsks,
	}, nil
}

// GetMatch retrieves a match
func (h *BiddingHandler) GetMatch(ctx context.Context, req *pb.GetMatchRequest) (*pb.GetMatchResponse, error) {
	if req.MatchId == 0 {
		return nil, status.Error(codes.InvalidArgument, "match_id is required")
	}

	match, err := h.biddingService.GetMatch(ctx, req.MatchId)
	if err != nil {
		return &pb.GetMatchResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.GetMatchResponse{
		Match: modelMatchToProto(match),
	}, nil
}

// GetUserMatches retrieves user's matches
func (h *BiddingHandler) GetUserMatches(ctx context.Context, req *pb.GetUserMatchesRequest) (*pb.GetUserMatchesResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	matches, total, err := h.biddingService.GetUserMatches(
		ctx,
		req.UserId,
		req.AsBuyer,
		req.AsSeller,
		int(req.Page),
		int(req.PageSize),
	)

	if err != nil {
		return &pb.GetUserMatchesResponse{
			Error: err.Error(),
		}, nil
	}

	protoMatches := make([]*pb.Match, len(matches))
	for i, match := range matches {
		protoMatches[i] = modelMatchToProto(match)
	}

	return &pb.GetUserMatchesResponse{
		Matches: protoMatches,
		Total:   total,
	}, nil
}

// Helper functions to convert between model and proto

func modelBidToProto(bid *model.Bid) *pb.Bid {
	protoBid := &pb.Bid{
		Id:        bid.ID,
		UserId:    bid.UserID,
		ProductId: bid.ProductID,
		SizeId:    bid.SizeID,
		Price:     bid.Price,
		Quantity:  int32(bid.Quantity),
		Status:    bid.Status,
		CreatedAt: bid.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: bid.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if bid.ExpiresAt != nil {
		protoBid.ExpiresAt = bid.ExpiresAt.Format("2006-01-02T15:04:05Z07:00")
	}

	if bid.MatchedAt != nil {
		protoBid.MatchedAt = bid.MatchedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	return protoBid
}

func modelAskToProto(ask *model.Ask) *pb.Ask {
	protoAsk := &pb.Ask{
		Id:        ask.ID,
		UserId:    ask.UserID,
		ProductId: ask.ProductID,
		SizeId:    ask.SizeID,
		Price:     ask.Price,
		Quantity:  int32(ask.Quantity),
		Status:    ask.Status,
		CreatedAt: ask.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: ask.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if ask.ExpiresAt != nil {
		protoAsk.ExpiresAt = ask.ExpiresAt.Format("2006-01-02T15:04:05Z07:00")
	}

	if ask.MatchedAt != nil {
		protoAsk.MatchedAt = ask.MatchedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	return protoAsk
}

func modelMatchToProto(match *model.Match) *pb.Match {
	protoMatch := &pb.Match{
		Id:        match.ID,
		BidId:     match.BidID,
		AskId:     match.AskID,
		BuyerId:   match.BuyerID,
		SellerId:  match.SellerID,
		ProductId: match.ProductID,
		SizeId:    match.SizeID,
		Price:     match.Price,
		Quantity:  int32(match.Quantity),
		Status:    match.Status,
		CreatedAt: match.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if match.CompletedAt != nil {
		protoMatch.CompletedAt = match.CompletedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	return protoMatch
}
