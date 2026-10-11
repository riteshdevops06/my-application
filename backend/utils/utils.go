package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"ngo-donation-tracker/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// GenerateToken creates a new JWT token
func GenerateToken(userID uuid.UUID, email, role string, ngoID *uuid.UUID, secret string, expiryHours int) (string, error) {
	claims := models.Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		NGOID:  ngoID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a password with its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// WriteError writes an error response
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

// WriteSuccess writes a success response
func WriteSuccess(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

// WriteCreated writes a created response
func WriteCreated(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

// ParseUUID parses a string to UUID
func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// GetPaginationParams extracts pagination from query params
func GetPaginationParams(r *http.Request) models.Pagination {
	page := 1
	limit := 20

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := parseInt(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := parseInt(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	return models.Pagination{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}

func parseInt(s string) (int, error) {
	var result int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, ErrInvalidInt
		}
		result = result*10 + int(c-'0')
	}
	return result, nil
}

var ErrInvalidInt = &ParseError{"invalid integer"}

type ParseError struct{ msg string }

func (e *ParseError) Error() string { return e.msg }
// ParseToken parses and validates a JWT token
func ParseToken(tokenString, secret string) (*models.Claims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(*models.Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Check if token is expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token is expired")
	}

	return claims, nil
}