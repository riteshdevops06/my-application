package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// ProductionConfig holds production-specific configuration
type ProductionConfig struct {
	*Config
	
	// Security settings
	Security SecurityConfig `json:"security"`
	
	// Database settings
	Database DatabaseConfig `json:"database"`
	
	// Redis settings
	Redis RedisConfig `json:"redis"`
	
	// Logging settings
	Logging LoggingConfig `json:"logging"`
	
	// Monitoring settings
	Monitoring MonitoringConfig `json:"monitoring"`
	
	// Email settings
	Email EmailConfig `json:"email"`
	
	// File upload settings
	Upload UploadConfig `json:"upload"`
	
	// Rate limiting settings
	RateLimit RateLimitConfig `json:"rate_limit"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	JWTSecret           string        `json:"jwt_secret"`
	JWTExpiration      time.Duration `json:"jwt_expiration"`
	RefreshExpiration  time.Duration `json:"refresh_expiration"`
	BCryptCost         int           `json:"bcrypt_cost"`
	MaxLoginAttempts   int           `json:"max_login_attempts"`
	LoginLockoutTime   time.Duration `json:"login_lockout_time"`
	PasswordMinLength  int           `json:"password_min_length"`
	RequireHTTPS       bool          `json:"require_https"`
	TrustedProxies     []string      `json:"trusted_proxies"`
	AllowedOrigins     []string      `json:"allowed_origins"`
	APIKeys            map[string]string `json:"api_keys"`
	EncryptionKey      string        `json:"encryption_key"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string        `json:"host"`
	Port            int           `json:"port"`
	User            string        `json:"user"`
	Password        string        `json:"password"`
	Database        string        `json:"database"`
	SSLMode         string        `json:"ssl_mode"`
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Password    string `json:"password"`
	Database    int    `json:"database"`
	MaxRetries  int    `json:"max_retries"`
	PoolSize    int    `json:"pool_size"`
	DialTimeout time.Duration `json:"dial_timeout"`
	ReadTimeout time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level          string `json:"level"`
	Format         string `json:"format"`
	Output         string `json:"output"`
	File           string `json:"file"`
	MaxSize        int    `json:"max_size"`
	MaxAge         int    `json:"max_age"`
	MaxBackups     int    `json:"max_backups"`
	Compress       bool   `json:"compress"`
	EnableConsole  bool   `json:"enable_console"`
	EnableFile     bool   `json:"enable_file"`
	EnableSyslog   bool   `json:"enable_syslog"`
	SyslogNetwork  string `json:"syslog_network"`
	SyslogAddress  string `json:"syslog_address"`
}

// MonitoringConfig holds monitoring configuration
type MonitoringConfig struct {
	EnableMetrics    bool   `json:"enable_metrics"`
	MetricsPort     int    `json:"metrics_port"`
	EnablePprof     bool   `json:"enable_pprof"`
	PprofPort       int    `json:"pprof_port"`
	HealthCheckPath string `json:"health_check_path"`
	EnableTracing   bool   `json:"enable_tracing"`
	TracingEndpoint string `json:"tracing_endpoint"`
	ServiceName     string `json:"service_name"`
	ServiceVersion  string `json:"service_version"`
}

// EmailConfig holds email service configuration
type EmailConfig struct {
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUser     string `json:"smtp_user"`
	SMTPPassword string `json:"smtp_password"`
	FromAddress  string `json:"from_address"`
	FromName     string `json:"from_name"`
	UseTLS       bool   `json:"use_tls"`
	UseStartTLS  bool   `json:"use_starttls"`
}

// UploadConfig holds file upload configuration
type UploadConfig struct {
	MaxFileSize      int64    `json:"max_file_size"`
	AllowedMimeTypes []string `json:"allowed_mime_types"`
	UploadPath       string   `json:"upload_path"`
	TempPath         string   `json:"temp_path"`
	EnableVirusScan  bool     `json:"enable_virus_scan"`
	VirusScanAPI     string   `json:"virus_scan_api"`
	S3Bucket         string   `json:"s3_bucket"`
	S3Region         string   `json:"s3_region"`
	S3AccessKey      string   `json:"s3_access_key"`
	S3SecretKey      string   `json:"s3_secret_key"`
	UseS3            bool     `json:"use_s3"`
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	GlobalLimit    int           `json:"global_limit"`
	GlobalWindow   time.Duration `json:"global_window"`
	AuthLimit      int           `json:"auth_limit"`
	AuthWindow     time.Duration `json:"auth_window"`
	APILimit       int           `json:"api_limit"`
	APIWindow      time.Duration `json:"api_window"`
	UploadLimit    int           `json:"upload_limit"`
	UploadWindow   time.Duration `json:"upload_window"`
	EnableRedis    bool          `json:"enable_redis"`
}

