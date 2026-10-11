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

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page   int `json:"page"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// PaginationResult represents pagination result metadata
type PaginationResult struct {
	TotalCount  int `json:"total_count"`
	TotalPages  int `json:"total_pages"`
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
}

// DonationRepository interface defines methods for donation data operations (enhanced)
type DonationRepository interface {
	// Enhanced methods (main interface)
	CreateDonation(donation *models.Donation) error
	GetDonationByID(id uuid.UUID) (*models.Donation, error)
	Update(donation *models.Donation) error
	Delete(id uuid.UUID) error
	GetByDonorID(donorID uuid.UUID, pagination PaginationParams) ([]*models.Donation, *PaginationResult, error)
	GetByCampaignID(campaignID uuid.UUID, pagination PaginationParams) ([]*models.Donation, *PaginationResult, error)
	GetByStatus(status string, pagination PaginationParams) ([]*models.Donation, *PaginationResult, error)
	GetDonationStats(campaignID uuid.UUID) (*DonationStats, error)
	ListEnhanced(filters DonationFiltersExtended, pagination PaginationParams) ([]*models.Donation, *PaginationResult, error)
	GetDonationHistory(donorID uuid.UUID, pagination PaginationParams) ([]*models.DonationWithCampaign, error)
	UpdateStatus(id uuid.UUID, status string, transactionID string) error
	GetDonationsByDateRange(startDate, endDate time.Time, pagination PaginationParams) ([]*models.Donation, *PaginationResult, error)
	
	// Interface compatibility methods (DonationRepositoryInterface)
	Create(ctx context.Context, donation *models.Donation) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Donation, error)
	List(ctx context.Context, params models.Pagination, filters DonationFilters) ([]*models.Donation, int64, error)
	ListByCampaign(ctx context.Context, campaignID uuid.UUID, params models.Pagination) ([]*models.Donation, int64, error)
	ListByDonor(ctx context.Context, donorID uuid.UUID, params models.Pagination) ([]*models.Donation, int64, error)
	UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status string, transactionID *string) error
}

// DonationFilters represents filtering options for donation queries
type DonationFiltersExtended struct {
	DonorID     *uuid.UUID `json:"donor_id,omitempty"`
	CampaignID  *uuid.UUID `json:"campaign_id,omitempty"`
	NGOID       *uuid.UUID `json:"ngo_id,omitempty"`
	Status      string     `json:"status,omitempty"`
	MinAmount   *float64   `json:"min_amount,omitempty"`
	MaxAmount   *float64   `json:"max_amount,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Anonymous   *bool      `json:"anonymous,omitempty"`
	SearchTerm  string     `json:"search_term,omitempty"`
}

// DonationStats represents aggregated donation statistics
type DonationStats struct {
	TotalAmount      float64 `json:"total_amount"`
	TotalDonations   int     `json:"total_donations"`
	AverageAmount    float64 `json:"average_amount"`
	CompletedAmount  float64 `json:"completed_amount"`
	PendingAmount    float64 `json:"pending_amount"`
	RefundedAmount   float64 `json:"refunded_amount"`
	UniqueDonors     int     `json:"unique_donors"`
	RecentDonations  int     `json:"recent_donations"` // Last 30 days
}

// DonationWithCampaign represents donation with campaign details
type DonationWithCampaign struct {
	*models.Donation
	CampaignTitle       string `json:"campaign_title"`
	CampaignDescription string `json:"campaign_description"`
	NGOName            string `json:"ngo_name"`
}

// donationRepository implements both DonationRepository and DonationRepositoryInterface interfaces
type donationRepository struct {
	db *sql.DB
}

// NewDonationRepository creates a new instance of donation repository
func NewDonationRepository(db *sql.DB) DonationRepository {
	return &donationRepository{db: db}
}

// Interface methods (DonationRepositoryInterface)

// CreateDonation creates a new donation record (interface method)
func (r *donationRepository) CreateInterface(ctx context.Context, donation *models.Donation) error {
	return r.CreateDonation(donation)
}

