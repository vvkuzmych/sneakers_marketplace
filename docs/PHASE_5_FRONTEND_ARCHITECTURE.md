# 🎨 Phase 5: Frontend Architecture

**React + Redux + TypeScript для Sneakers Marketplace**

---

## 🎯 Цілі Phase 5

Створити сучасний, швидкий та зручний веб-інтерфейс для Sneakers Marketplace з використанням:
- **React 18** - Component-based UI
- **Redux Toolkit** - State management
- **RTK Query** - API communication
- **TypeScript** - Type safety
- **React Router v6** - Navigation
- **WebSocket** - Real-time notifications
- **Tailwind CSS** - Modern styling
- **Vite** - Fast build tool

---

## 📐 Архітектура Frontend

```
┌──────────────────────────────────────────────────────────────┐
│                     React Application                        │
│                        (Vite + TS)                           │
└───────────────────────┬──────────────────────────────────────┘
                        │
        ┌───────────────┼───────────────┐
        │               │               │
┌───────▼──────┐ ┌─────▼─────┐ ┌──────▼───────┐
│  Redux Store │ │  RTK Query│ │  WebSocket   │
│  (Toolkit)   │ │  (API)    │ │  (Real-time) │
└──────┬───────┘ └─────┬─────┘ └──────┬───────┘
       │               │               │
       └───────────────┼───────────────┘
                       │
                ┌──────▼──────┐
                │ API Gateway │
                │   :8080     │
                └─────────────┘
```

---

## 📁 Структура Проекту

```
frontend/
├── public/                      # Статичні файли
│   ├── favicon.ico
│   └── logo.png
├── src/
│   ├── app/                     # Redux store setup
│   │   ├── store.ts             # Redux store configuration
│   │   └── hooks.ts             # Typed hooks (useAppDispatch, useAppSelector)
│   │
│   ├── features/                # Feature-based modules
│   │   ├── auth/                # 🔐 Аутентифікація
│   │   │   ├── authSlice.ts     # Redux slice
│   │   │   ├── authApi.ts       # RTK Query API
│   │   │   ├── Login.tsx        # Login page
│   │   │   ├── Register.tsx     # Register page
│   │   │   └── ProtectedRoute.tsx
│   │   │
│   │   ├── products/            # 📦 Каталог продуктів
│   │   │   ├── productsSlice.ts
│   │   │   ├── productsApi.ts
│   │   │   ├── ProductList.tsx  # Список продуктів
│   │   │   ├── ProductCard.tsx  # Картка продукту
│   │   │   ├── ProductDetail.tsx# Деталі продукту
│   │   │   └── ProductSearch.tsx
│   │   │
│   │   ├── bidding/             # 🎯 Біддінг (Bid/Ask)
│   │   │   ├── biddingSlice.ts
│   │   │   ├── biddingApi.ts
│   │   │   ├── BidAskBoard.tsx  # Order book
│   │   │   ├── PlaceBid.tsx     # Форма для bid
│   │   │   ├── PlaceAsk.tsx     # Форма для ask
│   │   │   └── MarketPrice.tsx  # Ринкова ціна
│   │   │
│   │   ├── orders/              # 📦 Замовлення
│   │   │   ├── ordersSlice.ts
│   │   │   ├── ordersApi.ts
│   │   │   ├── OrderList.tsx    # Список замовлень
│   │   │   ├── OrderDetail.tsx  # Деталі замовлення
│   │   │   └── OrderStatus.tsx  # Статус компонент
│   │   │
│   │   ├── notifications/       # 🔔 Сповіщення
│   │   │   ├── notificationsSlice.ts
│   │   │   ├── notificationsApi.ts
│   │   │   ├── NotificationBell.tsx
│   │   │   ├── NotificationList.tsx
│   │   │   └── NotificationItem.tsx
│   │   │
│   │   └── user/                # 👤 Профіль користувача
│   │       ├── userSlice.ts
│   │       ├── userApi.ts
│   │       ├── Profile.tsx
│   │       ├── AddressList.tsx
│   │       └── Settings.tsx
│   │
│   ├── components/              # Shared components
│   │   ├── layout/
│   │   │   ├── Header.tsx       # Навігація + JWT
│   │   │   ├── Footer.tsx
│   │   │   └── Sidebar.tsx
│   │   ├── ui/                  # UI components
│   │   │   ├── Button.tsx
│   │   │   ├── Input.tsx
│   │   │   ├── Modal.tsx
│   │   │   ├── Card.tsx
│   │   │   ├── Spinner.tsx
│   │   │   └── Toast.tsx
│   │   └── common/
│   │       ├── ErrorBoundary.tsx
│   │       ├── LoadingSpinner.tsx
│   │       └── EmptyState.tsx
│   │
│   ├── services/                # Services
│   │   ├── api.ts               # Axios/Fetch setup
│   │   ├── websocket.ts         # WebSocket client
│   │   └── storage.ts           # LocalStorage utils
│   │
│   ├── hooks/                   # Custom hooks
│   │   ├── useAuth.ts           # Auth logic
│   │   ├── useWebSocket.ts      # WebSocket hook
│   │   ├── useLocalStorage.ts
│   │   └── useDebounce.ts
│   │
│   ├── utils/                   # Utility functions
│   │   ├── formatters.ts        # Date, price formatters
│   │   ├── validators.ts        # Form validation
│   │   └── constants.ts         # App constants
│   │
│   ├── types/                   # TypeScript types
│   │   ├── auth.types.ts
│   │   ├── product.types.ts
│   │   ├── bidding.types.ts
│   │   ├── order.types.ts
│   │   └── api.types.ts
│   │
│   ├── styles/                  # Global styles
│   │   ├── index.css            # Tailwind imports
│   │   └── variables.css        # CSS variables
│   │
│   ├── App.tsx                  # Main App component
│   ├── main.tsx                 # Entry point
│   └── vite-env.d.ts            # Vite types
│
├── index.html                   # HTML template
├── package.json                 # Dependencies
├── tsconfig.json                # TypeScript config
├── vite.config.ts               # Vite config
├── tailwind.config.js           # Tailwind config
├── postcss.config.js            # PostCSS config
└── .env.example                 # Environment variables

```

