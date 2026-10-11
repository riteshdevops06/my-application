// Performance optimization utilities

/**
 * Debounce function to limit function calls
 */
export const debounce = (func, delay) => {
  let timeoutId;
  return (...args) => {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => func.apply(null, args), delay);
  };
};

/**
 * Throttle function to limit function calls per time period
 */
export const throttle = (func, limit) => {
  let inThrottle;
  return (...args) => {
    if (!inThrottle) {
      func.apply(null, args);
      inThrottle = true;
      setTimeout(() => inThrottle = false, limit);
    }
  };
};

/**
 * Memoization for expensive calculations
 */
export const memoize = (func, getKey = (...args) => JSON.stringify(args)) => {
  const cache = new Map();
  return (...args) => {
    const key = getKey(...args);
    if (cache.has(key)) {
      return cache.get(key);
    }
    const result = func(...args);
    cache.set(key, result);
    return result;
  };
};

/**
 * Lazy loading utilities
 */
export const createLazyLoader = () => {
  const observer = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        const element = entry.target;
        const src = element.dataset.src;
        if (src) {
          element.src = src;
          element.removeAttribute('data-src');
          observer.unobserve(element);
        }
      }
    });
  });

  return {
    observe: (element) => observer.observe(element),
    unobserve: (element) => observer.unobserve(element),
    disconnect: () => observer.disconnect(),
  };
};

/**
 * Image optimization utilities
 */
export const optimizeImage = (file, options = {}) => {
  const {
    maxWidth = 1920,
    maxHeight = 1080,
    quality = 0.8,
    format = 'image/jpeg'
  } = options;

  return new Promise((resolve) => {
    const canvas = document.createElement('canvas');
    const ctx = canvas.getContext('2d');
    const img = new Image();

    img.onload = () => {
      // Calculate new dimensions
      let { width, height } = img;
      
      if (width > maxWidth) {
        height = (height * maxWidth) / width;
        width = maxWidth;
      }
      
      if (height > maxHeight) {
        width = (width * maxHeight) / height;
        height = maxHeight;
      }

      canvas.width = width;
      canvas.height = height;

      // Draw and compress
      ctx.drawImage(img, 0, 0, width, height);
      
      canvas.toBlob(resolve, format, quality);
    };

    img.src = URL.createObjectURL(file);
  });
};

/**
 * Virtual scrolling for large lists
 */
export const useVirtualScroll = (items, itemHeight, containerHeight) => {
  const [scrollTop, setScrollTop] = React.useState(0);
  
  const startIndex = Math.floor(scrollTop / itemHeight);
  const endIndex = Math.min(
    startIndex + Math.ceil(containerHeight / itemHeight) + 1,
    items.length - 1
  );
  
  const visibleItems = items.slice(startIndex, endIndex + 1);
  
  const totalHeight = items.length * itemHeight;
  const offsetY = startIndex * itemHeight;
  
  return {
    visibleItems,
    totalHeight,
    offsetY,
    startIndex,
    endIndex,
    onScroll: (e) => setScrollTop(e.target.scrollTop),
  };
};

/**
 * Performance monitoring
 */
class PerformanceMonitor {
  constructor() {
    this.metrics = new Map();
    this.observers = [];
  }

  // Start timing
  startTiming(name) {
    this.metrics.set(name, { startTime: performance.now() });
  }

  // End timing
  endTiming(name) {
    const metric = this.metrics.get(name);
    if (metric) {
      metric.endTime = performance.now();
      metric.duration = metric.endTime - metric.startTime;
      return metric.duration;
    }
    return null;
  }

  // Measure function execution time
  measureFunction(name, func) {
    return (...args) => {
      this.startTiming(name);
      const result = func(...args);
      const duration = this.endTiming(name);
      
      if (duration > 100) { // Log slow functions
        console.warn(`Slow function detected: ${name} took ${duration.toFixed(2)}ms`);
      }
      
      return result;
    };
  }

