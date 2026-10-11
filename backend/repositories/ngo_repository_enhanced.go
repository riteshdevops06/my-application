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

// Custom errors for NGO operations
var (
	ErrNGONotFound           = errors.New("NGO not found")
	ErrNGOAlreadyExists      = errors.New("NGO profile already exists for this user")
	ErrInvalidVerificationStatus = errors.New("invalid verification status")
	ErrNGONotVerified        = errors.New("NGO is not verified")
)

// NGORepositoryEnhanced handles all NGO-related database operations
type NGORepositoryEnhanced struct {
	db *sql.DB
}

// NewNGORepositoryEnhanced creates a new enhanced NGO repository instance
func NewNGORepositoryEnhanced(db *sql.DB) *NGORepositoryEnhanced {
	return &NGORepositoryEnhanced{db: db}
}

// Create implements the interface for backward compatibility
func (r *NGORepositoryEnhanced) Create(ctx context.Context, ngo *models.NGO) error {
	query := `
		INSERT INTO ngos (user_id, name, description, mission, website, registration_number,
			registration_document_url, logo_url, cover_image_url, address, city, country, verification_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		ngo.UserID, ngo.Name, ngo.Description, ngo.Mission, ngo.Website,
		ngo.RegistrationNumber, ngo.RegistrationDocumentURL, ngo.LogoURL,
		ngo.CoverImageURL, ngo.Address, ngo.City, ngo.Country, ngo.VerificationStatus,
	).Scan(&ngo.ID, &ngo.CreatedAt, &ngo.UpdatedAt)
}

// GetByID implements the interface
func (r *NGORepositoryEnhanced) GetByID(ctx context.Context, id uuid.UUID) (*models.NGO, error) {
	return r.GetNGOByID(ctx, id)
}

// GetByUserID implements the interface
func (r *NGORepositoryEnhanced) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.NGO, error) {
	return r.GetNGOByUserID(ctx, userID)
}

// Update implements the interface for backward compatibility
func (r *NGORepositoryEnhanced) Update(ctx context.Context, ngo *models.NGO) error {
	query := `
		UPDATE ngos SET
			name = $2, description = $3, mission = $4, website = $5,
			registration_number = $6, registration_document_url = $7,
			logo_url = $8, cover_image_url = $9, address = $10, city = $11, country = $12,
			verification_status = $13, verification_notes = $14, verified_at = $15, verified_by = $16,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		ngo.ID, ngo.Name, ngo.Description, ngo.Mission, ngo.Website,
		ngo.RegistrationNumber, ngo.RegistrationDocumentURL,
		ngo.LogoURL, ngo.CoverImageURL, ngo.Address, ngo.City, ngo.Country,
		ngo.VerificationStatus, ngo.VerificationNotes, ngo.VerifiedAt, ngo.VerifiedBy,
	).Scan(&ngo.UpdatedAt)
}

// List implements the interface for backward compatibility
func (r *NGORepositoryEnhanced) List(ctx context.Context, params models.Pagination, status string) ([]*models.NGO, int64, error) {
	filters := NGOFilters{}
	if status != "" {
		filters.VerificationStatus = status
	}
	return r.ListNGOs(ctx, params, filters)
}

