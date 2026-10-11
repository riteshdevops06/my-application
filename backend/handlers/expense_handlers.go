package handlers

import (
	"encoding/json"
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

// ExpenseHandlers contains all expense-related HTTP handlers
type ExpenseHandlers struct {
	expenseService *services.ExpenseService
}

// NewExpenseHandlers creates a new instance of ExpenseHandlers
func NewExpenseHandlers(expenseService *services.ExpenseService) *ExpenseHandlers {
	return &ExpenseHandlers{
		expenseService: expenseService,
	}
}

// NGO Expense Management Handlers

// CreateExpense handles expense creation
func (h *ExpenseHandlers) CreateExpense(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req services.CreateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create expense
	expense, err := h.expenseService.CreateExpense(user.UserID, &req)
	if err != nil {
		if err.Error() == "invalid campaign" {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid campaign")
			return
		}
		if err.Error() == "invalid milestone" {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid milestone")
			return
		}
		if err.Error() == "milestone does not belong to the specified campaign" {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "Milestone does not belong to campaign")
			return
		}
		if err.Error() == "unauthorized to access expenses for this campaign" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create expense")
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, expense)
}

// UpdateExpense handles expense updates
func (h *ExpenseHandlers) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract expense ID from URL
	vars := mux.Vars(r)
	expenseIDStr := vars["id"]
	expenseID, err := uuid.Parse(expenseIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid expense ID")
		return
	}

	var req services.UpdateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update expense
	expense, err := h.expenseService.UpdateExpense(expenseID, user.UserID, &req)
	if err != nil {
		if err.Error() == "expense not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, "Expense not found")
			return
		}
		if err.Error() == "unauthorized to access expenses for this campaign" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		if err.Error() == "approved or rejected expenses cannot be updated" {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "Approved or rejected expenses cannot be updated")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update expense")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, expense)
}

// DeleteExpense handles expense deletion
func (h *ExpenseHandlers) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract expense ID from URL
	vars := mux.Vars(r)
	expenseIDStr := vars["id"]
	expenseID, err := uuid.Parse(expenseIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid expense ID")
		return
	}

	// Delete expense
	err = h.expenseService.DeleteExpense(expenseID, user.UserID)
	if err != nil {
		if err.Error() == "expense not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, "Expense not found")
			return
		}
		if err.Error() == "unauthorized to access expenses for this campaign" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		if err.Error() == "approved expenses cannot be deleted" {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "Approved expenses cannot be deleted")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete expense")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Expense deleted successfully"})
}

// GetExpense handles retrieving a single expense
func (h *ExpenseHandlers) GetExpense(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract expense ID from URL
	vars := mux.Vars(r)
	expenseIDStr := vars["id"]
	expenseID, err := uuid.Parse(expenseIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid expense ID")
		return
	}

	// Get expense
	expense, err := h.expenseService.GetExpense(expenseID, user.UserID)
	if err != nil {
		if err.Error() == "expense not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, "Expense not found")
			return
		}
		if err.Error() == "unauthorized to access expenses for this campaign" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get expense")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, expense)
}