  // Measure async function execution time
  measureAsync(name, func) {
    return async (...args) => {
      this.startTiming(name);
      const result = await func(...args);
      const duration = this.endTiming(name);
      
      if (duration > 500) { // Log slow async operations
        console.warn(`Slow async operation: ${name} took ${duration.toFixed(2)}ms`);
      }
      
      return result;
    };
  }

  // Monitor Core Web Vitals
  observeWebVitals() {
    if ('PerformanceObserver' in window) {
      // Largest Contentful Paint (LCP)
      new PerformanceObserver((list) => {
        list.getEntries().forEach((entry) => {
          console.log('LCP:', entry.startTime);
          this.reportMetric('LCP', entry.startTime);
        });
      }).observe({ type: 'largest-contentful-paint', buffered: true });

      // First Input Delay (FID)
      new PerformanceObserver((list) => {
        list.getEntries().forEach((entry) => {
          const fid = entry.processingStart - entry.startTime;
          console.log('FID:', fid);
          this.reportMetric('FID', fid);
        });
      }).observe({ type: 'first-input', buffered: true });

      // Cumulative Layout Shift (CLS)
      let clsValue = 0;
      new PerformanceObserver((list) => {
        list.getEntries().forEach((entry) => {
          if (!entry.hadRecentInput) {
            clsValue += entry.value;
          }
        });
        console.log('CLS:', clsValue);
        this.reportMetric('CLS', clsValue);
      }).observe({ type: 'layout-shift', buffered: true });
    }
  }

  // Report metrics to analytics service
  reportMetric(name, value) {
    if (process.env.NODE_ENV === 'production') {
      // Send to analytics service
      fetch('/api/v1/analytics/performance', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify({
          metric: name,
          value,
          timestamp: Date.now(),
          url: window.location.pathname,
        }),
      }).catch(console.error);
    }
  }

  // Get all metrics
  getMetrics() {
    const result = {};
    this.metrics.forEach((value, key) => {
      result[key] = value;
    });
    return result;
  }

  // Clear all metrics
  clearMetrics() {
    this.metrics.clear();
  }
}

// Global performance monitor instance
export const performanceMonitor = new PerformanceMonitor();

/**
 * React hooks for performance optimization
 */

// Hook for debounced values
export const useDebounce = (value, delay) => {
  const [debouncedValue, setDebouncedValue] = React.useState(value);

  React.useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return debouncedValue;
};

// Hook for throttled callbacks
export const useThrottle = (callback, delay) => {
  const throttledCallback = React.useMemo(
    () => throttle(callback, delay),
    [callback, delay]
  );

  return throttledCallback;
};

// Hook for intersection observer
export const useIntersectionObserver = (options = {}) => {
  const [entry, setEntry] = React.useState(null);
  const [element, setElement] = React.useState(null);

  React.useEffect(() => {
    if (!element) return;

    const observer = new IntersectionObserver(([entry]) => {
      setEntry(entry);
    }, options);

    observer.observe(element);

    return () => observer.disconnect();
  }, [element, options]);

  return [setElement, entry];
};

// Hook for preloading resources
export const usePreloader = (resources = []) => {
  const [loadedCount, setLoadedCount] = React.useState(0);
  const [loading, setLoading] = React.useState(true);

  React.useEffect(() => {
    if (resources.length === 0) {
      setLoading(false);
      return;
    }

    let loaded = 0;
    
    const handleLoad = () => {
      loaded++;
      setLoadedCount(loaded);
      if (loaded === resources.length) {
        setLoading(false);
      }
    };

    resources.forEach((resource) => {
      if (resource.type === 'image') {
        const img = new Image();
        img.onload = handleLoad;
        img.onerror = handleLoad;
        img.src = resource.src;
      } else if (resource.type === 'script') {
        const script = document.createElement('script');
        script.onload = handleLoad;
        script.onerror = handleLoad;
        script.src = resource.src;
        document.head.appendChild(script);
      }
    });
  }, [resources]);

  return {
    loading,
    progress: resources.length > 0 ? loadedCount / resources.length : 1,
    loadedCount,
    totalCount: resources.length,
  };
};

/**
 * Bundle size optimization utilities
 */

