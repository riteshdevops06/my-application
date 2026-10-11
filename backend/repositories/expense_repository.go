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

// ExpenseFilters represents filtering options for expense queries
type ExpenseFilters struct {
	CampaignID   *uuid.UUID `json:"campaign_id,omitempty"`
	MilestoneID  *uuid.UUID `json:"milestone_id,omitempty"`
	Status       string     `json:"status,omitempty"`
	Category     string     `json:"category,omitempty"`
	MinAmount    *float64   `json:"min_amount,omitempty"`
	MaxAmount    *float64   `json:"max_amount,omitempty"`
	StartDate    *time.Time `json:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	ApprovedBy   *uuid.UUID `json:"approved_by,omitempty"`
	SearchTerm   string     `json:"search_term,omitempty"`
}

// ExpenseRepository interface defines methods for expense data operations
type ExpenseRepository interface {
	// Enhanced methods (main interface)
	CreateExpense(expense *models.Expense) error
	GetExpenseByID(id uuid.UUID) (*models.Expense, error)
	UpdateExpense(expense *models.Expense) error
	DeleteExpense(id uuid.UUID) error
	ListExpensesByCampaign(campaignID uuid.UUID, pagination PaginationParams) ([]*models.Expense, *PaginationResult, error)
	ListExpensesByMilestone(milestoneID uuid.UUID) ([]*models.Expense, error)
	UpdateExpenseStatus(id uuid.UUID, status string, approvedBy uuid.UUID) error
	ListExpenses(filters ExpenseFilters, pagination PaginationParams) ([]*models.Expense, *PaginationResult, error)
	GetExpenseStats(campaignID uuid.UUID) (*ExpenseStats, error)
	GetExpensesByStatus(status string, pagination PaginationParams) ([]*models.Expense, *PaginationResult, error)
	ApproveExpense(id uuid.UUID, approvedBy uuid.UUID) error
	RejectExpense(id uuid.UUID, approvedBy uuid.UUID, reason string) error
	
	// Interface compatibility methods (ExpenseRepositoryInterface)
	Create(ctx context.Context, expense *models.Expense) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Expense, error)
	Update(ctx context.Context, expense *models.Expense) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByCampaign(ctx context.Context, campaignID uuid.UUID, params models.Pagination) ([]*models.Expense, int64, error)
	ListByMilestone(ctx context.Context, milestoneID uuid.UUID) ([]*models.Expense, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, approvedBy uuid.UUID) error
}

// ExpenseStats represents aggregated expense statistics
type ExpenseStats struct {
	TotalExpenses     int     `json:"total_expenses"`
	PendingExpenses   int     `json:"pending_expenses"`
	ApprovedExpenses  int     `json:"approved_expenses"`
	RejectedExpenses  int     `json:"rejected_expenses"`
	TotalAmount       float64 `json:"total_amount"`
	ApprovedAmount    float64 `json:"approved_amount"`
	PendingAmount     float64 `json:"pending_amount"`
	RejectedAmount    float64 `json:"rejected_amount"`
	AverageAmount     float64 `json:"average_amount"`
	ExpensesByCategory map[string]float64 `json:"expenses_by_category"`
}

// ExpenseWithDetails represents expense with campaign and milestone details
type ExpenseWithDetails struct {
	*models.Expense
	CampaignTitle   string `json:"campaign_title"`
	MilestoneTitle  string `json:"milestone_title,omitempty"`
	ApproverName    string `json:"approver_name,omitempty"`
	NGOName         string `json:"ngo_name"`
}

// expenseRepository implements ExpenseRepository interface
type expenseRepository struct {
	db *sql.DB
}

// NewExpenseRepository creates a new instance of expense repository
func NewExpenseRepository(db *sql.DB) ExpenseRepository {
	return &expenseRepository{db: db}
}

// Interface compatibility methods

func (r *expenseRepository) Create(ctx context.Context, expense *models.Expense) error {
	return r.CreateExpense(expense)
}

func (r *expenseRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Expense, error) {
	return r.GetExpenseByID(id)
}

func (r *expenseRepository) Update(ctx context.Context, expense *models.Expense) error {
	return r.UpdateExpense(expense)
}

func (r *expenseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DeleteExpense(id)
}

func (r *expenseRepository) ListByCampaign(ctx context.Context, campaignID uuid.UUID, params models.Pagination) ([]*models.Expense, int64, error) {
	pagination := PaginationParams{
		Page:   params.Page,
		Limit:  params.Limit,
		Offset: (params.Page - 1) * params.Limit,
	}

	expenses, result, err := r.ListExpensesByCampaign(campaignID, pagination)
	if err != nil {
		return nil, 0, err
	}

	return expenses, int64(result.TotalCount), nil
}

func (r *expenseRepository) ListByMilestone(ctx context.Context, milestoneID uuid.UUID) ([]*models.Expense, error) {
	return r.ListExpensesByMilestone(milestoneID)
}

func (r *expenseRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, approvedBy uuid.UUID) error {
	return r.UpdateExpenseStatus(id, status, approvedBy)
}

// Enhanced methods

// CreateExpense creates a new expense record
func (r *expenseRepository) CreateExpense(expense *models.Expense) error {
	query := `
		INSERT INTO expenses (
			id, campaign_id, milestone_id, title, description, amount, currency,
			category, receipt_url, expense_date, status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)`

	_, err := r.db.Exec(query,
		expense.ID, expense.CampaignID, expense.MilestoneID, expense.Title,
		expense.Description, expense.Amount, expense.Currency, expense.Category,
		expense.ReceiptURL, expense.ExpenseDate, expense.Status,
		expense.CreatedAt, expense.UpdatedAt)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503": // Foreign key violation
				return fmt.Errorf("invalid campaign or milestone ID")
			case "23505": // Unique violation
				return fmt.Errorf("expense with this title already exists for this campaign")
			default:
				return fmt.Errorf("failed to create expense: %v", err)
			}
		}
		return fmt.Errorf("failed to create expense: %v", err)
	}

	return nil
}

// GetExpenseByID retrieves an expense by its ID
func (r *expenseRepository) GetExpenseByID(id uuid.UUID) (*models.Expense, error) {
	query := `
		SELECT 
			id, campaign_id, milestone_id, title, description, amount, currency,
			category, receipt_url, expense_date, status, approved_by, approved_at,
			created_at, updated_at
		FROM expenses 
		WHERE id = $1`

	expense := &models.Expense{}
	err := r.db.QueryRow(query, id).Scan(
		&expense.ID, &expense.CampaignID, &expense.MilestoneID, &expense.Title,
		&expense.Description, &expense.Amount, &expense.Currency, &expense.Category,
		&expense.ReceiptURL, &expense.ExpenseDate, &expense.Status,
		&expense.ApprovedBy, &expense.ApprovedAt, &expense.CreatedAt, &expense.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("expense not found")
		}
		return nil, fmt.Errorf("failed to get expense: %v", err)
	}

	return expense, nil
}

// UpdateExpense updates an existing expense
func (r *expenseRepository) UpdateExpense(expense *models.Expense) error {
	query := `
		UPDATE expenses SET
			title = $2, description = $3, amount = $4, currency = $5, category = $6,
			receipt_url = $7, expense_date = $8, status = $9, approved_by = $10,
			approved_at = $11, updated_at = $12
		WHERE id = $1`

	result, err := r.db.Exec(query,
		expense.ID, expense.Title, expense.Description, expense.Amount,
		expense.Currency, expense.Category, expense.ReceiptURL, expense.ExpenseDate,
		expense.Status, expense.ApprovedBy, expense.ApprovedAt, expense.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update expense: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("expense not found")
	}

	return nil
}

// DeleteExpense deletes an expense (soft delete - updates status to cancelled)
func (r *expenseRepository) DeleteExpense(id uuid.UUID) error {
	query := `UPDATE expenses SET status = 'cancelled', updated_at = $2 WHERE id = $1`
	
	result, err := r.db.Exec(query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete expense: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("expense not found")
	}

	return nil
}

// ListExpensesByCampaign retrieves expenses by campaign ID with pagination
func (r *expenseRepository) ListExpensesByCampaign(campaignID uuid.UUID, pagination PaginationParams) ([]*models.Expense, *PaginationResult, error) {
	countQuery := `SELECT COUNT(*) FROM expenses WHERE campaign_id = $1 AND status != 'cancelled'`
	
	var totalCount int
	err := r.db.QueryRow(countQuery, campaignID).Scan(&totalCount)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count expenses: %v", err)
	}

	query := `
		SELECT 
			id, campaign_id, milestone_id, title, description, amount, currency,
			category, receipt_url, expense_date, status, approved_by, approved_at,
			created_at, updated_at
		FROM expenses 
		WHERE campaign_id = $1 AND status != 'cancelled'
		ORDER BY expense_date DESC, created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, campaignID, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query expenses: %v", err)
	}
	defer rows.Close()

	var expenses []*models.Expense
	for rows.Next() {
		expense := &models.Expense{}
		err := rows.Scan(
			&expense.ID, &expense.CampaignID, &expense.MilestoneID, &expense.Title,
			&expense.Description, &expense.Amount, &expense.Currency, &expense.Category,
			&expense.ReceiptURL, &expense.ExpenseDate, &expense.Status,
			&expense.ApprovedBy, &expense.ApprovedAt, &expense.CreatedAt, &expense.UpdatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan expense: %v", err)
		}
		expenses = append(expenses, expense)
	}

	result := &PaginationResult{
		TotalCount:  totalCount,
		TotalPages:  (totalCount + pagination.Limit - 1) / pagination.Limit,
		CurrentPage: pagination.Page,
		PageSize:    pagination.Limit,
	}

	return expenses, result, nil
}

