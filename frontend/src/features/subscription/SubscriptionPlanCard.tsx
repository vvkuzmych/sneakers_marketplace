import type { SubscriptionPlan } from '../../types/subscription.types';

interface SubscriptionPlanCardProps {
  plan: SubscriptionPlan;
  isCurrentPlan?: boolean;
  onSelectPlan: (planId: number) => void;
  billingCycle: 'monthly' | 'yearly';
}

export default function SubscriptionPlanCard({
  plan,
  isCurrentPlan = false,
  onSelectPlan,
  billingCycle,
}: SubscriptionPlanCardProps) {
  const isFree = plan.name === 'free';
  const isRecommended = plan.name === 'pro';

  // Calculate yearly price (assume monthly price is stored)
  const monthlyPrice = plan.price_monthly;
  const yearlyPrice = plan.price_yearly || monthlyPrice * 12 * 0.85; // 15% discount
  const displayPrice = billingCycle === 'monthly' ? monthlyPrice : yearlyPrice;
  const pricePerMonth = billingCycle === 'yearly' ? yearlyPrice / 12 : monthlyPrice;

  return (
    <div
      className={`relative rounded-lg border-2 p-6 ${
        isRecommended
          ? 'border-blue-500 shadow-lg'
          : isCurrentPlan
          ? 'border-green-500'
          : 'border-gray-200'
      } ${isFree ? 'bg-gray-50' : 'bg-white'}`}
    >
      {/* Recommended Badge */}
      {isRecommended && !isCurrentPlan && (
        <div className="absolute -top-4 left-1/2 -translate-x-1/2 transform">
          <span className="rounded-full bg-blue-500 px-4 py-1 text-sm font-semibold text-white">
            Most Popular
          </span>
        </div>
      )}

      {/* Current Plan Badge */}
      {isCurrentPlan && (
        <div className="absolute -top-4 left-1/2 -translate-x-1/2 transform">
          <span className="rounded-full bg-green-500 px-4 py-1 text-sm font-semibold text-white">
            Current Plan
          </span>
        </div>
      )}

      {/* Plan Name */}
      <div className="mb-4">
        <h3 className="text-2xl font-bold text-gray-900">{plan.display_name}</h3>
        <p className="mt-2 text-sm text-gray-600">
          {isFree
            ? 'Get started with basic features'
            : plan.name === 'pro'
            ? 'Best for growing sellers'
            : 'Maximum savings for professionals'}
        </p>
      </div>

      {/* Price */}
      <div className="mb-6">
        {isFree ? (
          <div>
            <span className="text-4xl font-bold text-gray-900">$0</span>
            <span className="text-gray-600">/month</span>
          </div>
        ) : (
          <div>
            <div className="flex items-baseline">
              <span className="text-4xl font-bold text-gray-900">
                ${displayPrice.toFixed(0)}
              </span>
              <span className="ml-2 text-gray-600">
                /{billingCycle === 'monthly' ? 'month' : 'year'}
              </span>
            </div>
            {billingCycle === 'yearly' && (
              <div className="mt-1 text-sm text-green-600">
                ${pricePerMonth.toFixed(2)}/month (save 15%)
              </div>
            )}
          </div>
        )}
      </div>

      {/* Seller Fee */}
      <div className="mb-6">
        <div className="flex items-center justify-between rounded-lg bg-blue-50 p-3">
          <span className="text-sm font-medium text-gray-700">Seller Fee</span>
          <span className="text-lg font-bold text-blue-600">
            {plan.seller_fee_percent}%
          </span>
        </div>
        <p className="mt-2 text-xs text-gray-500">
          Platform fee on your sales
        </p>
      </div>

      {/* Features */}
      <ul className="mb-6 space-y-3">
        {plan.features.map((feature, index) => (
          <li key={index} className="flex items-start">
            <svg
              className="mr-3 h-5 w-5 flex-shrink-0 text-green-500"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fillRule="evenodd"
                d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                clipRule="evenodd"
              />
            </svg>
            <span className="text-sm text-gray-700">{feature}</span>
          </li>
        ))}
      </ul>

      {/* Savings Example */}
      {!isFree && (
        <div className="mb-6 rounded-lg bg-green-50 p-3">
          <p className="text-xs font-medium text-gray-700">
            Example: On $10,000 in sales
          </p>
          <p className="mt-1 text-lg font-bold text-green-600">
            Save ${((1.0 - plan.seller_fee_percent / 100) * 10000 - plan.seller_fee_percent).toFixed(0)} vs Free
          </p>
        </div>
      )}

      {/* CTA Button */}
      <button
        onClick={() => onSelectPlan(plan.id)}
        disabled={isCurrentPlan}
        className={`w-full rounded-lg px-4 py-3 font-semibold transition-colors ${
          isCurrentPlan
            ? 'cursor-not-allowed bg-gray-300 text-gray-500'
            : isFree
            ? 'bg-gray-600 text-white hover:bg-gray-700'
            : isRecommended
            ? 'bg-blue-600 text-white hover:bg-blue-700'
            : 'bg-gray-900 text-white hover:bg-gray-800'
        }`}
      >
        {isCurrentPlan
          ? 'Current Plan'
          : isFree
          ? 'Get Started'
          : `Upgrade to ${plan.display_name}`}
      </button>
    </div>
  );
}