---

## 🔧 Tech Stack

### Core
- **React 18** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool (faster than CRA)

### State Management
- **Redux Toolkit** - State management
- **RTK Query** - Data fetching & caching
- **Redux Persist** - Persist auth state

### Routing
- **React Router v6** - Client-side routing

### Styling
- **Tailwind CSS** - Utility-first CSS
- **HeadlessUI** - Unstyled UI components
- **Heroicons** - SVG icons

### Real-time
- **WebSocket API** - Real-time notifications
- **React hook** - Custom useWebSocket

### Forms & Validation
- **React Hook Form** - Form management
- **Zod** - Schema validation

### Utils
- **Axios** - HTTP client
- **date-fns** - Date formatting
- **classnames** - Conditional classes

---

## 🎨 Pages & Routes

### Public Routes
```tsx
/                       → HomePage (Landing)
/products               → ProductList (Каталог)
/products/:id           → ProductDetail (Деталі)
/login                  → Login
/register               → Register
```

### Protected Routes (Authentication Required)
```tsx
/dashboard              → Dashboard (User overview)
/bidding/:productId     → BidAskBoard (Order book)
/orders                 → OrderList (Мої замовлення)
/orders/:id             → OrderDetail (Деталі замовлення)
/notifications          → NotificationList
/profile                → Profile (Settings)
/profile/addresses      → AddressList
```

---

## 🔐 Authentication Flow

### Login Process
```
1. User enters email + password
2. POST /api/v1/auth/login → API Gateway
3. Receive: { access_token, refresh_token, user }
4. Store tokens in Redux + LocalStorage
5. Redirect to /dashboard
6. WebSocket connects with JWT
```

### Token Management
```tsx
// Redux state
{
  auth: {
    user: { id, email, firstName, lastName } | null,
    accessToken: string | null,
    refreshToken: string | null,
    isAuthenticated: boolean,
    isLoading: boolean,
    error: string | null
  }
}
```

