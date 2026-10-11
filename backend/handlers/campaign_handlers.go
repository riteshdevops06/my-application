package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"ngo-donation-tracker/middleware"
	"ngo-donation-tracker/models"
	"ngo-donation-tracker/repositories"
	"ngo-donation-tracker/services"
	"ngo-donation-tracker/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// CampaignHandlers contains all campaign-related HTTP handlers
type CampaignHandlers struct {
	campaignService *services.CampaignService
}

// NewCampaignHandlers creates a new instance of CampaignHandlers
func NewCampaignHandlers(campaignService *services.CampaignService) *CampaignHandlers {
	return &CampaignHandlers{
		campaignService: campaignService,
	}
}

// Public Campaign Handlers (no authentication required)

// ListPublicCampaigns handles public campaign listing
func (h *CampaignHandlers) ListPublicCampaigns(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	params := utils.GetPaginationParams(r)

	// Parse filter parameters
	filters := repositories.CampaignFilters{
		Status:     "active", // Only show active campaigns to public
		Category:   r.URL.Query().Get("category"),
		SearchTerm: r.URL.Query().Get("search"),
	}

	// Parse featured filter
	if featured := r.URL.Query().Get("featured"); featured == "true" {
		isFeatured := true
		filters.IsFeatured = &isFeatured
	}

	// List campaigns
	campaigns, total, err := h.campaignService.ListCampaigns(r.Context(), params, filters)
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
		Data:       campaigns,
		Pagination: params,
	}

	utils.WriteSuccess(w, response)
}

// GetPublicCampaign handles public campaign retrieval by ID
func (h *CampaignHandlers) GetPublicCampaign(w http.ResponseWriter, r *http.Request) {
	// Extract campaign ID from URL
	vars := mux.Vars(r)
	campaignIDStr := vars["id"]
	
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid campaign ID")
		return
	}

	// Get campaign with NGO information
	campaignData, err := h.campaignService.GetCampaignWithNGO(r.Context(), campaignID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "campaign not found" {
			status = http.StatusNotFound
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteSuccess(w, campaignData)
}

// GetFeaturedCampaigns handles featured campaign listing
func (h *CampaignHandlers) GetFeaturedCampaigns(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	params := utils.GetPaginationParams(r)

	// Get featured campaigns
	campaigns, total, err := h.campaignService.GetFeaturedCampaigns(r.Context(), params)
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
		Data:       campaigns,
		Pagination: params,
	}

	utils.WriteSuccess(w, response)
}

// SearchCampaigns handles campaign search
func (h *CampaignHandlers) SearchCampaigns(w http.ResponseWriter, r *http.Request) {
	// Parse search term
	searchTerm := r.URL.Query().Get("q")
	if searchTerm == "" {
		utils.WriteError(w, http.StatusBadRequest, "Search term is required")
		return
	}

	// Parse pagination parameters
	params := utils.GetPaginationParams(r)

	// Search campaigns
	campaigns, total, err := h.campaignService.SearchCampaigns(r.Context(), searchTerm, params)
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
		Data:       campaigns,
		Pagination: params,
	}

	utils.WriteSuccess(w, response)
}

