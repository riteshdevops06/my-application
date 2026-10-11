package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"ngo-donation-tracker/models"
	"ngo-donation-tracker/repositories"

	"github.com/google/uuid"
)

// NGOService handles NGO business logic
type NGOService struct {
	ngoRepo  repositories.NGORepositoryInterface
	userRepo repositories.UserRepositoryInterface
}

// NewNGOService creates a new NGO service
func NewNGOService(ngoRepo repositories.NGORepositoryInterface, userRepo repositories.UserRepositoryInterface) *NGOService {
	return &NGOService{
		ngoRepo:  ngoRepo,
		userRepo: userRepo,
	}
}

// CreateNGOProfile creates a new NGO profile for a user
func (s *NGOService) CreateNGOProfile(ctx context.Context, userID uuid.UUID, req *models.NGOCreate) (*models.NGO, error) {
	// Validate input
	if err := s.validateNGOInput(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Verify user exists and is NGO admin
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.Role != "ngo_admin" {
		return nil, fmt.Errorf("user must have ngo_admin role to create NGO profile")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("user account is not active")
	}

	// Create NGO profile using enhanced repository
	if enhancedRepo, ok := s.ngoRepo.(*repositories.NGORepositoryEnhanced); ok {
		return enhancedRepo.CreateNGO(ctx, userID, req)
	}

	// Fallback to interface method (for backward compatibility)
	ngo := &models.NGO{
		UserID:                  userID,
		Name:                    req.Name,
		Description:             &req.Description,
		Mission:                 &req.Mission,
		Website:                 &req.Website,
		RegistrationNumber:      &req.RegistrationNumber,
		RegistrationDocumentURL: &req.RegistrationDocumentURL,
		LogoURL:                 &req.LogoURL,
		CoverImageURL:           &req.CoverImageURL,
		Address:                 &req.Address,
		City:                    &req.City,
		Country:                 &req.Country,
		VerificationStatus:      "pending",
	}

	err = s.ngoRepo.Create(ctx, ngo)
	if err != nil {
		return nil, fmt.Errorf("failed to create NGO: %w", err)
	}

	return ngo, nil
}

// GetNGOProfile retrieves NGO profile by ID
func (s *NGOService) GetNGOProfile(ctx context.Context, id uuid.UUID) (*models.NGO, error) {
	ngo, err := s.ngoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get NGO: %w", err)
	}
	if ngo == nil {
		return nil, fmt.Errorf("NGO not found")
	}
	return ngo, nil
}

// GetNGOByUserID retrieves NGO profile by user ID
func (s *NGOService) GetNGOByUserID(ctx context.Context, userID uuid.UUID) (*models.NGO, error) {
	ngo, err := s.ngoRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get NGO: %w", err)
	}
	if ngo == nil {
		return nil, fmt.Errorf("NGO profile not found for user")
	}
	return ngo, nil
}

// UpdateNGOProfile updates NGO profile information
func (s *NGOService) UpdateNGOProfile(ctx context.Context, userID uuid.UUID, req *models.NGOCreate) (*models.NGO, error) {
	// Validate input
	if err := s.validateNGOInput(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Get existing NGO to verify ownership
	existingNGO, err := s.ngoRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing NGO: %w", err)
	}
	if existingNGO == nil {
		return nil, fmt.Errorf("NGO profile not found for user")
	}

	// Use enhanced repository if available
	if enhancedRepo, ok := s.ngoRepo.(*repositories.NGORepositoryEnhanced); ok {
		return enhancedRepo.UpdateNGO(ctx, existingNGO.ID, req)
	}

	// Fallback to interface method
	existingNGO.Name = req.Name
	existingNGO.Description = &req.Description
	existingNGO.Mission = &req.Mission
	existingNGO.Website = &req.Website
	existingNGO.RegistrationNumber = &req.RegistrationNumber
	existingNGO.RegistrationDocumentURL = &req.RegistrationDocumentURL
	existingNGO.LogoURL = &req.LogoURL
	existingNGO.CoverImageURL = &req.CoverImageURL
	existingNGO.Address = &req.Address
	existingNGO.City = &req.City
	existingNGO.Country = &req.Country

	err = s.ngoRepo.Update(ctx, existingNGO)
	if err != nil {
		return nil, fmt.Errorf("failed to update NGO: %w", err)
	}

	return existingNGO, nil
}

// ListNGOs retrieves paginated list of NGOs with optional filters
func (s *NGOService) ListNGOs(ctx context.Context, params models.Pagination, verificationStatus string) ([]*models.NGO, int64, error) {
	ngos, total, err := s.ngoRepo.List(ctx, params, verificationStatus)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list NGOs: %w", err)
	}
	return ngos, total, nil
}
// GetVerifiedNGOs returns all verified NGOs
func (s *NGOService) GetVerifiedNGOs(ctx context.Context, params models.Pagination) ([]*models.NGO, int64, error) {
	return s.ListNGOs(ctx, params, "verified")
}

// GetPendingVerifications returns NGOs pending verification
func (s *NGOService) GetPendingVerifications(ctx context.Context, params models.Pagination) ([]*models.NGO, int64, error) {
	return s.ListNGOs(ctx, params, "pending")
}