// CreateNGO creates a new NGO profile
func (r *NGORepositoryEnhanced) CreateNGO(ctx context.Context, userID uuid.UUID, ngoCreate *models.NGOCreate) (*models.NGO, error) {
	// Check if NGO already exists for this user
	existing, err := r.GetNGOByUserID(ctx, userID)
	if err != nil && err != ErrNGONotFound {
		return nil, fmt.Errorf("failed to check existing NGO: %w", err)
	}
	if existing != nil {
		return nil, ErrNGOAlreadyExists
	}

	ngo := &models.NGO{
		UserID:                  userID,
		Name:                    strings.TrimSpace(ngoCreate.Name),
		Description:             stringPtr(strings.TrimSpace(ngoCreate.Description)),
		Mission:                 stringPtr(strings.TrimSpace(ngoCreate.Mission)),
		Website:                 stringPtr(strings.TrimSpace(ngoCreate.Website)),
		RegistrationNumber:      stringPtr(strings.TrimSpace(ngoCreate.RegistrationNumber)),
		RegistrationDocumentURL: stringPtr(ngoCreate.RegistrationDocumentURL),
		LogoURL:                 stringPtr(ngoCreate.LogoURL),
		CoverImageURL:           stringPtr(ngoCreate.CoverImageURL),
		Address:                 stringPtr(strings.TrimSpace(ngoCreate.Address)),
		City:                    stringPtr(strings.TrimSpace(ngoCreate.City)),
		Country:                 stringPtr(strings.TrimSpace(ngoCreate.Country)),
		VerificationStatus:      "pending",
	}

	query := `
		INSERT INTO ngos (user_id, name, description, mission, website, registration_number,
			registration_document_url, logo_url, cover_image_url, address, city, country, verification_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`

	err = r.db.QueryRowContext(ctx, query,
		ngo.UserID, ngo.Name, ngo.Description, ngo.Mission, ngo.Website,
		ngo.RegistrationNumber, ngo.RegistrationDocumentURL, ngo.LogoURL,
		ngo.CoverImageURL, ngo.Address, ngo.City, ngo.Country, ngo.VerificationStatus,
	).Scan(&ngo.ID, &ngo.CreatedAt, &ngo.UpdatedAt)

	if err != nil {
		// Handle unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, ErrNGOAlreadyExists
		}
		return nil, fmt.Errorf("failed to create NGO: %w", err)
	}

	return ngo, nil
}
// GetNGOByID retrieves an NGO by its ID
func (r *NGORepositoryEnhanced) GetNGOByID(ctx context.Context, id uuid.UUID) (*models.NGO, error) {
	query := `
		SELECT id, user_id, name, description, mission, website, registration_number,
			registration_document_url, logo_url, cover_image_url, address, city, country,
			verification_status, verification_notes, verified_at, verified_by, created_at, updated_at
		FROM ngos WHERE id = $1
	`

	ngo := &models.NGO{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ngo.ID, &ngo.UserID, &ngo.Name, &ngo.Description, &ngo.Mission,
		&ngo.Website, &ngo.RegistrationNumber, &ngo.RegistrationDocumentURL,
		&ngo.LogoURL, &ngo.CoverImageURL, &ngo.Address, &ngo.City, &ngo.Country,
		&ngo.VerificationStatus, &ngo.VerificationNotes, &ngo.VerifiedAt, &ngo.VerifiedBy,
		&ngo.CreatedAt, &ngo.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNGONotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get NGO by ID: %w", err)
	}

	return ngo, nil
}

