# Redux - Як це працює? 🚀

## 🎯 Що таке Redux?

**Redux** — це бібліотека для управління станом (state) застосунку.

### Аналогія:

Уяви **банк** 🏦:
- **Store** (сховище) — це сейф у банку
- **State** (стан) — гроші в сейфі
- **Actions** (дії) — заявки "покласти" або "зняти" гроші
- **Reducers** (редюсери) — касири, які обробляють заявки
- **Dispatch** — відправка заявки касиру

---

## 📊 Як це працює?

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│  1. USER CLICKS BUTTON                                      │
│     ↓                                                       │
│  2. DISPATCH ACTION  → "LOGIN_SUCCESS"                      │
│     ↓                                                       │
│  3. REDUCER receives action                                 │
│     ↓                                                       │
│  4. REDUCER updates STATE                                   │
│     ↓                                                       │
│  5. REACT RE-RENDERS components using that state            │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔥 Приклад з нашого проекту

### 1️⃣ **STATE (стан)** — що зберігаємо?

```typescript
// src/features/auth/authSlice.ts
const initialState = {
  user: null,      // інформація про користувача
  token: null,     // JWT токен
};
```

**Це як змінна, але глобальна для всього застосунку!**

---

### 2️⃣ **ACTION (дія)** — що хочемо зробити?

```typescript
// Приклад: користувач залогінився
const action = {
  type: 'auth/setCredentials',
  payload: {
    user: { id: 1, email: 'test@example.com' },
    token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...'
  }
};
```

**Action = повідомлення "що сталося"**

---

### 3️⃣ **REDUCER (редюсер)** — як змінити state?

```typescript
// src/features/auth/authSlice.ts
reducers: {
  setCredentials: (state, action) => {
    // Отримали action → змінюємо state
    state.user = action.payload.user;
    state.token = action.payload.token;
    
    // Зберігаємо в localStorage
    localStorage.setItem('token', action.payload.token);
  },
  
  logout: (state) => {
    // Очищаємо state
    state.user = null;
    state.token = null;
    localStorage.removeItem('token');
  },
}
```

**Reducer = функція, яка каже "як змінити state"**

---

### 4️⃣ **DISPATCH (відправка)** — як запустити action?

```typescript
// src/features/auth/Login.tsx
import { useAppDispatch } from '../../app/hooks';
import { setCredentials } from './authSlice';

const dispatch = useAppDispatch();

// Коли користувач логінується:
const result = await login({ email, password }).unwrap();

// Відправляємо action
dispatch(setCredentials({
  user: result.user,
  token: result.access_token
}));
```

**Dispatch = "запусти цю дію"**

---

### 5️⃣ **SELECTOR (читання)** — як прочитати state?

```typescript
// src/components/layout/Header.tsx
import { useAppSelector } from '../../app/hooks';

const user = useAppSelector((state) => state.auth.user);
const token = useAppSelector((state) => state.auth.token);

// Тепер можна використати:
{user ? (
  <p>Hello, {user.first_name}!</p>
) : (
  <Link to="/login">Login</Link>
)}
```

**Selector = "дай мені частину state"**

---

## 🎯 Повний приклад: LOGIN FLOW

### Крок 1: Користувач натискає "Login"

```typescript
// Login.tsx
const handleSubmit = async (e: React.FormEvent) => {
  e.preventDefault();
  
  try {
    // 1️⃣ Викликаємо API
    const result = await login({ email, password }).unwrap();
    
    // 2️⃣ Відправляємо action
    dispatch(setCredentials(result));
    
    // 3️⃣ Переходимо на іншу сторінку
    navigate('/products');
  } catch (err) {
    console.error('Login failed:', err);
  }
};
```

---

### Крок 2: Reducer обробляє action

```typescript
// authSlice.ts
setCredentials: (state, action) => {
  state.user = action.payload.user;        // { id, email, first_name, ... }
  state.token = action.payload.access_token; // "eyJhbGciOiJIUzI1..."
  
  localStorage.setItem('token', action.payload.access_token);
}
```

---

### Крок 3: Header бачить зміни

```typescript
// Header.tsx
const user = useAppSelector((state) => state.auth.user);

// Автоматично ре-рендериться!
{user ? (
  <div>Welcome, {user.first_name}!</div>
) : (
  <Link to="/login">Login</Link>
)}
```

---

## 🔄 Redux Toolkit Query (RTK Query)

**RTK Query** — це розширення Redux для API запитів.

### Без RTK Query (старий спосіб):

```typescript
// ❌ Багато коду:
const [products, setProducts] = useState([]);
const [loading, setLoading] = useState(false);
const [error, setError] = useState(null);

useEffect(() => {
  setLoading(true);
  fetch('/api/products')
    .then(res => res.json())
    .then(data => setProducts(data))
    .catch(err => setError(err))
    .finally(() => setLoading(false));
}, []);
```

### З RTK Query (новий спосіб):

```typescript
// ✅ Одна строка:
const { data, isLoading, error } = useGetProductsQuery({ page: 1 });

// Все автоматично!
// - loading state
// - error handling
// - caching
// - re-fetching
```

---

## 📁 Структура Redux в нашому проекті

```
src/
├── app/
│   ├── store.ts           ← Redux Store (глобальний state)
│   └── hooks.ts           ← useAppDispatch, useAppSelector
│
├── features/
│   ├── auth/
│   │   ├── authSlice.ts   ← state для auth (user, token)
│   │   └── authApi.ts     ← RTK Query для login/register
│   │
│   ├── products/
│   │   └── productsApi.ts ← RTK Query для products
│   │
│   ├── bidding/
│   │   └── biddingApi.ts  ← RTK Query для bid/ask
│   │
│   └── ...
```