// LoadProductionConfig loads configuration for production environment
func LoadProductionConfig() (*ProductionConfig, error) {
	baseConfig, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	
	prodConfig := &ProductionConfig{
		Config: baseConfig,
	}
	
	// Load security configuration
	prodConfig.Security = SecurityConfig{
		JWTSecret:         getEnvString("JWT_SECRET", generateSecretKey()),
		JWTExpiration:     getEnvDuration("JWT_EXPIRATION", 24*time.Hour),
		RefreshExpiration: getEnvDuration("REFRESH_EXPIRATION", 7*24*time.Hour),
		BCryptCost:        getEnvInt("BCRYPT_COST", 12),
		MaxLoginAttempts:  getEnvInt("MAX_LOGIN_ATTEMPTS", 5),
		LoginLockoutTime:  getEnvDuration("LOGIN_LOCKOUT_TIME", 15*time.Minute),
		PasswordMinLength: getEnvInt("PASSWORD_MIN_LENGTH", 8),
		RequireHTTPS:      getEnvBool("REQUIRE_HTTPS", true),
		TrustedProxies:    getEnvStringSlice("TRUSTED_PROXIES", []string{}),
		AllowedOrigins:    getEnvStringSlice("ALLOWED_ORIGINS", []string{}),
		EncryptionKey:     getEnvString("ENCRYPTION_KEY", generateEncryptionKey()),
	}
	
	// Load database configuration
	prodConfig.Database = DatabaseConfig{
		Host:            getEnvString("DB_HOST", "localhost"),
		Port:            getEnvInt("DB_PORT", 5432),
		User:            getEnvString("DB_USER", "postgres"),
		Password:        getEnvString("DB_PASSWORD", ""),
		Database:        getEnvString("DB_NAME", "ngo_tracker"),
		SSLMode:         getEnvString("DB_SSL_MODE", "require"),
		MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
		ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", time.Hour),
		ConnMaxIdleTime: getEnvDuration("DB_CONN_MAX_IDLE_TIME", 30*time.Minute),
	}
	
	// Load Redis configuration
	prodConfig.Redis = RedisConfig{
		Host:         getEnvString("REDIS_HOST", "localhost"),
		Port:         getEnvInt("REDIS_PORT", 6379),
		Password:     getEnvString("REDIS_PASSWORD", ""),
		Database:     getEnvInt("REDIS_DB", 0),
		MaxRetries:   getEnvInt("REDIS_MAX_RETRIES", 3),
		PoolSize:     getEnvInt("REDIS_POOL_SIZE", 10),
		DialTimeout:  getEnvDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
		ReadTimeout:  getEnvDuration("REDIS_READ_TIMEOUT", 3*time.Second),
		WriteTimeout: getEnvDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),
	}
	
	// Load logging configuration
	prodConfig.Logging = LoggingConfig{
		Level:         getEnvString("LOG_LEVEL", "info"),
		Format:        getEnvString("LOG_FORMAT", "json"),
		Output:        getEnvString("LOG_OUTPUT", "file"),
		File:          getEnvString("LOG_FILE", "/var/log/ngo-tracker/app.log"),
		MaxSize:       getEnvInt("LOG_MAX_SIZE", 100),
		MaxAge:        getEnvInt("LOG_MAX_AGE", 30),
		MaxBackups:    getEnvInt("LOG_MAX_BACKUPS", 10),
		Compress:      getEnvBool("LOG_COMPRESS", true),
		EnableConsole: getEnvBool("LOG_ENABLE_CONSOLE", false),
		EnableFile:    getEnvBool("LOG_ENABLE_FILE", true),
		EnableSyslog:  getEnvBool("LOG_ENABLE_SYSLOG", false),
		SyslogNetwork: getEnvString("LOG_SYSLOG_NETWORK", "udp"),
		SyslogAddress: getEnvString("LOG_SYSLOG_ADDRESS", "localhost:514"),
	}
	
	// Load monitoring configuration
	prodConfig.Monitoring = MonitoringConfig{
		EnableMetrics:    getEnvBool("ENABLE_METRICS", true),
		MetricsPort:     getEnvInt("METRICS_PORT", 9090),
		EnablePprof:     getEnvBool("ENABLE_PPROF", false),
		PprofPort:       getEnvInt("PPROF_PORT", 6060),
		HealthCheckPath: getEnvString("HEALTH_CHECK_PATH", "/health"),
		EnableTracing:   getEnvBool("ENABLE_TRACING", true),
		TracingEndpoint: getEnvString("TRACING_ENDPOINT", ""),
		ServiceName:     getEnvString("SERVICE_NAME", "ngo-donation-tracker"),
		ServiceVersion:  getEnvString("SERVICE_VERSION", "1.0.0"),
	}
	
	// Load email configuration
	prodConfig.Email = EmailConfig{
		SMTPHost:     getEnvString("SMTP_HOST", ""),
		SMTPPort:     getEnvInt("SMTP_PORT", 587),
		SMTPUser:     getEnvString("SMTP_USER", ""),
		SMTPPassword: getEnvString("SMTP_PASSWORD", ""),
		FromAddress:  getEnvString("FROM_ADDRESS", "noreply@ngotracker.com"),
		FromName:     getEnvString("FROM_NAME", "NGO Tracker"),
		UseTLS:       getEnvBool("SMTP_USE_TLS", true),
		UseStartTLS:  getEnvBool("SMTP_USE_STARTTLS", true),
	}
	
	// Load upload configuration
	prodConfig.Upload = UploadConfig{
		MaxFileSize:      getEnvInt64("UPLOAD_MAX_FILE_SIZE", 10<<20), // 10MB
		AllowedMimeTypes: getEnvStringSlice("UPLOAD_ALLOWED_TYPES", []string{
			"image/jpeg", "image/png", "image/gif", "image/webp",
			"application/pdf", "application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		}),
		UploadPath:       getEnvString("UPLOAD_PATH", "/var/uploads"),
		TempPath:         getEnvString("TEMP_PATH", "/tmp"),
		EnableVirusScan:  getEnvBool("ENABLE_VIRUS_SCAN", true),
		VirusScanAPI:     getEnvString("VIRUS_SCAN_API", ""),
		S3Bucket:         getEnvString("S3_BUCKET", ""),
		S3Region:         getEnvString("S3_REGION", "us-east-1"),
		S3AccessKey:      getEnvString("S3_ACCESS_KEY", ""),
		S3SecretKey:      getEnvString("S3_SECRET_KEY", ""),
		UseS3:            getEnvBool("USE_S3", false),
	}
	
	// Load rate limiting configuration
	prodConfig.RateLimit = RateLimitConfig{
		GlobalLimit:  getEnvInt("RATE_LIMIT_GLOBAL", 1000),
		GlobalWindow: getEnvDuration("RATE_LIMIT_GLOBAL_WINDOW", time.Hour),
		AuthLimit:    getEnvInt("RATE_LIMIT_AUTH", 10),
		AuthWindow:   getEnvDuration("RATE_LIMIT_AUTH_WINDOW", time.Minute),
		APILimit:     getEnvInt("RATE_LIMIT_API", 100),
		APIWindow:    getEnvDuration("RATE_LIMIT_API_WINDOW", time.Minute),
		UploadLimit:  getEnvInt("RATE_LIMIT_UPLOAD", 5),
		UploadWindow: getEnvDuration("RATE_LIMIT_UPLOAD_WINDOW", time.Minute),
		EnableRedis:  getEnvBool("RATE_LIMIT_ENABLE_REDIS", true),
	}
	
	return prodConfig, nil
}

// Utility functions for environment variable parsing
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvStringSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

// generateSecretKey generates a secure random key for JWT signing
func generateSecretKey() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(fmt.Sprintf("Failed to generate secret key: %v", err))
	}
	return hex.EncodeToString(bytes)
}

