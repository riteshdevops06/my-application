package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"ngo-donation-tracker/models"
	"ngo-donation-tracker/repositories"

	"github.com/google/uuid"
)

// CampaignService handles campaign business logic
type CampaignService struct {
	campaignRepo repositories.CampaignRepositoryInterface
	ngoRepo      repositories.NGORepositoryInterface
	userRepo     repositories.UserRepositoryInterface
}

// NewCampaignService creates a new campaign service
func NewCampaignService(
	campaignRepo repositories.CampaignRepositoryInterface,
	ngoRepo repositories.NGORepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
) *CampaignService {
	return &CampaignService{
		campaignRepo: campaignRepo,
		ngoRepo:      ngoRepo,
		userRepo:     userRepo,
	}
}

// CreateCampaign creates a new campaign for an NGO
func (s *CampaignService) CreateCampaign(ctx context.Context, userID uuid.UUID, req *models.CampaignCreate) (*models.Campaign, error) {
	// Validate input
	if err := s.validateCampaignInput(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Get user and verify they are NGO admin
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.Role != "ngo_admin" {
		return nil, fmt.Errorf("only NGO administrators can create campaigns")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("user account is not active")
	}

	// Get NGO profile
	ngo, err := s.ngoRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get NGO profile: %w", err)
	}
	if ngo == nil {
		return nil, fmt.Errorf("NGO profile not found. Please create your NGO profile first")
	}

	// Verify NGO is verified
	if ngo.VerificationStatus != "verified" {
		return nil, fmt.Errorf("NGO must be verified before creating campaigns")
	}

	// Create campaign using enhanced repository
	if enhancedRepo, ok := s.campaignRepo.(*repositories.CampaignRepository); ok {
		return enhancedRepo.CreateCampaign(ctx, ngo.ID, req)
	}

	// Fallback to interface method
	campaign := &models.Campaign{
		NGOID:            ngo.ID,
		Title:            req.Title,
		Description:      &req.Description,
		ShortDescription: &req.ShortDescription,
		GoalAmount:       req.GoalAmount,
		RaisedAmount:     0.0,
		Currency:         req.Currency,
		Status:           "draft",
		Category:         &req.Category,
		Tags:             req.Tags,
		ImageURL:         &req.ImageURL,
		GalleryImages:    req.GalleryImages,
		Location:         &req.Location,
		IsFeatured:       false,
	}

	err = s.campaignRepo.Create(ctx, campaign)
	if err != nil {
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	return campaign, nil
}

// GetCampaign retrieves a campaign by ID
func (s *CampaignService) GetCampaign(ctx context.Context, id uuid.UUID) (*models.Campaign, error) {
	campaign, err := s.campaignRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrCampaignNotFound) {
			return nil, fmt.Errorf("campaign not found")
		}
		return nil, fmt.Errorf("failed to get campaign: %w", err)
	}

	return campaign, nil
}

