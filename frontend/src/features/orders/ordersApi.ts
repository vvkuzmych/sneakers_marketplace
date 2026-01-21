import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import type { Order, OrdersResponse, OrdersRequest } from '../../types/order.types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

export const ordersApi = createApi({
  reducerPath: 'ordersApi',
  baseQuery: fetchBaseQuery({
    baseUrl: API_BASE_URL,
    prepareHeaders: (headers, { getState }) => {
      const token = (getState() as any).auth.accessToken;
      if (token) {
        headers.set('Authorization', `Bearer ${token}`);
      }
      return headers;
    },
  }),
  tagTypes: ['Order'],
  endpoints: (builder) => ({
    getOrder: builder.query<Order, string>({
      query: (orderId) => `/orders/${orderId}`,
      providesTags: (_result, _error, id) => [{ type: 'Order', id }],
    }),
    getBuyerOrders: builder.query<OrdersResponse, { buyerId: string } & OrdersRequest>({
      query: ({ buyerId, page, pageSize, status }) => {
        const searchParams = new URLSearchParams();
        if (page) searchParams.append('page', page.toString());
        if (pageSize) searchParams.append('page_size', pageSize.toString());
        if (status) searchParams.append('status', status);
        
        return `/orders/buyer/${buyerId}?${searchParams.toString()}`;
      },
      providesTags: [{ type: 'Order', id: 'BUYER_LIST' }],
    }),
    getSellerOrders: builder.query<OrdersResponse, { sellerId: string } & OrdersRequest>({
      query: ({ sellerId, page, pageSize, status }) => {
        const searchParams = new URLSearchParams();
        if (page) searchParams.append('page', page.toString());
        if (pageSize) searchParams.append('page_size', pageSize.toString());
        if (status) searchParams.append('status', status);
        
        return `/orders/seller/${sellerId}?${searchParams.toString()}`;
      },
      providesTags: [{ type: 'Order', id: 'SELLER_LIST' }],
    }),
  }),
});

export const {
  useGetOrderQuery,
  useGetBuyerOrdersQuery,
  useGetSellerOrdersQuery,
} = ordersApi;
