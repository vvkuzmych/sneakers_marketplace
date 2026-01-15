# 💳 Payment Service - Stripe Integration Guide

**Status:** Hybrid Mode (Demo + Real Stripe)  
**Version:** 1.0  
**Last Updated:** 2026-01-15

---

## 🎯 Overview

Payment Service підтримує **2 режими роботи**:
1. **Demo Mode** - для локальної розробки без Stripe
2. **Real Mode** - повна інтеграція зі Stripe API

---

## 🔧 Configuration

### 1. Environment Variables

Додай в `.env`:

```bash
# Stripe Configuration
STRIPE_MODE=demo              # "demo" or "real"
STRIPE_SECRET_KEY=sk_test_... # Your Stripe test secret key
STRIPE_PUBLISHABLE_KEY=pk_test_... # Your Stripe publishable key
STRIPE_WEBHOOK_SECRET=whsec_... # Webhook signing secret (optional for now)
```

### 2. Get Stripe Keys

**Option A: Stripe Dashboard**
1. Go to: https://dashboard.stripe.com/test/apikeys
2. Copy **Secret key** (starts with `sk_test_`)
3. Copy **Publishable key** (starts with `pk_test_`)

**Option B: Use Demo Mode**
- Set `STRIPE_MODE=demo`
- No keys needed!

---

## 🚀 Usage

### Demo Mode (Default)

```bash
# In .env
STRIPE_MODE=demo
```

**Features:**
- ✅ Creates fake PaymentIntents (`pi_demo_123`)
- ✅ Simulates successful payments
- ✅ Works without internet
- ✅ Perfect for development

**Limitations:**
- ❌ No real Stripe dashboard
- ❌ Can't test real card processing
- ❌ Webhooks don't work

---

### Real Stripe Mode

```bash
# In .env
STRIPE_MODE=real
STRIPE_SECRET_KEY=sk_test_your_actual_key
STRIPE_PUBLISHABLE_KEY=pk_test_your_actual_key
```

**Features:**
- ✅ Real Stripe PaymentIntents
- ✅ View in Stripe Dashboard
- ✅ Test with real test cards (4242 4242 4242 4242)
- ✅ Real refunds & transfers

**Requirements:**
- ✅ Stripe account (free)
- ✅ Internet connection
- ⏳ Webhook secret (optional, for production)

---

## 📋 API Methods

### 1. Create Payment Intent

**Demo Mode:**
```json
{
  "stripe_payment_intent_id": "pi_demo_123",
  "client_secret": "pi_demo_123_secret"
}
```

**Real Mode:**
```json
{
  "stripe_payment_intent_id": "pi_3ABC123...",
  "client_secret": "pi_3ABC123..._secret_xyz"
}
```

### 2. Confirm Payment

**Demo Mode:**
- Automatically succeeds
- Returns fake charge ID: `ch_demo_456`

**Real Mode:**
- Verifies with Stripe API
- Returns real charge ID
- Extracts card details (last4, brand)

### 3. Create Refund

**Demo Mode:**
- Returns fake refund ID: `re_demo_789`

**Real Mode:**
- Creates real Stripe refund
- Money returned to customer
- Visible in Stripe Dashboard

### 4. Create Payout (Stripe Connect)

**Demo Mode:**
- Returns fake transfer ID: `tr_demo_immediate`

**Real Mode:**
- Creates real Stripe Transfer
- Money sent to seller's Stripe account
- Requires seller to have Stripe Connect account

---

## 🧪 Testing

### Test Cards (Real Mode Only)

```
Success:
4242 4242 4242 4242  - Visa
5555 5555 5555 4444  - Mastercard

Decline:
4000 0000 0000 0002  - Generic decline
4000 0000 0000 9995  - Insufficient funds
```

**Expiry:** Any future date (e.g., 12/34)  
**CVC:** Any 3 digits (e.g., 123)  
**ZIP:** Any 5 digits (e.g., 12345)

---

## 🔄 Switching Modes

### From Demo → Real:

1. Update `.env`:
   ```bash
   STRIPE_MODE=real
   STRIPE_SECRET_KEY=sk_test_your_key
   ```

2. Restart Payment Service:
   ```bash
   pkill payment-service
   ./bin/payment-service
   ```

3. Test with real card!

### From Real → Demo:

1. Update `.env`:
   ```bash
   STRIPE_MODE=demo
   ```

2. Restart service

---

## 📊 Comparison

| Feature | Demo Mode | Real Mode |
|---------|-----------|-----------|
| Internet required | ❌ No | ✅ Yes |
| Stripe keys required | ❌ No | ✅ Yes |
| Development speed | ⚡ Fast | 🐢 Slower |
| Stripe Dashboard | ❌ No | ✅ Yes |
| Real card testing | ❌ No | ✅ Yes |
| Webhooks | ❌ No | ✅ Yes (with CLI) |
| Production ready | ❌ No | ✅ Yes |

---

## 🔐 Security

### Demo Mode:
- ✅ Safe - no real money
- ✅ No sensitive data

### Real Mode:
- ⚠️ Use **Test Mode** keys only!
- ⚠️ Never commit keys to git
- ⚠️ Use `.env` (in .gitignore)
- ✅ Stripe handles card data (PCI compliant)

---

## 🐛 Troubleshooting

### "Stripe API key not set"
```bash
# Check .env
STRIPE_MODE=real
STRIPE_SECRET_KEY=sk_test_...  # Must start with sk_test_
```

### "Payment Intent not found"
- In demo mode: This is OK
- In real mode: Check Stripe Dashboard

### "Invalid API Key"
- Verify key starts with `sk_test_`
- Regenerate key in Stripe Dashboard

---

## 🚀 Next Steps

### For Development:
- Use **Demo Mode** - fast & easy

### For Testing:
- Use **Real Mode** with test cards
- View payments in Stripe Dashboard

### For Production:
- Use **Real Mode** with live keys (`sk_live_`)
- Setup webhooks (need `STRIPE_WEBHOOK_SECRET`)
- Enable Stripe Connect for payouts

---

## 📝 Code Examples

### Check Current Mode

```go
service := NewPaymentService(repo)
mode := service.GetStripeMode()
fmt.Printf("Running in %s mode\n", mode) // "demo" or "real"
```

### Create Payment

```go
payment, clientSecret, err := service.CreatePaymentIntent(
    ctx, orderID, userID, 100.00, "USD", "",
)

// Demo: clientSecret = "pi_demo_123_secret"
// Real: clientSecret = "pi_3ABC..._secret_xyz"
```

---

**Made with ❤️ for Sneakers Marketplace**
