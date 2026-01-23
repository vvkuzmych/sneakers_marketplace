# 📧 Production Email Configuration

## ✅ ТАК! В PRODUCTION EMAIL ВІДПРАВЛЯЮТЬСЯ КОРИСТУВАЧАМ!

---

## 🎯 Як це працює в Production:

1. **User1** створює BID (email: `user1@gmail.com`)
2. **User2** створює ASK (email: `user2@yahoo.com`)
3. **MATCH** створено! 🎯
4. **Notification Service** відправляє:
   - Email на `user1@gmail.com`: "Your BID matched!"
   - Email на `user2@yahoo.com`: "Your ASK matched!"
5. ✅ **Користувачі ОТРИМУЮТЬ email у своїх поштових скриньках!**

---

## 🔧 Що потрібно змінити для Production

### 1️⃣ Оновити код `email_service.go`

**Файл:** `internal/notification/email/email_service.go`

```go
func NewEmailService() *EmailService {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		host = "localhost" // DEV: Mailhog
	}

	port := os.Getenv("SMTP_PORT")
	if port == "" {
		port = "1025" // DEV: Mailhog
	}

	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = "noreply@sneakersmarketplace.com"
	}

	// 🔴 ЗМІНИТИ ЦЕ ДЛЯ PRODUCTION:
	var auth smtp.Auth
	
	// Для production додати SMTP аутентифікацію
	username := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")
	
	if username != "" && password != "" {
		// Production SMTP with authentication
		auth = smtp.PlainAuth("", username, password, host)
	} else {
		// DEV: Mailhog без authentication
		auth = nil
	}

	return &EmailService{
		host: host,
		port: port,
		from: from,
		auth: auth, // ← Тепер буде auth для production
	}
}
```

---

### 2️⃣ Налаштувати `.env` для Production

**DEV (зараз):**
```bash
# .env (Development)
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_FROM=noreply@sneakersmarketplace.com
# SMTP_USER та SMTP_PASS не потрібні для Mailhog
```

**PRODUCTION (Gmail):**
```bash
# .env (Production)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_FROM=noreply@sneakersmarketplace.com
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password  # НЕ звичайний пароль!
```

**PRODUCTION (SendGrid):**
```bash
# .env (Production - SendGrid)
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_FROM=noreply@sneakersmarketplace.com
SMTP_USER=apikey
SMTP_PASS=your-sendgrid-api-key
```

**PRODUCTION (AWS SES):**
```bash
# .env (Production - AWS SES)
SMTP_HOST=email-smtp.us-east-1.amazonaws.com
SMTP_PORT=587
SMTP_FROM=noreply@sneakersmarketplace.com
SMTP_USER=your-aws-smtp-username
SMTP_PASS=your-aws-smtp-password
```

---

### 3️⃣ Як отримати App Password для Gmail

1. Перейти: https://myaccount.google.com/security
2. Увімкнути **2-Step Verification**
3. Перейти до **App passwords**
4. Створити новий App Password для "Mail"
5. Скопіювати пароль у `.env` → `SMTP_PASS`

---

## 🚀 Як запустити в Production

### Крок 1: Оновити код
```bash
cd /Users/vkuzm/GolandProjects/sneakers_marketplace
# Змінити код в email_service.go (як показано вище)
```

### Крок 2: Оновити .env
```bash
# Відредагувати .env файл з production SMTP налаштуваннями
nano .env
```

### Крок 3: Перезбудувати Notification Service
```bash
make build-notification  # або go build
```

### Крок 4: Перезапустити Notification Service
```bash
source .env
./bin/notification-service
```

### Крок 5: Зупинити Mailhog (більше не потрібен)
```bash
pkill -f MailHog
```

---

## ✅ Що відбувається в Production

```
MATCH створено
   ↓
Notification Service готує email
   ↓
Відправляє через smtp.gmail.com:587 (Production SMTP)
   ↓
Gmail відправляє email користувачам
   ↓
✅ Користувачі ОТРИМУЮТЬ email у своїй поштовій скриньці
```

---

## 📊 Порівняння DEV vs PRODUCTION

| Параметр | DEV (зараз) | PRODUCTION |
|----------|-------------|------------|
| SMTP Host | `localhost` | `smtp.gmail.com` |
| SMTP Port | `1025` | `587` |
| Authentication | ❌ Не потрібна | ✅ Потрібна (username+password) |
| Email доставка | ❌ Mailhog (локально) | ✅ Користувачам |
| Перегляд | `http://localhost:8025` | Користувачі у своїй пошті |
| Вартість | 🆓 Безкоштовно | Gmail 🆓 / SendGrid/SES 💰 |

---

## 🎯 Висновок

**ТАК!** В production email **БУДУТЬ відправлятися користувачам**! 🎉

Потрібно тільки:
1. ✅ Додати SMTP authentication в код (показано вище)
2. ✅ Налаштувати `.env` з production SMTP
3. ✅ Перезбудувати та перезапустити Notification Service

**Код вже готовий** — він використовує справжній `smtp.SendMail()`, який працює з будь-яким SMTP сервером! 🚀