// GetNGOByUserID retrieves an NGO by user ID
func (r *NGORepositoryEnhanced) GetNGOByUserID(ctx context.Context, userID uuid.UUID) (*models.NGO, error) {
	query := `
		SELECT id, user_id, name, description, mission, website, registration_number,
			registration_document_url, logo_url, cover_image_url, address, city, country,
			verification_status, verification_notes, verified_at, verified_by, created_at, updated_at
		FROM ngos WHERE user_id = $1
	`

	ngo := &models.NGO{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&ngo.ID, &ngo.UserID, &ngo.Name, &ngo.Description, &ngo.Mission,
		&ngo.Website, &ngo.RegistrationNumber, &ngo.RegistrationDocumentURL,
		&ngo.LogoURL, &ngo.CoverImageURL, &ngo.Address, &ngo.City, &ngo.Country,
		&ngo.VerificationStatus, &ngo.VerificationNotes, &ngo.VerifiedAt, &ngo.VerifiedBy,
		&ngo.CreatedAt, &ngo.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNGONotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get NGO by user ID: %w", err)
	}

	return ngo, nil
}
// UpdateNGO updates NGO profile information
func (r *NGORepositoryEnhanced) UpdateNGO(ctx context.Context, id uuid.UUID, update *models.NGOCreate) (*models.NGO, error) {
	// Build dynamic query based on fields to update
	setParts := []string{}
	args := []interface{}{id}
	argIndex := 2

	if update.Name != "" {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, strings.TrimSpace(update.Name))
		argIndex++
	}

	if update.Description != "" {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, stringPtr(strings.TrimSpace(update.Description)))
		argIndex++
	}

	if update.Mission != "" {
		setParts = append(setParts, fmt.Sprintf("mission = $%d", argIndex))
		args = append(args, stringPtr(strings.TrimSpace(update.Mission)))
		argIndex++
	}

	if update.Website != "" {
		setParts = append(setParts, fmt.Sprintf("website = $%d", argIndex))
		args = append(args, stringPtr(strings.TrimSpace(update.Website)))
		argIndex++
	}

	if update.RegistrationNumber != "" {
		setParts = append(setParts, fmt.Sprintf("registration_number = $%d", argIndex))
		args = append(args, stringPtr(strings.TrimSpace(update.RegistrationNumber)))
		argIndex++
	}

	if update.RegistrationDocumentURL != "" {
		setParts = append(setParts, fmt.Sprintf("registration_document_url = $%d", argIndex))
		args = append(args, stringPtr(update.RegistrationDocumentURL))
		argIndex++
	}

	if update.LogoURL != "" {
		setParts = append(setParts, fmt.Sprintf("logo_url = $%d", argIndex))
		args = append(args, stringPtr(update.LogoURL))
		argIndex++
	}

	if update.CoverImageURL != "" {
		setParts = append(setParts, fmt.Sprintf("cover_image_url = $%d", argIndex))
		args = append(args, stringPtr(update.CoverImageURL))
		argIndex++
	}

	if update.Address != "" {
		setParts = append(setParts, fmt.Sprintf("address = $%d", argIndex))
		args = append(args, stringPtr(strings.TrimSpace(update.Address)))
		argIndex++
	}

	if update.City != "" {
		setParts = append(setParts, fmt.Sprintf("city = $%d", argIndex))
		args = append(args, stringPtr(strings.TrimSpace(update.City)))
		argIndex++
	}

	if update.Country != "" {
		setParts = append(setParts, fmt.Sprintf("country = $%d", argIndex))
		args = append(args, stringPtr(strings.TrimSpace(update.Country)))
		argIndex++
	}

	if len(setParts) == 0 {
		// No fields to update, just return current NGO
		return r.GetNGOByID(ctx, id)
	}

	setParts = append(setParts, "updated_at = NOW()")

	query := fmt.Sprintf(`
		UPDATE ngos SET %s
		WHERE id = $1
		RETURNING id, user_id, name, description, mission, website, registration_number,
			registration_document_url, logo_url, cover_image_url, address, city, country,
			verification_status, verification_notes, verified_at, verified_by, created_at, updated_at
	`, strings.Join(setParts, ", "))

	ngo := &models.NGO{}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&ngo.ID, &ngo.UserID, &ngo.Name, &ngo.Description, &ngo.Mission,
		&ngo.Website, &ngo.RegistrationNumber, &ngo.RegistrationDocumentURL,
		&ngo.LogoURL, &ngo.CoverImageURL, &ngo.Address, &ngo.City, &ngo.Country,
		&ngo.VerificationStatus, &ngo.VerificationNotes, &ngo.VerifiedAt, &ngo.VerifiedBy,
		&ngo.CreatedAt, &ngo.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNGONotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update NGO: %w", err)
	}

	return ngo, nil
}
// ListNGOs retrieves paginated list of NGOs with filters
func (r *NGORepositoryEnhanced) ListNGOs(ctx context.Context, params models.Pagination, filters NGOFilters) ([]*models.NGO, int64, error) {
	// Build WHERE clause
	whereClause, filterArgs := r.buildNGOFilters(filters)

	// Get total count
	countQuery := "SELECT COUNT(*) FROM ngos" + whereClause
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, filterArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count NGOs: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, user_id, name, description, mission, website, registration_number,
			registration_document_url, logo_url, cover_image_url, address, city, country,
			verification_status, verification_notes, verified_at, verified_by, created_at, updated_at
		FROM ngos` + whereClause + `
		ORDER BY created_at DESC
		LIMIT $` + fmt.Sprintf("%d", len(filterArgs)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(filterArgs)+2)

	allArgs := append(filterArgs, params.Limit, params.Offset)

	rows, err := r.db.QueryContext(ctx, query, allArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list NGOs: %w", err)
	}
	defer rows.Close()

	var ngos []*models.NGO
	for rows.Next() {
		ngo := &models.NGO{}
		err := rows.Scan(
			&ngo.ID, &ngo.UserID, &ngo.Name, &ngo.Description, &ngo.Mission,
			&ngo.Website, &ngo.RegistrationNumber, &ngo.RegistrationDocumentURL,
			&ngo.LogoURL, &ngo.CoverImageURL, &ngo.Address, &ngo.City, &ngo.Country,
			&ngo.VerificationStatus, &ngo.VerificationNotes, &ngo.VerifiedAt, &ngo.VerifiedBy,
			&ngo.CreatedAt, &ngo.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan NGO: %w", err)
		}
		ngos = append(ngos, ngo)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating NGOs: %w", err)
	}

	return ngos, total, nil
}

// UpdateVerificationStatus updates NGO verification status
func (r *NGORepositoryEnhanced) UpdateVerificationStatus(ctx context.Context, id uuid.UUID, status string, notes *string, verifiedBy *uuid.UUID) error {
	// Validate status
	validStatuses := []string{"pending", "verified", "rejected"}
	isValidStatus := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		return ErrInvalidVerificationStatus
	}

	query := `
		UPDATE ngos SET
			verification_status = $2, verification_notes = $3, verified_by = $4,
			verified_at = CASE WHEN $2 = 'verified' THEN NOW() ELSE NULL END,
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, status, notes, verifiedBy)
	if err != nil {
		return fmt.Errorf("failed to update verification status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNGONotFound
	}

	return nil
}

