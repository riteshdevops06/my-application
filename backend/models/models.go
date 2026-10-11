package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Claims represents JWT claims
type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	NGOID    *uuid.UUID `json:"ngo_id,omitempty"`
	jwt.RegisteredClaims
}

// User represents a platform user (donor, ngo_admin, or platform_admin)
type User struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	Email          string     `json:"email" db:"email" validate:"required,email"`
	PasswordHash   string     `json:"-" db:"password_hash"` // Never expose password hash
	Name           string     `json:"name" db:"name" validate:"required,min=2,max=255"`
	Role           string     `json:"role" db:"role" validate:"required,oneof=donor ngo_admin platform_admin"`
	AvatarURL      *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	Phone          *string    `json:"phone,omitempty" db:"phone"`
	Address        *string    `json:"address,omitempty" db:"address"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	EmailVerified  bool       `json:"email_verified" db:"email_verified"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// UserCreate represents the request body for user registration
type UserCreate struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2,max=255"`
	Role     string `json:"role" validate:"required,oneof=donor ngo_admin platform_admin"`
	Phone    string `json:"phone,omitempty"`
	Address  string `json:"address,omitempty"`
}

// UserLogin represents the login request
type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserUpdate represents the user profile update request
type UserUpdate struct {
	Name      string `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// PasswordChange represents a password change request
type PasswordChange struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// NGO represents a non-governmental organization
type NGO struct {
	ID                       uuid.UUID  `json:"id" db:"id"`
	UserID                   uuid.UUID  `json:"user_id" db:"user_id"`
	Name                     string     `json:"name" db:"name" validate:"required,min=2,max=255"`
	Description              *string    `json:"description,omitempty" db:"description"`
	Mission                  *string    `json:"mission,omitempty" db:"mission"`
	Website                  *string    `json:"website,omitempty" db:"website" validate:"omitempty,url"`
	RegistrationNumber       *string    `json:"registration_number,omitempty" db:"registration_number"`
	RegistrationDocumentURL  *string    `json:"registration_document_url,omitempty" db:"registration_document_url"`
	LogoURL                  *string    `json:"logo_url,omitempty" db:"logo_url"`
	CoverImageURL            *string    `json:"cover_image_url,omitempty" db:"cover_image_url"`
	Address                  *string    `json:"address,omitempty" db:"address"`
	City                     *string    `json:"city,omitempty" db:"city"`
	Country                  *string    `json:"country,omitempty" db:"country"`
	VerificationStatus       string     `json:"verification_status" db:"verification_status"`
	VerificationNotes        *string    `json:"verification_notes,omitempty" db:"verification_notes"`
	VerifiedAt               *time.Time `json:"verified_at,omitempty" db:"verified_at"`
	VerifiedBy               *uuid.UUID `json:"verified_by,omitempty" db:"verified_by"`
	CreatedAt                time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at" db:"updated_at"`
	
	// Relations
	User     *User `json:"user,omitempty"`
	Verifier *User `json:"verifier,omitempty"`
}

// NGOCreate represents the request to create/update NGO profile
type NGOCreate struct {
	Name                    string `json:"name" validate:"required,min=2,max=255"`
	Description             string `json:"description,omitempty"`
	Mission                 string `json:"mission,omitempty"`
	Website                 string `json:"website,omitempty" validate:"omitempty,url"`
	RegistrationNumber      string `json:"registration_number,omitempty"`
	RegistrationDocumentURL string `json:"registration_document_url,omitempty"`
	LogoURL                 string `json:"logo_url,omitempty"`
	CoverImageURL           string `json:"cover_image_url,omitempty"`
	Address                 string `json:"address,omitempty"`
	City                    string `json:"city,omitempty"`
	Country                 string `json:"country,omitempty"`
}

// StringArray is a custom type to handle PostgreSQL text arrays
type StringArray []string

