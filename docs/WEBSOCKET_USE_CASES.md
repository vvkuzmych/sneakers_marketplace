# WebSocket - Де і чому потрібен? 🔌

## 🎯 Що таке WebSocket?

**WebSocket** - це двостороння комунікація між клієнтом і сервером в реальному часі.

### HTTP vs WebSocket:

```
HTTP (звичайний API):
Client → Request  → Server
Client ← Response ← Server
(кожен раз нове з'єднання)

WebSocket:
Client ↔ Server (постійне з'єднання)
```

---

## 💡 ДЕ ПОТРІБЕН У SNEAKERS MARKETPLACE?

### 1️⃣ **Bidding System (НАЙВАЖЛИВІШЕ!)** 🎯

**Ситуація:**
Ти на сторінці продукту (Air Jordan 1) і бачиш:
- Highest BID: $200 (хтось хоче купити за $200)
- Lowest ASK: $220 (хтось хоче продати за $220)

**Проблема без WebSocket:**
- Ці ціни можуть змінитися будь-коли
- Хтось може розмістити новий BID $210
- Хтось може розмістити новий ASK $215
- Може відбутися MATCH (покупка)
- Ти не побачиш це, поки не оновиш сторінку

**З WebSocket:**
```
User1: Places BID $210
  → WebSocket → All users see: "New BID: $210"

User2: Places ASK $215
  → WebSocket → All users see: "New ASK: $215"

User3: Places BID $220 (matches ASK $215!)
  → WebSocket → All users see: "MATCH! Sold at $215"
```

**Приклад коду (якби залишили websocket.ts):**

```typescript
// Connect to WebSocket
const ws = new WebSocket('ws://localhost:8080/ws');

// Listen for market updates
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  
  if (data.type === 'BID_UPDATED') {
    updateHighestBid(data.price); // $200 → $210
  }
  
  if (data.type === 'ASK_UPDATED') {
    updateLowestAsk(data.price); // $220 → $215
  }
  
  if (data.type === 'MATCH_CREATED') {
    showNotification('Match found! Price: $' + data.price);
  }
};
```

---

### 2️⃣ **Real-Time Notifications** 🔔

**Сценарій:**
1. Ти розмістив ASK: "Продам Air Jordan за $200"
2. Хтось розміщує BID: $205 (matcher!)
3. Match створено
4. **Сповіщення приходить миттєво через WebSocket**

```typescript
// WebSocket notification
{
  type: "MATCH_CREATED",
  message: "Your sneakers sold for $205!",
  match_id: 123,
  timestamp: "2026-01-21T20:00:00Z"
}
```

**Без WebSocket:**
- Треба постійно робити API запити (polling)
- Навантаження на сервер
- Затримка 5-10 секунд

**З WebSocket:**
- Миттєве сповіщення
- Мінімальне навантаження
- Real-time!

---

### 3️⃣ **Live Product Feed** 📊

**Сторінка "Hot Deals":**
- Нові продукти додаються
- Ціни змінюються
- Популярні продукти оновлюються

```typescript
ws.onmessage = (event) => {
  if (event.type === 'NEW_PRODUCT') {
    addProductToTop(event.product); // Новий продукт з'являється
  }
  
  if (event.type === 'PRICE_DROP') {
    highlightProduct(event.product_id); // "Price dropped!"
  }
};
```

---

### 4️⃣ **Order Status Updates** 📦

**Твоє замовлення:**
1. ✅ Order Placed
2. ✅ Payment Confirmed
3. ✅ Shipped
4. ✅ Delivered

**З WebSocket:**
Кожен статус приходить миттєво → оновлення на екрані без refresh!

```typescript
ws.onmessage = (event) => {
  if (event.type === 'ORDER_STATUS_UPDATED') {
    updateOrderStatus(event.order_id, event.status);
    showNotification(`Order ${event.status}!`);
  }
};
```

---

### 5️⃣ **Live Chat (Support)** 💬

**Якщо додамо підтримку:**
- Real-time чат з support
- Миттєві відповіді
- Без затримок

---

## 🏗️ Архітектура WebSocket у проекті

### Backend (вже є!):

```
Phase 3: WebSocket Integration
├── Notification Service (gRPC + WebSocket)
├── API Gateway (WebSocket proxy)
└── Redis Pub/Sub (broadcast до всіх клієнтів)
```

**Файл:** `cmd/notification-service/main.go`

```go
// WebSocket handler (simplified)
func handleWebSocket(c *gin.Context) {
    conn, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
    
    // Authenticate user
    token := c.Query("token")
    user := validateToken(token)
    
    // Subscribe to user's channel
    sub := redis.Subscribe("notifications:" + user.ID)
    
    // Send messages
    for msg := range sub.Channel() {
        conn.WriteJSON(msg)
    }
}
```

