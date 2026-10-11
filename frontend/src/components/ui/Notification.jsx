import { useState, useEffect, createContext, useContext } from 'react';
import { createPortal } from 'react-dom';
import { cn } from '../../utils/cn';

// Notification Context
const NotificationContext = createContext({});

// Notification Provider
const NotificationProvider = ({ children, position = 'top-right', maxNotifications = 5 }) => {
  const [notifications, setNotifications] = useState([]);

  // Add notification
  const addNotification = (notification) => {
    const id = Date.now() + Math.random();
    const newNotification = {
      id,
      type: 'info',
      duration: 5000,
      ...notification,
    };

    setNotifications(prev => {
      const updated = [newNotification, ...prev];
      return updated.slice(0, maxNotifications);
    });

    // Auto remove after duration
    if (newNotification.duration > 0) {
      setTimeout(() => {
        removeNotification(id);
      }, newNotification.duration);
    }

    return id;
  };

  // Remove notification
  const removeNotification = (id) => {
    setNotifications(prev => prev.filter(n => n.id !== id));
  };

  // Clear all notifications
  const clearNotifications = () => {
    setNotifications([]);
  };

  // Notification functions
  const success = (message, options = {}) => 
    addNotification({ ...options, type: 'success', message });
    
  const error = (message, options = {}) => 
    addNotification({ ...options, type: 'error', message, duration: 0 });
    
  const warning = (message, options = {}) => 
    addNotification({ ...options, type: 'warning', message });
    
  const info = (message, options = {}) => 
    addNotification({ ...options, type: 'info', message });

  const contextValue = {
    notifications,
    addNotification,
    removeNotification,
    clearNotifications,
    success,
    error,
    warning,
    info,
  };

  return (
    <NotificationContext.Provider value={contextValue}>
      {children}
      <NotificationContainer position={position} />
    </NotificationContext.Provider>
  );
};

// Notification Container
const NotificationContainer = ({ position }) => {
  const { notifications } = useContext(NotificationContext);

  const positionClasses = {
    'top-right': 'top-4 right-4',
    'top-left': 'top-4 left-4',
    'top-center': 'top-4 left-1/2 transform -translate-x-1/2',
    'bottom-right': 'bottom-4 right-4',
    'bottom-left': 'bottom-4 left-4',
    'bottom-center': 'bottom-4 left-1/2 transform -translate-x-1/2',
  };

  if (notifications.length === 0) return null;

  return createPortal(
    <div className={cn('fixed z-50 space-y-2', positionClasses[position])}>
      {notifications.map((notification) => (
        <NotificationItem key={notification.id} notification={notification} />
      ))}
    </div>,
    document.body
  );
};

