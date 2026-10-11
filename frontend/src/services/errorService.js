// Error handling service for centralized error management
import React from 'react';

class ErrorService {
  constructor() {
    this.errorHandlers = new Map();
    this.globalErrorHandler = null;
    this.setupGlobalErrorHandlers();
  }

  // Setup global error handlers
  setupGlobalErrorHandlers() {
    // Capture unhandled JavaScript errors
    window.addEventListener('error', (event) => {
      this.handleError({
        type: 'javascript_error',
        message: event.message,
        filename: event.filename,
        lineno: event.lineno,
        colno: event.colno,
        stack: event.error?.stack,
        timestamp: new Date().toISOString(),
      });
    });

    // Capture unhandled promise rejections
    window.addEventListener('unhandledrejection', (event) => {
      this.handleError({
        type: 'unhandled_promise_rejection',
        message: event.reason?.message || 'Unhandled promise rejection',
        stack: event.reason?.stack,
        timestamp: new Date().toISOString(),
      });
    });

    // Capture network errors
    window.addEventListener('offline', () => {
      this.handleError({
        type: 'network_error',
        message: 'Network connection lost',
        timestamp: new Date().toISOString(),
      });
    });
  }

  // Register error handler for specific error types
  registerErrorHandler(errorType, handler) {
    if (!this.errorHandlers.has(errorType)) {
      this.errorHandlers.set(errorType, []);
    }
    this.errorHandlers.get(errorType).push(handler);
  }

  // Set global error handler
  setGlobalErrorHandler(handler) {
    this.globalErrorHandler = handler;
  }

  // Handle errors
  handleError(error) {
    // Enrich error with context
    const enrichedError = this.enrichError(error);
    
    // Log to console in development
    if (process.env.NODE_ENV === 'development') {
      console.error('Error Service:', enrichedError);
    }

    // Call specific error handlers
    const handlers = this.errorHandlers.get(error.type) || [];
    handlers.forEach(handler => {
      try {
        handler(enrichedError);
      } catch (handlerError) {
        console.error('Error in error handler:', handlerError);
      }
    });

    // Call global error handler
    if (this.globalErrorHandler) {
      try {
        this.globalErrorHandler(enrichedError);
      } catch (handlerError) {
        console.error('Error in global error handler:', handlerError);
      }
    }

    // Send to logging service
    this.logError(enrichedError);
  }

  // Enrich error with additional context
  enrichError(error) {
    return {
      ...error,
      id: this.generateErrorId(),
      url: window.location.href,
      userAgent: navigator.userAgent,
      timestamp: error.timestamp || new Date().toISOString(),
      sessionId: this.getSessionId(),
      userId: this.getUserId(),
      buildVersion: process.env.REACT_APP_VERSION || 'unknown',
      environment: process.env.NODE_ENV || 'unknown',
    };
  }

  // Generate unique error ID
  generateErrorId() {
    return Date.now().toString(36) + Math.random().toString(36).substr(2);
  }

  // Get session ID
  getSessionId() {
    let sessionId = sessionStorage.getItem('sessionId');
    if (!sessionId) {
      sessionId = this.generateErrorId();
      sessionStorage.setItem('sessionId', sessionId);
    }
    return sessionId;
  }

  // Get user ID
  getUserId() {
    try {
      const token = localStorage.getItem('token');
      if (token) {
        const payload = JSON.parse(atob(token.split('.')[1]));
        return payload.user_id || payload.sub;
      }
    } catch (e) {
      // Token parsing failed
    }
    return null;
  }

