package services

import (
	"context"
	"fmt"
	"time"

	"ngo-donation-tracker/models"
	"ngo-donation-tracker/repositories"

	"github.com/google/uuid"
)

// ExpenseService handles business logic for expenses
type ExpenseService struct {
	expenseRepo   repositories.ExpenseRepository
	milestoneRepo repositories.MilestoneRepository
	campaignRepo  repositories.CampaignRepositoryInterface
	userRepo      repositories.UserRepositoryInterface
	ngoRepo       repositories.NGORepositoryInterface
}

// NewExpenseService creates a new instance of ExpenseService
func NewExpenseService(
	expenseRepo repositories.ExpenseRepository,
	milestoneRepo repositories.MilestoneRepository,
	campaignRepo repositories.CampaignRepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
	ngoRepo repositories.NGORepositoryInterface,
) *ExpenseService {
	return &ExpenseService{
		expenseRepo:   expenseRepo,
		milestoneRepo: milestoneRepo,
		campaignRepo:  campaignRepo,
		userRepo:      userRepo,
		ngoRepo:       ngoRepo,
	}
}

// CreateExpenseRequest represents request to create an expense
type CreateExpenseRequest struct {
	CampaignID  uuid.UUID `json:"campaign_id" validate:"required"`
	MilestoneID uuid.UUID `json:"milestone_id,omitempty"`
	Title       string    `json:"title" validate:"required,min=3,max=255"`
	Description string    `json:"description,omitempty"`
	Amount      float64   `json:"amount" validate:"required,gt=0"`
	Currency    string    `json:"currency,omitempty"`
	Category    string    `json:"category,omitempty"`
	ReceiptURL  string    `json:"receipt_url,omitempty"`
	ExpenseDate string    `json:"expense_date,omitempty"`
}

// UpdateExpenseRequest represents request to update an expense
type UpdateExpenseRequest struct {
	Title       string  `json:"title,omitempty" validate:"omitempty,min=3,max=255"`
	Description string  `json:"description,omitempty"`
	Amount      float64 `json:"amount,omitempty" validate:"omitempty,gt=0"`
	Category    string  `json:"category,omitempty"`
	ReceiptURL  string  `json:"receipt_url,omitempty"`
	ExpenseDate string  `json:"expense_date,omitempty"`
}

// ApproveExpenseRequest represents request to approve/reject an expense
type ApproveExpenseRequest struct {
	Status string `json:"status" validate:"required,oneof=approved rejected"`
	Notes  string `json:"notes,omitempty"`
}

// ExpenseResponse represents expense response with additional details
type ExpenseResponse struct {
	*models.Expense
	CampaignTitle  string `json:"campaign_title,omitempty"`
	MilestoneTitle string `json:"milestone_title,omitempty"`
	NGOName        string `json:"ngo_name,omitempty"`
	ApproverName   string `json:"approver_name,omitempty"`
}

// CreateExpense creates a new expense
func (s *ExpenseService) CreateExpense(requesterID uuid.UUID, req *CreateExpenseRequest) (*models.Expense, error) {
	// Validate campaign exists and user has permission
	_, err := s.campaignRepo.GetByID(context.Background(), req.CampaignID)
	if err != nil {
		return nil, fmt.Errorf("invalid campaign: %v", err)
	}

	// Check if requester is authorized to create expenses for this campaign
	err = s.validateExpenseAccess(req.CampaignID, requesterID)
	if err != nil {
		return nil, err
	}

	// Validate milestone if provided
	var milestoneID *uuid.UUID
	if req.MilestoneID != uuid.Nil {
		milestone, err := s.milestoneRepo.GetMilestoneByID(req.MilestoneID)
		if err != nil {
			return nil, fmt.Errorf("invalid milestone: %v", err)
		}

		// Ensure milestone belongs to the same campaign
		if milestone.CampaignID != req.CampaignID {
			return nil, fmt.Errorf("milestone does not belong to the specified campaign")
		}

		milestoneID = &req.MilestoneID
	}

	// Parse expense date if provided, otherwise use current date
	expenseDate := time.Now()
	if req.ExpenseDate != "" {
		parsed, err := time.Parse("2006-01-02", req.ExpenseDate)
		if err != nil {
			return nil, fmt.Errorf("invalid expense date format (use YYYY-MM-DD): %v", err)
		}
		expenseDate = parsed
	}

	// Set default currency if not provided
	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}

	// Create expense
	expense := &models.Expense{
		ID:          uuid.New(),
		CampaignID:  req.CampaignID,
		MilestoneID: milestoneID,
		Title:       req.Title,
		Description: &req.Description,
		Amount:      req.Amount,
		Currency:    currency,
		Category:    &req.Category,
		ReceiptURL:  &req.ReceiptURL,
		ExpenseDate: expenseDate,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.expenseRepo.CreateExpense(expense)
	if err != nil {
		return nil, fmt.Errorf("failed to create expense: %v", err)
	}

	return expense, nil
}

