#!/bin/bash

set -e

echo "🔔 Testing Notification Service..."
echo ""

# Test 1: Send notification
echo "1️⃣ Send Notification..."
grpcurl -plaintext -d '{
  "user_id": 1,
  "type": "match_created",
  "title": "Your bid has been matched!",
  "message": "Your bid for Nike Air Jordan 1 (US 9) has been matched at $220",
  "send_email": true,
  "send_push": false
}' localhost:50056 notification.NotificationService/SendNotification

echo ""

# Test 2: Get notifications
echo "2️⃣ Get User Notifications..."
grpcurl -plaintext -d '{
  "user_id": 1,
  "page": 1,
  "page_size": 10
}' localhost:50056 notification.NotificationService/GetNotifications

echo ""

# Test 3: Get unread count
echo "3️⃣ Get Unread Count..."
grpcurl -plaintext -d '{
  "user_id": 1
}' localhost:50056 notification.NotificationService/GetUnreadCount

echo ""

# Test 4: Get preferences
echo "4️⃣ Get Notification Preferences..."
grpcurl -plaintext -d '{
  "user_id": 1
}' localhost.NotificationService/GetPreferences

echo ""

# Test 5: Notify match created (event notification)
echo "5️⃣ Notify Match Created (buyer + seller)..."
grpcurl -plaintext -d '{
  "match_id": 1,
  "buyer_id": 1,
  "seller_id": 2,
  "product_id": 1,
  "product_name": "Nike Air Jordan 1 Retro High OG",
  "size": "US 9",
  "price": 220
}' localhost:50056 notification.NotificationService/NotifyMatchCreated

echo ""

echo "✅ Notification Service Test Complete! 🎉"
echo ""
echo "📧 Check Mailhog for emails: http://localhost:8025"
echo ""
