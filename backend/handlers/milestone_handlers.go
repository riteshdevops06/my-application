package handlers

import (
	"encoding/json"
	"net/http"

	"ngo-donation-tracker/middleware"
	"ngo-donation-tracker/repositories"
	"ngo-donation-tracker/services"
	"ngo-donation-tracker/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// MilestoneHandlers contains all milestone-related HTTP handlers
type MilestoneHandlers struct {
	milestoneService *services.MilestoneService
}

// NewMilestoneHandlers creates a new instance of MilestoneHandlers
func NewMilestoneHandlers(milestoneService *services.MilestoneService) *MilestoneHandlers {
	return &MilestoneHandlers{
		milestoneService: milestoneService,
	}
}

// NGO Milestone Management Handlers

// CreateMilestone handles milestone creation
func (h *MilestoneHandlers) CreateMilestone(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req services.CreateMilestoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create milestone
	milestone, err := h.milestoneService.CreateMilestone(user.UserID, &req)
	if err != nil {
		if err.Error() == "invalid campaign" {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid campaign")
			return
		}
		if err.Error() == "unauthorized to create milestones" || 
		   err.Error() == "unauthorized to create milestones for this campaign" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create milestone")
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, milestone)
}

// UpdateMilestone handles milestone updates
func (h *MilestoneHandlers) UpdateMilestone(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract milestone ID from URL
	vars := mux.Vars(r)
	milestoneIDStr := vars["id"]
	milestoneID, err := uuid.Parse(milestoneIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid milestone ID")
		return
	}

	var req services.UpdateMilestoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update milestone
	milestone, err := h.milestoneService.UpdateMilestone(milestoneID, user.UserID, &req)
	if err != nil {
		if err.Error() == "milestone not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, "Milestone not found")
			return
		}
		if err.Error() == "unauthorized to access this milestone" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		if err.Error() == "completed milestones cannot be updated" {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "Completed milestones cannot be updated")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update milestone")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, milestone)
}

// DeleteMilestone handles milestone deletion
func (h *MilestoneHandlers) DeleteMilestone(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract milestone ID from URL
	vars := mux.Vars(r)
	milestoneIDStr := vars["id"]
	milestoneID, err := uuid.Parse(milestoneIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid milestone ID")
		return
	}

	// Delete milestone
	err = h.milestoneService.DeleteMilestone(milestoneID, user.UserID)
	if err != nil {
		if err.Error() == "milestone not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, "Milestone not found")
			return
		}
		if err.Error() == "unauthorized to access this milestone" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		if err.Error() == "completed milestones cannot be deleted" {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "Completed milestones cannot be deleted")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete milestone")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Milestone deleted successfully"})
}

// GetMilestone handles retrieving a single milestone
func (h *MilestoneHandlers) GetMilestone(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract milestone ID from URL
	vars := mux.Vars(r)
	milestoneIDStr := vars["id"]
	milestoneID, err := uuid.Parse(milestoneIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid milestone ID")
		return
	}

	// Get milestone
	milestone, err := h.milestoneService.GetMilestone(milestoneID, user.UserID)
	if err != nil {
		if err.Error() == "milestone not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, "Milestone not found")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get milestone")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, milestone)
}

// ListCampaignMilestones handles retrieving milestones for a campaign
func (h *MilestoneHandlers) ListCampaignMilestones(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract campaign ID from URL
	vars := mux.Vars(r)
	campaignIDStr := vars["campaignId"]
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid campaign ID")
		return
	}

	// List milestones
	milestones, err := h.milestoneService.ListMilestonesByCampaign(campaignID, user.UserID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get milestones")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, milestones)
}

