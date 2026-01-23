import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useAppSelector } from '../../app/hooks';
import { useGetProductQuery } from '../products/productsApi';
import { useGetMarketPriceQuery, usePlaceBidMutation, usePlaceAskMutation, useGetBidsQuery, useGetAsksQuery } from './biddingApi';
import { websocketService, WS_MESSAGE_TYPES } from '../../services/websocket';
import ActivityFeed from './ActivityFeed';
import { FeeBreakdown } from '../../components/fees/FeeBreakdown';
import styles from './BiddingPage.module.css';

interface MarketUpdate {
  highestBid?: number;
  lowestAsk?: number;
  totalBids?: string;
  totalAsks?: string;
}

interface ActivityItem {
  id: string;
  type: 'bid' | 'ask' | 'match';
  price: number;
  timestamp: string;
  userId?: string;
  userName?: string;
}

export default function BiddingPage() {
  const { productId } = useParams<{ productId: string }>();
  const token = useAppSelector((state) => state.auth.accessToken);
  
  // Validate productId
  const productIdNum = productId ? parseInt(productId, 10) : NaN;
  if (!productId || isNaN(productIdNum)) {
    return (
      <div className="container mx-auto p-6">
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <h2 className="text-red-800 font-semibold">❌ Invalid Product ID</h2>
          <p className="text-red-600 mt-2">
            The product ID in the URL is invalid. Please go back to the{' '}
            <a href="/products" className="underline">products page</a> and select a product.
          </p>
        </div>
      </div>
    );
  }
  
  const { data: product, isLoading: productLoading } = useGetProductQuery(productId!);
  const { data: initialMarketPrice } = useGetMarketPriceQuery(
    { productId: productId!, sizeId: '1' },
    { skip: !productId }
  );
  
  // Load historical bids and asks
  const { data: bidsData } = useGetBidsQuery(
    { productId: productId!, sizeId: '1' },
    { skip: !productId }
  );
  const { data: asksData } = useGetAsksQuery(
    { productId: productId!, sizeId: '1' },
    { skip: !productId }
  );

  const [placeBid] = usePlaceBidMutation();
  const [placeAsk] = usePlaceAskMutation();

  // Real-time market data
  const [marketPrice, setMarketPrice] = useState<MarketUpdate>({});
  const [bidPrice, setBidPrice] = useState<string>('');
  const [askPrice, setAskPrice] = useState<string>('');
  const [wsStatus, setWsStatus] = useState<string>('Disconnected');
  const [notification, setNotification] = useState<string>('');
  const [activities, setActivities] = useState<ActivityItem[]>([]);

  const addActivity = (type: 'bid' | 'ask' | 'match', price: number, data?: any) => {
    const newActivity: ActivityItem = {
      id: `${Date.now()}-${Math.random()}`,
      type,
      price,
      timestamp: new Date().toISOString(),
      userId: data?.userId,
      userName: data?.userName || 'User',
    };
    setActivities((prev) => [newActivity, ...prev]);
  };

  const showNotification = (message: string) => {
    setNotification(message);
    setTimeout(() => setNotification(''), 5000);
  };

  // Initialize WebSocket connection
  useEffect(() => {
    if (!token) {
      console.warn('⚠️ No token available, skipping WebSocket connection');
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setWsStatus('Disconnected');
      return;
    }

    console.log('🚀 Initializing WebSocket with token...');
    
    websocketService.connect(token)
      .then(() => {
        setWsStatus('Connected');
        console.log('✅ WebSocket ready for real-time updates');
      })
      .catch((error) => {
        console.error('❌ WebSocket connection failed:', error);
        setWsStatus('Failed');
      });

    // Subscribe to market updates
    const handleMarketUpdate = (data: MarketUpdate) => {
      console.log('📊 Market price updated:', data);
      setMarketPrice(data);
    };

    const handleBidPlaced = (data: any) => {
      console.log('💰 New BID placed:', data);
      const price = data.bid?.price || data.price;
      showNotification(`New BID: $${price}`);
      addActivity('bid', price, data);
    };

    const handleAskPlaced = (data: any) => {
      console.log('💰 New ASK placed:', data);
      const price = data.ask?.price || data.price;
      showNotification(`New ASK: $${price}`);
      addActivity('ask', price, data);
    };

    const handleMatchCreated = (data: any) => {
      console.log('⚡ MATCH CREATED:', data);
      showNotification(`🎉 MATCH! Sold at $${data.price}`);
      addActivity('match', data.price, data);
    };

    websocketService.on(WS_MESSAGE_TYPES.MARKET_PRICE_UPDATED, handleMarketUpdate);
    websocketService.on(WS_MESSAGE_TYPES.BID_PLACED, handleBidPlaced);
    websocketService.on(WS_MESSAGE_TYPES.ASK_PLACED, handleAskPlaced);
    websocketService.on(WS_MESSAGE_TYPES.MATCH_CREATED, handleMatchCreated);

    return () => {
      websocketService.off(WS_MESSAGE_TYPES.MARKET_PRICE_UPDATED);
      websocketService.off(WS_MESSAGE_TYPES.BID_PLACED);
      websocketService.off(WS_MESSAGE_TYPES.ASK_PLACED);
      websocketService.off(WS_MESSAGE_TYPES.MATCH_CREATED);
    };
  }, [token]);

  // Update market price from initial query
  useEffect(() => {
    if (initialMarketPrice) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setMarketPrice(initialMarketPrice);
    }
  }, [initialMarketPrice]);

  // Load historical bids and asks into activity feed
  useEffect(() => {
    const historicalActivities: ActivityItem[] = [];
    
    // Add historical bids
    if (bidsData?.bids) {
      bidsData.bids.forEach((bid) => {
        historicalActivities.push({
          id: `bid-${bid.id}`,
          type: 'bid',
          price: bid.price,
          timestamp: bid.createdAt,
          userId: bid.userId,
          userName: 'User',
        });
      });
    }
    
    // Add historical asks
    if (asksData?.asks) {
      asksData.asks.forEach((ask) => {
        historicalActivities.push({
          id: `ask-${ask.id}`,
          type: 'ask',
          price: ask.price,
          timestamp: ask.createdAt,
          userId: ask.userId,
          userName: 'User',
        });
      });
    }
    
    // Sort by timestamp (newest first)
    historicalActivities.sort((a, b) => 
      new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
    );
    
    if (historicalActivities.length > 0) {
      console.warn(`📚 Loaded ${historicalActivities.length} historical activities`);
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setActivities(historicalActivities);
    }
  }, [bidsData, asksData]);

  const handlePlaceBid = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!productId || !bidPrice) return;

    console.warn('🔍 DEBUG: Placing BID with:', {
      productId: productIdNum,
      productIdRaw: productId,
      sizeId: 1,
      price: parseFloat(bidPrice),
    });

    try {
      const result = await placeBid({
        productId: productIdNum,
        sizeId: 1, // TODO: Add size selector
        price: parseFloat(bidPrice),
        quantity: 1,
      }).unwrap();

      console.log('✅ BID placed:', result);
      const price = parseFloat(bidPrice);
      setBidPrice('');

      if (result.match) {
        showNotification(`🎉 INSTANT MATCH! Bought at $${result.match.price}`);
        addActivity('match', result.match.price, result);
      } else {
        showNotification(`✅ BID placed at $${price}`);
        addActivity('bid', price, result);
      }
    } catch (error: any) {
      console.error('❌ Failed to place BID:', error);
      showNotification(`❌ Error: ${error?.data?.error || 'Failed to place BID'}`);
    }
  };

  const handlePlaceAsk = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!productId || !askPrice) return;

    console.warn('🔍 DEBUG: Placing ASK with:', {
      productId: productIdNum,
      productIdRaw: productId,
      sizeId: 1,
      price: parseFloat(askPrice),
    });

    try {
      const result = await placeAsk({
        productId: productIdNum,
        sizeId: 1,
        price: parseFloat(askPrice),
        quantity: 1,
      }).unwrap();

      console.log('✅ ASK placed:', result);
      const price = parseFloat(askPrice);
      setAskPrice('');

      if (result.match) {
        showNotification(`🎉 INSTANT MATCH! Sold at $${result.match.price}`);
        addActivity('match', result.match.price, result);
      } else {
        showNotification(`✅ ASK placed at $${price}`);
        addActivity('ask', price, result);
      }
    } catch (error: any) {
      console.error('❌ Failed to place ASK:', error);
      showNotification(`❌ Error: ${error?.data?.error || 'Failed to place ASK'}`);
    }
  };

  if (productLoading) {
    return <div className={styles.loading}>Loading product...</div>;
  }

  if (!product) {
    return <div className={styles.error}>Product not found</div>;
  }

  const highestBid = marketPrice.highestBid || 0;
  const lowestAsk = marketPrice.lowestAsk || 0;
  const spread = lowestAsk && highestBid ? lowestAsk - highestBid : 0;

  return (
    <div className={styles.container}>
      {/* Header */}
      <div className={styles.header}>
        <div className={styles.productInfo}>
          <h1 className={styles.title}>{product.name}</h1>
          <p className={styles.subtitle}>
            {product.brand} • {product.model} • {product.color}
          </p>
        </div>
        <div className={styles.wsStatus}>
          <span className={`${styles.indicator} ${wsStatus === 'Connected' ? styles.connected : styles.disconnected}`}>
            ●
          </span>
          <span className={styles.statusText}>{wsStatus}</span>
        </div>
      </div>

      {/* Notification */}
      {notification && (
        <div className={styles.notification}>
          {notification}
        </div>
      )}

      {/* Market Price */}
      <div className={styles.marketPrice}>
        <h2 className={styles.sectionTitle}>📊 Live Market Price</h2>
        <div className={styles.priceGrid}>
          <div className={styles.priceCard}>
            <div className={styles.priceLabel}>Highest BID</div>
            <div className={styles.priceValue}>
              {highestBid > 0 ? `$${highestBid}` : '—'}
            </div>
            <div className={styles.priceCount}>
              {marketPrice.totalBids || 0} bids
            </div>
          </div>

          <div className={styles.spread}>
            <div className={styles.spreadLabel}>Spread</div>
            <div className={styles.spreadValue}>
              {spread > 0 ? `$${spread}` : '—'}
            </div>
          </div>

          <div className={styles.priceCard}>
            <div className={styles.priceLabel}>Lowest ASK</div>
            <div className={styles.priceValue}>
              {lowestAsk > 0 ? `$${lowestAsk}` : '—'}
            </div>
            <div className={styles.priceCount}>
              {marketPrice.totalAsks || 0} asks
            </div>
          </div>
        </div>
      </div>

      {/* Trading Forms */}
      <div className={styles.tradingForms}>
        {/* BID Form (Buy) */}
        <div className={styles.formCard}>
          <h3 className={styles.formTitle}>💰 Place BID (Buy)</h3>
          <form onSubmit={handlePlaceBid}>
            <div className={styles.formGroup}>
              <label className={styles.label}>Your BID Price</label>
              <div className={styles.inputWrapper}>
                <span className={styles.currency}>$</span>
                <input
                  type="number"
                  className={styles.input}
                  placeholder="200"
                  value={bidPrice}
                  onChange={(e) => setBidPrice(e.target.value)}
                  min="0"
                  step="0.01"
                  required
                />
              </div>
              <p className={styles.hint}>
                {lowestAsk > 0 && parseFloat(bidPrice) >= lowestAsk
                  ? '⚡ This will create INSTANT MATCH!'
                  : 'Enter your maximum buy price'}
              </p>
            </div>
            
            {/* Fee Breakdown for Buyer */}
            {bidPrice && parseFloat(bidPrice) > 0 && (
              <div className="mb-4">
                <FeeBreakdown
                  vertical="sneakers"
                  price={parseFloat(bidPrice)}
                  includeAuth={true}
                  role="buyer"
                />
              </div>
            )}
            
            <button type="submit" className={`${styles.button} ${styles.buttonBuy}`}>
              Place BID
            </button>
          </form>
        </div>

        {/* ASK Form (Sell) */}
        <div className={styles.formCard}>
          <h3 className={styles.formTitle}>💵 Place ASK (Sell)</h3>
          <form onSubmit={handlePlaceAsk}>
            <div className={styles.formGroup}>
              <label className={styles.label}>Your ASK Price</label>
              <div className={styles.inputWrapper}>
                <span className={styles.currency}>$</span>
                <input
                  type="number"
                  className={styles.input}
                  placeholder="220"
                  value={askPrice}
                  onChange={(e) => setAskPrice(e.target.value)}
                  min="0"
                  step="0.01"
                  required
                />
              </div>
              <p className={styles.hint}>
                {highestBid > 0 && parseFloat(askPrice) <= highestBid
                  ? '⚡ This will create INSTANT MATCH!'
                  : 'Enter your minimum sell price'}
              </p>
            </div>
            
            {/* Fee Breakdown for Seller */}
            {askPrice && parseFloat(askPrice) > 0 && (
              <div className="mb-4">
                <FeeBreakdown
                  vertical="sneakers"
                  price={parseFloat(askPrice)}
                  includeAuth={true}
                  role="seller"
                />
              </div>
            )}
            
            <button type="submit" className={`${styles.button} ${styles.buttonSell}`}>
              Place ASK
            </button>
          </form>
        </div>
      </div>

      {/* Activity Feed */}
      <div className={styles.activitySection}>
        <ActivityFeed activities={activities} />
      </div>

      {/* Info Section */}
      <div className={styles.infoBox}>
        <h3>ℹ️ How it works</h3>
        <ul>
          <li><strong>BID:</strong> Your maximum buy price. If a seller's ASK is ≤ your BID, instant match!</li>
          <li><strong>ASK:</strong> Your minimum sell price. If a buyer's BID is ≥ your ASK, instant match!</li>
          <li><strong>Real-time:</strong> All prices update live via WebSocket. No refresh needed!</li>
          <li><strong>Match:</strong> When BID ≥ ASK, trade happens automatically at the ASK price.</li>
        </ul>
      </div>
    </div>
  );
}