---

### Frontend (треба було б додати):

```typescript
// services/websocket.ts
class WebSocketClient {
  private ws: WebSocket;
  
  connect(token: string) {
    this.ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);
    
    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      this.handleMessage(data);
    };
  }
  
  handleMessage(data: any) {
    switch(data.type) {
      case 'BID_UPDATED':
        store.dispatch(updateBid(data));
        break;
      case 'ASK_UPDATED':
        store.dispatch(updateAsk(data));
        break;
      case 'MATCH_CREATED':
        store.dispatch(addMatch(data));
        showNotification('Match found!');
        break;
    }
  }
}
```

---

## 📊 Порівняння підходів

### ❌ Без WebSocket (HTTP Polling):

```typescript
// Кожні 5 секунд запит
setInterval(() => {
  fetch('/api/market-price')
    .then(res => res.json())
    .then(data => updatePrice(data));
}, 5000);
```

**Недоліки:**
- ❌ Затримка 5 секунд
- ❌ 12 запитів на хвилину
- ❌ Навантаження на сервер
- ❌ Марнування bandwidth
- ❌ Не real-time!

---

### ✅ З WebSocket:

```typescript
ws.onmessage = (event) => {
  updatePrice(event.data);
};
```

**Переваги:**
- ✅ Миттєво (0 затримки)
- ✅ Тільки коли є зміни
- ✅ Мінімальне навантаження
- ✅ Real-time!

---

## 🎯 Коли WebSocket в цьому проекті?

### Phase 3: WebSocket Integration ✅ (Backend готовий)

**Вже реалізовано (Backend):**
1. ✅ Notification Service з WebSocket
2. ✅ API Gateway WebSocket proxy
3. ✅ Redis Pub/Sub
4. ✅ JWT authentication для WebSocket

**Треба додати (Frontend):**
1. ❌ WebSocket client (`services/websocket.ts`) - **видалили під час очистки**
2. ❌ Redux integration для real-time updates
3. ❌ UI components для notifications
4. ❌ Bidding page з live updates

---

## 💡 Чи потрібен зараз?

### Для поточного функціоналу:

**Login, Register, ProductList** → ❌ WebSocket НЕ потрібен
- Це статичні сторінки
- Дані завантажуються один раз
- HTTP API достатньо

### Для Bidding System:

**Bid/Ask сторінка** → ✅ WebSocket ОБОВ'ЯЗКОВИЙ!
- Real-time ціни
- Live матчі
- Миттєві сповіщення

---

## 🚀 Як це виглядатиме для користувача?

### Сценарій: Покупка кросівок

1. **Відкриваєш продукт:**
   ```
   Air Jordan 1 Chicago
   
   Market Price:
   📈 Highest BID: $200 (live)
   📉 Lowest ASK: $220 (live)
   ```

2. **Хтось розміщує новий BID $210:**
   ```
   💚 New BID: $210 (animated update)
   ```

3. **Розміщуєш свій BID $225:**
   ```
   ⚡ INSTANT MATCH!
   🎉 You bought for $220!
   
   → Order created
   → Notification sent
   → Email sent
   ```

**Все це відбувається миттєво через WebSocket!**

---

## 📝 Підсумок

### WebSocket потрібен для:

1. ✅ **Bidding System** (найважливіше!)
2. ✅ **Real-time notifications**
3. ✅ **Live market prices**
4. ✅ **Order status updates**

### Зараз використовується?

❌ **НІ** - ми в Phase 5 (Frontend Basic)
- Login/Register - HTTP API
- ProductList - HTTP API

✅ **ТАК** - коли дійдемо до Bidding
- Phase 3 backend готовий
- Треба додати frontend WebSocket client

---

## 🤔 Чи варто було видаляти websocket.ts?

### Мої думки:

**За видалення:**
- ✅ Зараз не використовується
- ✅ Чистіший код
- ✅ Менше плутанини

**Проти видалення:**
- ❌ Доведеться створювати знову для Phase 3
- ❌ Backend вже готовий

### Рішення:

**Якщо плануєш Phase 3 (Bidding):**
→ Треба буде створити websocket.ts знову

**Якщо НЕ плануєш:**
→ Правильно видалили, не потрібен

---

## 💬 Хочеш побачити як це працює?

Можу:
1. Створити websocket.ts знову
2. Додати Bidding сторінку з live updates
3. Показати як це працює з backend

**Скажи якщо цікаво!** 🚀

---

**TLDR:** WebSocket потрібен для Bidding System (real-time ціни і матчі). Зараз не використовується, тому видалили. Коли робитимеш Bidding - треба буде додати знову.