// generateEncryptionKey generates a secure random key for data encryption
func generateEncryptionKey() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(fmt.Sprintf("Failed to generate encryption key: %v", err))
	}
	return hex.EncodeToString(bytes)
}

// Validate validates the production configuration
func (pc *ProductionConfig) Validate() error {
	// Validate required security settings
	if pc.Security.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	
	if len(pc.Security.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}
	
	// Validate database settings
	if pc.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	
	if pc.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	
	// Validate email settings if email features are enabled
	if pc.Email.SMTPHost != "" {
		if pc.Email.SMTPUser == "" || pc.Email.SMTPPassword == "" {
			return fmt.Errorf("SMTP credentials are required when SMTP_HOST is set")
		}
	}
	
	// Validate S3 settings if S3 is enabled
	if pc.Upload.UseS3 {
		if pc.Upload.S3Bucket == "" || pc.Upload.S3AccessKey == "" || pc.Upload.S3SecretKey == "" {
			return fmt.Errorf("S3 credentials are required when USE_S3 is enabled")
		}
	}
	
	return nil
}

// GetDatabaseURL returns the formatted database connection URL
func (pc *ProductionConfig) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		pc.Database.User,
		pc.Database.Password,
		pc.Database.Host,
		pc.Database.Port,
		pc.Database.Database,
		pc.Database.SSLMode,
	)
}

// GetRedisURL returns the formatted Redis connection URL
func (pc *ProductionConfig) GetRedisURL() string {
	if pc.Redis.Password != "" {
		return fmt.Sprintf("redis://:%s@%s:%d/%d",
			pc.Redis.Password,
			pc.Redis.Host,
			pc.Redis.Port,
			pc.Redis.Database,
		)
	}
	return fmt.Sprintf("redis://%s:%d/%d",
		pc.Redis.Host,
		pc.Redis.Port,
		pc.Redis.Database,
	)
}