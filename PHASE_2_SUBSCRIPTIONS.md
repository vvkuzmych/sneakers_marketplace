# 💎 Phase 2: Premium Subscriptions System

**Goal:** Implement tiered subscription model with dynamic fees based on user subscription level.

---

## 🎯 Business Model

### **Subscription Tiers:**

| Plan | Monthly Price | Buyer Fee | Seller Fee | Benefits |
|------|---------------|-----------|------------|----------|
| **Free** | $0 | 1% | 2% | Basic marketplace access |
| **Pro** | $10/month | 1% | 1% | Reduced seller fees, priority support |
| **Elite** | $50/month | 0.5% | 0.5% | Lowest fees, early access, verified badge |

### **Revenue Model:**

```
Transaction Fees:
  Free users:  2% from seller = $2 per $100 sale
  Pro users:   1% from seller = $1 per $100 sale
  Elite users: 0.5% from seller = $0.50 per $100 sale

Subscription Revenue:
  Pro:   $10/month × users
  Elite: $50/month × users

Total Revenue = Transaction Fees + Subscription Fees
```

### **Example (100 sales @ $100 each):**

**Free User:**
- Transaction revenue: $200 (100 × $2)
- Subscription: $0
- **Total: $200**

**Pro User:**
- Transaction revenue: $100 (100 × $1)
- Subscription: $10/month
- **Total: $110/month**

**Elite User:**
- Transaction revenue: $50 (100 × $0.50)
- Subscription: $50/month
- **Total: $100/month**

---

## 📊 Phase 2 Implementation Plan

### **Day 1: Database Schema** (2 hours)
- `subscription_plans` table
- `user_subscriptions` table
- `subscription_transactions` table (Stripe payments)
- Migrations (up/down)

### **Day 2: Backend - Subscription Service** (3 hours)
- Subscription models
- Subscription repository
- Subscription service (create, cancel, upgrade)
- gRPC service definition

### **Day 3: Stripe Integration** (4 hours)
- Stripe API setup
- Webhook handlers (payment success/failed)
- Subscription creation flow
- Payment processing
- Webhook validation

### **Day 4: Dynamic Fees** (2 hours)
- Update Fee Service to check user subscription
- Apply tiered fees based on plan
- Update FeeBreakdown UI to show subscription discount
- "Upgrade to Pro" prompts

### **Day 5: Frontend UI** (3 hours)
- Subscription page (`/subscription`)
- Plan comparison cards
- Stripe Checkout integration
- Subscription management (cancel, upgrade)
- Billing history

### **Day 6: Testing & Polish** (2 hours)
- End-to-end testing
- Webhook testing
- Error handling
- Documentation

**Total: ~16 hours** (2 days of full-time work)

---

## 🗄️ Database Schema

### **subscription_plans**
```sql
CREATE TABLE subscription_plans (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,        -- 'free', 'pro', 'elite'
    display_name VARCHAR(100) NOT NULL,      -- 'Pro Plan'
    price_monthly DECIMAL(10,2) NOT NULL,    -- 10.00, 50.00
    buyer_fee_percent DECIMAL(5,2) NOT NULL, -- 1.00, 0.50
    seller_fee_percent DECIMAL(5,2) NOT NULL,-- 2.00, 1.00, 0.50
    stripe_price_id VARCHAR(255),            -- Stripe Price ID
    features JSONB,                          -- List of features
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Insert default plans
INSERT INTO subscription_plans (name, display_name, price_monthly, buyer_fee_percent, seller_fee_percent, features) VALUES
('free', 'Free', 0.00, 1.00, 2.00, '["Basic marketplace access", "Standard support"]'),
('pro', 'Pro', 10.00, 1.00, 1.00, '["Reduced seller fees (1%)", "Priority support", "Analytics dashboard"]'),
('elite', 'Elite', 50.00, 0.50, 0.50, '["Lowest fees (0.5%)", "Verified seller badge", "Early access to features", "Dedicated account manager"]');
```

### **user_subscriptions**
```sql
CREATE TABLE user_subscriptions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plan_id INTEGER NOT NULL REFERENCES subscription_plans(id),
    stripe_subscription_id VARCHAR(255),     -- Stripe Subscription ID
    stripe_customer_id VARCHAR(255),         -- Stripe Customer ID
    status VARCHAR(50) NOT NULL,             -- 'active', 'canceled', 'past_due', 'trialing'
    current_period_start TIMESTAMP,
    current_period_end TIMESTAMP,
    cancel_at_period_end BOOLEAN DEFAULT false,
    canceled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id)                          -- One subscription per user
);

CREATE INDEX idx_user_subscriptions_user_id ON user_subscriptions(user_id);
CREATE INDEX idx_user_subscriptions_status ON user_subscriptions(status);
```

### **subscription_transactions**
```sql
CREATE TABLE subscription_transactions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subscription_id BIGINT REFERENCES user_subscriptions(id),
    stripe_invoice_id VARCHAR(255),
    stripe_payment_intent_id VARCHAR(255),
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    status VARCHAR(50) NOT NULL,             -- 'succeeded', 'failed', 'pending'
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_subscription_transactions_user_id ON subscription_transactions(user_id);
```

---

## 🔧 Backend Architecture

### **Subscription Service (Go)**

