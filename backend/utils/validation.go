package utils

import (
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

// ValidationRule represents a validation rule
type ValidationRule struct {
	Field   string
	Value   interface{}
	Rules   []string
	Message string
}

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Value   string `json:"value"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(messages, "; ")
}

// Validator provides validation functionality
type Validator struct {
	errors ValidationErrors
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		errors: make(ValidationErrors, 0),
	}
}

// ValidateField validates a single field with rules
func (v *Validator) ValidateField(field string, value interface{}, rules ...string) *Validator {
	for _, rule := range rules {
		if err := v.validateRule(field, value, rule); err != nil {
			v.errors = append(v.errors, *err)
		}
	}
	return v
}

// ValidateEmail validates an email address
func (v *Validator) ValidateEmail(field string, email string) *Validator {
	if email == "" {
		return v
	}
	
	if _, err := mail.ParseAddress(email); err != nil {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   email,
			Rule:    "email",
			Message: "Must be a valid email address",
		})
	}
	return v
}

// ValidatePassword validates password strength
func (v *Validator) ValidatePassword(field string, password string) *Validator {
	if len(password) < 8 {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   "***",
			Rule:    "password_length",
			Message: "Password must be at least 8 characters long",
		})
		return v
	}
	
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	
	if !hasUpper {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   "***",
			Rule:    "password_uppercase",
			Message: "Password must contain at least one uppercase letter",
		})
	}
	
	if !hasLower {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   "***",
			Rule:    "password_lowercase",
			Message: "Password must contain at least one lowercase letter",
		})
	}
	
	if !hasNumber {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   "***",
			Rule:    "password_number",
			Message: "Password must contain at least one number",
		})
	}
	
	if !hasSpecial {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   "***",
			Rule:    "password_special",
			Message: "Password must contain at least one special character",
		})
	}
	
	return v
}

// ValidateUUID validates a UUID string
func (v *Validator) ValidateUUID(field string, id string) *Validator {
	if id == "" {
		return v
	}
	
	if _, err := uuid.Parse(id); err != nil {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   id,
			Rule:    "uuid",
			Message: "Must be a valid UUID",
		})
	}
	return v
}

// ValidateURL validates a URL
func (v *Validator) ValidateURL(field string, url string) *Validator {
	if url == "" {
		return v
	}
	
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	if !urlRegex.MatchString(url) {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   url,
			Rule:    "url",
			Message: "Must be a valid URL",
		})
	}
	return v
}

// ValidatePhone validates a phone number
func (v *Validator) ValidatePhone(field string, phone string) *Validator {
	if phone == "" {
		return v
	}
	
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	if !phoneRegex.MatchString(phone) {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   phone,
			Rule:    "phone",
			Message: "Must be a valid phone number",
		})
	}
	return v
}

// ValidateDate validates a date string
func (v *Validator) ValidateDate(field string, dateStr string) *Validator {
	if dateStr == "" {
		return v
	}
	
	layouts := []string{
		"2006-01-02",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05-07:00",
	}
	
	var validDate bool
	for _, layout := range layouts {
		if _, err := time.Parse(layout, dateStr); err == nil {
			validDate = true
			break
		}
	}
	
	if !validDate {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   dateStr,
			Rule:    "date",
			Message: "Must be a valid date",
		})
	}
	return v
}

// ValidateAmount validates monetary amounts
func (v *Validator) ValidateAmount(field string, amount float64, min, max float64) *Validator {
	if amount < min {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   fmt.Sprintf("%.2f", amount),
			Rule:    "min_amount",
			Message: fmt.Sprintf("Amount must be at least $%.2f", min),
		})
	}
	
	if max > 0 && amount > max {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   fmt.Sprintf("%.2f", amount),
			Rule:    "max_amount",
			Message: fmt.Sprintf("Amount must not exceed $%.2f", max),
		})
	}
	
	return v
}

// ValidateEnum validates that a value is in a list of allowed values
func (v *Validator) ValidateEnum(field string, value string, allowedValues []string) *Validator {
	if value == "" {
		return v
	}
	
	for _, allowed := range allowedValues {
		if value == allowed {
			return v
		}
	}
	
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Value:   value,
		Rule:    "enum",
		Message: fmt.Sprintf("Must be one of: %s", strings.Join(allowedValues, ", ")),
	})
	return v
}

// ValidateArrayLength validates array length
func (v *Validator) ValidateArrayLength(field string, array []interface{}, min, max int) *Validator {
	length := len(array)
	
	if length < min {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   fmt.Sprintf("length:%d", length),
			Rule:    "min_length",
			Message: fmt.Sprintf("Must contain at least %d items", min),
		})
	}
	
	if max > 0 && length > max {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   fmt.Sprintf("length:%d", length),
			Rule:    "max_length",
			Message: fmt.Sprintf("Must contain no more than %d items", max),
		})
	}
	
	return v
}

// validateRule validates a single rule
func (v *Validator) validateRule(field string, value interface{}, rule string) *ValidationError {
	strValue := fmt.Sprintf("%v", value)
	
	// Parse rule and parameter
	parts := strings.Split(rule, ":")
	ruleName := parts[0]
	var param string
	if len(parts) > 1 {
		param = parts[1]
	}
	
	switch ruleName {
	case "required":
		if value == nil || strValue == "" || strValue == "0" {
			return &ValidationError{
				Field:   field,
				Value:   strValue,
				Rule:    ruleName,
				Message: "This field is required",
			}
		}
		
	case "min":
		if param != "" {
			minLen, _ := strconv.Atoi(param)
			if len(strValue) < minLen {
				return &ValidationError{
					Field:   field,
					Value:   strValue,
					Rule:    ruleName,
					Message: fmt.Sprintf("Must be at least %d characters", minLen),
				}
			}
		}
		
	case "max":
		if param != "" {
			maxLen, _ := strconv.Atoi(param)
			if len(strValue) > maxLen {
				return &ValidationError{
					Field:   field,
					Value:   strValue,
					Rule:    ruleName,
					Message: fmt.Sprintf("Must be no more than %d characters", maxLen),
				}
			}
		}
		
	case "numeric":
		if _, err := strconv.ParseFloat(strValue, 64); err != nil {
			return &ValidationError{
				Field:   field,
				Value:   strValue,
				Rule:    ruleName,
				Message: "Must be a number",
			}
		}
		
	case "alpha":
		alphaRegex := regexp.MustCompile(`^[a-zA-Z]+$`)
		if !alphaRegex.MatchString(strValue) {
			return &ValidationError{
				Field:   field,
				Value:   strValue,
				Rule:    ruleName,
				Message: "Must contain only letters",
			}
		}
		
	case "alphanum":
		alphanumRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
		if !alphanumRegex.MatchString(strValue) {
			return &ValidationError{
				Field:   field,
				Value:   strValue,
				Rule:    ruleName,
				Message: "Must contain only letters and numbers",
			}
		}
	}
	
	return nil
}

// HasErrors checks if there are any validation errors
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// Errors returns the validation errors
func (v *Validator) Errors() ValidationErrors {
	return v.errors
}

// ErrorsAsMap returns errors as a field->message map
func (v *Validator) ErrorsAsMap() map[string]string {
	errorMap := make(map[string]string)
	for _, err := range v.errors {
		errorMap[err.Field] = err.Message
	}
	return errorMap
}

// Clear clears all validation errors
func (v *Validator) Clear() {
	v.errors = make(ValidationErrors, 0)
}

// Sanitization functions

// SanitizeString removes potentially dangerous characters
func SanitizeString(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")
	
	// Trim whitespace
	input = strings.TrimSpace(input)
	
	return input
}

// SanitizeHTML removes HTML tags (basic sanitization)
func SanitizeHTML(input string) string {
	// Remove HTML tags
	htmlRegex := regexp.MustCompile(`<[^>]*>`)
	return htmlRegex.ReplaceAllString(input, "")
}

// SanitizeSQL escapes single quotes for basic SQL injection prevention
func SanitizeSQL(input string) string {
	return strings.ReplaceAll(input, "'", "''")
}

// NormalizeEmail normalizes email addresses
func NormalizeEmail(email string) string {
	email = strings.TrimSpace(strings.ToLower(email))
	return email
}

// NormalizePhone removes non-numeric characters from phone numbers
func NormalizePhone(phone string) string {
	phoneRegex := regexp.MustCompile(`[^\d+]`)
	return phoneRegex.ReplaceAllString(phone, "")
}

// Custom validation functions for specific business rules

// ValidateCampaignData validates campaign creation/update data
func ValidateCampaignData(data map[string]interface{}) ValidationErrors {
	v := NewValidator()
	
	// Required fields
	v.ValidateField("title", data["title"], "required", "min:3", "max:200")
	v.ValidateField("description", data["description"], "required", "min:10", "max:5000")
	v.ValidateField("goal_amount", data["goal_amount"], "required", "numeric")
	v.ValidateField("category", data["category"], "required")
	
	// Optional fields with validation
	if endDate, ok := data["end_date"].(string); ok && endDate != "" {
		v.ValidateDate("end_date", endDate)
	}
	
	if website, ok := data["website"].(string); ok && website != "" {
		v.ValidateURL("website", website)
	}
	
	// Business rule: goal amount should be reasonable
	if goalAmount, ok := data["goal_amount"].(float64); ok {
		v.ValidateAmount("goal_amount", goalAmount, 100, 10000000) // $100 to $10M
	}
	
	return v.Errors()
}

// ValidateUserRegistration validates user registration data
func ValidateUserRegistration(data map[string]interface{}) ValidationErrors {
	v := NewValidator()
	
	// Required fields
	v.ValidateField("name", data["name"], "required", "min:2", "max:100")
	v.ValidateField("email", data["email"], "required")
	v.ValidateField("password", data["password"], "required")
	v.ValidateField("role", data["role"], "required")
	
	// Email validation
	if email, ok := data["email"].(string); ok {
		v.ValidateEmail("email", email)
	}
	
	// Password validation
	if password, ok := data["password"].(string); ok {
		v.ValidatePassword("password", password)
	}
	
	// Role validation
	if role, ok := data["role"].(string); ok {
		allowedRoles := []string{"donor", "ngo_admin", "platform_admin"}
		v.ValidateEnum("role", role, allowedRoles)
	}
	
	return v.Errors()
}

// ValidateDonationData validates donation data
func ValidateDonationData(data map[string]interface{}) ValidationErrors {
	v := NewValidator()
	
	// Required fields
	v.ValidateField("campaign_id", data["campaign_id"], "required")
	v.ValidateField("amount", data["amount"], "required", "numeric")
	
	// UUID validation
	if campaignID, ok := data["campaign_id"].(string); ok {
		v.ValidateUUID("campaign_id", campaignID)
	}
	
	// Amount validation
	if amount, ok := data["amount"].(float64); ok {
		v.ValidateAmount("amount", amount, 1, 1000000) // $1 to $1M per donation
	}
	
	return v.Errors()
}