// Code splitting helper
export const loadComponent = (importFn) => {
  return React.lazy(() => {
    return importFn().catch((error) => {
      console.error('Component loading failed:', error);
      return { default: () => <div>Failed to load component</div> };
    });
  });
};

// Dynamic imports with retry
export const importWithRetry = (importFn, retries = 3, delay = 1000) => {
  return new Promise((resolve, reject) => {
    const attempt = (attemptsLeft) => {
      importFn()
        .then(resolve)
        .catch((error) => {
          if (attemptsLeft === 0) {
            reject(error);
          } else {
            setTimeout(() => attempt(attemptsLeft - 1), delay);
          }
        });
    };
    attempt(retries);
  });
};

/**
 * Memory management utilities
 */

// Cleanup manager for preventing memory leaks
export class CleanupManager {
  constructor() {
    this.cleanupFunctions = [];
  }

  add(cleanupFn) {
    this.cleanupFunctions.push(cleanupFn);
  }

  cleanup() {
    this.cleanupFunctions.forEach((fn) => {
      try {
        fn();
      } catch (error) {
        console.error('Cleanup function failed:', error);
      }
    });
    this.cleanupFunctions = [];
  }
}

// Hook for automatic cleanup
export const useCleanup = () => {
  const cleanupManager = React.useRef(new CleanupManager());

  React.useEffect(() => {
    return () => {
      cleanupManager.current.cleanup();
    };
  }, []);

  return (cleanupFn) => {
    cleanupManager.current.add(cleanupFn);
  };
};

/**
 * Caching utilities
 */

// Simple in-memory cache with TTL
class MemoryCache {
  constructor(defaultTTL = 5 * 60 * 1000) { // 5 minutes
    this.cache = new Map();
    this.defaultTTL = defaultTTL;
  }

  set(key, value, ttl = this.defaultTTL) {
    const expiresAt = Date.now() + ttl;
    this.cache.set(key, { value, expiresAt });
  }

  get(key) {
    const item = this.cache.get(key);
    if (!item) return null;

    if (Date.now() > item.expiresAt) {
      this.cache.delete(key);
      return null;
    }

    return item.value;
  }

  has(key) {
    const item = this.cache.get(key);
    if (!item) return false;

    if (Date.now() > item.expiresAt) {
      this.cache.delete(key);
      return false;
    }

    return true;
  }

  delete(key) {
    return this.cache.delete(key);
  }

  clear() {
    this.cache.clear();
  }

  // Cleanup expired entries
  cleanup() {
    const now = Date.now();
    for (const [key, item] of this.cache.entries()) {
      if (now > item.expiresAt) {
        this.cache.delete(key);
      }
    }
  }
}

export const memoryCache = new MemoryCache();

// Auto-cleanup expired cache entries every 5 minutes
setInterval(() => {
  memoryCache.cleanup();
}, 5 * 60 * 1000);

/**
 * Initialize performance monitoring
 */
export const initializePerformanceMonitoring = () => {
  // Monitor Core Web Vitals
  performanceMonitor.observeWebVitals();

  // Monitor bundle loading
  if ('PerformanceObserver' in window) {
    new PerformanceObserver((list) => {
      list.getEntries().forEach((entry) => {
        if (entry.name.includes('.js') || entry.name.includes('.css')) {
          console.log(`Resource loaded: ${entry.name} (${entry.duration.toFixed(2)}ms)`);
        }
      });
    }).observe({ type: 'resource', buffered: true });
  }

  // Monitor memory usage
  if ('memory' in performance) {
    const logMemoryUsage = () => {
      const memory = performance.memory;
      console.log('Memory usage:', {
        used: Math.round(memory.usedJSHeapSize / 1024 / 1024),
        total: Math.round(memory.totalJSHeapSize / 1024 / 1024),
        limit: Math.round(memory.jsHeapSizeLimit / 1024 / 1024),
      });
    };

    // Log memory usage every 30 seconds in development
    if (process.env.NODE_ENV === 'development') {
      setInterval(logMemoryUsage, 30000);
    }
  }
};