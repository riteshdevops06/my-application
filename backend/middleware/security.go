package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// SecurityHeaders adds security-related headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")
		
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		
		// XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Content Security Policy
		csp := []string{
			"default-src 'self'",
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'",
			"style-src 'self' 'unsafe-inline'",
			"img-src 'self' data: https:",
			"connect-src 'self'",
			"font-src 'self'",
			"object-src 'none'",
			"media-src 'self'",
			"form-action 'self'",
		}
		c.Header("Content-Security-Policy", strings.Join(csp, "; "))
		
		// HSTS (only in production with HTTPS)
		if gin.Mode() == gin.ReleaseMode {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}
		
		// Permissions Policy
		permissions := []string{
			"camera=()",
			"microphone=()",
			"geolocation=()",
			"accelerometer=()",
			"gyroscope=()",
			"magnetometer=()",
			"payment=()",
		}
		c.Header("Permissions-Policy", strings.Join(permissions, ", "))
		
		c.Next()
	}
}

// RateLimiter implements rate limiting middleware
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	requests int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		requests: requests,
		window:   window,
	}
}

// Middleware returns the rate limiting middleware
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client identifier (IP address)
		clientID := c.ClientIP()
		
		// Get or create limiter for this client
		limiter, exists := rl.limiters[clientID]
		if !exists {
			limiter = rate.NewLimiter(rate.Every(rl.window), rl.requests)
			rl.limiters[clientID] = limiter
		}
		
		// Check if request is allowed
		if !limiter.Allow() {
			c.Header("X-RateLimit-Limit", strconv.Itoa(rl.requests))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(rl.window).Unix(), 10))
			
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{
					"type":    "RATE_LIMIT_EXCEEDED",
					"message": "Too many requests. Please try again later.",
				},
			})
			c.Abort()
			return
		}
		
		// Add rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(rl.requests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(rl.requests-1))
		
		c.Next()
	}
}

// IPWhitelist middleware allows only whitelisted IP addresses
func IPWhitelist(allowedIPs []string) gin.HandlerFunc {
	ipMap := make(map[string]bool)
	for _, ip := range allowedIPs {
		ipMap[ip] = true
	}
	
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		
		if !ipMap[clientIP] {
			c.JSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"type":    "IP_NOT_ALLOWED",
					"message": "Access denied from this IP address",
				},
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// InputSanitization sanitizes user input to prevent injection attacks
func InputSanitization() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sanitize query parameters
		for key, values := range c.Request.URL.Query() {
			for i, value := range values {
				values[i] = sanitizeInput(value)
			}
			c.Request.URL.Query()[key] = values
		}
		
		// Sanitize form data
		if c.Request.Form != nil {
			for key, values := range c.Request.Form {
				for i, value := range values {
					values[i] = sanitizeInput(value)
				}
				c.Request.Form[key] = values
			}
		}
		
		c.Next()
	}
}

// sanitizeInput removes potentially dangerous characters
func sanitizeInput(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")
	
	// Remove script tags (basic XSS prevention)
	input = strings.ReplaceAll(input, "<script", "&lt;script")
	input = strings.ReplaceAll(input, "</script>", "&lt;/script&gt;")
	
	// Remove potentially dangerous SQL characters (basic SQL injection prevention)
	dangerous := []string{"'", "\"", ";", "--", "/*", "*/", "xp_", "sp_"}
	for _, char := range dangerous {
		input = strings.ReplaceAll(input, char, "")
	}
	
	return input
}

// RequestSizeLimit limits the size of request bodies
func RequestSizeLimit(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": gin.H{
					"type":    "REQUEST_TOO_LARGE",
					"message": fmt.Sprintf("Request body too large. Maximum size is %d bytes", maxSize),
				},
			})
			c.Abort()
			return
		}
		
		// Set max bytes reader
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		
		c.Next()
	}
}

// SecureHeaders adds additional security headers for API responses
func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		
		// Prevent caching of sensitive data
		if strings.Contains(c.Request.URL.Path, "/auth/") || 
		   strings.Contains(c.Request.URL.Path, "/admin/") {
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}
		
		// Add security headers to JSON responses
		if strings.Contains(c.GetHeader("Content-Type"), "application/json") {
			c.Header("X-Content-Type-Options", "nosniff")
		}
	}
}

// APIKeyAuth middleware for API key authentication
func APIKeyAuth(validAPIKeys map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"type":    "MISSING_API_KEY",
					"message": "API key is required",
				},
			})
			c.Abort()
			return
		}
		
		// Validate API key
		if clientName, valid := validAPIKeys[apiKey]; valid {
			c.Set("api_client", clientName)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"type":    "INVALID_API_KEY",
					"message": "Invalid API key",
				},
			})
			c.Abort()
		}
	}
}

// RequestTimeout middleware adds timeout to requests
func RequestTimeout(timeout time.Duration) gin.HandlerFunc {
	return gin.TimeoutWithHandler(timeout, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusRequestTimeout)
		w.Write([]byte(`{"error":{"type":"REQUEST_TIMEOUT","message":"Request timeout"}}`))
	}))
}

// HTTPS redirect middleware
func HTTPSRedirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Forwarded-Proto") == "http" {
			httpsURL := "https://" + c.Request.Host + c.Request.RequestURI
			c.Redirect(http.StatusPermanentRedirect, httpsURL)
			c.Abort()
			return
		}
		c.Next()
	}
}

// Maintenance mode middleware
func MaintenanceMode(enabled bool, message string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if enabled {
			// Allow health checks during maintenance
			if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/api/v1/health" {
				c.Next()
				return
			}
			
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": gin.H{
					"type":    "MAINTENANCE_MODE",
					"message": message,
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequestID middleware adds unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		
		c.Next()
	}
}

// generateRequestID creates a unique request ID
func generateRequestID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().Unix())
}

// TrustedProxies sets trusted proxy IPs
func TrustedProxies(trustedProxies []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// This would typically be handled by gin.Engine.SetTrustedProxies()
		// but we can add additional proxy validation here if needed
		c.Next()
	}
}

// CORS with security considerations
func SecureCORS(allowedOrigins []string, allowCredentials bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		
		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}
		
		if allowed || len(allowedOrigins) == 0 {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		
		if allowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
		c.Header("Access-Control-Max-Age", "86400") // 24 hours
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}