import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import type { RootState } from '../../app/store';
import type {
  SubscriptionPlan,
  UserSubscriptionWithPlan,
  SubscriptionTransaction,
  SubscribeRequest,
  SubscribeResponse,
  CancelSubscriptionRequest,
  UpdateSubscriptionRequest,
  FeeSavingsCalculation,
} from '../../types/subscription.types';

export const subscriptionApi = createApi({
  reducerPath: 'subscriptionApi',
  baseQuery: fetchBaseQuery({
    baseUrl: '/api/v1/subscriptions',
    prepareHeaders: (headers, { getState }) => {
      const token = (getState() as RootState).auth.accessToken;
      if (token) {
        headers.set('Authorization', `Bearer ${token}`);
      }
      return headers;
    },
  }),
  tagTypes: ['SubscriptionPlans', 'UserSubscription', 'Transactions'],
  endpoints: (builder) => ({
    // Get all subscription plans
    getSubscriptionPlans: builder.query<SubscriptionPlan[], void>({
      query: () => '/plans',
      providesTags: ['SubscriptionPlans'],
    }),

    // Get user's current subscription
    getCurrentSubscription: builder.query<UserSubscriptionWithPlan, void>({
      query: () => '/current',
      providesTags: ['UserSubscription'],
    }),

    // Subscribe to a plan
    subscribe: builder.mutation<SubscribeResponse, SubscribeRequest>({
      query: (data) => ({
        url: '/subscribe',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: ['UserSubscription', 'Transactions'],
    }),

    // Cancel subscription
    cancelSubscription: builder.mutation<{ message: string }, CancelSubscriptionRequest>({
      query: (data) => ({
        url: '/cancel',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: ['UserSubscription'],
    }),

    // Update subscription (upgrade/downgrade)
    updateSubscription: builder.mutation<UserSubscriptionWithPlan, UpdateSubscriptionRequest>({
      query: (data) => ({
        url: '/update',
        method: 'PUT',
        body: data,
      }),
      invalidatesTags: ['UserSubscription', 'Transactions'],
    }),

    // Get subscription transactions history
    getTransactions: builder.query<SubscriptionTransaction[], void>({
      query: () => '/transactions',
      providesTags: ['Transactions'],
    }),

    // Calculate fee savings for a plan
    calculateSavings: builder.query<FeeSavingsCalculation, { planId: number; salePrice: number }>({
      query: ({ planId, salePrice }) => `/savings?plan_id=${planId}&sale_price=${salePrice}`,
    }),
  }),
});

export const {
  useGetSubscriptionPlansQuery,
  useGetCurrentSubscriptionQuery,
  useSubscribeMutation,
  useCancelSubscriptionMutation,
  useUpdateSubscriptionMutation,
  useGetTransactionsQuery,
  useCalculateSavingsQuery,
} = subscriptionApi;