// UpdateExpense updates an existing expense
func (s *ExpenseService) UpdateExpense(expenseID uuid.UUID, requesterID uuid.UUID, req *UpdateExpenseRequest) (*models.Expense, error) {
	// Get existing expense
	expense, err := s.expenseRepo.GetExpenseByID(expenseID)
	if err != nil {
		return nil, fmt.Errorf("expense not found: %v", err)
	}

	// Check authorization
	err = s.validateExpenseAccess(expense.CampaignID, requesterID)
	if err != nil {
		return nil, err
	}

	// Prevent updates if expense is approved or rejected
	if expense.Status == "approved" || expense.Status == "rejected" {
		return nil, fmt.Errorf("approved or rejected expenses cannot be updated")
	}

	// Update fields if provided
	if req.Title != "" {
		expense.Title = req.Title
	}
	if req.Description != "" {
		expense.Description = &req.Description
	}
	if req.Amount > 0 {
		expense.Amount = req.Amount
	}
	if req.Category != "" {
		expense.Category = &req.Category
	}
	if req.ReceiptURL != "" {
		expense.ReceiptURL = &req.ReceiptURL
	}
	if req.ExpenseDate != "" {
		parsed, err := time.Parse("2006-01-02", req.ExpenseDate)
		if err != nil {
			return nil, fmt.Errorf("invalid expense date format (use YYYY-MM-DD): %v", err)
		}
		expense.ExpenseDate = parsed
	}

	expense.UpdatedAt = time.Now()

	err = s.expenseRepo.UpdateExpense(expense)
	if err != nil {
		return nil, fmt.Errorf("failed to update expense: %v", err)
	}

	return expense, nil
}

// DeleteExpense deletes an expense
func (s *ExpenseService) DeleteExpense(expenseID uuid.UUID, requesterID uuid.UUID) error {
	// Get existing expense
	expense, err := s.expenseRepo.GetExpenseByID(expenseID)
	if err != nil {
		return fmt.Errorf("expense not found: %v", err)
	}

	// Check authorization
	err = s.validateExpenseAccess(expense.CampaignID, requesterID)
	if err != nil {
		return err
	}

	// Prevent deletion if expense is approved
	if expense.Status == "approved" {
		return fmt.Errorf("approved expenses cannot be deleted")
	}

	err = s.expenseRepo.DeleteExpense(expenseID)
	if err != nil {
		return fmt.Errorf("failed to delete expense: %v", err)
	}

	return nil
}

// GetExpense retrieves an expense with additional details
func (s *ExpenseService) GetExpense(expenseID uuid.UUID, requesterID uuid.UUID) (*ExpenseResponse, error) {
	expense, err := s.expenseRepo.GetExpenseByID(expenseID)
	if err != nil {
		return nil, fmt.Errorf("expense not found: %v", err)
	}

	// Check authorization
	err = s.validateExpenseAccess(expense.CampaignID, requesterID)
	if err != nil {
		return nil, err
	}

	response := &ExpenseResponse{
		Expense: expense,
	}

	// Add campaign details
	campaign, err := s.campaignRepo.GetByID(context.Background(), expense.CampaignID)
	if err == nil {
		response.CampaignTitle = campaign.Title
		if campaign.NGO != nil {
			response.NGOName = campaign.NGO.Name
		}
	}

	// Add milestone details if applicable
	if expense.MilestoneID != nil {
		milestone, err := s.milestoneRepo.GetMilestoneByID(*expense.MilestoneID)
		if err == nil {
			response.MilestoneTitle = milestone.Title
		}
	}

	// Add approver details if applicable
	if expense.ApprovedBy != nil {
		approver, err := s.userRepo.GetUserByID(context.Background(), *expense.ApprovedBy)
		if err == nil {
			response.ApproverName = approver.Name
		}
	}

	return response, nil
}

