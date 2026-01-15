# 🚀 Getting Started - Sneakers Marketplace

## Ініціалізація GitHub репозиторію

### 1. Створи репозиторій на GitHub

Перейди на https://github.com/new та створи новий репозиторій:
- **Repository name:** `sneakers_marketplace`
- **Description:** "Production-ready microservices platform for sneaker trading with real-time auction system"
- **Visibility:** Public (або Private, як забажаєш)
- ❌ **НЕ** обирай "Initialize this repository with README"

### 2. Ініціалізуй локальний Git

```bash
cd /Users/vkuzm/GolandProjects/sneakers_marketplace

# Ініціалізувати git
git init

# Додати всі файли
git add .

# Перший коміт
git commit -m "🎉 Initial commit: Project structure and documentation

- Setup microservices architecture (9 services)
- Configure Docker Compose for infrastructure
- Add Makefile for development workflow
- Create comprehensive documentation
- Initialize Go module with github.com/vvkuzmych/sneakers_marketplace"

# Перейменувати branch на main
git branch -M main

# Додати remote
git remote add origin https://github.com/vvkuzmych/sneakers_marketplace.git

# Push
git push -u origin main
```

### 3. Перевір що все працює

```bash
# Перевірити remote
git remote -v

# Output має бути:
# origin  https://github.com/vvkuzmych/sneakers_marketplace.git (fetch)
# origin  https://github.com/vvkuzmych/sneakers_marketplace.git (push)
```

---

## 🏃 Швидкий Старт (Розробка)

### Prerequisites

Переконайся що встановлено:
- **Go 1.25+** - `go version`
- **Docker & Docker Compose** - `docker --version`
- **Make** - `make --version`

### Крок 1: Створи .env файл

```bash
# Скопіюй шаблон (коли буде створений)
cp .env.example .env

# Або створи вручну:
cat > .env << 'EOF'
DATABASE_URL=postgres://postgres:postgres@localhost:5432/sneakers_marketplace?sslmode=disable
REDIS_URL=redis://localhost:6379/0
KAFKA_BROKERS=localhost:9092
JWT_SECRET=your-super-secret-key-change-me
EOF
```

### Крок 2: Запусти інфраструктуру

```bash
# Запустити PostgreSQL, Redis, Kafka, Elasticsearch, etc.
make docker-up

# Перевірити що все запустилось
docker-compose ps

# Подивитись логи
make docker-logs
```

### Крок 3: Завантаж залежності

```bash
# Завантажити Go модулі
make deps

# Встановити protoc plugins
make proto-install
```

### Крок 4: (Майбутнє) Запусти міграції

```bash
# Коли будуть створені міграції:
make migrate-up
```

### Крок 5: (Майбутнє) Запусти сервіси

```bash
# В окремих терміналах:
make run-user-service
make run-product-service
make run-bidding-service
# ... etc
```

---

## 📁 Структура Проекту

```
sneakers_marketplace/
├── cmd/                        # Main applications (entry points)
│   ├── user-service/
│   ├── product-service/
│   ├── bidding-service/        ← Matching Engine 🔥
│   ├── order-service/
│   ├── payment-service/
│   ├── notification-service/
│   ├── search-service/
│   ├── analytics-service/
│   └── auth-service/
│
├── internal/                   # Private application code
│   ├── user/
│   │   ├── handler/           # HTTP/gRPC handlers
│   │   ├── service/           # Business logic
│   │   ├── repository/        # Database layer
│   │   └── model/             # Domain models
│   ├── product/
│   ├── bidding/
│   └── ...
│
├── pkg/                        # Public shared code
│   ├── proto/                 # gRPC definitions (.proto files)
│   ├── middleware/            # Shared middleware (auth, logging)
│   └── utils/                 # Helper functions
│
├── migrations/                # SQL migrations
│   ├── 000001_init.up.sql
│   └── 000001_init.down.sql
│
├── scripts/                   # Helper scripts
│   └── seed/                  # Database seeding
│
├── tests/
│   ├── integration/           # Integration tests
│   └── e2e/                   # End-to-end tests
│
├── docs/                      # Documentation
│   ├── ARCHITECTURE.md        ← ГОТОВО ✅
│   ├── DATABASE_SCHEMA.md     ← TODO
│   ├── API.md                 ← TODO
│   └── MATCHING_ENGINE.md     ← TODO
│
├── deployments/
│   ├── docker-compose.yml     ← ГОТОВО ✅
│   ├── kubernetes/
│   └── terraform/
│
├── README.md                  ← ГОТОВО ✅
├── Makefile                   ← ГОТОВО ✅
├── .gitignore                 ← ГОТОВО ✅
├── .gitattributes             ← ГОТОВО ✅
├── LICENSE                    ← ГОТОВО ✅
└── go.mod                     ← ГОТОВО ✅
```

