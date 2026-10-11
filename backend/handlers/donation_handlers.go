package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"ngo-donation-tracker/middleware"
	"ngo-donation-tracker/repositories"
	"ngo-donation-tracker/services"
	"ngo-donation-tracker/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// DonationHandlers contains all donation-related HTTP handlers
type DonationHandlers struct {
	donationService *services.DonationService
}

// NewDonationHandlers creates a new instance of DonationHandlers
func NewDonationHandlers(donationService *services.DonationService) *DonationHandlers {
	return &DonationHandlers{
		donationService: donationService,
	}
}

// Public Donation Handlers

// CreateDonation handles donation creation
func (h *DonationHandlers) CreateDonation(w http.ResponseWriter, r *http.Request) {
	// Extract user from context (must be authenticated)
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req services.CreateDonationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create donation
	donation, err := h.donationService.CreateDonation(user.UserID, &req)
	if err != nil {
		if err.Error() == "invalid campaign" || err.Error() == "campaign is not active" || 
		   err.Error() == "campaign has ended" || err.Error() == "campaign has reached its funding goal" {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create donation")
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, donation)
}

// ProcessDonation handles donation processing (payment completion)
func (h *DonationHandlers) ProcessDonation(w http.ResponseWriter, r *http.Request) {
	// Extract donation ID from URL
	vars := mux.Vars(r)
	donationIDStr := vars["id"]
	donationID, err := uuid.Parse(donationIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid donation ID")
		return
	}

	var req services.ProcessDonationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Process donation
	donation, err := h.donationService.ProcessDonation(donationID, &req)
	if err != nil {
		if err.Error() == "donation not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, "Donation not found")
			return
		}
		if err.Error() == "donation is not pending" {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "Donation is not pending")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to process donation")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, donation)
}

// GetDonation handles retrieving a single donation
func (h *DonationHandlers) GetDonation(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract donation ID from URL
	vars := mux.Vars(r)
	donationIDStr := vars["id"]
	donationID, err := uuid.Parse(donationIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid donation ID")
		return
	}

	// Get donation
	donation, err := h.donationService.GetDonation(donationID, user.UserID)
	if err != nil {
		if err.Error() == "donation not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, "Donation not found")
			return
		}
		if err.Error() == "unauthorized to view this donation" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get donation")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, donation)
}

// Donor-specific Handlers

// GetMyDonations handles retrieving donations for the authenticated donor
func (h *DonationHandlers) GetMyDonations(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Parse pagination parameters
	pagination := utils.GetPaginationParams(r)

	// Get donor donations
	donations, result, err := h.donationService.GetDonationsByDonor(user.UserID, pagination)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get donations")
		return
	}

	response := map[string]interface{}{
		"donations":  donations,
		"pagination": result,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// GetMyDonationHistory handles retrieving donation history with campaign details
func (h *DonationHandlers) GetMyDonationHistory(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Parse pagination parameters
	pagination := utils.GetPaginationParams(r)

	// Get donation history
	history, err := h.donationService.GetDonationHistory(user.UserID, pagination)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get donation history")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, history)
}

// RefundMyDonation handles donation refund request by donor
func (h *DonationHandlers) RefundMyDonation(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract donation ID from URL
	vars := mux.Vars(r)
	donationIDStr := vars["id"]
	donationID, err := uuid.Parse(donationIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid donation ID")
		return
	}

	var req services.RefundDonationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Process refund
	donation, err := h.donationService.RefundDonation(donationID, &req, user.UserID)
	if err != nil {
		if err.Error() == "donation not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, "Donation not found")
			return
		}
		if err.Error() == "unauthorized to refund this donation" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		if err.Error() == "only completed donations can be refunded" {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "Only completed donations can be refunded")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to process refund")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, donation)
}

// NGO-specific Handlers

// GetCampaignDonations handles retrieving donations for NGO campaigns
func (h *DonationHandlers) GetCampaignDonations(w http.ResponseWriter, r *http.Request) {
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

	// Parse pagination parameters
	pagination := utils.GetPaginationParams(r)

	// Get campaign donations
	donations, result, err := h.donationService.GetDonationsByCampaign(campaignID, user.UserID, pagination)
	if err != nil {
		if err.Error() == "invalid campaign" {
			utils.WriteErrorResponse(w, http.StatusNotFound, "Campaign not found")
			return
		}
		if err.Error() == "unauthorized to view campaign donations" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get campaign donations")
		return
	}

	response := map[string]interface{}{
		"donations":  donations,
		"pagination": result,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// GetCampaignDonationStats handles retrieving donation statistics for a campaign
func (h *DonationHandlers) GetCampaignDonationStats(w http.ResponseWriter, r *http.Request) {
	// Extract campaign ID from URL
	vars := mux.Vars(r)
	campaignIDStr := vars["campaignId"]
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid campaign ID")
		return
	}

	// Get donation stats (public endpoint)
	stats, err := h.donationService.GetDonationStats(campaignID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get donation stats")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// Admin-specific Handlers

// ListAllDonations handles listing all donations (admin only)
func (h *DonationHandlers) ListAllDonations(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Parse pagination parameters
	pagination := utils.GetPaginationParams(r)

	// Parse filters
	filters := repositories.DonationFiltersExtended{}
	
	if donorIDStr := r.URL.Query().Get("donor_id"); donorIDStr != "" {
		if donorID, err := uuid.Parse(donorIDStr); err == nil {
			filters.DonorID = &donorID
		}
	}
	
	if campaignIDStr := r.URL.Query().Get("campaign_id"); campaignIDStr != "" {
		if campaignID, err := uuid.Parse(campaignIDStr); err == nil {
			filters.CampaignID = &campaignID
		}
	}
	
	if ngoIDStr := r.URL.Query().Get("ngo_id"); ngoIDStr != "" {
		if ngoID, err := uuid.Parse(ngoIDStr); err == nil {
			filters.NGOID = &ngoID
		}
	}

	if status := r.URL.Query().Get("status"); status != "" {
		filters.Status = status
	}

	if minAmountStr := r.URL.Query().Get("min_amount"); minAmountStr != "" {
		if minAmount, err := strconv.ParseFloat(minAmountStr, 64); err == nil {
			filters.MinAmount = &minAmount
		}
	}

	if maxAmountStr := r.URL.Query().Get("max_amount"); maxAmountStr != "" {
		if maxAmount, err := strconv.ParseFloat(maxAmountStr, 64); err == nil {
			filters.MaxAmount = &maxAmount
		}
	}

	if startDateStr := r.URL.Query().Get("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filters.StartDate = &startDate
		}
	}

	if endDateStr := r.URL.Query().Get("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filters.EndDate = &endDate
		}
	}

	if anonymousStr := r.URL.Query().Get("anonymous"); anonymousStr != "" {
		if anonymous, err := strconv.ParseBool(anonymousStr); err == nil {
			filters.Anonymous = &anonymous
		}
	}

	if searchTerm := r.URL.Query().Get("search"); searchTerm != "" {
		filters.SearchTerm = searchTerm
	}

	// List donations
	donations, result, err := h.donationService.ListDonations(filters, pagination, user.UserID)
	if err != nil {
		if err.Error() == "unauthorized to list all donations" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to list donations")
		return
	}

	response := map[string]interface{}{
		"donations":  donations,
		"pagination": result,
		"filters":    filters,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// GetDonationsByDateRange handles getting donations within date range (admin only)
func (h *DonationHandlers) GetDonationsByDateRange(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Parse date parameters
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "start_date and end_date are required")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid start_date format (use YYYY-MM-DD)")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid end_date format (use YYYY-MM-DD)")
		return
	}

	// Validate date range
	if endDate.Before(startDate) {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "end_date must be after start_date")
		return
	}

	// Parse pagination parameters
	pagination := utils.GetPaginationParams(r)

	// Get donations by date range
	donations, result, err := h.donationService.GetDonationsByDateRange(startDate, endDate, user.UserID, pagination)
	if err != nil {
		if err.Error() == "unauthorized to view donations by date range" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get donations by date range")
		return
	}

	response := map[string]interface{}{
		"donations":   donations,
		"pagination":  result,
		"date_range": map[string]string{
			"start_date": startDateStr,
			"end_date":   endDateStr,
		},
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// AdminRefundDonation handles admin refund of donations
func (h *DonationHandlers) AdminRefundDonation(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract donation ID from URL
	vars := mux.Vars(r)
	donationIDStr := vars["id"]
	donationID, err := uuid.Parse(donationIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid donation ID")
		return
	}

	var req services.RefundDonationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Process refund
	donation, err := h.donationService.RefundDonation(donationID, &req, user.UserID)
	if err != nil {
		if err.Error() == "donation not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, "Donation not found")
			return
		}
		if err.Error() == "unauthorized to refund this donation" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		if err.Error() == "only completed donations can be refunded" {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "Only completed donations can be refunded")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to process refund")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, donation)
}

// GetDonationDetails handles getting detailed donation information (admin only)
func (h *DonationHandlers) GetDonationDetails(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract donation ID from URL
	vars := mux.Vars(r)
	donationIDStr := vars["id"]
	donationID, err := uuid.Parse(donationIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid donation ID")
		return
	}

	// Get donation details
	donation, err := h.donationService.GetDonation(donationID, user.UserID)
	if err != nil {
		if err.Error() == "donation not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, "Donation not found")
			return
		}
		if err.Error() == "unauthorized to view this donation" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get donation details")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, donation)
}