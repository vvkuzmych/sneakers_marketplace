package model

import "time"

// PlatformStats represents overall platform statistics
type PlatformStats struct {
	TotalUsers         int32
	ActiveUsersToday   int32
	TotalProducts      int32
	ActiveProducts     int32
	TotalOrders        int32
	OrdersToday        int32
	TotalRevenue       float64
	RevenueToday       float64
	TotalFeesCollected float64
	TotalMatches       int32
	MatchesToday       int32
}

// RevenueDataPoint represents a single point in revenue report
type RevenueDataPoint struct {
	Label      string // Date or period label
	Revenue    float64
	Fees       float64
	OrderCount int32
}

// RevenueReport represents a revenue report
type RevenueReport struct {
	DateFrom     time.Time
	DateTo       time.Time
	GroupBy      string
	DataPoints   []RevenueDataPoint
	TotalRevenue float64
	TotalFees    float64
}

// UserActivityReport represents user activity statistics
type UserActivityReport struct {
	DateFrom            time.Time
	DateTo              time.Time
	NewUsers            int32
	ActiveUsers         int32
	TotalBidsPlaced     int32
	TotalAsksPlaced     int32
	TotalMatchesCreated int32
}

// OrderSummary represents an order summary for admin view
type OrderSummary struct {
	CreatedAt   time.Time
	ProductName string
	OrderNumber string
	Status      string
	BuyerEmail  string
	SellerEmail string
	ProductID   int64
	ID          int64
	Subtotal    float64
	BuyerFee    float64
	SellerFee   float64
	Total       float64
	SellerID    int64
	BuyerID     int64
}

// OrderStatusChange represents a status change in order history
type OrderStatusChange struct {
	ChangedAt time.Time
	Status    string
	ChangedBy string
	Notes     string
}

// ProductSummary represents a product summary for admin view
type ProductSummary struct {
	CreatedAt    time.Time
	SKU          string
	Name         string
	Brand        string
	Model        string
	ID           int64
	RetailPrice  float64
	LowestAsk    float64
	HighestBid   float64
	TotalAsks    int32
	TotalMatches int32
	TotalBids    int32
	IsFeatured   bool
	IsActive     bool
}

// ServiceHealth represents health status of a service
type ServiceHealth struct {
	ServiceName   string
	Status        string
	Address       string
	Version       string
	UptimeSeconds int64
}

// DatabaseHealth represents database health metrics
type DatabaseHealth struct {
	Healthy           bool
	ActiveConnections int32
	IdleConnections   int32
	MaxConnections    int32
	QueryAvgTimeMs    float64
}

// ServiceMetric represents metrics for a service
type ServiceMetric struct {
	ServiceName       string
	TotalRequests     int64
	FailedRequests    int64
	AvgResponseTimeMs float64
	RequestsPerSecond float64
	ErrorRate         float64 // percentage
}

// GetRevenueReportParams contains parameters for revenue report
type GetRevenueReportParams struct {
	DateFrom time.Time
	DateTo   time.Time
	GroupBy  string // "day", "week", "month"
}

// GetUserActivityReportParams contains parameters for user activity report
type GetUserActivityReportParams struct {
	DateFrom time.Time
	DateTo   time.Time
}

// ListOrdersParams contains parameters for listing orders
type ListOrdersParams struct {
	DateFrom  *time.Time
	DateTo    *time.Time
	Status    string
	SortBy    string
	SortOrder string
	Page      int32
	PageSize  int32
}

// ListProductsParams contains parameters for listing products
type ListProductsParams struct {
	Status   string
	Search   string
	Page     int32
	PageSize int32
}
