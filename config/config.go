// DESIGN PATTERN: Factory Pattern + Service Integration Configuration
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config - Centralized configuration container for sender service
type Config struct {
	Port        string         // Service port (8002)
	Environment string         // Runtime environment
	Database    DatabaseConfig // Database configuration
	AuthService string         // URL for Auth Service (Service Integration)
	Email       EmailConfig    // Email service configuration (Strategy Pattern)
	Frontend    FrontendConfig // Frontend application configuration
	Cors        CorsConfig     // CORS settings
}

// DatabaseConfig - Encapsulates database connection details
type DatabaseConfig struct {
	Host     string // Database host address
	Port     string // Database port
	Name     string // Database name
	User     string // Database username
	Password string // Database password
	SSLMode  string // SSL mode for secure connection
}

// EmailConfig - Encapsulates email service configuration (Strategy Pattern)
type EmailConfig struct {
	GmailAddress string // Gmail account for sending emails
	GmailAppPass string // Gmail app password
	From         string // Sender email address
	SMTPHost     string // SMTP server host
	SMTPPort     string // SMTP server port
}

// FrontendConfig - Encapsulates frontend application settings
type FrontendConfig struct {
	URL string // Frontend application URL for claim links
}

// CorsConfig - Encapsulates CORS policy settings
type CorsConfig struct {
	AllowedOrigins string // Allowed frontend domains
}

// LoadConfig - Factory method that creates configured Config instance
func LoadConfig() *Config {
	// Load environment variables with fallback to OS environment
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Factory construction with sensible defaults
	return &Config{
		Port:        getEnv("PORT", "8002"), // Sender service default port
		Environment: getEnv("ENVIRONMENT", "development"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "point_transfer"),
			User:     getEnv("DB_USER", "point_user"),
			Password: getEnv("DB_PASSWORD", "password123"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		AuthService: getEnv("AUTH_SERVICE_URL", "http://localhost:8001"), // Service integration
		Email: EmailConfig{
			GmailAddress: getEnv("GMAIL_ADDRESS", ""),      // Email strategy configuration
			GmailAppPass: getEnv("GMAIL_APP_PASSWORD", ""), // Email strategy configuration
			From:         getEnv("EMAIL_FROM", "noreply@pointtransfer.com"),
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"), // Default to Gmail
			SMTPPort:     getEnv("SMTP_PORT", "587"),            // Default TLS port
		},
		Frontend: FrontendConfig{
			URL: getEnv("FRONTEND_URL", "http://localhost:3000"), // Frontend URL for claim links
		},
		Cors: CorsConfig{
			AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
		},
	}
}

// getEnv - Helper with fallback values (Null Object Pattern)
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
