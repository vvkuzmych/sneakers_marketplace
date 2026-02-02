# 💰 Dynamic Fees Based on Subscription

Complete implementation of subscription-based dynamic fee pricing for the sneakers marketplace.

---

## 🎯 Overview

The platform now uses **dynamic fees** based on the seller's subscription tier:

| Subscription Plan | Seller Fee | Buyer Fee | Monthly Cost |
|-------------------|------------|-----------|--------------|
| **Free**          | 1.0%       | 1.0%      | $0           |
| **Pro**           | 0.75%      | 1.0%      | $29          |
| **Elite**         | 0.5%       | 1.0%      | $99          |

### Key Points:

✅ **Seller pays platform fee** based on their subscription tier  
✅ **Buyer always pays 1%** processing fee (fixed)  
✅ **Free tier** = 1% for both buyer and seller  
✅ **Automatic fee calculation** at match creation  

---

## 🏗️ Architecture

```
┌─────────────────────┐
│  Bidding Service    │
│  (Match Created)    │
└──────────┬──────────┘
           │
           │ Calculate Fees
           ↓
┌─────────────────────┐
│   Fee Service       │───────┐
│                     │       │
└──────────┬──────────┘       │
           │                  │
           │ Get User's       │ Get Subscription
           │ Subscription     │ Fee Percentages
           ↓                  ↓
┌─────────────────────┐   ┌──────────────────────┐
│ Subscription        │   │  SubscriptionFee     │
│ Repository          │   │  Provider            │
└──────────┬──────────┘   └──────────────────────┘
           │
           │ Get Active
           │ Subscription
           │ with Plan
           ↓
┌─────────────────────┐
│  PostgreSQL         │
│  - subscription     │
│    _plans           │
│  - user_            │
│    subscriptions    │
└─────────────────────┘
```

---

## 📁 Files Changed

### 1. **Fee Service** (`internal/fees/service/`)

#### `subscription_provider.go` (New)
```go
// Interface for getting user's fee percentages
type SubscriptionFeeProvider interface {
    GetUserFeePercentages(ctx context.Context, userID int64) (float64, float64, error)
}

// Default provider (Free tier: 1%)
type DefaultFeeProvider struct{}
```

#### `fee_service.go` (Modified)
- Added `subscriptionProvider` dependency
- Modified `CalculateFees()` to:
  - Accept `sellerUserID` instead of `includeAuth`
  - Query seller's subscription tier
  - Apply dynamic fees based on subscription
  - Support `-1` for preview/default fees

**Before:**
```go
func CalculateFees(ctx context.Context, vertical string, price float64, includeAuth bool)
```

**After:**
```go
func CalculateFees(ctx context.Context, vertical string, price float64, sellerUserID int64)
```

### 2. **Subscription Service** (`internal/subscription/service/`)

#### `fee_provider.go` (New)
```go
type FeeProvider struct {
    repo repository.SubscriptionRepository
}

// Returns (sellerFee%, buyerFee%, error)
func (p *FeeProvider) GetUserFeePercentages(
    ctx context.Context,
    userID int64,
) (float64, float64, error) {
    // Get user's active subscription with plan
    subscription, err := p.repo.GetUserActiveSubscriptionWithPlan(ctx, userID)
    
    // Return fee percentages from plan
    return subscription.Plan.SellerFeePercent, 1.0, nil
}
```

### 3. **Bidding Service** (`internal/bidding/service/bidding_service.go`)

**Before:**
```go
feeBreakdown, err := s.feeService.CalculateFees(ctx, vertical, matchPrice, includeAuth)
```

**After:**
```go
// Pass seller's user ID for dynamic fee calculation
feeBreakdown, err := s.feeService.CalculateFees(ctx, vertical, matchPrice, match.SellerID)
```

### 4. **API Gateway**

#### `handlers/fee_handler.go`
- Updated `/api/v1/fees/calculate` endpoint
- Changed `include_auth` parameter to `seller_user_id`
- Support `-1` for preview mode (uses default Free tier fees)

#### `router/router.go`
- Initialize `subscriptionRepo`
- Create `subscriptionFeeProvider`
- Pass to `NewFeeHandler()`

**Before:**
```go
feeRepo := feeRepository.NewFeeRepository(db)
feeHandler := handlers.NewFeeHandler(feeRepo, log)
```

**After:**
```go
feeRepo := feeRepository.NewFeeRepository(db)
subscriptionRepo := subscriptionRepository.NewPostgresSubscriptionRepository(db)
subscriptionFeeProvider := subscriptionService.NewFeeProvider(subscriptionRepo)
feeHandler := handlers.NewFeeHandler(feeRepo, log, subscriptionFeeProvider)
```