```
internal/subscription/
├── model/
│   ├── subscription_plan.go
│   ├── user_subscription.go
│   └── subscription_transaction.go
├── repository/
│   └── subscription_repository.go
├── service/
│   └── subscription_service.go
└── handler/
    └── grpc_handler.go
```

### **Key Methods:**

```go
// Subscription Service
func (s *SubscriptionService) CreateSubscription(userID int64, planName string) (*UserSubscription, error)
func (s *SubscriptionService) CancelSubscription(userID int64) error
func (s *SubscriptionService) UpgradeSubscription(userID int64, newPlanName string) error
func (s *SubscriptionService) GetUserSubscription(userID int64) (*UserSubscription, error)
func (s *SubscriptionService) HandleStripeWebhook(event stripe.Event) error
```

### **Fee Service Integration:**

```go
// Modified Fee Service
func (s *FeeService) CalculateFees(ctx context.Context, userID int64, vertical string, salePrice float64) (*FeeBreakdown, error) {
    // Get user subscription
    subscription := s.subscriptionService.GetUserSubscription(ctx, userID)
    
    // Get plan fees
    plan := subscription.Plan
    
    // Apply tiered fees
    breakdown.BuyerProcessingFee = salePrice * (plan.BuyerFeePercent / 100)
    breakdown.SellerTransactionFee = salePrice * (plan.SellerFeePercent / 100)
    
    return breakdown, nil
}
```

---

## 💳 Stripe Integration

### **Setup:**

1. **Create Stripe Account** (test mode)
2. **Create Products in Stripe:**
   - Pro Subscription ($10/month)
   - Elite Subscription ($50/month)
3. **Get API Keys:**
   - Publishable Key (frontend)
   - Secret Key (backend)
4. **Setup Webhook:**
   - Endpoint: `https://yourdomain.com/api/v1/webhooks/stripe`
   - Events: `customer.subscription.created`, `customer.subscription.updated`, `customer.subscription.deleted`, `invoice.payment_succeeded`, `invoice.payment_failed`

### **Environment Variables:**

```bash
STRIPE_SECRET_KEY=sk_test_...
STRIPE_PUBLISHABLE_KEY=pk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...
STRIPE_PRO_PRICE_ID=price_...
STRIPE_ELITE_PRICE_ID=price_...
```

### **Webhook Handler:**

```go
func (h *SubscriptionHandler) HandleStripeWebhook(c *gin.Context) {
    payload, _ := ioutil.ReadAll(c.Request.Body)
    event, err := webhook.ConstructEvent(payload, c.GetHeader("Stripe-Signature"), webhookSecret)
    
    switch event.Type {
    case "customer.subscription.created":
        // Update user_subscriptions to 'active'
    case "invoice.payment_succeeded":
        // Record transaction in subscription_transactions
    case "invoice.payment_failed":
        // Update subscription to 'past_due', send email
    case "customer.subscription.deleted":
        // Update subscription to 'canceled'
    }
}
```

---

## 🎨 Frontend UI

### **Subscription Page (`/subscription`)**

```typescript
// SubscriptionPlans.tsx
interface Plan {
  name: string;
  displayName: string;
  price: number;
  buyerFee: number;
  sellerFee: number;
  features: string[];
  popular?: boolean;
}

const plans: Plan[] = [
  {
    name: 'free',
    displayName: 'Free',
    price: 0,
    buyerFee: 1.0,
    sellerFee: 2.0,
    features: ['Basic access', 'Standard support'],
  },
  {
    name: 'pro',
    displayName: 'Pro',
    price: 10,
    buyerFee: 1.0,
    sellerFee: 1.0,
    features: ['Reduced fees (1%)', 'Priority support', 'Analytics'],
    popular: true,
  },
  {
    name: 'elite',
    displayName: 'Elite',
    price: 50,
    buyerFee: 0.5,
    sellerFee: 0.5,
    features: ['Lowest fees (0.5%)', 'Verified badge', 'Dedicated manager'],
  },
];
```

### **Stripe Checkout Flow:**

```typescript
// Subscribe button
const handleSubscribe = async (planName: string) => {
  const response = await fetch('/api/v1/subscriptions/checkout', {
    method: 'POST',
    body: JSON.stringify({ plan: planName }),
  });
  
  const { sessionId } = await response.json();
  
  // Redirect to Stripe Checkout
  const stripe = await loadStripe(STRIPE_PUBLISHABLE_KEY);
  stripe.redirectToCheckout({ sessionId });
};
```

---

## ✅ Success Criteria

- [ ] Free users pay 2% seller fee
- [ ] Pro users pay 1% seller fee
- [ ] Elite users pay 0.5% seller fee
- [ ] Stripe payments work (test mode)
- [ ] Webhooks update subscription status
- [ ] Users can upgrade/downgrade plans
- [ ] Users can cancel subscription
- [ ] FeeBreakdown shows subscription discount
- [ ] Billing history visible
- [ ] Email notifications on payment events

---

## 🚀 Next Steps

1. Review this plan
2. Confirm business logic
3. Start Day 1: Database migrations

**Ready to start?** Say "**так**" or "**починаємо**"!

---

## 📚 Resources

- [Stripe Subscriptions Docs](https://stripe.com/docs/billing/subscriptions/overview)
- [Stripe Webhooks Guide](https://stripe.com/docs/webhooks)
- [Stripe Checkout](https://stripe.com/docs/payments/checkout)
