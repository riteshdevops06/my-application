package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ngo-donation-tracker/config"
	"ngo-donation-tracker/routes"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// TestEnvironment holds test dependencies
type TestEnvironment struct {
	DB     *sql.DB
	Router *mux.Router
	Config *config.Config
}

// SetupTestEnvironment initializes test environment
func SetupTestEnvironment(t *testing.T) *TestEnvironment {
	// Test configuration
	cfg := &config.Config{
		DatabaseURL: "postgres://localhost/ngo_tracker_test?sslmode=disable",
		JWTSecret:   "test-secret-key",
		Port:        "8081",
	}

	// Test database connection
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		t.Skipf("Skipping test: could not connect to test database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Skipf("Skipping test: test database not available: %v", err)
	}

	// Setup routes
	router := routes.Setup(db, cfg)

	return &TestEnvironment{
		DB:     db,
		Router: router,
		Config: cfg,
	}
}

// TeardownTestEnvironment cleans up test environment
func (env *TestEnvironment) TeardownTestEnvironment() {
	env.DB.Close()
}

// Helper function to make HTTP requests
func (env *TestEnvironment) makeRequest(method, path string, body interface{}, headers map[string]string) (*httptest.ResponseRecorder, error) {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, path, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	// Set default headers
	req.Header.Set("Content-Type", "application/json")
	
	// Set custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Record response
	rr := httptest.NewRecorder()
	env.Router.ServeHTTP(rr, req)

	return rr, nil
}

// Helper function to authenticate and get JWT token
func (env *TestEnvironment) authenticate(email, password string) (string, error) {
	loginData := map[string]string{
		"email":    email,
		"password": password,
	}

	rr, err := env.makeRequest("POST", "/api/v1/auth/login", loginData, nil)
	if err != nil {
		return "", err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		return "", err
	}

	token, _ := response["token"].(string)
	return token, nil
}

// Test health check endpoint
func TestHealthCheck(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.TeardownTestEnvironment()

	rr, err := env.makeRequest("GET", "/api/v1/health", nil, nil)
	if err != nil {
		t.Fatalf("Could not make request: %v", err)
	}

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check response body
	expected := `{"status":"ok","timestamp":`
	if !bytes.Contains(rr.Body.Bytes(), []byte(expected[:len(expected)-1])) {
		t.Errorf("Handler returned unexpected body: got %v", rr.Body.String())
	}
}

