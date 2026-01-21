# WebSocket + Bidding Page - Testing Guide 🧪

## 🚀 Quick Start

### 1. Ensure Backend Services Are Running

```bash
cd /Users/vkuzm/GolandProjects/sneakers_marketplace

# Check running services
ps aux | grep -E "(notification-service|api-gateway|bidding-service)"

# If not running, start them:
./bin/notification-service &
./bin/api-gateway &
./bin/bidding-service &
```

### 2. Open Frontend

Frontend should already be running on:
```
http://localhost:5173
```

If not:
```bash
cd /Users/vkuzm/GolandProjects/sneakers_marketplace/frontend
npm run dev
```

---

## 🧪 Test Scenarios

### Scenario 1: Basic WebSocket Connection

1. Open browser: `http://localhost:5173/login`
2. Login with credentials:
   - Email: `test@example.com`
   - Password: `password123`
3. Click on any product
4. You should see:
   - WebSocket status: **Connected** (green ●)
   - Console: `✅ WebSocket ready for real-time updates`

**Expected Result:** ✅ Green "Connected" indicator

---

### Scenario 2: Place BID (No Match)

1. Navigate to any product's bidding page
2. Current market:
   - Highest BID: —
   - Lowest ASK: —
3. Place BID: **$200**
4. Click "Place BID"

**Expected Result:**
- ✅ Notification: "BID placed at $200"
- ✅ Highest BID updates to: $200
- ✅ Console: `✅ BID placed: { bid: { price: 200, ... } }`

---

### Scenario 3: Place ASK (No Match)

1. Still on the same product
2. Current market:
   - Highest BID: $200
   - Lowest ASK: —
3. Place ASK: **$220**
4. Click "Place ASK"

**Expected Result:**
- ✅ Notification: "ASK placed at $220"
- ✅ Lowest ASK updates to: $220
- ✅ Spread shows: $20
- ✅ Console: `✅ ASK placed: { ask: { price: 220, ... } }`

---

### Scenario 4: Instant Match

1. Current market:
   - Highest BID: $200
   - Lowest ASK: $220
2. Place BID: **$225** (higher than ASK!)
3. Click "Place BID"

**Expected Result:**
- 🎉 Notification: "INSTANT MATCH! Bought at $220"
- ✅ Match created
- ✅ Console: `⚡ MATCH CREATED: { match: { price: 220, ... } }`
- ✅ Form hint shows: "⚡ This will create INSTANT MATCH!" (when typing $225)

---

### Scenario 5: Real-Time Updates (Multi-User)

**This is the MOST IMPORTANT test for WebSocket!**

1. Open **TWO browser tabs** (or use Incognito mode)
2. Login in both tabs
3. Navigate to the **SAME product** in both tabs

**Tab 1:**
- Place BID: $200

**Tab 2:**
- Should immediately see: "New BID: $200"
- Highest BID updates to: $200

**Tab 2:**
- Place ASK: $200 (equal to BID)

**Tab 1:**
- Should immediately see: "MATCH! Sold at $200"
- 🎉 Match notification

**Expected Result:**
- ✅ Both tabs show real-time updates
- ✅ No page refresh needed
- ✅ Instant synchronization

---

## 🐛 Troubleshooting

### Problem: WebSocket shows "Disconnected" (red ●)

**Possible causes:**
1. ❌ Notification Service not running
2. ❌ API Gateway not running
3. ❌ JWT token expired or invalid
4. ❌ WebSocket upgrade failed

**Solutions:**
```bash
# Check if services are running
lsof -i :8080  # API Gateway
ps aux | grep notification-service

# Restart services
killall notification-service api-gateway
./bin/notification-service &
./bin/api-gateway &

# Clear localStorage and re-login
localStorage.clear()
```

---

### Problem: "Failed to place BID"

**Possible causes:**
1. ❌ Bidding Service not running
2. ❌ Invalid product/size ID
3. ❌ Invalid price (negative, zero)