// Scan implements the sql.Scanner interface for StringArray
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	case pq.StringArray:
		*s = StringArray(v)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into StringArray", value)
	}
}

// Value implements the driver.Valuer interface for StringArray
func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return pq.StringArray(s).Value()
}

// Campaign represents a fundraising campaign
type Campaign struct {
	ID               uuid.UUID    `json:"id" db:"id"`
	NGOID            uuid.UUID    `json:"ngo_id" db:"ngo_id"`
	Title            string       `json:"title" db:"title" validate:"required,min=5,max=255"`
	Description      *string      `json:"description,omitempty" db:"description"`
	ShortDescription *string      `json:"short_description,omitempty" db:"short_description"`
	GoalAmount       float64      `json:"goal_amount" db:"goal_amount" validate:"required,gt=0"`
	RaisedAmount     float64      `json:"raised_amount" db:"raised_amount"`
	Currency         string       `json:"currency" db:"currency"`
	Status           string       `json:"status" db:"status"`
	Category         *string      `json:"category,omitempty" db:"category"`
	Tags             StringArray  `json:"tags,omitempty" db:"tags"`
	StartDate        *time.Time   `json:"start_date,omitempty" db:"start_date"`
	EndDate          *time.Time   `json:"end_date,omitempty" db:"end_date"`
	ImageURL         *string      `json:"image_url,omitempty" db:"image_url"`
	GalleryImages    StringArray  `json:"gallery_images,omitempty" db:"gallery_images"`
	Location         *string      `json:"location,omitempty" db:"location"`
	IsFeatured       bool         `json:"is_featured" db:"is_featured"`
	CreatedAt        time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at" db:"updated_at"`
	
	// Relations
	NGO        *NGO         `json:"ngo,omitempty"`
	Donations  []Donation   `json:"donations,omitempty"`
	Milestones []Milestone  `json:"milestones,omitempty"`
	Expenses   []Expense    `json:"expenses,omitempty"`
}

// CampaignCreate represents the request to create a campaign
type CampaignCreate struct {
	Title            string    `json:"title" validate:"required,min=5,max=255"`
	Description      string    `json:"description,omitempty"`
	ShortDescription string    `json:"short_description,omitempty"`
	GoalAmount       float64   `json:"goal_amount" validate:"required,gt=0"`
	Currency         string    `json:"currency,omitempty"`
	Category         string    `json:"category,omitempty"`
	Tags             []string  `json:"tags,omitempty"`
	StartDate        string    `json:"start_date,omitempty"`
	EndDate          string    `json:"end_date,omitempty"`
	ImageURL         string    `json:"image_url,omitempty"`
	GalleryImages    []string  `json:"gallery_images,omitempty"`
	Location         string    `json:"location,omitempty"`
}

// CampaignUpdate represents the request to update a campaign
type CampaignUpdate struct {
	Title            string   `json:"title,omitempty" validate:"omitempty,min=5,max=255"`
	Description      string   `json:"description,omitempty"`
	ShortDescription string   `json:"short_description,omitempty"`
	GoalAmount       float64  `json:"goal_amount,omitempty" validate:"omitempty,gt=0"`
	Category         string   `json:"category,omitempty"`
	Tags             []string `json:"tags,omitempty"`
	StartDate        string   `json:"start_date,omitempty"`
	EndDate          string   `json:"end_date,omitempty"`
	ImageURL         string   `json:"image_url,omitempty"`
	GalleryImages    []string `json:"gallery_images,omitempty"`
	Location         string   `json:"location,omitempty"`
}