// ListDonations lists donations with filters (interface method)
func (r *donationRepository) ListDonations(ctx context.Context, params models.Pagination, filters DonationFilters) ([]*models.Donation, int64, error) {
	// Convert models.Pagination to PaginationParams
	pagination := PaginationParams{
		Page:   params.Page,
		Limit:  params.Limit,
		Offset: (params.Page - 1) * params.Limit,
	}

	// Convert DonationFilters to DonationFiltersExtended
	extendedFilters := DonationFiltersExtended{
		DonorID:    filters.DonorID,
		CampaignID: filters.CampaignID,
		Status:     filters.PaymentStatus, // Map payment status to status
		MinAmount:  filters.MinAmount,
		MaxAmount:  filters.MaxAmount,
		StartDate:  filters.DonatedAfter,
		EndDate:    filters.DonatedBefore,
		Anonymous:  filters.IsAnonymous,
	}

	donations, result, err := r.ListEnhanced(extendedFilters, pagination)
	if err != nil {
		return nil, 0, err
	}

	return donations, int64(result.TotalCount), nil
}

// ListByCampaignInterface lists donations by campaign (interface method)
func (r *donationRepository) ListByCampaignInterface(ctx context.Context, campaignID uuid.UUID, params models.Pagination) ([]*models.Donation, int64, error) {
	pagination := PaginationParams{
		Page:   params.Page,
		Limit:  params.Limit,
		Offset: (params.Page - 1) * params.Limit,
	}

	donations, result, err := r.GetByCampaignID(campaignID, pagination)
	if err != nil {
		return nil, 0, err
	}

	return donations, int64(result.TotalCount), nil
}

// ListByDonorInterface lists donations by donor (interface method)
func (r *donationRepository) ListByDonorInterface(ctx context.Context, donorID uuid.UUID, params models.Pagination) ([]*models.Donation, int64, error) {
	pagination := PaginationParams{
		Page:   params.Page,
		Limit:  params.Limit,
		Offset: (params.Page - 1) * params.Limit,
	}

	donations, result, err := r.GetByDonorID(donorID, pagination)
	if err != nil {
		return nil, 0, err
	}

	return donations, int64(result.TotalCount), nil
}

// UpdatePaymentStatusInterface updates payment status (interface method)
func (r *donationRepository) UpdatePaymentStatusInterface(ctx context.Context, id uuid.UUID, status string, transactionID *string) error {
	txID := ""
	if transactionID != nil {
		txID = *transactionID
	}
	return r.UpdateStatus(id, status, txID)
}

// Interface method implementations for DonationRepositoryInterface
func (r *donationRepository) List(ctx context.Context, params models.Pagination, filters DonationFilters) ([]*models.Donation, int64, error) {
	return r.ListDonations(ctx, params, filters)
}

func (r *donationRepository) ListByCampaign(ctx context.Context, campaignID uuid.UUID, params models.Pagination) ([]*models.Donation, int64, error) {
	return r.ListByCampaignInterface(ctx, campaignID, params)
}

func (r *donationRepository) ListByDonor(ctx context.Context, donorID uuid.UUID, params models.Pagination) ([]*models.Donation, int64, error) {
	return r.ListByDonorInterface(ctx, donorID, params)
}

func (r *donationRepository) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status string, transactionID *string) error {
	return r.UpdatePaymentStatusInterface(ctx, id, status, transactionID)
}

func (r *donationRepository) Create(ctx context.Context, donation *models.Donation) error {
	return r.CreateDonation(donation)
}

func (r *donationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Donation, error) {
	return r.GetDonationByID(id)
}

