import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import type {
  NotificationsResponse,
  NotificationPreference,
} from '../../types/notification.types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

export const notificationsApi = createApi({
  reducerPath: 'notificationsApi',
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
  tagTypes: ['Notification', 'NotificationPreference'],
  endpoints: (builder) => ({
    getNotifications: builder.query<NotificationsResponse, { userId: string; page?: number; pageSize?: number }>({
      query: ({ userId, page, pageSize }) => {
        const searchParams = new URLSearchParams();
        if (page) searchParams.append('page', page.toString());
        if (pageSize) searchParams.append('page_size', pageSize.toString());
        
        return `/notifications/user/${userId}?${searchParams.toString()}`;
      },
      providesTags: ['Notification'],
    }),
    getUnreadCount: builder.query<{ count: number }, string>({
      query: (userId) => `/notifications/user/${userId}/unread/count`,
      providesTags: ['Notification'],
    }),
    markAsRead: builder.mutation<void, string>({
      query: (notificationId) => ({
        url: `/notifications/${notificationId}/read`,
        method: 'POST',
      }),
      invalidatesTags: ['Notification'],
    }),
    getPreferences: builder.query<NotificationPreference, string>({
      query: (userId) => `/notifications/preferences/${userId}`,
      providesTags: ['NotificationPreference'],
    }),
    updatePreferences: builder.mutation<NotificationPreference, { userId: string; preferences: Partial<NotificationPreference> }>({
      query: ({ userId, preferences }) => ({
        url: `/notifications/preferences/${userId}`,
        method: 'PUT',
        body: preferences,
      }),
      invalidatesTags: ['NotificationPreference'],
    }),
  }),
});

// Export hooks when needed
// export const {
//   useGetNotificationsQuery,
//   useGetUnreadCountQuery,
//   useMarkAsReadMutation,
//   useGetPreferencesQuery,
//   useUpdatePreferencesMutation,
// } = notificationsApi;
