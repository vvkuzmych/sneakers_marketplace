# ⚡ GORM Quick Start Guide

**Швидкий старт для тестування GORM паралельно з існуючим кодом**

---

## 🎯 Що створено?

Ми створили **паралельний GORM пакет** без зміни основного коду:

```
sneakers_marketplace/
├── internal/user/
│   ├── model/              # ✅ Існуючий (Raw SQL)
│   ├── model_gorm/         # 🆕 Новий (GORM)
│   ├── repository/         # ✅ Існуючий (Raw SQL)
│   └── repository_gorm/    # 🆕 Новий (GORM)
├── examples/gorm_vs_raw/   # 🆕 Демо і бенчмарки
│   ├── main.go
│   ├── benchmark_test.go
│   └── README.md
└── docs/
    ├── GORM_INVESTIGATION.md  # Детальний аналіз
    └── GORM_QUICKSTART.md     # Цей файл
```

**✨ Основний код не змінювався! Це окремі пакети для порівняння.**

---

## 🚀 Як запустити?

### 1️⃣ Запустити через Shell Script (Найпростіше) ⭐

```bash
cd /Users/vkuzm/GolandProjects/sneakers_marketplace/examples/gorm_vs_raw

# Інтерактивне меню
./run_comparison.sh

# Або прямо:
./run_comparison.sh demo      # Demo
./run_comparison.sh bench     # Benchmarks
./run_comparison.sh all       # Все разом

# Швидкі shortcuts:
./demo.sh     # Demo
./bench.sh    # Benchmarks
```

**Що робить скрипт:**
- ✅ Перевіряє підключення до PostgreSQL
- ✅ Автоматично будує проєкт
- ✅ Красиво форматує вивід
- ✅ Аналізує benchmarks (показує overhead %)
- ✅ Підсвічує результати кольорами

### 2️⃣ Запустити через Go (Альтернатива)

```bash
cd /Users/vkuzm/GolandProjects/sneakers_marketplace/examples/gorm_vs_raw
go run main.go
```

**Що він робить:**
- ✅ Створює користувачів (Raw SQL vs GORM)
- ✅ Читає з бази (GetByEmail)
- ✅ Оновлює користувачів
- ✅ Показує pagination
- ✅ Демонструє soft delete
- ✅ Порівнює продуктивність

**Очікуваний результат:**
```
╔══════════════════════════════════════════════════════════════════╗
║         🔬 GORM vs Raw SQL Comparison Demo                      ║
╚══════════════════════════════════════════════════════════════════╝

📦 Setting up Raw SQL (pgx) connection...
✅ Raw SQL repository ready

📦 Setting up GORM connection...
✅ GORM repository ready

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📝 Demo 1: CREATE USER
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🔹 Raw SQL (pgx):
✅ Created user ID: 42 (took 2.5ms)

🔹 GORM:
✅ Created user ID: 43 (took 3.8ms)

📊 Performance: Raw SQL 2.5ms vs GORM 3.8ms (1.5x)

... (більше тестів)

╔══════════════════════════════════════════════════════════════════╗
║                        📊 SUMMARY                                ║
╚══════════════════════════════════════════════════════════════════╝

✨ GORM Advantages:
   • Less boilerplate code (3-5x shorter)
   • Auto timestamps (CreatedAt, UpdatedAt)
   • Built-in soft deletes
   • Scopes for reusable queries
   • Automatic scanning (no manual Scan())

⚡ Raw SQL (pgx) Advantages:
   • 20-60% faster performance
   • Full control over queries
   • Better for complex queries
   • More transparent
```

---

### 2️⃣ Запустити Benchmarks

```bash
cd /Users/vkuzm/GolandProjects/sneakers_marketplace/examples/gorm_vs_raw

# Всі бенчмарки
go test -bench=. -benchmem

# Тільки CREATE операції
go test -bench=BenchmarkCreate -benchmem

# З більше ітерацій (точніше)
go test -bench=. -benchmem -benchtime=5s
```

**Очікуваний результат:**
```
goos: darwin
goarch: arm64
BenchmarkCreate_RawSQL-10       1000    1250000 ns/op    1024 B/op    15 allocs/op
BenchmarkCreate_GORM-10          800    1875000 ns/op    2048 B/op    28 allocs/op
BenchmarkGetByEmail_RawSQL-10   3000     450000 ns/op     512 B/op     8 allocs/op
BenchmarkGetByEmail_GORM-10     2500     650000 ns/op     896 B/op    14 allocs/op
BenchmarkUpdate_RawSQL-10       2000     700000 ns/op     768 B/op    12 allocs/op
BenchmarkUpdate_GORM-10         1500    1050000 ns/op    1280 B/op    22 allocs/op
```

**Інтерпретація:**
- GORM **~50% повільніше** за Raw SQL
- Але різниця в абсолютних числах: **< 1ms** (несуттєво для більшості операцій)

---

## 📊 Швидке порівняння

### Code Simplicity

