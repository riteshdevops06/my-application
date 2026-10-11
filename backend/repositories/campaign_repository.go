package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"ngo-donation-tracker/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Custom errors for campaign operations
var (
	ErrCampaignNotFound    = errors.New("campaign not found")
	ErrCampaignExists      = errors.New("campaign with this title already exists for this NGO")
	ErrInvalidCampaignStatus = errors.New("invalid campaign status")
	ErrCampaignNotActive   = errors.New("campaign is not active")
	ErrGoalAmountInvalid   = errors.New("goal amount must be greater than 0")
)

// CampaignRepository handles all campaign-related database operations
type CampaignRepository struct {
	db *sql.DB
}

// NewCampaignRepository creates a new campaign repository instance
func NewCampaignRepository(db *sql.DB) *CampaignRepository {
	return &CampaignRepository{db: db}
}

// CreateCampaign creates a new campaign
func (r *CampaignRepository) CreateCampaign(ctx context.Context, ngoID uuid.UUID, campaignCreate *models.CampaignCreate) (*models.Campaign, error) {
	// Validate input
	if err := r.validateCampaignInput(campaignCreate); err != nil {
		return nil, err
	}

	// Check for duplicate campaign title within the same NGO
	exists, err := r.campaignTitleExists(ctx, ngoID, campaignCreate.Title, uuid.Nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check campaign title: %w", err)
	}
	if exists {
		return nil, ErrCampaignExists
	}

	campaign := &models.Campaign{
		NGOID:            ngoID,
		Title:            strings.TrimSpace(campaignCreate.Title),
		Description:      stringPtrOrNil(strings.TrimSpace(campaignCreate.Description)),
		ShortDescription: stringPtrOrNil(strings.TrimSpace(campaignCreate.ShortDescription)),
		GoalAmount:       campaignCreate.GoalAmount,
		RaisedAmount:     0.0, // Always start at 0
		Currency:         getDefaultCurrency(campaignCreate.Currency),
		Status:           "draft", // Always start as draft
		Category:         stringPtrOrNil(campaignCreate.Category),
		Tags:             campaignCreate.Tags,
		ImageURL:         stringPtrOrNil(campaignCreate.ImageURL),
		GalleryImages:    campaignCreate.GalleryImages,
		Location:         stringPtrOrNil(strings.TrimSpace(campaignCreate.Location)),
		IsFeatured:       false, // Only admins can set featured
	}

	// Parse dates if provided
	if campaignCreate.StartDate != "" {
		if startDate, err := time.Parse("2006-01-02", campaignCreate.StartDate); err == nil {
			campaign.StartDate = &startDate
		}
	}
	if campaignCreate.EndDate != "" {
		if endDate, err := time.Parse("2006-01-02", campaignCreate.EndDate); err == nil {
			campaign.EndDate = &endDate
		}
	}

	// Validate date logic
	if err := r.validateCampaignDates(campaign.StartDate, campaign.EndDate); err != nil {
		return nil, err
	}

	query := `
		INSERT INTO campaigns (ngo_id, title, description, short_description, goal_amount, raised_amount,
			currency, status, category, tags, start_date, end_date, image_url, gallery_images, 
			location, is_featured)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, created_at, updated_at
	`

	err = r.db.QueryRowContext(ctx, query,
		campaign.NGOID, campaign.Title, campaign.Description, campaign.ShortDescription,
		campaign.GoalAmount, campaign.RaisedAmount, campaign.Currency, campaign.Status,
		campaign.Category, campaign.Tags, campaign.StartDate, campaign.EndDate,
		campaign.ImageURL, campaign.GalleryImages, campaign.Location, campaign.IsFeatured,
	).Scan(&campaign.ID, &campaign.CreatedAt, &campaign.UpdatedAt)

	if err != nil {
		// Handle unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, ErrCampaignExists
		}
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	return campaign, nil
}