// GetVerifiedNGOs returns all verified NGOs
func (r *NGORepositoryEnhanced) GetVerifiedNGOs(ctx context.Context, params models.Pagination) ([]*models.NGO, int64, error) {
	filters := NGOFilters{
		VerificationStatus: "verified",
	}
	return r.ListNGOs(ctx, params, filters)
}
// GetPendingVerifications returns NGOs pending verification
func (r *NGORepositoryEnhanced) GetPendingVerifications(ctx context.Context, params models.Pagination) ([]*models.NGO, int64, error) {
	filters := NGOFilters{
		VerificationStatus: "pending",
	}
	return r.ListNGOs(ctx, params, filters)
}

// GetNGOStats returns NGO statistics for dashboard
func (r *NGORepositoryEnhanced) GetNGOStats(ctx context.Context) (*NGOStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_ngos,
			COUNT(CASE WHEN verification_status = 'verified' THEN 1 END) as verified_ngos,
			COUNT(CASE WHEN verification_status = 'pending' THEN 1 END) as pending_ngos,
			COUNT(CASE WHEN verification_status = 'rejected' THEN 1 END) as rejected_ngos,
			COUNT(CASE WHEN created_at >= NOW() - INTERVAL '30 days' THEN 1 END) as new_ngos_30d
		FROM ngos
	`

	stats := &NGOStats{}
	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalNGOs, &stats.VerifiedNGOs, &stats.PendingNGOs,
		&stats.RejectedNGOs, &stats.NewNGOs30Days,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get NGO stats: %w", err)
	}

	return stats, nil
}

// GetRecentNGOs returns recently registered NGOs
func (r *NGORepositoryEnhanced) GetRecentNGOs(ctx context.Context, limit int) ([]*models.NGO, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	query := `
		SELECT id, user_id, name, description, mission, website, registration_number,
			registration_document_url, logo_url, cover_image_url, address, city, country,
			verification_status, verification_notes, verified_at, verified_by, created_at, updated_at
		FROM ngos
		ORDER BY created_at DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent NGOs: %w", err)
	}
	defer rows.Close()

	var ngos []*models.NGO
	for rows.Next() {
		ngo := &models.NGO{}
		err := rows.Scan(
			&ngo.ID, &ngo.UserID, &ngo.Name, &ngo.Description, &ngo.Mission,
			&ngo.Website, &ngo.RegistrationNumber, &ngo.RegistrationDocumentURL,
			&ngo.LogoURL, &ngo.CoverImageURL, &ngo.Address, &ngo.City, &ngo.Country,
			&ngo.VerificationStatus, &ngo.VerificationNotes, &ngo.VerifiedAt, &ngo.VerifiedBy,
			&ngo.CreatedAt, &ngo.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recent NGO: %w", err)
		}
		ngos = append(ngos, ngo)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating recent NGOs: %w", err)
	}

	return ngos, nil
}