// GetCampaignsByCategory handles category-based campaign listing
func (h *CampaignHandlers) GetCampaignsByCategory(w http.ResponseWriter, r *http.Request) {
	// Extract category from URL
	vars := mux.Vars(r)
	category := vars["category"]

	// Parse pagination parameters
	params := utils.GetPaginationParams(r)

	// Get campaigns by category
	campaigns, total, err := h.campaignService.GetCampaignsByCategory(r.Context(), category, params)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "category is required" {
			status = http.StatusBadRequest
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	// Calculate pagination metadata
	params.Total = int(total)
	if params.Limit > 0 {
		params.TotalPages = (params.Total + params.Limit - 1) / params.Limit
	}

	response := models.PaginatedResponse{
		Data:       campaigns,
		Pagination: params,
	}

	utils.WriteSuccess(w, response)
}

// NGO Campaign Management Handlers (requires ngo_admin role)

// CreateCampaign handles campaign creation by NGO admins
func (h *CampaignHandlers) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Verify NGO admin role
	if claims.Role != "ngo_admin" {
		utils.WriteError(w, http.StatusForbidden, "NGO admin access required")
		return
	}

	var req models.CampaignCreate

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create campaign
	campaign, err := h.campaignService.CreateCampaign(r.Context(), claims.UserID, &req)
	if err != nil {
		// Determine appropriate status code based on error
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		} else if err.Error() == "only NGO administrators can create campaigns" ||
				  err.Error() == "user account is not active" ||
				  err.Error() == "NGO must be verified before creating campaigns" {
			status = http.StatusForbidden
		} else if err.Error() == "NGO profile not found. Please create your NGO profile first" {
			status = http.StatusPreconditionFailed
		} else if err.Error() == "campaign with this title already exists for this NGO" {
			status = http.StatusConflict
		} else if err.Error() == "validation error" {
			status = http.StatusBadRequest
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteCreated(w, campaign)
}

// GetNGOCampaigns handles listing campaigns for the authenticated NGO
func (h *CampaignHandlers) GetNGOCampaigns(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Verify NGO admin role
	if claims.Role != "ngo_admin" {
		utils.WriteError(w, http.StatusForbidden, "NGO admin access required")
		return
	}

	// Parse pagination parameters
	params := utils.GetPaginationParams(r)

	// List NGO campaigns
	campaigns, total, err := h.campaignService.ListNGOCampaigns(r.Context(), claims.UserID, params)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "NGO profile not found" {
			status = http.StatusNotFound
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	// Calculate pagination metadata
	params.Total = int(total)
	if params.Limit > 0 {
		params.TotalPages = (params.Total + params.Limit - 1) / params.Limit
	}

	response := models.PaginatedResponse{
		Data:       campaigns,
		Pagination: params,
	}

	utils.WriteSuccess(w, response)
}

// GetCampaign handles campaign retrieval by ID (with ownership check for NGO admins)
func (h *CampaignHandlers) GetCampaign(w http.ResponseWriter, r *http.Request) {
	// Extract campaign ID from URL
	vars := mux.Vars(r)
	campaignIDStr := vars["id"]
	
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid campaign ID")
		return
	}

	// Get campaign with NGO information
	campaignData, err := h.campaignService.GetCampaignWithNGO(r.Context(), campaignID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "campaign not found" {
			status = http.StatusNotFound
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	// If user is authenticated and is NGO admin, they can see their own campaigns regardless of status
	claims := middleware.GetUserFromContext(r)
	if claims != nil && claims.Role == "ngo_admin" {
		// NGO admins can see their own campaigns
		utils.WriteSuccess(w, campaignData)
		return
	}

	// For public access, only show active campaigns
	campaign := campaignData["campaign"].(*models.Campaign)
	if campaign.Status != "active" {
		utils.WriteError(w, http.StatusNotFound, "Campaign not found")
		return
	}

	utils.WriteSuccess(w, campaignData)
}
// UpdateCampaign handles campaign updates by NGO admins
func (h *CampaignHandlers) UpdateCampaign(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Verify NGO admin role
	if claims.Role != "ngo_admin" {
		utils.WriteError(w, http.StatusForbidden, "NGO admin access required")
		return
	}

	// Extract campaign ID from URL
	vars := mux.Vars(r)
	campaignIDStr := vars["id"]
	
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid campaign ID")
		return
	}

	var req models.CampaignUpdate

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update campaign
	campaign, err := h.campaignService.UpdateCampaign(r.Context(), claims.UserID, campaignID, &req)
	if err != nil {
		// Determine appropriate status code based on error
		status := http.StatusInternalServerError
		if err.Error() == "campaign not found" {
			status = http.StatusNotFound
		} else if err.Error() == "campaign does not belong to your NGO" ||
				  err.Error() == "only NGO administrators can manage campaigns" {
			status = http.StatusForbidden
		} else if err.Error() == "cannot modify completed or cancelled campaigns" ||
				  err.Error() == "campaign with this title already exists for this NGO" {
			status = http.StatusConflict
		} else if err.Error() == "validation error" {
			status = http.StatusBadRequest
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteSuccess(w, campaign)
}

// DeleteCampaign handles campaign deletion by NGO admins
func (h *CampaignHandlers) DeleteCampaign(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Verify NGO admin role
	if claims.Role != "ngo_admin" {
		utils.WriteError(w, http.StatusForbidden, "NGO admin access required")
		return
	}

	// Extract campaign ID from URL
	vars := mux.Vars(r)
	campaignIDStr := vars["id"]
	
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid campaign ID")
		return
	}

	// Delete campaign
	err = h.campaignService.DeleteCampaign(r.Context(), claims.UserID, campaignID)
	if err != nil {
		// Determine appropriate status code based on error
		status := http.StatusInternalServerError
		if err.Error() == "campaign not found" {
			status = http.StatusNotFound
		} else if err.Error() == "campaign does not belong to your NGO" ||
				  err.Error() == "only NGO administrators can manage campaigns" {
			status = http.StatusForbidden
		} else if err.Error() == "cannot delete campaigns that have received donations" {
			status = http.StatusConflict
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "Campaign deleted successfully"})
}

// PublishCampaign handles campaign publishing (draft to active)
func (h *CampaignHandlers) PublishCampaign(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Verify NGO admin role
	if claims.Role != "ngo_admin" {
		utils.WriteError(w, http.StatusForbidden, "NGO admin access required")
		return
	}

	// Extract campaign ID from URL
	vars := mux.Vars(r)
	campaignIDStr := vars["id"]
	
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid campaign ID")
		return
	}

	// Publish campaign
	err = h.campaignService.PublishCampaign(r.Context(), claims.UserID, campaignID)
	if err != nil {
		// Determine appropriate status code based on error
		status := http.StatusInternalServerError
		if err.Error() == "campaign not found" {
			status = http.StatusNotFound
		} else if err.Error() == "campaign does not belong to your NGO" ||
				  err.Error() == "only NGO administrators can manage campaigns" {
			status = http.StatusForbidden
		} else if err.Error() == "campaign cannot be published" {
			status = http.StatusBadRequest
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "Campaign published successfully"})
}

// PauseCampaign handles campaign pausing
func (h *CampaignHandlers) PauseCampaign(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Verify NGO admin role
	if claims.Role != "ngo_admin" {
		utils.WriteError(w, http.StatusForbidden, "NGO admin access required")
		return
	}

	// Extract campaign ID from URL
	vars := mux.Vars(r)
	campaignIDStr := vars["id"]
	
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid campaign ID")
		return
	}

	// Pause campaign
	err = h.campaignService.PauseCampaign(r.Context(), claims.UserID, campaignID)
	if err != nil {
		// Determine appropriate status code based on error
		status := http.StatusInternalServerError
		if err.Error() == "campaign not found" {
			status = http.StatusNotFound
		} else if err.Error() == "campaign does not belong to your NGO" ||
				  err.Error() == "only NGO administrators can manage campaigns" {
			status = http.StatusForbidden
		} else if err.Error() == "only active campaigns can be paused" {
			status = http.StatusBadRequest
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "Campaign paused successfully"})
}

// ResumeCampaign handles campaign resuming
func (h *CampaignHandlers) ResumeCampaign(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Verify NGO admin role
	if claims.Role != "ngo_admin" {
		utils.WriteError(w, http.StatusForbidden, "NGO admin access required")
		return
	}

	// Extract campaign ID from URL
	vars := mux.Vars(r)
	campaignIDStr := vars["id"]
	
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid campaign ID")
		return
	}

	// Resume campaign
	err = h.campaignService.ResumeCampaign(r.Context(), claims.UserID, campaignID)
	if err != nil {
		// Determine appropriate status code based in error
		status := http.StatusInternalServerError
		if err.Error() == "campaign not found" {
			status = http.StatusNotFound
		} else if err.Error() == "campaign does not belong to your NGO" ||
				  err.Error() == "only NGO administrators can manage campaigns" {
			status = http.StatusForbidden
		} else if err.Error() == "only paused campaigns can be resumed" {
			status = http.StatusBadRequest
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "Campaign resumed successfully"})
}

