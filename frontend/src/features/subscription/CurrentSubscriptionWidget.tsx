import { Link } from 'react-router-dom';
import { useGetCurrentSubscriptionQuery } from './subscriptionApi';

export default function CurrentSubscriptionWidget() {
  const { data: subscription, isLoading, error } = useGetCurrentSubscriptionQuery();

  if (isLoading) {
    return (
      <div className="animate-pulse rounded-lg bg-white p-6 shadow">
        <div className="h-4 w-1/3 rounded bg-gray-200"></div>
        <div className="mt-4 h-8 w-1/2 rounded bg-gray-200"></div>
      </div>
    );
  }

  if (error || !subscription) {
    // User has no subscription (Free tier by default)
    return (
      <div className="rounded-lg bg-gradient-to-r from-blue-500 to-blue-600 p-6 text-white shadow-lg">
        <div className="flex items-start justify-between">
          <div>
            <h3 className="text-lg font-semibold">Free Plan</h3>
            <p className="mt-2 text-sm text-blue-100">
              You're currently on the Free plan (1% seller fee)
            </p>
            <p className="mt-4 text-sm font-medium">
              Upgrade to Pro or Elite to save on every sale! 💰
            </p>
          </div>
          <svg
            className="h-12 w-12 opacity-80"
            fill="currentColor"
            viewBox="0 0 20 20"
          >
            <path d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-11a1 1 0 10-2 0v2H7a1 1 0 100 2h2v2a1 1 0 102 0v-2h2a1 1 0 100-2h-2V7z" />
          </svg>
        </div>
        <Link
          to="/subscription/plans"
          className="mt-4 inline-block rounded-lg bg-white px-6 py-2 font-semibold text-blue-600 transition-colors hover:bg-blue-50"
        >
          View Plans
        </Link>
      </div>
    );
  }

  const { plan, status, currentPeriodEnd, cancelAtPeriodEnd } = subscription;
  const isActive = status === 'active' || status === 'trialing';
  const isFree = plan.name === 'free';

  // Format date
  const endDate = new Date(currentPeriodEnd).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });

  return (
    <div
      className={`rounded-lg p-6 shadow-lg ${
        isFree
          ? 'bg-gradient-to-r from-gray-400 to-gray-500'
          : plan.name === 'pro'
          ? 'bg-gradient-to-r from-blue-500 to-blue-600'
          : 'bg-gradient-to-r from-purple-500 to-purple-600'
      } text-white`}
    >
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="flex items-center">
            <h3 className="text-2xl font-bold">{plan.display_name}</h3>
            {isActive && !cancelAtPeriodEnd && (
              <span className="ml-3 rounded-full bg-green-400 px-3 py-1 text-xs font-semibold text-white">
                Active
              </span>
            )}
            {cancelAtPeriodEnd && (
              <span className="ml-3 rounded-full bg-yellow-400 px-3 py-1 text-xs font-semibold text-gray-900">
                Canceling
              </span>
            )}
          </div>

          <div className="mt-4 grid grid-cols-2 gap-4">
            <div>
              <p className="text-sm opacity-90">Seller Fee</p>
              <p className="text-3xl font-bold">{plan.seller_fee_percent}%</p>
            </div>
            <div>
              <p className="text-sm opacity-90">Status</p>
              <p className="text-lg font-semibold capitalize">{status}</p>
            </div>
          </div>

          <div className="mt-4">
            <p className="text-sm opacity-90">
              {cancelAtPeriodEnd ? 'Active until' : 'Renews on'}
            </p>
            <p className="text-sm font-medium">{endDate}</p>
          </div>

          {!isFree && !cancelAtPeriodEnd && plan.seller_fee_percent < 1.0 && (
            <div className="mt-4 rounded-lg bg-white bg-opacity-20 p-3">
              <p className="text-xs font-medium">💰 Monthly Savings</p>
              <p className="text-sm">
                On $10,000 sales: Save{' '}
                <span className="font-bold">
                  ${((1.0 - plan.seller_fee_percent / 100) * 10000 - plan.seller_fee_percent).toFixed(0)}
                </span>{' '}
                vs Free
              </p>
            </div>
          )}
        </div>

        {/* Plan Icon */}
        <div className="ml-4">
          {isFree ? (
            <svg
              className="h-16 w-16 opacity-80"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-11a1 1 0 10-2 0v2H7a1 1 0 100 2h2v2a1 1 0 102 0v-2h2a1 1 0 100-2h-2V7z" />
            </svg>
          ) : plan.name === 'pro' ? (
            <svg
              className="h-16 w-16 opacity-80"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
            </svg>
          ) : (
            <svg
              className="h-16 w-16 opacity-80"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fillRule="evenodd"
                d="M5 2a1 1 0 011 1v1h1a1 1 0 010 2H6v1a1 1 0 01-2 0V6H3a1 1 0 010-2h1V3a1 1 0 011-1zm0 10a1 1 0 011 1v1h1a1 1 0 110 2H6v1a1 1 0 11-2 0v-1H3a1 1 0 110-2h1v-1a1 1 0 011-1zM12 2a1 1 0 01.967.744L14.146 7.2 17.5 9.134a1 1 0 010 1.732l-3.354 1.935-1.18 4.455a1 1 0 01-1.933 0L9.854 12.8 6.5 10.866a1 1 0 010-1.732l3.354-1.935 1.18-4.455A1 1 0 0112 2z"
                clipRule="evenodd"
              />
            </svg>
          )}
        </div>
      </div>

      {/* Actions */}
      <div className="mt-6 flex space-x-3">
        {!isFree && (
          <Link
            to="/subscription/manage"
            className="rounded-lg bg-white bg-opacity-20 px-4 py-2 text-sm font-medium transition-colors hover:bg-opacity-30"
          >
            Manage Subscription
          </Link>
        )}
        <Link
          to="/subscription/plans"
          className="rounded-lg bg-white px-4 py-2 text-sm font-medium text-gray-900 transition-colors hover:bg-gray-100"
        >
          {isFree ? 'Upgrade Plan' : 'View All Plans'}
        </Link>
      </div>
    </div>
  );
}
