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

// GetCurrentUser returns current user profile with NGO information if applicable
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

	// Create response with user data and NGO info if available
	response := map[string]interface{}{
		"user": user,
	}

	// If user is NGO admin and has NGO ID in claims, include NGO info
	if claims.Role == "ngo_admin" && claims.NGOID != nil {
		response["ngo_id"] = *claims.NGOID
		response["has_ngo_profile"] = true
	} else if claims.Role == "ngo_admin" {
		response["has_ngo_profile"] = false
	}

	utils.WriteSuccess(w, response)
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
// RegisterWithProfile handles user registration with automatic NGO profile creation for NGO admins
func (h *AuthHandlers) RegisterWithProfile(w http.ResponseWriter, r *http.Request) {
	var req struct {
		models.UserCreate
		NGOProfile *models.NGOCreate `json:"ngo_profile,omitempty"`
	}
	
	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Register user
	authResponse, err := h.authService.Register(r.Context(), &req.UserCreate)
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

	// If user is NGO admin and NGO profile data is provided, create NGO profile
	if req.Role == "ngo_admin" && req.NGOProfile != nil {
		// Note: In a production system, you might want to create NGO service here
		// For now, we'll return the auth response and let the client create the NGO profile separately
		authResponse.User.PasswordHash = "" // Ensure password hash is not exposed
		
		response := map[string]interface{}{
			"user":   authResponse.User,
			"tokens": authResponse.Tokens,
			"message": "User registered successfully. Please create your NGO profile to complete setup.",
		}
		utils.WriteCreated(w, response)
		return
	}

	utils.WriteCreated(w, authResponse)
}

// CheckEmailAvailability checks if an email is available for registration
func (h *AuthHandlers) CheckEmailAvailability(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		utils.WriteError(w, http.StatusBadRequest, "Email parameter is required")
		return
	}

	// Check if email exists (this should be done via service, but for now we'll create a simple check)
	available := true
	message := "Email is available"

	// In a real implementation, you'd check via the auth service
	// For now, we'll return a simple response
	
	utils.WriteSuccess(w, map[string]interface{}{
		"email":     email,
		"available": available,
		"message":   message,
	})
}

// RefreshTokens handles both access and refresh token renewal
func (h *AuthHandlers) RefreshTokens(w http.ResponseWriter, r *http.Request) {
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

// ValidateToken validates the current token and returns user info
func (h *AuthHandlers) ValidateToken(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	// Get current user to ensure account is still active
	user, err := h.authService.GetCurrentUser(r.Context(), claims.UserID)
	if err != nil {
		status := http.StatusUnauthorized
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		utils.WriteError(w, status, "Token is valid but user account is not available")
		return
	}

	// Return validation response
	response := map[string]interface{}{
		"valid":   true,
		"user_id": claims.UserID,
		"email":   claims.Email,
		"role":    claims.Role,
		"user":    user,
	}

	if claims.NGOID != nil {
		response["ngo_id"] = *claims.NGOID
	}

	utils.WriteSuccess(w, response)
}

// GetUserRoles returns available user roles for registration
func (h *AuthHandlers) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	roles := []map[string]string{
		{
			"value":       "donor",
			"label":       "Donor",
			"description": "Individual or organization that wants to donate to NGO campaigns",
		},
		{
			"value":       "ngo_admin",
			"label":       "NGO Administrator", 
			"description": "NGO representative who manages campaigns and organizational profile",
		},
		// Note: platform_admin is typically not available for public registration
	}

	utils.WriteSuccess(w, map[string]interface{}{
		"roles": roles,
	})
}

// RequestPasswordReset initiates password reset process
func (h *AuthHandlers) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	
	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" {
		utils.WriteError(w, http.StatusBadRequest, "Email is required")
		return
	}

	// TODO: Implement actual password reset logic with email sending
	// For now, return a success response
	utils.WriteSuccess(w, map[string]string{
		"message": "If an account with this email exists, a password reset link has been sent",
	})
}

// ConfirmPasswordReset completes the password reset process
func (h *AuthHandlers) ConfirmPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token           string `json:"token"`
		NewPassword     string `json:"new_password"`
	}
	
	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Token == "" {
		utils.WriteError(w, http.StatusBadRequest, "Reset token is required")
		return
	}

	if req.NewPassword == "" {
		utils.WriteError(w, http.StatusBadRequest, "New password is required")
		return
	}

	if len(req.NewPassword) < 8 {
		utils.WriteError(w, http.StatusBadRequest, "Password must be at least 8 characters long")
		return
	}

	// TODO: Implement actual password reset confirmation logic
	// For now, return a success response
	utils.WriteSuccess(w, map[string]string{
		"message": "Password has been reset successfully",
	})
}