// GetCampaignWithNGO retrieves a campaign with its NGO information
func (s *CampaignService) GetCampaignWithNGO(ctx context.Context, id uuid.UUID) (map[string]interface{}, error) {
	campaign, err := s.GetCampaign(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get NGO information
	ngo, err := s.ngoRepo.GetByID(ctx, campaign.NGOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get NGO information: %w", err)
	}

	response := map[string]interface{}{
		"campaign": campaign,
		"ngo":      ngo,
	}

	return response, nil
}

// UpdateCampaign updates a campaign (NGO admin only)
func (s *CampaignService) UpdateCampaign(ctx context.Context, userID, campaignID uuid.UUID, req *models.CampaignUpdate) (*models.Campaign, error) {
	// Validate input
	if err := s.validateCampaignUpdateInput(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Verify ownership
	if err := s.verifyCampaignOwnership(ctx, userID, campaignID); err != nil {
		return nil, err
	}

	// Get existing campaign to check status
	existingCampaign, err := s.campaignRepo.GetByID(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign: %w", err)
	}

	// Check if campaign can be modified
	if existingCampaign.Status == "completed" || existingCampaign.Status == "cancelled" {
		return nil, fmt.Errorf("cannot modify completed or cancelled campaigns")
	}

	// Update campaign using enhanced repository
	if enhancedRepo, ok := s.campaignRepo.(*repositories.CampaignRepository); ok {
		return enhancedRepo.UpdateCampaign(ctx, campaignID, req)
	}

	// Fallback to basic update (limited functionality)
	existingCampaign.Title = req.Title
	existingCampaign.Description = &req.Description
	existingCampaign.ShortDescription = &req.ShortDescription
	if req.GoalAmount > 0 {
		existingCampaign.GoalAmount = req.GoalAmount
	}
	existingCampaign.Category = &req.Category

	err = s.campaignRepo.Update(ctx, existingCampaign)
	if err != nil {
		return nil, fmt.Errorf("failed to update campaign: %w", err)
	}

	return existingCampaign, nil
}

// DeleteCampaign soft deletes a campaign (NGO admin only)
func (s *CampaignService) DeleteCampaign(ctx context.Context, userID, campaignID uuid.UUID) error {
	// Verify ownership
	if err := s.verifyCampaignOwnership(ctx, userID, campaignID); err != nil {
		return err
	}

	// Get existing campaign to check if it can be deleted
	campaign, err := s.campaignRepo.GetByID(ctx, campaignID)
	if err != nil {
		return fmt.Errorf("failed to get campaign: %w", err)
	}

	// Only allow deletion of draft campaigns or campaigns with no donations
	if campaign.Status != "draft" && campaign.RaisedAmount > 0 {
		return fmt.Errorf("cannot delete campaigns that have received donations")
	}

	// Delete campaign
	err = s.campaignRepo.Delete(ctx, campaignID)
	if err != nil {
		return fmt.Errorf("failed to delete campaign: %w", err)
	}

	return nil
}

// PublishCampaign changes campaign status from draft to active
func (s *CampaignService) PublishCampaign(ctx context.Context, userID, campaignID uuid.UUID) error {
	// Verify ownership
	if err := s.verifyCampaignOwnership(ctx, userID, campaignID); err != nil {
		return err
	}

	// Get campaign to validate it's ready for publishing
	campaign, err := s.campaignRepo.GetByID(ctx, campaignID)
	if err != nil {
		return fmt.Errorf("failed to get campaign: %w", err)
	}

	// Validate campaign is ready for publishing
	if err := s.validateCampaignForPublishing(campaign); err != nil {
		return fmt.Errorf("campaign cannot be published: %w", err)
	}

	// Update status to active
	err = s.campaignRepo.UpdateStatus(ctx, campaignID, "active")
	if err != nil {
		return fmt.Errorf("failed to publish campaign: %w", err)
	}

	return nil
}

// PauseCampaign pauses an active campaign
func (s *CampaignService) PauseCampaign(ctx context.Context, userID, campaignID uuid.UUID) error {
	// Verify ownership
	if err := s.verifyCampaignOwnership(ctx, userID, campaignID); err != nil {
		return err
	}

	// Get campaign to check current status
	campaign, err := s.campaignRepo.GetByID(ctx, campaignID)
	if err != nil {
		return fmt.Errorf("failed to get campaign: %w", err)
	}

	if campaign.Status != "active" {
		return fmt.Errorf("only active campaigns can be paused")
	}

	// Update status to paused
	err = s.campaignRepo.UpdateStatus(ctx, campaignID, "paused")
	if err != nil {
		return fmt.Errorf("failed to pause campaign: %w", err)
	}

	return nil
}

// ResumeCampaign resumes a paused campaign
func (s *CampaignService) ResumeCampaign(ctx context.Context, userID, campaignID uuid.UUID) error {
	// Verify ownership
	if err := s.verifyCampaignOwnership(ctx, userID, campaignID); err != nil {
		return err
	}

	// Get campaign to check current status
	campaign, err := s.campaignRepo.GetByID(ctx, campaignID)
	if err != nil {
		return fmt.Errorf("failed to get campaign: %w", err)
	}

	if campaign.Status != "paused" {
		return fmt.Errorf("only paused campaigns can be resumed")
	}

	// Update status to active
	err = s.campaignRepo.UpdateStatus(ctx, campaignID, "active")
	if err != nil {
		return fmt.Errorf("failed to resume campaign: %w", err)
	}

	return nil
}

// ListCampaigns retrieves paginated list of campaigns with filters
func (s *CampaignService) ListCampaigns(ctx context.Context, params models.Pagination, filters repositories.CampaignFilters) ([]*models.Campaign, int64, error) {
	if enhancedRepo, ok := s.campaignRepo.(*repositories.CampaignRepository); ok {
		return enhancedRepo.ListCampaigns(ctx, params, filters)
	}

	// Fallback to interface method with limited filtering
	return s.campaignRepo.List(ctx, params, filters)
}

// ListNGOCampaigns retrieves campaigns for a specific NGO
func (s *CampaignService) ListNGOCampaigns(ctx context.Context, userID uuid.UUID, params models.Pagination) ([]*models.Campaign, int64, error) {
	// Get NGO profile
	ngo, err := s.ngoRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get NGO profile: %w", err)
	}
	if ngo == nil {
		return nil, 0, fmt.Errorf("NGO profile not found")
	}

	return s.campaignRepo.ListByNGO(ctx, ngo.ID, params)
}

// GetActiveCampaigns returns all active campaigns
func (s *CampaignService) GetActiveCampaigns(ctx context.Context, params models.Pagination) ([]*models.Campaign, int64, error) {
	if enhancedRepo, ok := s.campaignRepo.(*repositories.CampaignRepository); ok {
		return enhancedRepo.GetActiveCampaigns(ctx, params)
	}

	// Fallback using filters
	filters := repositories.CampaignFilters{Status: "active"}
	return s.campaignRepo.List(ctx, params, filters)
}

// GetFeaturedCampaigns returns all featured campaigns
func (s *CampaignService) GetFeaturedCampaigns(ctx context.Context, params models.Pagination) ([]*models.Campaign, int64, error) {
	if enhancedRepo, ok := s.campaignRepo.(*repositories.CampaignRepository); ok {
		return enhancedRepo.GetFeaturedCampaigns(ctx, params)
	}

	// Fallback - return active campaigns (basic implementation)
	filters := repositories.CampaignFilters{Status: "active"}
	return s.campaignRepo.List(ctx, params, filters)
}

// SearchCampaigns searches campaigns by title and description
func (s *CampaignService) SearchCampaigns(ctx context.Context, searchTerm string, params models.Pagination) ([]*models.Campaign, int64, error) {
	if searchTerm = strings.TrimSpace(searchTerm); searchTerm == "" {
		return nil, 0, fmt.Errorf("search term is required")
	}

	if enhancedRepo, ok := s.campaignRepo.(*repositories.CampaignRepository); ok {
		return enhancedRepo.SearchCampaigns(ctx, searchTerm, params)
	}

	// Fallback using filters
	filters := repositories.CampaignFilters{SearchTerm: searchTerm}
	return s.campaignRepo.List(ctx, params, filters)
}

// GetCampaignsByCategory returns campaigns in a specific category
func (s *CampaignService) GetCampaignsByCategory(ctx context.Context, category string, params models.Pagination) ([]*models.Campaign, int64, error) {
	if category = strings.TrimSpace(category); category == "" {
		return nil, 0, fmt.Errorf("category is required")
	}

	if enhancedRepo, ok := s.campaignRepo.(*repositories.CampaignRepository); ok {
		return enhancedRepo.GetCampaignsByCategory(ctx, category, params)
	}

	// Fallback using filters
	filters := repositories.CampaignFilters{Category: category}
	return s.campaignRepo.List(ctx, params, filters)
}
// Admin Methods

// SetFeaturedStatus sets or unsets a campaign as featured (platform admin only)
func (s *CampaignService) SetFeaturedStatus(ctx context.Context, adminUserID, campaignID uuid.UUID, isFeatured bool) error {
	// Verify user is platform admin
	user, err := s.userRepo.GetUserByID(ctx, adminUserID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.Role != "platform_admin" {
		return fmt.Errorf("only platform administrators can set featured status")
	}

	// Verify campaign exists and is active
	campaign, err := s.campaignRepo.GetByID(ctx, campaignID)
	if err != nil {
		return fmt.Errorf("failed to get campaign: %w", err)
	}

	if campaign.Status != "active" {
		return fmt.Errorf("only active campaigns can be featured")
	}

	// Set featured status using enhanced repository
	if enhancedRepo, ok := s.campaignRepo.(*repositories.CampaignRepository); ok {
		return enhancedRepo.SetFeaturedStatus(ctx, campaignID, isFeatured)
	}

	return fmt.Errorf("featured status management not available with current repository")
}

// AdminSuspendCampaign suspends a campaign (platform admin only)
func (s *CampaignService) AdminSuspendCampaign(ctx context.Context, adminUserID, campaignID uuid.UUID, reason string) error {
	// Verify user is platform admin
	user, err := s.userRepo.GetUserByID(ctx, adminUserID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.Role != "platform_admin" {
		return fmt.Errorf("only platform administrators can suspend campaigns")
	}

	if strings.TrimSpace(reason) == "" {
		return fmt.Errorf("suspension reason is required")
	}

	// Update campaign status to cancelled
	err = s.campaignRepo.UpdateStatus(ctx, campaignID, "cancelled")
	if err != nil {
		return fmt.Errorf("failed to suspend campaign: %w", err)
	}

	// TODO: In a production system, you might want to:
	// 1. Store the suspension reason in a separate table
	// 2. Send notification to the NGO
	// 3. Log the admin action

	return nil
}

// GetCampaignStats returns campaign statistics
func (s *CampaignService) GetCampaignStats(ctx context.Context) (*repositories.CampaignStats, error) {
	if enhancedRepo, ok := s.campaignRepo.(*repositories.CampaignRepository); ok {
		return enhancedRepo.GetCampaignStats(ctx)
	}

	// Fallback implementation with basic stats
	_, activeTotal, err := s.campaignRepo.List(ctx, models.Pagination{Limit: 1}, repositories.CampaignFilters{Status: "active"})
	if err != nil {
		return nil, fmt.Errorf("failed to get active campaigns count: %w", err)
	}

	_, draftTotal, err := s.campaignRepo.List(ctx, models.Pagination{Limit: 1}, repositories.CampaignFilters{Status: "draft"})
	if err != nil {
		return nil, fmt.Errorf("failed to get draft campaigns count: %w", err)
	}

	_, completedTotal, err := s.campaignRepo.List(ctx, models.Pagination{Limit: 1}, repositories.CampaignFilters{Status: "completed"})
	if err != nil {
		return nil, fmt.Errorf("failed to get completed campaigns count: %w", err)
	}

	_, total, err := s.campaignRepo.List(ctx, models.Pagination{Limit: 1}, repositories.CampaignFilters{})
	if err != nil {
		return nil, fmt.Errorf("failed to get total campaigns count: %w", err)
	}

	return &repositories.CampaignStats{
		TotalCampaigns:     int(total),
		ActiveCampaigns:    int(activeTotal),
		DraftCampaigns:     int(draftTotal),
		CompletedCampaigns: int(completedTotal),
		FeaturedCampaigns:  0, // Cannot calculate without enhanced repository
		TotalGoalAmount:    0, // Cannot calculate without enhanced repository
		TotalRaisedAmount:  0, // Cannot calculate without enhanced repository
		SuccessRate:        0, // Cannot calculate without enhanced repository
		NewCampaigns30Days: 0, // Cannot calculate without enhanced repository
	}, nil
}

// Helper and validation methods

// verifyCampaignOwnership verifies that a user owns a campaign
func (s *CampaignService) verifyCampaignOwnership(ctx context.Context, userID, campaignID uuid.UUID) error {
	// Get user and verify they are NGO admin
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.Role != "ngo_admin" {
		return fmt.Errorf("only NGO administrators can manage campaigns")
	}

	// Get NGO profile
	ngo, err := s.ngoRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get NGO profile: %w", err)
	}
	if ngo == nil {
		return fmt.Errorf("NGO profile not found")
	}

	// Get campaign and verify ownership
	campaign, err := s.campaignRepo.GetByID(ctx, campaignID)
	if err != nil {
		if errors.Is(err, repositories.ErrCampaignNotFound) {
			return fmt.Errorf("campaign not found")
		}
		return fmt.Errorf("failed to get campaign: %w", err)
	}

	if campaign.NGOID != ngo.ID {
		return fmt.Errorf("campaign does not belong to your NGO")
	}

	return nil
}

// validateCampaignInput validates campaign creation input
func (s *CampaignService) validateCampaignInput(req *models.CampaignCreate) error {
	if req.Title == "" {
		return errors.New("campaign title is required")
	}

	if len(strings.TrimSpace(req.Title)) < 5 {
		return errors.New("campaign title must be at least 5 characters long")
	}

	if len(req.Title) > 255 {
		return errors.New("campaign title must not exceed 255 characters")
	}

	if req.GoalAmount <= 0 {
		return errors.New("goal amount must be greater than 0")
	}

	if req.GoalAmount > 10000000 { // 10 million limit
		return errors.New("goal amount cannot exceed 10,000,000")
	}

	// Validate category if provided
	if req.Category != "" && len(req.Category) > 100 {
		return errors.New("category must not exceed 100 characters")
	}

	// Validate location if provided
	if req.Location != "" && len(req.Location) > 255 {
		return errors.New("location must not exceed 255 characters")
	}

	return nil
}

// validateCampaignUpdateInput validates campaign update input
func (s *CampaignService) validateCampaignUpdateInput(req *models.CampaignUpdate) error {
	if req.Title != "" {
		if len(strings.TrimSpace(req.Title)) < 5 {
			return errors.New("campaign title must be at least 5 characters long")
		}
		if len(req.Title) > 255 {
			return errors.New("campaign title must not exceed 255 characters")
		}
	}

	if req.GoalAmount > 0 {
		if req.GoalAmount > 10000000 { // 10 million limit
			return errors.New("goal amount cannot exceed 10,000,000")
		}
	}

	if req.Category != "" && len(req.Category) > 100 {
		return errors.New("category must not exceed 100 characters")
	}

	if req.Location != "" && len(req.Location) > 255 {
		return errors.New("location must not exceed 255 characters")
	}

	return nil
}

// validateCampaignForPublishing validates that a campaign is ready to be published
func (s *CampaignService) validateCampaignForPublishing(campaign *models.Campaign) error {
	if campaign.Status != "draft" {
		return errors.New("only draft campaigns can be published")
	}

	if campaign.Title == "" || len(strings.TrimSpace(campaign.Title)) < 5 {
		return errors.New("campaign must have a valid title")
	}

	if campaign.Description == nil || *campaign.Description == "" {
		return errors.New("campaign must have a description")
	}

	if campaign.GoalAmount <= 0 {
		return errors.New("campaign must have a valid goal amount")
	}

	// Optional: Check if NGO is still verified
	// This would require additional repository call

	return nil
}

// GetAvailableCategories returns available campaign categories
func (s *CampaignService) GetAvailableCategories() []map[string]string {
	return []map[string]string{
		{"value": "education", "label": "Education"},
		{"value": "healthcare", "label": "Healthcare"},
		{"value": "environment", "label": "Environment"},
		{"value": "poverty", "label": "Poverty Alleviation"},
		{"value": "disaster_relief", "label": "Disaster Relief"},
		{"value": "animal_welfare", "label": "Animal Welfare"},
		{"value": "human_rights", "label": "Human Rights"},
		{"value": "community_development", "label": "Community Development"},
		{"value": "technology", "label": "Technology for Good"},
		{"value": "arts_culture", "label": "Arts & Culture"},
		{"value": "other", "label": "Other"},
	}
}

// GetCampaignStatusOptions returns available campaign status options
func (s *CampaignService) GetCampaignStatusOptions() []map[string]string {
	return []map[string]string{
		{"value": "draft", "label": "Draft", "description": "Campaign is being prepared"},
		{"value": "active", "label": "Active", "description": "Campaign is live and accepting donations"},
		{"value": "paused", "label": "Paused", "description": "Campaign is temporarily paused"},
		{"value": "completed", "label": "Completed", "description": "Campaign has reached its goal or ended successfully"},
		{"value": "cancelled", "label": "Cancelled", "description": "Campaign has been cancelled"},
	}
}