// Donation represents a donation made to a campaign
type Donation struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	CampaignID       uuid.UUID  `json:"campaign_id" db:"campaign_id"`
	DonorID          uuid.UUID  `json:"donor_id" db:"donor_id"`
	Amount           float64    `json:"amount" db:"amount" validate:"required,gt=0"`
	Currency         string     `json:"currency" db:"currency"`
	PaymentMethod    *string    `json:"payment_method,omitempty" db:"payment_method"`
	PaymentReference *string    `json:"payment_reference,omitempty" db:"payment_reference"`
	PaymentStatus    string     `json:"payment_status" db:"payment_status"`
	TransactionID    *string    `json:"transaction_id,omitempty" db:"transaction_id"`
	IsAnonymous      bool       `json:"is_anonymous" db:"is_anonymous"`
	Message          *string    `json:"message,omitempty" db:"message"`
	DonatedAt        time.Time  `json:"donated_at" db:"donated_at"`
	CompletedAt      *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	
	// Relations
	Campaign *Campaign `json:"campaign,omitempty"`
	Donor    *User     `json:"donor,omitempty"`
}

// DonationCreate represents the request to create a donation
type DonationCreate struct {
	CampaignID    string  `json:"campaign_id" validate:"required,uuid"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	Currency      string  `json:"currency,omitempty"`
	PaymentMethod string  `json:"payment_method,omitempty"`
	IsAnonymous   bool    `json:"is_anonymous"`
	Message       string  `json:"message,omitempty"`
}

// Milestone represents a campaign milestone
type Milestone struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	CampaignID        uuid.UUID  `json:"campaign_id" db:"campaign_id"`
	Title             string     `json:"title" db:"title" validate:"required,min=5,max=255"`
	Description       *string    `json:"description,omitempty" db:"description"`
	TargetAmount      float64    `json:"target_amount" db:"target_amount" validate:"required,gt=0"`
	RaisedAmount      float64    `json:"raised_amount" db:"raised_amount"`
	Status            string     `json:"status" db:"status"`
	DueDate           *time.Time `json:"due_date,omitempty" db:"due_date"`
	CompletedAt       *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	ProofDocumentURL  *string    `json:"proof_document_url,omitempty" db:"proof_document_url"`
	ProofNotes        *string    `json:"proof_notes,omitempty" db:"proof_notes"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
	
	// Relations
	Campaign *Campaign `json:"campaign,omitempty"`
	Expenses []Expense `json:"expenses,omitempty"`
}

// MilestoneCreate represents the request to create a milestone
type MilestoneCreate struct {
	Title        string `json:"title" validate:"required,min=5,max=255"`
	Description  string `json:"description,omitempty"`
	TargetAmount float64 `json:"target_amount" validate:"required,gt=0"`
	DueDate      string `json:"due_date,omitempty"`
}

// MilestoneUpdate represents the request to update a milestone
type MilestoneUpdate struct {
	Title             string `json:"title,omitempty" validate:"omitempty,min=5,max=255"`
	Description       string `json:"description,omitempty"`
	TargetAmount      float64 `json:"target_amount,omitempty" validate:"omitempty,gt=0"`
	DueDate           string `json:"due_date,omitempty"`
	ProofDocumentURL  string `json:"proof_document_url,omitempty"`
	ProofNotes        string `json:"proof_notes,omitempty"`
}

// Expense represents an expense made for a campaign
type Expense struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	CampaignID  uuid.UUID  `json:"campaign_id" db:"campaign_id"`
	MilestoneID *uuid.UUID `json:"milestone_id,omitempty" db:"milestone_id"`
	Title       string     `json:"title" db:"title" validate:"required,min=3,max=255"`
	Description *string    `json:"description,omitempty" db:"description"`
	Amount      float64    `json:"amount" db:"amount" validate:"required,gt=0"`
	Currency    string     `json:"currency" db:"currency"`
	Category    *string    `json:"category,omitempty" db:"category"`
	ReceiptURL  *string    `json:"receipt_url,omitempty" db:"receipt_url"`
	ExpenseDate time.Time  `json:"expense_date" db:"expense_date"`
	Status      string     `json:"status" db:"status"`
	ApprovedBy  *uuid.UUID `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty" db:"approved_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	
	// Relations
	Campaign  *Campaign  `json:"campaign,omitempty"`
	Milestone *Milestone `json:"milestone,omitempty"`
	Approver  *User      `json:"approver,omitempty"`
}

// ExpenseCreate represents the request to create an expense
type ExpenseCreate struct {
	MilestoneID string  `json:"milestone_id,omitempty"`
	Title       string  `json:"title" validate:"required,min=3,max=255"`
	Description string  `json:"description,omitempty"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Currency    string  `json:"currency,omitempty"`
	Category    string  `json:"category,omitempty"`
	ReceiptURL  string  `json:"receipt_url,omitempty"`
	ExpenseDate string  `json:"expense_date,omitempty"`
}

