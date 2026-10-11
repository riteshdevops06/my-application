package services

import (
	"context"
	"fmt"
	"time"

	"ngo-donation-tracker/models"
	"ngo-donation-tracker/repositories"

	"github.com/google/uuid"
)

// MilestoneService handles business logic for milestones
type MilestoneService struct {
	milestoneRepo repositories.MilestoneRepository
	campaignRepo  repositories.CampaignRepositoryInterface
	userRepo      repositories.UserRepositoryInterface
	ngoRepo       repositories.NGORepositoryInterface
}

// NewMilestoneService creates a new instance of MilestoneService
func NewMilestoneService(
	milestoneRepo repositories.MilestoneRepository,
	campaignRepo repositories.CampaignRepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
	ngoRepo repositories.NGORepositoryInterface,
) *MilestoneService {
	return &MilestoneService{
		milestoneRepo: milestoneRepo,
		campaignRepo:  campaignRepo,
		userRepo:      userRepo,
		ngoRepo:       ngoRepo,
	}
}

// CreateMilestoneRequest represents request to create a milestone
type CreateMilestoneRequest struct {
	CampaignID   uuid.UUID `json:"campaign_id" validate:"required"`
	Title        string    `json:"title" validate:"required,min=5,max=255"`
	Description  string    `json:"description,omitempty"`
	TargetAmount float64   `json:"target_amount" validate:"required,gt=0"`
	DueDate      string    `json:"due_date,omitempty"`
}

// UpdateMilestoneRequest represents request to update a milestone
type UpdateMilestoneRequest struct {
	Title        string  `json:"title,omitempty" validate:"omitempty,min=5,max=255"`
	Description  string  `json:"description,omitempty"`
	TargetAmount float64 `json:"target_amount,omitempty" validate:"omitempty,gt=0"`
	DueDate      string  `json:"due_date,omitempty"`
}

// CompleteMilestoneRequest represents request to complete a milestone
type CompleteMilestoneRequest struct {
	ProofDocumentURL string `json:"proof_document_url" validate:"required,url"`
	ProofNotes       string `json:"proof_notes" validate:"required,min=10"`
}

// MilestoneResponse represents milestone response with additional details
type MilestoneResponse struct {
	*models.Milestone
	CampaignTitle string `json:"campaign_title,omitempty"`
	NGOName       string `json:"ngo_name,omitempty"`
	IsOverdue     bool   `json:"is_overdue"`
}

// CreateMilestone creates a new milestone
func (s *MilestoneService) CreateMilestone(requesterID uuid.UUID, req *CreateMilestoneRequest) (*models.Milestone, error) {
	// Validate campaign exists and user has permission
	campaign, err := s.campaignRepo.GetByID(context.Background(), req.CampaignID)
	if err != nil {
		return nil, fmt.Errorf("invalid campaign: %v", err)
	}

	// Check if requester is authorized to create milestones for this campaign
	requester, err := s.userRepo.GetUserByID(context.Background(), requesterID)
	if err != nil {
		return nil, fmt.Errorf("invalid requester: %v", err)
	}

	if requester.Role != "platform_admin" {
		if requester.Role != "ngo_admin" {
			return nil, fmt.Errorf("unauthorized to create milestones")
		}

		// Check if user is admin of the NGO that owns this campaign
		ngo, err := s.ngoRepo.GetByUserID(context.Background(), requesterID)
		if err != nil || ngo.ID != campaign.NGOID {
			return nil, fmt.Errorf("unauthorized to create milestones for this campaign")
		}
	}

	// Parse due date if provided
	var dueDate *time.Time
	if req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			return nil, fmt.Errorf("invalid due date format (use YYYY-MM-DD): %v", err)
		}
		dueDate = &parsed
	}

	// Validate due date is in the future
	if dueDate != nil && dueDate.Before(time.Now()) {
		return nil, fmt.Errorf("due date must be in the future")
	}

	// Create milestone
	milestone := &models.Milestone{
		ID:           uuid.New(),
		CampaignID:   req.CampaignID,
		Title:        req.Title,
		Description:  &req.Description,
		TargetAmount: req.TargetAmount,
		RaisedAmount: 0,
		Status:       "pending",
		DueDate:      dueDate,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.milestoneRepo.CreateMilestone(milestone)
	if err != nil {
		return nil, fmt.Errorf("failed to create milestone: %v", err)
	}

	return milestone, nil
}