// GetCampaignByID retrieves a campaign by its ID
func (r *CampaignRepository) GetCampaignByID(ctx context.Context, id uuid.UUID) (*models.Campaign, error) {
	query := `
		SELECT id, ngo_id, title, description, short_description, goal_amount, raised_amount,
			currency, status, category, tags, start_date, end_date, image_url, gallery_images,
			location, is_featured, created_at, updated_at
		FROM campaigns WHERE id = $1
	`

	campaign := &models.Campaign{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&campaign.ID, &campaign.NGOID, &campaign.Title, &campaign.Description,
		&campaign.ShortDescription, &campaign.GoalAmount, &campaign.RaisedAmount,
		&campaign.Currency, &campaign.Status, &campaign.Category, &campaign.Tags,
		&campaign.StartDate, &campaign.EndDate, &campaign.ImageURL, &campaign.GalleryImages,
		&campaign.Location, &campaign.IsFeatured, &campaign.CreatedAt, &campaign.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrCampaignNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign by ID: %w", err)
	}

	return campaign, nil
}

// UpdateCampaign updates a campaign
func (r *CampaignRepository) UpdateCampaign(ctx context.Context, id uuid.UUID, campaignUpdate *models.CampaignUpdate) (*models.Campaign, error) {
	// Get existing campaign
	existingCampaign, err := r.GetCampaignByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Build dynamic query based on fields to update
	setParts := []string{}
	args := []interface{}{id}
	argIndex := 2

	if campaignUpdate.Title != "" {
		// Check for duplicate title
		exists, err := r.campaignTitleExists(ctx, existingCampaign.NGOID, campaignUpdate.Title, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check campaign title: %w", err)
		}
		if exists {
			return nil, ErrCampaignExists
		}
		
		setParts = append(setParts, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, strings.TrimSpace(campaignUpdate.Title))
		argIndex++
	}

	if campaignUpdate.Description != "" {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, stringPtrOrNil(strings.TrimSpace(campaignUpdate.Description)))
		argIndex++
	}

	if campaignUpdate.ShortDescription != "" {
		setParts = append(setParts, fmt.Sprintf("short_description = $%d", argIndex))
		args = append(args, stringPtrOrNil(strings.TrimSpace(campaignUpdate.ShortDescription)))
		argIndex++
	}

	if campaignUpdate.GoalAmount > 0 {
		setParts = append(setParts, fmt.Sprintf("goal_amount = $%d", argIndex))
		args = append(args, campaignUpdate.GoalAmount)
		argIndex++
	}

	if campaignUpdate.Category != "" {
		setParts = append(setParts, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, stringPtrOrNil(campaignUpdate.Category))
		argIndex++
	}

	if campaignUpdate.Tags != nil {
		setParts = append(setParts, fmt.Sprintf("tags = $%d", argIndex))
		args = append(args, campaignUpdate.Tags)
		argIndex++
	}

	if campaignUpdate.StartDate != "" {
		if startDate, err := time.Parse("2006-01-02", campaignUpdate.StartDate); err == nil {
			setParts = append(setParts, fmt.Sprintf("start_date = $%d", argIndex))
			args = append(args, &startDate)
			argIndex++
		}
	}

	if campaignUpdate.EndDate != "" {
		if endDate, err := time.Parse("2006-01-02", campaignUpdate.EndDate); err == nil {
			setParts = append(setParts, fmt.Sprintf("end_date = $%d", argIndex))
			args = append(args, &endDate)
			argIndex++
		}
	}

	if campaignUpdate.ImageURL != "" {
		setParts = append(setParts, fmt.Sprintf("image_url = $%d", argIndex))
		args = append(args, stringPtrOrNil(campaignUpdate.ImageURL))
		argIndex++
	}

	if campaignUpdate.GalleryImages != nil {
		setParts = append(setParts, fmt.Sprintf("gallery_images = $%d", argIndex))
		args = append(args, campaignUpdate.GalleryImages)
		argIndex++
	}

	if campaignUpdate.Location != "" {
		setParts = append(setParts, fmt.Sprintf("location = $%d", argIndex))
		args = append(args, stringPtrOrNil(strings.TrimSpace(campaignUpdate.Location)))
		argIndex++
	}

	if len(setParts) == 0 {
		// No fields to update, just return current campaign
		return existingCampaign, nil
	}

	setParts = append(setParts, "updated_at = NOW()")

	query := fmt.Sprintf(`
		UPDATE campaigns SET %s
		WHERE id = $1
		RETURNING id, ngo_id, title, description, short_description, goal_amount, raised_amount,
			currency, status, category, tags, start_date, end_date, image_url, gallery_images,
			location, is_featured, created_at, updated_at
	`, strings.Join(setParts, ", "))

	campaign := &models.Campaign{}
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&campaign.ID, &campaign.NGOID, &campaign.Title, &campaign.Description,
		&campaign.ShortDescription, &campaign.GoalAmount, &campaign.RaisedAmount,
		&campaign.Currency, &campaign.Status, &campaign.Category, &campaign.Tags,
		&campaign.StartDate, &campaign.EndDate, &campaign.ImageURL, &campaign.GalleryImages,
		&campaign.Location, &campaign.IsFeatured, &campaign.CreatedAt, &campaign.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrCampaignNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update campaign: %w", err)
	}

	return campaign, nil
}
// DeleteCampaign soft deletes a campaign (changes status to deleted)
func (r *CampaignRepository) DeleteCampaign(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE campaigns SET status = 'deleted', updated_at = NOW()
		WHERE id = $1 AND status != 'deleted'
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete campaign: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrCampaignNotFound
	}

	return nil
}