---

## 🎯 Поточний Статус

### ✅ Завершено (Phase 0):
- [x] Структура проекту
- [x] README.md з повним описом
- [x] Architecture документація
- [x] Docker Compose для інфраструктури
- [x] Makefile для розробки
- [x] Git setup
- [x] Go module ініціалізація
- [x] LICENSE

### 📝 Наступні кроки (Phase 1 - Week 1):

#### 1. Створити базову інфраструктуру
- [ ] Додати SQL міграції (users, products таблиці)
- [ ] Створити базовий config пакет
- [ ] Додати logger (zerolog або zap)
- [ ] Налаштувати database connection pooling

#### 2. User Service (JWT Auth)
- [ ] Створити gRPC protobuf definitions
- [ ] Реалізувати Register/Login
- [ ] JWT generation & validation
- [ ] Password hashing (bcrypt)
- [ ] Unit tests

#### 3. Product Service (Catalog)
- [ ] gRPC protobuf definitions
- [ ] CRUD операції для products
- [ ] Size-based inventory
- [ ] Redis caching
- [ ] Unit tests

---

## 🛠️ Корисні команди

```bash
# Development
make run-user-service          # Запустити User Service
make build                     # Зібрати всі сервіси
make test                      # Запустити тести
make test-coverage             # Тести з coverage

# Docker
make docker-up                 # Запустити інфраструктуру
make docker-down               # Зупинити
make docker-logs               # Логи

# Database
make migrate-up                # Запустити міграції
make migrate-down              # Відкотити
make seed                      # Заповнити тестовими даними

# Protobuf
make proto                     # Генерувати Go код з .proto

# Monitoring
make prometheus                # Відкрити Prometheus
make grafana                   # Відкрити Grafana
make jaeger                    # Відкрити Jaeger

# Cleanup
make clean                     # Очистити build артефакти
make clean-all                 # Очистити все включно з Docker
```

---

## 📚 Документація

- [README.md](./README.md) - Огляд проекту ✅
- [ARCHITECTURE.md](./docs/ARCHITECTURE.md) - Детальна архітектура ✅
- [DATABASE_SCHEMA.md](./docs/DATABASE_SCHEMA.md) - Database design (TODO)
- [MATCHING_ENGINE.md](./docs/MATCHING_ENGINE.md) - Bid/Ask matching logic (TODO)
- [API.md](./docs/API.md) - API documentation (TODO)
- [DEVELOPMENT_PLAN.md](./docs/DEVELOPMENT_PLAN.md) - Week-by-week plan (TODO)

---

## 🎓 Навчальні ресурси

### Go
- [Effective Go](https://golang.org/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
- [Concurrency in Go (книга)](https://www.oreilly.com/library/view/concurrency-in-go/9781491941294/)

### gRPC
- [gRPC Go Tutorial](https://grpc.io/docs/languages/go/basics/)
- `/golang_practice/GRPC_GUIDE.md` ✅

### Microservices
- [Building Microservices (книга)](https://www.oreilly.com/library/view/building-microservices-2nd/9781492034018/)
- [Microservices Patterns (книга)](https://www.manning.com/books/microservices-patterns)

---

## 💡 Tips

1. **Працюй послідовно** - не намагайся зробити все одразу
2. **Пиши тести** - TDD допоможе уникнути багів
3. **Коміть часто** - невеликі коміти легше review
4. **Документуй рішення** - чому обрав саме цей підхід?
5. **Запитуй AI** - використовуй для code review та ідей

---

## 🆘 Troubleshooting

### Docker containers не запускаються

```bash
# Перевірити логи
docker-compose logs

# Очистити все і запустити заново
make docker-clean
make docker-up
```

### Port вже зайнятий

```bash
# Знайти процес що використовує port 5432 (PostgreSQL)
lsof -i :5432

# Вбити процес
kill -9 <PID>
```

### Go module проблеми

```bash
# Оновити залежності
go mod tidy

# Очистити cache
go clean -modcache
```

---

## 📞 Підтримка

**Project Maintainer:** vvkuzmych
- GitHub: [@vvkuzmych](https://github.com/vvkuzmych)
- Repository: [sneakers_marketplace](https://github.com/vvkuzmych/sneakers_marketplace)

---

**Готовий почати? Запускай `make setup` і вперед! 🚀**
