# 🧪 Testing Phase 3 - Notification Service + WebSocket

**Comprehensive testing guide for all Phase 3 features**

---

## 📋 Pre-requisites Checklist

Before testing, ensure:

- ✅ Docker containers running (`make docker-up`)
- ✅ Database migrations applied (`migrate up`)
- ✅ `.env` file configured with all variables

---

## 🚀 Step 1: Start All Services

### Terminal 1: User Service
```bash
cd /Users/vkuzm/GolandProjects/sneakers_marketplace
export $(cat .env | grep -v '^#' | xargs)
./bin/user-service
```

**Expected output:**
```
✅ Connected to PostgreSQL
✅ User Service listening on :50051
```

### Terminal 2: Product Service
```bash
cd /Users/vkuzm/GolandProjects/sneakers_marketplace
export $(cat .env | grep -v '^#' | xargs)
./bin/product-service
```

### Terminal 3: Bidding Service
```bash
cd /Users/vkuzm/GolandProjects/sneakers_marketplace
export $(cat .env | grep -v '^#' | xargs)
./bin/bidding-service
```

### Terminal 4: Order Service
```bash
cd /Users/vkuzm/GolandProjects/sneakers_marketplace
export $(cat .env | grep -v '^#' | xargs)
./bin/order-service
```

### Terminal 5: Payment Service
```bash
cd /Users/vkuzm/GolandProjects/sneakers_marketplace
export $(cat .env | grep -v '^#' | xargs)
./bin/payment-service
```

### Terminal 6: Notification Service ⭐ NEW
```bash
cd /Users/vkuzm/GolandProjects/sneakers_marketplace
export $(cat .env | grep -v '^#' | xargs)
./bin/notification-service
```

**Expected output:**
```
Starting Notification Service
✅ Connected to PostgreSQL
✅ Notification Service listening on port 50056
```

### Terminal 7: API Gateway (with WebSocket) ⭐ UPDATED
```bash
cd /Users/vkuzm/GolandProjects/sneakers_marketplace
export $(cat .env | grep -v '^#' | xargs)
./bin/api-gateway
```

**Expected output:**
```
✅ Connected to all gRPC services
✅ WebSocket Hub started
✅ API Gateway listening on :8080
```

---

## 🔍 Step 2: Verify Services are Running

### Check Service Ports

```bash
# Check all services
lsof -i :50051 | grep LISTEN  # User Service
lsof -i :50052 | grep LISTEN  # Product Service
lsof -i :50053 | grep LISTEN  # Bidding Service
lsof -i :50054 | grep LISTEN  # Order Service
lsof -i :50055 | grep LISTEN  # Payment Service
lsof -i :50056 | grep LISTEN  # Notification Service ⭐
lsof -i :8080  | grep LISTEN  # API Gateway
```

**All should return process information.**

### Check API Gateway Health (includes WebSocket status)

```bash
curl http://localhost:8080/health
```

**Expected response:**
```json
{
  "status": "healthy",
  "service": "api-gateway",
  "ws_connections": 0
}
```

---

## 📧 Step 3: Test Email Notifications

### 3.1 Open Mailhog UI

```bash
open http://localhost:8025
```

**You should see Mailhog interface (no emails yet)**

### 3.2 Send Test Notification via gRPC

```bash
grpcurl -plaintext -d '{
  "user_id": 1,
  "type": "match_created",
  "title": "Test Notification",
  "message": "This is a test email notification!",
  "send_email": true,
  "send_push": false
}' localhost:50056 notification.NotificationService/SendNotification
```

**Expected response:**
```json
{
  "notification": {
    "id": "1",
    "user_id": "1",
    "type": "match_created",
    "title": "Test Notification",
    "message": "This is a test email notification!",
    "email_sent": false,
    "push_sent": false,
    "is_read": false,
    "created_at": "2026-01-15T..."
  }
}
```

### 3.3 Check Email in Mailhog