// UpdateMilestone updates an existing milestone
func (s *MilestoneService) UpdateMilestone(milestoneID uuid.UUID, requesterID uuid.UUID, req *UpdateMilestoneRequest) (*models.Milestone, error) {
	// Get existing milestone
	milestone, err := s.milestoneRepo.GetMilestoneByID(milestoneID)
	if err != nil {
		return nil, fmt.Errorf("milestone not found: %v", err)
	}

	// Check authorization
	err = s.validateMilestoneAccess(milestone.CampaignID, requesterID)
	if err != nil {
		return nil, err
	}

	// Prevent updates if milestone is completed
	if milestone.Status == "completed" {
		return nil, fmt.Errorf("completed milestones cannot be updated")
	}

	// Update fields if provided
	if req.Title != "" {
		milestone.Title = req.Title
	}
	if req.Description != "" {
		milestone.Description = &req.Description
	}
	if req.TargetAmount > 0 {
		milestone.TargetAmount = req.TargetAmount
	}
	if req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			return nil, fmt.Errorf("invalid due date format (use YYYY-MM-DD): %v", err)
		}
		
		// Validate due date is in the future
		if parsed.Before(time.Now()) {
			return nil, fmt.Errorf("due date must be in the future")
		}
		
		milestone.DueDate = &parsed
	}

	milestone.UpdatedAt = time.Now()

	err = s.milestoneRepo.UpdateMilestone(milestone)
	if err != nil {
		return nil, fmt.Errorf("failed to update milestone: %v", err)
	}

	return milestone, nil
}

// DeleteMilestone deletes a milestone
func (s *MilestoneService) DeleteMilestone(milestoneID uuid.UUID, requesterID uuid.UUID) error {
	// Get existing milestone
	milestone, err := s.milestoneRepo.GetMilestoneByID(milestoneID)
	if err != nil {
		return fmt.Errorf("milestone not found: %v", err)
	}

	// Check authorization
	err = s.validateMilestoneAccess(milestone.CampaignID, requesterID)
	if err != nil {
		return err
	}

	// Prevent deletion if milestone is completed
	if milestone.Status == "completed" {
		return fmt.Errorf("completed milestones cannot be deleted")
	}

	err = s.milestoneRepo.DeleteMilestone(milestoneID)
	if err != nil {
		return fmt.Errorf("failed to delete milestone: %v", err)
	}

	return nil
}

// GetMilestone retrieves a milestone with additional details
func (s *MilestoneService) GetMilestone(milestoneID uuid.UUID, requesterID uuid.UUID) (*MilestoneResponse, error) {
	milestone, err := s.milestoneRepo.GetMilestoneByID(milestoneID)
	if err != nil {
		return nil, fmt.Errorf("milestone not found: %v", err)
	}

	// Check authorization for detailed view
	requester, err := s.userRepo.GetUserByID(context.Background(), requesterID)
	if err != nil {
		return nil, fmt.Errorf("invalid requester: %v", err)
	}

	response := &MilestoneResponse{
		Milestone: milestone,
		IsOverdue: s.isMilestoneOverdue(milestone),
	}

	// Add campaign details if user has access
	campaign, err := s.campaignRepo.GetByID(context.Background(), milestone.CampaignID)
	if err == nil {
		response.CampaignTitle = campaign.Title
		
		// Add NGO name
		if campaign.NGO != nil {
			response.NGOName = campaign.NGO.Name
		}
	}

	// Only NGO admin and platform admin can see detailed milestone info
	if requester.Role == "platform_admin" {
		return response, nil
	}

	if requester.Role == "ngo_admin" {
		ngo, err := s.ngoRepo.GetByUserID(context.Background(), requesterID)
		if err == nil && campaign != nil && ngo.ID == campaign.NGOID {
			return response, nil
		}
	}

	// For donors, return limited information
	response.ProofDocumentURL = nil
	response.ProofNotes = nil

	return response, nil
}