// ListExpensesByCampaign retrieves expenses for a campaign
func (s *ExpenseService) ListExpensesByCampaign(campaignID uuid.UUID, pagination repositories.PaginationParams, requesterID uuid.UUID) ([]*ExpenseResponse, *repositories.PaginationResult, error) {
	// Check authorization
	err := s.validateExpenseAccess(campaignID, requesterID)
	if err != nil {
		return nil, nil, err
	}

	expenses, result, err := s.expenseRepo.ListExpensesByCampaign(campaignID, pagination)
	if err != nil {
		return nil, nil, err
	}

	var responses []*ExpenseResponse
	for _, expense := range expenses {
		response := &ExpenseResponse{
			Expense: expense,
		}

		// Add milestone details if applicable
		if expense.MilestoneID != nil {
			milestone, err := s.milestoneRepo.GetMilestoneByID(*expense.MilestoneID)
			if err == nil {
				response.MilestoneTitle = milestone.Title
			}
		}

		// Add approver details if applicable
		if expense.ApprovedBy != nil {
			approver, err := s.userRepo.GetUserByID(context.Background(), *expense.ApprovedBy)
			if err == nil {
				response.ApproverName = approver.Name
			}
		}

		responses = append(responses, response)
	}

	return responses, result, nil
}

// ListExpensesByMilestone retrieves expenses for a milestone
func (s *ExpenseService) ListExpensesByMilestone(milestoneID uuid.UUID, requesterID uuid.UUID) ([]*ExpenseResponse, error) {
	// Get milestone to check campaign access
	milestone, err := s.milestoneRepo.GetMilestoneByID(milestoneID)
	if err != nil {
		return nil, fmt.Errorf("milestone not found: %v", err)
	}

	// Check authorization
	err = s.validateExpenseAccess(milestone.CampaignID, requesterID)
	if err != nil {
		return nil, err
	}

	expenses, err := s.expenseRepo.ListExpensesByMilestone(milestoneID)
	if err != nil {
		return nil, err
	}

	var responses []*ExpenseResponse
	for _, expense := range expenses {
		response := &ExpenseResponse{
			Expense:        expense,
			MilestoneTitle: milestone.Title,
		}

		// Add campaign details
		campaign, err := s.campaignRepo.GetByID(context.Background(), expense.CampaignID)
		if err == nil {
			response.CampaignTitle = campaign.Title
			if campaign.NGO != nil {
				response.NGOName = campaign.NGO.Name
			}
		}

		// Add approver details if applicable
		if expense.ApprovedBy != nil {
			approver, err := s.userRepo.GetUserByID(context.Background(), *expense.ApprovedBy)
			if err == nil {
				response.ApproverName = approver.Name
			}
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// ApproveExpense approves or rejects an expense
func (s *ExpenseService) ApproveExpense(expenseID uuid.UUID, requesterID uuid.UUID, req *ApproveExpenseRequest) (*models.Expense, error) {
	// Get existing expense
	expense, err := s.expenseRepo.GetExpenseByID(expenseID)
	if err != nil {
		return nil, fmt.Errorf("expense not found: %v", err)
	}

	// Check authorization - only platform admin can approve expenses
	requester, err := s.userRepo.GetUserByID(context.Background(), requesterID)
	if err != nil {
		return nil, fmt.Errorf("invalid requester: %v", err)
	}

	if requester.Role != "platform_admin" {
		return nil, fmt.Errorf("unauthorized to approve expenses")
	}

	// Validate expense can be approved/rejected
	if expense.Status != "pending" {
		return nil, fmt.Errorf("only pending expenses can be approved or rejected")
	}

	// Update expense status
	if req.Status == "approved" {
		err = s.expenseRepo.ApproveExpense(expenseID, requesterID)
	} else {
		reason := req.Notes
		if reason == "" {
			reason = "No reason provided"
		}
		err = s.expenseRepo.RejectExpense(expenseID, requesterID, reason)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to %s expense: %v", req.Status, err)
	}

	// Get updated expense
	updatedExpense, err := s.expenseRepo.GetExpenseByID(expenseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated expense: %v", err)
	}

	// TODO: Send notification to NGO about expense approval/rejection

	return updatedExpense, nil
}

// GetExpenseStats retrieves expense statistics for a campaign
func (s *ExpenseService) GetExpenseStats(campaignID uuid.UUID, requesterID uuid.UUID) (*repositories.ExpenseStats, error) {
	// Check authorization
	err := s.validateExpenseAccess(campaignID, requesterID)
	if err != nil {
		return nil, err
	}

	return s.expenseRepo.GetExpenseStats(campaignID)
}

// ListExpenses retrieves expenses based on filters (admin only)
func (s *ExpenseService) ListExpenses(filters repositories.ExpenseFilters, pagination repositories.PaginationParams, requesterID uuid.UUID) ([]*ExpenseResponse, *repositories.PaginationResult, error) {
	// Check authorization - only platform admin can list all expenses
	requester, err := s.userRepo.GetUserByID(context.Background(), requesterID)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid requester: %v", err)
	}

	if requester.Role != "platform_admin" {
		return nil, nil, fmt.Errorf("unauthorized to list all expenses")
	}

	expenses, result, err := s.expenseRepo.ListExpenses(filters, pagination)
	if err != nil {
		return nil, nil, err
	}

	var responses []*ExpenseResponse
	for _, expense := range expenses {
		response := &ExpenseResponse{
			Expense: expense,
		}

		// Add campaign details
		campaign, err := s.campaignRepo.GetByID(context.Background(), expense.CampaignID)
		if err == nil {
			response.CampaignTitle = campaign.Title
			if campaign.NGO != nil {
				response.NGOName = campaign.NGO.Name
			}
		}

		// Add milestone details if applicable
		if expense.MilestoneID != nil {
			milestone, err := s.milestoneRepo.GetMilestoneByID(*expense.MilestoneID)
			if err == nil {
				response.MilestoneTitle = milestone.Title
			}
		}

		// Add approver details if applicable
		if expense.ApprovedBy != nil {
			approver, err := s.userRepo.GetUserByID(context.Background(), *expense.ApprovedBy)
			if err == nil {
				response.ApproverName = approver.Name
			}
		}

		responses = append(responses, response)
	}

	return responses, result, nil
}

// GetExpensesByStatus retrieves expenses by status (admin only)
func (s *ExpenseService) GetExpensesByStatus(status string, pagination repositories.PaginationParams, requesterID uuid.UUID) ([]*ExpenseResponse, *repositories.PaginationResult, error) {
	// Check authorization - only platform admin can view all expenses
	requester, err := s.userRepo.GetUserByID(context.Background(), requesterID)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid requester: %v", err)
	}

	if requester.Role != "platform_admin" {
		return nil, nil, fmt.Errorf("unauthorized to view all expenses")
	}

	expenses, result, err := s.expenseRepo.GetExpensesByStatus(status, pagination)
	if err != nil {
		return nil, nil, err
	}

	var responses []*ExpenseResponse
	for _, expense := range expenses {
		response := &ExpenseResponse{
			Expense: expense,
		}

		// Add campaign details
		campaign, err := s.campaignRepo.GetByID(context.Background(), expense.CampaignID)
		if err == nil {
			response.CampaignTitle = campaign.Title
			if campaign.NGO != nil {
				response.NGOName = campaign.NGO.Name
			}
		}

		// Add milestone details if applicable
		if expense.MilestoneID != nil {
			milestone, err := s.milestoneRepo.GetMilestoneByID(*expense.MilestoneID)
			if err == nil {
				response.MilestoneTitle = milestone.Title
			}
		}

		// Add approver details if applicable
		if expense.ApprovedBy != nil {
			approver, err := s.userRepo.GetUserByID(context.Background(), *expense.ApprovedBy)
			if err == nil {
				response.ApproverName = approver.Name
			}
		}

		responses = append(responses, response)
	}

	return responses, result, nil
}

// validateExpenseAccess checks if user can access/modify expenses for a campaign
func (s *ExpenseService) validateExpenseAccess(campaignID uuid.UUID, requesterID uuid.UUID) error {
	requester, err := s.userRepo.GetUserByID(context.Background(), requesterID)
	if err != nil {
		return fmt.Errorf("invalid requester: %v", err)
	}

	// Platform admin has access to all expenses
	if requester.Role == "platform_admin" {
		return nil
	}

	// NGO admin can access expenses for their campaigns
	if requester.Role == "ngo_admin" {
		ngo, err := s.ngoRepo.GetByUserID(context.Background(), requesterID)
		if err != nil {
			return fmt.Errorf("failed to get NGO: %v", err)
		}

		campaign, err := s.campaignRepo.GetByID(context.Background(), campaignID)
		if err != nil {
			return fmt.Errorf("invalid campaign: %v", err)
		}

		if ngo.ID == campaign.NGOID {
			return nil
		}
	}

	return fmt.Errorf("unauthorized to access expenses for this campaign")
}