// CreateDonation creates a new donation record
func (r *donationRepository) CreateDonation(donation *models.Donation) error {
	query := `
		INSERT INTO donations (
			id, donor_id, campaign_id, amount, currency, status,
			payment_method, transaction_id, anonymous, message,
			receipt_email, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)`

	_, err := r.db.Exec(query,
		donation.ID, donation.DonorID, donation.CampaignID, donation.Amount,
		donation.Currency, donation.Status, donation.PaymentMethod,
		donation.TransactionID, donation.Anonymous, donation.Message,
		donation.ReceiptEmail, donation.CreatedAt, donation.UpdatedAt)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503": // Foreign key violation
				return fmt.Errorf("invalid donor or campaign ID")
			case "23505": // Unique violation
				return fmt.Errorf("donation with this transaction ID already exists")
			default:
				return fmt.Errorf("failed to create donation: %v", err)
			}
		}
		return fmt.Errorf("failed to create donation: %v", err)
	}

	return nil
}

// GetDonationByID retrieves a donation by its ID
func (r *donationRepository) GetDonationByID(id uuid.UUID) (*models.Donation, error) {
	query := `
		SELECT 
			id, donor_id, campaign_id, amount, currency, status,
			payment_method, transaction_id, anonymous, message,
			receipt_email, created_at, updated_at
		FROM donations 
		WHERE id = $1`

	donation := &models.Donation{}
	err := r.db.QueryRow(query, id).Scan(
		&donation.ID, &donation.DonorID, &donation.CampaignID,
		&donation.Amount, &donation.Currency, &donation.Status,
		&donation.PaymentMethod, &donation.TransactionID,
		&donation.Anonymous, &donation.Message, &donation.ReceiptEmail,
		&donation.CreatedAt, &donation.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("donation not found")
		}
		return nil, fmt.Errorf("failed to get donation: %v", err)
	}

	return donation, nil
}

// Update updates an existing donation
func (r *donationRepository) Update(donation *models.Donation) error {
	query := `
		UPDATE donations SET
			status = $2, payment_method = $3, transaction_id = $4,
			anonymous = $5, message = $6, receipt_email = $7, updated_at = $8
		WHERE id = $1`

	result, err := r.db.Exec(query,
		donation.ID, donation.Status, donation.PaymentMethod,
		donation.TransactionID, donation.Anonymous, donation.Message,
		donation.ReceiptEmail, donation.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update donation: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("donation not found")
	}

	return nil
}

// Delete deletes a donation (soft delete - updates status to cancelled)
func (r *donationRepository) Delete(id uuid.UUID) error {
	query := `UPDATE donations SET status = 'cancelled', updated_at = $2 WHERE id = $1`
	
	result, err := r.db.Exec(query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete donation: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("donation not found")
	}

	return nil
}