// Individual Notification Item
const NotificationItem = ({ notification }) => {
  const { removeNotification } = useContext(NotificationContext);
  const [isVisible, setIsVisible] = useState(false);
  const [isRemoving, setIsRemoving] = useState(false);

  useEffect(() => {
    // Animate in
    const timer = setTimeout(() => setIsVisible(true), 10);
    return () => clearTimeout(timer);
  }, []);

  const handleRemove = () => {
    setIsRemoving(true);
    setTimeout(() => {
      removeNotification(notification.id);
    }, 300);
  };

  const typeStyles = {
    success: {
      bg: 'bg-green-50 border-green-200',
      icon: '✅',
      iconBg: 'bg-green-500',
      text: 'text-green-800',
    },
    error: {
      bg: 'bg-red-50 border-red-200',
      icon: '❌',
      iconBg: 'bg-red-500',
      text: 'text-red-800',
    },
    warning: {
      bg: 'bg-yellow-50 border-yellow-200',
      icon: '⚠️',
      iconBg: 'bg-yellow-500',
      text: 'text-yellow-800',
    },
    info: {
      bg: 'bg-blue-50 border-blue-200',
      icon: 'ℹ️',
      iconBg: 'bg-blue-500',
      text: 'text-blue-800',
    },
  };

  const style = typeStyles[notification.type] || typeStyles.info;

  return (
    <div
      className={cn(
        'max-w-sm w-full shadow-lg rounded-lg pointer-events-auto border transition-all duration-300 transform',
        style.bg,
        isVisible && !isRemoving ? 'translate-x-0 opacity-100' : 'translate-x-full opacity-0',
        isRemoving && '-translate-x-full opacity-0'
      )}
    >
      <div className="p-4">
        <div className="flex items-start">
          <div className="flex-shrink-0">
            <div className={cn('flex items-center justify-center w-8 h-8 rounded-full text-white text-sm', style.iconBg)}>
              {style.icon}
            </div>
          </div>
          
          <div className="ml-3 w-0 flex-1 pt-0.5">
            {notification.title && (
              <p className={cn('text-sm font-medium', style.text)}>
                {notification.title}
              </p>
            )}
            <p className={cn('text-sm', style.text, notification.title && 'mt-1')}>
              {notification.message}
            </p>
            
            {notification.actions && (
              <div className="mt-3 flex space-x-2">
                {notification.actions.map((action, index) => (
                  <button
                    key={index}
                    onClick={() => {
                      action.onClick();
                      if (action.closeOnClick !== false) {
                        handleRemove();
                      }
                    }}
                    className={cn(
                      'text-sm font-medium rounded-md px-3 py-1',
                      action.style === 'primary' 
                        ? `${style.iconBg} text-white hover:opacity-90`
                        : `border ${style.text} hover:bg-white hover:bg-opacity-50`
                    )}
                  >
                    {action.label}
                  </button>
                ))}
              </div>
            )}
          </div>
          
          <div className="ml-4 flex-shrink-0 flex">
            <button
              onClick={handleRemove}
              className={cn('rounded-md inline-flex text-gray-400 hover:text-gray-500 focus:outline-none')}
            >
              <span className="sr-only">Close</span>
              <svg className="w-5 h-5" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
              </svg>
            </button>
          </div>
        </div>
      </div>
      
      {/* Progress bar for auto-dismiss */}
      {notification.duration > 0 && (
        <div className="h-1 bg-gray-200">
          <div
            className={cn('h-full', style.iconBg.replace('bg-', 'bg-opacity-60 bg-'))}
            style={{
              animation: `notification-progress ${notification.duration}ms linear forwards`
            }}
          />
        </div>
      )}
    </div>
  );
};

// Hook to use notifications
const useNotifications = () => {
  const context = useContext(NotificationContext);
  if (!context) {
    throw new Error('useNotifications must be used within a NotificationProvider');
  }
  return context;
};

// Real-time notification system
const useRealTimeNotifications = (userId) => {
  const notifications = useNotifications();
  
  useEffect(() => {
    // This would connect to WebSocket or Server-Sent Events
    // For now, we'll simulate with polling
    
    const checkForUpdates = async () => {
      try {
        const response = await fetch(`/api/v1/notifications/${userId}`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
          },
        });
        
        if (response.ok) {
          const data = await response.json();
          
          // Process new notifications
          data.notifications?.forEach(notif => {
            notifications.addNotification({
              type: notif.type,
              title: notif.title,
              message: notif.message,
              duration: notif.persistent ? 0 : 5000,
              actions: notif.actions,
            });
          });
        }
      } catch (error) {
        console.error('Error fetching notifications:', error);
      }
    };

    // Check immediately and then every 30 seconds
    checkForUpdates();
    const interval = setInterval(checkForUpdates, 30000);

    return () => clearInterval(interval);
  }, [userId, notifications]);
};

// Progress notification for long operations
const useProgressNotification = () => {
  const notifications = useNotifications();
  
  const showProgress = (title, initialProgress = 0) => {
    const id = notifications.addNotification({
      type: 'info',
      title,
      message: `Progress: ${initialProgress}%`,
      duration: 0,
      progress: initialProgress,
    });

    const updateProgress = (progress, message) => {
      notifications.updateNotification(id, {
        message: message || `Progress: ${progress}%`,
        progress,
      });
    };

    const complete = (message = 'Complete!') => {
      notifications.updateNotification(id, {
        type: 'success',
        message,
        progress: 100,
        duration: 3000,
      });
    };

    const error = (message = 'An error occurred') => {
      notifications.updateNotification(id, {
        type: 'error',
        message,
        duration: 0,
      });
    };

    return { updateProgress, complete, error };
  };

  return { showProgress };
};

// Add CSS for progress animation
const style = document.createElement('style');
style.textContent = `
  @keyframes notification-progress {
    from { width: 100%; }
    to { width: 0%; }
  }
`;
document.head.appendChild(style);

export default NotificationProvider;
export { useNotifications, useRealTimeNotifications, useProgressNotification };