// ListCampaigns retrieves paginated list of campaigns with filters
func (r *CampaignRepository) ListCampaigns(ctx context.Context, params models.Pagination, filters CampaignFilters) ([]*models.Campaign, int64, error) {
	// Build WHERE clause
	whereClause, filterArgs := r.buildCampaignFilters(filters)

	// Get total count
	countQuery := "SELECT COUNT(*) FROM campaigns" + whereClause
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, filterArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count campaigns: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, ngo_id, title, description, short_description, goal_amount, raised_amount,
			currency, status, category, tags, start_date, end_date, image_url, gallery_images,
			location, is_featured, created_at, updated_at
		FROM campaigns` + whereClause + `
		ORDER BY created_at DESC
		LIMIT $` + fmt.Sprintf("%d", len(filterArgs)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(filterArgs)+2)

	allArgs := append(filterArgs, params.Limit, params.Offset)

	rows, err := r.db.QueryContext(ctx, query, allArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list campaigns: %w", err)
	}
	defer rows.Close()

	var campaigns []*models.Campaign
	for rows.Next() {
		campaign := &models.Campaign{}
		err := rows.Scan(
			&campaign.ID, &campaign.NGOID, &campaign.Title, &campaign.Description,
			&campaign.ShortDescription, &campaign.GoalAmount, &campaign.RaisedAmount,
			&campaign.Currency, &campaign.Status, &campaign.Category, &campaign.Tags,
			&campaign.StartDate, &campaign.EndDate, &campaign.ImageURL, &campaign.GalleryImages,
			&campaign.Location, &campaign.IsFeatured, &campaign.CreatedAt, &campaign.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan campaign: %w", err)
		}
		campaigns = append(campaigns, campaign)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating campaigns: %w", err)
	}

	return campaigns, total, nil
}

// ListCampaignsByNGO retrieves campaigns for a specific NGO
func (r *CampaignRepository) ListCampaignsByNGO(ctx context.Context, ngoID uuid.UUID, params models.Pagination) ([]*models.Campaign, int64, error) {
	filters := CampaignFilters{
		NGOID: &ngoID,
	}
	return r.ListCampaigns(ctx, params, filters)
}

// UpdateCampaignStatus updates only the campaign status
func (r *CampaignRepository) UpdateCampaignStatus(ctx context.Context, id uuid.UUID, status string) error {
	// Validate status
	validStatuses := []string{"draft", "active", "paused", "completed", "cancelled", "deleted"}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}
	if !isValid {
		return ErrInvalidCampaignStatus
	}

	query := `
		UPDATE campaigns SET status = $2, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update campaign status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrCampaignNotFound
	}

	return nil
}

