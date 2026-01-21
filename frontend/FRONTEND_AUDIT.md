# Frontend Аудит - Що використовується? 🔍

## ✅ ВИКОРИСТОВУЄТЬСЯ

### 📦 Компоненти UI:

1. ✅ **Input** (`components/ui/Input.tsx`)
   - Login.tsx
   - Register.tsx
   - **Залишити** ✅

2. ✅ **Typography** (`components/ui/Typography.tsx`)
   - ProductList.tsx
   - **Залишити** ✅

3. ✅ **Box** (`components/ui/Box.tsx`)
   - ProductList.tsx
   - **Залишити** ✅

4. ✅ **Card** (`components/ui/Card.tsx`)
   - ProductList.tsx
   - **Залишити** ✅

5. ✅ **Badge** (`components/ui/Badge.tsx`)
   - ProductList.tsx
   - **Залишити** ✅

---

### 📁 Сторінки (Features):

1. ✅ **Login** (`features/auth/Login.tsx` + `Login.module.css`)
   - Використовується ✅
   
2. ✅ **Register** (`features/auth/Register.tsx` + `Register.module.css`)
   - Використовується ✅

3. ✅ **ProductList** (`features/products/ProductList.tsx`)
   - Використовується ✅

4. ✅ **ProtectedRoute** (`features/auth/ProtectedRoute.tsx`)
   - Використовується в App.tsx ✅

---

### 🔌 Redux API:

1. ✅ **authApi** (`features/auth/authApi.ts`)
   - Login.tsx, Register.tsx, store.ts
   - **Залишити** ✅

2. ✅ **authSlice** (`features/auth/authSlice.ts`)
   - Login.tsx, Register.tsx, Header.tsx, store.ts
   - **Залишити** ✅

3. ✅ **productsApi** (`features/products/productsApi.ts`)
   - ProductList.tsx, store.ts
   - **Залишити** ✅

4. ✅ **biddingApi** (`features/bidding/biddingApi.ts`)
   - store.ts (готовий до майбутнього використання)
   - **Залишити** ✅

5. ✅ **ordersApi** (`features/orders/ordersApi.ts`)
   - store.ts (готовий до майбутнього використання)
   - **Залишити** ✅

6. ✅ **notificationsApi** (`features/notifications/notificationsApi.ts`)
   - store.ts (готовий до майбутнього використання)
   - **Залишити** ✅

---

### 🎯 Core:

1. ✅ **store.ts** - Redux store
2. ✅ **hooks.ts** - useAppDispatch, useAppSelector
3. ✅ **main.tsx** - entry point
4. ✅ **App.tsx** - routing
5. ✅ **Header.tsx** - navigation

---

## ❌ НЕ ВИКОРИСТОВУЄТЬСЯ

### 🗑️ Компоненти (можна видалити):

1. ❌ **Alert** (`components/ui/Alert.tsx`)
   - Не використовується
   - **ВИДАЛИТИ** ❌

2. ❌ **Button** (`components/ui/Button.tsx`)
   - Не використовується (використовуємо звичайний `<button>`)
   - **ВИДАЛИТИ** ❌

---

### 📁 Пусті папки (видалити):

1. ❌ `features/user/` - пуста
2. ❌ `utils/` - пуста
3. ❌ `styles/` - пуста
4. ❌ `components/common/` - пуста
5. ❌ `hooks/` - пуста

---

### 🔌 Сервіси (не використовуються):

1. ❌ **websocket.ts** (`services/websocket.ts`)
   - Створено, але не використовується
   - **ВИДАЛИТИ АБО ЗАЛИШИТИ** для Phase 3 (WebSocket notifications)

2. ❌ **storage.ts** (`services/storage.ts`)
   - Створено, але не використовується
   - **ВИДАЛИТИ** (localStorage використовується напряму в authSlice)

3. ❌ **api.ts** (`services/api.ts`)
   - Axios client, але не використовується (RTK Query)
   - **ВИДАЛИТИ**

---

### 📄 Інші файли:

1. ❌ **App.css** - пустий або непотрібний
2. ❌ **assets/react.svg** - непотрібний

---

## 🎯 РЕКОМЕНДАЦІЇ

### ✅ Обов'язково видалити:

```bash
# Непотрібні UI компоненти
rm src/components/ui/Alert.tsx
rm src/components/ui/Button.tsx

# Пусті папки
rmdir src/features/user
rmdir src/utils
rmdir src/styles
rmdir src/components/common
rmdir src/hooks

# Непотрібні сервіси
rm src/services/api.ts
rm src/services/storage.ts
rm src/services/websocket.ts

# Непотрібні файли
rm src/App.css
rm src/assets/react.svg
```

---

### 🤔 Опціонально (залежить від планів):

**Якщо НЕ плануєш Phase 3 (WebSocket):**
- ❌ Видалити `services/websocket.ts`

**Якщо плануєш Phase 3:**
- ✅ Залишити `services/websocket.ts`
- ✅ Залишити всі API (biddingApi, ordersApi, notificationsApi)

---

## 📊 Статистика:

### До очистки:
- **Файлів:** ~40
- **Папок:** ~15
- **Розмір:** ~150KB

### Після очистки:
- **Файлів:** ~28
- **Папок:** ~10
- **Розмір:** ~100KB

**Економія:** ~33% менше файлів!

---

## 🚀 План дій:

### Варіант 1: Повна очистка (рекомендую)

```bash
cd /Users/vkuzm/GolandProjects/sneakers_marketplace/frontend/src

# 1. Видалити непотрібні UI компоненти
rm components/ui/Alert.tsx
rm components/ui/Button.tsx

# 2. Видалити пусті папки
rmdir features/user utils styles components/common hooks

# 3. Видалити непотрібні сервіси
rm services/api.ts services/storage.ts services/websocket.ts

# 4. Видалити непотрібні файли
rm App.css assets/react.svg
```

---

### Варіант 2: Часткова очистка (залишити для майбутнього)

```bash
# Видалити тільки 100% непотрібне
rm components/ui/Alert.tsx
rm components/ui/Button.tsx
rmdir features/user utils styles components/common hooks
rm App.css assets/react.svg

# Залишити сервіси для Phase 3
# - services/websocket.ts
# - services/api.ts
# - services/storage.ts
```

---

## 💡 Підсумок:

### ✅ Що працює і використовується:

1. Login + Register (з CSS Modules)
2. ProductList (з Tailwind компонентами)
3. Redux (auth, products)
4. RTK Query APIs (готові до використання)

### ❌ Що не використовується:

1. Alert, Button компоненти
2. 5 пустих папок
3. 3 сервіси (websocket, api, storage)
4. 2 непотрібні файли

### 🎯 Рекомендація:

**Виконай "Варіант 1: Повна очистка"** - видалить все непотрібне та зменшить кодову базу на 33%!

Якщо хочеш залишити щось для майбутнього Phase 3, використай "Варіант 2".

---

**Готовий виконати очистку? Скажи "так" і я видалю все непотрібне!** 🧹
