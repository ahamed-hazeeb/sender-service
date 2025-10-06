// DESIGN PATTERN: Service Layer + Saga Pattern + Observer Pattern
package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sender-service/config"
	"sender-service/models"
	"sender-service/repositories"
	"time"
)

// TransferService - Orchestrates transfer business logic and coordinates with other services
type TransferService struct {
	transferRepo *repositories.TransferRepository // Composition: HAS-A repository
	emailService *EmailService                    // Composition: HAS-A email service
	config       *config.Config                   // Composition: HAS-A configuration
}

// NewTransferService - Factory method with dependency injection
func NewTransferService(transferRepo *repositories.TransferRepository,
	emailService *EmailService,
	config *config.Config) *TransferService {
	return &TransferService{
		transferRepo: transferRepo,
		emailService: emailService,
		config:       config,
	}
}

// InitiateTransfer - Business logic for creating a new points transfer
func (s *TransferService) InitiateTransfer(senderID string, req models.TransferRequest) (*models.Transfer, error) {
	// 1. SERVICE INTEGRATION: Get sender details from Auth Service
	sender, err := s.getUser(senderID)
	if err != nil {
		return nil, errors.New("failed to get sender details")
	}

	// 2. BUSINESS VALIDATION: Check transfer feasibility
	if err := s.validateTransfer(sender, req); err != nil {
		return nil, err
	}

	// 3. ENTITY CREATION: Create transfer record (points NOT deducted yet - Saga Pattern)
	transfer := &models.Transfer{
		ID:            generateID(),                   // Unique identifier
		SenderID:      senderID,                       // Sender user ID
		SenderEmail:   sender.Email,                   // Sender email
		ReceiverEmail: req.ReceiverEmail,              // Receiver email
		ReceiverName:  req.ReceiverName,               // Receiver name
		Points:        req.Points,                     // Points amount
		Status:        "pending",                      // Initial status
		Token:         generateToken(),                // Unique claim token
		ExpiresAt:     time.Now().Add(24 * time.Hour), // 24-hour expiration
		CreatedAt:     time.Now(),                     // Creation timestamp
		UpdatedAt:     time.Now(),                     // Update timestamp
	}

	// 4. PERSISTENCE: Save transfer to database
	if err := s.transferRepo.Create(transfer); err != nil {
		return nil, errors.New("failed to create transfer")
	}

	// 🎯 SAGA PATTERN: Points are NOT deducted here - only when receiver claims
	// This ensures points remain with sender if receiver doesn't claim

	// 5. OBSERVER PATTERN: Send email notification asynchronously
	go func() {
		if err := s.emailService.SendTransferEmail(transfer); err != nil {
			fmt.Printf("Failed to send email to %s: %v\n", transfer.ReceiverEmail, err)
		} else {
			fmt.Printf("Email sent successfully to: %s\n", transfer.ReceiverEmail)
		}
	}()

	return transfer, nil
}

// GetUserTransfers - Business logic to retrieve user's transfer history
func (s *TransferService) GetUserTransfers(userID string) ([]models.Transfer, error) {
	return s.transferRepo.FindBySenderID(userID)
}

// CompleteTransfer - SAGA PATTERN: Finalize transfer when receiver claims points
func (s *TransferService) CompleteTransfer(transferID string) error {
	transfer, err := s.transferRepo.FindByID(transferID)
	if err != nil {
		return errors.New("transfer not found")
	}

	// 1. SERVICE INTEGRATION: Get current sender details
	sender, err := s.getUser(transfer.SenderID)
	if err != nil {
		return errors.New("failed to get sender details")
	}

	// 2. VALIDATION: Ensure sender still has sufficient points
	if sender.Points < transfer.Points {
		// Mark transfer as failed due to insufficient points
		transfer.Status = "failed"
		s.transferRepo.Update(transfer)
		return errors.New("sender no longer has sufficient points")
	}

	// 3. POINT DEDUCTION: Deduct points from sender (Saga commitment)
	if err := s.updateUserPoints(transfer.SenderID, sender.Points-transfer.Points); err != nil {
		return errors.New("failed to deduct points from sender")
	}

	// 4. STATUS UPDATE: Mark transfer as completed
	transfer.Status = "completed"
	if err := s.transferRepo.Update(transfer); err != nil {
		// ⚠️ SAGA COMPENSATION: Points deducted but transfer not completed
		// In production, implement compensation logic here
		return errors.New("failed to complete transfer")
	}

	return nil
}

// validateTransfer - Business rules validation
func (s *TransferService) validateTransfer(sender *models.User, req models.TransferRequest) error {
	// Business Rule 1: Sufficient points
	if sender.Points < req.Points {
		return errors.New("insufficient points")
	}

	// Business Rule 2: Cannot transfer to self
	if sender.Email == req.ReceiverEmail {
		return errors.New("cannot transfer points to yourself")
	}

	// Business Rule 3: Positive points amount
	if req.Points <= 0 {
		return errors.New("points must be greater than zero")
	}

	return nil
}

// getUser - Service-to-service call to Auth Service
func (s *TransferService) getUser(userID string) (*models.User, error) {
	resp, err := http.Get(s.config.AuthService + "/users/" + userID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("user not found")
	}

	var response struct {
		Success bool         `json:"success"`
		Data    *models.User `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil || !response.Success {
		return nil, errors.New("failed to get user data")
	}

	return response.Data, nil
}

// updateUserPoints - Service-to-service call to update user points
func (s *TransferService) updateUserPoints(userID string, points int) error {
	requestBody := map[string]int{"points": points}
	jsonData, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("PUT", s.config.AuthService+"/users/"+userID+"/points",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to update points")
	}

	return nil
}

// generateID - Utility function for unique ID generation
func generateID() string {
	return fmt.Sprintf("transfer_%d", time.Now().UnixNano())
}

// generateToken - Utility function for unique token generation
func generateToken() string {
	return fmt.Sprintf("token_%d", time.Now().UnixNano())
}
