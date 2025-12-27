package controllers

import (
	"investment-tracker-backend/config"
	"investment-tracker-backend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetGoals retrieves all goals for the authenticated user
func GetGoals(c *gin.Context) {
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

	var goals []models.Goal
	if err := config.DB.Where("user_id = ?", uint(userID)).Find(&goals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goals)
}

// GetGoal retrieves a single goal by ID for the authenticated user
func GetGoal(c *gin.Context) {
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
	goalID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var goal models.Goal
	if err := config.DB.Where("id = ? AND user_id = ?", uint(goalID), uint(userID)).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	c.JSON(http.StatusOK, goal)
}

// CreateGoal creates a new goal
func CreateGoal(c *gin.Context) {
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

	var goal models.Goal
	if err := c.ShouldBindJSON(&goal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	goal.UserID = uint(userID)

	// Manual validation
	if goal.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Goal name is required"})
		return
	}
	if goal.TargetAmount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Target amount must be greater than 0"})
		return
	}

	// Set default deadline if not provided
	if goal.Deadline.IsZero() {
		goal.Deadline = time.Now().AddDate(1, 0, 0) // 1 year from now
	}

	// Set default status if not provided
	if goal.Status == "" {
		goal.Status = "Planned"
	}

	// Update status based on progress
	goal.UpdateStatus()

	if err := config.DB.Create(&goal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create goal: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, goal)
}

// UpdateGoal updates an existing goal
func UpdateGoal(c *gin.Context) {
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
	goalID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Verify ownership
	var existingGoal models.Goal
	if err := config.DB.Where("id = ? AND user_id = ?", uint(goalID), uint(userID)).First(&existingGoal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	var goal models.Goal
	if err := c.ShouldBindJSON(&goal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal.ID = uint(goalID)
	goal.UserID = uint(userID)

	// Update status based on progress
	goal.UpdateStatus()

	if err := config.DB.Save(&goal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goal)
}

// DeleteGoal deletes a goal
func DeleteGoal(c *gin.Context) {
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
	goalID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Verify ownership before deleting
	var goal models.Goal
	if err := config.DB.Where("id = ? AND user_id = ?", uint(goalID), uint(userID)).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	if err := config.DB.Delete(&goal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Goal deleted successfully"})
}

//
