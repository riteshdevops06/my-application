package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// AppError represents a structured application error
type AppError struct {
	Type       string            `json:"type"`
	Message    string            `json:"message"`
	Code       int               `json:"code"`
	Details    map[string]string `json:"details,omitempty"`
	Timestamp  string            `json:"timestamp"`
	RequestID  string            `json:"request_id,omitempty"`
	StackTrace string            `json:"stack_trace,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Success bool     `json:"success"`
	Error   AppError `json:"error"`
}

// Common error types
const (
	ErrorTypeValidation     = "VALIDATION_ERROR"
	ErrorTypeAuthentication = "AUTHENTICATION_ERROR"
	ErrorTypeAuthorization  = "AUTHORIZATION_ERROR"
	ErrorTypeNotFound       = "NOT_FOUND"
	ErrorTypeConflict       = "CONFLICT"
	ErrorTypeInternal       = "INTERNAL_ERROR"
	ErrorTypeRateLimit      = "RATE_LIMIT_ERROR"
	ErrorTypeBadRequest     = "BAD_REQUEST"
	ErrorTypeTimeout        = "TIMEOUT_ERROR"
	ErrorTypeServiceUnavail = "SERVICE_UNAVAILABLE"
)

// NewAppError creates a new application error
func NewAppError(errorType, message string, code int) *AppError {
	return &AppError{
		Type:      errorType,
		Message:   message,
		Code:      code,
		Timestamp: GetCurrentTimestamp(),
	}
}

// NewValidationError creates a validation error with field details
func NewValidationError(details map[string]string) *AppError {
	return &AppError{
		Type:      ErrorTypeValidation,
		Message:   "Validation failed",
		Code:      http.StatusBadRequest,
		Details:   details,
		Timestamp: GetCurrentTimestamp(),
	}
}

// NewAuthenticationError creates an authentication error
func NewAuthenticationError(message string) *AppError {
	if message == "" {
		message = "Authentication required"
	}
	return NewAppError(ErrorTypeAuthentication, message, http.StatusUnauthorized)
}

// NewAuthorizationError creates an authorization error
func NewAuthorizationError(message string) *AppError {
	if message == "" {
		message = "Access denied"
	}
	return NewAppError(ErrorTypeAuthorization, message, http.StatusForbidden)
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string) *AppError {
	message := fmt.Sprintf("%s not found", resource)
	return NewAppError(ErrorTypeNotFound, message, http.StatusNotFound)
}

// NewConflictError creates a conflict error
func NewConflictError(message string) *AppError {
	return NewAppError(ErrorTypeConflict, message, http.StatusConflict)
}

// NewInternalError creates an internal server error
func NewInternalError(message string) *AppError {
	if message == "" {
		message = "Internal server error occurred"
	}
	return NewAppError(ErrorTypeInternal, message, http.StatusInternalServerError)
}

// NewRateLimitError creates a rate limit error
func NewRateLimitError() *AppError {
	return NewAppError(ErrorTypeRateLimit, "Rate limit exceeded", http.StatusTooManyRequests)
}

// NewBadRequestError creates a bad request error
func NewBadRequestError(message string) *AppError {
	return NewAppError(ErrorTypeBadRequest, message, http.StatusBadRequest)
}

// WriteErrorResponse writes an error response to the client
func WriteErrorResponse(c *gin.Context, err error) {
	var appErr *AppError
	
	switch e := err.(type) {
	case *AppError:
		appErr = e
	case validator.ValidationErrors:
		// Handle validator errors
		details := make(map[string]string)
		for _, validationErr := range e {
			field := strings.ToLower(validationErr.Field())
			details[field] = getValidationErrorMessage(validationErr)
		}
		appErr = NewValidationError(details)
	default:
		// Unknown error - log and return internal error
		log.Printf("Unknown error: %v", err)
		appErr = NewInternalError("An unexpected error occurred")
	}
	
	// Add request ID if available
	if requestID := c.GetString("request_id"); requestID != "" {
		appErr.RequestID = requestID
	}
	
	// Add stack trace in development
	if gin.Mode() == gin.DebugMode {
		appErr.StackTrace = getStackTrace()
	}
	
	// Log error for monitoring
	logError(c, appErr)
	
	response := ErrorResponse{
		Success: false,
		Error:   *appErr,
	}
	
	c.JSON(appErr.Code, response)
}

// WriteJSONError writes a JSON error response
func WriteJSONError(c *gin.Context, errorType, message string, code int) {
	err := NewAppError(errorType, message, code)
	WriteErrorResponse(c, err)
}

// WriteValidationError writes validation errors
func WriteValidationError(c *gin.Context, details map[string]string) {
	err := NewValidationError(details)
	WriteErrorResponse(c, err)
}

// RecoverFromPanic recovers from panics and converts them to errors
func RecoverFromPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				
				// Log stack trace
				buf := make([]byte, 1024)
				n := runtime.Stack(buf, false)
				log.Printf("Stack trace: %s", buf[:n])
				
				// Return internal server error
				appErr := NewInternalError("Internal server error")
				if gin.Mode() == gin.DebugMode {
					appErr.Details = map[string]string{
						"panic": fmt.Sprintf("%v", err),
					}
					appErr.StackTrace = string(buf[:n])
				}
				
				WriteErrorResponse(c, appErr)
				c.Abort()
			}
		}()
		c.Next()
	}
}

// getValidationErrorMessage returns user-friendly validation error messages
func getValidationErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return fmt.Sprintf("Must be at least %s characters", err.Param())
	case "max":
		return fmt.Sprintf("Must be no more than %s characters", err.Param())
	case "len":
		return fmt.Sprintf("Must be exactly %s characters", err.Param())
	case "gte":
		return fmt.Sprintf("Must be at least %s", err.Param())
	case "lte":
		return fmt.Sprintf("Must be no more than %s", err.Param())
	case "gt":
		return fmt.Sprintf("Must be greater than %s", err.Param())
	case "lt":
		return fmt.Sprintf("Must be less than %s", err.Param())
	case "alphanum":
		return "Must contain only letters and numbers"
	case "alpha":
		return "Must contain only letters"
	case "numeric":
		return "Must be a number"
	case "url":
		return "Must be a valid URL"
	case "uuid":
		return "Must be a valid UUID"
	default:
		return fmt.Sprintf("Validation failed for %s", err.Tag())
	}
}

// getStackTrace returns the current stack trace
func getStackTrace() string {
	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// logError logs errors for monitoring and debugging
func logError(c *gin.Context, err *AppError) {
	logData := map[string]interface{}{
		"error_type":   err.Type,
		"error_code":   err.Code,
		"message":      err.Message,
		"method":       c.Request.Method,
		"path":         c.Request.URL.Path,
		"ip":           c.ClientIP(),
		"user_agent":   c.Request.UserAgent(),
		"timestamp":    err.Timestamp,
	}
	
	if err.RequestID != "" {
		logData["request_id"] = err.RequestID
	}
	
	if userID := c.GetString("user_id"); userID != "" {
		logData["user_id"] = userID
	}
	
	if err.Details != nil {
		logData["details"] = err.Details
	}
	
	// Convert to JSON for structured logging
	logJSON, _ := json.Marshal(logData)
	
	// Log based on severity
	switch err.Code {
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		log.Printf("ERROR: %s", string(logJSON))
	case http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound:
		log.Printf("WARN: %s", string(logJSON))
	default:
		log.Printf("INFO: %s", string(logJSON))
	}
}

// ValidateStruct validates a struct using the validator package
func ValidateStruct(s interface{}) error {
	validate := validator.New()
	return validate.Struct(s)
}

// ValidateField validates a single field
func ValidateField(field interface{}, tag string) error {
	validate := validator.New()
	return validate.Var(field, tag)
}

// ErrorMiddleware adds error handling middleware
func ErrorMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(error); ok {
			WriteErrorResponse(c, err)
		} else {
			WriteErrorResponse(c, NewInternalError(fmt.Sprintf("Panic: %v", recovered)))
		}
		c.Abort()
	})
}

// RequestIDMiddleware adds request ID to context
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := GenerateUUID()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}
// GetCurrentTimestamp returns current timestamp in ISO format
func GetCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

// GenerateUUID generates a new UUID string
func GenerateUUID() string {
	return uuid.New().String()
}