// Test user registration
func TestUserRegistration(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.TeardownTestEnvironment()

	testCases := []struct {
		name         string
		requestBody  map[string]interface{}
		expectedCode int
	}{
		{
			name: "Valid registration",
			requestBody: map[string]interface{}{
				"name":     "Test User",
				"email":    "test@example.com",
				"password": "Password123!",
				"role":     "donor",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "Missing required fields",
			requestBody: map[string]interface{}{
				"email": "test2@example.com",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid email format",
			requestBody: map[string]interface{}{
				"name":     "Test User",
				"email":    "invalid-email",
				"password": "Password123!",
				"role":     "donor",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Weak password",
			requestBody: map[string]interface{}{
				"name":     "Test User",
				"email":    "test3@example.com",
				"password": "weak",
				"role":     "donor",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr, err := env.makeRequest("POST", "/api/v1/auth/register", tc.requestBody, nil)
			if err != nil {
				t.Fatalf("Could not make request: %v", err)
			}

			if status := rr.Code; status != tc.expectedCode {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tc.expectedCode)
				t.Errorf("Response body: %s", rr.Body.String())
			}
		})
	}
}

// Test user login
func TestUserLogin(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.TeardownTestEnvironment()

	// First register a user
	registrationData := map[string]interface{}{
		"name":     "Test User",
		"email":    "login-test@example.com",
		"password": "Password123!",
		"role":     "donor",
	}

	_, err := env.makeRequest("POST", "/api/v1/auth/register", registrationData, nil)
	if err != nil {
		t.Fatalf("Could not register user: %v", err)
	}

	testCases := []struct {
		name         string
		requestBody  map[string]interface{}
		expectedCode int
	}{
		{
			name: "Valid login",
			requestBody: map[string]interface{}{
				"email":    "login-test@example.com",
				"password": "Password123!",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Invalid email",
			requestBody: map[string]interface{}{
				"email":    "nonexistent@example.com",
				"password": "Password123!",
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "Invalid password",
			requestBody: map[string]interface{}{
				"email":    "login-test@example.com",
				"password": "WrongPassword",
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr, err := env.makeRequest("POST", "/api/v1/auth/login", tc.requestBody, nil)
			if err != nil {
				t.Fatalf("Could not make request: %v", err)
			}

			if status := rr.Code; status != tc.expectedCode {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tc.expectedCode)
				t.Errorf("Response body: %s", rr.Body.String())
			}

			// Check for token in successful login
			if tc.expectedCode == http.StatusOK {
				var response map[string]interface{}
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Errorf("Could not parse response: %v", err)
				}

				if _, exists := response["token"]; !exists {
					t.Error("Token not found in successful login response")
				}
			}
		})
	}
}

// Test campaign CRUD operations
func TestCampaignCRUD(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.TeardownTestEnvironment()

	// Register and authenticate NGO admin user
	registrationData := map[string]interface{}{
		"name":     "NGO Admin",
		"email":    "ngo@example.com",
		"password": "Password123!",
		"role":     "ngo_admin",
	}

	_, err := env.makeRequest("POST", "/api/v1/auth/register", registrationData, nil)
	if err != nil {
		t.Fatalf("Could not register NGO admin: %v", err)
	}

	token, err := env.authenticate("ngo@example.com", "Password123!")
	if err != nil {
		t.Fatalf("Could not authenticate NGO admin: %v", err)
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	// Test campaign creation
	t.Run("Create campaign", func(t *testing.T) {
		campaignData := map[string]interface{}{
			"title":       "Test Campaign",
			"description": "This is a test campaign for automated testing",
			"goal_amount": 10000.0,
			"category":    "education",
		}

		rr, err := env.makeRequest("POST", "/api/v1/ngo/campaigns", campaignData, headers)
		if err != nil {
			t.Fatalf("Could not make request: %v", err)
		}

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
			t.Errorf("Response body: %s", rr.Body.String())
		}
	})

	// Test campaign listing
	t.Run("List campaigns", func(t *testing.T) {
		rr, err := env.makeRequest("GET", "/api/v1/ngo/campaigns", nil, headers)
		if err != nil {
			t.Fatalf("Could not make request: %v", err)
		}

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
			t.Errorf("Response body: %s", rr.Body.String())
		}
	})
}

// Test donation flow
func TestDonationFlow(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.TeardownTestEnvironment()

	// Register donor user
	donorData := map[string]interface{}{
		"name":     "Donor User",
		"email":    "donor@example.com",
		"password": "Password123!",
		"role":     "donor",
	}

	_, err := env.makeRequest("POST", "/api/v1/auth/register", donorData, nil)
	if err != nil {
		t.Fatalf("Could not register donor: %v", err)
	}

	donorToken, err := env.authenticate("donor@example.com", "Password123!")
	if err != nil {
		t.Fatalf("Could not authenticate donor: %v", err)
	}

	donorHeaders := map[string]string{
		"Authorization": "Bearer " + donorToken,
	}

	// Create a campaign first (using NGO admin)
	ngoData := map[string]interface{}{
		"name":     "NGO Admin",
		"email":    "ngo-donation@example.com",
		"password": "Password123!",
		"role":     "ngo_admin",
	}

	_, err = env.makeRequest("POST", "/api/v1/auth/register", ngoData, nil)
	if err != nil {
		t.Fatalf("Could not register NGO admin: %v", err)
	}

	ngoToken, err := env.authenticate("ngo-donation@example.com", "Password123!")
	if err != nil {
		t.Fatalf("Could not authenticate NGO admin: %v", err)
	}

	ngoHeaders := map[string]string{
		"Authorization": "Bearer " + ngoToken,
	}

	// Create campaign
	campaignData := map[string]interface{}{
		"title":       "Donation Test Campaign",
		"description": "Campaign for testing donations",
		"goal_amount": 5000.0,
		"category":    "health",
	}

	campaignRR, err := env.makeRequest("POST", "/api/v1/ngo/campaigns", campaignData, ngoHeaders)
	if err != nil {
		t.Fatalf("Could not create campaign: %v", err)
	}

	var campaignResponse map[string]interface{}
	if err := json.Unmarshal(campaignRR.Body.Bytes(), &campaignResponse); err != nil {
		t.Fatalf("Could not parse campaign response: %v", err)
	}

	campaignID, _ := campaignResponse["id"].(string)

	// Test donation
	t.Run("Make donation", func(t *testing.T) {
		donationData := map[string]interface{}{
			"campaign_id": campaignID,
			"amount":      100.0,
			"message":     "Test donation",
		}

		rr, err := env.makeRequest("POST", "/api/v1/donor/donations", donationData, donorHeaders)
		if err != nil {
			t.Fatalf("Could not make donation request: %v", err)
		}

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("Donation failed: got status %v want %v", status, http.StatusCreated)
			t.Errorf("Response body: %s", rr.Body.String())
		}
	})

	// Test donation history
	t.Run("Get donation history", func(t *testing.T) {
		rr, err := env.makeRequest("GET", "/api/v1/donor/donations", nil, donorHeaders)
		if err != nil {
			t.Fatalf("Could not get donation history: %v", err)
		}

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})
}

