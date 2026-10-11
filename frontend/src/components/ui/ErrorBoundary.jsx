import React from 'react';
import Button from './Button';
import Card from './Card';

class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      hasError: false,
      error: null,
      errorInfo: null,
      errorId: null,
    };
  }

  static getDerivedStateFromError(error) {
    // Update state so the next render will show the fallback UI
    return {
      hasError: true,
      errorId: Date.now().toString(36) + Math.random().toString(36).substr(2),
    };
  }

  componentDidCatch(error, errorInfo) {
    // Log error details
    this.setState({
      error,
      errorInfo,
    });

    // Log to external service in production
    if (process.env.NODE_ENV === 'production') {
      this.logErrorToService(error, errorInfo);
    } else {
      console.error('Error Boundary caught an error:', error, errorInfo);
    }
  }

  logErrorToService = (error, errorInfo) => {
    try {
      const errorData = {
        message: error.message,
        stack: error.stack,
        componentStack: errorInfo.componentStack,
        timestamp: new Date().toISOString(),
        url: window.location.href,
        userAgent: navigator.userAgent,
        errorId: this.state.errorId,
      };

      // Send to error tracking service
      fetch('/api/v1/errors/client', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify(errorData),
      }).catch(err => {
        console.error('Failed to log error:', err);
      });
    } catch (err) {
      console.error('Error in error logging:', err);
    }
  };

  handleRetry = () => {
    this.setState({
      hasError: false,
      error: null,
      errorInfo: null,
      errorId: null,
    });
  };

  handleReload = () => {
    window.location.reload();
  };

  handleReportError = () => {
    const subject = encodeURIComponent(`Error Report - ${this.state.errorId}`);
    const body = encodeURIComponent(`
Error ID: ${this.state.errorId}
Timestamp: ${new Date().toISOString()}
URL: ${window.location.href}
Error: ${this.state.error?.message}

Please describe what you were doing when this error occurred:


Technical Details (do not modify):
${this.state.error?.stack}
${this.state.errorInfo?.componentStack}
    `.trim());

    window.open(`mailto:support@ngotracker.com?subject=${subject}&body=${body}`);
  };

  render() {
    if (this.state.hasError) {
      // Custom fallback UI
      if (this.props.fallback) {
        return this.props.fallback(
          this.state.error,
          this.state.errorInfo,
          this.handleRetry
        );
      }

      // Default error UI
      return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50 p-4">
          <Card className="max-w-lg w-full p-8 text-center">
            <div className="mb-6">
              <div className="mx-auto w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mb-4">
                <svg className="w-8 h-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16c-.77.833.192 2.5 1.732 2.5z" />
                </svg>
              </div>
              <h1 className="text-2xl font-bold text-gray-900 mb-2">
                Oops! Something went wrong
              </h1>
              <p className="text-gray-600 mb-4">
                We're sorry, but an unexpected error occurred. Our team has been notified.
              </p>
              {this.state.errorId && (
                <p className="text-sm text-gray-500 mb-6">
                  Error ID: <code className="bg-gray-100 px-1 py-0.5 rounded">{this.state.errorId}</code>
                </p>
              )}
            </div>

            <div className="space-y-3">
              <Button 
                onClick={this.handleRetry} 
                className="w-full"
              >
                Try Again
              </Button>
              
              <Button 
                onClick={this.handleReload} 
                variant="outline"
                className="w-full"
              >
                Refresh Page
              </Button>
              
              <Button 
                onClick={this.handleReportError} 
                variant="ghost"
                className="w-full text-sm"
              >
                Report This Error
              </Button>
            </div>

            {process.env.NODE_ENV === 'development' && (
              <details className="mt-8 text-left">
                <summary className="cursor-pointer text-sm font-medium text-gray-700 mb-2">
                  Technical Details (Development Mode)
                </summary>
                <div className="bg-gray-100 p-4 rounded-lg text-xs overflow-auto max-h-40">
                  <div className="mb-2">
                    <strong>Error:</strong> {this.state.error?.message}
                  </div>
                  <div className="mb-2">
                    <strong>Stack:</strong>
                    <pre className="whitespace-pre-wrap">{this.state.error?.stack}</pre>
                  </div>
                  <div>
                    <strong>Component Stack:</strong>
                    <pre className="whitespace-pre-wrap">{this.state.errorInfo?.componentStack}</pre>
                  </div>
                </div>
              </details>
            )}
          </Card>
        </div>
      );
    }

    return this.props.children;
  }
}

// HOC for wrapping components with error boundaries
export const withErrorBoundary = (Component, errorFallback) => {
  return function WrappedComponent(props) {
    return (
      <ErrorBoundary fallback={errorFallback}>
        <Component {...props} />
      </ErrorBoundary>
    );
  };
};

// Hook for error boundary context
export const useErrorHandler = () => {
  const [error, setError] = React.useState(null);

  const resetError = () => setError(null);

  const captureError = (error, errorInfo = {}) => {
    console.error('Captured error:', error);
    setError({ error, errorInfo });

    // Log to service
    if (process.env.NODE_ENV === 'production') {
      fetch('/api/v1/errors/client', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify({
          message: error.message || error,
          stack: error.stack,
          timestamp: new Date().toISOString(),
          url: window.location.href,
          ...errorInfo,
        }),
      }).catch(console.error);
    }
  };

  // Throw error to be caught by nearest error boundary
  React.useEffect(() => {
    if (error) {
      throw error.error;
    }
  }, [error]);

  return { captureError, resetError };
};

// Async error boundary for handling promise rejections
export class AsyncErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      hasError: false,
      error: null,
    };
  }

  static getDerivedStateFromError(error) {
    return {
      hasError: true,
      error,
    };
  }

  componentDidCatch(error, errorInfo) {
    console.error('Async Error Boundary:', error, errorInfo);
  }

  componentDidMount() {
    // Listen for unhandled promise rejections
    window.addEventListener('unhandledrejection', this.handleUnhandledRejection);
  }

  componentWillUnmount() {
    window.removeEventListener('unhandledrejection', this.handleUnhandledRejection);
  }

  handleUnhandledRejection = (event) => {
    this.setState({
      hasError: true,
      error: event.reason,
    });
    event.preventDefault();
  };

  handleRetry = () => {
    this.setState({
      hasError: false,
      error: null,
    });
  };

  render() {
    if (this.state.hasError) {
      return (
        <div className="p-4 border border-red-200 rounded-lg bg-red-50">
          <div className="flex items-center mb-2">
            <svg className="w-5 h-5 text-red-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <h3 className="text-red-800 font-medium">Async Operation Failed</h3>
          </div>
          <p className="text-red-700 text-sm mb-3">
            {this.state.error?.message || 'An error occurred during an async operation'}
          </p>
          <Button size="sm" onClick={this.handleRetry}>
            Retry
          </Button>
        </div>
      );
    }

    return this.props.children;
  }
}

export default ErrorBoundary;