// Platform Admin Handlers (requires platform_admin role)

// AdminListCampaigns handles admin campaign listing with filters
func (h *CampaignHandlers) AdminListCampaigns(w http.ResponseWriter, r *http.Request) {
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

	// Parse filter parameters
	filters := repositories.CampaignFilters{
		Status:     r.URL.Query().Get("status"),
		Category:   r.URL.Query().Get("category"),
		SearchTerm: r.URL.Query().Get("search"),
	}

	// Parse NGO ID filter if provided
	if ngoIDStr := r.URL.Query().Get("ngo_id"); ngoIDStr != "" {
		if ngoID, err := uuid.Parse(ngoIDStr); err == nil {
			filters.NGOID = &ngoID
		}
	}

	// List campaigns
	campaigns, total, err := h.campaignService.ListCampaigns(r.Context(), params, filters)
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
		Data:       campaigns,
		Pagination: params,
	}

	utils.WriteSuccess(w, response)
}

// AdminSetFeaturedStatus handles setting/unsetting campaign as featured
func (h *CampaignHandlers) AdminSetFeaturedStatus(w http.ResponseWriter, r *http.Request) {
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

	// Extract campaign ID from URL
	vars := mux.Vars(r)
	campaignIDStr := vars["id"]
	
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid campaign ID")
		return
	}

	var req struct {
		IsFeatured bool `json:"is_featured"`
	}

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Set featured status
	err = h.campaignService.SetFeaturedStatus(r.Context(), claims.UserID, campaignID, req.IsFeatured)
	if err != nil {
		// Determine appropriate status code based on error
		status := http.StatusInternalServerError
		if err.Error() == "campaign not found" {
			status = http.StatusNotFound
		} else if err.Error() == "only platform administrators can set featured status" {
			status = http.StatusForbidden
		} else if err.Error() == "only active campaigns can be featured" {
			status = http.StatusBadRequest
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	message := "Campaign unfeatured successfully"
	if req.IsFeatured {
		message = "Campaign featured successfully"
	}

	utils.WriteSuccess(w, map[string]string{"message": message})
}

// AdminSuspendCampaign handles campaign suspension by admin
func (h *CampaignHandlers) AdminSuspendCampaign(w http.ResponseWriter, r *http.Request) {
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

	// Extract campaign ID from URL
	vars := mux.Vars(r)
	campaignIDStr := vars["id"]
	
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid campaign ID")
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

	// Suspend campaign
	err = h.campaignService.AdminSuspendCampaign(r.Context(), claims.UserID, campaignID, req.Reason)
	if err != nil {
		// Determine appropriate status code based on error
		status := http.StatusInternalServerError
		if err.Error() == "campaign not found" {
			status = http.StatusNotFound
		} else if err.Error() == "only platform administrators can suspend campaigns" {
			status = http.StatusForbidden
		} else if err.Error() == "suspension reason is required" {
			status = http.StatusBadRequest
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "Campaign suspended successfully"})
}
// Utility and Configuration Handlers

// GetCampaignStats handles campaign statistics for admin dashboard
func (h *CampaignHandlers) GetCampaignStats(w http.ResponseWriter, r *http.Request) {
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

	// Get campaign statistics
	stats, err := h.campaignService.GetCampaignStats(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(w, stats)
}

// GetCampaignCategories handles available campaign categories
func (h *CampaignHandlers) GetCampaignCategories(w http.ResponseWriter, r *http.Request) {
	categories := h.campaignService.GetAvailableCategories()
	utils.WriteSuccess(w, map[string]interface{}{
		"categories": categories,
	})
}

// GetCampaignStatusOptions handles available campaign status options
func (h *CampaignHandlers) GetCampaignStatusOptions(w http.ResponseWriter, r *http.Request) {
	// Get user from context to determine which statuses to show
	claims := middleware.GetUserFromContext(r)
	
	statusOptions := h.campaignService.GetCampaignStatusOptions()
	
	// Filter status options based on user role
	if claims == nil || claims.Role == "donor" {
		// Public users only see active campaigns
		filteredOptions := []map[string]string{}
		for _, option := range statusOptions {
			if option["value"] == "active" {
				filteredOptions = append(filteredOptions, option)
			}
		}
		statusOptions = filteredOptions
	} else if claims.Role == "ngo_admin" {
		// NGO admins see most statuses but not admin-specific ones
		filteredOptions := []map[string]string{}
		for _, option := range statusOptions {
			if option["value"] != "cancelled" { // Only admins can see cancelled status
				filteredOptions = append(filteredOptions, option)
			}
		}
		statusOptions = filteredOptions
	}
	// Platform admins see all statuses

	utils.WriteSuccess(w, map[string]interface{}{
		"status_options": statusOptions,
	})
}

// Helper Methods for Integration with Legacy Handlers

// GetCampaignForDonation retrieves a campaign for donation purposes (validates it's active)
func (h *CampaignHandlers) GetCampaignForDonation(w http.ResponseWriter, r *http.Request) {
	// Extract campaign ID from URL
	vars := mux.Vars(r)
	campaignIDStr := vars["id"]
	
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid campaign ID")
		return
	}

	// Get campaign
	campaign, err := h.campaignService.GetCampaign(r.Context(), campaignID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "campaign not found" {
			status = http.StatusNotFound
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	// Verify campaign is active and can accept donations
	if campaign.Status != "active" {
		utils.WriteError(w, http.StatusBadRequest, "Campaign is not currently accepting donations")
		return
	}

	// Calculate campaign progress
	progressPercentage := 0.0
	if campaign.GoalAmount > 0 {
		progressPercentage = (campaign.RaisedAmount / campaign.GoalAmount) * 100
	}

	response := map[string]interface{}{
		"campaign":            campaign,
		"progress_percentage": progressPercentage,
		"remaining_amount":    campaign.GoalAmount - campaign.RaisedAmount,
		"can_donate":         campaign.Status == "active",
	}

	utils.WriteSuccess(w, response)
}

// BatchUpdateCampaignStatuses handles batch status updates (admin only)
func (h *CampaignHandlers) BatchUpdateCampaignStatuses(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		CampaignIDs []string `json:"campaign_ids"`
		NewStatus   string   `json:"new_status"`
		Reason      string   `json:"reason,omitempty"`
	}

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.CampaignIDs) == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Campaign IDs are required")
		return
	}

	if req.NewStatus == "" {
		utils.WriteError(w, http.StatusBadRequest, "New status is required")
		return
	}

	// Process each campaign
	results := []map[string]interface{}{}
	successCount := 0

	for _, campaignIDStr := range req.CampaignIDs {
		campaignID, err := uuid.Parse(campaignIDStr)
		if err != nil {
			results = append(results, map[string]interface{}{
				"campaign_id": campaignIDStr,
				"success":     false,
				"error":       "Invalid campaign ID",
			})
			continue
		}

		// For suspension, use AdminSuspendCampaign
		if req.NewStatus == "cancelled" {
			err = h.campaignService.AdminSuspendCampaign(r.Context(), claims.UserID, campaignID, req.Reason)
		} else {
			// For other status changes, we would need to implement a generic admin status update
			// For now, we'll return an error for unsupported operations
			err = fmt.Errorf("batch status update to %s not implemented", req.NewStatus)
		}

		if err != nil {
			results = append(results, map[string]interface{}{
				"campaign_id": campaignIDStr,
				"success":     false,
				"error":       err.Error(),
			})
		} else {
			results = append(results, map[string]interface{}{
				"campaign_id": campaignIDStr,
				"success":     true,
			})
			successCount++
		}
	}

	response := map[string]interface{}{
		"results":       results,
		"total_count":   len(req.CampaignIDs),
		"success_count": successCount,
		"failure_count": len(req.CampaignIDs) - successCount,
	}

	utils.WriteSuccess(w, response)
}