  // Log error to external service
  async logError(error) {
    try {
      // Only log in production or when explicitly enabled
      if (process.env.NODE_ENV !== 'production' && !process.env.REACT_APP_LOG_ERRORS) {
        return;
      }

      await fetch('/api/v1/errors/client', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify(error),
      });
    } catch (logError) {
      console.error('Failed to log error to server:', logError);
    }
  }

  // Handle API errors
  handleApiError(error, context = {}) {
    const apiError = {
      type: 'api_error',
      message: error.message || 'API request failed',
      status: error.status,
      statusText: error.statusText,
      url: error.url || context.url,
      method: error.method || context.method,
      requestData: context.requestData,
      responseData: error.responseData,
      timestamp: new Date().toISOString(),
    };

    this.handleError(apiError);
    return apiError;
  }

  // Handle validation errors
  handleValidationError(errors, context = {}) {
    const validationError = {
      type: 'validation_error',
      message: 'Form validation failed',
      errors: errors,
      formName: context.formName,
      timestamp: new Date().toISOString(),
    };

    this.handleError(validationError);
    return validationError;
  }

  // Handle authentication errors
  handleAuthError(error, context = {}) {
    const authError = {
      type: 'auth_error',
      message: error.message || 'Authentication failed',
      action: context.action,
      timestamp: new Date().toISOString(),
    };

    this.handleError(authError);
    return authError;
  }

  // Create error from HTTP response
  static async createFromResponse(response, requestContext = {}) {
    let responseData = null;
    let message = response.statusText;

    try {
      responseData = await response.json();
      message = responseData.error?.message || responseData.message || message;
    } catch (e) {
      // Response is not JSON
    }

    return {
      message,
      status: response.status,
      statusText: response.statusText,
      url: response.url,
      method: requestContext.method,
      requestData: requestContext.data,
      responseData,
    };
  }

  // Retry mechanism for failed operations
  async retry(operation, maxAttempts = 3, delay = 1000) {
    for (let attempt = 1; attempt <= maxAttempts; attempt++) {
      try {
        return await operation();
      } catch (error) {
        if (attempt === maxAttempts) {
          this.handleError({
            type: 'retry_failed',
            message: `Operation failed after ${maxAttempts} attempts`,
            originalError: error.message,
            attempts: maxAttempts,
            timestamp: new Date().toISOString(),
          });
          throw error;
        }

        // Wait before retry
        await new Promise(resolve => setTimeout(resolve, delay * attempt));
      }
    }
  }

  // Get error statistics
  getErrorStats() {
    const errors = this.getStoredErrors();
    const stats = {
      total: errors.length,
      byType: {},
      recent: errors.filter(e => 
        new Date(e.timestamp) > new Date(Date.now() - 24 * 60 * 60 * 1000)
      ).length,
    };

    errors.forEach(error => {
      stats.byType[error.type] = (stats.byType[error.type] || 0) + 1;
    });

    return stats;
  }

  // Store error locally for debugging
  storeError(error) {
    try {
      const stored = this.getStoredErrors();
      stored.push(error);
      
      // Keep only last 50 errors
      const limited = stored.slice(-50);
      
      localStorage.setItem('errorLog', JSON.stringify(limited));
    } catch (e) {
      // Storage failed
    }
  }

  // Get stored errors
  getStoredErrors() {
    try {
      const stored = localStorage.getItem('errorLog');
      return stored ? JSON.parse(stored) : [];
    } catch (e) {
      return [];
    }
  }

  // Clear stored errors
  clearStoredErrors() {
    localStorage.removeItem('errorLog');
  }
}

// Create singleton instance
const errorService = new ErrorService();

// Export factory functions for common error types
export const createApiError = (response, context) => 
  errorService.handleApiError(response, context);

export const createValidationError = (errors, context) => 
  errorService.handleValidationError(errors, context);

export const createAuthError = (error, context) => 
  errorService.handleAuthError(error, context);

// Export utility functions
export const captureError = (error, context = {}) => 
  errorService.handleError({ ...error, ...context });

export const withErrorHandling = (fn) => async (...args) => {
  try {
    return await fn(...args);
  } catch (error) {
    errorService.handleError({
      type: 'function_error',
      message: error.message,
      stack: error.stack,
      functionName: fn.name,
      arguments: args.length,
      timestamp: new Date().toISOString(),
    });
    throw error;
  }
};

export const retry = (operation, maxAttempts, delay) => 
  errorService.retry(operation, maxAttempts, delay);

// React hook for error handling
export const useErrorHandler = () => {
  const handleError = React.useCallback((error, context = {}) => {
    errorService.handleError({ ...error, ...context });
  }, []);

  const handleApiError = React.useCallback((response, context = {}) => {
    return errorService.handleApiError(response, context);
  }, []);

  const handleValidationError = React.useCallback((errors, context = {}) => {
    return errorService.handleValidationError(errors, context);
  }, []);

  return {
    handleError,
    handleApiError,
    handleValidationError,
    retry: errorService.retry.bind(errorService),
  };
};

// Higher-order component for error handling
export const withErrorBoundary = (WrappedComponent, fallbackComponent = null) => {
  return function ErrorBoundaryWrapper(props) {
    return React.createElement(
      ErrorBoundary,
      { fallback: fallbackComponent },
      React.createElement(WrappedComponent, props)
    );
  };
};

export default errorService;