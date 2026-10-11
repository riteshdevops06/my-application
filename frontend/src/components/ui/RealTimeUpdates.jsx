import { useState, useEffect, useRef, createContext, useContext } from 'react';

// Real-time Context
const RealTimeContext = createContext({});

// Real-time Provider
const RealTimeProvider = ({ children, wsUrl, userId, onConnect, onDisconnect, onError }) => {
  const [isConnected, setIsConnected] = useState(false);
  const [connectionStatus, setConnectionStatus] = useState('disconnected');
  const [lastUpdate, setLastUpdate] = useState(null);
  const [subscribers, setSubscribers] = useState(new Map());
  
  const wsRef = useRef(null);
  const reconnectTimeoutRef = useRef(null);
  const heartbeatRef = useRef(null);
  const reconnectAttempts = useRef(0);
  const maxReconnectAttempts = 5;

  // Initialize WebSocket connection
  const connect = () => {
    if (wsRef.current?.readyState === WebSocket.OPEN) return;

    setConnectionStatus('connecting');
    
    try {
      const token = localStorage.getItem('token');
      const url = `${wsUrl}?token=${token}&userId=${userId}`;
      
      wsRef.current = new WebSocket(url);

      wsRef.current.onopen = () => {
        console.log('WebSocket connected');
        setIsConnected(true);
        setConnectionStatus('connected');
        reconnectAttempts.current = 0;
        
        // Start heartbeat
        startHeartbeat();
        
        if (onConnect) onConnect();
      };

      wsRef.current.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          handleMessage(data);
        } catch (error) {
          console.error('Error parsing WebSocket message:', error);
        }
      };

      wsRef.current.onclose = (event) => {
        console.log('WebSocket disconnected:', event.code, event.reason);
        setIsConnected(false);
        setConnectionStatus('disconnected');
        stopHeartbeat();
        
        if (onDisconnect) onDisconnect(event);
        
        // Attempt reconnection if not intentional
        if (event.code !== 1000 && reconnectAttempts.current < maxReconnectAttempts) {
          scheduleReconnect();
        }
      };

      wsRef.current.onerror = (error) => {
        console.error('WebSocket error:', error);
        setConnectionStatus('error');
        if (onError) onError(error);
      };
      
    } catch (error) {
      console.error('Error creating WebSocket:', error);
      setConnectionStatus('error');
      if (onError) onError(error);
    }
  };

  // Disconnect WebSocket
  const disconnect = () => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
    }
    
    stopHeartbeat();
    
    if (wsRef.current) {
      wsRef.current.close(1000, 'User disconnect');
    }
  };

  // Schedule reconnection
  const scheduleReconnect = () => {
    if (reconnectTimeoutRef.current) return;
    
    reconnectAttempts.current++;
    const delay = Math.min(1000 * Math.pow(2, reconnectAttempts.current), 30000);
    
    console.log(`Scheduling reconnect attempt ${reconnectAttempts.current} in ${delay}ms`);
    setConnectionStatus('reconnecting');
    
    reconnectTimeoutRef.current = setTimeout(() => {
      reconnectTimeoutRef.current = null;
      connect();
    }, delay);
  };

  // Start heartbeat
  const startHeartbeat = () => {
    heartbeatRef.current = setInterval(() => {
      if (wsRef.current?.readyState === WebSocket.OPEN) {
        wsRef.current.send(JSON.stringify({ type: 'ping' }));
      }
    }, 30000); // 30 seconds
  };

  // Stop heartbeat
  const stopHeartbeat = () => {
    if (heartbeatRef.current) {
      clearInterval(heartbeatRef.current);
      heartbeatRef.current = null;
    }
  };

  // Handle incoming messages
  const handleMessage = (data) => {
    setLastUpdate(Date.now());
    
    const { type, channel, payload } = data;
    
    // Handle different message types
    switch (type) {
      case 'pong':
        // Heartbeat response
        break;
        
      case 'update':
        // Broadcast to subscribers
        if (subscribers.has(channel)) {
          subscribers.get(channel).forEach(callback => {
            try {
              callback(payload);
            } catch (error) {
              console.error('Error in subscriber callback:', error);
            }
          });
        }
        break;
        
      case 'notification':
        // Handle real-time notifications
        if (window.NotificationAPI) {
          window.NotificationAPI.addNotification(payload);
        }
        break;
        
      default:
        console.log('Unknown message type:', type);
    }
  };

  // Subscribe to channel updates
  const subscribe = (channel, callback) => {
    if (!subscribers.has(channel)) {
      subscribers.set(channel, new Set());
    }
    
    subscribers.get(channel).add(callback);
    
    // Send subscription request
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({
        type: 'subscribe',
        channel
      }));
    }
    
    // Return unsubscribe function
    return () => {
      const channelSubs = subscribers.get(channel);
      if (channelSubs) {
        channelSubs.delete(callback);
        
        if (channelSubs.size === 0) {
          subscribers.delete(channel);
          
          // Send unsubscription request
          if (wsRef.current?.readyState === WebSocket.OPEN) {
            wsRef.current.send(JSON.stringify({
              type: 'unsubscribe',
              channel
            }));
          }
        }
      }
    };
  };

  // Send message
  const send = (message) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message));
      return true;
    }
    return false;
  };

  // Initialize connection on mount
  useEffect(() => {
    if (wsUrl && userId) {
      connect();
    }
    
    return () => {
      disconnect();
    };
  }, [wsUrl, userId]);

  const contextValue = {
    isConnected,
    connectionStatus,
    lastUpdate,
    connect,
    disconnect,
    subscribe,
    send,
  };

  return (
    <RealTimeContext.Provider value={contextValue}>
      {children}
    </RealTimeContext.Provider>
  );
};

