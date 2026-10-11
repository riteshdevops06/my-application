package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"ngo-donation-tracker/models"
	"ngo-donation-tracker/repositories"

	"github.com/google/uuid"
)

// AuthUtils provides utility functions for authentication
type AuthUtils struct {
	userRepo repositories.UserRepositoryInterface
}

// NewAuthUtils creates a new authentication utilities instance
func NewAuthUtils(userRepo repositories.UserRepositoryInterface) *AuthUtils {
	return &AuthUtils{
		userRepo: userRepo,
	}
}

// ValidateEmail performs comprehensive email validation
func (au *AuthUtils) ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	// Basic format validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	// Length validation
	if len(email) > 254 {
		return errors.New("email address is too long")
	}

	// Domain validation (basic)
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errors.New("invalid email format")
	}

	domain := parts[1]
	if len(domain) > 63 {
		return errors.New("domain name is too long")
	}

	return nil
}

// ValidatePassword performs comprehensive password validation
func (au *AuthUtils) ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password is required")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if len(password) > 128 {
		return errors.New("password is too long (maximum 128 characters)")
	}

	// Check for at least one letter
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	if !hasLetter {
		return errors.New("password must contain at least one letter")
	}

	// Check for at least one number
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}

	// Check for common passwords (basic check)
	commonPasswords := []string{"password", "12345678", "password123", "admin123"}
	lowerPassword := strings.ToLower(password)
	for _, common := range commonPasswords {
		if lowerPassword == common {
			return errors.New("password is too common, please choose a stronger password")
		}
	}

	return nil
}

// ValidateUserRole validates if a role is allowed for registration
func (au *AuthUtils) ValidateUserRole(role string) error {
	allowedRoles := []string{"donor", "ngo_admin"}
	
	for _, allowedRole := range allowedRoles {
		if role == allowedRole {
			return nil
		}
	}

	return errors.New("invalid role for registration")
}

// ValidateName validates user name
func (au *AuthUtils) ValidateName(name string) error {
	if name == "" {
		return errors.New("name is required")
	}

	trimmed := strings.TrimSpace(name)
	if len(trimmed) < 2 {
		return errors.New("name must be at least 2 characters long")
	}

	if len(trimmed) > 100 {
		return errors.New("name is too long (maximum 100 characters)")
	}

	// Check for valid characters (letters, spaces, hyphens, apostrophes)
	nameRegex := regexp.MustCompile(`^[a-zA-Z\s\-'\.]+$`)
	if !nameRegex.MatchString(trimmed) {
		return errors.New("name contains invalid characters")
	}

	return nil
}

// CheckEmailAvailability checks if an email is available for registration
func (au *AuthUtils) CheckEmailAvailability(ctx context.Context, email string) (bool, error) {
	// Validate email first
	if err := au.ValidateEmail(email); err != nil {
		return false, err
	}

	// Check if email exists
	exists, err := au.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to check email availability: %w", err)
	}

	return !exists, nil
}

// ValidateRegistrationData performs comprehensive validation of registration data
func (au *AuthUtils) ValidateRegistrationData(ctx context.Context, req *models.UserCreate) error {
	// Validate email
	if err := au.ValidateEmail(req.Email); err != nil {
		return fmt.Errorf("email validation failed: %w", err)
	}

	// Check email availability
	available, err := au.CheckEmailAvailability(ctx, req.Email)
	if err != nil {
		return err
	}
	if !available {
		return errors.New("email address is already registered")
	}

	// Validate password
	if err := au.ValidatePassword(req.Password); err != nil {
		return fmt.Errorf("password validation failed: %w", err)
	}

	// Validate name
	if err := au.ValidateName(req.Name); err != nil {
		return fmt.Errorf("name validation failed: %w", err)
	}

	// Validate role
	if err := au.ValidateUserRole(req.Role); err != nil {
		return fmt.Errorf("role validation failed: %w", err)
	}

	// Validate phone if provided
	if req.Phone != "" {
		if err := au.ValidatePhone(req.Phone); err != nil {
			return fmt.Errorf("phone validation failed: %w", err)
		}
	}

	return nil
}

