import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import type {
  PlaceBidRequest,
  PlaceAskRequest,
  BidResponse,
  AskResponse,
  MarketPrice,
  BidsResponse,
  AsksResponse,
} from '../../types/bidding.types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

export const biddingApi = createApi({
  reducerPath: 'biddingApi',
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
  tagTypes: ['Bid', 'Ask', 'MarketPrice'],
  endpoints: (builder) => ({
    placeBid: builder.mutation<BidResponse, PlaceBidRequest>({
      query: (bid) => ({
        url: '/bids',
        method: 'POST',
        body: bid,
      }),
      invalidatesTags: ['Bid', 'MarketPrice'],
    }),
    placeAsk: builder.mutation<AskResponse, PlaceAskRequest>({
      query: (ask) => ({
        url: '/asks',
        method: 'POST',
        body: ask,
      }),
      invalidatesTags: ['Ask', 'MarketPrice'],
    }),
    getMarketPrice: builder.query<MarketPrice, { productId: string; sizeId: string }>({
      query: ({ productId, sizeId }) => `/market/${productId}/${sizeId}`,
      providesTags: ['MarketPrice'],
    }),
    getBids: builder.query<BidsResponse, { productId: string; sizeId: string }>({
      query: ({ productId, sizeId }) => `/bids/product/${productId}?size_id=${sizeId}`,
      providesTags: ['Bid'],
    }),
    getAsks: builder.query<AsksResponse, { productId: string; sizeId: string }>({
      query: ({ productId, sizeId }) => `/asks/product/${productId}?size_id=${sizeId}`,
      providesTags: ['Ask'],
    }),
    cancelBid: builder.mutation<void, string>({
      query: (bidId) => ({
        url: `/bids/${bidId}/cancel`,
        method: 'POST',
      }),
      invalidatesTags: ['Bid', 'MarketPrice'],
    }),
    cancelAsk: builder.mutation<void, string>({
      query: (askId) => ({
        url: `/asks/${askId}/cancel`,
        method: 'POST',
      }),
      invalidatesTags: ['Ask', 'MarketPrice'],
    }),
  }),
});

export const {
  usePlaceBidMutation,
  usePlaceAskMutation,
  useGetMarketPriceQuery,
  useGetBidsQuery,
  useGetAsksQuery,
  useCancelBidMutation,
  useCancelAskMutation,
} = biddingApi;