// ExpenseUpdate represents the request to update an expense
type ExpenseUpdate struct {
	Title       string  `json:"title,omitempty" validate:"omitempty,min=3,max=255"`
	Description string  `json:"description,omitempty"`
	Amount      float64 `json:"amount,omitempty" validate:"omitempty,gt=0"`
	Category    string  `json:"category,omitempty"`
	ReceiptURL  string  `json:"receipt_url,omitempty"`
	ExpenseDate string  `json:"expense_date,omitempty"`
}

// VerificationRequest represents an NGO verification request
type VerificationRequest struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	NGOID      uuid.UUID  `json:"ngo_id" db:"ngo_id"`
	RequestedBy uuid.UUID `json:"requested_by" db:"requested_by"`
	Status     string     `json:"status" db:"status"`
	Documents  StringArray `json:"documents,omitempty" db:"documents"`
	AdminNotes *string    `json:"admin_notes,omitempty" db:"admin_notes"`
	ReviewedBy *uuid.UUID `json:"reviewed_by,omitempty" db:"reviewed_by"`
	ReviewedAt *time.Time `json:"reviewed_at,omitempty" db:"reviewed_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	
	// Relations
	NGO       *NGO  `json:"ngo,omitempty"`
	Requester *User `json:"requester,omitempty"`
	Reviewer  *User `json:"reviewer,omitempty"`
}

// VerificationSubmit represents the request to submit verification documents
type VerificationSubmit struct {
	Documents []string `json:"documents" validate:"required,min=1"`
}

// VerificationReview represents admin review of verification request
type VerificationReview struct {
	Status     string `json:"status" validate:"required,oneof=approved rejected"`
	AdminNotes string `json:"admin_notes,omitempty"`
}

// Dashboard stats structures
type DonorDashboard struct {
	TotalDonated     float64       `json:"total_donated"`
	DonationCount    int           `json:"donation_count"`
	CampaignsSupported int         `json:"campaigns_supported"`
	RecentDonations  []Donation    `json:"recent_donations"`
	FeaturedCampaigns []Campaign   `json:"featured_campaigns"`
}

type NGODashboard struct {
	TotalRaised      float64     `json:"total_raised"`
	CampaignCount    int         `json:"campaign_count"`
	DonorCount       int         `json:"donor_count"`
	ActiveCampaigns  []Campaign  `json:"active_campaigns"`
	PendingMilestones []Milestone `json:"pending_milestones"`
	VerificationStatus string    `json:"verification_status"`
}

type PlatformAdminDashboard struct {
	TotalUsers        int         `json:"total_users"`
	TotalNGOs         int         `json:"total_ngos"`
	TotalCampaigns    int         `json:"total_campaigns"`
	TotalDonations    float64     `json:"total_donations"`
	PendingVerifications int      `json:"pending_verifications"`
	RecentActivity    []Activity  `json:"recent_activity"`
}

// Activity represents platform activity for admin dashboard
type Activity struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	UserID      uuid.UUID `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	
	// Relations
	User *User `json:"user,omitempty"`
}

// Pagination represents pagination parameters and metadata
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Offset     int `json:"offset"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// TokenPair represents JWT access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User   User      `json:"user"`
	Tokens TokenPair `json:"tokens"`
}