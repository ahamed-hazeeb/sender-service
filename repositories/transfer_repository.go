// DESIGN PATTERN: Repository Pattern + CRUD Operations
package repositories

import (
	"sender-service/models"

	"gorm.io/gorm"
)

// TransferRepository - Abstracts all database operations for Transfer entity
type TransferRepository struct {
	db *gorm.DB // Composition: HAS-A database connection
}

// NewTransferRepository - Factory method for repository
func NewTransferRepository(db *gorm.DB) *TransferRepository {
	return &TransferRepository{db: db}
}

// Create - Persists new transfer to database
func (r *TransferRepository) Create(transfer *models.Transfer) error {
	// GORM: INSERT INTO transfers (...) VALUES (...)
	return r.db.Create(transfer).Error
}

// FindBySenderID - Finds all transfers for a specific sender
func (r *TransferRepository) FindBySenderID(senderID string) ([]models.Transfer, error) {
	var transfers []models.Transfer
	// GORM: SELECT * FROM transfers WHERE sender_id = ? ORDER BY created_at DESC
	err := r.db.Where("sender_id = ?", senderID).
		Order("created_at DESC").
		Find(&transfers).Error
	return transfers, err
}

// FindByToken - Finds transfer by unique claim token
func (r *TransferRepository) FindByToken(token string) (*models.Transfer, error) {
	var transfer models.Transfer
	// GORM: SELECT * FROM transfers WHERE token = ? LIMIT 1
	err := r.db.Where("token = ?", token).First(&transfer).Error
	return &transfer, err
}

// Update - Updates transfer entity in database
func (r *TransferRepository) Update(transfer *models.Transfer) error {
	// GORM: UPDATE transfers SET ... WHERE id = ?
	return r.db.Save(transfer).Error
}

// Delete - Removes transfer from database (for rollback scenarios)
func (r *TransferRepository) Delete(transfer *models.Transfer) error {
	// GORM: DELETE FROM transfers WHERE id = ?
	return r.db.Delete(transfer).Error
}

// FindByID - Finds transfer by unique identifier (for Saga completion)
func (r *TransferRepository) FindByID(transferID string) (*models.Transfer, error) {
	var transfer models.Transfer
	// GORM: SELECT * FROM transfers WHERE id = ? LIMIT 1
	err := r.db.Where("id = ?", transferID).First(&transfer).Error
	return &transfer, err
}
