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

// Custom errors for better error handling
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserExists        = errors.New("user with this email already exists")
	ErrUserInactive      = errors.New("user account is inactive")
	ErrUserNotVerified   = errors.New("user email not verified")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// UserRepository handles all user-related database operations
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(ctx context.Context, userCreate *models.UserCreate, passwordHash string) (*models.User, error) {
	// Check if user already exists
	exists, err := r.ExistsByEmail(ctx, userCreate.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, ErrUserExists
	}

	user := &models.User{
		Email:         strings.ToLower(strings.TrimSpace(userCreate.Email)),
		PasswordHash:  passwordHash,
		Name:          strings.TrimSpace(userCreate.Name),
		Role:          userCreate.Role,
		Phone:         stringPtr(strings.TrimSpace(userCreate.Phone)),
		Address:       stringPtr(strings.TrimSpace(userCreate.Address)),
		IsActive:      true,
		EmailVerified: false, // Default to false, require email verification
	}

	query := `
		INSERT INTO users (email, password_hash, name, role, phone, address, is_active, email_verified)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	
	err = r.db.QueryRowContext(ctx, query,
		user.Email, user.PasswordHash, user.Name, user.Role,
		user.Phone, user.Address, user.IsActive, user.EmailVerified,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		// Handle unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, ErrUserExists
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, name, role, avatar_url, phone, address, 
		       is_active, email_verified, created_at, updated_at
		FROM users WHERE id = $1
	`
	
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
		&user.AvatarURL, &user.Phone, &user.Address, &user.IsActive, 
		&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	
	return user, nil
}

// GetUserByEmail retrieves a user by their email (for login)
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	
	query := `
		SELECT id, email, password_hash, name, role, avatar_url, phone, address, 
		       is_active, email_verified, created_at, updated_at
		FROM users WHERE email = $1
	`
	
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
		&user.AvatarURL, &user.Phone, &user.Address, &user.IsActive, 
		&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	
	return user, nil
}

// GetUserForAuth retrieves a user for authentication and validates their status
func (r *UserRepository) GetUserForAuth(ctx context.Context, email string) (*models.User, error) {
	user, err := r.GetUserByEmail(ctx, email)
	if err != nil {
		if err == ErrUserNotFound {
			return nil, ErrInvalidCredentials // Don't reveal if email exists
		}
		return nil, err
	}
	
	// Check if account is active
	if !user.IsActive {
		return nil, ErrUserInactive
	}
	
	return user, nil
}

// UpdateUser updates user profile information
func (r *UserRepository) UpdateUser(ctx context.Context, id uuid.UUID, update *models.UserUpdate) (*models.User, error) {
	// Build dynamic query based on fields to update
	setParts := []string{}
	args := []interface{}{id}
	argIndex := 2

	if update.Name != "" {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, strings.TrimSpace(update.Name))
		argIndex++
	}
	
	if update.Phone != "" {
		setParts = append(setParts, fmt.Sprintf("phone = $%d", argIndex))
		args = append(args, stringPtr(strings.TrimSpace(update.Phone)))
		argIndex++
	}
	
	if update.Address != "" {
		setParts = append(setParts, fmt.Sprintf("address = $%d", argIndex))
		args = append(args, stringPtr(strings.TrimSpace(update.Address)))
		argIndex++
	}
	
	if update.AvatarURL != "" {
		setParts = append(setParts, fmt.Sprintf("avatar_url = $%d", argIndex))
		args = append(args, stringPtr(update.AvatarURL))
		argIndex++
	}

	if len(setParts) == 0 {
		// No fields to update, just return current user
		return r.GetUserByID(ctx, id)
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = NOW()"))
	
	query := fmt.Sprintf(`
		UPDATE users SET %s
		WHERE id = $1
		RETURNING id, email, password_hash, name, role, avatar_url, phone, address, 
		          is_active, email_verified, created_at, updated_at
	`, strings.Join(setParts, ", "))
	
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
		&user.AvatarURL, &user.Phone, &user.Address, &user.IsActive, 
		&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	
	return user, nil
}

// ChangePassword updates user's password
func (r *UserRepository) ChangePassword(ctx context.Context, id uuid.UUID, newPasswordHash string) error {
	query := `
		UPDATE users SET password_hash = $2, updated_at = NOW() 
		WHERE id = $1 AND is_active = true
	`
	
	result, err := r.db.ExecContext(ctx, query, id, newPasswordHash)
	if err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	
	return nil
}