// ListMilestonesByCampaign retrieves milestones for a campaign
func (s *MilestoneService) ListMilestonesByCampaign(campaignID uuid.UUID, requesterID uuid.UUID) ([]*MilestoneResponse, error) {
	milestones, err := s.milestoneRepo.ListMilestonesByCampaign(campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to get milestones: %v", err)
	}

	// Check if requester has detailed access
	requester, err := s.userRepo.GetUserByID(context.Background(), requesterID)
	if err != nil {
		return nil, fmt.Errorf("invalid requester: %v", err)
	}

	hasDetailedAccess := false
	if requester.Role == "platform_admin" {
		hasDetailedAccess = true
	} else if requester.Role == "ngo_admin" {
		ngo, err := s.ngoRepo.GetByUserID(context.Background(), requesterID)
		if err == nil {
			campaign, err := s.campaignRepo.GetByID(context.Background(), campaignID)
			if err == nil && ngo.ID == campaign.NGOID {
				hasDetailedAccess = true
			}
		}
	}

	var responses []*MilestoneResponse
	for _, milestone := range milestones {
		response := &MilestoneResponse{
			Milestone: milestone,
			IsOverdue: s.isMilestoneOverdue(milestone),
		}

		// Limit information for non-authorized users
		if !hasDetailedAccess {
			response.ProofDocumentURL = nil
			response.ProofNotes = nil
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// CompleteMilestone marks a milestone as complete with proof
func (s *MilestoneService) CompleteMilestone(milestoneID uuid.UUID, requesterID uuid.UUID, req *CompleteMilestoneRequest) (*models.Milestone, error) {
	// Get existing milestone
	milestone, err := s.milestoneRepo.GetMilestoneByID(milestoneID)
	if err != nil {
		return nil, fmt.Errorf("milestone not found: %v", err)
	}

	// Check authorization
	err = s.validateMilestoneAccess(milestone.CampaignID, requesterID)
	if err != nil {
		return nil, err
	}

	// Validate milestone can be completed
	if milestone.Status == "completed" {
		return nil, fmt.Errorf("milestone is already completed")
	}

	if milestone.Status == "cancelled" {
		return nil, fmt.Errorf("cancelled milestones cannot be completed")
	}

	// Mark milestone as complete
	err = s.milestoneRepo.MarkMilestoneComplete(milestoneID, req.ProofDocumentURL, req.ProofNotes)
	if err != nil {
		return nil, fmt.Errorf("failed to complete milestone: %v", err)
	}

	// Get updated milestone
	updatedMilestone, err := s.milestoneRepo.GetMilestoneByID(milestoneID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated milestone: %v", err)
	}

	// TODO: Send notification to donors about milestone completion
	// TODO: Update campaign progress if needed

	return updatedMilestone, nil
}

// UpdateMilestoneStatus updates milestone status
func (s *MilestoneService) UpdateMilestoneStatus(milestoneID uuid.UUID, status string, requesterID uuid.UUID) (*models.Milestone, error) {
	// Get existing milestone
	milestone, err := s.milestoneRepo.GetMilestoneByID(milestoneID)
	if err != nil {
		return nil, fmt.Errorf("milestone not found: %v", err)
	}

	// Check authorization
	err = s.validateMilestoneAccess(milestone.CampaignID, requesterID)
	if err != nil {
		return nil, err
	}

	// Validate status transition
	validStatuses := []string{"pending", "in_progress", "completed", "cancelled"}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}

	if !isValid {
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	// Prevent certain status changes
	if milestone.Status == "completed" && status != "completed" {
		return nil, fmt.Errorf("completed milestones cannot change status")
	}

	err = s.milestoneRepo.UpdateMilestoneStatus(milestoneID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to update milestone status: %v", err)
	}

	// Get updated milestone
	updatedMilestone, err := s.milestoneRepo.GetMilestoneByID(milestoneID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated milestone: %v", err)
	}

	return updatedMilestone, nil
}

// GetMilestoneStats retrieves milestone statistics for a campaign
func (s *MilestoneService) GetMilestoneStats(campaignID uuid.UUID, requesterID uuid.UUID) (*repositories.MilestoneStats, error) {
	// Check authorization
	err := s.validateMilestoneAccess(campaignID, requesterID)
	if err != nil {
		return nil, err
	}

	return s.milestoneRepo.GetMilestoneStats(campaignID)
}

// GetMilestonesByStatus retrieves milestones by status (admin only)
func (s *MilestoneService) GetMilestonesByStatus(status string, pagination repositories.PaginationParams, requesterID uuid.UUID) ([]*MilestoneResponse, *repositories.PaginationResult, error) {
	// Check authorization - only platform admin can view all milestones
	requester, err := s.userRepo.GetUserByID(context.Background(), requesterID)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid requester: %v", err)
	}

	if requester.Role != "platform_admin" {
		return nil, nil, fmt.Errorf("unauthorized to view all milestones")
	}

	milestones, result, err := s.milestoneRepo.GetMilestonesByStatus(status, pagination)
	if err != nil {
		return nil, nil, err
	}

	var responses []*MilestoneResponse
	for _, milestone := range milestones {
		response := &MilestoneResponse{
			Milestone: milestone,
			IsOverdue: s.isMilestoneOverdue(milestone),
		}

		// Add campaign details
		campaign, err := s.campaignRepo.GetByID(context.Background(), milestone.CampaignID)
		if err == nil {
			response.CampaignTitle = campaign.Title
			if campaign.NGO != nil {
				response.NGOName = campaign.NGO.Name
			}
		}

		responses = append(responses, response)
	}

	return responses, result, nil
}

// validateMilestoneAccess checks if user can access/modify milestone
func (s *MilestoneService) validateMilestoneAccess(campaignID uuid.UUID, requesterID uuid.UUID) error {
	requester, err := s.userRepo.GetUserByID(context.Background(), requesterID)
	if err != nil {
		return fmt.Errorf("invalid requester: %v", err)
	}

	// Platform admin has access to all milestones
	if requester.Role == "platform_admin" {
		return nil
	}

	// NGO admin can access milestones for their campaigns
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

	return fmt.Errorf("unauthorized to access this milestone")
}

// isMilestoneOverdue checks if a milestone is overdue
func (s *MilestoneService) isMilestoneOverdue(milestone *models.Milestone) bool {
	if milestone.DueDate == nil || milestone.Status == "completed" {
		return false
	}
	return time.Now().After(*milestone.DueDate)
}