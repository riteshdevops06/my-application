package services

import (
	"context"
	"fmt"
	"time"

	"ngo-donation-tracker/models"
	"ngo-donation-tracker/repositories"

	"github.com/google/uuid"
)

// DonationService handles business logic for donations
type DonationService struct {
	donationRepo repositories.DonationRepository
	campaignRepo repositories.CampaignRepository
	userRepo     repositories.UserRepository
	ngoRepo      repositories.NGORepositoryInterface
}

// NewDonationService creates a new instance of DonationService
func NewDonationService(
	donationRepo repositories.DonationRepository,
	campaignRepo repositories.CampaignRepository,
	userRepo repositories.UserRepository,
	ngoRepo repositories.NGORepositoryInterface,
) *DonationService {
	return &DonationService{
		donationRepo: donationRepo,
		campaignRepo: campaignRepo,
		userRepo:     userRepo,
		ngoRepo:      ngoRepo,
	}
}

// CreateDonationRequest represents request to create a donation
type CreateDonationRequest struct {
	CampaignID    uuid.UUID `json:"campaign_id" validate:"required"`
	Amount        float64   `json:"amount" validate:"required,gt=0"`
	Currency      string    `json:"currency" validate:"required,len=3"`
	PaymentMethod string    `json:"payment_method" validate:"required,oneof=card bank_transfer mobile_money"`
	Anonymous     bool      `json:"anonymous"`
	Message       string    `json:"message,omitempty"`
	ReceiptEmail  string    `json:"receipt_email,omitempty" validate:"omitempty,email"`
}

// ProcessDonationRequest represents request to process a pending donation
type ProcessDonationRequest struct {
	TransactionID string `json:"transaction_id" validate:"required"`
	Status        string `json:"status" validate:"required,oneof=completed failed"`
}

// RefundDonationRequest represents request to refund a donation
type RefundDonationRequest struct {
	Reason string `json:"reason" validate:"required"`
}

// DonationResponse represents donation response with additional details
type DonationResponse struct {
	*models.Donation
	CampaignTitle string `json:"campaign_title,omitempty"`
	DonorName     string `json:"donor_name,omitempty"`
	NGOName       string `json:"ngo_name,omitempty"`
}