### Axios Interceptor
```ts
// Automatically attach JWT to all requests
axios.interceptors.request.use(config => {
  const token = store.getState().auth.accessToken;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Refresh token on 401
axios.interceptors.response.use(
  response => response,
  async error => {
    if (error.response?.status === 401) {
      // Try refresh token
      // If fails → logout
    }
    return Promise.reject(error);
  }
);
```

---

## 📡 RTK Query API Structure

### authApi.ts
```ts
export const authApi = createApi({
  reducerPath: 'authApi',
  baseQuery: fetchBaseQuery({ baseUrl: '/api/v1' }),
  endpoints: (builder) => ({
    login: builder.mutation<LoginResponse, LoginRequest>({...}),
    register: builder.mutation<RegisterResponse, RegisterRequest>({...}),
    logout: builder.mutation<void, void>({...}),
    getProfile: builder.query<User, string>({...})
  })
});
```

### productsApi.ts
```ts
export const productsApi = createApi({
  endpoints: (builder) => ({
    getProducts: builder.query<ProductsResponse, ProductsRequest>({...}),
    getProduct: builder.query<Product, string>({...}),
    searchProducts: builder.query<ProductsResponse, string>({...})
  })
});
```

### biddingApi.ts
```ts
export const biddingApi = createApi({
  endpoints: (builder) => ({
    placeBid: builder.mutation<BidResponse, PlaceBidRequest>({...}),
    placeAsk: builder.mutation<AskResponse, PlaceAskRequest>({...}),
    getMarketPrice: builder.query<MarketPrice, { productId, sizeId }>({...}),
    getBids: builder.query<BidsResponse, { productId, sizeId }>({...}),
    getAsks: builder.query<AsksResponse, { productId, sizeId }>({...})
  })
});
```

---

## 🔌 WebSocket Integration

### useWebSocket Hook
```tsx
function useWebSocket() {
  const { accessToken } = useAppSelector(state => state.auth);
  const dispatch = useAppDispatch();
  const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    if (!accessToken) return;

    const socket = new WebSocket(
      `ws://localhost:8080/ws?token=${accessToken}`
    );

    socket.onopen = () => {
      console.log('WebSocket connected');
    };

    socket.onmessage = (event) => {
      const message = JSON.parse(event.data);
      
      if (message.type === 'notification') {
        // Dispatch to Redux
        dispatch(addNotification(message.data));
        // Show toast
        toast.success(message.data.title);
      }
    };

    socket.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    socket.onclose = () => {
      console.log('WebSocket disconnected');
      // Auto-reconnect logic
    };

    setWs(socket);

    return () => {
      socket.close();
    };
  }, [accessToken]);

  return { ws, isConnected: ws?.readyState === WebSocket.OPEN };
}
```

---

## 🎨 UI Components (Tailwind CSS)

### ProductCard Example
```tsx
<div className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-xl transition-shadow">
  <img
    src={product.imageUrl}
    alt={product.name}
    className="w-full h-48 object-cover"
  />
  <div className="p-4">
    <h3 className="text-lg font-semibold text-gray-900">
      {product.name}
    </h3>
    <p className="text-sm text-gray-500 mt-1">
      {product.brand} - {product.model}
    </p>
    <div className="mt-4 flex items-center justify-between">
      <span className="text-xl font-bold text-green-600">
        ${product.retailPrice}
      </span>
      <button className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">
        View
      </button>
    </div>
  </div>
</div>
```

### BidAskBoard Example
```tsx
<div className="grid grid-cols-2 gap-4">
  {/* Bids (Buyers) */}
  <div className="bg-green-50 p-4 rounded-lg">
    <h3 className="text-lg font-semibold mb-4">Bids (Buy)</h3>
    {bids.map(bid => (
      <div key={bid.id} className="flex justify-between py-2 border-b">
        <span className="text-green-700 font-medium">${bid.price}</span>
        <span className="text-gray-600">Qty: {bid.quantity}</span>
      </div>
    ))}
  </div>

  {/* Asks (Sellers) */}
  <div className="bg-red-50 p-4 rounded-lg">
    <h3 className="text-lg font-semibold mb-4">Asks (Sell)</h3>
    {asks.map(ask => (
      <div key={ask.id} className="flex justify-between py-2 border-b">
        <span className="text-red-700 font-medium">${ask.price}</span>
        <span className="text-gray-600">Qty: {ask.quantity}</span>
      </div>
    ))}
  </div>
