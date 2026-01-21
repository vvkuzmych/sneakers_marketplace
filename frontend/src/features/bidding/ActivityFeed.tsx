interface ActivityItem {
  id: string;
  type: 'bid' | 'ask' | 'match';
  price: number;
  timestamp: string;
  userId?: string;
  userName?: string;
}

interface ActivityFeedProps {
  activities: ActivityItem[];
}

export default function ActivityFeed({ activities }: ActivityFeedProps) {
  const getActivityIcon = (type: string) => {
    switch (type) {
      case 'bid':
        return '📈';
      case 'ask':
        return '📉';
      case 'match':
        return '⚡';
      default:
        return '📊';
    }
  };

  const getActivityColor = (type: string) => {
    switch (type) {
      case 'bid':
        return 'text-green-600';
      case 'ask':
        return 'text-red-600';
      case 'match':
        return 'text-purple-600';
      default:
        return 'text-gray-600';
    }
  };

  const getActivityLabel = (type: string) => {
    switch (type) {
      case 'bid':
        return 'BID placed';
      case 'ask':
        return 'ASK placed';
      case 'match':
        return 'MATCH created';
      default:
        return 'Activity';
    }
  };

  const formatTime = (timestamp: string) => {
    const date = new Date(timestamp);
    return date.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    });
  };

  if (activities.length === 0) {
    return (
      <div className="bg-white rounded-lg border border-gray-200 p-8 text-center">
        <p className="text-gray-500">No activity yet. Be the first to place a bid or ask!</p>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg border border-gray-200">
      <div className="px-6 py-4 border-b border-gray-200">
        <h3 className="text-lg font-bold text-gray-900 flex items-center gap-2">
          📊 Live Activity Feed
          <span className="ml-auto text-sm font-normal text-gray-500">
            {activities.length} {activities.length === 1 ? 'event' : 'events'}
          </span>
        </h3>
      </div>
      
      <div className="max-h-96 overflow-y-auto">
        <div className="divide-y divide-gray-100">
          {activities.map((activity) => (
            <div
              key={activity.id}
              className="px-6 py-3 hover:bg-gray-50 transition-colors"
            >
              <div className="flex items-start gap-3">
                <span className="text-2xl mt-0.5" role="img" aria-label={activity.type}>
                  {getActivityIcon(activity.type)}
                </span>
                
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 mb-1">
                    <span className={`font-semibold ${getActivityColor(activity.type)}`}>
                      {getActivityLabel(activity.type)}
                    </span>
                    {activity.type === 'match' && (
                      <span className="px-2 py-0.5 text-xs font-bold text-white bg-purple-600 rounded-full animate-pulse">
                        NEW MATCH!
                      </span>
                    )}
                  </div>
                  
                  <div className="flex items-center gap-2 text-sm">
                    <span className="font-bold text-gray-900">
                      ${activity.price}
                    </span>
                    {activity.userName && (
                      <span className="text-gray-500">
                        • by {activity.userName}
                      </span>
                    )}
                  </div>
                </div>
                
                <div className="text-xs text-gray-400 whitespace-nowrap">
                  {formatTime(activity.timestamp)}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
      
      {activities.length > 5 && (
        <div className="px-6 py-3 bg-gray-50 text-center text-sm text-gray-500 border-t border-gray-200">
          Scroll to see all activity
        </div>
      )}
    </div>
  );
}