Refresh Mailhog (http://localhost:8025) - you should see an email:

**From:** `noreply@sneakersmarketplace.com`  
**To:** `user1@example.com`  
**Subject:** Test Notification  
**Body:** This is a test email notification!

✅ **Email notifications working!**

---

## 🌐 Step 4: Test WebSocket Real-Time Notifications

### 4.1 Get JWT Token

```bash
# Register or login
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' | jq -r '.access_token')

echo $TOKEN
```

**Copy the token (it will look like: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`)**

### 4.2 Open WebSocket Test Page

```bash
open test_websocket.html
```

**In the browser:**
1. Paste your JWT token in the input field
2. Click "Connect"
3. You should see:
   - Status: "Connected" (green)
   - Message: "System: Connected to WebSocket server"
   - Message: "System: Welcome! User ID: 1"

✅ **WebSocket connection established!**

### 4.3 Test Ping/Pong

In the browser, click "Send Ping"

**Expected:**
- Message: "Sent: Ping"
- Message: "System: Pong received"

✅ **Ping/Pong working!**

### 4.4 Send Real-Time Notification

**Keep WebSocket page open**, then in terminal:

```bash
grpcurl -plaintext -d '{
  "user_id": 1,
  "type": "match_created",
  "title": "🎯 Real-Time Test!",
  "message": "This notification should appear instantly in your browser!",
  "send_email": false,
  "send_push": true
}' localhost:50056 notification.NotificationService/SendNotification
```

**In WebSocket browser page, you should INSTANTLY see:**

```
🔔 Notification:
Real-Time Test!
This notification should appear instantly in your browser!
[timestamp]
```

✅ **Real-time WebSocket notifications working!**

---

## 🎯 Step 5: Test Full Flow (Match Created → Notifications)

### 5.1 Create a Match (via Bidding Service)

```bash
# This will create a match and trigger notifications
./scripts/test_bidding_service.sh
```

**Expected:**
1. Match created between buyer and seller
2. **Email sent to buyer:** "Your bid has been matched!"
3. **Email sent to seller:** "Your ask has been matched!"
4. **WebSocket notifications** sent to both (if connected)

### 5.2 Check Mailhog

Open http://localhost:8025 - you should see **2 new emails**:

1. **Buyer email** - "Your bid has been matched at $..."
2. **Seller email** - "Your ask has been matched at $..."

✅ **Match notification flow working!**

---

## 📊 Step 6: Test Notification Features

### 6.1 Get User Notifications

```bash
grpcurl -plaintext -d '{
  "user_id": 1,
  "page": 1,
  "page_size": 10
}' localhost:50056 notification.NotificationService/GetNotifications
```

**Expected:** List of all notifications for user

### 6.2 Get Unread Count

```bash
grpcurl -plaintext -d '{
  "user_id": 1
}' localhost:50056 notification.NotificationService/GetUnreadCount
```

**Expected:** `{"count": "3"}` (or however many you created)

### 6.3 Mark Notification as Read

```bash
grpcurl -plaintext -d '{
  "notification_id": 1,
  "user_id": 1
}' localhost:50056 notification.NotificationService/MarkAsRead
```

**Expected:** `{"success": true}`

### 6.4 Get Notification Preferences

```bash
grpcurl -plaintext -d '{
  "user_id": 1
}' localhost:50056 notification.NotificationService/GetPreferences
```

**Expected:** User's notification preferences (all enabled by default)

---

## 🔥 Step 7: Stress Test WebSocket

### 7.1 Open Multiple WebSocket Connections

Open `test_websocket.html` in **3 different browser tabs**:

1. Tab 1: Connect with user 1's token
2. Tab 2: Connect with user 2's token (different user)
3. Tab 3: Connect with user 1's token (same user)

**Expected:**
- Tab 3 connects successfully
- Tab 1 disconnects (only one connection per user)
- Tab 2 stays connected

Check API Gateway health:
```bash
curl http://localhost:8080/health
```

**Expected:** `"ws_connections": 2`

✅ **Multiple connections managed correctly!**

### 7.2 Broadcast to All

```bash
# Send notification to user 1
grpcurl -plaintext -d '{
  "user_id": 1,
  "type": "test",
  "title": "Test for User 1",
  "message": "This goes to user 1 only",
  "send_push": true
}' localhost:50056 notification.NotificationService/SendNotification
```

**Expected:** Only user 1's browser receives notification

---

## 🧹 Step 8: Test Cleanup

### 8.1 Disconnect WebSocket

In browser, click "Disconnect"

**Expected:**
- Status: "Disconnected" (red)
- Message: "System: Disconnected from server"

Check health again:
```bash
curl http://localhost:8080/health
```

**Expected:** `"ws_connections": 1` (or 0 if all disconnected)

---

## 📝 Step 9: Run Automated Test Script

```bash
./scripts/test_notification_service.sh
```

**This will test:**
- Send notification
- Get notifications
- Get unread count
- Get preferences
- Notify match created

---

## ✅ Success Criteria

All tests pass if:

| Feature | Status |
|---------|--------|
| Notification Service running | ✅ Port 50056 |
| Email sending (Mailhog) | ✅ Emails received |
| WebSocket connection | ✅ Connected with JWT |
| Real-time notifications | ✅ Instant delivery |
| Ping/Pong keep-alive | ✅ Working |
| Multiple connections | ✅ One per user |
| Match notification flow | ✅ Buyer + Seller notified |
| Notification preferences | ✅ Get/Update working |
| Mark as read | ✅ Working |

---

## 🐛 Troubleshooting

### Email not received in Mailhog

**Check:**
```bash
# Mailhog running?
lsof -i :1025  # SMTP port
lsof -i :8025  # Web UI port

# Check Notification Service logs for errors
```

### WebSocket connection fails

**Check:**
1. JWT token valid? (not expired)
2. API Gateway running?
3. Browser console for errors (F12)

**Fix:**
```bash
# Get fresh token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### Notification not appearing in WebSocket

**Check:**
1. WebSocket still connected? (green status)
2. Correct user_id in notification?
3. `send_push: true` in request?

**Debug:**
```bash
# Check API Gateway logs (should show "Sent message to UserID=...")
# Check Notification Service logs
```

---

## 🎉 All Tests Passing?

**Congratulations!** Phase 3 is fully working:

✅ Notification Service  
✅ Email Notifications  
✅ WebSocket Real-Time  
✅ User Preferences  
✅ Full Integration

**Ready for production! 🚀**

---

## 📚 Next Steps

- Read `docs/WEBSOCKET_GUIDE.md` for integration examples
- Read `docs/PHASE_3_ARCHITECTURE.md` for architecture details
- Integrate notifications into your frontend
- Add HTML email templates (optional)
- Setup Stripe webhooks for real Payment events (optional)

**Made with ❤️ for Sneakers Marketplace**