---

## 🎯 Коли використовувати Redux?

### ✅ ВИКОРИСТОВУЙ REDUX:

1. **Глобальний state** — user, auth token, theme
2. **Спільні дані** — дані, які потрібні багатьом компонентам
3. **Складна логіка** — багато взаємозв'язаних станів

### ❌ НЕ ВИКОРИСТОВУЙ REDUX:

1. **Локальний state** — стан форми, модалки (використай `useState`)
2. **Простий застосунок** — 1-2 сторінки без спільного state
3. **Server state** — дані з API (використай RTK Query)

---

## 🔍 Redux DevTools

**Відкрий Chrome DevTools → Redux Tab**

Там побачиш:
- 📜 Всі actions (що відбулося)
- 🔍 State до і після кожної action
- ⏮️ Time travel (повернутися назад)
- 🐛 Debugging

---

## 💡 Ключові концепції

### 1. **Immutability** (незмінність)

```typescript
// ❌ НЕ РОБИ ТАК:
state.user = { ...state.user, name: 'New Name' };

// ✅ РОБИ ТАК (Redux Toolkit робить це автоматично):
state.user.name = 'New Name';
```

**Redux Toolkit використовує Immer, тому можна писати "мутуючий" код!**

---

### 2. **Single Source of Truth**

```typescript
// ❌ БЕЗ REDUX:
// Header.tsx
const [user, setUser] = useState(null);

// Sidebar.tsx
const [user, setUser] = useState(null);

// Profile.tsx
const [user, setUser] = useState(null);

// ❗ Проблема: 3 різні копії user!


// ✅ З REDUX:
// Один store → один user → всі компоненти використовують його
const user = useAppSelector((state) => state.auth.user);
```

---

### 3. **Unidirectional Data Flow**

```
┌─────────────────────────────────────────────────┐
│                                                 │
│  VIEW → ACTION → REDUCER → STATE → VIEW        │
│   ↑                                        ↓    │
│   └────────────────────────────────────────┘    │
│                                                 │
└─────────────────────────────────────────────────┘
```

**Дані течуть тільки в одному напрямку!**

---

## 🎓 Приклади з коду

### Приклад 1: Читання state

```typescript
// Будь-який компонент
import { useAppSelector } from '../../app/hooks';

function MyComponent() {
  const user = useAppSelector((state) => state.auth.user);
  const token = useAppSelector((state) => state.auth.token);
  
  return <div>Hello, {user?.first_name}</div>;
}
```

---

### Приклад 2: Зміна state

```typescript
import { useAppDispatch } from '../../app/hooks';
import { setCredentials, logout } from './authSlice';

function MyComponent() {
  const dispatch = useAppDispatch();
  
  const handleLogin = () => {
    dispatch(setCredentials({ user, token }));
  };
  
  const handleLogout = () => {
    dispatch(logout());
  };
}
```

---

### Приклад 3: API запит (RTK Query)

```typescript
import { useGetProductsQuery } from './productsApi';

function ProductList() {
  const { data, isLoading, error } = useGetProductsQuery({ 
    page: 1, 
    pageSize: 12 
  });
  
  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error!</div>;
  
  return (
    <div>
      {data?.products.map(product => (
        <div key={product.id}>{product.name}</div>
      ))}
    </div>
  );
}
```

---

## 📚 Корисні ресурси

- [Redux Toolkit Docs](https://redux-toolkit.js.org/)
- [RTK Query Tutorial](https://redux-toolkit.js.org/tutorials/rtk-query)
- [Redux DevTools](https://github.com/reduxjs/redux-devtools)

---

## ❓ Часті питання

### 1. Redux vs Context API?

**Context API:**
- ✅ Простіше для невеликих застосунків
- ❌ Re-renders всі компоненти при зміні

**Redux:**
- ✅ Оптимізовано (re-renders тільки потрібні компоненти)
- ✅ Redux DevTools для debugging
- ✅ Middleware (logging, async)

---

### 2. Коли використовувати Redux Toolkit Query?

**Використовуй RTK Query для:**
- ✅ CRUD операцій (Create, Read, Update, Delete)
- ✅ Автоматичний caching
- ✅ Автоматичний re-fetching

**Використовуй звичайний Redux для:**
- ✅ Локальний state (theme, sidebar open/close)
- ✅ Складна бізнес-логіка

---

### 3. Чи потрібен Redux для малих проектів?

**НІ!** Для малих проектів використовуй:
- `useState` — локальний state
- `useContext` — глобальний state
- `React Query` — server state

**Redux потрібен тільки для складних застосунків!**

---

## 🎯 Підсумок

```typescript
// 1. Створюємо store (один раз)
const store = configureStore({ reducer: { auth: authReducer } });

// 2. Підключаємо до React (один раз)
<Provider store={store}>
  <App />
</Provider>

// 3. Читаємо state (скрізь де потрібно)
const user = useAppSelector((state) => state.auth.user);

// 4. Змінюємо state (через actions)
dispatch(setCredentials({ user, token }));

// 5. RTK Query для API (автоматично!)
const { data } = useGetProductsQuery({ page: 1 });
```

---

**Redux = єдине джерело правди для твого застосунку! 🎉**