// VerifyEmail marks user's email as verified
func (r *UserRepository) VerifyEmail(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users SET email_verified = true, updated_at = NOW() 
		WHERE id = $1 AND is_active = true
	`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	
	return nil
}

// DeactivateUser deactivates a user account
func (r *UserRepository) DeactivateUser(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users SET is_active = false, updated_at = NOW() 
		WHERE id = $1
	`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	
	return nil
}

// ReactivateUser reactivates a user account
func (r *UserRepository) ReactivateUser(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users SET is_active = true, updated_at = NOW() 
		WHERE id = $1
	`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to reactivate user: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	
	return nil
}

// ListUsers retrieves paginated list of users with optional filters
func (r *UserRepository) ListUsers(ctx context.Context, params models.Pagination, filters UserFilters) ([]*models.User, int64, error) {
	// Build WHERE clause
	whereClause, filterArgs := r.buildUserFilters(filters)
	
	// Get total count
	countQuery := "SELECT COUNT(*) FROM users" + whereClause
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, filterArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, email, password_hash, name, role, avatar_url, phone, address, 
		       is_active, email_verified, created_at, updated_at
		FROM users` + whereClause + `
		ORDER BY created_at DESC
		LIMIT $` + fmt.Sprintf("%d", len(filterArgs)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(filterArgs)+2)
	
	allArgs := append(filterArgs, params.Limit, params.Offset)
	
	rows, err := r.db.QueryContext(ctx, query, allArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
			&user.AvatarURL, &user.Phone, &user.Address, &user.IsActive, 
			&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating users: %w", err)
	}

	return users, total, nil
}

// DeleteUser permanently deletes a user (use with caution)
func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	
	return nil
}

// ExistsByEmail checks if a user exists with the given email
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	
	return exists, nil
}

// GetUserStats returns user statistics for dashboard
func (r *UserRepository) GetUserStats(ctx context.Context) (*UserStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_users,
			COUNT(CASE WHEN is_active = true THEN 1 END) as active_users,
			COUNT(CASE WHEN email_verified = true THEN 1 END) as verified_users,
			COUNT(CASE WHEN role = 'donor' THEN 1 END) as donors,
			COUNT(CASE WHEN role = 'ngo_admin' THEN 1 END) as ngo_admins,
			COUNT(CASE WHEN created_at >= NOW() - INTERVAL '30 days' THEN 1 END) as new_users_30d
		FROM users
	`
	
	stats := &UserStats{}
	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalUsers, &stats.ActiveUsers, &stats.VerifiedUsers,
		&stats.Donors, &stats.NGOAdmins, &stats.NewUsers30Days,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}
	
	return stats, nil
}

// GetRecentUsers returns recently registered users
func (r *UserRepository) GetRecentUsers(ctx context.Context, limit int) ([]*models.User, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	
	query := `
		SELECT id, email, password_hash, name, role, avatar_url, phone, address, 
		       is_active, email_verified, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1
	`
	
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
			&user.AvatarURL, &user.Phone, &user.Address, &user.IsActive, 
			&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recent user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating recent users: %w", err)
	}

	return users, nil
}

// UserFilters represents filters for user queries
type UserFilters struct {
	Role         string
	IsActive     *bool
	IsVerified   *bool
	SearchTerm   string // Search in name and email
	CreatedAfter *time.Time
	CreatedBefore *time.Time
}

// UserStats represents user statistics
type UserStats struct {
	TotalUsers     int `json:"total_users"`
	ActiveUsers    int `json:"active_users"`
	VerifiedUsers  int `json:"verified_users"`
	Donors         int `json:"donors"`
	NGOAdmins      int `json:"ngo_admins"`
	NewUsers30Days int `json:"new_users_30_days"`
}

// Helper functions

// buildUserFilters builds WHERE clause and arguments for user filters
func (r *UserRepository) buildUserFilters(filters UserFilters) (string, []interface{}) {
	conditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if filters.Role != "" {
		conditions = append(conditions, fmt.Sprintf("role = $%d", argIndex))
		args = append(args, filters.Role)
		argIndex++
	}

	if filters.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *filters.IsActive)
		argIndex++
	}

	if filters.IsVerified != nil {
		conditions = append(conditions, fmt.Sprintf("email_verified = $%d", argIndex))
		args = append(args, *filters.IsVerified)
		argIndex++
	}

	if filters.SearchTerm != "" {
		searchTerm := "%" + strings.ToLower(filters.SearchTerm) + "%"
		conditions = append(conditions, fmt.Sprintf("(LOWER(name) LIKE $%d OR LOWER(email) LIKE $%d)", argIndex, argIndex))
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

// stringPtr returns a pointer to string, or nil if string is empty
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}