### 5. **Main Services**

#### `cmd/bidding-service/main.go`
- Added subscription repository initialization
- Created fee provider with subscription support

#### `cmd/api-gateway/main.go`
- Updated through router changes

---

## 🔄 Fee Calculation Flow

### 1. **Match Created** (Bidding Service)
```go
// User A (buyer) bids $100
// User B (seller, Pro plan) asks $100
// Match created!

match := &model.Match{
    BuyerID:  userA_ID,  // 14
    SellerID: userB_ID,  // 42
    Price:    100.00,
}

// Calculate fees based on seller's subscription
feeBreakdown, err := s.feeService.CalculateFees(
    ctx,
    "sneakers",
    100.00,
    match.SellerID,  // 42 - User B (seller)
)
```

### 2. **Get Seller's Subscription** (Fee Service → Subscription Service)
```go
// Fee Service asks: "What are the fees for user 42?"
sellerFee, buyerFee, err := s.subscriptionProvider.GetUserFeePercentages(ctx, 42)

// Subscription Service returns:
// sellerFee = 0.75% (Pro plan)
// buyerFee = 1.0% (fixed)
```

### 3. **Calculate Fees**
```go
Sale Price: $100.00

Seller Fees (Pro plan):
  Platform Fee (0.75%): $0.75

Buyer Fees (fixed):
  Processing Fee (1%): $1.00

Totals:
  Buyer Pays: $101.00 (price + $1.00)
  Seller Receives: $99.25 (price - $0.75)
  Platform Revenue: $1.75
```

### 4. **Record Transaction** (Database)
```sql
INSERT INTO transaction_fees (
    match_id,
    vertical,
    sale_price,
    seller_transaction_fee,
    buyer_processing_fee,
    platform_revenue,
    seller_fee_config_snapshot,  -- {"fee_percent": 0.75, "plan": "Pro"}
    ...
) VALUES (...);
```

---

## 🧪 Testing

### Test Case 1: Free Tier User

```bash
# User with no subscription (Free tier default)
curl -X POST http://localhost:8080/api/v1/bids \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"product_id": 7, "size_id": 1, "price": 100, "quantity": 1}'

# Expected Fees:
# Seller Fee: 1% = $1.00
# Buyer Fee: 1% = $1.00
# Platform Revenue: $2.00
```

### Test Case 2: Pro Tier User

```sql
-- Set user to Pro plan
UPDATE user_subscriptions
SET plan_id = (SELECT id FROM subscription_plans WHERE name = 'pro'),
    status = 'active'
WHERE user_id = 42;
```

```bash
# User 42 places an ASK (seller)
# Buyer bids and match happens

# Expected Fees:
# Seller Fee: 0.75% = $0.75
# Buyer Fee: 1% = $1.00
# Platform Revenue: $1.75
```

### Test Case 3: Elite Tier User

```sql
-- Set user to Elite plan
UPDATE user_subscriptions
SET plan_id = (SELECT id FROM subscription_plans WHERE name = 'elite'),
    status = 'active'
WHERE user_id = 42;
```

```bash
# Expected Fees:
# Seller Fee: 0.5% = $0.50
# Buyer Fee: 1% = $1.00
# Platform Revenue: $1.50
```

### Test Case 4: Preview Mode (No User)

```bash
# Calculate fees without specific user (uses Free tier default)
curl "http://localhost:8080/api/v1/fees/calculate?vertical=sneakers&price=100"

# Or explicitly with seller_user_id=-1
curl "http://localhost:8080/api/v1/fees/calculate?vertical=sneakers&price=100&seller_user_id=-1"

# Expected:
# {
#   "sale_price": 100.00,
#   "seller_transaction_fee": 1.00,   // Free tier: 1%
#   "buyer_processing_fee": 1.00,     // Fixed: 1%
#   "platform_revenue": 2.00
# }
```

---

## 📊 Database Queries

### Check User's Current Subscription

```sql
SELECT 
    us.user_id,
    sp.name AS plan_name,
    sp.seller_fee_percent,
    us.status,
    us.current_period_start,
    us.current_period_end
FROM user_subscriptions us
JOIN subscription_plans sp ON us.plan_id = sp.id
WHERE us.user_id = 42
  AND us.status IN ('active', 'trialing')
ORDER BY us.created_at DESC
LIMIT 1;
```

### Transaction Fees by Subscription Plan

