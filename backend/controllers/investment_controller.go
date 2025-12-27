package controllers

import (
	"investment-tracker-backend/config"
	"investment-tracker-backend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetInvestments retrieves all investments for the authenticated user
func GetInvestments(c *gin.Context) {
	// Get user_id from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	var investments []models.Investment
	if err := config.DB.Where("user_id = ?", uint(userID)).Find(&investments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, investments)
}

// GetInvestment retrieves a single investment by ID for the authenticated user
func GetInvestment(c *gin.Context) {
	// Get user_id from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	id := c.Param("id")
	investmentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var investment models.Investment
	if err := config.DB.Where("id = ? AND user_id = ?", uint(investmentID), uint(userID)).First(&investment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment not found"})
		return
	}

	c.JSON(http.StatusOK, investment)
}

// CreateInvestment creates a new investment
func CreateInvestment(c *gin.Context) {
	var investment models.Investment
	if err := c.BindJSON(&investment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data: " + err.Error(),
			"hint":  "Make sure all required fields are provided: name, type, invested, current_value",
		})
		return
	}

	// Get user_id from context (set by AuthMiddleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}
	investment.UserID = uint(userID)

	// Manual validation
	if investment.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Investment name is required"})
		return
	}
	if investment.Type == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Investment type is required"})
		return
	}
	if investment.Invested <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invested amount must be greater than 0"})
		return
	}

	// Set default purchase date if not provided
	if investment.PurchaseDate.IsZero() {
		investment.PurchaseDate = time.Now()
	}

	// Calculate returns and status
	investment.CalculateReturns()
	investment.UpdateStatus()

	if err := config.DB.Create(&investment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create investment: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, investment)
}

// UpdateInvestment updates an existing investment
func UpdateInvestment(c *gin.Context) {
	// Get user_id from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	id := c.Param("id")
	investmentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Get the old investment to check if goal_id changed and verify ownership
	var oldInvestment models.Investment
	if err := config.DB.Where("id = ? AND user_id = ?", uint(investmentID), uint(userID)).First(&oldInvestment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment not found"})
		return
	}

	var investment models.Investment
	if err := c.ShouldBindJSON(&investment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	investment.ID = uint(investmentID)

	// Recalculate returns and status
	investment.CalculateReturns()
	investment.UpdateStatus()

	if err := config.DB.Save(&investment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update goal current_amount if goal_id changed or current_value changed
	if oldInvestment.GoalID != nil {
		updateGoalCurrentAmount(*oldInvestment.GoalID)
	}
	if investment.GoalID != nil && (oldInvestment.GoalID == nil || *oldInvestment.GoalID != *investment.GoalID) {
		updateGoalCurrentAmount(*investment.GoalID)
	}

	c.JSON(http.StatusOK, investment)
}

// DeleteInvestment deletes an investment
func DeleteInvestment(c *gin.Context) {
	// Get user_id from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	id := c.Param("id")
	investmentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Get the investment before deleting to check if it's linked to a goal and verify ownership
	var investment models.Investment
	if err := config.DB.Where("id = ? AND user_id = ?", uint(investmentID), uint(userID)).First(&investment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment not found"})
		return
	}

	// Store the goal_id before deletion
	goalID := investment.GoalID

	if err := config.DB.Delete(&investment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update the goal's current_amount if the investment was linked to a goal
	if goalID != nil {
		updateGoalCurrentAmount(*goalID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Investment deleted successfully"})
}

// LinkInvestmentToGoal links an investment to a goal
func LinkInvestmentToGoal(c *gin.Context) {
	// Get user_id from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	investmentID := c.Param("id")
	investmentIDUint, err := strconv.ParseUint(investmentID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid investment ID format"})
		return
	}

	var requestBody struct {
		GoalID string `json:"goal_id"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Get the investment and verify ownership
	var investment models.Investment
	if err := config.DB.Where("id = ? AND user_id = ?", uint(investmentIDUint), uint(userID)).First(&investment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment not found"})
		return
	}

	var goalIDPtr *uint
	var newGoalID uint
	if requestBody.GoalID != "" {
		gid, err := strconv.ParseUint(requestBody.GoalID, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID format"})
			return
		}
		goalIDUint := uint(gid)
		goalIDPtr = &goalIDUint
		newGoalID = goalIDUint

		// Verify goal exists and belongs to the user
		var goal models.Goal
		if err := config.DB.Where("id = ? AND user_id = ?", goalIDUint, uint(userID)).First(&goal).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
			return
		}
	}

	// Update investment with goal_id
	investment.GoalID = goalIDPtr

	if err := config.DB.Save(&investment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update the old goal's current_amount (if there was one)
	if investment.GoalID != nil {
		updateGoalCurrentAmount(*investment.GoalID)
	}

	// Update the new goal's current_amount (if linking to a goal)
	if newGoalID != 0 {
		updateGoalCurrentAmount(newGoalID)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Investment linked to goal successfully",
		"goal_id": requestBody.GoalID,
	})
}

// updateGoalCurrentAmount recalculates and updates a goal's current_amount based on linked investments
func updateGoalCurrentAmount(goalID uint) error {
	// Find all investments linked to this goal
	var investments []models.Investment
	if err := config.DB.Where("goal_id = ?", goalID).Find(&investments).Error; err != nil {
		return err
	}

	// Calculate total current value
	totalCurrentValue := 0.0
	for _, inv := range investments {
		totalCurrentValue += inv.CurrentValue
	}

	// Update the goal's current_amount
	return config.DB.Model(&models.Goal{}).Where("id = ?", goalID).Updates(map[string]interface{}{
		"current_amount": totalCurrentValue,
		"updated_at":     time.Now(),
	}).Error
}

// UnlinkInvestmentFromGoal removes the goal link from an investment
func UnlinkInvestmentFromGoal(c *gin.Context) {
	// Get user_id from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	investmentID := c.Param("id")
	investmentIDUint, err := strconv.ParseUint(investmentID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid investment ID format"})
		return
	}

	// Get the investment and verify ownership
	var investment models.Investment
	if err := config.DB.Where("id = ? AND user_id = ?", uint(investmentIDUint), uint(userID)).First(&investment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment not found"})
		return
	}

	// Store the old goal_id before unlinking
	oldGoalID := investment.GoalID

	// Update investment to remove goal_id
	investment.GoalID = nil

	if err := config.DB.Save(&investment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update the old goal's current_amount
	if oldGoalID != nil {
		updateGoalCurrentAmount(*oldGoalID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Investment unlinked from goal successfully"})
}

// GetInvestmentsByGoal retrieves all investments linked to a specific goal
func GetInvestmentsByGoal(c *gin.Context) {
	// Get user_id from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	goalID := c.Param("goal_id")
	goalIDUint, err := strconv.ParseUint(goalID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID format"})
		return
	}

	// Verify the goal belongs to the user
	var goal models.Goal
	if err := config.DB.Where("id = ? AND user_id = ?", uint(goalIDUint), uint(userID)).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	// Get investments for this goal (they should all belong to the user since the goal does)
	var investments []models.Investment
	if err := config.DB.Where("goal_id = ? AND user_id = ?", uint(goalIDUint), uint(userID)).Find(&investments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate total
	total := 0.0
	for _, inv := range investments {
		total += inv.CurrentValue
	}

	c.JSON(http.StatusOK, gin.H{
		"investments": investments,
		"total":       total,
		"count":       len(investments),
	})
}

//
