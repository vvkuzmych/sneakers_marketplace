# 📢 Phase 3 - Notification Service Architecture

**Goal:** Real-time notifications for users via Email and WebSocket

---

## 🎯 Overview

The Notification Service will:
1. Listen to events from other services (Kafka or direct gRPC calls)
2. Send email notifications (via Mailhog/SMTP)
3. Provide WebSocket for real-time browser notifications
4. Track notification history
5. Support notification preferences (email, push, in-app)

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Notification Service                     │
│                        (Port 50056)                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │   gRPC API   │  │   Kafka      │  │  WebSocket   │    │
│  │   Handler    │  │   Consumer   │  │    Server    │    │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘    │
│         │                  │                  │            │
│         └──────────────────┴──────────────────┘            │
│                           ↓                                │
│                  ┌─────────────────┐                       │
│                  │  Notification    │                       │
│                  │     Service      │                       │
│                  └────────┬─────────┘                       │
│                           │                                │
│         ┌─────────────────┴──────────────────┐            │
│         ↓                                     ↓            │
│  ┌──────────────┐                    ┌──────────────┐     │
│  │  Email       │                    │  WebSocket   │     │
│  │  Sender      │                    │  Publisher   │     │
│  │  (SMTP)      │                    │              │     │
│  └──────────────┘                    └──────────────┘     │
│                                                            │
└────────────────────────────────────────────────────────────┘
         ↓                                      ↓
    📧 Mailhog                           🌐 Browser
   (localhost:8025)                    WebSocket Client
