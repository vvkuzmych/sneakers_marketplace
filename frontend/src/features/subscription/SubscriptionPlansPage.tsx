import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  useGetSubscriptionPlansQuery,
  useGetCurrentSubscriptionQuery,
} from './subscriptionApi';
import SubscriptionPlanCard from './SubscriptionPlanCard';

export default function SubscriptionPlansPage() {
  const navigate = useNavigate();
  const [billingCycle, setBillingCycle] = useState<'monthly' | 'yearly'>('monthly');

  const { data: plans = [], isLoading: plansLoading } = useGetSubscriptionPlansQuery();
  const { data: currentSubscription } = useGetCurrentSubscriptionQuery();

  const handleSelectPlan = (planId: number) => {
    // Navigate to checkout page with selected plan
    navigate(`/subscription/checkout?plan_id=${planId}&billing_cycle=${billingCycle}`);
  };

  if (plansLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="h-12 w-12 animate-spin rounded-full border-4 border-gray-200 border-t-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-12">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="text-center">
          <h1 className="text-4xl font-bold tracking-tight text-gray-900 sm:text-5xl">
            Choose Your Plan
          </h1>
          <p className="mt-4 text-lg text-gray-600">
            Lower your fees and grow your business with our premium plans
          </p>
        </div>

        {/* Billing Cycle Toggle */}
        <div className="mt-8 flex justify-center">
          <div className="inline-flex rounded-lg bg-white p-1 shadow">
            <button
              onClick={() => setBillingCycle('monthly')}
              className={`rounded-md px-6 py-2 text-sm font-medium transition-colors ${
                billingCycle === 'monthly'
                  ? 'bg-blue-600 text-white'
                  : 'text-gray-700 hover:text-gray-900'
              }`}
            >
              Monthly
            </button>
            <button
              onClick={() => setBillingCycle('yearly')}
              className={`rounded-md px-6 py-2 text-sm font-medium transition-colors ${
                billingCycle === 'yearly'
                  ? 'bg-blue-600 text-white'
                  : 'text-gray-700 hover:text-gray-900'
              }`}
            >
              Yearly
              <span className="ml-2 rounded-full bg-green-100 px-2 py-0.5 text-xs font-semibold text-green-800">
                Save 15%
              </span>
            </button>
          </div>
        </div>

        {/* Plans Grid */}
        <div className="mt-12 grid gap-8 md:grid-cols-3">
          {[...plans]
            .sort((a, b) => a.sort_order - b.sort_order)
            .map((plan) => (
              <SubscriptionPlanCard
                key={plan.id}
                plan={plan}
                isCurrentPlan={currentSubscription?.planId === plan.id}
                onSelectPlan={handleSelectPlan}
                billingCycle={billingCycle}
              />
            ))}
        </div>

        {/* FAQ / Additional Info */}
        <div className="mt-16">
          <h2 className="text-center text-2xl font-bold text-gray-900">
            Frequently Asked Questions
          </h2>
          <div className="mx-auto mt-8 max-w-3xl space-y-6">
            <div className="rounded-lg bg-white p-6 shadow">
              <h3 className="text-lg font-semibold text-gray-900">
                Can I change my plan anytime?
              </h3>
              <p className="mt-2 text-gray-600">
                Yes! You can upgrade or downgrade your plan at any time. Changes take effect
                immediately, and we'll prorate the cost.
              </p>
            </div>
            <div className="rounded-lg bg-white p-6 shadow">
              <h3 className="text-lg font-semibold text-gray-900">
                What happens if I downgrade?
              </h3>
              <p className="mt-2 text-gray-600">
                If you downgrade, you'll continue to enjoy your current plan benefits until the
                end of your billing period. The new plan will take effect on your next billing
                date.
              </p>
            </div>
            <div className="rounded-lg bg-white p-6 shadow">
              <h3 className="text-lg font-semibold text-gray-900">
                How are fees calculated?
              </h3>
              <p className="mt-2 text-gray-600">
                Fees are calculated as a percentage of each sale. With Free, you pay 1%. With
                Pro, you pay 0.75%. With Elite, you pay only 0.5%. The savings add up quickly!
              </p>
            </div>
          </div>
        </div>

        {/* Back to Dashboard */}
        <div className="mt-12 text-center">
          <button
            onClick={() => navigate('/')}
            className="text-blue-600 hover:text-blue-700 hover:underline"
          >
            ← Back to Dashboard
          </button>
        </div>
      </div>
    </div>
  );
}
