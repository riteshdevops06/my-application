package handlers

import (
	"encoding/json"
	"net/http"

	"ngo-donation-tracker/middleware"
	"ngo-donation-tracker/models"
	"ngo-donation-tracker/services"
	"ngo-donation-tracker/utils"

	"github.com/google/uuid"
)

// AuthHandlers contains all authentication-related HTTP handlers
type AuthHandlers struct {
	authService *services.AuthService
}

// NewAuthHandlers creates a new instance of AuthHandlers
func NewAuthHandlers(authService *services.AuthService) *AuthHandlers {
	return &AuthHandlers{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	var req models.UserCreate
	
	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Register user
	response, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		// Determine appropriate status code based on error
		status := http.StatusInternalServerError
		if err.Error() == "account with this email already exists" {
			status = http.StatusConflict
		} else if err.Error() == "validation error" || 
		         err.Error() == "invalid role" ||
				 err.Error() == "invalid email format" {
			status = http.StatusBadRequest
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteCreated(w, response)
}

// Login handles user authentication
func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var req models.UserLogin
	
	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Authenticate user
	response, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		// Determine appropriate status code based on error
		status := http.StatusInternalServerError
		if err.Error() == "invalid email or password" {
			status = http.StatusUnauthorized
		} else if err.Error() == "account has been deactivated" {
			status = http.StatusForbidden
		} else if err.Error() == "validation error" {
			status = http.StatusBadRequest
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteSuccess(w, response)
}

// RefreshToken handles token refresh
func (h *AuthHandlers) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Generate new tokens
	tokens, err := h.authService.RefreshToken(r.Context(), claims.UserID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		} else if err.Error() == "account is deactivated" {
			status = http.StatusForbidden
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteSuccess(w, tokens)
}

// GetCurrentUser returns current user profile
func (h *AuthHandlers) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Get user profile
	user, err := h.authService.GetCurrentUser(r.Context(), claims.UserID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteSuccess(w, user)
}

// UpdateProfile updates user profile
func (h *AuthHandlers) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req models.UserUpdate
	
	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update profile
	user, err := h.authService.UpdateProfile(r.Context(), claims.UserID, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		} else if err.Error() == "validation error" {
			status = http.StatusBadRequest
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteSuccess(w, user)
}

// ChangePassword handles password change
func (h *AuthHandlers) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req models.PasswordChange
	
	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Change password
	err := h.authService.ChangePassword(r.Context(), claims.UserID, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		} else if err.Error() == "current password is incorrect" {
			status = http.StatusUnauthorized
		} else if err.Error() == "validation error" {
			status = http.StatusBadRequest
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "Password changed successfully"})
}

// VerifyEmail handles email verification
func (h *AuthHandlers) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Verify email
	err := h.authService.VerifyEmail(r.Context(), claims.UserID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "Email verified successfully"})
}

// DeactivateAccount handles account deactivation
func (h *AuthHandlers) DeactivateAccount(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Deactivate account
	err := h.authService.DeactivateAccount(r.Context(), claims.UserID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "Account deactivated successfully"})
}

// ForgotPassword handles password reset request
func (h *AuthHandlers) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	
	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Implement password reset logic with email sending
	// For now, return a placeholder response
	utils.WriteSuccess(w, map[string]string{
		"message": "Password reset email sent (not implemented yet)",
	})
}

// ResetPassword handles password reset
func (h *AuthHandlers) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	
	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Implement password reset logic with token validation
	// For now, return a placeholder response
	utils.WriteSuccess(w, map[string]string{
		"message": "Password reset successful (not implemented yet)",
	})
}

// Logout handles user logout (client-side token removal)
func (h *AuthHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	// JWT tokens are stateless, so logout is handled client-side by removing the token
	// In a production system, you might want to implement token blacklisting
	utils.WriteSuccess(w, map[string]string{
		"message": "Logged out successfully",
	})
}

// Helper function to extract UUID from URL path
func extractUUIDFromPath(r *http.Request, paramName string) (uuid.UUID, error) {
	// This would typically use a router like gorilla/mux to extract path parameters
	// For now, we'll implement this when we update the routes
	return uuid.Nil, nil
}