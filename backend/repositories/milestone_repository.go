package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ngo-donation-tracker/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// MilestoneRepository interface defines methods for milestone data operations
type MilestoneRepository interface {
	// Enhanced methods (main interface)
	CreateMilestone(milestone *models.Milestone) error
	GetMilestoneByID(id uuid.UUID) (*models.Milestone, error)
	UpdateMilestone(milestone *models.Milestone) error
	DeleteMilestone(id uuid.UUID) error
	ListMilestonesByCampaign(campaignID uuid.UUID) ([]*models.Milestone, error)
	UpdateMilestoneStatus(id uuid.UUID, status string) error
	MarkMilestoneComplete(id uuid.UUID, proofURL, proofNotes string) error
	GetMilestonesByStatus(status string, pagination PaginationParams) ([]*models.Milestone, *PaginationResult, error)
	GetMilestoneStats(campaignID uuid.UUID) (*MilestoneStats, error)
	
	// Interface compatibility methods (MilestoneRepositoryInterface)
	Create(ctx context.Context, milestone *models.Milestone) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Milestone, error)
	Update(ctx context.Context, milestone *models.Milestone) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByCampaign(ctx context.Context, campaignID uuid.UUID) ([]*models.Milestone, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	MarkComplete(ctx context.Context, id uuid.UUID, proofURL, proofNotes string) error
}

// MilestoneStats represents aggregated milestone statistics
type MilestoneStats struct {
	TotalMilestones     int     `json:"total_milestones"`
	CompletedMilestones int     `json:"completed_milestones"`
	PendingMilestones   int     `json:"pending_milestones"`
	OverdueMilestones   int     `json:"overdue_milestones"`
	TotalTargetAmount   float64 `json:"total_target_amount"`
	CompletedAmount     float64 `json:"completed_amount"`
	CompletionRate      float64 `json:"completion_rate"`
}

// MilestoneWithCampaign represents milestone with campaign details
type MilestoneWithCampaign struct {
	*models.Milestone
	CampaignTitle string `json:"campaign_title"`
	NGOName       string `json:"ngo_name"`
}

// milestoneRepository implements MilestoneRepository interface
type milestoneRepository struct {
	db *sql.DB
}

// NewMilestoneRepository creates a new instance of milestone repository
func NewMilestoneRepository(db *sql.DB) MilestoneRepository {
	return &milestoneRepository{db: db}
}

// Interface compatibility methods

// CreateInterface creates a new milestone (interface method)
func (r *milestoneRepository) Create(ctx context.Context, milestone *models.Milestone) error {
	return r.CreateMilestone(milestone)
}

func (r *milestoneRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Milestone, error) {
	return r.GetMilestoneByID(id)
}

func (r *milestoneRepository) Update(ctx context.Context, milestone *models.Milestone) error {
	return r.UpdateMilestone(milestone)
}

func (r *milestoneRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DeleteMilestone(id)
}

func (r *milestoneRepository) ListByCampaign(ctx context.Context, campaignID uuid.UUID) ([]*models.Milestone, error) {
	return r.ListMilestonesByCampaign(campaignID)
}

func (r *milestoneRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.UpdateMilestoneStatus(id, status)
}

func (r *milestoneRepository) MarkComplete(ctx context.Context, id uuid.UUID, proofURL, proofNotes string) error {
	return r.MarkMilestoneComplete(id, proofURL, proofNotes)
}

// Enhanced methods

// CreateMilestone creates a new milestone record
func (r *milestoneRepository) CreateMilestone(milestone *models.Milestone) error {
	query := `
		INSERT INTO milestones (
			id, campaign_id, title, description, target_amount, raised_amount,
			status, due_date, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)`

	_, err := r.db.Exec(query,
		milestone.ID, milestone.CampaignID, milestone.Title, milestone.Description,
		milestone.TargetAmount, milestone.RaisedAmount, milestone.Status,
		milestone.DueDate, milestone.CreatedAt, milestone.UpdatedAt)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503": // Foreign key violation
				return fmt.Errorf("invalid campaign ID")
			case "23505": // Unique violation
				return fmt.Errorf("milestone with this title already exists for this campaign")
			default:
				return fmt.Errorf("failed to create milestone: %v", err)
			}
		}
		return fmt.Errorf("failed to create milestone: %v", err)
	}

	return nil
}

// GetMilestoneByID retrieves a milestone by its ID
func (r *milestoneRepository) GetMilestoneByID(id uuid.UUID) (*models.Milestone, error) {
	query := `
		SELECT 
			id, campaign_id, title, description, target_amount, raised_amount,
			status, due_date, completed_at, proof_document_url, proof_notes,
			created_at, updated_at
		FROM milestones 
		WHERE id = $1`

	milestone := &models.Milestone{}
	err := r.db.QueryRow(query, id).Scan(
		&milestone.ID, &milestone.CampaignID, &milestone.Title, &milestone.Description,
		&milestone.TargetAmount, &milestone.RaisedAmount, &milestone.Status,
		&milestone.DueDate, &milestone.CompletedAt, &milestone.ProofDocumentURL,
		&milestone.ProofNotes, &milestone.CreatedAt, &milestone.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("milestone not found")
		}
		return nil, fmt.Errorf("failed to get milestone: %v", err)
	}

	return milestone, nil
}

// UpdateMilestone updates an existing milestone
func (r *milestoneRepository) UpdateMilestone(milestone *models.Milestone) error {
	query := `
		UPDATE milestones SET
			title = $2, description = $3, target_amount = $4, raised_amount = $5,
			status = $6, due_date = $7, completed_at = $8, proof_document_url = $9,
			proof_notes = $10, updated_at = $11
		WHERE id = $1`

	result, err := r.db.Exec(query,
		milestone.ID, milestone.Title, milestone.Description, milestone.TargetAmount,
		milestone.RaisedAmount, milestone.Status, milestone.DueDate, milestone.CompletedAt,
		milestone.ProofDocumentURL, milestone.ProofNotes, milestone.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update milestone: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("milestone not found")
	}

	return nil
}

