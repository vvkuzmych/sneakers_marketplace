# 💳 Stripe Integration Setup Guide

Complete guide for setting up Stripe payment integration for the subscription system.

---

## 🔑 Step 1: Get Stripe API Keys

1. Go to [Stripe Dashboard](https://dashboard.stripe.com)
2. Create account or log in
3. Navigate to **Developers → API keys**
4. Copy your keys:
   - **Publishable key** (starts with `pk_test_` or `pk_live_`)
   - **Secret key** (starts with `sk_test_` or `sk_live_`)

---

## ⚙️ Step 2: Configure Environment Variables

Add these to your `.env` file:

```bash
# ====================================
# STRIPE CONFIGURATION
# ====================================

# API Keys (from Stripe Dashboard)
STRIPE_SECRET_KEY=sk_test_51ABC...xyz
STRIPE_PUBLISHABLE_KEY=pk_test_51ABC...xyz

# Webhook signing secret (we'll get this in Step 3)
STRIPE_WEBHOOK_SECRET=whsec_...

# Subscription Plan Price IDs (we'll create these in Step 4)
STRIPE_PRICE_PRO_MONTHLY=price_1ABC...xyz
STRIPE_PRICE_PRO_YEARLY=price_1ABC...xyz
STRIPE_PRICE_ELITE_MONTHLY=price_1ABC...xyz
STRIPE_PRICE_ELITE_YEARLY=price_1ABC...xyz
```

---

## 🪝 Step 3: Setup Webhooks

Webhooks allow Stripe to notify your application about events (payments, subscription changes, etc.).

### Local Development (using Stripe CLI)

1. **Install Stripe CLI:**

```bash
# macOS
brew install stripe/stripe-cli/stripe

# Linux/Windows: Download from https://github.com/stripe/stripe-cli/releases
```

2. **Login to Stripe:**

```bash
stripe login
```

3. **Forward webhooks to local server:**

```bash
# Forward webhooks to your local API Gateway
stripe listen --forward-to localhost:8080/api/v1/webhooks/stripe

# You'll see output like:
# > Ready! Your webhook signing secret is whsec_abc123...
```

4. **Copy the webhook secret** and add to `.env`:

```bash
STRIPE_WEBHOOK_SECRET=whsec_abc123...
```

### Production Setup

1. Go to **Stripe Dashboard → Developers → Webhooks**
2. Click **Add endpoint**
3. Enter your webhook URL:
   ```
   https://yourdomain.com/api/v1/webhooks/stripe
   ```
4. Select events to listen to:
   - `customer.subscription.created`
   - `customer.subscription.updated`
   - `customer.subscription.deleted`
   - `invoice.payment_succeeded`
   - `invoice.payment_failed`
5. Copy the **Signing secret** and add to production `.env`

---

## 💰 Step 4: Create Subscription Products & Prices

### Option A: Via Stripe Dashboard (Recommended for beginners)

1. Go to **Products → Add product**

2. **Create Pro Plan:**
   - **Name:** `Pro Plan`
   - **Description:** `Professional tier with reduced fees`
   - **Pricing:**
     - Monthly: `$29/month` → Copy Price ID (e.g., `price_1ABC...xyz`)
     - Yearly: `$290/year` (save $58/year) → Copy Price ID

3. **Create Elite Plan:**
   - **Name:** `Elite Plan`
   - **Description:** `Premium tier with lowest fees`
   - **Pricing:**
     - Monthly: `$99/month` → Copy Price ID
     - Yearly: `$990/year` (save $198/year) → Copy Price ID

4. **Add Price IDs to `.env`:**

```bash
STRIPE_PRICE_PRO_MONTHLY=price_1ABC...xyz
STRIPE_PRICE_PRO_YEARLY=price_1DEF...xyz
STRIPE_PRICE_ELITE_MONTHLY=price_1GHI...xyz
STRIPE_PRICE_ELITE_YEARLY=price_1JKL...xyz
```

### Option B: Via Stripe API (Advanced)

```bash
# Create Pro Monthly Price
stripe prices create \
  --unit-amount=2900 \
  --currency=usd \
  --recurring[interval]=month \
  --product_data[name]="Pro Plan"

# Create Pro Yearly Price
stripe prices create \
  --unit-amount=29000 \
  --currency=usd \
  --recurring[interval]=year \
  --product_data[name]="Pro Plan"

# Repeat for Elite plan...
```

---

## 🗄️ Step 5: Update Database with Price IDs

Run SQL to update `subscription_plans` table:

```sql
UPDATE subscription_plans
SET 
  stripe_price_id_monthly = 'price_1ABC...xyz',
  stripe_price_id_yearly = 'price_1DEF...xyz'
WHERE name = 'pro';

UPDATE subscription_plans
SET 
  stripe_price_id_monthly = 'price_1GHI...xyz',
  stripe_price_id_yearly = 'price_1JKL...xyz'
WHERE name = 'elite';
```

Or via script:

```bash
psql $DATABASE_URL -c "
UPDATE subscription_plans SET 
  stripe_price_id_monthly = '$STRIPE_PRICE_PRO_MONTHLY',
  stripe_price_id_yearly = '$STRIPE_PRICE_PRO_YEARLY'
WHERE name = 'pro';

UPDATE subscription_plans SET 
  stripe_price_id_monthly = '$STRIPE_PRICE_ELITE_MONTHLY',
  stripe_price_id_yearly = '$STRIPE_PRICE_ELITE_YEARLY'
WHERE name = 'elite';
"
```

---

## 🧪 Step 6: Test the Integration

### Test Card Numbers

Stripe provides test cards for development:

| Card Number | Description |
|-------------|-------------|
| `4242 4242 4242 4242` | Successful payment |
| `4000 0000 0000 0002` | Card declined |
| `4000 0025 0000 3155` | Requires authentication (3D Secure) |

- **Expiry:** Any future date (e.g., `12/34`)
- **CVC:** Any 3 digits (e.g., `123`)
- **ZIP:** Any 5 digits (e.g., `12345`)

### Test Subscription Flow

1. Start your services:
   ```bash
   make dev
   ```

2. In another terminal, start Stripe webhook forwarding:
   ```bash
   stripe listen --forward-to localhost:8080/api/v1/webhooks/stripe
   ```

3. Open frontend:
   ```
   http://localhost:5173
   ```

4. Navigate to subscription page and try upgrading to Pro

5. Use test card: `4242 4242 4242 4242`

6. Check logs:
   - API Gateway should log webhook events
   - Database should have new subscription record
   - Transaction should be recorded

---

## 🔄 Step 7: Handle Webhook Events

The application automatically handles these events:

✅ **customer.subscription.created** - New subscription started  
✅ **customer.subscription.updated** - Subscription changed (upgrade/downgrade)  
✅ **customer.subscription.deleted** - Subscription canceled  
✅ **invoice.payment_succeeded** - Payment successful  
✅ **invoice.payment_failed** - Payment failed  

Check logs for webhook processing:

```bash
tail -f /tmp/api-gateway.log | grep webhook
```

---

## 🚀 Step 8: Go Live

### Pre-launch Checklist

- [ ] Switch to **Live Mode** in Stripe Dashboard
- [ ] Get **Live API Keys** (starts with `pk_live_`, `sk_live_`)
- [ ] Update `.env` with live keys
- [ ] Setup production webhook endpoint
- [ ] Test with real card (small amount)
- [ ] Configure email notifications for failed payments
- [ ] Setup subscription billing alerts

### Switch to Live Keys

```bash
# Production .env
STRIPE_SECRET_KEY=sk_live_51XYZ...
STRIPE_PUBLISHABLE_KEY=pk_live_51XYZ...
STRIPE_WEBHOOK_SECRET=whsec_live_...
```

---

## 📊 Monitoring & Analytics

### Stripe Dashboard

Monitor in real-time:
- **Payments** - All transactions
- **Customers** - User subscriptions
- **Billing** - Invoices and receipts
- **Logs** - API requests

### Application Metrics

Query database for insights:

```sql
-- Active subscriptions by plan
SELECT p.name, COUNT(us.*) 
FROM user_subscriptions us
JOIN subscription_plans p ON us.plan_id = p.id
WHERE us.status = 'active'
GROUP BY p.name;

-- Monthly recurring revenue (MRR)
SELECT 
  SUM(CASE WHEN us.billing_cycle = 'monthly' THEN p.price_monthly ELSE p.price_yearly / 12 END) AS mrr
FROM user_subscriptions us
JOIN subscription_plans p ON us.plan_id = p.id
WHERE us.status = 'active';

-- Revenue by month
SELECT 
  DATE_TRUNC('month', created_at) AS month,
  SUM(amount) AS revenue
FROM subscription_transactions
WHERE status = 'succeeded'
GROUP BY month
ORDER BY month DESC;
```

---

## 🛠️ Troubleshooting

### Issue: Webhook signature verification failed

**Solution:** Make sure `STRIPE_WEBHOOK_SECRET` matches the secret from Stripe CLI or Dashboard.

### Issue: No Stripe price ID configured for plan

**Solution:** Update database with Stripe price IDs (Step 5).

### Issue: Payment fails immediately

**Solution:** 
- Check if using test mode keys with test cards
- Verify card number, expiry, and CVC
- Check Stripe Dashboard logs

### Issue: Subscription not created in database

**Solution:**
- Check API Gateway logs
- Verify database connection
- Ensure migrations ran successfully

---

## 📚 Additional Resources

- [Stripe API Documentation](https://stripe.com/docs/api)
- [Stripe Billing Documentation](https://stripe.com/docs/billing)
- [Stripe CLI](https://stripe.com/docs/stripe-cli)
- [Test Cards](https://stripe.com/docs/testing)
- [Webhooks Guide](https://stripe.com/docs/webhooks)

---

## 🔐 Security Best Practices

1. ✅ **Never commit `.env` files** to git
2. ✅ **Always verify webhook signatures**
3. ✅ **Use environment variables** for keys
4. ✅ **Separate test and live keys**
5. ✅ **Rotate keys periodically**
6. ✅ **Monitor for suspicious activity**
7. ✅ **Enable Stripe Radar** (fraud detection)
8. ✅ **Use HTTPS** in production
9. ✅ **Log all payment events**
10. ✅ **Implement rate limiting**

---

**🎉 You're all set! Your Stripe integration is ready to accept subscription payments!**

For support: [Stripe Support](https://support.stripe.com)