</div>
```

---

## 🔄 State Management Examples

### Redux Slice (authSlice.ts)
```ts
const authSlice = createSlice({
  name: 'auth',
  initialState: {
    user: null,
    accessToken: null,
    refreshToken: null,
    isAuthenticated: false,
    isLoading: false,
    error: null
  },
  reducers: {
    setCredentials: (state, action) => {
      state.user = action.payload.user;
      state.accessToken = action.payload.accessToken;
      state.refreshToken = action.payload.refreshToken;
      state.isAuthenticated = true;
    },
    logout: (state) => {
      state.user = null;
      state.accessToken = null;
      state.refreshToken = null;
      state.isAuthenticated = false;
    }
  }
});
```

---

## 📊 Feature Priority

### Phase 5.1 - Core (Week 1)
- ✅ Project setup (Vite + React + TS + Tailwind)
- ✅ Redux store setup (Toolkit + RTK Query)
- ✅ Auth pages (Login, Register)
- ✅ Header with navigation
- ✅ Protected routes
- ✅ API integration (auth endpoints)

### Phase 5.2 - Products (Week 2)
- ✅ Product list page with pagination
- ✅ Product detail page
- ✅ Product search
- ✅ Image gallery
- ✅ Size selector

### Phase 5.3 - Bidding (Week 3)
- ✅ Bid/Ask board (order book)
- ✅ Place Bid form
- ✅ Place Ask form
- ✅ Market price display
- ✅ Real-time updates (WebSocket)

### Phase 5.4 - Orders & Profile (Week 4)
- ✅ Order list (buyer + seller views)
- ✅ Order detail page
- ✅ Order status tracking
- ✅ User profile page
- ✅ Address management

### Phase 5.5 - Notifications & Polish (Week 5)
- ✅ WebSocket integration
- ✅ Notification bell
- ✅ Notification list
- ✅ Toast notifications
- ✅ Error handling
- ✅ Loading states
- ✅ Responsive design

---

## 🧪 Testing Strategy

### Unit Tests (Vitest)
```bash
npm run test
```
- Redux slices
- Utility functions
- Custom hooks

### Component Tests (React Testing Library)
- User interactions
- Form validation
- API mocking

### E2E Tests (Playwright) - Optional
- User flows
- Critical paths

---

## 🚀 Development Workflow

### 1. Setup Project
```bash
cd /Users/vkuzm/GolandProjects/sneakers_marketplace
npm create vite@latest frontend -- --template react-ts
cd frontend
npm install
```

### 2. Install Dependencies
```bash
npm install @reduxjs/toolkit react-redux react-router-dom
npm install axios react-hook-form zod
npm install -D tailwindcss postcss autoprefixer
npm install @headlessui/react @heroicons/react
npm install date-fns classnames
```

### 3. Run Dev Server
```bash
npm run dev
# Open http://localhost:5173
```

### 4. Build for Production
```bash
npm run build
npm run preview
```

---

## 🔗 API Integration

### Environment Variables (.env)
```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_WS_URL=ws://localhost:8080/ws
```

### Axios Setup
```ts
import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
});

export default api;
```

---

## 📚 Resources

**React:**
- https://react.dev/
- https://react-typescript-cheatsheet.netlify.app/

**Redux Toolkit:**
- https://redux-toolkit.js.org/
- https://redux-toolkit.js.org/rtk-query/overview

**Tailwind CSS:**
- https://tailwindcss.com/docs
- https://tailwindui.com/components

**React Router:**
- https://reactrouter.com/

---

## ✅ Success Criteria

- [ ] User can register and login
- [ ] JWT tokens are stored and used
- [ ] Product list loads and displays
- [ ] Product detail shows images and sizes
- [ ] User can place bids and asks
- [ ] Order book updates in real-time
- [ ] Orders are visible in user dashboard
- [ ] Notifications work via WebSocket
- [ ] Responsive on mobile/tablet/desktop
- [ ] Error handling with user feedback
- [ ] Loading states for all async operations

---

**Created:** 2026-01-21  
**Phase:** 5 - Frontend  
**Status:** 🚧 Planning
