package routes

import (
	"database/sql"
	"net/http"

	"ngo-donation-tracker/config"
	"ngo-donation-tracker/handlers"
	"ngo-donation-tracker/middleware"

	"github.com/gorilla/mux"
)

func Setup(db *sql.DB, cfg *config.Config) *mux.Router {
	router := mux.NewRouter()

	// Middleware
	router.Use(middleware.Logging)
	router.Use(middleware.CORS)
	router.Use(middleware.Recovery)

	// API version prefix
	api := router.PathPrefix("/api/v1").Subrouter()

	// Public routes (no auth required)
	public := api.PathPrefix("").Subrouter()
	public.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	// Auth routes
	auth := public.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", handlers.Login(db, cfg)).Methods("POST")
	auth.HandleFunc("/register", handlers.Register(db, cfg)).Methods("POST")
	auth.HandleFunc("/refresh", handlers.RefreshToken(db, cfg)).Methods("POST")
	auth.HandleFunc("/forgot-password", handlers.ForgotPassword(db)).Methods("POST")
	auth.HandleFunc("/reset-password", handlers.ResetPassword(db)).Methods("POST")

	// Public campaign browsing
	campaigns := public.PathPrefix("/campaigns").Subrouter()
	campaigns.HandleFunc("", handlers.ListCampaigns(db)).Methods("GET")
	campaigns.HandleFunc("/{id}", handlers.GetCampaign(db)).Methods("GET")
	campaigns.HandleFunc("/{id}/milestones", handlers.GetCampaignMilestones(db)).Methods("GET")

	// Protected routes (auth required)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.Auth(db, cfg))

	// User profile
	protected.HandleFunc("/auth/me", handlers.GetCurrentUser(db)).Methods("GET")
	protected.HandleFunc("/auth/me", handlers.UpdateProfile(db)).Methods("PUT")
	protected.HandleFunc("/auth/change-password", handlers.ChangePassword(db)).Methods("POST")

	// Donor routes
	donor := protected.PathPrefix("/donor").Subrouter()
	donor.Use(middleware.RequireRole("donor"))
	donor.HandleFunc("/dashboard", handlers.DonorDashboard(db)).Methods("GET")
	donor.HandleFunc("/donations", handlers.ListDonorDonations(db)).Methods("GET")
	donor.HandleFunc("/donations/{id}", handlers.GetDonation(db)).Methods("GET")
	donor.HandleFunc("/donations", handlers.CreateDonation(db, cfg)).Methods("POST")

	// NGO Admin routes
	ngo := protected.PathPrefix("/ngo").Subrouter()
	ngo.Use(middleware.RequireRole("ngo_admin"))
	ngo.HandleFunc("/dashboard", handlers.NGODashboard(db)).Methods("GET")
	ngo.HandleFunc("/profile", handlers.GetNGOProfile(db)).Methods("GET")
	ngo.HandleFunc("/profile", handlers.UpdateNGOProfile(db)).Methods("PUT")

	// NGO Campaigns
	ngoCampaigns := ngo.PathPrefix("/campaigns").Subrouter()
	ngoCampaigns.HandleFunc("", handlers.ListNGOCampaigns(db)).Methods("GET")
	ngoCampaigns.HandleFunc("", handlers.CreateCampaign(db)).Methods("POST")
	ngoCampaigns.HandleFunc("/{id}", handlers.GetNGOCampaign(db)).Methods("GET")
	ngoCampaigns.HandleFunc("/{id}", handlers.UpdateCampaign(db)).Methods("PUT")
	ngoCampaigns.HandleFunc("/{id}", handlers.DeleteCampaign(db)).Methods("DELETE")
	ngoCampaigns.HandleFunc("/{id}/publish", handlers.PublishCampaign(db)).Methods("POST")

	// NGO Milestones
	ngoMilestones := ngo.PathPrefix("/campaigns/{campaign_id}/milestones").Subrouter()
	ngoMilestones.HandleFunc("", handlers.ListMilestones(db)).Methods("GET")
	ngoMilestones.HandleFunc("", handlers.CreateMilestone(db)).Methods("POST")
	ngoMilestones.HandleFunc("/{id}", handlers.GetMilestone(db)).Methods("GET")
	ngoMilestones.HandleFunc("/{id}", handlers.UpdateMilestone(db)).Methods("PUT")
	ngoMilestones.HandleFunc("/{id}", handlers.DeleteMilestone(db)).Methods("DELETE")
	ngoMilestones.HandleFunc("/{id}/complete", handlers.CompleteMilestone(db)).Methods("POST")

	// NGO Expenses
	ngoExpenses := ngo.PathPrefix("/campaigns/{campaign_id}/expenses").Subrouter()
	ngoExpenses.HandleFunc("", handlers.ListExpenses(db)).Methods("GET")
	ngoExpenses.HandleFunc("", handlers.CreateExpense(db)).Methods("POST")
	ngoExpenses.HandleFunc("/{id}", handlers.GetExpense(db)).Methods("GET")
	ngoExpenses.HandleFunc("/{id}", handlers.UpdateExpense(db)).Methods("PUT")
	ngoExpenses.HandleFunc("/{id}", handlers.DeleteExpense(db)).Methods("DELETE")

	// NGO Verification
	ngo.HandleFunc("/verification", handlers.GetVerificationStatus(db)).Methods("GET")
	ngo.HandleFunc("/verification", handlers.SubmitVerification(db)).Methods("POST")

	// Platform Admin routes
	admin := protected.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.RequireRole("platform_admin"))
	admin.HandleFunc("/dashboard", handlers.AdminDashboard(db)).Methods("GET")

	// Admin NGO management
	adminNGOs := admin.PathPrefix("/ngos").Subrouter()
	adminNGOs.HandleFunc("", handlers.AdminListNGOs(db)).Methods("GET")
	adminNGOs.HandleFunc("/{id}", handlers.AdminGetNGO(db)).Methods("GET")
	adminNGOs.HandleFunc("/{id}/verify", handlers.AdminVerifyNGO(db)).Methods("POST")
	adminNGOs.HandleFunc("/{id}/reject", handlers.AdminRejectNGO(db)).Methods("POST")

	// Admin Campaign oversight
	adminCampaigns := admin.PathPrefix("/campaigns").Subrouter()
	adminCampaigns.HandleFunc("", handlers.AdminListCampaigns(db)).Methods("GET")
	adminCampaigns.HandleFunc("/{id}", handlers.AdminGetCampaign(db)).Methods("GET")
	adminCampaigns.HandleFunc("/{id}/suspend", handlers.AdminSuspendCampaign(db)).Methods("POST")

	// Admin Milestone oversight
	adminMilestones := admin.PathPrefix("/milestones").Subrouter()
	adminMilestones.HandleFunc("", handlers.AdminListMilestones(db)).Methods("GET")
	adminMilestones.HandleFunc("/{id}/approve", handlers.AdminApproveMilestone(db)).Methods("POST")
	adminMilestones.HandleFunc("/{id}/reject", handlers.AdminRejectMilestone(db)).Methods("POST")

	// Admin Reports
	admin.HandleFunc("/reports/donations", handlers.AdminDonationReport(db)).Methods("GET")
	admin.HandleFunc("/reports/ngos", handlers.AdminNGOReport(db)).Methods("GET")
	admin.HandleFunc("/reports/campaigns", handlers.AdminCampaignReport(db)).Methods("GET")

	// Admin Settings
	admin.HandleFunc("/settings", handlers.AdminGetSettings(db)).Methods("GET")
	admin.HandleFunc("/settings", handlers.AdminUpdateSettings(db)).Methods("PUT")

	// File upload routes (protected)
	fileHandler := handlers.NewFileHandler()
	
	// File upload endpoints
	files := protected.PathPrefix("/files").Subrouter()
	files.HandleFunc("/upload", fileHandler.UploadFile).Methods("POST")
	files.HandleFunc("/upload/multiple", fileHandler.UploadMultipleFiles).Methods("POST")
	files.HandleFunc("/validate-image", fileHandler.ValidateImageFile).Methods("POST")
	files.HandleFunc("/progress/{upload_id}", fileHandler.GetUploadProgress).Methods("GET")
	
	// File serving and management
	files.HandleFunc("/{type}/{entity_type}/{filename}", fileHandler.ServeFile).Methods("GET")
	files.HandleFunc("/{type}/{entity_type}/{filename}", fileHandler.DeleteFile).Methods("DELETE")
	files.HandleFunc("/{type}/{entity_type}/{filename}/info", fileHandler.GetFileInfo).Methods("GET")

	// 404 handler
	router.NotFoundHandler = http.HandlerFunc(handlers.NotFound)

	return router
}