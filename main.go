// DESIGN PATTERN: Dependency Injection + Composition Root + Factory Pattern
package main

import (
	"fmt"
	"log"
	"sender-service/config"
	"sender-service/handlers"
	"sender-service/models"
	"sender-service/repositories"
	"sender-service/services"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// FACTORY PATTERN: Load configuration from environment
	cfg := config.LoadConfig()

	// 🗄️ DATABASE CONNECTION: Using GORM with PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// DATABASE MIGRATION: Auto-create transfer table
	db.AutoMigrate(&models.Transfer{})

	// DEPENDENCY INJECTION: Building the complete object graph
	// Repository Layer (Data Access)
	transferRepo := repositories.NewTransferRepository(db)

	// Service Layer (Business Logic + Email Integration)
	emailService := services.NewEmailService(cfg)
	transferService := services.NewTransferService(transferRepo, emailService, cfg)

	// Handler Layer (HTTP Interface)
	transferHandler := handlers.NewTransferHandler(transferService)

	// WEB SERVER CONFIGURATION
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode) // Optimized for production
	}

	r := gin.Default()

	// CORS MIDDLEWARE: Enable cross-origin requests
	setupCORS(r, cfg)

	// ROUTE SETUP: Define API endpoints for transfer operations
	setupRoutes(r, transferHandler)

	// START THE SENDER SERVICE
	log.Printf("Sender Service running on :%s in %s mode", cfg.Port, cfg.Environment)
	r.Run(":" + cfg.Port)
}

// setupCORS - Middleware for Cross-Origin Resource Sharing
func setupCORS(r *gin.Engine, cfg *config.Config) {
	r.Use(func(c *gin.Context) {
		// Set CORS headers to allow frontend communication
		c.Writer.Header().Set("Access-Control-Allow-Origin", cfg.Cors.AllowedOrigins)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight OPTIONS requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // No Content response
			return
		}
		c.Next()
	})
}

// setupRoutes - Router configuration (Front Controller Pattern)
func setupRoutes(r *gin.Engine, transferHandler *handlers.TransferHandler) {
	// TRANSFER MANAGEMENT ENDPOINTS
	r.POST("/transfer", transferHandler.InitiateTransfer)              // Create new transfer
	r.GET("/transfers/:userId", transferHandler.GetTransfers)          // Get user's transfer history
	r.POST("/transfer/:id/complete", transferHandler.CompleteTransfer) // Complete transfer (Saga step)
}
