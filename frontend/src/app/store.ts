import { configureStore } from '@reduxjs/toolkit';
import { setupListeners } from '@reduxjs/toolkit/query';
import authReducer from '../features/auth/authSlice';
import { authApi } from '../features/auth/authApi';
import { productsApi } from '../features/products/productsApi';
import { biddingApi } from '../features/bidding/biddingApi';
import { ordersApi } from '../features/orders/ordersApi';
import { notificationsApi } from '../features/notifications/notificationsApi';
import { subscriptionApi } from '../features/subscription/subscriptionApi';

export const store = configureStore({
  reducer: {
    auth: authReducer,
    [authApi.reducerPath]: authApi.reducer,
    [productsApi.reducerPath]: productsApi.reducer,
    [biddingApi.reducerPath]: biddingApi.reducer,
    [ordersApi.reducerPath]: ordersApi.reducer,
    [notificationsApi.reducerPath]: notificationsApi.reducer,
    [subscriptionApi.reducerPath]: subscriptionApi.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(
      authApi.middleware,
      productsApi.middleware,
      biddingApi.middleware,
      ordersApi.middleware,
      notificationsApi.middleware,
      subscriptionApi.middleware
    ),
});

setupListeners(store.dispatch);

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