// ValidatePhone validates phone number format
func (au *AuthUtils) ValidatePhone(phone string) error {
	if phone == "" {
		return nil // Phone is optional
	}

	// Remove spaces, dashes, and parentheses for validation
	cleaned := regexp.MustCompile(`[\s\-\(\)]+`).ReplaceAllString(phone, "")
	
	// Check if it starts with + for international format
	if strings.HasPrefix(cleaned, "+") {
		cleaned = cleaned[1:]
	}

	// Should contain only digits after cleaning
	if !regexp.MustCompile(`^\d+$`).MatchString(cleaned) {
		return errors.New("phone number should contain only digits, spaces, dashes, parentheses, and optionally start with +")
	}

	// Length validation
	if len(cleaned) < 7 || len(cleaned) > 15 {
		return errors.New("phone number should be between 7 and 15 digits")
	}

	return nil
}

// GeneratePasswordResetToken generates a secure token for password reset
func (au *AuthUtils) GeneratePasswordResetToken() (string, error) {
	// In a production system, you would generate a cryptographically secure token
	// For now, we'll use a UUID as a placeholder
	token := uuid.New()
	return token.String(), nil
}

// ValidatePasswordResetToken validates a password reset token
func (au *AuthUtils) ValidatePasswordResetToken(token string) (bool, uuid.UUID, error) {
	// In a production system, you would:
	// 1. Parse the token
	// 2. Check if it exists in the database
	// 3. Check if it's not expired
	// 4. Return the associated user ID
	
	// For now, we'll just validate the UUID format
	userID, err := uuid.Parse(token)
	if err != nil {
		return false, uuid.Nil, errors.New("invalid token format")
	}

	// TODO: Add actual token validation logic with database lookup
	
	return true, userID, nil
}

// AuthResponse represents an enhanced authentication response
type AuthResponse struct {
	User         *models.User      `json:"user"`
	Tokens       *models.TokenPair `json:"tokens"`
	NGOProfile   *models.NGO       `json:"ngo_profile,omitempty"`
	FirstLogin   bool              `json:"first_login"`
	RequiresSetup bool             `json:"requires_setup"`
	Message      string            `json:"message,omitempty"`
}

// LoginAttemptTracker tracks failed login attempts
type LoginAttemptTracker struct {
	attempts    map[string][]time.Time
	maxAttempts int
	lockoutDuration time.Duration
}

// NewLoginAttemptTracker creates a new login attempt tracker
func NewLoginAttemptTracker(maxAttempts int, lockoutDuration time.Duration) *LoginAttemptTracker {
	return &LoginAttemptTracker{
		attempts:        make(map[string][]time.Time),
		maxAttempts:     maxAttempts,
		lockoutDuration: lockoutDuration,
	}
}

// IsBlocked checks if an email is currently blocked due to too many failed attempts
func (lat *LoginAttemptTracker) IsBlocked(email string) bool {
	attempts, exists := lat.attempts[email]
	if !exists || len(attempts) < lat.maxAttempts {
		return false
	}

	// Check if the oldest attempt is within the lockout duration
	oldestAttempt := attempts[0]
	return time.Since(oldestAttempt) < lat.lockoutDuration
}

// RecordFailedAttempt records a failed login attempt
func (lat *LoginAttemptTracker) RecordFailedAttempt(email string) {
	now := time.Now()
	attempts := lat.attempts[email]
	
	// Clean up old attempts outside the lockout window
	var recentAttempts []time.Time
	for _, attempt := range attempts {
		if now.Sub(attempt) < lat.lockoutDuration {
			recentAttempts = append(recentAttempts, attempt)
		}
	}
	
	// Add the new attempt
	recentAttempts = append(recentAttempts, now)
	lat.attempts[email] = recentAttempts
}

// ClearAttempts clears failed attempts for an email (called on successful login)
func (lat *LoginAttemptTracker) ClearAttempts(email string) {
	delete(lat.attempts, email)
}

// GetRemainingAttempts returns the number of remaining attempts before lockout
func (lat *LoginAttemptTracker) GetRemainingAttempts(email string) int {
	attempts := lat.attempts[email]
	if attempts == nil {
		return lat.maxAttempts
	}
	
	remaining := lat.maxAttempts - len(attempts)
	if remaining < 0 {
		return 0
	}
	return remaining
}