// CreateDonation creates a new donation
func (s *DonationService) CreateDonation(donorID uuid.UUID, req *CreateDonationRequest) (*models.Donation, error) {
	// Validate campaign exists and is active
	campaign, err := s.campaignRepo.GetByID(context.Background(), req.CampaignID)
	if err != nil {
		return nil, fmt.Errorf("invalid campaign: %v", err)
	}

	if campaign.Status != "active" {
		return nil, fmt.Errorf("campaign is not active")
	}

	// Check campaign deadline
	if campaign.EndDate != nil && time.Now().After(*campaign.EndDate) {
		return nil, fmt.Errorf("campaign has ended")
	}

	// Validate donor exists
	_, err = s.userRepo.GetUserByID(context.Background(), donorID)
	if err != nil {
		return nil, fmt.Errorf("invalid donor: %v", err)
	}

	// Validate minimum donation amount
	if req.Amount < 1.0 {
		return nil, fmt.Errorf("minimum donation amount is $1.00")
	}

	// Check if campaign has reached its goal (optional check)
	if campaign.GoalAmount > 0 {
		stats, err := s.donationRepo.GetDonationStats(req.CampaignID)
		if err == nil && stats.CompletedAmount >= campaign.GoalAmount {
			return nil, fmt.Errorf("campaign has reached its funding goal")
		}
	}

	// Create donation with pending status
	paymentMethod := &req.PaymentMethod
	message := &req.Message
	receiptEmail := &req.ReceiptEmail
	if req.PaymentMethod == "" {
		paymentMethod = nil
	}
	if req.Message == "" {
		message = nil
	}
	if req.ReceiptEmail == "" {
		receiptEmail = nil
	}

	donation := &models.Donation{
		ID:            uuid.New(),
		DonorID:       donorID,
		CampaignID:    req.CampaignID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Status:        "pending",
		PaymentMethod: paymentMethod,
		Anonymous:     req.Anonymous,
		Message:       message,
		ReceiptEmail:  receiptEmail,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = s.donationRepo.CreateDonation(donation)
	if err != nil {
		return nil, fmt.Errorf("failed to create donation: %v", err)
	}

	return donation, nil
}

// ProcessDonation processes a pending donation (updates status after payment)
func (s *DonationService) ProcessDonation(donationID uuid.UUID, req *ProcessDonationRequest) (*models.Donation, error) {
	// Get existing donation
	donation, err := s.donationRepo.GetDonationByID(donationID)
	if err != nil {
		return nil, fmt.Errorf("donation not found: %v", err)
	}

	if donation.Status != "pending" {
		return nil, fmt.Errorf("donation is not pending")
	}

	// Update donation status and transaction ID
	err = s.donationRepo.UpdateStatus(donationID, req.Status, req.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to update donation status: %v", err)
	}

	// Get updated donation
	updatedDonation, err := s.donationRepo.GetDonationByID(donationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated donation: %v", err)
	}

	// TODO: Send notification email if donation completed
	// TODO: Update campaign statistics cache if needed

	return updatedDonation, nil
}

// RefundDonation processes a donation refund
func (s *DonationService) RefundDonation(donationID uuid.UUID, req *RefundDonationRequest, requesterID uuid.UUID) (*models.Donation, error) {
	// Get existing donation
	donation, err := s.donationRepo.GetDonationByID(donationID)
	if err != nil {
		return nil, fmt.Errorf("donation not found: %v", err)
	}

	if donation.Status != "completed" {
		return nil, fmt.Errorf("only completed donations can be refunded")
	}

	// Check authorization - only the donor or admin can request refund
	requester, err := s.userRepo.GetUserByID(context.Background(), requesterID)
	if err != nil {
		return nil, fmt.Errorf("invalid requester: %v", err)
	}

	if requester.Role != "platform_admin" && donation.DonorID != requesterID {
		return nil, fmt.Errorf("unauthorized to refund this donation")
	}

	// Update donation status to refunded
	txID := ""
	if donation.TransactionID != nil {
		txID = *donation.TransactionID
	}
	err = s.donationRepo.UpdateStatus(donationID, "refunded", txID)
	if err != nil {
		return nil, fmt.Errorf("failed to refund donation: %v", err)
	}

	// Get updated donation
	updatedDonation, err := s.donationRepo.GetDonationByID(donationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated donation: %v", err)
	}

	// TODO: Process actual refund with payment gateway
	// TODO: Send refund notification email
	// TODO: Update campaign statistics

	return updatedDonation, nil
}

// GetDonation retrieves a donation with additional details
func (s *DonationService) GetDonation(donationID uuid.UUID, requesterID uuid.UUID) (*DonationResponse, error) {
	donation, err := s.donationRepo.GetDonationByID(donationID)
	if err != nil {
		return nil, fmt.Errorf("donation not found: %v", err)
	}

	// Check authorization
	requester, err := s.userRepo.GetUserByID(context.Background(), requesterID)
	if err != nil {
		return nil, fmt.Errorf("invalid requester: %v", err)
	}

	// Only donor, campaign NGO admin, or platform admin can view donation details
	if requester.Role != "platform_admin" && donation.DonorID != requesterID {
		// Check if requester is NGO admin of the campaign
		campaign, err := s.campaignRepo.GetByID(context.Background(), donation.CampaignID)
		if err != nil {
			return nil, fmt.Errorf("invalid campaign: %v", err)
		}

		if requester.Role != "ngo_admin" {
			return nil, fmt.Errorf("unauthorized to view this donation")
		}

		// Get NGO by user ID to check if user is admin of the NGO that owns this campaign
		ngo, err := s.ngoRepo.GetByUserID(context.Background(), requesterID)
		if err != nil || ngo.ID != campaign.NGOID {
			return nil, fmt.Errorf("unauthorized to view this donation")
		}
	}

	response := &DonationResponse{
		Donation: donation,
	}

	// Add campaign details
	campaign, err := s.campaignRepo.GetByID(context.Background(), donation.CampaignID)
	if err == nil {
		response.CampaignTitle = campaign.Title
		
		// Add NGO name
		if campaign.NGO != nil {
			response.NGOName = campaign.NGO.Name
		}
	}

	// Add donor name (only if not anonymous or requester has permission)
	if !donation.Anonymous || requester.Role == "platform_admin" {
		donor, err := s.userRepo.GetUserByID(context.Background(), donation.DonorID)
		if err == nil {
			response.DonorName = donor.Name
		}
	}

	return response, nil
}

// GetDonationsByDonor retrieves donations for a specific donor
func (s *DonationService) GetDonationsByDonor(donorID uuid.UUID, pagination repositories.PaginationParams) ([]*models.Donation, *repositories.PaginationResult, error) {
	return s.donationRepo.GetByDonorID(donorID, pagination)
}

// GetDonationsByCampaign retrieves donations for a specific campaign
func (s *DonationService) GetDonationsByCampaign(campaignID uuid.UUID, requesterID uuid.UUID, pagination repositories.PaginationParams) ([]*DonationResponse, *repositories.PaginationResult, error) {
	// Check authorization
	requester, err := s.userRepo.GetUserByID(context.Background(), requesterID)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid requester: %v", err)
	}

	// Check if requester can view campaign donations
	if requester.Role != "platform_admin" {
		campaign, err := s.campaignRepo.GetByID(context.Background(), campaignID)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid campaign: %v", err)
		}

		if requester.Role != "ngo_admin" {
			return nil, nil, fmt.Errorf("unauthorized to view campaign donations")
		}

		// Check if user is admin of the NGO that owns this campaign
		ngo, err := s.ngoRepo.GetByUserID(context.Background(), requesterID)
		if err != nil || ngo.ID != campaign.NGOID {
			return nil, nil, fmt.Errorf("unauthorized to view campaign donations")
		}
	}

	donations, result, err := s.donationRepo.GetByCampaignID(campaignID, pagination)
	if err != nil {
		return nil, nil, err
	}

	var responses []*DonationResponse
	for _, donation := range donations {
		response := &DonationResponse{
			Donation: donation,
		}

		// Add donor name if not anonymous
		if !donation.Anonymous {
			donor, err := s.userRepo.GetUserByID(context.Background(), donation.DonorID)
			if err == nil {
				response.DonorName = donor.Name
			}
		}

		responses = append(responses, response)
	}

	return responses, result, nil
}

