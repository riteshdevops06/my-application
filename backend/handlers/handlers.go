package handlers

import (
	"database/sql"
	"net/http"

	"ngo-donation-tracker/config"
	"ngo-donation-tracker/repositories"
	"ngo-donation-tracker/services"
	"ngo-donation-tracker/utils"
)

// HealthCheck returns server health status
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, map[string]string{"status": "ok"})
}

// NotFound handles 404 errors
func NotFound(w http.ResponseWriter, r *http.Request) {
	utils.WriteError(w, http.StatusNotFound, "Route not found")
}

// Legacy auth handlers - these will be replaced by the new AuthHandlers
// Keep them for backward compatibility during migration

func Login(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	// Create repositories and services
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg)
	authHandlers := NewAuthHandlers(authService)
	
	return authHandlers.Login
}

func Register(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	// Create repositories and services
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg)
	authHandlers := NewAuthHandlers(authService)
	
	return authHandlers.Register
}

func RefreshToken(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	// Create repositories and services
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg)
	authHandlers := NewAuthHandlers(authService)
	
	return authHandlers.RefreshToken
}

func ForgotPassword(db *sql.DB) http.HandlerFunc {
	// Create repositories and services with default config
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, &config.Config{})
	authHandlers := NewAuthHandlers(authService)
	
	return authHandlers.ForgotPassword
}

func ResetPassword(db *sql.DB) http.HandlerFunc {
	// Create repositories and services with default config
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, &config.Config{})
	authHandlers := NewAuthHandlers(authService)
	
	return authHandlers.ResetPassword
}

func GetCurrentUser(db *sql.DB) http.HandlerFunc {
	// Create repositories and services with default config
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, &config.Config{})
	authHandlers := NewAuthHandlers(authService)
	
	return authHandlers.GetCurrentUser
}

func UpdateProfile(db *sql.DB) http.HandlerFunc {
	// Create repositories and services with default config
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, &config.Config{})
	authHandlers := NewAuthHandlers(authService)
	
	return authHandlers.UpdateProfile
}

func ChangePassword(db *sql.DB) http.HandlerFunc {
	// Create repositories and services with default config
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, &config.Config{})
	authHandlers := NewAuthHandlers(authService)
	
	return authHandlers.ChangePassword
}

// Campaign handlers
func ListCampaigns(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "List campaigns not implemented yet")
	}
}

func GetCampaign(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Get campaign not implemented yet")
	}
}

func GetCampaignMilestones(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Get campaign milestones not implemented yet")
	}
}

// Donor handlers
func DonorDashboard(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Donor dashboard not implemented yet")
	}
}

func ListDonorDonations(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "List donor donations not implemented yet")
	}
}

func GetDonation(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Get donation not implemented yet")
	}
}

func CreateDonation(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Create donation not implemented yet")
	}
}

// NGO handlers
func NGODashboard(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "NGO dashboard not implemented yet")
	}
}

func GetNGOProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Get NGO profile not implemented yet")
	}
}

func UpdateNGOProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Update NGO profile not implemented yet")
	}
}

func ListNGOCampaigns(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "List NGO campaigns not implemented yet")
	}
}

func CreateCampaign(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Create campaign not implemented yet")
	}
}

func GetNGOCampaign(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Get NGO campaign not implemented yet")
	}
}

func UpdateCampaign(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Update campaign not implemented yet")
	}
}

func DeleteCampaign(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Delete campaign not implemented yet")
	}
}

func PublishCampaign(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Publish campaign not implemented yet")
	}
}

// Milestone handlers
func ListMilestones(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "List milestones not implemented yet")
	}
}

func CreateMilestone(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Create milestone not implemented yet")
	}
}

func GetMilestone(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Get milestone not implemented yet")
	}
}

func UpdateMilestone(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Update milestone not implemented yet")
	}
}

func DeleteMilestone(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Delete milestone not implemented yet")
	}
}

func CompleteMilestone(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Complete milestone not implemented yet")
	}
}

// Expense handlers
func ListExpenses(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "List expenses not implemented yet")
	}
}

func CreateExpense(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Create expense not implemented yet")
	}
}

func GetExpense(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Get expense not implemented yet")
	}
}

func UpdateExpense(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Update expense not implemented yet")
	}
}

func DeleteExpense(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Delete expense not implemented yet")
	}
}

// Verification handlers
func GetVerificationStatus(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Get verification status not implemented yet")
	}
}

func SubmitVerification(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Submit verification not implemented yet")
	}
}

// Admin handlers
func AdminDashboard(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Admin dashboard not implemented yet")
	}
}

func AdminListNGOs(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Admin list NGOs not implemented yet")
	}
}

func AdminGetNGO(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Admin get NGO not implemented yet")
	}
}

func AdminVerifyNGO(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Admin verify NGO not implemented yet")
	}
}

func AdminRejectNGO(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Admin reject NGO not implemented yet")
	}
}

func AdminListCampaigns(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Admin list campaigns not implemented yet")
	}
}

func AdminGetCampaign(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Admin get campaign not implemented yet")
	}
}

func AdminSuspendCampaign(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Admin suspend campaign not implemented yet")
	}
}

func AdminListMilestones(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Admin list milestones not implemented yet")
	}
}

func AdminApproveMilestone(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Admin approve milestone not implemented yet")
	}
}

func AdminRejectMilestone(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Admin reject milestone not implemented yet")
	}
}

func AdminDonationReport(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Admin donation report not implemented yet")
	}
}

func AdminNGOReport(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Admin NGO report not implemented yet")
	}
}

func AdminCampaignReport(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Admin campaign report not implemented yet")
	}
}

func AdminGetSettings(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Admin get settings not implemented yet")
	}
}

func AdminUpdateSettings(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotImplemented, "Admin update settings not implemented yet")
	}
}