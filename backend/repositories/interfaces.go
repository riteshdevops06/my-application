package repositories

import (
	"context"
	"time"

	"ngo-donation-tracker/models"

	"github.com/google/uuid"
)

// UserRepositoryInterface defines the contract for user repository
type UserRepositoryInterface interface {
	// User creation and authentication
	CreateUser(ctx context.Context, userCreate *models.UserCreate, passwordHash string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserForAuth(ctx context.Context, email string) (*models.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// User profile management
	UpdateUser(ctx context.Context, id uuid.UUID, update *models.UserUpdate) (*models.User, error)
	ChangePassword(ctx context.Context, id uuid.UUID, newPasswordHash string) error
	VerifyEmail(ctx context.Context, id uuid.UUID) error

	// User status management
	DeactivateUser(ctx context.Context, id uuid.UUID) error
	ReactivateUser(ctx context.Context, id uuid.UUID) error
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// User listing and statistics
	ListUsers(ctx context.Context, params models.Pagination, filters UserFilters) ([]*models.User, int64, error)
	GetUserStats(ctx context.Context) (*UserStats, error)
	GetRecentUsers(ctx context.Context, limit int) ([]*models.User, error)
}

// NGORepositoryInterface defines the contract for NGO repository
type NGORepositoryInterface interface {
	Create(ctx context.Context, ngo *models.NGO) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.NGO, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.NGO, error)
	Update(ctx context.Context, ngo *models.NGO) error
	List(ctx context.Context, params models.Pagination, status string) ([]*models.NGO, int64, error)
	UpdateVerificationStatus(ctx context.Context, id uuid.UUID, status string, notes *string, verifiedBy *uuid.UUID) error
}

// CampaignRepositoryInterface defines the contract for campaign repository
type CampaignRepositoryInterface interface {
	Create(ctx context.Context, campaign *models.Campaign) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Campaign, error)
	Update(ctx context.Context, campaign *models.Campaign) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, params models.Pagination, filters CampaignFilters) ([]*models.Campaign, int64, error)
	ListByNGO(ctx context.Context, ngoID uuid.UUID, params models.Pagination) ([]*models.Campaign, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	UpdateRaisedAmount(ctx context.Context, id uuid.UUID, amount float64) error
}

// DonationRepositoryInterface defines the contract for donation repository
type DonationRepositoryInterface interface {
	Create(ctx context.Context, donation *models.Donation) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Donation, error)
	List(ctx context.Context, params models.Pagination, filters DonationFilters) ([]*models.Donation, int64, error)
	ListByCampaign(ctx context.Context, campaignID uuid.UUID, params models.Pagination) ([]*models.Donation, int64, error)
	ListByDonor(ctx context.Context, donorID uuid.UUID, params models.Pagination) ([]*models.Donation, int64, error)
	UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status string, transactionID *string) error
}

// MilestoneRepositoryInterface defines the contract for milestone repository
type MilestoneRepositoryInterface interface {
	Create(ctx context.Context, milestone *models.Milestone) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Milestone, error)
	Update(ctx context.Context, milestone *models.Milestone) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByCampaign(ctx context.Context, campaignID uuid.UUID) ([]*models.Milestone, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	MarkComplete(ctx context.Context, id uuid.UUID, proofURL, proofNotes string) error
}

// ExpenseRepositoryInterface defines the contract for expense repository
type ExpenseRepositoryInterface interface {
	Create(ctx context.Context, expense *models.Expense) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Expense, error)
	Update(ctx context.Context, expense *models.Expense) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByCampaign(ctx context.Context, campaignID uuid.UUID, params models.Pagination) ([]*models.Expense, int64, error)
	ListByMilestone(ctx context.Context, milestoneID uuid.UUID) ([]*models.Expense, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, approvedBy uuid.UUID) error
}

// VerificationRequestRepositoryInterface defines the contract for verification repository
type VerificationRequestRepositoryInterface interface {
	Create(ctx context.Context, request *models.VerificationRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.VerificationRequest, error)
	GetByNGOID(ctx context.Context, ngoID uuid.UUID) (*models.VerificationRequest, error)
	List(ctx context.Context, params models.Pagination, status string) ([]*models.VerificationRequest, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status, adminNotes string, reviewedBy uuid.UUID) error
}

// Filter structs for repository interfaces

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

// DonationFilters represents filters for donation queries
type DonationFilters struct {
	CampaignID      *uuid.UUID
	DonorID         *uuid.UUID
	PaymentStatus   string
	MinAmount       *float64
	MaxAmount       *float64
	DonatedAfter    *time.Time
	DonatedBefore   *time.Time
	IsAnonymous     *bool
}