// VerifyNGO approves an NGO's verification request
func (s *NGOService) VerifyNGO(ctx context.Context, ngoID uuid.UUID, verifierID uuid.UUID, notes string) error {
	// Verify that the verifier is a platform admin
	verifier, err := s.userRepo.GetUserByID(ctx, verifierID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return fmt.Errorf("verifier not found")
		}
		return fmt.Errorf("failed to get verifier: %w", err)
	}

	if verifier.Role != "platform_admin" {
		return fmt.Errorf("only platform admins can verify NGOs")
	}

	// Verify NGO exists
	ngo, err := s.ngoRepo.GetByID(ctx, ngoID)
	if err != nil {
		return fmt.Errorf("failed to get NGO: %w", err)
	}
	if ngo == nil {
		return fmt.Errorf("NGO not found")
	}

	// Update verification status
	notesPtr := &notes
	if notes == "" {
		notesPtr = nil
	}

	err = s.ngoRepo.UpdateVerificationStatus(ctx, ngoID, "verified", notesPtr, &verifierID)
	if err != nil {
		return fmt.Errorf("failed to verify NGO: %w", err)
	}

	return nil
}

// RejectNGO rejects an NGO's verification request
func (s *NGOService) RejectNGO(ctx context.Context, ngoID uuid.UUID, verifierID uuid.UUID, reason string) error {
	// Verify that the verifier is a platform admin
	verifier, err := s.userRepo.GetUserByID(ctx, verifierID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return fmt.Errorf("verifier not found")
		}
		return fmt.Errorf("failed to get verifier: %w", err)
	}

	if verifier.Role != "platform_admin" {
		return fmt.Errorf("only platform admins can reject NGOs")
	}

	// Verify NGO exists
	ngo, err := s.ngoRepo.GetByID(ctx, ngoID)
	if err != nil {
		return fmt.Errorf("failed to get NGO: %w", err)
	}
	if ngo == nil {
		return fmt.Errorf("NGO not found")
	}

	// Require rejection reason
	if strings.TrimSpace(reason) == "" {
		return fmt.Errorf("rejection reason is required")
	}

	// Update verification status
	err = s.ngoRepo.UpdateVerificationStatus(ctx, ngoID, "rejected", &reason, &verifierID)
	if err != nil {
		return fmt.Errorf("failed to reject NGO: %w", err)
	}

	return nil
}

// ResubmitForVerification resets NGO status to pending for re-verification
func (s *NGOService) ResubmitForVerification(ctx context.Context, userID uuid.UUID) error {
	// Get NGO by user ID
	ngo, err := s.ngoRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get NGO: %w", err)
	}
	if ngo == nil {
		return fmt.Errorf("NGO profile not found for user")
	}

	// Only allow resubmission if currently rejected
	if ngo.VerificationStatus != "rejected" {
		return fmt.Errorf("can only resubmit rejected NGOs for verification")
	}

	// Reset verification status
	err = s.ngoRepo.UpdateVerificationStatus(ctx, ngo.ID, "pending", nil, nil)
	if err != nil {
		return fmt.Errorf("failed to resubmit for verification: %w", err)
	}

	return nil
}

// GetNGOStats returns NGO statistics
func (s *NGOService) GetNGOStats(ctx context.Context) (*repositories.NGOStats, error) {
	if enhancedRepo, ok := s.ngoRepo.(*repositories.NGORepositoryEnhanced); ok {
		return enhancedRepo.GetNGOStats(ctx)
	}

	// Fallback implementation using basic interface
	_, verifiedTotal, err := s.ngoRepo.List(ctx, models.Pagination{Limit: 1}, "verified")
	if err != nil {
		return nil, fmt.Errorf("failed to get verified count: %w", err)
	}

	_, pendingTotal, err := s.ngoRepo.List(ctx, models.Pagination{Limit: 1}, "pending")
	if err != nil {
		return nil, fmt.Errorf("failed to get pending count: %w", err)
	}

	_, rejectedTotal, err := s.ngoRepo.List(ctx, models.Pagination{Limit: 1}, "rejected")
	if err != nil {
		return nil, fmt.Errorf("failed to get rejected count: %w", err)
	}

	_, total, err := s.ngoRepo.List(ctx, models.Pagination{Limit: 1}, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	return &repositories.NGOStats{
		TotalNGOs:     int(total),
		VerifiedNGOs:  int(verifiedTotal),
		PendingNGOs:   int(pendingTotal),
		RejectedNGOs:  int(rejectedTotal),
		NewNGOs30Days: 0, // Cannot calculate without enhanced repository
	}, nil
}

// SearchNGOs searches NGOs by name
func (s *NGOService) SearchNGOs(ctx context.Context, searchTerm string, params models.Pagination) ([]*models.NGO, int64, error) {
	if enhancedRepo, ok := s.ngoRepo.(*repositories.NGORepositoryEnhanced); ok {
		return enhancedRepo.SearchNGOsByName(ctx, searchTerm, params)
	}

	// Fallback - return all NGOs (basic search not available)
	return s.ngoRepo.List(ctx, params, "")
}

// Validation methods

func (s *NGOService) validateNGOInput(req *models.NGOCreate) error {
	if req.Name == "" {
		return errors.New("NGO name is required")
	}
	if len(strings.TrimSpace(req.Name)) < 2 {
		return errors.New("NGO name must be at least 2 characters long")
	}
	if len(req.Name) > 255 {
		return errors.New("NGO name must not exceed 255 characters")
	}

	// Website validation
	if req.Website != "" && !isValidURL(req.Website) {
		return errors.New("invalid website URL format")
	}

	return nil
}

// Helper function for basic URL validation
func isValidURL(url string) bool {
	return strings.HasPrefix(strings.ToLower(url), "http://") || 
		   strings.HasPrefix(strings.ToLower(url), "https://")
}