// CompleteMilestone handles marking a milestone as complete
func (h *MilestoneHandlers) CompleteMilestone(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract milestone ID from URL
	vars := mux.Vars(r)
	milestoneIDStr := vars["id"]
	milestoneID, err := uuid.Parse(milestoneIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid milestone ID")
		return
	}

	var req services.CompleteMilestoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Complete milestone
	milestone, err := h.milestoneService.CompleteMilestone(milestoneID, user.UserID, &req)
	if err != nil {
		if err.Error() == "milestone not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, "Milestone not found")
			return
		}
		if err.Error() == "unauthorized to access this milestone" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		if err.Error() == "milestone is already completed" {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "Milestone is already completed")
			return
		}
		if err.Error() == "cancelled milestones cannot be completed" {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "Cancelled milestones cannot be completed")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to complete milestone")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, milestone)
}

// UpdateMilestoneStatus handles updating milestone status
func (h *MilestoneHandlers) UpdateMilestoneStatus(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract milestone ID from URL
	vars := mux.Vars(r)
	milestoneIDStr := vars["id"]
	milestoneID, err := uuid.Parse(milestoneIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid milestone ID")
		return
	}

	var reqBody struct {
		Status string `json:"status" validate:"required,oneof=pending in_progress completed cancelled"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&reqBody); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update milestone status
	milestone, err := h.milestoneService.UpdateMilestoneStatus(milestoneID, reqBody.Status, user.UserID)
	if err != nil {
		if err.Error() == "milestone not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, "Milestone not found")
			return
		}
		if err.Error() == "unauthorized to access this milestone" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		if err.Error() == "completed milestones cannot change status" {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "Completed milestones cannot change status")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update milestone status")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, milestone)
}

// GetCampaignMilestoneStats handles retrieving milestone statistics for a campaign
func (h *MilestoneHandlers) GetCampaignMilestoneStats(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract campaign ID from URL
	vars := mux.Vars(r)
	campaignIDStr := vars["campaignId"]
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid campaign ID")
		return
	}

	// Get milestone stats
	stats, err := h.milestoneService.GetMilestoneStats(campaignID, user.UserID)
	if err != nil {
		if err.Error() == "unauthorized to access this milestone" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get milestone stats")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// Admin-specific Handlers

// ListMilestonesByStatus handles listing milestones by status (admin only)
func (h *MilestoneHandlers) ListMilestonesByStatus(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Parse status parameter
	status := r.URL.Query().Get("status")
	if status == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Status parameter is required")
		return
	}

	// Parse pagination parameters
	pagination := utils.GetPaginationParams(r)

	// List milestones
	milestones, result, err := h.milestoneService.GetMilestonesByStatus(status, pagination, user.UserID)
	if err != nil {
		if err.Error() == "unauthorized to view all milestones" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to list milestones")
		return
	}

	response := map[string]interface{}{
		"milestones": milestones,
		"pagination": result,
		"status":     status,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// Public Handlers

// GetPublicMilestone handles retrieving public milestone information
func (h *MilestoneHandlers) GetPublicMilestone(w http.ResponseWriter, r *http.Request) {
	// Extract milestone ID from URL
	vars := mux.Vars(r)
	milestoneIDStr := vars["id"]
	milestoneID, err := uuid.Parse(milestoneIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid milestone ID")
		return
	}

	// For public access, we use a placeholder user ID (donors can see basic milestone info)
	// This will return limited information without sensitive details
	milestone, err := h.milestoneService.GetMilestone(milestoneID, uuid.Nil)
	if err != nil {
		if err.Error() == "milestone not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, "Milestone not found")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get milestone")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, milestone)
}

// ListPublicCampaignMilestones handles retrieving public milestones for a campaign
func (h *MilestoneHandlers) ListPublicCampaignMilestones(w http.ResponseWriter, r *http.Request) {
	// Extract campaign ID from URL
	vars := mux.Vars(r)
	campaignIDStr := vars["campaignId"]
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid campaign ID")
		return
	}

	// For public access, use placeholder user ID to get limited information
	milestones, err := h.milestoneService.ListMilestonesByCampaign(campaignID, uuid.Nil)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get milestones")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, milestones)
}