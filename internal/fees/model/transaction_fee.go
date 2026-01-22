package model

import "time"

// TransactionFee represents a record of fees charged for a transaction
type TransactionFee struct {
	ID                      int64                  `json:"id" db:"id"`
	MatchID                 int64                  `json:"match_id" db:"match_id"`
	OrderID                 *int64                 `json:"order_id,omitempty" db:"order_id"`
	Vertical                string                 `json:"vertical" db:"vertical"`
	SalePrice               float64                `json:"sale_price" db:"sale_price"`
	BuyerProcessingFee      float64                `json:"buyer_processing_fee" db:"buyer_processing_fee"`
	BuyerShippingFee        float64                `json:"buyer_shipping_fee" db:"buyer_shipping_fee"`
	BuyerTotal              float64                `json:"buyer_total" db:"buyer_total"`
	SellerTransactionFee    float64                `json:"seller_transaction_fee" db:"seller_transaction_fee"`
	SellerAuthenticationFee float64                `json:"seller_authentication_fee" db:"seller_authentication_fee"`
	SellerShippingCost      float64                `json:"seller_shipping_cost" db:"seller_shipping_cost"`
	SellerPayout            float64                `json:"seller_payout" db:"seller_payout"`
	PlatformRevenue         float64                `json:"platform_revenue" db:"platform_revenue"`
	FeeConfigSnapshot       map[string]interface{} `json:"fee_config_snapshot" db:"fee_config_snapshot"`
	CreatedAt               time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time              `json:"updated_at" db:"updated_at"`
}
