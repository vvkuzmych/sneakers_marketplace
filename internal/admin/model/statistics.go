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
	DataPoints   []RevenueDataPoint
	TotalRevenue float64
	TotalFees    float64
	DateFrom     time.Time
	DateTo       time.Time
	GroupBy      string // "day", "week", "month"
}

// UserActivityReport represents user activity statistics
type UserActivityReport struct {
	NewUsers            int32
	ActiveUsers         int32
	TotalBidsPlaced     int32
	TotalAsksPlaced     int32
	TotalMatchesCreated int32
	DateFrom            time.Time
	DateTo              time.Time
}

// OrderSummary represents an order summary for admin view
type OrderSummary struct {
	ID          int64
	OrderNumber string
	BuyerID     int64
	SellerID    int64
	BuyerEmail  string
	SellerEmail string
	ProductID   int64
	ProductName string
	Subtotal    float64
	BuyerFee    float64
	SellerFee   float64
	Total       float64
	Status      string
	CreatedAt   time.Time
}

// OrderStatusChange represents a status change in order history
type OrderStatusChange struct {
	Status    string
	ChangedBy string
	Notes     string
	ChangedAt time.Time
}

// ProductSummary represents a product summary for admin view
type ProductSummary struct {
	ID           int64
	SKU          string
	Name         string
	Brand        string
	Model        string
	RetailPrice  float64
	IsActive     bool
	IsFeatured   bool
	TotalBids    int32
	TotalAsks    int32
	TotalMatches int32
	HighestBid   float64
	LowestAsk    float64
	CreatedAt    time.Time
}

// ServiceHealth represents health status of a service
type ServiceHealth struct {
	ServiceName   string
	Status        string // "healthy", "degraded", "down"
	Address       string
	UptimeSeconds int64
	Version       string
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
	Page      int32
	PageSize  int32
	Status    string
	SortBy    string
	SortOrder string
	DateFrom  *time.Time
	DateTo    *time.Time
}

// ListProductsParams contains parameters for listing products
type ListProductsParams struct {
	Page     int32
	PageSize int32
	Status   string // "all", "active", "hidden", "featured"
	Search   string
}