```

---

## 📊 Database Schema

### Table: `notifications`

```sql
CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    type VARCHAR(50) NOT NULL,  -- 'match_created', 'order_created', 'order_shipped', etc.
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data JSONB,  -- Additional context (order_id, match_id, etc.)
    
    -- Channels
    email_sent BOOLEAN DEFAULT FALSE,
    email_sent_at TIMESTAMP,
    push_sent BOOLEAN DEFAULT FALSE,
    push_sent_at TIMESTAMP,
    
    -- Status
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Indexes
    INDEX idx_notifications_user_id (user_id),
    INDEX idx_notifications_type (type),
    INDEX idx_notifications_created_at (created_at DESC),
    INDEX idx_notifications_is_read (is_read)
);
```

### Table: `notification_preferences`

```sql
CREATE TABLE notification_preferences (
    user_id BIGINT PRIMARY KEY,
    
    -- Email preferences
    email_enabled BOOLEAN DEFAULT TRUE,
    email_match_created BOOLEAN DEFAULT TRUE,
    email_order_created BOOLEAN DEFAULT TRUE,
    email_order_shipped BOOLEAN DEFAULT TRUE,
    email_payment_received BOOLEAN DEFAULT TRUE,
    
    -- Push preferences
    push_enabled BOOLEAN DEFAULT TRUE,
    push_match_created BOOLEAN DEFAULT TRUE,
    push_order_updates BOOLEAN DEFAULT TRUE,
    
    -- In-app preferences
    inapp_enabled BOOLEAN DEFAULT TRUE,
    
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

---

## 🎬 Event Types & Triggers

### 1. **Match Created** 🎯

**Trigger:** Bidding Service creates a match  
**Recipients:** Buyer + Seller  
**Email Subject:** "Your bid/ask has been matched!"

**Buyer Email:**
```
Congratulations! Your bid for Nike Air Jordan 1 (US 9) 
has been matched at $220.

Order #ORD-20260115-0001 has been created.
Please complete payment within 24 hours.

View Order: [Link]
```

**Seller Email:**
```
Good news! Your ask for Nike Air Jordan 1 (US 9) 
has been matched at $220.

Please prepare the item for shipment.
You'll receive payout after delivery confirmation.

View Order: [Link]
```

---

### 2. **Order Status Changes** 📦

**Trigger:** Order Service updates status  
**Recipients:** Buyer and/or Seller

**Status Transitions:**
- `pending` → `paid`: Notify **Seller** (prepare for shipment)
- `paid` → `processing`: Notify **Seller** (reminder)
- `processing` → `shipped`: Notify **Buyer** (tracking number)
- `shipped` → `in_transit`: Notify **Buyer** (tracking updates)
- `in_transit` → `delivered`: Notify **Buyer & Seller** (confirm delivery)
- `delivered` → `completed`: Notify **Seller** (payout initiated)

---

### 3. **Payment Events** 💳

**Trigger:** Payment Service events  
**Recipients:** Buyer or Seller

**Events:**
- **Payment Successful:** Notify Buyer (receipt)
- **Payment Failed:** Notify Buyer (retry)
- **Refund Issued:** Notify Buyer (refund details)
- **Payout Completed:** Notify Seller (funds transferred)

---

### 4. **Bid/Ask Expiration** ⏰

**Trigger:** Scheduled job (every hour)  
**Recipients:** User who placed bid/ask

**Email:**
```
Your bid for Nike Air Jordan 1 (US 9) at $200 
has expired without a match.

You can place a new bid at any time.

Place New Bid: [Link]
```

---

## 🔌 Integration Methods

### Option 1: **Kafka Events** (Recommended for Production)

**Pros:**
- ✅ Decoupled services
- ✅ Event replay capability
- ✅ Scalable (multiple consumers)
- ✅ Async processing

**Flow:**
```
Bidding Service → Kafka Topic: "matches"
Order Service → Kafka Topic: "orders"
Payment Service → Kafka Topic: "payments"
                    ↓
        Notification Service (Kafka Consumer)
```

---

### Option 2: **Direct gRPC Calls** (Simpler for MVP)

**Pros:**
- ✅ Simple implementation
- ✅ Immediate feedback
- ✅ No additional infrastructure

**Flow:**
```
Bidding Service → grpc.NotificationService.NotifyMatchCreated()
Order Service → grpc.NotificationService.NotifyOrderUpdate()
Payment Service → grpc.NotificationService.NotifyPaymentEvent()
```

**Decision:** Start with **Option 2** (gRPC), migrate to Kafka in Phase 4.

---

## 📧 Email Service

### SMTP Configuration

Using **Mailhog** (already in docker-compose):

```go
type EmailConfig struct {
    Host     string // "localhost"
    Port     int    // 1025 (Mailhog SMTP)
    From     string // "noreply@sneakersmarketplace.com"
    Username string // "" (optional)
    Password string // "" (optional)
}
```

### Email Templates

Using Go's `html/template`:

```
internal/notification/templates/
├── match_created_buyer.html
├── match_created_seller.html
├── order_shipped.html
├── payment_received.html
└── payout_completed.html
```

---

## 🌐 WebSocket Service

### Connection Management

```go
type WebSocketHub struct {
    clients    map[int64]*Client  // userID -> connection
    broadcast  chan *Notification
    register   chan *Client
    unregister chan *Client
}
```

### Client Connection

```
WebSocket URL: ws://localhost:8080/ws

Authentication: JWT token in query param or header
Example: ws://localhost:8080/ws?token=eyJhbGci...
```

### Message Format

```json
{
  "type": "notification",
  "data": {
    "id": "123",
    "title": "Your bid has been matched!",
    "message": "Order #ORD-20260115-0001 created",
    "link": "/orders/1",
    "timestamp": "2026-01-15T20:00:00Z"
  }
}
```

---

## 🔔 gRPC API

### Proto Definition

```protobuf
service NotificationService {
  // Send notification
  rpc SendNotification(SendNotificationRequest) returns (SendNotificationResponse);
  
  // Get user notifications
  rpc GetNotifications(GetNotificationsRequest) returns (GetNotificationsResponse);
  
  // Mark as read
  rpc MarkAsRead(MarkAsReadRequest) returns (MarkAsReadResponse);
  
  // Get preferences
  rpc GetPreferences(GetPreferencesRequest) returns (GetPreferencesResponse);
  
  // Update preferences
  rpc UpdatePreferences(UpdatePreferencesRequest) returns (UpdatePreferencesResponse);
  
  // Batch notify (for match events)
  rpc NotifyMatchCreated(NotifyMatchCreatedRequest) returns (NotifyMatchCreatedResponse);
  rpc NotifyOrderUpdate(NotifyOrderUpdateRequest) returns (NotifyOrderUpdateResponse);
  rpc NotifyPaymentEvent(NotifyPaymentEventRequest) returns (NotifyPaymentEventResponse);
}
```

---

## 📊 Implementation Plan

### Phase 3.1 - Core Notification Service (Week 1)

1. ✅ Database migration (notifications, preferences)
2. ✅ Proto definition
3. ✅ Models & Repository
4. ✅ Email sender (SMTP/Mailhog)
5. ✅ gRPC handler
6. ✅ Integration with Bidding/Order/Payment services

### Phase 3.2 - WebSocket (Week 1-2)

1. ✅ WebSocket hub implementation
2. ✅ JWT authentication for WebSocket
3. ✅ Client connection management
4. ✅ Real-time message broadcasting
5. ✅ API Gateway WebSocket endpoint

### Phase 3.3 - Advanced Features (Week 2)

1. ✅ Email templates (HTML)
2. ✅ Notification preferences UI (via API Gateway)
3. ✅ Batch notifications
4. ✅ Notification history pagination
5. ✅ Mark all as read

---

## 🧪 Testing Strategy

**Email Testing:**
- View emails in Mailhog UI: http://localhost:8025

**WebSocket Testing:**
- Use `websocat` or browser console
- Test reconnection logic
- Test JWT authentication

**Integration Testing:**
- Create match → verify emails sent
- Update order status → verify notifications
- Payment event → verify buyer/seller notified

---

## 🚀 Quick Start (After Implementation)

```bash
# Start Notification Service
./bin/notification-service

# View emails in Mailhog
open http://localhost:8025

# Test WebSocket (browser console)
const ws = new WebSocket('ws://localhost:8080/ws?token=YOUR_JWT');
ws.onmessage = (event) => console.log(JSON.parse(event.data));
```

---

**Ready to implement!** 🎉
