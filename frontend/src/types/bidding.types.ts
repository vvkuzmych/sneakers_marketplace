export interface Bid {
  id: string;
  userId: string;
  productId: string;
  sizeId: string;
  price: number;
  quantity: number;
  status: 'active' | 'matched' | 'cancelled' | 'expired';
  expiresAt: string;
  createdAt: string;
  updatedAt: string;
}

export interface Ask {
  id: string;
  userId: string;
  productId: string;
  sizeId: string;
  price: number;
  quantity: number;
  status: 'active' | 'matched' | 'cancelled' | 'expired';
  expiresAt: string;
  createdAt: string;
  updatedAt: string;
}

export interface Match {
  id: string;
  bidId: string;
  askId: string;
  buyerId: string;
  sellerId: string;
  productId: string;
  sizeId: string;
  price: number;
  quantity: number;
  status: 'pending' | 'completed' | 'cancelled';
  createdAt: string;
}

export interface MarketPrice {
  highestBid?: number;
  lowestAsk?: number;
  lastSalePrice?: number;
  totalBids: string;
  totalAsks: string;
}

export interface PlaceBidRequest {
  productId: number;
  sizeId: number;
  price: number;
  quantity: number;
  expiresInHours?: number;
}

export interface PlaceAskRequest {
  productId: number;
  sizeId: number;
  price: number;
  quantity: number;
  expiresInHours?: number;
}

export interface BidResponse {
  bid: Bid;
  match?: Match;
}

export interface AskResponse {
  ask: Ask;
  match?: Match;
}

export interface BidsResponse {
  bids: Bid[];
  total: string;
}

export interface AsksResponse {
  asks: Ask[];
  total: string;
}