**Raw SQL (23 рядки):**
```go
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
    query := `SELECT id, email, password_hash, first_name, last_name, phone,
              COALESCE(role, 'user') as role, is_verified, is_active, 
              created_at, updated_at FROM users WHERE email = $1`
    
    user := &model.User{}
    err := r.db.QueryRow(ctx, query, email).Scan(
        &user.ID, &user.Email, &user.PasswordHash,
        &user.FirstName, &user.LastName, &user.Phone, &user.Role,
        &user.IsVerified, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
    )
    return user, err
}
```

**GORM (6 рядків):**
```go
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model_gorm.User, error) {
    var user model_gorm.User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    return &user, err
}
```

**Різниця:** 4x менше коду! 🎉

---

## 🎯 Коли використовувати?

### ✅ Використовуй GORM для:
- Simple CRUD (Get, Create, Update, Delete)
- Admin панелі (низький трафік)
- Прототипування
- Коли потрібно швидко написати код
- Soft deletes

### ✅ Використовуй Raw SQL для:
- **Bidding Service** (matching engine - критична продуктивність)
- **Analytics** (складні запити з JOINs)
- **Order processing** (транзакції)
- Bulk operations (1000+ записів)
- Коли потрібен повний контроль

### 🎯 **Рекомендація: Hybrid Approach**
Використовуй обидва в одному проєкті!

---

## 🔄 Як інтегрувати GORM у проєкт?

### Варіант 1: Залишити як є ✅
Продовжити з Raw SQL (pgx) - він працює відмінно!

### Варіант 2: Додати GORM до Admin Service
```go
// cmd/admin-service/main.go

// Додати GORM connection
gormDB, err := gorm.Open(postgres.Open(cfg.Database.URL), &gorm.Config{})
if err != nil {
    log.Fatal().Err(err).Msg("Failed to connect to database (GORM)")
}

// Використовувати GORM repository
adminRepo := repository_gorm.NewUserRepository(gormDB)
```

### Варіант 3: Hybrid Repository
```go
type UserRepository struct {
    pgx  *pgxpool.Pool  // Для складних запитів
    gorm *gorm.DB       // Для простих CRUD
}

// Simple CRUD - use GORM
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
    var user User
    err := r.gorm.WithContext(ctx).First(&user, id).Error
    return &user, err
}

// Complex query - use pgx
func (r *UserRepository) GetUserStats(ctx context.Context, userID int64) (*Stats, error) {
    query := `SELECT ... complex JOIN query ...`
    return r.pgx.QueryRow(ctx, query, userID).Scan(...)
}
```

---

## 📚 Більше інформації

- **Детальний аналіз:** `docs/GORM_INVESTIGATION.md`
- **Приклади коду:** `examples/gorm_vs_raw/main.go`
- **README з benchmarks:** `examples/gorm_vs_raw/README.md`
- **GORM Docs:** https://gorm.io/docs/

---

## 🤔 FAQ

**Q: Чи вплине GORM на продуктивність?**  
A: GORM ~50% повільніше, але в абсолютних числах це < 1ms. Для більшості операцій це несуттєво.

**Q: Чи можна використовувати обидва в одному проєкті?**  
A: Так! Це **рекомендований підхід**. GORM для простих CRUD, Raw SQL для складних запитів.

**Q: Чи потрібно переписувати існуючий код?**  
A: Ні! Новий GORM пакет повністю окремий. Можна додавати поступово.

**Q: Чи підтримує GORM PostgreSQL?**  
A: Так, повністю. Використовує той же `pgx` драйвер всередині.

---

## 🎓 Наступні кроки

1. ✅ **Запусти demo:** `go run examples/gorm_vs_raw/main.go`
2. ✅ **Запусти benchmarks:** `go test -bench=. -benchmem`
3. ✅ **Прочитай код:** Порівняй `repository` vs `repository_gorm`
4. ✅ **Вивчи документацію:** `docs/GORM_INVESTIGATION.md`
5. 🤔 **Вирішуй:** Залишити Raw SQL, додати GORM, чи hybrid?

---

## 💡 Моя Рекомендація

**Для Sneakers Marketplace:**

1. **Залиши Raw SQL для критичних сервісів:**
   - ✅ Bidding Service (matching engine)
   - ✅ Order Service (transactions)
   - ✅ Payment Service

2. **Спробуй GORM для Admin Service:**
   - ✅ Низький трафік
   - ✅ Прості CRUD операції
   - ✅ Хороше місце для експериментів

3. **Вивчи обидва підходи:**
   - ✅ Raw SQL - глибоке розуміння SQL
   - ✅ GORM - швидкість розробки
   - ✅ Обидва - цінні навички!

**Поточний підхід (Raw SQL) вже ЧУДОВИЙ!** 🎉  
GORM - це додатковий інструмент, не заміна.

---

**Питання? Запускай demo і дивись результати!** 🚀

```bash
cd examples/gorm_vs_raw && go run main.go
```