```sql
SELECT 
    sp.name AS plan_name,
    COUNT(*) AS transactions,
    AVG(tf.seller_transaction_fee) AS avg_seller_fee,
    SUM(tf.platform_revenue) AS total_revenue
FROM transaction_fees tf
JOIN matches m ON tf.match_id = m.id
JOIN user_subscriptions us ON m.seller_id = us.user_id
JOIN subscription_plans sp ON us.plan_id = sp.id
WHERE tf.created_at >= NOW() - INTERVAL '30 days'
GROUP BY sp.name
ORDER BY total_revenue DESC;
```

### Revenue Impact by Plan

```sql
WITH plan_revenue AS (
    SELECT 
        sp.name AS plan,
        COUNT(*) AS sales,
        SUM(tf.sale_price) AS gmv,
        SUM(tf.platform_revenue) AS revenue,
        AVG(tf.seller_transaction_fee / tf.sale_price * 100) AS avg_fee_percent
    FROM transaction_fees tf
    JOIN matches m ON tf.match_id = m.id
    JOIN user_subscriptions us ON m.seller_id = us.user_id
    JOIN subscription_plans sp ON us.plan_id = sp.id
    WHERE tf.created_at >= NOW() - INTERVAL '30 days'
      AND us.status = 'active'
    GROUP BY sp.name
)
SELECT 
    plan,
    sales,
    TO_CHAR(gmv, 'FM$999,999,999.00') AS gmv,
    TO_CHAR(revenue, 'FM$999,999,999.00') AS revenue,
    ROUND(avg_fee_percent, 2) || '%' AS avg_fee
FROM plan_revenue
ORDER BY revenue DESC;
```

---

## 🎓 Key Concepts

### 1. **Dependency Injection**

Fee Service accepts `SubscriptionFeeProvider` interface:
```go
type FeeService struct {
    repo                 *repository.FeeRepository
    log                  *logger.Logger
    subscriptionProvider SubscriptionFeeProvider  // ← Interface!
}
```

Benefits:
- ✅ Testable (mock provider)
- ✅ Flexible (swap implementations)
- ✅ Decoupled (fees doesn't depend on subscription concrete types)

### 2. **Interface Segregation**

`SubscriptionFeeProvider` only exposes what's needed:
```go
type SubscriptionFeeProvider interface {
    GetUserFeePercentages(ctx context.Context, userID int64) (float64, float64, error)
}
```

Subscription Service implements it, but Fee Service doesn't know about:
- Stripe integration
- Subscription management
- Plan changes
- etc.

### 3. **Default Values**

Special handling for missing subscriptions:
```go
// No subscription found → Use Free tier (1%)
if err != nil {
    return 1.0, 1.0, nil
}

// Inactive subscription → Use Free tier (1%)
if subscription.Status != "active" && subscription.Status != "trialing" {
    return 1.0, 1.0, nil
}
```

### 4. **Preview Mode**

Support calculating fees without a real user:
```go
// sellerUserID = -1 → Use default Free tier fees
if sellerUserID == -1 {
    sellerFeePercent = 1.0
    buyerFeePercent = 1.0
}
```

---

## 🚀 Next Steps

### Phase 2, Day 5: Frontend Subscription UI

1. **Subscription Plans Page**
   - Display Free, Pro, Elite plans
   - Show savings calculator
   - "Upgrade" buttons

2. **Checkout Flow**
   - Stripe payment form
   - Subscription confirmation
   - Receipt/invoice

3. **User Dashboard**
   - Current plan display
   - Upgrade/downgrade options
   - Billing history
   - Fee savings tracker

4. **Real-time Fee Preview**
   - When creating ASK, show expected fees
   - "With Pro plan, you'd save $X"
   - Upsell opportunities

---

## ✅ Summary

**Completed:**
- ✅ Subscription-based fee calculation
- ✅ Dynamic fees at match creation
- ✅ API endpoint for fee preview
- ✅ Default fees for non-subscribers
- ✅ Integration with all services
- ✅ Database tracking

**Benefits:**
- 📈 Incentivizes subscriptions
- 💰 Tiered pricing model
- 🎯 Seller-focused monetization
- 📊 Transparent fee structure
- 🔄 Automatic fee application

**Example:**

Free tier seller selling $1000 item:
- Seller pays: $10 (1%)
- Buyer pays: $10 (1%)
- **Total platform revenue: $20**

Pro tier seller selling $1000 item:
- Seller pays: $7.50 (0.75%)
- Buyer pays: $10 (1%)
- **Total platform revenue: $17.50**
- **Seller saves: $2.50**

Elite tier seller selling $1000 item:
- Seller pays: $5 (0.5%)
- Buyer pays: $10 (1%)
- **Total platform revenue: $15**
- **Seller saves: $5**

---

**Phase 2, Day 4: Dynamic Fees ✅ COMPLETE!**

Next: Phase 2, Day 5 - Frontend Subscription UI 🚀