// ListExpensesByMilestone retrieves all expenses for a milestone
func (r *expenseRepository) ListExpensesByMilestone(milestoneID uuid.UUID) ([]*models.Expense, error) {
	query := `
		SELECT 
			id, campaign_id, milestone_id, title, description, amount, currency,
			category, receipt_url, expense_date, status, approved_by, approved_at,
			created_at, updated_at
		FROM expenses 
		WHERE milestone_id = $1 AND status != 'cancelled'
		ORDER BY expense_date DESC, created_at DESC`

	rows, err := r.db.Query(query, milestoneID)
	if err != nil {
		return nil, fmt.Errorf("failed to query expenses: %v", err)
	}
	defer rows.Close()

	var expenses []*models.Expense
	for rows.Next() {
		expense := &models.Expense{}
		err := rows.Scan(
			&expense.ID, &expense.CampaignID, &expense.MilestoneID, &expense.Title,
			&expense.Description, &expense.Amount, &expense.Currency, &expense.Category,
			&expense.ReceiptURL, &expense.ExpenseDate, &expense.Status,
			&expense.ApprovedBy, &expense.ApprovedAt, &expense.CreatedAt, &expense.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan expense: %v", err)
		}
		expenses = append(expenses, expense)
	}

	return expenses, nil
}