// GetByDonorID retrieves donations by donor ID with pagination
func (r *donationRepository) GetByDonorID(donorID uuid.UUID, pagination PaginationParams) ([]*models.Donation, *PaginationResult, error) {
	countQuery := `SELECT COUNT(*) FROM donations WHERE donor_id = $1`
	
	var totalCount int
	err := r.db.QueryRow(countQuery, donorID).Scan(&totalCount)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count donations: %v", err)
	}

	query := `
		SELECT 
			id, donor_id, campaign_id, amount, currency, status,
			payment_method, transaction_id, anonymous, message,
			receipt_email, created_at, updated_at
		FROM donations 
		WHERE donor_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, donorID, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query donations: %v", err)
	}
	defer rows.Close()

	var donations []*models.Donation
	for rows.Next() {
		donation := &models.Donation{}
		err := rows.Scan(
			&donation.ID, &donation.DonorID, &donation.CampaignID,
			&donation.Amount, &donation.Currency, &donation.Status,
			&donation.PaymentMethod, &donation.TransactionID,
			&donation.Anonymous, &donation.Message, &donation.ReceiptEmail,
			&donation.CreatedAt, &donation.UpdatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan donation: %v", err)
		}
		donations = append(donations, donation)
	}

	result := &PaginationResult{
		TotalCount:  totalCount,
		TotalPages:  (totalCount + pagination.Limit - 1) / pagination.Limit,
		CurrentPage: pagination.Page,
		PageSize:    pagination.Limit,
	}

	return donations, result, nil
}

// GetByCampaignID retrieves donations by campaign ID with pagination
func (r *donationRepository) GetByCampaignID(campaignID uuid.UUID, pagination PaginationParams) ([]*models.Donation, *PaginationResult, error) {
	countQuery := `SELECT COUNT(*) FROM donations WHERE campaign_id = $1`
	
	var totalCount int
	err := r.db.QueryRow(countQuery, campaignID).Scan(&totalCount)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count donations: %v", err)
	}

	query := `
		SELECT 
			id, donor_id, campaign_id, amount, currency, status,
			payment_method, transaction_id, anonymous, message,
			receipt_email, created_at, updated_at
		FROM donations 
		WHERE campaign_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, campaignID, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query donations: %v", err)
	}
	defer rows.Close()

	var donations []*models.Donation
	for rows.Next() {
		donation := &models.Donation{}
		err := rows.Scan(
			&donation.ID, &donation.DonorID, &donation.CampaignID,
			&donation.Amount, &donation.Currency, &donation.Status,
			&donation.PaymentMethod, &donation.TransactionID,
			&donation.Anonymous, &donation.Message, &donation.ReceiptEmail,
			&donation.CreatedAt, &donation.UpdatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan donation: %v", err)
		}
		donations = append(donations, donation)
	}

	result := &PaginationResult{
		TotalCount:  totalCount,
		TotalPages:  (totalCount + pagination.Limit - 1) / pagination.Limit,
		CurrentPage: pagination.Page,
		PageSize:    pagination.Limit,
	}

	return donations, result, nil
}

// GetByStatus retrieves donations by status with pagination
func (r *donationRepository) GetByStatus(status string, pagination PaginationParams) ([]*models.Donation, *PaginationResult, error) {
	countQuery := `SELECT COUNT(*) FROM donations WHERE status = $1`
	
	var totalCount int
	err := r.db.QueryRow(countQuery, status).Scan(&totalCount)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count donations: %v", err)
	}

	query := `
		SELECT 
			id, donor_id, campaign_id, amount, currency, status,
			payment_method, transaction_id, anonymous, message,
			receipt_email, created_at, updated_at
		FROM donations 
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, status, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query donations: %v", err)
	}
	defer rows.Close()

	var donations []*models.Donation
	for rows.Next() {
		donation := &models.Donation{}
		err := rows.Scan(
			&donation.ID, &donation.DonorID, &donation.CampaignID,
			&donation.Amount, &donation.Currency, &donation.Status,
			&donation.PaymentMethod, &donation.TransactionID,
			&donation.Anonymous, &donation.Message, &donation.ReceiptEmail,
			&donation.CreatedAt, &donation.UpdatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan donation: %v", err)
		}
		donations = append(donations, donation)
	}

	result := &PaginationResult{
		TotalCount:  totalCount,
		TotalPages:  (totalCount + pagination.Limit - 1) / pagination.Limit,
		CurrentPage: pagination.Page,
		PageSize:    pagination.Limit,
	}

	return donations, result, nil
}

// GetDonationStats retrieves aggregated donation statistics for a campaign
func (r *donationRepository) GetDonationStats(campaignID uuid.UUID) (*DonationStats, error) {
	query := `
		SELECT 
			COALESCE(SUM(amount), 0) as total_amount,
			COUNT(*) as total_donations,
			COALESCE(AVG(amount), 0) as average_amount,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN amount ELSE 0 END), 0) as completed_amount,
			COALESCE(SUM(CASE WHEN status = 'pending' THEN amount ELSE 0 END), 0) as pending_amount,
			COALESCE(SUM(CASE WHEN status = 'refunded' THEN amount ELSE 0 END), 0) as refunded_amount,
			COUNT(DISTINCT donor_id) as unique_donors,
			COUNT(CASE WHEN created_at > NOW() - INTERVAL '30 days' THEN 1 END) as recent_donations
		FROM donations 
		WHERE campaign_id = $1 AND status != 'cancelled'`

	stats := &DonationStats{}
	err := r.db.QueryRow(query, campaignID).Scan(
		&stats.TotalAmount, &stats.TotalDonations, &stats.AverageAmount,
		&stats.CompletedAmount, &stats.PendingAmount, &stats.RefundedAmount,
		&stats.UniqueDonors, &stats.RecentDonations)

	if err != nil {
		return nil, fmt.Errorf("failed to get donation stats: %v", err)
	}

	return stats, nil
}

// ListEnhanced retrieves donations based on filters and pagination (enhanced version)
func (r *donationRepository) ListEnhanced(filters DonationFiltersExtended, pagination PaginationParams) ([]*models.Donation, *PaginationResult, error) {
	whereClause, args := r.buildWhereClause(filters)
	
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM donations d JOIN campaigns c ON d.campaign_id = c.id %s", whereClause)
	
	var totalCount int
	err := r.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count donations: %v", err)
	}

	query := fmt.Sprintf(`
		SELECT 
			d.id, d.donor_id, d.campaign_id, d.amount, d.currency, d.status,
			d.payment_method, d.transaction_id, d.anonymous, d.message,
			d.receipt_email, d.created_at, d.updated_at
		FROM donations d 
		JOIN campaigns c ON d.campaign_id = c.id
		%s
		ORDER BY d.created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pagination.Limit, pagination.Offset)
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query donations: %v", err)
	}
	defer rows.Close()

	var donations []*models.Donation
	for rows.Next() {
		donation := &models.Donation{}
		err := rows.Scan(
			&donation.ID, &donation.DonorID, &donation.CampaignID,
			&donation.Amount, &donation.Currency, &donation.Status,
			&donation.PaymentMethod, &donation.TransactionID,
			&donation.Anonymous, &donation.Message, &donation.ReceiptEmail,
			&donation.CreatedAt, &donation.UpdatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan donation: %v", err)
		}
		donations = append(donations, donation)
	}

	result := &PaginationResult{
		TotalCount:  totalCount,
		TotalPages:  (totalCount + pagination.Limit - 1) / pagination.Limit,
		CurrentPage: pagination.Page,
		PageSize:    pagination.Limit,
	}

	return donations, result, nil
}