// GetDonationHistory retrieves donation history for a donor
func (s *DonationService) GetDonationHistory(donorID uuid.UUID, pagination repositories.PaginationParams) ([]*models.DonationWithCampaign, error) {
	return s.donationRepo.GetDonationHistory(donorID, pagination)
}

// GetDonationStats retrieves aggregated donation statistics for a campaign
func (s *DonationService) GetDonationStats(campaignID uuid.UUID) (*repositories.DonationStats, error) {
	return s.donationRepo.GetDonationStats(campaignID)
}

// ListDonations retrieves donations based on filters (admin only)
func (s *DonationService) ListDonations(filters repositories.DonationFiltersExtended, pagination repositories.PaginationParams, requesterID uuid.UUID) ([]*DonationResponse, *repositories.PaginationResult, error) {
	// Check authorization - only platform admin can list all donations
	requester, err := s.userRepo.GetUserByID(context.Background(), requesterID)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid requester: %v", err)
	}

	if requester.Role != "platform_admin" {
		return nil, nil, fmt.Errorf("unauthorized to list all donations")
	}

	donations, result, err := s.donationRepo.ListEnhanced(filters, pagination)
	if err != nil {
		return nil, nil, err
	}

	var responses []*DonationResponse
	for _, donation := range donations {
		response := &DonationResponse{
			Donation: donation,
		}

		// Add campaign details
		campaign, err := s.campaignRepo.GetByID(context.Background(), donation.CampaignID)
		if err == nil {
			response.CampaignTitle = campaign.Title
			if campaign.NGO != nil {
				response.NGOName = campaign.NGO.Name
			}
		}

		// Add donor name if not anonymous
		if !donation.Anonymous {
			donor, err := s.userRepo.GetUserByID(context.Background(), donation.DonorID)
			if err == nil {
				response.DonorName = donor.Name
			}
		}

		responses = append(responses, response)
	}

	return responses, result, nil
}

// GetDonationsByDateRange retrieves donations within a date range (admin only)
func (s *DonationService) GetDonationsByDateRange(startDate, endDate time.Time, requesterID uuid.UUID, pagination repositories.PaginationParams) ([]*DonationResponse, *repositories.PaginationResult, error) {
	// Check authorization
	requester, err := s.userRepo.GetUserByID(context.Background(), requesterID)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid requester: %v", err)
	}

	if requester.Role != "platform_admin" {
		return nil, nil, fmt.Errorf("unauthorized to view donations by date range")
	}

	donations, result, err := s.donationRepo.GetDonationsByDateRange(startDate, endDate, pagination)
	if err != nil {
		return nil, nil, err
	}

	var responses []*DonationResponse
	for _, donation := range donations {
		response := &DonationResponse{
			Donation: donation,
		}

		// Add campaign details
		campaign, err := s.campaignRepo.GetByID(context.Background(), donation.CampaignID)
		if err == nil {
			response.CampaignTitle = campaign.Title
			if campaign.NGO != nil {
				response.NGOName = campaign.NGO.Name
			}
		}

		// Add donor name if not anonymous
		if !donation.Anonymous {
			donor, err := s.userRepo.GetUserByID(context.Background(), donation.DonorID)
			if err == nil {
				response.DonorName = donor.Name
			}
		}

		responses = append(responses, response)
	}

	return responses, result, nil
}

// ValidateDonationAccess checks if user can access donation information
func (s *DonationService) ValidateDonationAccess(donationID uuid.UUID, requesterID uuid.UUID) error {
	donation, err := s.donationRepo.GetDonationByID(donationID)
	if err != nil {
		return fmt.Errorf("donation not found: %v", err)
	}

	requester, err := s.userRepo.GetUserByID(context.Background(), requesterID)
	if err != nil {
		return fmt.Errorf("invalid requester: %v", err)
	}

	// Platform admin has access to all donations
	if requester.Role == "platform_admin" {
		return nil
	}

	// Donor has access to their own donations
	if donation.DonorID == requesterID {
		return nil
	}

	// NGO admin has access to donations for their campaigns
	if requester.Role == "ngo_admin" {
		ngo, err := s.ngoRepo.GetByUserID(context.Background(), requesterID)
		if err == nil {
			campaign, err := s.campaignRepo.GetByID(context.Background(), donation.CampaignID)
			if err == nil && ngo.ID == campaign.NGOID {
				return nil
			}
		}
	}

	return fmt.Errorf("unauthorized to access this donation")
}