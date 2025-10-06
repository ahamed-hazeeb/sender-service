// DESIGN PATTERN: Controller Pattern + Request Handler
package handlers

import (
	"net/http"
	"sender-service/models"
	"sender-service/services"

	"github.com/gin-gonic/gin"
)

// TransferHandler - Handles HTTP requests for transfer operations
type TransferHandler struct {
	transferService *services.TransferService // Composition: HAS-A business service
}

// NewTransferHandler - Factory method with dependency injection
func NewTransferHandler(transferService *services.TransferService) *TransferHandler {
	return &TransferHandler{transferService: transferService}
}

// InitiateTransfer - HTTP handler to create a new points transfer
func (h *TransferHandler) InitiateTransfer(c *gin.Context) {
	var req models.TransferRequest

	// 1. REQUEST VALIDATION: Parse and validate JSON input
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request data",
			"details": err.Error(), // Development details
		})
		return
	}

	// 2. AUTHENTICATION: Extract user ID from header (simplified JWT)
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User authentication required",
		})
		return
	}

	// 3. BUSINESS LOGIC: Delegate to service layer
	transfer, err := h.transferService.InitiateTransfer(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(), // Business error
		})
		return
	}

	// 4. SUCCESS RESPONSE
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Transfer initiated successfully",
		"data":    transfer,
	})
}

// GetTransfers - HTTP handler to get user's transfer history
func (h *TransferHandler) GetTransfers(c *gin.Context) {
	userID := c.Param("userId") // Extract user ID from URL path

	transfers, err := h.transferService.GetUserTransfers(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch transfers",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    transfers,
	})
}

// CompleteTransfer - HTTP handler for completing transfer (Saga Pattern step)
func (h *TransferHandler) CompleteTransfer(c *gin.Context) {
	transferID := c.Param("id") // Extract transfer ID from URL path

	// Delegate to service layer for business logic
	err := h.transferService.CompleteTransfer(transferID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Transfer completed successfully",
	})
}
