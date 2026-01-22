package model

// FeeBreakdown provides a detailed breakdown of all fees for a transaction
// Used for displaying fee information to users before they commit to a transaction
type FeeBreakdown struct {
	SalePrice            float64 `json:"sale_price"`
	BuyerProcessingFee   float64 `json:"buyer_processing_fee"`
	BuyerShippingFee     float64 `json:"buyer_shipping_fee"`
	BuyerTotal           float64 `json:"buyer_total"`
	SellerTransactionFee float64 `json:"seller_transaction_fee"`
	SellerAuthFee        float64 `json:"seller_auth_fee"`
	SellerShippingCost   float64 `json:"seller_shipping_cost"`
	SellerPayout         float64 `json:"seller_payout"`
	PlatformRevenue      float64 `json:"platform_revenue"`
}

// NewFeeBreakdown creates a new FeeBreakdown from basic parameters
func NewFeeBreakdown(salePrice float64) *FeeBreakdown {
	return &FeeBreakdown{
		SalePrice: salePrice,
	}
}

// CalculateTotals calculates buyer total, seller payout, and platform revenue
func (fb *FeeBreakdown) CalculateTotals() {
	// Buyer pays: sale price + fees
	fb.BuyerTotal = fb.SalePrice + fb.BuyerProcessingFee + fb.BuyerShippingFee

	// Seller receives: sale price - fees
	fb.SellerPayout = fb.SalePrice - fb.SellerTransactionFee - fb.SellerAuthFee - fb.SellerShippingCost

	// Platform earns: all fees
	fb.PlatformRevenue = fb.BuyerProcessingFee + fb.SellerTransactionFee + fb.SellerAuthFee

	// Add shipping markup (buyer pays more than seller costs)
	shippingMarkup := fb.BuyerShippingFee - fb.SellerShippingCost
	if shippingMarkup > 0 {
		fb.PlatformRevenue += shippingMarkup
	}
}

// ToTransactionFee converts FeeBreakdown to TransactionFee model for database storage
func (fb *FeeBreakdown) ToTransactionFee(matchID int64, vertical string, feeConfig *FeeConfig) *TransactionFee {
	snapshot := map[string]interface{}{
		"transaction_fee_percent": feeConfig.TransactionFeePercent,
		"processing_fee_fixed":    feeConfig.ProcessingFeeFixed,
		"authentication_fee":      feeConfig.AuthenticationFee,
		"shipping_buyer_charge":   feeConfig.ShippingBuyerCharge,
		"shipping_seller_cost":    feeConfig.ShippingSellerCost,
	}

	return &TransactionFee{
		MatchID:                 matchID,
		Vertical:                vertical,
		SalePrice:               fb.SalePrice,
		BuyerProcessingFee:      fb.BuyerProcessingFee,
		BuyerShippingFee:        fb.BuyerShippingFee,
		BuyerTotal:              fb.BuyerTotal,
		SellerTransactionFee:    fb.SellerTransactionFee,
		SellerAuthenticationFee: fb.SellerAuthFee,
		SellerShippingCost:      fb.SellerShippingCost,
		SellerPayout:            fb.SellerPayout,
		PlatformRevenue:         fb.PlatformRevenue,
		FeeConfigSnapshot:       snapshot,
	}
}