// GetCampaignMilestones handles milestone listing for a campaign (placeholder for future implementation)
func (h *CampaignHandlers) GetCampaignMilestones(w http.ResponseWriter, r *http.Request) {
	// Extract campaign ID from URL
	vars := mux.Vars(r)
	campaignIDStr := vars["id"]
	
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid campaign ID")
		return
	}

	// Verify campaign exists
	_, err = h.campaignService.GetCampaign(r.Context(), campaignID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "campaign not found" {
			status = http.StatusNotFound
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	// TODO: Implement milestone listing when milestone repository is ready
	// For now, return empty list
	utils.WriteSuccess(w, map[string]interface{}{
		"campaign_id": campaignID,
		"milestones":  []interface{}{},
		"message":     "Milestone management not yet implemented",
	})
}

// GetCampaignExpenses handles expense listing for a campaign (placeholder for future implementation)  
func (h *CampaignHandlers) GetCampaignExpenses(w http.ResponseWriter, r *http.Request) {
	// Extract campaign ID from URL
	vars := mux.Vars(r)
	campaignIDStr := vars["id"]
	
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid campaign ID")
		return
	}

	// Verify campaign exists
	_, err = h.campaignService.GetCampaign(r.Context(), campaignID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "campaign not found" {
			status = http.StatusNotFound
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	// TODO: Implement expense listing when expense repository is ready
	// For now, return empty list
	utils.WriteSuccess(w, map[string]interface{}{
		"campaign_id": campaignID,
		"expenses":    []interface{}{},
		"message":     "Expense management not yet implemented",
	})
}