// UpdateRaisedAmount updates the raised amount for a campaign
func (r *CampaignRepository) UpdateRaisedAmount(ctx context.Context, id uuid.UUID, amount float64) error {
	if amount < 0 {
		return fmt.Errorf("raised amount cannot be negative")
	}

	query := `
		UPDATE campaigns SET raised_amount = $2, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, amount)
	if err != nil {
		return fmt.Errorf("failed to update raised amount: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrCampaignNotFound
	}

	return nil
}

// IncrementRaisedAmount adds to the current raised amount
func (r *CampaignRepository) IncrementRaisedAmount(ctx context.Context, id uuid.UUID, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("increment amount must be positive")
	}

	query := `
		UPDATE campaigns SET raised_amount = raised_amount + $2, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, amount)
	if err != nil {
		return fmt.Errorf("failed to increment raised amount: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrCampaignNotFound
	}

	return nil
}

// GetActiveCampaigns returns all active campaigns
func (r *CampaignRepository) GetActiveCampaigns(ctx context.Context, params models.Pagination) ([]*models.Campaign, int64, error) {
	filters := CampaignFilters{
		Status: "active",
	}
	return r.ListCampaigns(ctx, params, filters)
}

// GetFeaturedCampaigns returns all featured campaigns
func (r *CampaignRepository) GetFeaturedCampaigns(ctx context.Context, params models.Pagination) ([]*models.Campaign, int64, error) {
	isFeatured := true
	filters := CampaignFilters{
		IsFeatured: &isFeatured,
		Status:     "active", // Only show active featured campaigns
	}
	return r.ListCampaigns(ctx, params, filters)
}

// SetFeaturedStatus sets or unsets a campaign as featured (admin only)
func (r *CampaignRepository) SetFeaturedStatus(ctx context.Context, id uuid.UUID, isFeatured bool) error {
	query := `
		UPDATE campaigns SET is_featured = $2, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, isFeatured)
	if err != nil {
		return fmt.Errorf("failed to update featured status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrCampaignNotFound
	}

	return nil
}

// SearchCampaigns searches campaigns by title and description
func (r *CampaignRepository) SearchCampaigns(ctx context.Context, searchTerm string, params models.Pagination) ([]*models.Campaign, int64, error) {
	filters := CampaignFilters{
		SearchTerm: searchTerm,
	}
	return r.ListCampaigns(ctx, params, filters)
}

// GetCampaignsByCategory returns campaigns in a specific category
func (r *CampaignRepository) GetCampaignsByCategory(ctx context.Context, category string, params models.Pagination) ([]*models.Campaign, int64, error) {
	filters := CampaignFilters{
		Category: category,
		Status:   "active", // Only show active campaigns
	}
	return r.ListCampaigns(ctx, params, filters)
}

// GetCampaignStats returns campaign statistics
func (r *CampaignRepository) GetCampaignStats(ctx context.Context) (*CampaignStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_campaigns,
			COUNT(CASE WHEN status = 'active' THEN 1 END) as active_campaigns,
			COUNT(CASE WHEN status = 'draft' THEN 1 END) as draft_campaigns,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_campaigns,
			COUNT(CASE WHEN is_featured = true THEN 1 END) as featured_campaigns,
			COALESCE(SUM(goal_amount), 0) as total_goal_amount,
			COALESCE(SUM(raised_amount), 0) as total_raised_amount,
			COUNT(CASE WHEN created_at >= NOW() - INTERVAL '30 days' THEN 1 END) as new_campaigns_30d
		FROM campaigns
		WHERE status != 'deleted'
	`

	stats := &CampaignStats{}
	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalCampaigns, &stats.ActiveCampaigns, &stats.DraftCampaigns,
		&stats.CompletedCampaigns, &stats.FeaturedCampaigns, &stats.TotalGoalAmount,
		&stats.TotalRaisedAmount, &stats.NewCampaigns30Days,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get campaign stats: %w", err)
	}

	// Calculate success rate
	if stats.CompletedCampaigns > 0 {
		stats.SuccessRate = float64(stats.CompletedCampaigns) / float64(stats.TotalCampaigns) * 100
	}

	return stats, nil
}
// CampaignFilters represents filters for campaign queries
type CampaignFilters struct {
	Status       string
	Category     string
	NGOID        *uuid.UUID
	SearchTerm   string
	IsFeatured   *bool
	MinAmount    *float64
	MaxAmount    *float64
	StartAfter   *time.Time
	EndBefore    *time.Time
}

