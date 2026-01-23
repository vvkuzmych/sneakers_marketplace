import React, { useEffect, useState } from 'react';

interface FeeBreakdownData {
  sale_price: number;
  buyer_processing_fee: number;
  buyer_shipping_fee: number;
  buyer_total: number;
  seller_transaction_fee: number;
  seller_auth_fee: number;
  seller_shipping_cost: number;
  seller_payout: number;
  platform_revenue: number;
}

interface FeeBreakdownProps {
  vertical: string;
  price: number;
  includeAuth: boolean;
  role: 'buyer' | 'seller';
}

export const FeeBreakdown: React.FC<FeeBreakdownProps> = ({
  vertical,
  price,
  includeAuth,
  role,
}) => {
  const [breakdown, setBreakdown] = useState<FeeBreakdownData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (price > 0) {
      fetchFees();
    } else {
      setBreakdown(null);
    }
  }, [price, vertical, includeAuth]);

  const fetchFees = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch(
        `/api/v1/fees/calculate?vertical=${vertical}&price=${price}&include_auth=${includeAuth}`
      );
      
      if (!response.ok) {
        throw new Error('Failed to fetch fees');
      }
      
      const data = await response.json();
      setBreakdown(data);
    } catch (err) {
      console.error('Failed to fetch fees:', err);
      setError('Failed to calculate fees');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="bg-white rounded-lg border p-4">
        <div className="animate-pulse">
          <div className="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
          <div className="h-4 bg-gray-200 rounded w-1/2"></div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-lg p-4">
        <p className="text-red-600 text-sm">{error}</p>
      </div>
    );
  }

  if (!breakdown) {
    return null;
  }

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-4 shadow-sm">
      <h3 className="font-semibold text-lg mb-3 text-gray-800">
        {role === 'buyer' ? '💳 Total Cost' : '💰 Your Payout'}
      </h3>

      {role === 'buyer' ? (
        <div className="space-y-2">
          <div className="flex justify-between items-center">
            <span className="text-gray-700">Item Price</span>
            <span className="font-mono font-medium">
              ${breakdown.sale_price.toFixed(2)}
            </span>
          </div>

          {breakdown.buyer_processing_fee > 0 && (
            <div className="flex justify-between items-center text-sm text-gray-600">
              <span>Platform Fee (1%)</span>
              <span className="font-mono">
                ${breakdown.buyer_processing_fee.toFixed(2)}
              </span>
            </div>
          )}

          {breakdown.buyer_shipping_fee > 0 && (
            <div className="flex justify-between items-center text-sm text-gray-600">
              <span>Shipping</span>
              <span className="font-mono">
                ${breakdown.buyer_shipping_fee.toFixed(2)}
              </span>
            </div>
          )}

          <div className="border-t pt-2 mt-2 flex justify-between items-center font-bold text-lg">
            <span className="text-gray-900">Total</span>
            <span className="font-mono text-green-600">
              ${breakdown.buyer_total.toFixed(2)}
            </span>
          </div>

          {breakdown.buyer_processing_fee > 0 && (
            <div className="mt-3 pt-3 border-t">
              <p className="text-xs text-gray-500 italic">
                💡 Platform fee helps us maintain a safe and reliable marketplace
              </p>
            </div>
          )}
        </div>
      ) : (
        <div className="space-y-2">
          <div className="flex justify-between items-center">
            <span className="text-gray-700">Sale Price</span>
            <span className="font-mono font-medium">
              ${breakdown.sale_price.toFixed(2)}
            </span>
          </div>

          {breakdown.seller_transaction_fee > 0 && (
            <div className="flex justify-between items-center text-sm text-red-600">
              <span>
                Transaction Fee ({((breakdown.seller_transaction_fee / breakdown.sale_price) * 100).toFixed(1)}%)
              </span>
              <span className="font-mono">
                -${breakdown.seller_transaction_fee.toFixed(2)}
              </span>
            </div>
          )}


          <div className="border-t pt-2 mt-2 flex justify-between items-center font-bold text-lg">
            <span className="text-gray-900">Your Payout</span>
            <span className="font-mono text-green-600">
              ${breakdown.seller_payout.toFixed(2)}
            </span>
          </div>

          <div className="mt-3 pt-3 border-t">
            <p className="text-xs text-gray-500 italic">
              💡 You receive the full sale price! Buyer pays the platform fee.
            </p>
          </div>
        </div>
      )}
    </div>
  );
};
