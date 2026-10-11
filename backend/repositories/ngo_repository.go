package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"ngo-donation-tracker/models"

	"github.com/google/uuid"
)

type NGORepository struct {
	db *sql.DB
}

func NewNGORepository(db *sql.DB) *NGORepository {
	return &NGORepository{db: db}
}

func (r *NGORepository) Create(ctx context.Context, ngo *models.NGO) error {
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

func (r *NGORepository) GetByID(ctx context.Context, id uuid.UUID) (*models.NGO, error) {
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
		return nil, nil
	}
	return ngo, err
}

func (r *NGORepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.NGO, error) {
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
		return nil, nil
	}
	return ngo, err
}

func (r *NGORepository) Update(ctx context.Context, ngo *models.NGO) error {
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

func (r *NGORepository) List(ctx context.Context, params models.Pagination, status string) ([]*models.NGO, int64, error) {
	// Build query
	countQuery := `SELECT COUNT(*) FROM ngos`
	listQuery := `
		SELECT id, user_id, name, description, mission, website, registration_number,
			registration_document_url, logo_url, cover_image_url, address, city, country,
			verification_status, verification_notes, verified_at, verified_by, created_at, updated_at
		FROM ngos
	`
	args := []interface{}{}
	argIndex := 1

	if status != "" {
		countQuery += fmt.Sprintf(" WHERE verification_status = $%d", argIndex)
		listQuery += fmt.Sprintf(" WHERE verification_status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	listQuery += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", argIndex) + ` OFFSET $` + fmt.Sprintf("%d", argIndex+1)
	args = append(args, params.Limit, params.Offset)

	// Get total count
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args[:argIndex-1]...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
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
			return nil, 0, err
		}
		ngos = append(ngos, ngo)
	}

	return ngos, total, nil
}

func (r *NGORepository) UpdateVerificationStatus(ctx context.Context, id uuid.UUID, status string, notes *string, verifiedBy *uuid.UUID) error {
	query := `
		UPDATE ngos SET
			verification_status = $2, verification_notes = $3, verified_by = $4,
			verified_at = CASE WHEN $2 = 'verified' THEN NOW() ELSE NULL END,
			updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, status, notes, verifiedBy)
	return err
}