package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Connect to database
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/ngo_tracker?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Seeding database...")

	// Hash password for seed users
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Seed platform admin
	_, err = db.Exec(`
		INSERT INTO users (email, password_hash, name, role, email_verified)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (email) DO NOTHING
	`, "admin@ngotracker.org", string(hashedPassword), "Platform Admin", "platform_admin", true)
	if err != nil {
		log.Printf("Warning: Failed to seed platform admin: %v", err)
	}

	// Seed test donor
	_, err = db.Exec(`
		INSERT INTO users (email, password_hash, name, role, email_verified)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (email) DO NOTHING
	`, "donor@example.com", string(hashedPassword), "Test Donor", "donor", true)
	if err != nil {
		log.Printf("Warning: Failed to seed donor: %v", err)
	}

	// Seed test NGO admin
	_, err = db.Exec(`
		INSERT INTO users (email, password_hash, name, role, email_verified)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (email) DO NOTHING
	`, "ngo@example.com", string(hashedPassword), "Test NGO Admin", "ngo_admin", true)
	if err != nil {
		log.Printf("Warning: Failed to seed NGO admin: %v", err)
	}

	// Get NGO admin user ID for NGO creation
	var ngoAdminID string
	err = db.QueryRow(`SELECT id FROM users WHERE email = $1`, "ngo@example.com").Scan(&ngoAdminID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Warning: Failed to get NGO admin ID: %v", err)
	}

	if ngoAdminID != "" {
		// Seed test NGO
		_, err = db.Exec(`
			INSERT INTO ngos (user_id, name, description, mission, website, registration_number, city, country, verification_status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (user_id) DO NOTHING
		`, ngoAdminID, "Hope Foundation", "Providing education to underprivileged children", "Every child deserves quality education", "https://hopefoundation.org", "REG-12345", "New York", "USA", "verified")
		if err != nil {
			log.Printf("Warning: Failed to seed NGO: %v", err)
		}

		// Get NGO ID
		var ngoID string
		err = db.QueryRow(`SELECT id FROM ngos WHERE user_id = $1`, ngoAdminID).Scan(&ngoID)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Warning: Failed to get NGO ID: %v", err)
		}

		if ngoID != "" {
			// Seed test campaign
			_, err = db.Exec(`
				INSERT INTO campaigns (ngo_id, title, description, short_description, goal_amount, currency, status, category, start_date, end_date, location)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
				ON CONFLICT DO NOTHING
			`, ngoID, "Build a School in Rural Kenya", "Help us build a primary school for 200 children in rural Kenya", "Building a school for 200 children in Kenya", 50000.00, "USD", "active", "education", "2024-01-01", "2024-12-31", "Nairobi, Kenya")
			if err != nil {
				log.Printf("Warning: Failed to seed campaign: %v", err)
			}
		}
	}

	fmt.Println("Database seeding completed!")
}