// GetDonationHistory retrieves donation history with campaign details
func (r *donationRepository) GetDonationHistory(donorID uuid.UUID, pagination PaginationParams) ([]*models.DonationWithCampaign, error) {
	query := `
		SELECT 
			d.id, d.donor_id, d.campaign_id, d.amount, d.currency, d.status,
			d.payment_method, d.transaction_id, d.anonymous, d.message,
			d.receipt_email, d.created_at, d.updated_at,
			c.title, c.description, n.name as ngo_name
		FROM donations d
		JOIN campaigns c ON d.campaign_id = c.id
		JOIN ngos n ON c.ngo_id = n.id
		WHERE d.donor_id = $1
		ORDER BY d.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, donorID, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query donation history: %v", err)
	}
	defer rows.Close()

	var donations []*models.DonationWithCampaign
	for rows.Next() {
		donation := &models.DonationWithCampaign{
			Donation: &models.Donation{},
		}
		err := rows.Scan(
			&donation.ID, &donation.DonorID, &donation.CampaignID,
			&donation.Amount, &donation.Currency, &donation.Status,
			&donation.PaymentMethod, &donation.TransactionID,
			&donation.Anonymous, &donation.Message, &donation.ReceiptEmail,
			&donation.CreatedAt, &donation.UpdatedAt,
			&donation.CampaignTitle, &donation.CampaignDescription, &donation.NGOName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan donation history: %v", err)
		}
		donations = append(donations, donation)
	}

	return donations, nil
}

// UpdateStatus updates donation status and transaction ID
func (r *donationRepository) UpdateStatus(id uuid.UUID, status string, transactionID string) error {
	query := `UPDATE donations SET status = $2, transaction_id = $3, updated_at = $4 WHERE id = $1`
	
	result, err := r.db.Exec(query, id, status, transactionID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update donation status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("donation not found")
	}

	return nil
}

// GetDonationsByDateRange retrieves donations within a date range
func (r *donationRepository) GetDonationsByDateRange(startDate, endDate time.Time, pagination PaginationParams) ([]*models.Donation, *PaginationResult, error) {
	countQuery := `SELECT COUNT(*) FROM donations WHERE created_at BETWEEN $1 AND $2`
	
	var totalCount int
	err := r.db.QueryRow(countQuery, startDate, endDate).Scan(&totalCount)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count donations: %v", err)
	}

	query := `
		SELECT 
			id, donor_id, campaign_id, amount, currency, status,
			payment_method, transaction_id, anonymous, message,
			receipt_email, created_at, updated_at
		FROM donations 
		WHERE created_at BETWEEN $1 AND $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.Query(query, startDate, endDate, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query donations: %v", err)
	}
	defer rows.Close()

	var donations []*models.Donation
	for rows.Next() {
		donation := &models.Donation{}
		err := rows.Scan(
			&donation.ID, &donation.DonorID, &donation.CampaignID,
			&donation.Amount, &donation.Currency, &donation.Status,
			&donation.PaymentMethod, &donation.TransactionID,
			&donation.Anonymous, &donation.Message, &donation.ReceiptEmail,
			&donation.CreatedAt, &donation.UpdatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan donation: %v", err)
		}
		donations = append(donations, donation)
	}

	result := &PaginationResult{
		TotalCount:  totalCount,
		TotalPages:  (totalCount + pagination.Limit - 1) / pagination.Limit,
		CurrentPage: pagination.Page,
		PageSize:    pagination.Limit,
	}

	return donations, result, nil
}

