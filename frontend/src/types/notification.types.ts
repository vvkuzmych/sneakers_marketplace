export interface Notification {
  id: string;
  userId: string;
  type: NotificationType;
  title: string;
  message: string;
  data?: Record<string, any>;
  isRead: boolean;
  emailSent: boolean;
  emailSentAt?: string;
  createdAt: string;
}

export type NotificationType =
  | 'match_created'
  | 'order_created'
  | 'order_shipped'
  | 'order_delivered'
  | 'payment_succeeded'
  | 'payment_failed'
  | 'refund_issued'
  | 'payout_completed';

export interface NotificationPreference {
  userId: string;
  emailEnabled: boolean;
  pushEnabled: boolean;
  matchCreated: boolean;
  orderUpdates: boolean;
  paymentUpdates: boolean;
  marketingEmails: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface NotificationsResponse {
  notifications: Notification[];
  total: string;
  page: number;
  pageSize: number;
}

export interface WebSocketMessage {
  type: 'connected' | 'notification' | 'error';
  data?: any;
  message?: string;
}
