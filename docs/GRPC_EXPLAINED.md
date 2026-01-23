# 🚀 gRPC - Як це працює?

## 📖 Що таке gRPC?

**gRPC** = **g**oogle **R**emote **P**rocedure **C**all

Це високопродуктивний фреймворк для комунікації між сервісами, створений Google.

---

## 🎯 Проста аналогія

**REST API:**
```
Клієнт: "Дай мені дані про користувача 123"
Сервер: "Ось JSON з даними"
```

**gRPC:**
```
Клієнт: userService.GetUser(123)
Сервер: повертає User object (бінарний формат)
```

**Різниця:** gRPC працює як звичайний виклик функції, але функція виконується на іншому сервісі!

---

## 🏗️ Архітектура вашого проекту

```
┌─────────────────────────────────────────────────────────────────┐
│                         FRONTEND                                │
│                    (React / TypeScript)                         │
│                                                                 │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                     │
│  │ Login    │  │ Products │  │ Bidding  │                     │
│  │ Page     │  │ Page     │  │ Page     │                     │
│  └──────────┘  └──────────┘  └──────────┘                     │
│       │              │              │                           │
└───────┼──────────────┼──────────────┼───────────────────────────┘
        │              │              │
        │         HTTP REST API       │
        │              │              │
        ▼              ▼              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API GATEWAY                                │
│                      (Port 8080)                                │
│                                                                 │
│  ┌────────────────────────────────────────────────────────────┐│
│  │  Gin HTTP Server                                           ││
│  │  - Handles REST requests from frontend                    ││
│  │  - JWT authentication                                      ││
│  │  - Routes to appropriate microservice                     ││
│  └────────────────────────────────────────────────────────────┘│
│                                                                 │
│       │                  │                  │                   │
└───────┼──────────────────┼──────────────────┼───────────────────┘
        │                  │                  │
        │    gRPC calls    │                  │
        │   (Binary/Fast)  │                  │
        ▼                  ▼                  ▼
┌────────────┐    ┌────────────┐    ┌────────────┐
│   USER     │    │  PRODUCT   │    │  BIDDING   │
│  SERVICE   │    │  SERVICE   │    │  SERVICE   │
│ Port 50051 │    │ Port 50052 │    │ Port 50053 │
│            │    │            │    │            │
│ ┌────────┐ │    │ ┌────────┐ │    │ ┌────────┐ │
│ │ gRPC   │ │    │ │ gRPC   │ │    │ │ gRPC   │ │
│ │ Server │ │    │ │ Server │ │    │ │ Server │ │
│ └────────┘ │    │ └────────┘ │    │ └────────┘ │
│     │      │    │     │      │    │     │      │
│     ▼      │    │     ▼      │    │     ▼      │
│ ┌────────┐ │    │ ┌────────┐ │    │ ┌────────┐ │
│ │Database│ │    │ │Database│ │    │ │Database│ │
│ │  Users │ │    │ │Products│ │    │ │Bids/Ask│ │
│ └────────┘ │    │ └────────┘ │    │ └────────┘ │
└────────────┘    └────────────┘    └────────────┘
```

---

## 🔄 Реальний приклад: User Login

### 1️⃣ Frontend → API Gateway (HTTP)

```typescript
// Frontend (BrowserScript)
const response = await fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST',
  body: JSON.stringify({
    email: 'user@test.com',
    password: 'password'
  })
});
```

### 2️⃣ API Gateway → User Service (gRPC)

```go
// API Gateway (internal/gateway/handlers/user_handler.go)
func (h *UserHandler) Login(c *gin.Context) {
    var req LoginRequest
    c.BindJSON(&req)
    
    // gRPC call to User Service
    response, err := h.userClient.Login(c.Request.Context(), &pb.LoginRequest{
        Email:    req.Email,
        Password: req.Password,
    })
    
    c.JSON(200, response)
}
```

### 3️⃣ User Service обробляє запит

```go
// User Service (internal/user/handler/grpc_handler.go)
func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
    // Validate credentials
    user, err := h.service.AuthenticateUser(req.Email, req.Password)
    
    // Generate JWT token
    token := generateToken(user)
    
    return &pb.LoginResponse{
        AccessToken: token,
        User: &pb.User{
            Id:    user.ID,
            Email: user.Email,
        },
    }, nil
}
```

### 4️⃣ Відповідь повертається назад

```
User Service → API Gateway → Frontend
(gRPC binary)  (JSON/HTTP)   (JavaScript)
```

---

## 📊 Детальна схема gRPC call

