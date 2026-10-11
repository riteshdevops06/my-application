package handlers

import (
	"encoding/json"
	"net/http"

	"ngo-donation-tracker/middleware"
	"ngo-donation-tracker/models"
	"ngo-donation-tracker/services"
	"ngo-donation-tracker/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// NGOHandlers contains all NGO-related HTTP handlers
type NGOHandlers struct {
	ngoService *services.NGOService
}

// NewNGOHandlers creates a new instance of NGOHandlers
func NewNGOHandlers(ngoService *services.NGOService) *NGOHandlers {
	return &NGOHandlers{
		ngoService: ngoService,
	}
}

// CreateNGOProfile handles NGO profile creation
func (h *NGOHandlers) CreateNGOProfile(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req models.NGOCreate

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create NGO profile
	ngo, err := h.ngoService.CreateNGOProfile(r.Context(), claims.UserID, &req)
	if err != nil {
		// Determine appropriate status code based on error
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		} else if err.Error() == "user must have ngo_admin role to create NGO profile" ||
				  err.Error() == "user account is not active" {
			status = http.StatusForbidden
		} else if err.Error() == "NGO profile already exists for this user" {
			status = http.StatusConflict
		} else if err.Error() == "validation error" {
			status = http.StatusBadRequest
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteCreated(w, ngo)
}

// GetNGOProfile handles NGO profile retrieval
func (h *NGOHandlers) GetNGOProfile(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Get NGO profile by user ID
	ngo, err := h.ngoService.GetNGOByUserID(r.Context(), claims.UserID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "NGO profile not found for user" {
			status = http.StatusNotFound
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteSuccess(w, ngo)
}

// UpdateNGOProfile handles NGO profile updates
func (h *NGOHandlers) UpdateNGOProfile(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req models.NGOCreate

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update NGO profile
	ngo, err := h.ngoService.UpdateNGOProfile(r.Context(), claims.UserID, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "NGO profile not found for user" {
			status = http.StatusNotFound
		} else if err.Error() == "validation error" {
			status = http.StatusBadRequest
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteSuccess(w, ngo)
}

// GetNGOByID handles public NGO profile retrieval by ID
func (h *NGOHandlers) GetNGOByID(w http.ResponseWriter, r *http.Request) {
	// Extract NGO ID from URL
	vars := mux.Vars(r)
	ngoIDStr := vars["id"]
	
	ngoID, err := uuid.Parse(ngoIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid NGO ID")
		return
	}

	// Get NGO profile
	ngo, err := h.ngoService.GetNGOProfile(r.Context(), ngoID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "NGO not found" {
			status = http.StatusNotFound
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteSuccess(w, ngo)
}

// ListNGOs handles NGO listing with optional filters
func (h *NGOHandlers) ListNGOs(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	params := utils.GetPaginationParams(r)

	// Parse verification status filter
	verificationStatus := r.URL.Query().Get("status")

	// List NGOs
	ngos, total, err := h.ngoService.ListNGOs(r.Context(), params, verificationStatus)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Calculate pagination metadata
	params.Total = int(total)
	if params.Limit > 0 {
		params.TotalPages = (params.Total + params.Limit - 1) / params.Limit
	}

	response := models.PaginatedResponse{
		Data:       ngos,
		Pagination: params,
	}

	utils.WriteSuccess(w, response)
}

// GetVerifiedNGOs handles verified NGO listing
func (h *NGOHandlers) GetVerifiedNGOs(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	params := utils.GetPaginationParams(r)

	// Get verified NGOs
	ngos, total, err := h.ngoService.GetVerifiedNGOs(r.Context(), params)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Calculate pagination metadata
	params.Total = int(total)
	if params.Limit > 0 {
		params.TotalPages = (params.Total + params.Limit - 1) / params.Limit
	}

	response := models.PaginatedResponse{
		Data:       ngos,
		Pagination: params,
	}

	utils.WriteSuccess(w, response)
}
// SearchNGOs handles NGO search
func (h *NGOHandlers) SearchNGOs(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	params := utils.GetPaginationParams(r)

	// Parse search term
	searchTerm := r.URL.Query().Get("q")
	if searchTerm == "" {
		utils.WriteError(w, http.StatusBadRequest, "Search term is required")
		return
	}

	// Search NGOs
	ngos, total, err := h.ngoService.SearchNGOs(r.Context(), searchTerm, params)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Calculate pagination metadata
	params.Total = int(total)
	if params.Limit > 0 {
		params.TotalPages = (params.Total + params.Limit - 1) / params.Limit
	}

	response := models.PaginatedResponse{
		Data:       ngos,
		Pagination: params,
	}

	utils.WriteSuccess(w, response)
}

// ResubmitForVerification handles NGO verification resubmission
func (h *NGOHandlers) ResubmitForVerification(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Resubmit for verification
	err := h.ngoService.ResubmitForVerification(r.Context(), claims.UserID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "NGO profile not found for user" {
			status = http.StatusNotFound
		} else if err.Error() == "can only resubmit rejected NGOs for verification" {
			status = http.StatusBadRequest
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "NGO resubmitted for verification"})
}

// Admin Handlers

// AdminListNGOs handles admin NGO listing with filters
func (h *NGOHandlers) AdminListNGOs(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Verify admin role
	if claims.Role != "platform_admin" {
		utils.WriteError(w, http.StatusForbidden, "Admin access required")
		return
	}

	// Parse pagination parameters
	params := utils.GetPaginationParams(r)

	// Parse verification status filter
	verificationStatus := r.URL.Query().Get("status")

	// List NGOs
	ngos, total, err := h.ngoService.ListNGOs(r.Context(), params, verificationStatus)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Calculate pagination metadata
	params.Total = int(total)
	if params.Limit > 0 {
		params.TotalPages = (params.Total + params.Limit - 1) / params.Limit
	}

	response := models.PaginatedResponse{
		Data:       ngos,
		Pagination: params,
	}

	utils.WriteSuccess(w, response)
}

// AdminGetPendingVerifications handles admin pending verification listing
func (h *NGOHandlers) AdminGetPendingVerifications(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Verify admin role
	if claims.Role != "platform_admin" {
		utils.WriteError(w, http.StatusForbidden, "Admin access required")
		return
	}

	// Parse pagination parameters
	params := utils.GetPaginationParams(r)

	// Get pending verifications
	ngos, total, err := h.ngoService.GetPendingVerifications(r.Context(), params)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Calculate pagination metadata
	params.Total = int(total)
	if params.Limit > 0 {
		params.TotalPages = (params.Total + params.Limit - 1) / params.Limit
	}

	response := models.PaginatedResponse{
		Data:       ngos,
		Pagination: params,
	}

	utils.WriteSuccess(w, response)
}

// AdminVerifyNGO handles NGO verification by admin
func (h *NGOHandlers) AdminVerifyNGO(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Verify admin role
	if claims.Role != "platform_admin" {
		utils.WriteError(w, http.StatusForbidden, "Admin access required")
		return
	}

	// Extract NGO ID from URL
	vars := mux.Vars(r)
	ngoIDStr := vars["id"]
	
	ngoID, err := uuid.Parse(ngoIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid NGO ID")
		return
	}

	var req struct {
		Notes string `json:"notes,omitempty"`
	}

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Verify NGO
	err = h.ngoService.VerifyNGO(r.Context(), ngoID, claims.UserID, req.Notes)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "NGO not found" {
			status = http.StatusNotFound
		} else if err.Error() == "only platform admins can verify NGOs" {
			status = http.StatusForbidden
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "NGO verified successfully"})
}

// AdminRejectNGO handles NGO rejection by admin
func (h *NGOHandlers) AdminRejectNGO(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Verify admin role
	if claims.Role != "platform_admin" {
		utils.WriteError(w, http.StatusForbidden, "Admin access required")
		return
	}

	// Extract NGO ID from URL
	vars := mux.Vars(r)
	ngoIDStr := vars["id"]
	
	ngoID, err := uuid.Parse(ngoIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid NGO ID")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Reject NGO
	err = h.ngoService.RejectNGO(r.Context(), ngoID, claims.UserID, req.Reason)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "NGO not found" {
			status = http.StatusNotFound
		} else if err.Error() == "only platform admins can reject NGOs" {
			status = http.StatusForbidden
		} else if err.Error() == "rejection reason is required" {
			status = http.StatusBadRequest
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "NGO rejected successfully"})
}

// AdminGetNGOStats handles NGO statistics for admin dashboard
func (h *NGOHandlers) AdminGetNGOStats(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Verify admin role
	if claims.Role != "platform_admin" {
		utils.WriteError(w, http.StatusForbidden, "Admin access required")
		return
	}

	// Get NGO statistics
	stats, err := h.ngoService.GetNGOStats(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(w, stats)
}