// UpdateExpenseStatus updates expense status and approver
func (r *expenseRepository) UpdateExpenseStatus(id uuid.UUID, status string, approvedBy uuid.UUID) error {
	query := `
		UPDATE expenses SET
			status = $2, approved_by = $3, approved_at = $4, updated_at = $5
		WHERE id = $1`
	
	approvedAt := time.Now()
	result, err := r.db.Exec(query, id, status, approvedBy, approvedAt, approvedAt)
	if err != nil {
		return fmt.Errorf("failed to update expense status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("expense not found")
	}

	return nil
}

// ListExpenses retrieves expenses based on filters and pagination
func (r *expenseRepository) ListExpenses(filters ExpenseFilters, pagination PaginationParams) ([]*models.Expense, *PaginationResult, error) {
	whereClause, args := r.buildWhereClause(filters)
	
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM expenses e %s", whereClause)
	
	var totalCount int
	err := r.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count expenses: %v", err)
	}

	query := fmt.Sprintf(`
		SELECT 
			e.id, e.campaign_id, e.milestone_id, e.title, e.description, e.amount,
			e.currency, e.category, e.receipt_url, e.expense_date, e.status,
			e.approved_by, e.approved_at, e.created_at, e.updated_at
		FROM expenses e
		%s
		ORDER BY e.expense_date DESC, e.created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pagination.Limit, pagination.Offset)
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query expenses: %v", err)
	}
	defer rows.Close()

	var expenses []*models.Expense
	for rows.Next() {
		expense := &models.Expense{}
		err := rows.Scan(
			&expense.ID, &expense.CampaignID, &expense.MilestoneID, &expense.Title,
			&expense.Description, &expense.Amount, &expense.Currency, &expense.Category,
			&expense.ReceiptURL, &expense.ExpenseDate, &expense.Status,
			&expense.ApprovedBy, &expense.ApprovedAt, &expense.CreatedAt, &expense.UpdatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan expense: %v", err)
		}
		expenses = append(expenses, expense)
	}

	result := &PaginationResult{
		TotalCount:  totalCount,
		TotalPages:  (totalCount + pagination.Limit - 1) / pagination.Limit,
		CurrentPage: pagination.Page,
		PageSize:    pagination.Limit,
	}

	return expenses, result, nil
}

// GetExpenseStats retrieves aggregated expense statistics for a campaign
func (r *expenseRepository) GetExpenseStats(campaignID uuid.UUID) (*ExpenseStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_expenses,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_expenses,
			COUNT(CASE WHEN status = 'approved' THEN 1 END) as approved_expenses,
			COUNT(CASE WHEN status = 'rejected' THEN 1 END) as rejected_expenses,
			COALESCE(SUM(amount), 0) as total_amount,
			COALESCE(SUM(CASE WHEN status = 'approved' THEN amount ELSE 0 END), 0) as approved_amount,
			COALESCE(SUM(CASE WHEN status = 'pending' THEN amount ELSE 0 END), 0) as pending_amount,
			COALESCE(SUM(CASE WHEN status = 'rejected' THEN amount ELSE 0 END), 0) as rejected_amount,
			COALESCE(AVG(amount), 0) as average_amount
		FROM expenses 
		WHERE campaign_id = $1 AND status != 'cancelled'`

	stats := &ExpenseStats{}
	err := r.db.QueryRow(query, campaignID).Scan(
		&stats.TotalExpenses, &stats.PendingExpenses, &stats.ApprovedExpenses,
		&stats.RejectedExpenses, &stats.TotalAmount, &stats.ApprovedAmount,
		&stats.PendingAmount, &stats.RejectedAmount, &stats.AverageAmount)

	if err != nil {
		return nil, fmt.Errorf("failed to get expense stats: %v", err)
	}

	// Get expenses by category
	categoryQuery := `
		SELECT COALESCE(category, 'Uncategorized') as category, SUM(amount) as amount
		FROM expenses 
		WHERE campaign_id = $1 AND status = 'approved'
		GROUP BY COALESCE(category, 'Uncategorized')`

	rows, err := r.db.Query(categoryQuery, campaignID)
	if err != nil {
		return stats, nil // Return stats without category breakdown if query fails
	}
	defer rows.Close()

	stats.ExpensesByCategory = make(map[string]float64)
	for rows.Next() {
		var category string
		var amount float64
		if err := rows.Scan(&category, &amount); err == nil {
			stats.ExpensesByCategory[category] = amount
		}
	}

	return stats, nil
}

// GetExpensesByStatus retrieves expenses by status with pagination
func (r *expenseRepository) GetExpensesByStatus(status string, pagination PaginationParams) ([]*models.Expense, *PaginationResult, error) {
	countQuery := `SELECT COUNT(*) FROM expenses WHERE status = $1`
	
	var totalCount int
	err := r.db.QueryRow(countQuery, status).Scan(&totalCount)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count expenses: %v", err)
	}

	query := `
		SELECT 
			id, campaign_id, milestone_id, title, description, amount, currency,
			category, receipt_url, expense_date, status, approved_by, approved_at,
			created_at, updated_at
		FROM expenses 
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, status, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query expenses: %v", err)
	}
	defer rows.Close()

	var expenses []*models.Expense
	for rows.Next() {
		expense := &models.Expense{}
		err := rows.Scan(
			&expense.ID, &expense.CampaignID, &expense.MilestoneID, &expense.Title,
			&expense.Description, &expense.Amount, &expense.Currency, &expense.Category,
			&expense.ReceiptURL, &expense.ExpenseDate, &expense.Status,
			&expense.ApprovedBy, &expense.ApprovedAt, &expense.CreatedAt, &expense.UpdatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan expense: %v", err)
		}
		expenses = append(expenses, expense)
	}

	result := &PaginationResult{
		TotalCount:  totalCount,
		TotalPages:  (totalCount + pagination.Limit - 1) / pagination.Limit,
		CurrentPage: pagination.Page,
		PageSize:    pagination.Limit,
	}

	return expenses, result, nil
}

// ApproveExpense approves an expense
func (r *expenseRepository) ApproveExpense(id uuid.UUID, approvedBy uuid.UUID) error {
	return r.UpdateExpenseStatus(id, "approved", approvedBy)
}

// RejectExpense rejects an expense with reason
func (r *expenseRepository) RejectExpense(id uuid.UUID, approvedBy uuid.UUID, reason string) error {
	query := `
		UPDATE expenses SET
			status = 'rejected', approved_by = $2, approved_at = $3, 
			description = CONCAT(COALESCE(description, ''), '\n\nRejection Reason: ', $4),
			updated_at = $5
		WHERE id = $1`
	
	approvedAt := time.Now()
	result, err := r.db.Exec(query, id, approvedBy, approvedAt, reason, approvedAt)
	if err != nil {
		return fmt.Errorf("failed to reject expense: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("expense not found")
	}

	return nil
}

// buildWhereClause constructs WHERE clause and arguments for filtering
func (r *expenseRepository) buildWhereClause(filters ExpenseFilters) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Always exclude cancelled expenses
	conditions = append(conditions, "e.status != 'cancelled'")

	if filters.CampaignID != nil {
		conditions = append(conditions, fmt.Sprintf("e.campaign_id = $%d", argIndex))
		args = append(args, *filters.CampaignID)
		argIndex++
	}

	if filters.MilestoneID != nil {
		conditions = append(conditions, fmt.Sprintf("e.milestone_id = $%d", argIndex))
		args = append(args, *filters.MilestoneID)
		argIndex++
	}

	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("e.status = $%d", argIndex))
		args = append(args, filters.Status)
		argIndex++
	}

	if filters.Category != "" {
		conditions = append(conditions, fmt.Sprintf("e.category = $%d", argIndex))
		args = append(args, filters.Category)
		argIndex++
	}

	if filters.MinAmount != nil {
		conditions = append(conditions, fmt.Sprintf("e.amount >= $%d", argIndex))
		args = append(args, *filters.MinAmount)
		argIndex++
	}

	if filters.MaxAmount != nil {
		conditions = append(conditions, fmt.Sprintf("e.amount <= $%d", argIndex))
		args = append(args, *filters.MaxAmount)
		argIndex++
	}

	if filters.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("e.expense_date >= $%d", argIndex))
		args = append(args, *filters.StartDate)
		argIndex++
	}

	if filters.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("e.expense_date <= $%d", argIndex))
		args = append(args, *filters.EndDate)
		argIndex++
	}

	if filters.ApprovedBy != nil {
		conditions = append(conditions, fmt.Sprintf("e.approved_by = $%d", argIndex))
		args = append(args, *filters.ApprovedBy)
		argIndex++
	}

	if filters.SearchTerm != "" {
		conditions = append(conditions, fmt.Sprintf("(e.title ILIKE $%d OR e.description ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+filters.SearchTerm+"%")
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			whereClause += " AND " + conditions[i]
		}
	}

	return whereClause, args
}