// CampaignStats represents campaign statistics
type CampaignStats struct {
	TotalCampaigns      int     `json:"total_campaigns"`
	ActiveCampaigns     int     `json:"active_campaigns"`
	DraftCampaigns      int     `json:"draft_campaigns"`
	CompletedCampaigns  int     `json:"completed_campaigns"`
	FeaturedCampaigns   int     `json:"featured_campaigns"`
	TotalGoalAmount     float64 `json:"total_goal_amount"`
	TotalRaisedAmount   float64 `json:"total_raised_amount"`
	SuccessRate         float64 `json:"success_rate"`
	NewCampaigns30Days  int     `json:"new_campaigns_30_days"`
}

// Helper methods

// buildCampaignFilters builds WHERE clause and arguments for campaign filters
func (r *CampaignRepository) buildCampaignFilters(filters CampaignFilters) (string, []interface{}) {
	conditions := []string{}
	args := []interface{}{}
	argIndex := 1

	// Always exclude deleted campaigns unless explicitly searching for them
	if filters.Status != "deleted" {
		conditions = append(conditions, "status != 'deleted'")
	}

	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filters.Status)
		argIndex++
	}

	if filters.Category != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(category) = $%d", argIndex))
		args = append(args, strings.ToLower(filters.Category))
		argIndex++
	}

	if filters.NGOID != nil {
		conditions = append(conditions, fmt.Sprintf("ngo_id = $%d", argIndex))
		args = append(args, *filters.NGOID)
		argIndex++
	}

	if filters.SearchTerm != "" {
		searchTerm := "%" + strings.ToLower(filters.SearchTerm) + "%"
		conditions = append(conditions, fmt.Sprintf("(LOWER(title) LIKE $%d OR LOWER(description) LIKE $%d)", argIndex, argIndex))
		args = append(args, searchTerm)
		argIndex++
	}

	if filters.IsFeatured != nil {
		conditions = append(conditions, fmt.Sprintf("is_featured = $%d", argIndex))
		args = append(args, *filters.IsFeatured)
		argIndex++
	}

	if filters.MinAmount != nil {
		conditions = append(conditions, fmt.Sprintf("goal_amount >= $%d", argIndex))
		args = append(args, *filters.MinAmount)
		argIndex++
	}

	if filters.MaxAmount != nil {
		conditions = append(conditions, fmt.Sprintf("goal_amount <= $%d", argIndex))
		args = append(args, *filters.MaxAmount)
		argIndex++
	}

	if filters.StartAfter != nil {
		conditions = append(conditions, fmt.Sprintf("start_date >= $%d", argIndex))
		args = append(args, *filters.StartAfter)
		argIndex++
	}

	if filters.EndBefore != nil {
		conditions = append(conditions, fmt.Sprintf("end_date <= $%d", argIndex))
		args = append(args, *filters.EndBefore)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	return whereClause, args
}

// validateCampaignInput validates campaign creation input
func (r *CampaignRepository) validateCampaignInput(campaign *models.CampaignCreate) error {
	if campaign.Title == "" {
		return errors.New("campaign title is required")
	}

	if len(strings.TrimSpace(campaign.Title)) < 5 {
		return errors.New("campaign title must be at least 5 characters long")
	}

	if len(campaign.Title) > 255 {
		return errors.New("campaign title must not exceed 255 characters")
	}

	if campaign.GoalAmount <= 0 {
		return ErrGoalAmountInvalid
	}

	if campaign.GoalAmount > 10000000 { // 10 million limit
		return errors.New("goal amount cannot exceed 10,000,000")
	}

	return nil
}