```
┌─────────────────────────────────────────────────────────────────┐
│                     gRPC COMMUNICATION                          │
└─────────────────────────────────────────────────────────────────┘

1. DEFINE SERVICE (Protocol Buffers - .proto file)
   ┌───────────────────────────────────────────────────────────┐
   │ syntax = "proto3";                                        │
   │                                                           │
   │ service UserService {                                     │
   │   rpc Login(LoginRequest) returns (LoginResponse);       │
   │ }                                                         │
   │                                                           │
   │ message LoginRequest {                                    │
   │   string email = 1;                                       │
   │   string password = 2;                                    │
   │ }                                                         │
   │                                                           │
   │ message LoginResponse {                                   │
   │   string access_token = 1;                                │
   │   User user = 2;                                          │
   │ }                                                         │
   └───────────────────────────────────────────────────────────┘

2. GENERATE CODE (protoc compiler)
   ┌───────────────────────────────────────────────────────────┐
   │ protoc --go_out=. --go-grpc_out=. user.proto            │
   └───────────────────────────────────────────────────────────┘
                           │
                           ▼
   ┌───────────────────────────────────────────────────────────┐
   │ Generated files:                                          │
   │ - user.pb.go         (message structs)                    │
   │ - user_grpc.pb.go    (client & server interfaces)        │
   └───────────────────────────────────────────────────────────┘

3. SERVER IMPLEMENTATION
   ┌───────────────────────────────────────────────────────────┐
   │ type UserHandler struct {                                 │
   │   pb.UnimplementedUserServiceServer                       │
   │ }                                                         │
   │                                                           │
   │ func (h *UserHandler) Login(ctx, req) (*pb.LoginResponse, error) {│
   │   // Your business logic here                            │
   │   return &pb.LoginResponse{...}, nil                      │
   │ }                                                         │
   │                                                           │
   │ // Start gRPC server                                      │
   │ grpcServer := grpc.NewServer()                            │
   │ pb.RegisterUserServiceServer(grpcServer, userHandler)    │
   │ grpcServer.Serve(listener)                                │
   └───────────────────────────────────────────────────────────┘

4. CLIENT CALLS
   ┌───────────────────────────────────────────────────────────┐
   │ // Connect to service                                     │
   │ conn, _ := grpc.Dial("localhost:50051", grpc.WithInsecure())│
   │ client := pb.NewUserServiceClient(conn)                   │
   │                                                           │
   │ // Call method (like a regular function!)                │
   │ response, err := client.Login(ctx, &pb.LoginRequest{     │
   │   Email: "user@test.com",                                 │
   │   Password: "password",                                   │
   │ })                                                        │
   └───────────────────────────────────────────────────────────┘

5. NETWORK TRANSMISSION (HTTP/2 + Protocol Buffers)
   ┌───────────────────────────────────────────────────────────┐
   │                                                           │
   │  Client                        Server                     │
   │    │                              │                       │
   │    │ 1. Serialize request         │                       │
   │    │    to binary (protobuf)      │                       │
   │    │──────────────────────────────>│                      │
   │    │                              │ 2. Deserialize       │
   │    │                              │ 3. Execute method    │
   │    │                              │ 4. Serialize response│
   │    │<──────────────────────────────│                      │
   │    │ 5. Deserialize response      │                       │
   │                                                           │
   └───────────────────────────────────────────────────────────┘
```

---

## 🆚 gRPC vs REST

| Feature | REST API | gRPC |
|---------|----------|------|
| **Format** | JSON (text) | Protocol Buffers (binary) |
| **Speed** | Slower | **10x faster** ⚡ |
| **Size** | Larger | **Smaller** (60-80% less) |
| **Protocol** | HTTP/1.1 | **HTTP/2** (multiplexing) |
| **Type Safety** | ❌ Runtime errors | ✅ **Compile-time** |
| **Streaming** | ❌ Limited | ✅ **Bi-directional** |
| **Browser Support** | ✅ Native | ❌ Needs proxy (gRPC-Web) |

---

## 🔍 В вашому проекті

### Файлова структура:

```
pkg/proto/
├── user.proto              # User Service definition
├── user.pb.go              # Generated code (messages)
├── user_grpc.pb.go         # Generated code (client/server)
├── product.proto
├── product.pb.go
├── product_grpc.pb.go
├── bidding.proto
├── bidding.pb.go
└── bidding_grpc.pb.go
```

### Приклад .proto файлу:

```protobuf
// pkg/proto/user.proto
syntax = "proto3";

package user;
option go_package = "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/user";

service UserService {
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc GetProfile(GetProfileRequest) returns (User);
}

message LoginRequest {
  string email = 1;
  string password = 2;
}

message LoginResponse {
  string access_token = 1;
  string refresh_token = 2;
  User user = 3;
}

message User {
  int64 id = 1;
  string email = 2;
  string first_name = 3;
  string last_name = 4;
  bool is_active = 5;
}
```

