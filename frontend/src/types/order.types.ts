export interface Order {
  id: string;
  orderNumber: string;
  matchId: string;
  buyerId: string;
  sellerId: string;
  productId: string;
  sizeId: string;
  shippingAddressId: string;
  productPrice: number;
  buyerFee: number;
  sellerFee: number;
  shippingCost: number;
  totalAmount: number;
  status: OrderStatus;
  trackingNumber?: string;
  shippedAt?: string;
  deliveredAt?: string;
  createdAt: string;
  updatedAt: string;
}

export type OrderStatus =
  | 'pending'
  | 'paid'
  | 'processing'
  | 'shipped'
  | 'delivered'
  | 'cancelled'
  | 'refunded'
  | 'disputed'
  | 'completed'
  | 'failed'
  | 'on_hold';

export interface OrderStatusHistory {
  id: string;
  orderId: string;
  status: OrderStatus;
  notes?: string;
  createdAt: string;
}

export interface OrdersResponse {
  orders: Order[];
  total: string;
  page: number;
  pageSize: number;
}

export interface OrdersRequest {
  page?: number;
  pageSize?: number;
  status?: OrderStatus;
}