// buildWhereClause constructs WHERE clause and arguments for filtering
func (r *donationRepository) buildWhereClause(filters DonationFiltersExtended) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filters.DonorID != nil {
		conditions = append(conditions, fmt.Sprintf("d.donor_id = $%d", argIndex))
		args = append(args, *filters.DonorID)
		argIndex++
	}

	if filters.CampaignID != nil {
		conditions = append(conditions, fmt.Sprintf("d.campaign_id = $%d", argIndex))
		args = append(args, *filters.CampaignID)
		argIndex++
	}

	if filters.NGOID != nil {
		conditions = append(conditions, fmt.Sprintf("c.ngo_id = $%d", argIndex))
		args = append(args, *filters.NGOID)
		argIndex++
	}

	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("d.status = $%d", argIndex))
		args = append(args, filters.Status)
		argIndex++
	}

	if filters.MinAmount != nil {
		conditions = append(conditions, fmt.Sprintf("d.amount >= $%d", argIndex))
		args = append(args, *filters.MinAmount)
		argIndex++
	}

	if filters.MaxAmount != nil {
		conditions = append(conditions, fmt.Sprintf("d.amount <= $%d", argIndex))
		args = append(args, *filters.MaxAmount)
		argIndex++
	}

	if filters.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("d.created_at >= $%d", argIndex))
		args = append(args, *filters.StartDate)
		argIndex++
	}

	if filters.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("d.created_at <= $%d", argIndex))
		args = append(args, *filters.EndDate)
		argIndex++
	}

	if filters.Anonymous != nil {
		conditions = append(conditions, fmt.Sprintf("d.anonymous = $%d", argIndex))
		args = append(args, *filters.Anonymous)
		argIndex++
	}

	if filters.SearchTerm != "" {
		conditions = append(conditions, fmt.Sprintf("(c.title ILIKE $%d OR d.message ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+filters.SearchTerm+"%")
		argIndex++
	}

	if len(conditions) == 0 {
		return "", args
	}

	return "WHERE " + fmt.Sprintf("%s", conditions[0]) + 
		func() string {
			if len(conditions) > 1 {
				result := ""
				for i := 1; i < len(conditions); i++ {
					result += " AND " + conditions[i]
				}
				return result
			}
			return ""
		}(), args
}