// SearchNGOsByName searches NGOs by name
func (r *NGORepositoryEnhanced) SearchNGOsByName(ctx context.Context, searchTerm string, params models.Pagination) ([]*models.NGO, int64, error) {
	filters := NGOFilters{
		SearchTerm: searchTerm,
	}
	return r.ListNGOs(ctx, params, filters)
}

// GetNGOsByCountry returns NGOs filtered by country
func (r *NGORepositoryEnhanced) GetNGOsByCountry(ctx context.Context, country string, params models.Pagination) ([]*models.NGO, int64, error) {
	filters := NGOFilters{
		Country: country,
	}
	return r.ListNGOs(ctx, params, filters)
}

// DeleteNGO permanently deletes an NGO (use with caution)
func (r *NGORepositoryEnhanced) DeleteNGO(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM ngos WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete NGO: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNGONotFound
	}

	return nil
}
// NGOFilters represents filters for NGO queries
type NGOFilters struct {
	VerificationStatus string
	Country            string
	City               string
	SearchTerm         string // Search in name, description, mission
	CreatedAfter       *time.Time
	CreatedBefore      *time.Time
}

// NGOStats represents NGO statistics
type NGOStats struct {
	TotalNGOs     int `json:"total_ngos"`
	VerifiedNGOs  int `json:"verified_ngos"`
	PendingNGOs   int `json:"pending_ngos"`
	RejectedNGOs  int `json:"rejected_ngos"`
	NewNGOs30Days int `json:"new_ngos_30_days"`
}

// Helper methods

// buildNGOFilters builds WHERE clause and arguments for NGO filters
func (r *NGORepositoryEnhanced) buildNGOFilters(filters NGOFilters) (string, []interface{}) {
	conditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if filters.VerificationStatus != "" {
		conditions = append(conditions, fmt.Sprintf("verification_status = $%d", argIndex))
		args = append(args, filters.VerificationStatus)
		argIndex++
	}

	if filters.Country != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(country) = $%d", argIndex))
		args = append(args, strings.ToLower(filters.Country))
		argIndex++
	}

	if filters.City != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(city) = $%d", argIndex))
		args = append(args, strings.ToLower(filters.City))
		argIndex++
	}

	if filters.SearchTerm != "" {
		searchTerm := "%" + strings.ToLower(filters.SearchTerm) + "%"
		conditions = append(conditions, fmt.Sprintf("(LOWER(name) LIKE $%d OR LOWER(description) LIKE $%d OR LOWER(mission) LIKE $%d)", argIndex, argIndex, argIndex))
		args = append(args, searchTerm)
		argIndex++
	}

	if filters.CreatedAfter != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *filters.CreatedAfter)
		argIndex++
	}

	if filters.CreatedBefore != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *filters.CreatedBefore)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	return whereClause, args
}