// Hook to use real-time context
const useRealTime = () => {
  const context = useContext(RealTimeContext);
  if (!context) {
    throw new Error('useRealTime must be used within a RealTimeProvider');
  }
  return context;
};

// Campaign progress tracker
const useCampaignProgress = (campaignId) => {
  const [progress, setProgress] = useState(null);
  const [donations, setDonations] = useState([]);
  const [milestones, setMilestones] = useState([]);
  const { subscribe, isConnected } = useRealTime();

  useEffect(() => {
    if (!campaignId || !isConnected) return;

    // Subscribe to campaign updates
    const unsubscribe = subscribe(`campaign:${campaignId}`, (data) => {
      switch (data.type) {
        case 'donation':
          setDonations(prev => [data.donation, ...prev.slice(0, 9)]); // Keep last 10
          setProgress(prev => ({
            ...prev,
            currentAmount: data.newTotal,
            donorCount: data.donorCount,
            percentage: (data.newTotal / prev?.goalAmount) * 100
          }));
          break;
          
        case 'milestone':
          setMilestones(prev => 
            prev.map(m => m.id === data.milestone.id ? data.milestone : m)
          );
          break;
          
        case 'progress':
          setProgress(data.progress);
          break;
      }
    });

    // Fetch initial data
    fetchCampaignProgress();

    return unsubscribe;
  }, [campaignId, isConnected, subscribe]);

  const fetchCampaignProgress = async () => {
    try {
      const response = await fetch(`/api/v1/campaigns/${campaignId}/progress`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });
      
      if (response.ok) {
        const data = await response.json();
        setProgress(data.progress);
        setDonations(data.recentDonations || []);
        setMilestones(data.milestones || []);
      }
    } catch (error) {
      console.error('Error fetching campaign progress:', error);
    }
  };

  return {
    progress,
    donations,
    milestones,
    refresh: fetchCampaignProgress
  };
};

// Live donor feed component
const LiveDonorFeed = ({ campaignId, maxItems = 5, className }) => {
  const { donations } = useCampaignProgress(campaignId);

  if (!donations.length) {
    return (
      <div className={className}>
        <p className="text-gray-500 text-sm">No recent donations</p>
      </div>
    );
  }

  return (
    <div className={className}>
      <h4 className="text-sm font-medium text-gray-900 mb-3">Recent Donations</h4>
      <div className="space-y-2">
        {donations.slice(0, maxItems).map((donation, index) => (
          <div
            key={donation.id}
            className={`flex items-center justify-between p-2 rounded-lg transition-all duration-500 ${
              index === 0 ? 'bg-green-50 border border-green-200' : 'bg-gray-50'
            }`}
          >
            <div>
              <p className="text-sm font-medium text-gray-900">
                {donation.anonymous ? 'Anonymous Donor' : donation.donorName}
              </p>
              <p className="text-xs text-gray-500">
                {new Date(donation.createdAt).toLocaleTimeString()}
              </p>
            </div>
            <div className="text-right">
              <p className="text-sm font-semibold text-green-600">
                ${donation.amount.toLocaleString()}
              </p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

// Connection status indicator
const ConnectionStatus = ({ className }) => {
  const { connectionStatus, isConnected, lastUpdate } = useRealTime();

  const statusConfig = {
    connected: {
      color: 'bg-green-500',
      text: 'Connected',
      textColor: 'text-green-700'
    },
    connecting: {
      color: 'bg-yellow-500',
      text: 'Connecting...',
      textColor: 'text-yellow-700'
    },
    reconnecting: {
      color: 'bg-orange-500',
      text: 'Reconnecting...',
      textColor: 'text-orange-700'
    },
    disconnected: {
      color: 'bg-gray-500',
      text: 'Disconnected',
      textColor: 'text-gray-700'
    },
    error: {
      color: 'bg-red-500',
      text: 'Connection Error',
      textColor: 'text-red-700'
    }
  };

  const config = statusConfig[connectionStatus] || statusConfig.disconnected;

  return (
    <div className={`flex items-center space-x-2 ${className}`}>
      <div className={`w-2 h-2 rounded-full ${config.color} ${isConnected ? 'animate-pulse' : ''}`} />
      <span className={`text-xs ${config.textColor}`}>
        {config.text}
        {lastUpdate && isConnected && (
          <span className="ml-1 text-gray-500">
            (last update: {new Date(lastUpdate).toLocaleTimeString()})
          </span>
        )}
      </span>
    </div>
  );
};

// Auto-refreshing data hook
const useAutoRefresh = (fetcher, interval = 30000, dependencies = []) => {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const { isConnected } = useRealTime();

  useEffect(() => {
    let intervalId;

    const fetchData = async () => {
      try {
        setLoading(true);
        setError(null);
        const result = await fetcher();
        setData(result);
      } catch (err) {
        setError(err);
      } finally {
        setLoading(false);
      }
    };

    // Fetch immediately
    fetchData();

    // Set up interval only if not connected to real-time
    if (!isConnected) {
      intervalId = setInterval(fetchData, interval);
    }

    return () => {
      if (intervalId) {
        clearInterval(intervalId);
      }
    };
  }, [fetcher, interval, isConnected, ...dependencies]);

  return { data, loading, error, refresh: () => fetcher().then(setData) };
};

export default RealTimeProvider;
export { 
  useRealTime, 
  useCampaignProgress, 
  LiveDonorFeed, 
  ConnectionStatus, 
  useAutoRefresh 
};