// validateCampaignDates validates start and end dates
func (r *CampaignRepository) validateCampaignDates(startDate, endDate *time.Time) error {
	if startDate != nil && endDate != nil {
		if startDate.After(*endDate) {
			return errors.New("start date cannot be after end date")
		}
	}

	if startDate != nil && startDate.Before(time.Now().AddDate(0, 0, -1)) {
		return errors.New("start date cannot be in the past")
	}

	if endDate != nil && endDate.Before(time.Now()) {
		return errors.New("end date cannot be in the past")
	}

	return nil
}

// campaignTitleExists checks if a campaign title exists for the same NGO
func (r *CampaignRepository) campaignTitleExists(ctx context.Context, ngoID uuid.UUID, title string, excludeID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM campaigns WHERE ngo_id = $1 AND LOWER(title) = LOWER($2) AND id != $3 AND status != 'deleted')`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, ngoID, title, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check campaign title existence: %w", err)
	}
	return exists, nil
}

// Interface implementation methods

// Create implements CampaignRepositoryInterface
func (r *CampaignRepository) Create(ctx context.Context, campaign *models.Campaign) error {
	query := `
		INSERT INTO campaigns (ngo_id, title, description, short_description, goal_amount, raised_amount,
			currency, status, category, tags, start_date, end_date, image_url, gallery_images, 
			location, is_featured)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(ctx, query,
		campaign.NGOID, campaign.Title, campaign.Description, campaign.ShortDescription,
		campaign.GoalAmount, campaign.RaisedAmount, campaign.Currency, campaign.Status,
		campaign.Category, campaign.Tags, campaign.StartDate, campaign.EndDate,
		campaign.ImageURL, campaign.GalleryImages, campaign.Location, campaign.IsFeatured,
	).Scan(&campaign.ID, &campaign.CreatedAt, &campaign.UpdatedAt)
}

// GetByID implements CampaignRepositoryInterface
func (r *CampaignRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Campaign, error) {
	return r.GetCampaignByID(ctx, id)
}

// Update implements CampaignRepositoryInterface
func (r *CampaignRepository) Update(ctx context.Context, campaign *models.Campaign) error {
	query := `
		UPDATE campaigns SET
			title = $2, description = $3, short_description = $4, goal_amount = $5,
			currency = $6, status = $7, category = $8, tags = $9, start_date = $10,
			end_date = $11, image_url = $12, gallery_images = $13, location = $14,
			is_featured = $15, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	return r.db.QueryRowContext(ctx, query,
		campaign.ID, campaign.Title, campaign.Description, campaign.ShortDescription,
		campaign.GoalAmount, campaign.Currency, campaign.Status, campaign.Category,
		campaign.Tags, campaign.StartDate, campaign.EndDate, campaign.ImageURL,
		campaign.GalleryImages, campaign.Location, campaign.IsFeatured,
	).Scan(&campaign.UpdatedAt)
}

// Delete implements CampaignRepositoryInterface  
func (r *CampaignRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DeleteCampaign(ctx, id)
}

// List implements CampaignRepositoryInterface
func (r *CampaignRepository) List(ctx context.Context, params models.Pagination, filters CampaignFilters) ([]*models.Campaign, int64, error) {
	return r.ListCampaigns(ctx, params, filters)
}

// ListByNGO implements CampaignRepositoryInterface
func (r *CampaignRepository) ListByNGO(ctx context.Context, ngoID uuid.UUID, params models.Pagination) ([]*models.Campaign, int64, error) {
	return r.ListCampaignsByNGO(ctx, ngoID, params)
}

// UpdateStatus implements CampaignRepositoryInterface
func (r *CampaignRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.UpdateCampaignStatus(ctx, id, status)
}

// UpdateRaisedAmountInterface implements CampaignRepositoryInterface
func (r *CampaignRepository) UpdateRaisedAmountInterface(ctx context.Context, id uuid.UUID, amount float64) error {
	return r.UpdateRaisedAmount(ctx, id, amount)
}

// Helper utility functions

// stringPtrOrNil returns a pointer to string, or nil if string is empty
func stringPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// getDefaultCurrency returns default currency if not specified
func getDefaultCurrency(currency string) string {
	if currency == "" {
		return "USD"
	}
	return currency
}