---

## 💡 Чому gRPC для мікросервісів?

### ✅ Переваги:

1. **Швидкість** ⚡
   - Binary format (Protocol Buffers)
   - HTTP/2 multiplexing
   - 10x швидше за JSON REST

2. **Type Safety** 🛡️
   - Помилки виявляються на етапі компіляції
   - IDE autocomplete
   - Generated client/server code

3. **Streaming** 🌊
   - Server streaming (WebSocket альтернатива)
   - Client streaming
   - Bidirectional streaming

4. **Language Agnostic** 🌍
   - Go, Python, Java, C++, JavaScript...
   - Один .proto файл → код для всіх мов

### ❌ Недоліки:

1. **Browser Support** 🌐
   - Браузер не підтримує gRPC напряму
   - Потрібен gRPC-Web або REST proxy (API Gateway)

2. **Debugging** 🐛
   - Binary format важче читати
   - Потрібні спеціальні інструменти (grpcurl, Postman)

3. **Learning Curve** 📚
   - Protocol Buffers синтаксис
   - Code generation workflow

---

## 🎯 Best Practices (як у вашому проекті)

### 1. API Gateway Pattern

```
Frontend (Browser)
    │
    │ HTTP/JSON (зручно для браузера)
    ▼
API Gateway (Port 8080)
    │
    │ gRPC (швидко між серверами)
    ▼
Microservices (Ports 50051-50056)
```

**Чому:**
- Frontend працює з REST (простіше)
- Між серверами gRPC (швидше)
- Кращий баланс між зручністю та продуктивністю

### 2. Versioning

```protobuf
syntax = "proto3";

package user.v1;  // Version in package name
```

### 3. Error Handling

```go
return nil, status.Errorf(codes.NotFound, "user not found: %v", userID)
```

### 4. Context для Timeouts

```go
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

response, err := client.Login(ctx, req)
```

---

## 📦 Приклад повного flow

### Scenario: User places BID

```
1. Frontend
   ↓ HTTP POST /api/v1/bids
   
2. API Gateway (port 8080)
   ↓ gRPC: biddingClient.PlaceBid()
   
3. Bidding Service (port 50053)
   ↓ Validate & Save to DB
   ↓ gRPC: notificationClient.NotifyMatchCreated()
   
4. Notification Service (port 50056)
   ↓ Send email via SMTP
   
5. Response flows back:
   Notification → Bidding → API Gateway → Frontend
```

**Code:**

```go
// API Gateway → Bidding Service (gRPC)
response, err := h.biddingClient.PlaceBid(c.Request.Context(), &pb.PlaceBidRequest{
    UserId:    userID,
    ProductId: req.ProductID,
    Price:     req.Price,
})

// Bidding Service → Notification Service (gRPC)
_, err = h.notificationClient.NotifyMatchCreated(ctx, &notificationPb.NotifyMatchCreatedRequest{
    MatchId:   match.ID,
    BuyerId:   match.BuyerID,
    SellerId:  match.SellerID,
    Price:     match.Price,
})
```

---

## 🛠️ Інструменти

### 1. Protocol Buffers Compiler

```bash
# Install protoc
brew install protobuf

# Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate code
protoc --go_out=. --go-grpc_out=. user.proto
```

### 2. Testing

```bash
# grpcurl - like curl for gRPC
brew install grpcurl

# Test login
grpcurl -plaintext \
  -d '{"email":"test@test.com","password":"password"}' \
  localhost:50051 \
  user.UserService/Login
```

### 3. GUI Tools

- **BloomRPC** - Postman for gRPC
- **gRPC UI** - Web interface

---

## 📚 Корисні посилання

- [gRPC Official Docs](https://grpc.io/docs/)
- [Protocol Buffers Guide](https://developers.google.com/protocol-buffers)
- [gRPC in Go Tutorial](https://grpc.io/docs/languages/go/quickstart/)

---

## 🎓 Підсумок

**gRPC** - це як виклик функції на іншому сервері:

```go
// Виглядає як звичайний виклик функції
response, err := userClient.Login(ctx, &LoginRequest{...})

// Але насправді:
// 1. Запит серіалізується в binary
// 2. Відправляється по мережі (HTTP/2)
// 3. Сервер десеріалізує
// 4. Виконує метод
// 5. Серіалізує відповідь
// 6. Відправляє назад
// 7. Клієнт десеріалізує
// Все це відбувається автоматично!
```

**Швидко. Type-safe. Просто у використанні.**

---

💡 Є питання? Напиши яку частину хочеш дослідити детальніше!
