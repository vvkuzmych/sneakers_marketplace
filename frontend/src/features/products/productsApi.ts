import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import type {
  Product,
  ProductsResponse,
  ProductsRequest,
  SearchProductsRequest,
} from '../../types/product.types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

export const productsApi = createApi({
  reducerPath: 'productsApi',
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
  tagTypes: ['Product'],
  endpoints: (builder) => ({
    getProducts: builder.query<ProductsResponse, ProductsRequest | void>({
      query: (params) => {
        const searchParams = new URLSearchParams();
        if (params && params.page) searchParams.append('page', params.page.toString());
        if (params && params.pageSize) searchParams.append('page_size', params.pageSize.toString());
        if (params && params.category) searchParams.append('category', params.category);
        if (params && params.brand) searchParams.append('brand', params.brand);
        
        return `/products?${searchParams.toString()}`;
      },
      providesTags: (result) =>
        result
          ? [
              ...result.products.map(({ id }) => ({ type: 'Product' as const, id })),
              { type: 'Product', id: 'LIST' },
            ]
          : [{ type: 'Product', id: 'LIST' }],
    }),
    getProduct: builder.query<Product, string>({
      query: (id) => `/products/${id}`,
      providesTags: (_result, _error, id) => [{ type: 'Product', id }],
    }),
    searchProducts: builder.query<ProductsResponse, SearchProductsRequest>({
      query: ({ query, page, pageSize }) => {
        const searchParams = new URLSearchParams({ query });
        if (page) searchParams.append('page', page.toString());
        if (pageSize) searchParams.append('page_size', pageSize.toString());
        
        return `/products/search?${searchParams.toString()}`;
      },
      providesTags: [{ type: 'Product', id: 'SEARCH' }],
    }),
  }),
});

export const {
  useGetProductsQuery,
  useGetProductQuery,
  useSearchProductsQuery,
} = productsApi;