// ListCampaignExpenses handles retrieving expenses for a campaign
func (h *ExpenseHandlers) ListCampaignExpenses(w http.ResponseWriter, r *http.Request) {
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

	// List expenses
	expenses, result, err := h.expenseService.ListExpensesByCampaign(campaignID, pagination, user.UserID)
	if err != nil {
		if err.Error() == "unauthorized to access expenses for this campaign" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get expenses")
		return
	}

	response := map[string]interface{}{
		"expenses":   expenses,
		"pagination": result,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// ListMilestoneExpenses handles retrieving expenses for a milestone
func (h *ExpenseHandlers) ListMilestoneExpenses(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract milestone ID from URL
	vars := mux.Vars(r)
	milestoneIDStr := vars["milestoneId"]
	milestoneID, err := uuid.Parse(milestoneIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid milestone ID")
		return
	}

	// List expenses
	expenses, err := h.expenseService.ListExpensesByMilestone(milestoneID, user.UserID)
	if err != nil {
		if err.Error() == "milestone not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, "Milestone not found")
			return
		}
		if err.Error() == "unauthorized to access expenses for this campaign" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get expenses")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, expenses)
}

// GetCampaignExpenseStats handles retrieving expense statistics for a campaign
func (h *ExpenseHandlers) GetCampaignExpenseStats(w http.ResponseWriter, r *http.Request) {
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

	// Get expense stats
	stats, err := h.expenseService.GetExpenseStats(campaignID, user.UserID)
	if err != nil {
		if err.Error() == "unauthorized to access expenses for this campaign" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get expense stats")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// Admin-specific Handlers

// ApproveExpense handles expense approval/rejection
func (h *ExpenseHandlers) ApproveExpense(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract expense ID from URL
	vars := mux.Vars(r)
	expenseIDStr := vars["id"]
	expenseID, err := uuid.Parse(expenseIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid expense ID")
		return
	}

	var req services.ApproveExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Approve/reject expense
	expense, err := h.expenseService.ApproveExpense(expenseID, user.UserID, &req)
	if err != nil {
		if err.Error() == "expense not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, "Expense not found")
			return
		}
		if err.Error() == "unauthorized to approve expenses" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		if err.Error() == "only pending expenses can be approved or rejected" {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "Only pending expenses can be approved or rejected")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to process expense approval")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, expense)
}

// ListAllExpenses handles listing all expenses with filters (admin only)
func (h *ExpenseHandlers) ListAllExpenses(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Parse pagination parameters
	pagination := utils.GetPaginationParams(r)

	// Parse filters
	filters := repositories.ExpenseFilters{}
	
	if campaignIDStr := r.URL.Query().Get("campaign_id"); campaignIDStr != "" {
		if campaignID, err := uuid.Parse(campaignIDStr); err == nil {
			filters.CampaignID = &campaignID
		}
	}
	
	if milestoneIDStr := r.URL.Query().Get("milestone_id"); milestoneIDStr != "" {
		if milestoneID, err := uuid.Parse(milestoneIDStr); err == nil {
			filters.MilestoneID = &milestoneID
		}
	}

	if status := r.URL.Query().Get("status"); status != "" {
		filters.Status = status
	}

	if category := r.URL.Query().Get("category"); category != "" {
		filters.Category = category
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

	if approvedByStr := r.URL.Query().Get("approved_by"); approvedByStr != "" {
		if approvedBy, err := uuid.Parse(approvedByStr); err == nil {
			filters.ApprovedBy = &approvedBy
		}
	}

	if searchTerm := r.URL.Query().Get("search"); searchTerm != "" {
		filters.SearchTerm = searchTerm
	}

	// List expenses
	expenses, result, err := h.expenseService.ListExpenses(filters, pagination, user.UserID)
	if err != nil {
		if err.Error() == "unauthorized to list all expenses" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to list expenses")
		return
	}

	response := map[string]interface{}{
		"expenses":   expenses,
		"pagination": result,
		"filters":    filters,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// ListExpensesByStatus handles listing expenses by status (admin only)
func (h *ExpenseHandlers) ListExpensesByStatus(w http.ResponseWriter, r *http.Request) {
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

	// List expenses
	expenses, result, err := h.expenseService.GetExpensesByStatus(status, pagination, user.UserID)
	if err != nil {
		if err.Error() == "unauthorized to view all expenses" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to list expenses")
		return
	}

	response := map[string]interface{}{
		"expenses":   expenses,
		"pagination": result,
		"status":     status,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// GetExpenseDetails handles getting detailed expense information (admin only)
func (h *ExpenseHandlers) GetExpenseDetails(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract expense ID from URL
	vars := mux.Vars(r)
	expenseIDStr := vars["id"]
	expenseID, err := uuid.Parse(expenseIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid expense ID")
		return
	}

	// Get expense details
	expense, err := h.expenseService.GetExpense(expenseID, user.UserID)
	if err != nil {
		if err.Error() == "expense not found" {
			utils.WriteErrorResponse(w, http.StatusNotFound, "Expense not found")
			return
		}
		if err.Error() == "unauthorized to access expenses for this campaign" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get expense details")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, expense)
}