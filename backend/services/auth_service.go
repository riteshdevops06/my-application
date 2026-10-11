package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"ngo-donation-tracker/config"
	"ngo-donation-tracker/models"
	"ngo-donation-tracker/repositories"
	"ngo-donation-tracker/utils"

	"github.com/google/uuid"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo repositories.UserRepositoryInterface
	config   *config.Config
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo repositories.UserRepositoryInterface, config *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   config,
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *models.UserCreate) (*models.AuthResponse, error) {
	// Validate input
	if err := s.validateRegistrationInput(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user, err := s.userRepo.CreateUser(ctx, req, passwordHash)
	if err != nil {
		if errors.Is(err, repositories.ErrUserExists) {
			return nil, fmt.Errorf("account with this email already exists")
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Remove password hash from response
	user.PasswordHash = ""

	return &models.AuthResponse{
		User:   *user,
		Tokens: *tokens,
	}, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, req *models.UserLogin) (*models.AuthResponse, error) {
	// Validate input
	if err := s.validateLoginInput(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Get user for authentication
	user, err := s.userRepo.GetUserForAuth(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) || errors.Is(err, repositories.ErrInvalidCredentials) {
			return nil, fmt.Errorf("invalid email or password")
		}
		if errors.Is(err, repositories.ErrUserInactive) {
			return nil, fmt.Errorf("account has been deactivated")
		}
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate tokens
	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Remove password hash from response
	user.PasswordHash = ""

	return &models.AuthResponse{
		User:   *user,
		Tokens: *tokens,
	}, nil
}

// RefreshToken generates new access token from refresh token
func (s *AuthService) RefreshToken(ctx context.Context, userID uuid.UUID) (*models.TokenPair, error) {
	// Get user to ensure they still exist and are active
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Generate new tokens
	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return tokens, nil
}

// GetCurrentUser returns the current user's profile
func (s *AuthService) GetCurrentUser(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Remove password hash from response
	user.PasswordHash = ""
	return user, nil
}

// UpdateProfile updates user profile information
func (s *AuthService) UpdateProfile(ctx context.Context, userID uuid.UUID, req *models.UserUpdate) (*models.User, error) {
	// Validate input
	if err := s.validateProfileUpdate(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Update user
	user, err := s.userRepo.UpdateUser(ctx, userID, req)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// Remove password hash from response
	user.PasswordHash = ""
	return user, nil
}

// ChangePassword changes user's password
func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, req *models.PasswordChange) error {
	// Validate input
	if err := s.validatePasswordChange(req); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Get current user to verify current password
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify current password
	if !utils.CheckPasswordHash(req.CurrentPassword, user.PasswordHash) {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	newPasswordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	err = s.userRepo.ChangePassword(ctx, userID, newPasswordHash)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to change password: %w", err)
	}

	return nil
}

// VerifyEmail marks user's email as verified
func (s *AuthService) VerifyEmail(ctx context.Context, userID uuid.UUID) error {
	err := s.userRepo.VerifyEmail(ctx, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to verify email: %w", err)
	}
	return nil
}

// DeactivateAccount deactivates a user account
func (s *AuthService) DeactivateAccount(ctx context.Context, userID uuid.UUID) error {
	err := s.userRepo.DeactivateUser(ctx, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to deactivate account: %w", err)
	}
	return nil
}

// Helper methods

// generateTokens generates access and refresh tokens for a user
func (s *AuthService) generateTokens(user *models.User) (*models.TokenPair, error) {
	// Get NGO ID if user is NGO admin
	var ngoID *uuid.UUID
	if user.Role == "ngo_admin" {
		// Note: In a real implementation, you'd fetch the NGO ID from the database
		// For now, we'll leave it nil and handle it when we implement NGO repository
	}

	// Generate access token
	accessToken, err := utils.GenerateToken(
		user.ID,
		user.Email,
		user.Role,
		ngoID,
		s.config.JWTSecret,
		s.config.JWTExpiry,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &models.TokenPair{
		AccessToken: accessToken,
		ExpiresIn:   s.config.JWTExpiry * 3600, // Convert hours to seconds
		TokenType:   "Bearer",
	}, nil
}

// Validation methods

func (s *AuthService) validateRegistrationInput(req *models.UserCreate) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if req.Name == "" {
		return errors.New("name is required")
	}
	if len(strings.TrimSpace(req.Name)) < 2 {
		return errors.New("name must be at least 2 characters long")
	}
	if req.Role == "" {
		return errors.New("role is required")
	}
	if req.Role != "donor" && req.Role != "ngo_admin" && req.Role != "platform_admin" {
		return errors.New("invalid role")
	}

	// Basic email validation
	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		return errors.New("invalid email format")
	}

	return nil
}

func (s *AuthService) validateLoginInput(req *models.UserLogin) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (s *AuthService) validateProfileUpdate(req *models.UserUpdate) error {
	if req.Name != "" && len(strings.TrimSpace(req.Name)) < 2 {
		return errors.New("name must be at least 2 characters long")
	}
	return nil
}

func (s *AuthService) validatePasswordChange(req *models.PasswordChange) error {
	if req.CurrentPassword == "" {
		return errors.New("current password is required")
	}
	if req.NewPassword == "" {
		return errors.New("new password is required")
	}
	if len(req.NewPassword) < 8 {
		return errors.New("new password must be at least 8 characters long")
	}
	if req.CurrentPassword == req.NewPassword {
		return errors.New("new password must be different from current password")
	}
	return nil
}