// Benchmark tests
func BenchmarkHealthCheck(b *testing.B) {
	env := SetupTestEnvironment(&testing.T{})
	defer env.TeardownTestEnvironment()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr, _ := env.makeRequest("GET", "/api/v1/health", nil, nil)
		if rr.Code != http.StatusOK {
			b.Errorf("Unexpected status code: %d", rr.Code)
		}
	}
}

func BenchmarkUserLogin(b *testing.B) {
	env := SetupTestEnvironment(&testing.T{})
	defer env.TeardownTestEnvironment()

	// Register a user first
	registrationData := map[string]interface{}{
		"name":     "Benchmark User",
		"email":    "benchmark@example.com",
		"password": "Password123!",
		"role":     "donor",
	}

	env.makeRequest("POST", "/api/v1/auth/register", registrationData, nil)

	loginData := map[string]interface{}{
		"email":    "benchmark@example.com",
		"password": "Password123!",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr, _ := env.makeRequest("POST", "/api/v1/auth/login", loginData, nil)
		if rr.Code != http.StatusOK {
			b.Errorf("Unexpected status code: %d", rr.Code)
		}
	}
}

// Load testing helper
func TestConcurrentRequests(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.TeardownTestEnvironment()

	const numRequests = 100
	const numWorkers = 10

	// Channel for work
	requests := make(chan int, numRequests)
	results := make(chan int, numRequests)

	// Worker function
	worker := func() {
		for range requests {
			rr, err := env.makeRequest("GET", "/api/v1/health", nil, nil)
			if err != nil {
				results <- 0
				continue
			}
			results <- rr.Code
		}
	}

	// Start workers
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	// Send requests
	for i := 0; i < numRequests; i++ {
		requests <- i
	}
	close(requests)

	// Collect results
	successCount := 0
	for i := 0; i < numRequests; i++ {
		if code := <-results; code == http.StatusOK {
			successCount++
		}
	}

	successRate := float64(successCount) / float64(numRequests)
	t.Logf("Success rate: %.2f%% (%d/%d)", successRate*100, successCount, numRequests)

	if successRate < 0.95 { // 95% success rate threshold
		t.Errorf("Success rate too low: %.2f%%", successRate*100)
	}
}