// DeleteMilestone deletes a milestone (soft delete - updates status to cancelled)
func (r *milestoneRepository) DeleteMilestone(id uuid.UUID) error {
	query := `UPDATE milestones SET status = 'cancelled', updated_at = $2 WHERE id = $1`
	
	result, err := r.db.Exec(query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete milestone: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("milestone not found")
	}

	return nil
}

// ListMilestonesByCampaign retrieves all milestones for a campaign
func (r *milestoneRepository) ListMilestonesByCampaign(campaignID uuid.UUID) ([]*models.Milestone, error) {
	query := `
		SELECT 
			id, campaign_id, title, description, target_amount, raised_amount,
			status, due_date, completed_at, proof_document_url, proof_notes,
			created_at, updated_at
		FROM milestones 
		WHERE campaign_id = $1 AND status != 'cancelled'
		ORDER BY due_date ASC, created_at ASC`

	rows, err := r.db.Query(query, campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to query milestones: %v", err)
	}
	defer rows.Close()

	var milestones []*models.Milestone
	for rows.Next() {
		milestone := &models.Milestone{}
		err := rows.Scan(
			&milestone.ID, &milestone.CampaignID, &milestone.Title, &milestone.Description,
			&milestone.TargetAmount, &milestone.RaisedAmount, &milestone.Status,
			&milestone.DueDate, &milestone.CompletedAt, &milestone.ProofDocumentURL,
			&milestone.ProofNotes, &milestone.CreatedAt, &milestone.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan milestone: %v", err)
		}
		milestones = append(milestones, milestone)
	}

	return milestones, nil
}

// UpdateMilestoneStatus updates milestone status
func (r *milestoneRepository) UpdateMilestoneStatus(id uuid.UUID, status string) error {
	query := `UPDATE milestones SET status = $2, updated_at = $3 WHERE id = $1`
	
	result, err := r.db.Exec(query, id, status, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update milestone status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("milestone not found")
	}

	return nil
}

// MarkMilestoneComplete marks milestone as complete with proof
func (r *milestoneRepository) MarkMilestoneComplete(id uuid.UUID, proofURL, proofNotes string) error {
	query := `
		UPDATE milestones SET
			status = 'completed', completed_at = $2, proof_document_url = $3,
			proof_notes = $4, updated_at = $5
		WHERE id = $1`
	
	result, err := r.db.Exec(query, id, time.Now(), proofURL, proofNotes, time.Now())
	if err != nil {
		return fmt.Errorf("failed to mark milestone complete: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("milestone not found")
	}

	return nil
}

// GetMilestonesByStatus retrieves milestones by status with pagination
func (r *milestoneRepository) GetMilestonesByStatus(status string, pagination PaginationParams) ([]*models.Milestone, *PaginationResult, error) {
	countQuery := `SELECT COUNT(*) FROM milestones WHERE status = $1`
	
	var totalCount int
	err := r.db.QueryRow(countQuery, status).Scan(&totalCount)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count milestones: %v", err)
	}

	query := `
		SELECT 
			id, campaign_id, title, description, target_amount, raised_amount,
			status, due_date, completed_at, proof_document_url, proof_notes,
			created_at, updated_at
		FROM milestones 
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, status, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query milestones: %v", err)
	}
	defer rows.Close()

	var milestones []*models.Milestone
	for rows.Next() {
		milestone := &models.Milestone{}
		err := rows.Scan(
			&milestone.ID, &milestone.CampaignID, &milestone.Title, &milestone.Description,
			&milestone.TargetAmount, &milestone.RaisedAmount, &milestone.Status,
			&milestone.DueDate, &milestone.CompletedAt, &milestone.ProofDocumentURL,
			&milestone.ProofNotes, &milestone.CreatedAt, &milestone.UpdatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan milestone: %v", err)
		}
		milestones = append(milestones, milestone)
	}

	result := &PaginationResult{
		TotalCount:  totalCount,
		TotalPages:  (totalCount + pagination.Limit - 1) / pagination.Limit,
		CurrentPage: pagination.Page,
		PageSize:    pagination.Limit,
	}

	return milestones, result, nil
}

// GetMilestoneStats retrieves aggregated milestone statistics for a campaign
func (r *milestoneRepository) GetMilestoneStats(campaignID uuid.UUID) (*MilestoneStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_milestones,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_milestones,
			COUNT(CASE WHEN status = 'pending' OR status = 'in_progress' THEN 1 END) as pending_milestones,
			COUNT(CASE WHEN due_date < NOW() AND status != 'completed' THEN 1 END) as overdue_milestones,
			COALESCE(SUM(target_amount), 0) as total_target_amount,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN target_amount ELSE 0 END), 0) as completed_amount
		FROM milestones 
		WHERE campaign_id = $1 AND status != 'cancelled'`

	stats := &MilestoneStats{}
	err := r.db.QueryRow(query, campaignID).Scan(
		&stats.TotalMilestones, &stats.CompletedMilestones, &stats.PendingMilestones,
		&stats.OverdueMilestones, &stats.TotalTargetAmount, &stats.CompletedAmount)

	if err != nil {
		return nil, fmt.Errorf("failed to get milestone stats: %v", err)
	}

	// Calculate completion rate
	if stats.TotalMilestones > 0 {
		stats.CompletionRate = (float64(stats.CompletedMilestones) / float64(stats.TotalMilestones)) * 100
	}

	return stats, nil
}