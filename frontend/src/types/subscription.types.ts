// Subscription Types

export interface SubscriptionPlan {
  id: number;
  name: string;
  display_name: string;
  description?: string;
  price_monthly: number;
  price_yearly: number;
  seller_fee_percent: number;
  buyer_fee_percent: number;
  features: string[];
  max_active_listings?: number | null;
  max_monthly_transactions?: number | null;
  is_active: boolean;
  sort_order: number;
  stripe_price_id_monthly?: string;
  stripe_price_id_yearly?: string;
  created_at: string;
  updated_at: string;
}

export interface UserSubscription {
  id: number;
  userId: number;
  planId: number;
  status: 'active' | 'canceled' | 'past_due' | 'trialing' | 'incomplete' | 'incomplete_expired';
  currentPeriodStart: string;
  currentPeriodEnd: string;
  cancelAtPeriodEnd: boolean;
  billingCycle: 'monthly' | 'yearly';
  stripeCustomerId?: string;
  stripeSubscriptionId?: string;
  createdAt: string;
  updatedAt: string;
}

export interface UserSubscriptionWithPlan extends UserSubscription {
  plan: SubscriptionPlan;
}

export interface SubscriptionTransaction {
  id: number;
  subscriptionId: number;
  userId: number;
  transactionType: 'subscription' | 'upgrade' | 'downgrade' | 'renewal' | 'cancellation';
  status: 'pending' | 'succeeded' | 'failed' | 'refunded';
  amount: number;
  currency: string;
  stripePaymentIntentId?: string;
  createdAt: string;
  updatedAt: string;
}

export interface FeeBreakdown {
  salePrice: number;
  sellerTransactionFee: number;
  sellerAuthFee: number;
  sellerShippingCost: number;
  sellerPayout: number;
  buyerProcessingFee: number;
  buyerShippingFee: number;
  buyerTotal: number;
  platformRevenue: number;
}

export interface SubscribeRequest {
  planId: number;
  billingCycle: 'monthly' | 'yearly';
  paymentMethodId: string;
}

export interface SubscribeResponse {
  subscription: UserSubscriptionWithPlan;
  clientSecret?: string; // For 3D Secure
  requiresAction: boolean;
}

export interface CancelSubscriptionRequest {
  subscriptionId: number;
  cancelAtPeriodEnd: boolean;
}

export interface UpdateSubscriptionRequest {
  subscriptionId: number;
  newPlanId: number;
  newBillingCycle?: 'monthly' | 'yearly';
}

export interface FeeSavingsCalculation {
  currentPlan: string;
  currentFeePercent: number;
  targetPlan: string;
  targetFeePercent: number;
  salePrice: number;
  currentFee: number;
  targetFee: number;
  savings: number;
  savingsPercent: number;
}
