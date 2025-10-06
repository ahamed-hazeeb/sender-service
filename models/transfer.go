// DESIGN PATTERN: Data Transfer Object (DTO) + Entity Pattern
package models

import "time"

// Transfer - Entity representing a points transfer in the system
type Transfer struct {
	ID            string    `json:"id" gorm:"primaryKey"`                 // Primary key
	SenderID      string    `json:"sender_id" gorm:"not null;index"`      // Sender user ID with index
	SenderEmail   string    `json:"sender_email" gorm:"not null"`         // Sender's email
	ReceiverEmail string    `json:"receiver_email" gorm:"not null;index"` // Receiver email with index
	ReceiverName  string    `json:"receiver_name" gorm:"not null"`        // Receiver's name
	Points        int       `json:"points" gorm:"not null"`               // Points amount
	Status        string    `json:"status" gorm:"default:pending"`        // Transfer lifecycle: pending, completed, expired, cancelled
	Token         string    `json:"token" gorm:"uniqueIndex;not null"`    // Unique claim token
	ExpiresAt     time.Time `json:"expires_at" gorm:"not null"`           // Claim expiration time
	CreatedAt     time.Time `json:"created_at"`                           // Creation timestamp
	UpdatedAt     time.Time `json:"updated_at"`                           // Last update timestamp
}

// TransferRequest - DTO for transfer creation API input
type TransferRequest struct {
	ReceiverEmail string `json:"receiver_email" binding:"required,email"` // Must be valid email
	ReceiverName  string `json:"receiver_name" binding:"required,min=2"`  // Min 2 characters
	Points        int    `json:"points" binding:"required,min=1"`         // Must be positive
}

// User - External user model (from Auth Service) for service integration
type User struct {
	ID     string `json:"id"`     // User identifier
	Email  string `json:"email"`  // User email
	Name   string `json:"name"`   // User name
	Points int    `json:"points"` // Current points balance
}