**Solutions:**
```bash
# Check Bidding Service
lsof -i :50053  # Bidding Service gRPC

# Check logs
tail -f logs/bidding-service.log

# Restart
killall bidding-service
./bin/bidding-service &
```

---

### Problem: No Real-Time Updates

**Possible causes:**
1. ❌ WebSocket not connected
2. ❌ Redis not running (Pub/Sub)
3. ❌ Message handlers not registered

**Solutions:**
```bash
# Check Redis
redis-cli ping  # Should return "PONG"

# Check WebSocket in browser console
websocketService.isConnected()  # Should return true
websocketService.getState()     # Should return "OPEN"

# Re-connect WebSocket
websocketService.disconnect()
websocketService.connect(token)
```

---

## 📊 Console Debugging

### Check WebSocket Messages

Open browser console (F12) and watch for:

```javascript
// Connection
🔌 Connecting to WebSocket: ws://localhost:8080/ws?token=...
✅ WebSocket connected
📡 Subscribed to "BID_PLACED" messages
📡 Subscribed to "ASK_PLACED" messages
📡 Subscribed to "MATCH_CREATED" messages

// Messages
📨 WebSocket message: { type: "BID_PLACED", data: { ... } }
💰 New BID placed: { bid: { price: 200 } }

📨 WebSocket message: { type: "MATCH_CREATED", data: { ... } }
⚡ MATCH CREATED: { match: { price: 220 } }
```

---

## 🎯 Success Criteria

Your WebSocket + Bidding implementation is working if:

- ✅ WebSocket connects automatically on page load
- ✅ Green "Connected" indicator shows
- ✅ BID/ASK placement works
- ✅ Market prices update in real-time
- ✅ Instant match detection works
- ✅ Notifications appear and auto-hide
- ✅ Multi-tab synchronization works
- ✅ Auto-reconnect on disconnect

---

## 🔧 Advanced Testing

### Test WebSocket Manually

```javascript
// In browser console
const ws = new WebSocket('ws://localhost:8080/ws?token=YOUR_TOKEN');

ws.onopen = () => console.log('Connected');
ws.onmessage = (e) => console.log('Message:', JSON.parse(e.data));
ws.onerror = (e) => console.error('Error:', e);
ws.onclose = () => console.log('Disconnected');

// Send test message
ws.send(JSON.stringify({
  type: 'TEST',
  data: { message: 'Hello WebSocket!' }
}));
```

### Check Redis Pub/Sub

```bash
# Terminal 1: Subscribe to all channels
redis-cli
> PSUBSCRIBE *

# Terminal 2: Publish test message
redis-cli
> PUBLISH notifications:1 '{"type":"BID_PLACED","data":{"price":200}}'

# Terminal 1 should show the message
```

---

## 📈 Performance Testing

### Load Test (Multiple Connections)

1. Open 10+ browser tabs
2. Navigate to the same product
3. Place BIDs from different tabs rapidly
4. All tabs should show updates

**Expected:**
- ✅ All tabs stay synchronized
- ✅ No missed updates
- ✅ Latency < 100ms

---

## 🎉 What to Expect

When everything works:

1. **Instant Feedback**
   - BID placed → immediate confirmation
   - Market prices update in < 50ms

2. **Real-Time Sync**
   - Multiple users see same data
   - No polling, no refresh

3. **Smooth UX**
   - Animations for notifications
   - Pulse animation for status
   - Gradient backgrounds

4. **Reliable Connection**
   - Auto-reconnect on disconnect
   - Connection state always visible

---

## 📚 References

- **WebSocket Service:** `src/services/websocket.ts`
- **Bidding Page:** `src/features/bidding/BiddingPage.tsx`
- **Backend WebSocket:** `cmd/notification-service/main.go`
- **Backend API Gateway:** `cmd/api-gateway/main.go`

---

**Happy Testing! 🚀**
